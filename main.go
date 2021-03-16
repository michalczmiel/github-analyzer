package main

import (
	"context"
	"fmt"
	"time"

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
	options := &github.RepositoryListOptions{
		Sort:      "pushed",
		Direction: "desc",
	}
	repositories, _, err := client.Repositories.List(context.Background(), username, options)
	return repositories, err
}

func fetchUser(username string) (*github.User, error) {
	client := github.NewClient(nil)
	user, _, err := client.Users.Get(context.Background(), username)
	return user, err
}

type UserStats struct {
	forkCount                  int
	repositoriesCount          int
	averageCommitsCountPerRepo int
	uniqueLanguages            []string
	joinedAt                   time.Time
	portfolioUrl               string
	accountAgeInYears          int
}

func fetchUserStats(username string) (*UserStats, error) {
	const RepositoryCountToAnalyze = 6

	user, err := fetchUser(username)
	if err != nil {
		return nil, err
	}

	repositories, err := fetchRepositories(username)
	if err != nil {
		return nil, err
	}
	repositoriesCount := len(repositories)

	uniqueLanguages := make([]string, 20)
	commitsCount := 0
	forkCount := 0

	for index, repository := range repositories {
		if *repository.Fork {
			forkCount += 1
			continue
		}

		if index+1 > RepositoryCountToAnalyze {
			break
		}

		languages, err := fetchLanguages(username, *repository.Name)
		if err != nil {
			continue
		}

		for language, _ := range languages {
			if !contains(uniqueLanguages, language) {
				uniqueLanguages = append(uniqueLanguages, language)
			}
		}

		commits, err := fetchCommits(username, *repository.Name)
		if err != nil {
			continue
		}
		commitsCount += len(commits)
	}

	averageCommitsCountPerRepo := commitsCount / (repositoriesCount - forkCount)

	accountAgeInYears := (time.Now().Year() - user.CreatedAt.Year())

	userStats := UserStats{
		averageCommitsCountPerRepo: averageCommitsCountPerRepo,
		repositoriesCount:          repositoriesCount,
		forkCount:                  forkCount,
		uniqueLanguages:            uniqueLanguages,
		portfolioUrl:               *user.Blog,
		accountAgeInYears:          accountAgeInYears,
	}

	return &userStats, nil
}

func main() {
	var username string
	fmt.Print("Enter GitHub username: ")
	fmt.Scanf("%s", &username)

	userStats, err := fetchUserStats(username)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Portfolio url: %v \n", userStats.portfolioUrl)
	fmt.Printf("Account age in years: %v \n", userStats.accountAgeInYears)
	fmt.Printf("Number of repositories: %v \n", userStats.repositoriesCount)
	fmt.Printf("Number of forks: %v \n", userStats.forkCount)
	fmt.Print("Languages: ")
	for _, language := range userStats.uniqueLanguages {
		if language != "" {
			fmt.Printf("%v ", language)
		}
	}
	fmt.Println()
	fmt.Printf("Average number of commits in repository: %v \n", userStats.averageCommitsCountPerRepo)
}
