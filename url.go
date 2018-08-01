package main

type UrlPath []string

func NewUrlPath(urls ...string) UrlPath {
	path := UrlPath{}
	for i := range urls {
		path.AddLink(urls[i])
	}
	return path
}

func (path *UrlPath) AddLink(url string) int {
	hits := 0
	for i := 0; i < len(*path); i++ {
		if url == (*path)[i] {
			hits++
		}
	}
	*path = append(*path, url)
	return hits
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
