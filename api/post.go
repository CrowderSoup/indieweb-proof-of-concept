package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"

	"github.com/crowdersoup/indieweb-proof-of-concept/app/config"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Create an empty context that won't be concelled
	ctx := context.Background()

	// Attempt to load the config from the environment
	config, err := config.GetConfig()
	if err != nil {
		fmt.Fprintf(w, "There was an error loading the config: %s", err.Error())
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHubPersonalAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	commiter := github.CommitAuthor{
		Name:  &config.AuthorName,
		Email: &config.AuthorEmail,
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
