package main

//go:generate go run scripts/genResources.go

import (
	"fmt"
	"os"
	// "strconv"
)

func FindUniqueLinks(url string) int {
	chUrls := make(chan string)
	urlCount := 0

	go func() {
		for {
			<-chUrls
			urlCount += 1
		}
	}()

	RecursiveCrawl(url, 1000, func(url string) {
		fmt.Println(" - " + url)
		chUrls <- url
	})
	close(chUrls)
	return urlCount
}

func GenerateSiteMapDotFile(url string) {
	file, err := os.Create("out.dot")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	site := IndexWebsite(url, 1000)
	site.GenerateDOT(file)
}

func main() {
	seedUrl := os.Args[1]

	GenerateSiteMapDotFile(seedUrl)

	// urlCount := FindUniqueLinks(seedUrl)
	// fmt.Println("Unique Url count: " + strconv.Itoa(urlCount))
}
