package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	fs, _ := ioutil.ReadDir("./resources")
	out, _ := os.Create("resources.go")
	out.Write([]byte("package main\n\nconst (\n"))
	for _, f := range fs {
		out.Write([]byte("\t" + fixFileName(f.Name()) + " = \""))
		f, err := os.Open("./resources/" + f.Name())
		if err != nil {
			panic(err)
		}

		content := []byte{}
		_, _ = f.Read(content)

		r := encode(f)
		io.Copy(out, r)

		out.Write([]byte("\"\n"))
	}
	out.Write([]byte((")\n")))
}

func fixFileName(filename string) string {
	filename = strings.Replace(filename, "html", "HTML", -1)

	parts := strings.Split(filename, ".")
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func encode(r io.Reader) io.Reader {
	buffer := make([]byte, 256)
	content := bytes.NewBuffer([]byte{})
	for {
		n, _ := r.Read(buffer)
		if n <= 0 {
			break
		}
		encode := string(buffer)
		encode = strings.Replace(encode, "\\", "\\\\", -1)
		encode = strings.Replace(encode, "\"", "\\\"", -1)
		encode = strings.Replace(encode, "\n", "\\n", -1)
		content.Write([]byte(encode))
	}
	return content
}
