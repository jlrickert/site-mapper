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
const container = document.getElementById("network");
const data = vis.network.convertDot(dot);
const options = {
  configure: {
	enabled: true,
	container: document.getElementById('config'),
	filter: (option, path) => {
	  if (path.indexOf('physics') !== -1) {
		return true;
	  }
	  if (path.indexOf('smooth') !== -1 || option === 'smooth') {
		return true;
	  }
	  return false;
	},
  },
  edges: {
    smooth: true
  },
  physics: {
    barnesHut: {
      gravitationalConstant: -20000,
      centralGravity: 1,
      springLength: 10 * data["nodes"].length,
      springConstant: 0.10,
      damping: 0.25,
      avoidOverlap: 0.25
    }
  }
};
const network = new vis.Network(container, data, options);
`, backtick, dot, backtick)
	return javascript
}
