package main

//go:generate go run scripts/genResources.go

import (
	"fmt"
	"os"
)

func main() {
	seedUrls := os.Args[1:]

	chUrls := make(chan Url)
	chFinished := make(chan bool)
	uniqueUrls := []Url{}
	urls := make([][]*Url, len(seedUrls))

	for i, url := range seedUrls {
		go func(url string, index int) {
			crawler := NewCrawler()
			urls[index] = crawler.RecursiveCrawl(url)
			for _, v := range urls[index] {
				chUrls <- *v
			}
			chFinished <- true
		}(url, i)
	}

	for running := len(seedUrls); running != 0; {
		select {
		case url := <-chUrls:
			uniqueUrls = append(uniqueUrls, url)
		case <-chFinished:
			running--
		}
	}

	fmt.Println("Unique Urls: ")

	for i := range uniqueUrls {
		url := uniqueUrls[i]
		fmt.Println(" -", url.Path, url.Href)
	}
	fmt.Println(len(uniqueUrls), "Unique urls")

	fmt.Println("----------------------------------------")
	fmt.Println("Printing Site map")
	fmt.Println("----------------------------------------")
	for i := range urls {
		site := NewSiteMapFromSlice(seedUrls[i], urls[i])
		site.Display(3)
	}
}
