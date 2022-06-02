package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"

	"github.com/crowdersoup/indieweb-proof-of-concept/app/config"
	"github.com/crowdersoup/indieweb-proof-of-concept/app/micropub"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Authorize Request
	err := micropub.AuthorizeRequest(r)
	if err != nil {
		fmt.Fprintf(w, "Not Authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Ensure Content-Type is correct
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// Attempt to load the config from the environment
	config, err := config.GetConfig()
	if err != nil {
		fmt.Fprintf(w, "There was an error loading the config: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Attempt to parse the request x-www-form-urlencoded data
	err = r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Unable to parse form: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	post, renderedContent, err := micropub.GetPostFromForm(r.Form)
	if err != nil {
		fmt.Fprintf(w, "Unable to build post: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create an empty context that won't be concelled
	ctx := context.Background()

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHubPersonalAccessToken},
	)
	tokenClient := oauth2.NewClient(ctx, tokenSource)

	client := github.NewClient(tokenClient)

	commiter := github.CommitAuthor{
		Name:  &config.AuthorName,
		Email: &config.AuthorEmail,
	}
	message := fmt.Sprintf(
		"A new post titled \"%s\" on %s",
		post.Name,
		post.Date.Format("2006-01-02 15:04"),
	)
	content := []byte(renderedContent)
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
		fmt.Sprintf("content/posts/%s.md", post.Slug),
		&file,
	)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	contentResponseJson, err := json.Marshal(contentResponse)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(contentResponseJson))
	w.WriteHeader(http.StatusOK)
}
