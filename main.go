package main

import (
	"fmt"
	"githubComment/github"
)

func main() {
	user, err := github.GetGitHead()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Usename: %s\n", user)

	currentPullReq, err := github.GetCurrentPullReq()
	if err != nil {
		panic(err)
	}

	fmt.Printf("GetCurrentPullReq: %d\n", currentPullReq)
}
