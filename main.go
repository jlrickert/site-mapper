package main

//go:generate go run scripts/genResources.go

import (
	"fmt"
	"os"
	"strconv"
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

	RecursiveCrawl(url, func(url string) {
		fmt.Println(" - " + url)
		chUrls <- url
	})
	close(chUrls)
	return urlCount
}

func main() {
	seedUrl := os.Args[1]

	urlCount := FindUniqueLinks(seedUrl)
	fmt.Println("Unique Url count: " + strconv.Itoa(urlCount))

	// uniquUrls := make(map[string]bool)
	// for i := range urlPaths {
	// 	url := urlPaths[i]
	// 	uniquUrls[url] = true
	// }

	// fmt.Println(" -", url.Path, url.Href)

	// fmt.Println(len(uniqueUrls), "Unique urls")

	// fmt.Println("----------------------------------------")
	// fmt.Println("Printing Site map")
	// fmt.Println("----------------------------------------")
	// for i := range urls {
	// 	site := NewSiteMapFromSlice(seedUrls[i], urls[i])
	// 	site.Display(3)
	// 	site.Graph()
	// }
}
