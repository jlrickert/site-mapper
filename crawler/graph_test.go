package crawler

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGraphSpec(t *testing.T) {
	Convey("Given a site map", t, func() {
		site := NewSiteMapFromUrlPaths("https://google.com", []UrlPath{
			UrlPath{"https://google.com", "https://google.com"},
			UrlPath{"https://google.com", "https://google.com/maps"},
			UrlPath{"https://google.com", "https://google.com/maps", "https://google.com"},
			UrlPath{"https://google.com", "https://google.com/email"},
			UrlPath{"https://google.com", "https://google.com/email", "https://google.com/email"},
			UrlPath{"https://google.com", "https://google.com/email", "https://google.com"},
			UrlPath{"https://google.com", "https://google.com/images"},
			UrlPath{"https://google.com", "https://google.com/images", "https://google.com/images"},
			UrlPath{"https://google.com", "https://google.com/images", "https://google.com"},
		})

		Convey("should generate a valid standalone html file", func() {
			file := bytes.NewBuffer([]byte{})

			graph := newGraphHtmlTemplate(site)
			writeGraphIndex(file, graph)

			contents := string(file.Bytes())
			So(contents, ShouldContainSubstring, "google.com")
			So(contents, ShouldContainSubstring, "<html")
			So(contents, ShouldContainSubstring, "</html>")
			So(contents, ShouldContainSubstring, "<head")
			So(contents, ShouldContainSubstring, "</head>")
			So(contents, ShouldContainSubstring, "<body")
			So(contents, ShouldContainSubstring, "</body>")
		})
	})
}
