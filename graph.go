package main

import (
	"log"
	"os"
	"text/template"
)

type GraphIndex struct {
	Title   string
	Css     []string
	Scripts []string
}

func NewGraphIndex(url string) *GraphIndex {
	return &GraphIndex{
		Title:   url + " Sitemap Graph",
		Css:     []string{visMinCss},
		Scripts: []string{visMinJs},
	}
}

func GenerateGraphIndex(site *SiteMap) {
	graph := NewGraphIndex(site.Href)
	t := template.Must(template.New("graph").Parse(graphHTML))
	err := t.Execute(os.Stdout, graph)
	if err != nil {
		log.Println("executing template:", err)
	}
	return
}
