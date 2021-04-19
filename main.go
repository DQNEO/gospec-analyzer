// Usage:
//   go run main.go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jdkato/prose"
)

var specFile = os.Getenv("GOPATH") + "/src/github.com/DQNEO/go-samples/nlp/gospec/spec.html"

func main() {
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
		//fmt.Println(tok.Text, tok.Tag, tok.Label)
		// Go NNP B-GPE
		// is VBZ O
		// an DT O
		// ...
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
