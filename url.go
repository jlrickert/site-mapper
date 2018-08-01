package main

type UrlPath []string

func NewUrlPath(urls ...string) UrlPath {
	path := UrlPath{}
	for i := range urls {
		path.AddLink(urls[i])
	}
	return path
}

func (path *UrlPath) AddLink(url string) bool {
	hits := 0
	for i := 0; i < len(*path); i++ {
		if url == (*path)[i] {
			hits++
		}
		if hits >= 2 {
			return false
		}
	}
	*path = append(*path, url)
	return true
}

func (path *UrlPath) Href() string {
	return (*path)[len(*path)-1]
}

// func NewUrlFromSlice(path []string) *Url {
// 	return &Url{
// 		Path: path[:len(path)-1],
// 		Href: path[len(path)-1],
// 	}
// }

// func (u *Url) Link(url string) *Url {
// 	return &Url{
// 		Path: append(u.Path, u.Href),
// 		Href: url,
// 	}
// }
