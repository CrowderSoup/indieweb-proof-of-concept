package micropub

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

func AuthorizeRequest(r *http.Request) error {
	err := attemptHeaderAuth(r)
	if err == nil {
		return nil
	}

	err = attemptFormAuth(r)
	return err
}

func attemptHeaderAuth(r *http.Request) error {
	authHeader := r.Header.Get("authorization")
	if authHeader == "" {
		return errors.New("Missing authorization header")
	}

	valueParts := strings.Split(authHeader, " ")
	if len(valueParts) <= 1 {
		return errors.New("Malformed authorization header")
	}

	token := valueParts[1]
	sysToken := os.Getenv("AUTH_TOKEN")

	if token != sysToken {
		return errors.New("Not Authorized")
	}

	return nil
}

func attemptFormAuth(r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	token, ok := r.Form["access_token"]
	if !ok {
		return errors.New("Missing access_token")
	}

	sysToken := os.Getenv("AUTH_TOKEN")

	if token[0] != sysToken {
		return errors.New("Not Authorized")
	}

	return nil
}
