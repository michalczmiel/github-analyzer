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

func fetchCommits(username string, repository string) ([]*github.RepositoryCommit, error) {
	client := github.NewClient(nil)
	commits, _, err := client.Repositories.ListCommits(context.Background(), username, repository, nil)
	return commits, err
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
	commitsCount := 0
	forkCount := 0

	for _, repository := range repositories {
		language := *repository.Language

		if !contains(languages, language) {
			languages = append(languages, language)
		}

		if *repository.Fork {
			forkCount += 1
		} else {
			commits, err := fetchCommits(username, *repository.Name)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			commitsCount += len(commits)
		}
	}

	averageCommitsCountPerRepo := commitsCount / (len(repositories) - forkCount)

	fmt.Printf("Number of repositories: %v \n", len(repositories))
	fmt.Printf("Number of forks: %v \n", forkCount)
	fmt.Print("Languages: ")
	for _, language := range languages {
		fmt.Printf("%v ", language)
	}
	fmt.Println()
	fmt.Printf("Average number of commits in repository: %v \n", averageCommitsCountPerRepo)
}
