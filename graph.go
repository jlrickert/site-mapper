package main

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

type GraphHtmlTemplate struct {
	Title   string
	Styles  []string
	Scripts []string
	MainJs  string
}

func newGraphHtmlTemplate(site *SiteMap) *GraphHtmlTemplate {
	dot := bytes.NewBuffer([]byte{})
	site.GenerateDOT(dot)
	return &GraphHtmlTemplate{
		Title:   site.Href + " Sitemap Graph",
		Styles:  []string{visMinCss},
		Scripts: []string{visMinJs},
		MainJs:  createMainJs(dot.Bytes()),
	}
}

func writeGraphIndex(out io.Writer, graph *GraphHtmlTemplate) error {
	t := template.Must(template.New(graph.Title).Parse(graphHTML))
	return t.Execute(out, graph)
}

func createMainJs(dot []byte) string {
	backtick := "`"
	javascript := fmt.Sprintf(`const dot = %s%s%s;
const container = document.getElementById("network")
const data = vis.network.convertDot(dot);
console.log(container);
const network = new vis.Network(container, data);
`, backtick, dot, backtick)
	return javascript
}
