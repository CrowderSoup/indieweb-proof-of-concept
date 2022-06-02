package micropub

import (
	"net/url"
	"time"
)

type Post struct {
	Name       string
	Content    string
	Categories []string
	Tags       []string
	Date       time.Time
}

func GetPostFromForm(form url.Values) (*Post, string, error) {
	return nil, "", nil
}
