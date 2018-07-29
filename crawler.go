package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"sync"
)

type Crawler struct {
	ChUniqueUrls chan string
	ChFinished   chan bool

	running     int
	runningLock sync.Mutex

	urls     map[string]bool
	urlsLock sync.Mutex
}

func NewCrawler() *Crawler {
	return &Crawler{
		ChUniqueUrls: make(chan string),
		ChFinished:   make(chan bool),
		running:      0,
		runningLock:  sync.Mutex{},
		urls:         make(map[string]bool),
		urlsLock:     sync.Mutex{},
	}
}

func (crawler *Crawler) AddUrl(url string) bool {
	crawler.urlsLock.Lock()
	defer crawler.urlsLock.Unlock()
	if !crawler.urls[url] {
		crawler.urls[url] = true
		crawler.ChUniqueUrls <- url
		return true
	}
	return false
}

func (crawler *Crawler) updateRunning(fn func(n int) int) int {
	crawler.runningLock.Lock()
	count := fn(crawler.running)
	crawler.running = count
	crawler.runningLock.Unlock()
	return count
}

func (crawler *Crawler) Crawl(url string) {
	crawler.updateRunning(add1)

	chUrls := make(chan string)
	chFinished := make(chan bool)
	go crawler.crawlUrl("", url, chUrls, chFinished)

	for running := true; running; {
		select {
		case u := <-chUrls:
			crawler.AddUrl(u)
		case <-chFinished:
			if crawler.updateRunning(sub1) == 0 {
				crawler.ChFinished <- true
			}
			running = false
		}
	}
}

func (crawler *Crawler) RecursiveCrawl(url string) {
	crawler.updateRunning(add1)

	chUrls := make(chan string)
	chFin := make(chan bool)

	go crawler.crawlUrl("", url, chUrls, chFin)

	go func() {
		for running := true; running; {
			select {
			case u := <-chUrls:
				if crawler.AddUrl(u) && strings.Contains(u, url) {
					crawler.updateRunning(add1)
					go crawler.crawlUrl(url, u, chUrls, chFin)
				}
			case <-chFin:
				if crawler.updateRunning(sub1) == 0 {
					crawler.ChFinished <- true
					running = false
				}
			}
		}
	}()
}

func (crawler *Crawler) crawlUrl(rootUrl, url string, chUrl chan string, chFinished chan bool) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
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

			ok, url := getHref(token)
			if !ok {
				continue
			}

			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				chUrl <- url
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
