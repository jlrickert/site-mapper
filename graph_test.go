package main

import (
	// "bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGraphSpec(t *testing.T) {
	Convey("Given a site map", t, func() {
		// urls := []*Url{
		// 	NewUrlFromSlice([]string{"google.com", "google.com"}),
		// 	NewUrlFromSlice([]string{"google.com", "google.com/maps"}),
		// 	NewUrlFromSlice([]string{"google.com", "google.com/maps", "google.com"}),
		// 	NewUrlFromSlice([]string{"google.com", "google.com/email"}),
		// 	NewUrlFromSlice([]string{"google.com", "google.com/email", "google.com/email"}),
		// 	NewUrlFromSlice([]string{"google.com", "google.com/email", "google.com"}),
		// 	NewUrlFromSlice([]string{"google.com", "google.com/images"}),
		// 	NewUrlFromSlice([]string{"google.com", "google.com/images", "google.com/images"}),
		// 	NewUrlFromSlice([]string{"google.com", "google.com/images", "google.com"}),
		// }
		// site := NewSiteMapFromSlice("google.com", urls)

		// Convey("should generate a valid standalone html file", func() {
		// 	file := bytes.NewBuffer([]byte{})

		// 	WriteGraphIndex(file, site)

		// 	contents := string(file.Bytes())
		// 	So(contents, ShouldContainSubstring, "google.com")
		// 	So(contents, ShouldContainSubstring, "<html")
		// 	So(contents, ShouldContainSubstring, "</html>")
		// 	So(contents, ShouldContainSubstring, "<head")
		// 	So(contents, ShouldContainSubstring, "</head>")
		// 	So(contents, ShouldContainSubstring, "<body")
		// 	So(contents, ShouldContainSubstring, "</body>")
		// })
	})
}
