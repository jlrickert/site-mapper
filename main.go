package main

import (
	"fmt"
	"os"
)

func main() {
	seedUrls := os.Args[1:]

	crawler := NewCrawler()
	for _, url := range seedUrls {
		go crawler.RecursiveCrawl(url)
	}

	fmt.Println("Unique Urls: ")

	for running := true; running; {
		select {
		case url := <-crawler.ChUniqueUrls:
			fmt.Println(" - " + url)
		case <-crawler.ChFinished:
			running = false
		}
	}
}
