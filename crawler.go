package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"
	"time"
)

type CrawlerOptions struct {
	Throttle    time.Duration
	maxCrawlers int
}

func RecursiveCrawl(rootUrl string, options CrawlerOptions, fn func(url string)) {
	chUrls := make(chan string)
	chFinished := make(chan bool)
	crawlers := make(Semaphore, options.maxCrawlers)

	urls := make(map[string]bool)

	handleFoundUrl := func(url string) {
		chUrls <- url
	}

	var rateLimit <-chan time.Time
	if options.Throttle > 0 {
		rateLimit = time.Tick(options.Throttle)
	}

	running := 1
	go func() {
		Crawl(rootUrl, handleFoundUrl)
		chFinished <- true
	}()

	for running != 0 {
		select {
		case url := <-chUrls:
			if !hasProto(url) {
				url = fmt.Sprintf(
					"%s/%s",
					strings.TrimRight(rootUrl, "/"),
					strings.TrimRight(url, "/"),
				)
			}

			if !urls[url] {
				urls[url] = true
				fn(url)
				if strings.Contains(url, rootUrl) {
					running++
					log.Println("New crawlers queued", running)
					go func() {
						if options.Throttle > 0 {
							<-rateLimit
						}
						if options.maxCrawlers > 0 {
							crawlers.Wait(1)
						}
						Crawl(url, handleFoundUrl)
						if options.maxCrawlers > 0 {
							crawlers.Signal()
						}
						chFinished <- true
					}()
				}
			}
		case <-chFinished:
			running--
			log.Println("Crawlers finished", running)
		}
	}
}

func Crawl(url string, fn func(url string)) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("ERROR: Failed to crawl \""+url+"\"", err)
	}

	contentType := resp.Header["Content-Type"]
	if contentType != nil || len(contentType) != 0 {
		validTypes := []string{
			"text/html",
		}
		ok := false
		for cti := range contentType {
			ct := contentType[cti]
			for vti := range validTypes {
				vt := validTypes[vti]
				if strings.Contains(ct, vt) {
					ok = true
					break
				}
			}
		}
		if !ok {
			return nil
		}
	}

	log.Println("Crawling " + url)

	body := resp.Body
	defer body.Close()

	z := html.NewTokenizer(body)
	for {
		tokType := z.Next()
		switch {
		case tokType == html.ErrorToken:
			if z.Err().Error() == "EOF" {
				return nil
			}
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

			fn(href)
		}
	}
}

func IndexWebsite(rootUrl string, options CrawlerOptions) *SiteMap {
	log.Println("Indexing " + rootUrl)
	site := NewSiteMap(rootUrl)
	chUrls := make(chan UrlPath)
	chFinished := make(chan bool)

	var crawlers Semaphore
	if options.maxCrawlers == 0 {
		crawlers = make(Semaphore, 10)
	} else if options.maxCrawlers > 0 {
		crawlers = make(Semaphore, options.maxCrawlers)
	}

	mkHandler := func(path UrlPath) func(url string) {
		return func(url string) {
			p := path.Clone()
			hits := p.AddLink(url)
			if p.Href() == "" {
				return
			}
			if hits <= 1 {
				chUrls <- *p
			}
		}
	}

	rootPath := NewUrlPath(rootUrl)

	running := 1
	go func() {
		err := Crawl(rootUrl, mkHandler(rootPath))
		if err != nil {
			fmt.Println(err)
		}
		chFinished <- true
	}()

	var rateLimit <-chan time.Time
	if options.Throttle > 0 {
		rateLimit = time.Tick(options.Throttle)
	}

	for running != 0 {
		select {
		case path := <-chUrls:
			url := path.Href()
			addedToSiteMap := site.AddUrlPath(*path.Clone())
			partOfDomain := strings.Index(url, rootUrl) == 0
			fragment := strings.Index(url, "#") >= 0
			mailto := strings.Contains(url, "mailto:") && strings.Contains(url, "@")
			tele := strings.Contains(url, "tel:")
			crawlable := !(fragment || mailto || tele)
			if addedToSiteMap && partOfDomain && crawlable {
				running++
				log.Println("New crawler queued for", url)
				go func() {
					if options.maxCrawlers > 0 {
						crawlers.Acquire(1)
					}
					if options.Throttle > 0 {
						<-rateLimit
					}
					err := Crawl(url, mkHandler(path))
					if err != nil {
						fmt.Println(err)
					}
					if options.maxCrawlers > 0 {
						crawlers.Signal()
					}
					chFinished <- true
				}()
			}
		case <-chFinished:
			running--
			log.Println("Crawler finished.", running, "remaining")
		}
	}
	return site
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
