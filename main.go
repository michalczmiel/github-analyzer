package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	ForkCount                  int
	RepositoriesCount          int
	AverageCommitsCountPerRepo int
	UniqueLanguages            []string
	PortfolioUrl               string
	AccountAgeInYears          int
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

	var uniqueLanguages []string
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
		AverageCommitsCountPerRepo: averageCommitsCountPerRepo,
		RepositoriesCount:          repositoriesCount,
		ForkCount:                  forkCount,
		UniqueLanguages:            uniqueLanguages,
		PortfolioUrl:               *user.Blog,
		AccountAgeInYears:          accountAgeInYears,
	}

	return &userStats, nil
}

type BodyRequest struct {
	Username string `json:"username"`
}

type BodyResponse struct {
	PortfolioUrl               string   `json:"portfolio_url"`
	ForkCount                  int      `json:"fork_count"`
	RepositoriesCount          int      `json:"repositories_count"`
	AverageCommitsCountPerRepo int      `json:"average_commits_count_per_repo"`
	UniqueLanguages            []string `json:"unique_languages"`
	AccountAgeInYears          int      `json:"account_age_in_years"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bodyRequest := BodyRequest{
		Username: "",
	}

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	userStats, err := fetchUserStats(bodyRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	bodyResponse := BodyResponse{
		PortfolioUrl:               userStats.PortfolioUrl,
		ForkCount:                  userStats.ForkCount,
		RepositoriesCount:          userStats.RepositoriesCount,
		AverageCommitsCountPerRepo: userStats.AverageCommitsCountPerRepo,
		UniqueLanguages:            userStats.UniqueLanguages,
		AccountAgeInYears:          userStats.AccountAgeInYears,
	}

	response, err := json.Marshal(&bodyResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
