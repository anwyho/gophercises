package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/anwyho/gophercises/link"
)

var exampleHtml = `
<html>
<body>
  <h1>Hello!</hw>
  <a href="/other-page">
	  A link to another page
	  <span> some span </span>
  </a>
  <a href="/page-two">A link to a second page</a>
</body>
</html>
`

func main() {
	readDir := "./examples/"
	files, err := ioutil.ReadDir(readDir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
		r, err := os.Open(filepath.Join(readDir, f.Name()))
		if err != nil {
			panic(err)
		}
		links, err := link.Parse(r)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", links)
	}
}
