package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func RecursiveCrawl(rootUrl string, fn func(url string)) {
	chUrls := make(chan string)
	chFinished := make(chan bool)
	crawlers := make(Semaphore, 10)

	urls := make(map[string]bool)

	handleFoundUrl := func(url string) {
		chUrls <- url
	}

	running := 1
	go func() {
		Crawl(rootUrl, handleFoundUrl)
		chFinished <- true
	}()

	for running != 0 {
		select {
		case url := <-chUrls:
			if !urls[url] {
				urls[url] = true
				go fn(url)
				if strings.Contains(url, rootUrl) {
					running++
					go func() {
						crawlers.Wait(1)
						Crawl(url, handleFoundUrl)
						crawlers.Signal()
						chFinished <- true
					}()
				}
			}
		case <-chFinished:
			running--
		}
	}
}

func Crawl(url string, fn func(url string)) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: Failed to crawl \""+url+"\"", err)
		return
	}

	body := resp.Body
	defer body.Close()

	z := html.NewTokenizer(body)
	for {
		tokType := z.Next()
		switch {
		case tokType == html.ErrorToken:
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

			hasProto := strings.Index(href, "http") == 0 || strings.Index(href, "https") == 0
			if hasProto {
				go fn(href)
			}
		}
	}
}

// func IndexWebsite(url string, maxCrawlers int) *SiteMap {
// }

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			ok = true
			href = a.Val
		}
	}
	return
}
