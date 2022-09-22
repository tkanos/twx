package twtfile

import (
	"bytes"
	"html/template"
	"regexp"
	"strconv"
	"strings"
)

var tmpl = template.Must(template.New("twtxtTemplate").Parse(defaultTemplate))

const defaultTemplate = `# Twtxt is an open, distributed microblogging platform that
# uses human-readable text files, common transport protocols,
# and free software.
#
# Learn more about twtxt at  https://github.com/buckket/twtxt
#
# nick        = {{ .Nick }}
# url         = {{ .URL }}
# avatar      = {{ .Avatar }}
# description = {{ .Description }}
#
# followers   = {{ .Followers }}
# following   = {{ .Following }}
#
{{- if .Link}}
{{range $index, $element := .Link -}}
# link = {{ $index }} {{ $element }}
{{ end -}}
#{{ end -}}
{{- if .Follow }}
{{ range $index, $element := .Follow -}}
# follow = {{ $index }} {{ $element }}
{{ end -}}
{{ end -}}
# 
`

type TweetMetadata struct {
	Nick        string
	URL         string
	Avatar      string
	Description string

	Followers int
	Following int

	Link   map[string]string
	Follow map[string]string
}

func (t TweetMetadata) CreateTwtxtMetaTemplate() ([]byte, error) {
	buf := bytes.Buffer{}
	t.Following = len(t.Follow)
	if err := tmpl.Execute(&buf, t); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var reMeta = regexp.MustCompile(`^#\s+([a-zA-Z]+)\s+\=\s+(.+)$`)

func parseMeta(line string, meta TweetMetadata) TweetMetadata {
	if len(line) == 0 {
		return meta
	}

	groups := reMeta.FindStringSubmatch(line)
	if groups == nil || groups[0] == "" {
		return meta
	}

	key := strings.ToLower(groups[1])
	value := groups[2]

	switch key {
	case "nick":
		meta.Nick = value
	case "url":
		meta.URL = value
	case "avatar":
		meta.Avatar = value
	case "description":
		meta.Description = value
	case "followers":
		n, err := strconv.Atoi(value)
		if err == nil {
			meta.Followers = n
		}
	case "following":
		n, err := strconv.Atoi(value)
		if err == nil {
			meta.Following = n
		}
	case "link":
		if meta.Link == nil {
			meta.Link = map[string]string{}
		}
		v := strings.Split(strings.Trim(value, " "), " ")
		if len(v) > 1 {
			meta.Link[v[0]] = strings.Join(v[1:], " ")
		}
	case "follow":
		if meta.Follow == nil {
			meta.Follow = map[string]string{}
		}
		v := strings.Split(strings.Trim(value, " "), " ")
		if len(v) > 1 {
			meta.Follow[v[0]] = strings.Join(v[1:], " ")
		}
	}

	return meta
}
