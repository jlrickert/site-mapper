package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUrlSpec(t *testing.T) {
	Convey("Given a url path", t, func() {
		urlPath := []UrlPath{
			NewUrlPath("https://google.com", "https://github.com"),
			NewUrlPath("https://google.com", "https://github.com", "/users"),
		}
		Convey("should have an href", func() {
			So(urlPath[0].Href(), ShouldEqual, "https://github.com")
			So(urlPath[1].Href(), ShouldEqual, "https://github.com/users")
		})
	})

	Convey("Given a path of urls without a cycle", t, func() {
		urls := []string{"https://google.com", "https://github.com", "https://github.com/user"}
		Convey("should correctly create the path", func() {
			path1 := NewUrlPath()

			for i := range urls {
				path1.AddLink(urls[i])
			}
			for i := range urls {
				So(path1[i], ShouldEqual, urls[i])
			}

			path2 := NewUrlPath(urls[0], urls[1], urls[2])
			for i := range urls {
				So(path2[i], ShouldEqual, urls[i])
			}
		})
	})

	Convey("Given a path of urls that contain a cycle", t, func() {
		urls := []string{"https://google.com", "https://github.com/user", "https://google.com", "https://google.com"}
		Convey("should correctly detect number of cycle", func() {
			path := NewUrlPath()
			So(path.AddLink(urls[0]), ShouldEqual, 0)
			So(path.AddLink(urls[1]), ShouldEqual, 0)
			So(path.AddLink(urls[2]), ShouldEqual, 1)
			So(path.AddLink(urls[3]), ShouldEqual, 2)
			So(len(path), ShouldEqual, 4)
		})
	})
}
