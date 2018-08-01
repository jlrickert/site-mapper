package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUrlSpec(t *testing.T) {
	Convey("Given a url path", t, func() {
		UrlPath := NewUrlPath("google.com", "github.com")
		Convey("should have an href", func() {
			So(UrlPath.Href(), ShouldEqual, "github.com")
		})
	})

	Convey("Given a path of urls without a cycle", t, func() {
		urls := []string{"google.com", "google.com", "github.com/user"}
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
		urls := []string{"google.com", "google.com", "github.com/user", "google.com", "google.com"}
		path := NewUrlPath()
		So(path.AddLink(urls[0]), ShouldBeTrue)
		So(path.AddLink(urls[1]), ShouldBeTrue)
		So(path.AddLink(urls[2]), ShouldBeTrue)
		So(path.AddLink(urls[3]), ShouldBeFalse)
		So(path.AddLink(urls[4]), ShouldBeFalse)
		So(len(path), ShouldEqual, 3)
	})
}
