package crawler

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

func (sm *SiteMap) AddUrlPath(path UrlPath) bool {
	addedNode := false
	lastCur := sm.Root
	for i := 0; i < len(path); i++ {
		cur := sm.nodes[path[i]]
		if cur == nil {
			cur = NewNode(path[i])
			addedNode = true
			sm.nodes[path[i]] = cur
		}
		lastCur.addLink(cur)
		lastCur = cur
	}
	return addedNode
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

func (sm *SiteMap) GenerateDOT(w io.Writer) error {
	_, err := w.Write([]byte("digraph {\n"))
	if err != nil {
		return err
	}

	for _, node := range sm.nodes {
		if err != nil {
			return err
		}
		for i := range node.links {
			src := node.href
			dst := node.links[i].href
			w.Write([]byte(fmt.Sprintf("    \"%s\" -> \"%s\";\n", src, dst)))
			if err != nil {
				return err
			}
		}
	}
	_, err = w.Write([]byte("}"))
	if err != nil {
		return err
	}
	return err
}

func (site *SiteMap) GenerateIndexHtml(w io.Writer) error {
	graph := newGraphHtmlTemplate(site)
	return writeGraphIndex(w, graph)
}
