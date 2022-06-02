package micropub

import (
	"bytes"
	"errors"
	"html/template"
	"net/url"
	"time"

	"github.com/gosimple/slug"
)

type Post struct {
	HEntry     string
	Content    string
	Name       string
	Slug       string
	Categories []string
	Tags       []string
	Date       time.Time
}

const postTemplate = `---
title: "{{ .Name }}"
slug: "{{ .Slug }}"
date: {{ .Date.Format "2006-01-02T15:04:05" }}
type: {{ .HEntry }}
categories:{{ range .Categories }}
  - {{ . }}
{{ end }}
tags:{{ range .Tags }}
  - {{ . }}
{{ end }}
---

{{ .Content }}
`

func GetPostFromForm(form url.Values) (*Post, string, error) {
	hentry, ok := form["h-entry"]
	if !ok {
		return nil, "", errors.New("Missing h-entry, which is required")
	}

	content, ok := form["content"]
	if !ok {
		return nil, "", errors.New("Missing content, which is required")
	}

	name, ok := form["name"]
	if !ok {
		name = []string{""}
	}

	categories, ok := form["category"]
	if !ok {
		categories = []string{}
	}

	tags, ok := form["tag"]
	if !ok {
		tags = []string{}
	}

	post := &Post{
		HEntry:     hentry[0],
		Content:    content[0],
		Name:       name[0],
		Slug:       getSlug(name[0], content[0]),
		Categories: categories,
		Tags:       tags,
		Date:       time.Now(),
	}

	tpl, err := template.New("post").Parse(postTemplate)
	if err != nil {
		return post, "", err
	}

	var b bytes.Buffer
	err = tpl.Execute(&b, post)
	if err != nil {
		return post, "", err
	}

	return post, b.String(), nil
}

func getSlug(name, content string) string {
	if name == "" {
		limit := 50
		if len(content) < 50 {
			limit = len(content)
		}
		name = content[:limit]
	}

	urlSlug := slug.Make(name)
	return urlSlug
}
