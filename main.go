package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jlrickert/site-mapper/crawler"
)

func FindUniqueLinksCommand(crawler *crawler.Crawler) {
	seedUrl := os.Args[1]
	chHref := make(chan string)
	file, err := os.Create("sitemap")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	urlCount := 0

	file.Write([]byte(fmt.Sprintf("\"%s\" sitemap:\n", seedUrl)))
	go func() {
		for href := range chHref {
			file.Write([]byte(fmt.Sprintf(" - %s\n", href)))
			urlCount += 1
		}
	}()

	crawler.RecursiveCrawl(seedUrl, func(url string) {
		chHref <- url
	})
	close(chHref)
	file.Write([]byte(fmt.Sprintf("Unique url count: %d", urlCount)))
}

func GenerateSiteMapDotFileCommand(crawler *crawler.Crawler) {
	seedUrl := os.Args[1]
	file, err := os.Create("out.dot")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	site := crawler.IndexWebsite(seedUrl)
	site.GenerateDOT(file)
}

func GenerateSiteMapIndexCommand(crawler *crawler.Crawler) {
	seedUrl := os.Args[1]

	file, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	site := crawler.IndexWebsite(seedUrl)
	site.GenerateIndexHtml(file)
}

func main() {
	flag.Parse()
	crawler := crawler.New(1000, 1)
	// GenerateSiteMapDotFileCommand(crawler)
	GenerateSiteMapIndexCommand(crawler)
	FindUniqueLinksCommand(crawler)
}
