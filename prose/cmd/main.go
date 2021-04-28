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
	tokens := prose.Tokenize(fname)
	for _, tok := range tokens {
		fmt.Printf("%s\t%s\t%s\n", tok.Text, tok.Tag, tok.Label)
	}
}
