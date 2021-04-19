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
	var modeCount bool
	switch arg {
	case "dump":
		modeDump = true
	case "count":
		modeCount = true
	default:
		showUsage()
		return
	}

	text := extractText(specFile)
	//fmt.Print(text)

	// NLP
	doc, err := prose.NewDocument(text)
	if err != nil {
		panic(err)
	}

	var wordCount = map[string]int{}

	// Iterate over the gdoc's tokens:
	for _, tok := range doc.Tokens() {
		word := strings.ToLower(tok.Text)
		wordCount[word]++
		if modeDump {
			// Go NNP B-GPE
			// is VBZ O
			// an DT O
			// ...
			fmt.Println(tok.Text, tok.Tag, tok.Label)
		}
	}
	if modeCount {
		showCount(wordCount)
	}
}

func showCount(wordCount map[string]int) {
	type sameFreqGroup []string
	var frequency []sameFreqGroup = make([]sameFreqGroup, len(wordCount))
	for w, n := range wordCount {
		frequency[n] = append(frequency[n], w)
	}

	var total int
	for n, grp := range frequency {
		for _, w := range grp {
			first := []byte(w)[0]
			if first < 'a' || 'z' < first {
				// not alphabet
				continue
			}
			if len(w) <= 2 {
				// too short
				continue
			}
			total++
			fmt.Printf("%4d\t%s\n", n, w)
		}
	}
	fmt.Printf("%4d\tTotal\n", total)

}