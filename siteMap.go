package main

import (
	"fmt"
	"io"
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

func NewSiteMapFromUrlPaths(url string, paths []UrlPath) *SiteMap {
	site := NewSiteMap(url)
	for i := range paths {
		site.AddUrlPath(paths[i])
	}
	return site
}

func (sm *SiteMap) AddUrlPath(path UrlPath) {
	lastCur := sm.Root
	for i := 0; i < len(path); i++ {
		cur := sm.nodes[path[i]]
		if cur == nil {
			cur = NewNode(path[i])
			sm.nodes[path[i]] = cur
		}
		lastCur.addLink(cur)
		lastCur = cur
	}
}

func (sm *SiteMap) GetNode(url string) *Node {
	return sm.nodes[url]
}

// func (sm *SiteMap) Display(maxDepth int) {
// 	fmt.Println(sm.Href)
// 	sm.printNode(sm.Root, 1, maxDepth)
// }

// func (sm *SiteMap) printNode(node *Node, depth, maxDepth int) {
// 	if depth > maxDepth {
// 		return
// 	}
// 	for href := range node.links {
// 		for i := 0; i < depth; i++ {
// 			fmt.Print("    ")
// 		}
// 		n := node.links[href]
// 		fmt.Println(n.href)
// 		sm.printNode(n, depth+1, maxDepth)
// 	}
// }

// func (sm *SiteMap) Graph() {
// 	// GenerateGraphIndex(sm)
// }

func (sm *SiteMap) GenerateDOT(w io.Writer) (int, error) {
	count := 0
	var err error
	n, err := w.Write([]byte("digraph graphname {\n"))
	count += n
	if err != nil {
		return count, err
	}

	for _, node := range sm.nodes {
		count += n
		if err != nil {
			return count, err
		}
		for i := range node.links {
			src := node.href
			dst := node.links[i].href
			w.Write([]byte(fmt.Sprintf("    \"%s\" -> \"%s\";\n", src, dst)))
			count += n
			if err != nil {
				return count, err
			}
		}
	}
	n, err = w.Write([]byte("}"))
	count += n
	if err != nil {
		return count, err
	}
	return count, err
}
