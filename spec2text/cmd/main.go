package main

import (
	"fmt"
	"os"

	"github.com/DQNEO/go-samples/nlp/gospec/spec2text"
)

func usage() {
	fmt.Println("s2t is a tool which converts a spec html file into text")
	fmt.Println("Usage:")
	fmt.Println("s2t spec.html > spec.txt")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	var htmlFile = os.Args[1]
	f, err := os.Open(htmlFile)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	text := spec2text.GetTextFromHTML(f)
	fmt.Print(text)
}
