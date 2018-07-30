package main

import (
	"fmt"
	"os"
)

func main() {
	seedUrls := os.Args[1:]

	crawler := NewCrawler()
	chUrls := make(chan Url)
	chFinished := make(chan bool)
	urls := []Url{}

	for _, url := range seedUrls {
		go func(url string) {
			urls := crawler.RecursiveCrawl(url)
			for _, v := range urls {
				chUrls <- v
			}
			chFinished <- true
		}(url)
	}

	for running := len(seedUrls); running != 0; {
		select {
		case url := <-chUrls:
			urls = append(urls, url)
		case <-chFinished:
			running--
		}
	}

	fmt.Println("Unique Urls: ")

	for i := range urls {
		url := urls[i]
		fmt.Println(" -", url.Path, url.Href)
	}
	fmt.Println(len(urls), "Unique urls")
}
