package github

import (
	"context"
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

func GetGitHead() (string, error) {

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

func GetCurrentPullReq() (int, error){
	head, err := GetGitHead()
	if err != nil {
		return -1, err
	}

	res, _, err := Client.PullRequests.ListPullRequestsWithCommit(context.TODO(), "nunnatsa", "test-comments", head, nil)
	if err != nil {
		return -1, err
	}

	for _, pr := range res {
		fmt.Println("Got pull request")
		if *pr.Head.SHA == head {
			return *pr.Number, nil
		}
	}

	return -1, nil



}

type githubUser struct {
	Login string `json:"login,omitempty"`
}