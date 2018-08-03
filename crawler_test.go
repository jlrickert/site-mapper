package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCrawlerSpec(t *testing.T) {
	Convey("Given a crawler with the cache option set", t, func() {
		// crawler := NewCrawler(0, -1, true)
		// Convey("should properly cache already used", func() {
		// 	urls := []string{
		// 		"https://google.com",
		// 		"https://github.com",
		// 		"https://duckduckgo.com",
		// 	}
		// 	So(crawler.addLinks(urls[0], urls[1:3]), ShouldBeTrue)
		// 	So(crawler.addLinks(urls[0], urls[1:3]), ShouldBeFalse)
		// 	So(crawler.cache[urls[0]], ShouldEqual, urls[1:3])
		// })
	})
}
