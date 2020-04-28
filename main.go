package main

import (
	"flag"
	"fmt"
	"github.com/nunnatsa/github-commenter/githubclient"
	"os"
)

var (
	org, repo, comment string
	help bool
)

func init() {
	flag.StringVar(&org, "o", "", "github organization; mandatory")
	flag.StringVar(&repo, "r", "", "github repository; mandatory")
	flag.StringVar(&comment, "m", "", "new comment text; mandatory")
	flag.BoolVar(&help, "h", false, "display this help")

	flag.Parse()

	if help || org == "" || repo == "" || comment == "" {
		if !help {
			fmt.Println("Error: missing parameters")
			fmt.Println()
		}
		flag.Usage()
		os.Exit(-1)
	}
}

func main() {
	client, err := githubclient.NewClient()
	if err != nil {
		fmt.Println("can't access github;", err)
		os.Exit(-2)
	}

	url, err := client.AddComment(org, repo, comment)
	if err != nil {
		fmt.Println("failed to add a comment;", err)
		os.Exit(-3)
	}

	fmt.Printf("comment added: %s\n", url)
}
