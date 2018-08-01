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
		url := "http://github.com"

		Convey("should correctly create sitemap href", func() {
			site := NewSiteMap(url)
			So(site.Href, ShouldEqual, url)
		})
		Convey("should correctly set root node", func() {
			url := "http://github.com"
			site := NewSiteMap(url)
			So(site.Root.href, ShouldEqual, url)
		})
	})

	Convey("Given a set of urls from a website", t, func() {
		paths := []UrlPath{}
		_ = paths
		urls := []UrlPath{
			NewUrlPath("a", "a"),
			NewUrlPath("a", "b"),
			NewUrlPath("a", "b", "c"),
			NewUrlPath("a", "b", "c", "a"),
			NewUrlPath("a", "b", "d"),
			NewUrlPath("a", "b", "c", "d"),
			NewUrlPath("a", "b", "c", "c"),
			NewUrlPath("a", "b", "d", "c"),
		}

		Convey("should correctly build nodes", func() {
			nodeA := NewNode("a")
			nodeB := NewNode("b")
			nodeC := NewNode("c")
			nodeD := NewNode("d")

			nodeA.addLink(nodeB)
			nodeA.addLink(nodeA)

			nodeB.addLink(nodeC)
			nodeB.addLink(nodeD)

			nodeC.addLink(nodeA)
			nodeC.addLink(nodeD)

			nodeD.addLink(nodeC)

			site := NewSiteMapFromUrlPaths("a", urls)

			/////////
			// Node A
			/////////
			actualNodeA := site.GetNode("a")
			So(actualNodeA.links["a"].href, ShouldEqual, nodeA.href)
			So(actualNodeA.links["b"].href, ShouldEqual, nodeB.href)
			So(actualNodeA.links["c"], ShouldEqual, nil)
			So(actualNodeA.links["d"], ShouldEqual, nil)
			So(len(actualNodeA.links), ShouldEqual, 2)

			So(actualNodeA.linkedFrom["a"].href, ShouldEqual, nodeA.href)
			So(actualNodeA.linkedFrom["b"], ShouldEqual, nil)
			So(actualNodeA.linkedFrom["c"].href, ShouldEqual, nodeC.href)
			So(actualNodeA.linkedFrom["d"], ShouldEqual, nil)
			So(len(actualNodeA.linkedFrom), ShouldEqual, 2)

			/////////
			// Node B
			/////////
			actualNodeB := site.GetNode("b")
			So(actualNodeB.links["a"], ShouldEqual, nil)
			So(actualNodeB.links["b"], ShouldEqual, nil)
			So(actualNodeB.links["c"].href, ShouldEqual, nodeC.href)
			So(actualNodeB.links["d"].href, ShouldEqual, nodeD.href)
			So(len(actualNodeB.links), ShouldEqual, 2)

			So(actualNodeB.linkedFrom["a"].href, ShouldEqual, nodeA.href)
			So(actualNodeB.linkedFrom["b"], ShouldEqual, nil)
			So(actualNodeB.linkedFrom["c"], ShouldEqual, nil)
			So(actualNodeB.linkedFrom["d"], ShouldEqual, nil)
			So(len(actualNodeB.linkedFrom), ShouldEqual, 1)

			/////////
			// Node C
			/////////
			actualNodeC := site.GetNode("c")
			So(actualNodeC.links["a"].href, ShouldEqual, nodeA.href)
			So(actualNodeC.links["b"], ShouldEqual, nil)
			So(actualNodeC.links["c"].href, ShouldEqual, nodeC.href)
			So(actualNodeC.links["d"].href, ShouldEqual, nodeD.href)
			So(len(actualNodeC.links), ShouldEqual, 3)

			So(actualNodeC.linkedFrom["a"], ShouldEqual, nil)
			So(actualNodeC.linkedFrom["b"].href, ShouldEqual, nodeB.href)
			So(actualNodeC.linkedFrom["c"].href, ShouldEqual, nodeC.href)
			So(actualNodeC.linkedFrom["d"].href, ShouldEqual, nodeD.href)
			So(len(actualNodeC.linkedFrom), ShouldEqual, 3)

			/////////
			// Node D
			/////////
			actualNodeD := site.GetNode("d")
			So(actualNodeD.links["a"], ShouldEqual, nil)
			So(actualNodeD.links["b"], ShouldEqual, nil)
			So(actualNodeD.links["c"].href, ShouldEqual, nodeC.href)
			So(actualNodeD.links["d"], ShouldEqual, nil)
			So(len(actualNodeD.links), ShouldEqual, 1)

			So(actualNodeD.linkedFrom["a"], ShouldEqual, nil)
			So(actualNodeD.linkedFrom["b"].href, ShouldEqual, nodeB.href)
			So(actualNodeD.linkedFrom["c"].href, ShouldEqual, nodeC.href)
			So(actualNodeD.linkedFrom["d"], ShouldEqual, nil)
			So(len(actualNodeD.linkedFrom), ShouldEqual, 2)
		})
	})

	Convey("Given a valid sitemap", t, func() {
		Convey("should be able translate to DOT format", func() {
			site := NewSiteMapFromUrlPaths("google.com", []UrlPath{
				NewUrlPath("google.com"),
				NewUrlPath("duckduckgo.com", "google.com"),
				NewUrlPath("github.com", "github.com/user"),
			})
			expected := `digraph graphname {
    "google.com" -> "google.com";
    "google.com" -> "duckduckgo.com";
    "google.com" -> "github.com";
    "duckduckgo.com" -> "google.com";
    "github.com" -> "github.com/user";
}`

			var contents bytes.Buffer
			w := bufio.NewWriter(&contents)
			_, err := site.GenerateDOT(w)
			So(err, ShouldBeNil)

			err = w.Flush()
			So(err, ShouldBeNil)

			So(string(contents.Bytes()), ShouldEqual, expected)
		})
	})
}
