// Usage:
//   go run parseHTML.go
package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jdkato/prose"
	"os"
	"strings"
)

func main() {
	specFile := os.Getenv("GOPATH") + "/src/github.com/DQNEO/go-samples/nlp/gospec/spec.html"
	f, err := os.Open(specFile)
	if err != nil {
		panic(err)
	}
	gdoc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		panic(err)
	}
	gdoc.Find("pre").Remove()
	noCodeText := gdoc.Text()

	// NLP
	doc, err := prose.NewDocument(noCodeText)
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
