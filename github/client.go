package github

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
	git "gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func init() {
	home, _ := os.UserHomeDir()
	file, err := os.Open(path.Join(home, ".github.txt"))
	if err != nil {
		panic("Can't find the token file")
	}
	res, err := ioutil.ReadAll(file)
	if err != nil {
		panic("Can't read the token file")
	}
	token := strings.TrimRight(string(res), "\n")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	Client = github.NewClient(tc)
}

var Client *github.Client

func getGitHead() (string, error) {

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	repo, err := git.PlainOpen(wd)
	if err != nil {
		return "", err
	}
	head ,err := repo.Head()
	if err != nil {
		return "", err
	}

	return head.Hash().String(), nil
}

func getCurrentPRNum(org, repo string) (int, error){
	head, err := getGitHead()
	if err != nil {
		return -1, err
	}

	res, _, err := Client.PullRequests.ListPullRequestsWithCommit(context.TODO(), org, repo, head, nil)
	if err != nil {
		return -1, err
	}

	for _, pr := range res {
		fmt.Println("Got pull request", *pr.Number)
		if *pr.Head.SHA == head {
			return *pr.Number, nil
		}
	}

	return -1, errors.New("commit not found")
}

func AddComment(org, repo, msg string) (string, error) {
	prNum, err := getCurrentPRNum(org, repo)
	if err != nil {
		return "", err
	}

	cmnt := &github.IssueComment{Body: &msg}
	res, _, err := Client.Issues.CreateComment(context.TODO(), org, repo, prNum, cmnt)

	if err != nil {
		return "", err
	}

	return *res.URL, nil
}
