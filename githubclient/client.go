package githubclient

import (
	"context"
	"errors"
	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
	git "gopkg.in/src-d/go-git.v4"
	"os"
)

func NewClient() (Client, error) {
	token := os.Getenv("GITHUB_API_TOKEN")
	if token == "" {
		return nil, errors.New("GITHUB_API_TOKEN environment variable is not set")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &client{github.NewClient(tc)}, nil
}

type Client interface {
	AddComment(org, repo, msg string) (string, error)
}

type client struct {
	client *github.Client
}

func (c client) getGitHead() (string, error) {

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

func (c client) getCurrentPRNum(org, repo string) (int, error){
	head, err := c.getGitHead()
	if err != nil {
		return -1, err
	}

	res, _, err := c.client.PullRequests.ListPullRequestsWithCommit(context.TODO(), org, repo, head, nil)
	if err != nil {
		return -1, err
	}

	for _, pr := range res {
		if *pr.Head.SHA == head {
			return *pr.Number, nil
		}
	}

	return -1, errors.New("commit not found")
}

func (c client) AddComment(org, repo, msg string) (string, error) {
	prNum, err := c.getCurrentPRNum(org, repo)
	if err != nil {
		return "", err
	}

	cmnt := &github.IssueComment{Body: &msg}
	res, _, err := c.client.Issues.CreateComment(context.TODO(), org, repo, prNum, cmnt)

	if err != nil {
		return "", err
	}

	return *res.URL, nil
}
