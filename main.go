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
	var meaningfulTokens []prose.Token
	// exclude meaningless tokens
	for _, tok := range tokens {
		if len(tok.Text) == 1 { // exclue one letter
			continue
		}


		// Exclude tokens with punctuations
		first := tok.Text[0]
		if ! (('a' <= first && first <= 'z') || ('A' <= first && first <= 'Z')) {
			continue
		}
		if strings.Contains(tok.Text, "[") {
			continue
		}
		if strings.Contains(tok.Text, "(") {
			continue
		}
		if strings.Contains(tok.Text, "+") {
			continue
		}

		// Exclude tokens of DT (a,an,the,..)
		switch tok.Tag {
		case
			"DT", // determiner
			"IN", // conjunction, subordinating or preposition
			"CC", // conjunction, coordinating
			"PRP", // pronoun, personal
			"PRP$", // pronoun, possessive
			"WDT", // wh-determiner
			"MD": // verb, modal auxiliary
			continue
		default:
			meaningfulTokens = append(meaningfulTokens, tok)
		}
	}
	for _, tok := range meaningfulTokens {
		if modeDump {
			// Go NNP B-GPE
			// is VBZ O
			// an DT O
			// ...
			//fmt.Printf("%20s %5s %10s\n", tok.Text, tok.Tag, tok.Label)
			lowerText := strings.ToLower(tok.Text)
			fmt.Printf("%s\t%s\t%s\n", lowerText, tok.Tag, tok.Label)
		}
	}
	if modeCount {
		showCount(meaningfulTokens)
	}
}

func showCount(meaningfulTokens []prose.Token) {
	var wordCount = map[string]int{}
	for _, tok := range meaningfulTokens {
		lowerText := strings.ToLower(tok.Text)
		wordCount[lowerText]++
	}

	type sameFreqGroup []string
	var frequency []sameFreqGroup = make([]sameFreqGroup, len(wordCount) * 2)
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