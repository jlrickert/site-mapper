package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fs, _ := ioutil.ReadDir("resources")
	out, _ := os.Create("resources.go")
	out.Write([]byte("package crawler\n\nconst (\n"))
	for _, f := range fs {
		rawContent, err := readFile(f.Name())
		if err != nil {
			panic(err)
		}

		writeNewItem(out, fixFileName(f.Name()), rawContent)
	}
	out.Write([]byte((")\n")))
}

func readFile(filename string) ([]byte, error) {
	fullPath, err := filepath.Abs(filepath.Join("resources", filename))
	if err != nil {
		return []byte{}, err
	}
	return ioutil.ReadFile(fullPath)
}

func writeNewItem(out io.Writer, varName string, b []byte) (int, error) {
	contents := encode(string(b))
	return out.Write([]byte("\t" + varName + ` = "` + contents + "\"\n"))
}

func fixFileName(filename string) string {
	filename = strings.Replace(filename, "html", "HTML", -1)

	parts := strings.Split(filename, ".")
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func encode(str string) string {
	str = strings.Replace(str, `\`, `\\`, -1)
	str = strings.Replace(str, `"`, `\"`, -1)
	str = strings.Replace(str, "\n", `\n`, -1)
	return str
}
