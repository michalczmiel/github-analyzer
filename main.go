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

func fetchLanguages(username string, repository string) (map[string]int, error) {
	client := github.NewClient(nil)
	languages, _, err := client.Repositories.ListLanguages(context.Background(), username, repository)
	return languages, err
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

	uniqueLanguages := make([]string, 20)
	commitsCount := 0
	forkCount := 0

	for _, repository := range repositories {
		languages, err := fetchLanguages(username, *repository.Name)
		if err != nil {
			continue
		}

		for language, _ := range languages {
			if !contains(uniqueLanguages, language) {
				uniqueLanguages = append(uniqueLanguages, language)
			}
		}

		if *repository.Fork {
			forkCount += 1
		} else {
			commits, err := fetchCommits(username, *repository.Name)
			if err != nil {
				continue
			}
			commitsCount += len(commits)
		}
	}

	averageCommitsCountPerRepo := commitsCount / (len(repositories) - forkCount)

	fmt.Printf("Number of repositories: %v \n", len(repositories))
	fmt.Printf("Number of forks: %v \n", forkCount)
	fmt.Print("Languages: ")
	for _, language := range uniqueLanguages {
		if language != "" {
			fmt.Printf("%v ", language)
		}
	}
	fmt.Println()
	fmt.Printf("Average number of commits in repository: %v \n", averageCommitsCountPerRepo)
}
