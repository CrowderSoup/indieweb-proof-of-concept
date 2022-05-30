package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	gitHubAccessToken := os.Getenv("GITHUB_PAT")
	if gitHubAccessToken == "" {
		fmt.Fprintf(w, "GITHUB_PAT (%s) is invalid", gitHubAccessToken)
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_PAT")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	name := "Aaron Crowder"
	email := "aaron.crowder@gmail.com"

	commiter := github.CommitAuthor{
		Name:  &name,
		Email: &email,
	}
	message := "A new file"
	content := []byte("This is the content of the file")
	branch := "testing-api-endpoint"

	file := github.RepositoryContentFileOptions{
		Message:   &message,
		Content:   content,
		Branch:    &branch,
		Committer: &commiter,
	}

	contentResponse, _, err := client.Repositories.CreateFile(
		ctx,
		"CrowderSoup",
		"indieweb-proof-of-concept",
		"content/posts",
		&file,
	)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	contentResponseJson, err := json.Marshal(contentResponse)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	fmt.Fprintf(w, string(contentResponseJson))
}
