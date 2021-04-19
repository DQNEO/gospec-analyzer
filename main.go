package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jdkato/prose"
)

var specFile = os.Getenv("GOPATH") + "/src/github.com/DQNEO/go-samples/nlp/gospec/spec.html"

func showUsage() {
	help := `
Usage:
	text:  get raw text by parsing HTML
	dump : show tokens
	count: show statistics
`
	fmt.Print(help)
}
func main() {
	if len(os.Args) == 1 {
		showUsage()
		return
	}
	arg := os.Args[1]

	var modeDump bool
	var modeText bool
	var modeCount bool
	switch arg {
	case "text":
		modeText = true
	case "dump":
		modeDump = true
	case "count":
		modeCount = true
	default:
		showUsage()
		return
	}

	text := extractText(specFile)
	if modeText {
		fmt.Print(text)
		return
	}

	// NLP
	doc, err := prose.NewDocument(text)
	if err != nil {
		panic(err)
	}

	// Iterate over the gdoc's tokens:
	tokens := doc.Tokens()
	for _, tok := range tokens {
		if modeDump {
			// Go NNP B-GPE
			// is VBZ O
			// an DT O
			// ...
			fmt.Printf("%20s %5s %10s\n", tok.Text, tok.Tag, tok.Label)
		}
	}
	if modeCount {
		showCount(tokens)
	}
}

func showCount(tokens []prose.Token) {
	var wordCount = map[string]int{}
	for _, tok := range tokens {
		word := strings.ToLower(tok.Text)
		wordCount[word]++
	}

	type sameFreqGroup []string
	var frequency []sameFreqGroup = make([]sameFreqGroup, len(wordCount))
	for w, n := range wordCount {
		frequency[n] = append(frequency[n], w)
	}

	var total int
	for n, grp := range frequency {
		for _, w := range grp {
			total++
			fmt.Printf("%4d\t%s\n", n, w)
		}
	}
	fmt.Printf("%4d\tTotal\n", total)

}