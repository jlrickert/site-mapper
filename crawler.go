package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"sync"
)

type Crawler struct {
	urls        map[*Url]bool
	visitedUrls map[string]bool
	urlsLock    sync.Mutex
}

func NewCrawler() *Crawler {
	return &Crawler{
		visitedUrls: make(map[string]bool),
		urls:        make(map[*Url]bool),
		urlsLock:    sync.Mutex{},
	}
}

func (crawler *Crawler) AddUrl(url *Url) (newVisit bool, newPath bool) {
	crawler.urlsLock.Lock()
	defer crawler.urlsLock.Unlock()

	newVisit = false
	newPath = false

	if !crawler.visitedUrls[url.Href] {
		crawler.visitedUrls[url.Href] = true
		newVisit = true
	}

	if !crawler.urls[url] {
		crawler.urls[url] = true
		newPath = true
	}
	return
}

func (crawler *Crawler) Crawl(url string) []Url {
	chUrls := make(chan *Url)
	chFinished := make(chan bool)

	rootUrl := NewUrl(url)
	go crawler.crawlUrl(rootUrl, chUrls, chFinished)

	for running := 1; running != 0; {
		select {
		case u := <-chUrls:
			crawler.AddUrl(u)
		case <-chFinished:
			running--
		}
	}

	keys := make([]Url, len(crawler.urls))
	for k := range crawler.urls {
		keys = append(keys, *k)
	}
	return keys
}

func (crawler *Crawler) RecursiveCrawl(url string) []Url {
	chUrls := make(chan *Url)
	chFin := make(chan bool)

	rootUrl := NewUrl(url)
	go crawler.crawlUrl(rootUrl, chUrls, chFin)

	for running := 1; running != 0; {
		select {
		case u := <-chUrls:
			newVisit, _ := crawler.AddUrl(u)
			if newVisit && strings.Contains(u.Href, rootUrl.Href) {
				running++
				go crawler.crawlUrl(u, chUrls, chFin)
			}
		case <-chFin:
			running--
		}
	}

	keys := make([]Url, len(crawler.urls))
	for k := range crawler.urls {
		keys = append(keys, *k)
	}
	return keys
}

func (crawler *Crawler) crawlUrl(url *Url, chUrl chan *Url, chFinished chan bool) {
	resp, err := http.Get(url.Href)
	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url.Href + "\"")
		chFinished <- true
		return
	}

	body := resp.Body
	defer body.Close()

	z := html.NewTokenizer(body)

	for {
		tokType := z.Next()
		switch {
		case tokType == html.ErrorToken:
			chFinished <- true
			return

		case tokType == html.StartTagToken:
			token := z.Token()

			isAnchor := token.Data == "a"
			if !isAnchor {
				continue
			}

			ok, href := getHref(token)
			if !ok {
				continue
			}

			hasProto := strings.Index(href, "http") == 0
			if hasProto {
				chUrl <- url.Link(href)
			}
		}
	}
}

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			ok = true
			href = a.Val
		}
	}
	return
}
