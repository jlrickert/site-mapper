package main

import (
	"fmt"
	"regexp"
	"strings"
	"net/url"
	"log"
)

var (
	tarRE = regexp.MustCompile(`\.tar\.gz$`)
	pdfRE = regexp.MustCompile(`\.pdf$`)
)

type UrlPath []string

func NewUrlPath(rawUrls ...string) UrlPath {
	path := UrlPath{}
	for i := range rawUrls {
		url, err := url.Parse(rawUrls[i])
		if err == nil {
			path.AddLink(url.String())
		} else {
			log.Println("ERROR:", err)
			path.AddLink(rawUrls[i])
		}
	}
	return path
}

func (path *UrlPath) AddLink(rawUrl string) int {
	if len(*path) == 0 && !hasProto(rawUrl) {
		panic(fmt.Sprintf("Root of %s needs to have the protocol", rawUrl))
	}

	i := strings.Index(rawUrl, "..")
	if i >= 0 {
		rawUrl = rawUrl[:i]
	}

	hits := 0
	for i := 0; i < len(*path); i++ {
		if rawUrl == (*path)[i] {
			hits++
		}
	}

	if len(rawUrl) > 0 && len(*path) >= 1 && !hasProto(rawUrl) {
		if rawUrl[0] == '.' {
			rawUrl = strings.TrimLeft(rawUrl, ".")
		}

		if rawUrl[0] == '/' {
			rawUrl = strings.TrimLeft(rawUrl, "/")
		}
		root := (*path)[len(*path)-1]
		root = strings.TrimRight(root, "/")
		rawUrl = fmt.Sprintf("%s/%s", root, rawUrl)
	}
	rawUrl = strings.TrimRight(rawUrl, "/")
	*path = append(*path, rawUrl)
	return hits
}

func (path *UrlPath) Root() string {
	return (*path)[0]
}

func (path *UrlPath) Href() string {
	return (*path)[len(*path)-1]
}

func (path *UrlPath) Clone() *UrlPath {
	newPath := NewUrlPath()
	for i := range *path {
		newPath.AddLink((*path)[i])
	}
	return &newPath
}

func hasProto(href string) bool {
	return strings.Index(href, "http") == 0 || strings.Index(href, "https") == 0
}

func Crawlable(href string) bool {
	fragment := strings.Index(href, "#") >= 0
	mailto := strings.Contains(href, "mailto:") && strings.Contains(href, "@")
	tele := strings.Contains(href, "tel:")
	isTar := tarRE.Match([]byte(href))
	isPdf := pdfRE.Match([]byte(href))
	return !(fragment || mailto || tele || isTar || isPdf)
}
