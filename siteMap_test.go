package main

import (
	"bufio"
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNodeSpec(t *testing.T) {
	Convey("Given a couple of urls that each link to each other", t, func() {
		url1 := "https://github.com"
		url2 := "https://www.google.com"

		Convey("should create a single node with correct href", func() {
			node := NewNode(url1)
			So(node.href, ShouldEqual, url1)
		})

		Convey("should create a single node with correct links", func() {
			node1 := NewNode(url1)
			node2 := NewNode(url2)
			node1.addLink(node2)
			So(node1.links[url2], ShouldEqual, node2)
			So(node2.links[url1], ShouldNotEqual, node1)
		})
	})
}

func TestSiteMapSpec(t *testing.T) {
	Convey("Given a root url", t, func() {
		url := "https://github.com"

		Convey("should correctly create sitemap href", func() {
			site := NewSiteMap(url)
			So(site.Href, ShouldEqual, url)
		})
		Convey("should correctly set root node", func() {
			url := "https://github.com"
			site := NewSiteMap(url)
			So(site.Root.href, ShouldEqual, url)
		})
	})

	Convey("Given a set of urls from a website", t, func() {
		paths := []UrlPath{}
		_ = paths

		urlTable := make(map[string]string)
		urlTable["a"] = "https://exampleA.com"
		urlTable["b"] = "https://exampleB.com"
		urlTable["c"] = "https://exampleC.com"
		urlTable["d"] = "https://exampleD.com"
		urls := []UrlPath{
			NewUrlPath(urlTable["a"], urlTable["a"]),
			NewUrlPath(urlTable["a"], urlTable["b"]),
			NewUrlPath(urlTable["a"], urlTable["b"], urlTable["c"]),
			NewUrlPath(urlTable["a"], urlTable["b"], urlTable["c"], urlTable["a"]),
			NewUrlPath(urlTable["a"], urlTable["b"], urlTable["d"]),
			NewUrlPath(urlTable["a"], urlTable["b"], urlTable["c"], urlTable["d"]),
			NewUrlPath(urlTable["a"], urlTable["b"], urlTable["c"], urlTable["c"]),
			NewUrlPath(urlTable["a"], urlTable["b"], urlTable["d"], urlTable["c"]),
		}

		Convey("should correctly build nodes", func() {
			nodeA := NewNode(urlTable["a"])
			nodeB := NewNode(urlTable["b"])
			nodeC := NewNode(urlTable["c"])
			nodeD := NewNode(urlTable["d"])

			nodeA.addLink(nodeB)
			nodeA.addLink(nodeA)

			nodeB.addLink(nodeC)
			nodeB.addLink(nodeD)

			nodeC.addLink(nodeA)
			nodeC.addLink(nodeD)

			nodeD.addLink(nodeC)

			site := NewSiteMapFromUrlPaths(urlTable["a"], urls)

			/////////
			// Node A
			/////////
			actualNodeA := site.GetNode(urlTable["a"])
			So(actualNodeA.links[urlTable["a"]].href, ShouldEqual, nodeA.href)
			So(actualNodeA.links[urlTable["b"]].href, ShouldEqual, nodeB.href)
			So(actualNodeA.links[urlTable["c"]], ShouldEqual, nil)
			So(actualNodeA.links[urlTable["d"]], ShouldEqual, nil)
			So(len(actualNodeA.links), ShouldEqual, 2)

			So(actualNodeA.linkedFrom[urlTable["a"]].href, ShouldEqual, nodeA.href)
			So(actualNodeA.linkedFrom[urlTable["b"]], ShouldEqual, nil)
			So(actualNodeA.linkedFrom[urlTable["c"]].href, ShouldEqual, nodeC.href)
			So(actualNodeA.linkedFrom[urlTable["d"]], ShouldEqual, nil)
			So(len(actualNodeA.linkedFrom), ShouldEqual, 2)

			/////////
			// Node B
			/////////
			actualNodeB := site.GetNode(urlTable["b"])
			So(actualNodeB.links[urlTable["a"]], ShouldEqual, nil)
			So(actualNodeB.links[urlTable["b"]], ShouldEqual, nil)
			So(actualNodeB.links[urlTable["c"]].href, ShouldEqual, nodeC.href)
			So(actualNodeB.links[urlTable["d"]].href, ShouldEqual, nodeD.href)
			So(len(actualNodeB.links), ShouldEqual, 2)

			So(actualNodeB.linkedFrom[urlTable["a"]].href, ShouldEqual, nodeA.href)
			So(actualNodeB.linkedFrom[urlTable["b"]], ShouldEqual, nil)
			So(actualNodeB.linkedFrom[urlTable["c"]], ShouldEqual, nil)
			So(actualNodeB.linkedFrom[urlTable["d"]], ShouldEqual, nil)
			So(len(actualNodeB.linkedFrom), ShouldEqual, 1)

			/////////
			// Node C
			/////////
			actualNodeC := site.GetNode(urlTable["c"])
			So(actualNodeC.links[urlTable["a"]].href, ShouldEqual, nodeA.href)
			So(actualNodeC.links[urlTable["b"]], ShouldEqual, nil)
			So(actualNodeC.links[urlTable["c"]].href, ShouldEqual, nodeC.href)
			So(actualNodeC.links[urlTable["d"]].href, ShouldEqual, nodeD.href)
			So(len(actualNodeC.links), ShouldEqual, 3)

			So(actualNodeC.linkedFrom[urlTable["a"]], ShouldEqual, nil)
			So(actualNodeC.linkedFrom[urlTable["b"]].href, ShouldEqual, nodeB.href)
			So(actualNodeC.linkedFrom[urlTable["c"]].href, ShouldEqual, nodeC.href)
			So(actualNodeC.linkedFrom[urlTable["d"]].href, ShouldEqual, nodeD.href)
			So(len(actualNodeC.linkedFrom), ShouldEqual, 3)

			/////////
			// Node D
			/////////
			actualNodeD := site.GetNode(urlTable["d"])
			So(actualNodeD.links[urlTable["a"]], ShouldEqual, nil)
			So(actualNodeD.links[urlTable["b"]], ShouldEqual, nil)
			So(actualNodeD.links[urlTable["c"]].href, ShouldEqual, nodeC.href)
			So(actualNodeD.links[urlTable["d"]], ShouldEqual, nil)
			So(len(actualNodeD.links), ShouldEqual, 1)

			So(actualNodeD.linkedFrom[urlTable["a"]], ShouldEqual, nil)
			So(actualNodeD.linkedFrom[urlTable["b"]].href, ShouldEqual, nodeB.href)
			So(actualNodeD.linkedFrom[urlTable["c"]].href, ShouldEqual, nodeC.href)
			So(actualNodeD.linkedFrom[urlTable["d"]], ShouldEqual, nil)
			So(len(actualNodeD.linkedFrom), ShouldEqual, 2)
		})
	})

	Convey("Given a valid sitemap", t, func() {
		Convey("should be able translate to DOT format", func() {
			site := NewSiteMapFromUrlPaths("https://google.com", []UrlPath{
				NewUrlPath("https://github.com", "https://github.com/user1"),
				NewUrlPath("https://github.com/user1", "https://github.com/user1/project1"),
				NewUrlPath("https://github.com/user1", "https://github.com/user1/project2"),
				NewUrlPath("https://github.com/user1", "https://github.com/user1/project3"),
				NewUrlPath("https://github.com/user2"),
				NewUrlPath("https://github.com", "https://github.com/user1", "/project4"),
				NewUrlPath("https://github.com", "https://github.com/user2", "project1"),
				NewUrlPath("https://github.com", "https://github.com/user2", "project2"),
				NewUrlPath("https://github.com", "https://github.com/user2", "project3"),
			})

			lines := []string{
				"digraph {",
				`"https://github.com/user1" -> "https://github.com/user1/project3";`,
				`"https://github.com/user1" -> "https://github.com/user1/project4";`,
				`"https://github.com/user1" -> "https://github.com/user1/project1";`,
				`"https://github.com/user1" -> "https://github.com/user1/project2";`,
				`"https://github.com/user2" -> "https://github.com/user2/project2";`,
				`"https://github.com/user2" -> "https://github.com/user2/project3";`,
				`"https://github.com/user2" -> "https://github.com/user2/project1";`,
				`"https://google.com" -> "https://github.com";`,
				`"https://google.com" -> "https://github.com/user1";`,
				`"https://google.com" -> "https://github.com/user2";`,
				`"https://github.com" -> "https://github.com/user1";`,
				`"https://github.com" -> "https://github.com/user2";`,
			}

			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			err := site.GenerateDOT(w)
			So(err, ShouldBeNil)

			err = w.Flush()
			So(err, ShouldBeNil)

			contents := string(buf.Bytes())
			for i := range lines {
				So(contents, ShouldContainSubstring, lines[i])
			}

			So(contents, ShouldNotContainSubstring, `""`)
		})
	})
}
