package main

import (
	"fmt"
)

type SiteMap struct {
	Href  string
	Root  *Node
	nodes map[string]*Node
}

type Node struct {
	href       string
	links      map[string]*Node
	linkedFrom map[string]*Node
}

func NewNode(url string) *Node {
	return &Node{
		href:       url,
		links:      make(map[string]*Node),
		linkedFrom: make(map[string]*Node),
	}
}

func (this *Node) addLink(that *Node) {
	this.links[that.href] = that
	that.linkedFrom[this.href] = this
}

func NewSiteMap(rootUrl string) *SiteMap {
	rootNode := NewNode(rootUrl)

	nodes := make(map[string]*Node)
	nodes[rootUrl] = rootNode

	site := SiteMap{
		Href:  rootUrl,
		Root:  rootNode,
		nodes: nodes,
	}

	return &site
}

func NewSiteMapFromSlice(rootUrl string, urls []*Url) (site *SiteMap) {
	site = NewSiteMap(rootUrl)
	for i := range urls {
		site.AddUrl(urls[i])
	}
	return
}

func (sm *SiteMap) AddUrl(url *Url) {
	urls := url.Path[:]
	urls = append(urls, url.Href)

	lastCur := sm.Root
	for i := 0; i < len(urls); i++ {
		cur := sm.nodes[urls[i]]
		if cur == nil {
			cur = NewNode(urls[i])
			sm.nodes[urls[i]] = cur
		}
		lastCur.addLink(cur)
		lastCur = cur
	}
}

func (sm *SiteMap) GetNode(url string) *Node {
	return sm.nodes[url]
}

func (sm *SiteMap) Display(maxDepth int) {
	fmt.Println(sm.Href)
	sm.printNode(sm.Root, 1, maxDepth)
}

func (sm *SiteMap) printNode(node *Node, depth, maxDepth int) {
	if depth > maxDepth {
		return
	}
	for href := range node.links {
		for i := 0; i < depth; i++ {
			fmt.Print("    ")
		}
		n := node.links[href]
		fmt.Println(n.href)
		sm.printNode(n, depth+1, maxDepth)
	}
}

func (sm *SiteMap) Graph() {
	// GenerateGraphIndex(sm)
}
