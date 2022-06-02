---
title: "{{ .Name }}"
slug: "{{ .Slug }}"
date: {{ .Date.Format "2006-01-02 15:04" }}
type: {{ .HEntry }}
categories:
  {{ range .Categories }}
  - {{ . }}
  {{ end }}
tags:
  {{ range .Tags }}
  - {{ . }}
  {{ end }}
---

{{ .Content }}
