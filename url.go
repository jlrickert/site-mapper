package main

import (
	"fmt"
	"strings"
)

type UrlPath []string

func NewUrlPath(urls ...string) UrlPath {
	path := UrlPath{}
	for i := range urls {
		path.AddLink(urls[i])
	}
	return path
}

func (path *UrlPath) AddLink(url string) int {
	if len(*path) == 0 && !hasProto(url) {
		panic(fmt.Sprintf("Root of %s needs to have the protocol", url))
	}

	i := strings.Index(url, "..")
	if i >= 0 {
		url = url[:i]
	}

	hits := 0
	for i := 0; i < len(*path); i++ {
		if url == (*path)[i] {
			hits++
		}
	}

	if len(url) > 0 && len(*path) >= 1 && !hasProto(url) {
		if url[0] == '.' {
			url = strings.TrimLeft(url, ".")
		}

		if url[0] == '/' {
			url = strings.TrimLeft(url, "/")
		}
		root := (*path)[len(*path)-1]
		root = strings.TrimRight(root, "/")
		url = fmt.Sprintf("%s/%s", root, url)
	}

	*path = append(*path, url)
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

func hasProto(url string) bool {
	return strings.Index(url, "http") == 0 || strings.Index(url, "https") == 0
}
