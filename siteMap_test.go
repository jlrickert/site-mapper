package main

import (
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
		urls := []*Url{
			NewUrlFromSlice([]string{"a"}),
			NewUrlFromSlice([]string{"a", "a"}),
			NewUrlFromSlice([]string{"a", "b"}),
			NewUrlFromSlice([]string{"a", "b", "c"}),
			NewUrlFromSlice([]string{"a", "b", "c", "a"}),
			NewUrlFromSlice([]string{"a", "b", "d"}),
			NewUrlFromSlice([]string{"a", "b", "c", "d"}),
			NewUrlFromSlice([]string{"a", "b", "c", "c"}),
			NewUrlFromSlice([]string{"a", "b", "d", "c"}),
		}
		_ = urls
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

			site := NewSiteMapFromSlice("a", urls)

			urlMap := make(map[string]*Node)
			urlMap["a"] = nodeA
			urlMap["b"] = nodeB
			urlMap["c"] = nodeC
			urlMap["d"] = nodeD

			/////////
			// Node A
			/////////
			actualNodeA := site.GetNode("a")
			So(actualNodeA.links["a"].href, ShouldEqual, nodeA.href)
			So(actualNodeA.links["b"].href, ShouldEqual, nodeB.href)
			So(actualNodeA.links["c"], ShouldEqual, nil)
			So(actualNodeA.links["d"], ShouldEqual, nil)

			So(actualNodeA.linkedFrom["a"].href, ShouldEqual, nodeA.href)
			So(actualNodeA.linkedFrom["b"], ShouldEqual, nil)
			So(actualNodeA.linkedFrom["c"].href, ShouldEqual, nodeC.href)
			So(actualNodeA.linkedFrom["d"], ShouldEqual, nil)

			/////////
			// Node B
			/////////
			actualNodeB := site.GetNode("b")
			So(actualNodeB.links["a"], ShouldEqual, nil)
			So(actualNodeB.links["b"], ShouldEqual, nil)
			So(actualNodeB.links["c"].href, ShouldEqual, nodeC.href)
			So(actualNodeB.links["d"].href, ShouldEqual, nodeD.href)

			So(actualNodeB.linkedFrom["a"].href, ShouldEqual, nodeA.href)
			So(actualNodeB.linkedFrom["b"], ShouldEqual, nil)
			So(actualNodeB.linkedFrom["c"], ShouldEqual, nil)
			So(actualNodeB.linkedFrom["d"], ShouldEqual, nil)

			/////////
			// Node C
			/////////
			actualNodeC := site.GetNode("c")
			So(actualNodeC.links["a"].href, ShouldEqual, nodeA.href)
			So(actualNodeC.links["b"], ShouldEqual, nil)
			So(actualNodeC.links["c"].href, ShouldEqual, nodeC.href)
			So(actualNodeC.links["d"].href, ShouldEqual, nodeD.href)

			So(actualNodeC.linkedFrom["a"], ShouldEqual, nil)
			So(actualNodeC.linkedFrom["b"].href, ShouldEqual, nodeB.href)
			So(actualNodeC.linkedFrom["c"].href, ShouldEqual, nodeC.href)
			So(actualNodeC.linkedFrom["d"].href, ShouldEqual, nodeD.href)

			/////////
			// Node D
			/////////
			actualNodeD := site.GetNode("d")
			So(actualNodeD.links["a"], ShouldEqual, nil)
			So(actualNodeD.links["b"], ShouldEqual, nil)
			So(actualNodeD.links["c"].href, ShouldEqual, nodeC.href)
			So(actualNodeD.links["d"], ShouldEqual, nil)

			So(actualNodeD.linkedFrom["a"], ShouldEqual, nil)
			So(actualNodeD.linkedFrom["b"].href, ShouldEqual, nodeB.href)
			So(actualNodeD.linkedFrom["c"].href, ShouldEqual, nodeC.href)
			So(actualNodeD.linkedFrom["d"], ShouldEqual, nil)
		})
	})
}
