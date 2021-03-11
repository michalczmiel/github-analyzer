package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v33/github"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func fetchRepositories(username string) ([]*github.Repository, error) {
	client := github.NewClient(nil)
	repositories, _, err := client.Repositories.List(context.Background(), username, nil)
	return repositories, err
}

func main() {
	var username string
	fmt.Print("Enter GitHub username: ")
	fmt.Scanf("%s", &username)

	repositories, err := fetchRepositories(username)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	languages := make([]string, len(repositories))

	for _, repository := range repositories {
		language := *repository.Language
		if !contains(languages, language) {
			languages = append(languages, language)
		}
	}

	fmt.Printf("Number of repositories: %v \n", len(repositories))
	fmt.Print("Languages: ")
	for _, language := range languages {
		fmt.Printf("%v ", language)
	}
	fmt.Println()
}
