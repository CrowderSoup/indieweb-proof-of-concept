package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	AuthorName                string
	AuthorEmail               string
	GitHubPersonalAccessToken string
}

func GetConfig() (*Config, error) {
	gitHubAccessToken := os.Getenv("GITHUB_PAT")
	if gitHubAccessToken == "" {
		return &Config{}, errors.New(fmt.Sprintf("GITHUB_PAT (%s) is invalid", gitHubAccessToken))
	}

	authorName := os.Getenv("AUTHOR_NAME")
	if authorName == "" {
		return &Config{}, errors.New(fmt.Sprintf("AUTHOR_NAME (%s) is invalid", authorName))
	}

	authorEmail := os.Getenv("AUTHOR_EMAIL")
	if authorEmail == "" {
		return &Config{}, errors.New(fmt.Sprintf("AUTHOR_EMAIL (%s) is invalid", authorEmail))
	}

	return &Config{
		AuthorName:                authorName,
		AuthorEmail:               authorEmail,
		GitHubPersonalAccessToken: gitHubAccessToken,
	}, nil
}
