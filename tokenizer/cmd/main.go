package main

import (
	"fmt"
	"os"
	"encoding/json"
	"github.com/DQNEO/gospec-analyzer/tokenizer"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Println("	tokenizer textfile > tokens.tsv")
	fmt.Println("	tokenizer --json textfile > tokens.json")
}

type mode int
const ModeTsv mode = 0
const ModeJson mode = 1

var outputMode mode

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	if os.Args[1] == "--json" {
		outputMode = ModeJson
		os.Args = os.Args[1:]
	}

	fname := os.Args[1]
	tokens := tokenizer.Tokenize(fname)

	switch outputMode {
	case ModeTsv:
		for _, tok := range tokens {
			fmt.Println(tokenizer.String(&tok))
		}
	case ModeJson:
		b, err := json.Marshal(tokens)
		if err != nil {
			panic(err)
		}
		fmt.Print(string(b))
	}
}
