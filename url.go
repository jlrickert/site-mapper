package main

type Url struct {
	Path []string
	Href string
}

func NewUrl(url string) *Url {
	return &Url{
		Path: []string{},
		Href: url,
	}
}

func (u *Url) Link(url string) *Url {
	return &Url{
		Path: append(u.Path, u.Href),
		Href: url,
	}
}
