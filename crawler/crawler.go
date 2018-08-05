package crawler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Crawler struct {
	Throttle    time.Duration
	MaxCrawlers int

	running        int
	runningLock    sync.Mutex // lock for running counter
	activeCrawlers Semaphore
	chRateLimit    <-chan time.Time
}

func New(throttle time.Duration, maxCrawlers int) *Crawler {
	crawler := Crawler{
		Throttle:    throttle * time.Millisecond,
		MaxCrawlers: maxCrawlers,

		running:     0,
		runningLock: sync.Mutex{},
	}
	if maxCrawlers > 0 {
		crawler.activeCrawlers = make(Semaphore, maxCrawlers)
	}
	if throttle > 0 {
		crawler.chRateLimit = time.Tick(throttle * time.Millisecond)
	}
	return &crawler
}

func (crawler *Crawler) modRunner(cb func(running int) int) {
	crawler.runningLock.Lock()
	crawler.running = cb(crawler.running)
	crawler.runningLock.Unlock()
}

func (crawler *Crawler) Crawl(rootHref string, fn func(href string)) error {
	if !Crawlable(rootHref) {
		return nil
	}

	crawler.modRunner(func(n int) int {
		log.Printf(`Queueing runner for "%s". remaining: %d`, rootHref, n+1)
		return n + 1
	})
	if crawler.Throttle > 0 {
		<-crawler.chRateLimit
	}
	if crawler.MaxCrawlers > 0 {
		crawler.activeCrawlers.Acquire(1)
	}
	err := crawler.crawl(rootHref, func(href string) {
		if !hasProto(href) {
			href = fmt.Sprintf("%s/%s", strings.TrimRight(rootHref, "/"), href)
		}
		u, err := url.Parse(href)
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
		fn(u.String())
	})
	crawler.modRunner(func(n int) int {
		log.Printf(`Runner for "%s" finished. %d remaining`, rootHref, n-1)
		return n - 1
	})
	if crawler.MaxCrawlers > 0 {
		crawler.activeCrawlers.Release(1)
	}

	return err
}

func (crawler *Crawler) crawl(rootHref string, fn func(href string)) error {
	log.Println("Crawling", rootHref)

	resp, err := http.Get(rootHref)
	if err != nil {
		return fmt.Errorf("ERROR: Failed to crawl \""+rootHref+"\"", err)
	}

	// contentType := resp.Header["Content-Type"]
	// if contentType != nil || len(contentType) != 0 {
	// 	validTypes := []string{
	// 		"text/html",
	// 	}
	// 	ok := false
	// 	for cti := range contentType {
	// 		ct := contentType[cti]
	// 		for vti := range validTypes {
	// 			vt := validTypes[vti]
	// 			if strings.Contains(ct, vt) {
	// 				ok = true
	// 				break
	// 			}
	// 		}
	// 	}
	// 	if !ok {
	// 		return nil
	// 	}
	// }

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

func (crawler *Crawler) RecursiveCrawl(rootHref string, fn func(href string)) {
	chHrefs := make(chan string)
	chFinished := make(chan bool)

	hrefHandler := func(href string) {
		chHrefs <- href
	}

	go func() {
		crawler.Crawl(rootHref, hrefHandler)
		chFinished <- true
	}()

	running := 1
	hrefs := make(map[string]bool)
	for running > 0 {
		select {
		case href := <-chHrefs:
			if !hrefs[href] {
				hrefs[href] = true
				fn(href)
				if strings.Contains(href, rootHref) && !strings.Contains(href, "..") {
					running++
					go func() {
						crawler.Crawl(href, hrefHandler)
						chFinished <- true
					}()
				}
			}
		case <-chFinished:
			running--
		}
	}
}

func (crawler *Crawler) IndexWebsite(rootHref string) *SiteMap {
	log.Println("Indexing", rootHref)
	site := NewSiteMap(rootHref)
	chHref := make(chan UrlPath)
	chFinished := make(chan bool)

	mkHandler := func(path UrlPath) func(string) {
		return func(href string) {
			p := path.Clone()
			hits := p.AddLink(href)
			if p.Href() == "" {
				return
			}
			if hits <= 1 {
				chHref <- *p
			}
		}
	}

	rootPath := NewUrlPath(rootHref)

	running := 1
	go func() {
		err := crawler.Crawl(rootHref, mkHandler(rootPath))
		if err != nil {
			log.Println(err)
		}
		chFinished <- true
	}()

	for running != 0 {
		select {
		case path := <-chHref:
			href := path.Href()
			addedToSiteMap := site.AddUrlPath(*path.Clone())
			partOfDomain := strings.Index(href, rootHref) == 0
			if addedToSiteMap && partOfDomain {
				running++
				go func() {
					err := crawler.Crawl(href, mkHandler(path))
					if err != nil {
						log.Println("ERROR:", err)
					}
					chFinished <- true
				}()
			}
		case <-chFinished:
			running--
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
