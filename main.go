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

func FindUniqueLinksCommand() {
	seedUrl := os.Args[1]
	urlCount := FindUniqueLinks(seedUrl, CrawlerOptions{
		maxCrawlers: 100,
		Throttle:    0 * time.Millisecond,
	})
	fmt.Println("Unique Url count: " + strconv.Itoa(urlCount))
}

func GenerateSiteMapDotFileCommand() {
	seedUrl := os.Args[1]
	GenerateSiteMapDotFile(seedUrl, CrawlerOptions{
		maxCrawlers: 100,
		Throttle:    0 * time.Millisecond,
	})
}

func main() {
	GenerateSiteMapDotFileCommand()
	// FindUniqueLinksCommand()
}
