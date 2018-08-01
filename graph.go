package main

// import (
// 	"io"
// 	"log"
// 	"text/template"
// )

// type GraphHtml struct {
// 	Title   string
// 	Styles  []string
// 	Scripts []string
// }

// type Javascript struct {
// }

// func NewGraphHtml(url string) *GraphHtml {
// 	return &GraphHtml{
// 		Title:   url + " Sitemap Graph",
// 		Styles:  []string{visMinCss},
// 		Scripts: []string{visMinJs},
// 	}
// }

// func WriteGraphIndex(out io.Writer, site *SiteMap) {
// 	t := template.Must(template.New("graph").Parse(graphHTML))

// 	graph := NewGraphHtml(site.Href)
// 	err := t.Execute(out, graph)
// 	if err != nil {
// 		log.Println("executing template:", err)
// 	}
// }
