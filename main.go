package main

import (
	"flag"
	"fmt"
	"githubComment/github"
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
		}
		flag.Usage()
		os.Exit(-1)
	} else {
		fmt.Println(org, repo)
	}
}

func main() {

	url, err := github.AddComment(org, repo, comment)
	if err != nil {
		panic(err)
	}

	fmt.Printf("comment added: %s\n", url)
}
