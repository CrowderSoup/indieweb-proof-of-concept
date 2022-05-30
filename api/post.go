package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_PAT")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	for _, repo := range repos {
		fmt.Fprintf(w, *repo.FullName)
	}
}
