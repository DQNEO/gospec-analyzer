package main

import (
	"fmt"
	"os"

	"github.com/DQNEO/gospec-analyzer/prose"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Println("	prs textfile > tokens.tsv")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	fname := os.Args[1]
	prose.Tokenize(fname)
}
