package main

import (
	"bufio"
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestGraphSpec(t *testing.T) {
	Convey("Given a file", t, func() {
		file1 := "graph.html"
		file2 := "vis.min.js"
		file3 := "vis.min.css"
		Convey("should have its name changed", func() {
			So(fixFileName(file1), ShouldEqual, "graphHTML")
			So(fixFileName(file2), ShouldEqual, "visMinJs")
			So(fixFileName(file3), ShouldEqual, "visMinCss")
		})

		Convey("should have its contents properly read", func() {
			cwd, _ := os.Getwd()
			defer os.Chdir(cwd)
			os.Chdir("..")
			contents, _ := readFile(file1)
			So(len(contents), ShouldBeGreaterThan, 0)
		})

		Convey("should have properly add en entry", func() {
			fileContents := []byte("Hello World")

			var outputBuffer bytes.Buffer
			out := bufio.NewWriter(&outputBuffer)

			writeNewItem(out, "helloWorld", fileContents)
			out.Flush()

			actual := string(outputBuffer.Bytes())

			So(actual, ShouldContainSubstring, "helloWorld = ")
		})
	})

	Convey("Given the contents of a file", t, func() {
		fileContents :=
			`<!DOCTYPE html>
<html lang="en">
  <head>
    <meta name=viewport content=width=device-width, initial-scale=1, shrink-to-fit=no>
	<meta name=theme-color content=#000000>
	<title></title>
  </head>
  <body>
  </body>
</html>
`
		Convey("should properly encode", func() {
			expected := `<!DOCTYPE html>\n<html lang=\"en\">\n  <head>\n    <meta name=viewport content=width=device-width, initial-scale=1, shrink-to-fit=no>\n	<meta name=theme-color content=#000000>\n	<title></title>\n  </head>\n  <body>\n  </body>\n</html>\n`
			actual := encode(fileContents)

			So(actual, ShouldEqual, expected)
			So(len(actual), ShouldEqual, len(expected))
		})
	})
}
