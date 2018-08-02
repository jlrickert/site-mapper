package main

//go:generate go run scripts/genResources.go

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func FindUniqueLinks(url string, options CrawlerOptions) int {
	chUrls := make(chan string)
	urlCount := 0

	go func() {
		for {
			<-chUrls
			urlCount += 1
		}
	}()

	RecursiveCrawl(url, options, func(url string) {
		fmt.Println(" - " + url)
		chUrls <- url
	})
	close(chUrls)
	return urlCount
}

func GenerateSiteMapDotFile(url string, options CrawlerOptions) {
	file, err := os.Create("out.dot")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	site := IndexWebsite(url, options)
	site.GenerateDOT(file)
}

func GenerateSiteMapIndex(url string, options CrawlerOptions) {
	file, err := os.Create("index.html")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	site := IndexWebsite(url, options)
	site.GenerateIndexHtml(file)
}

func FindUniqueLinksCommand() {
	seedUrl := os.Args[1]
	urlCount := FindUniqueLinks(seedUrl, CrawlerOptions{
		maxCrawlers: 5,
		Throttle:    1000 * time.Millisecond,
	})
	fmt.Println("Unique Url count: " + strconv.Itoa(urlCount))
}

func GenerateSiteMapDotFileCommand() {
	seedUrl := os.Args[1]
	GenerateSiteMapDotFile(seedUrl, CrawlerOptions{
		maxCrawlers: 1,
		Throttle:    1000 * time.Millisecond,
	})
}

func GenerateSiteMapIndexCommand() {
	seedUrl := os.Args[1]
	GenerateSiteMapIndex(seedUrl, CrawlerOptions{
		maxCrawlers: 1,
		Throttle:    1000 * time.Millisecond,
	})
}

func main() {
	// GenerateSiteMapDotFileCommand()
	GenerateSiteMapIndexCommand()
	// FindUniqueLinksCommand()
}
