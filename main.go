package main

import (
	"fmt"
	"github.com/DQNEO/go-samples/nlp/gospec/spec2text"
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
	f, err := os.Open(specFile)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	text := spec2text.GetTextFromHTML(f)
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
		aggregate(meaningfulTokens)
	}
}

func explainConversion(old *prose.Token, new *prose.Token) {
	fmt.Fprintf(os.Stderr, "Converting [%4s] %20s => [%4s] %20s\n",
		old.Tag, old.Text, new.Tag, new.Text)
}

func singulifyToken(origTok *prose.Token) *prose.Token {
	if strings.HasSuffix(origTok.Text, "s") {
		tok := &prose.Token{}
		tok.Tag = "NN"
		if strings.HasSuffix(origTok.Text, "ies") {
			// "entries" => "entry"
			tok.Text = strings.TrimSuffix(origTok.Text, "ies") + "y"
			explainConversion(origTok, tok)
		} else if strings.HasSuffix(origTok.Text, "es") {
			// @TODO "resumes" => "resume"
			// "classes" => "class"
			tok.Text = strings.TrimSuffix(origTok.Text, "es")
			explainConversion(origTok, tok)
		} else {
			tok.Text = strings.TrimSuffix(origTok.Text, "s")
			explainConversion(origTok, tok)
		}
		return tok
	}
	return origTok
}

func manipulateToken(origTok prose.Token) prose.Token {
	// Manipulate token
	tok := origTok
	switch origTok.Tag {
	case "NNS", "NNPS":
		// dogs NNS -> dog NN

		switch origTok.Text {
		case "halves":
			tok.Tag = "NN"
			tok.Text = "half"
			explainConversion(&origTok, &tok)
			return tok
		case "applies":
			tok.Tag = "NN"
			tok.Text = "apply"
			explainConversion(&origTok, &tok)
			return tok
		}
		return *singulifyToken(&origTok)
	case "VBZ":
		switch origTok.Text {
		case "is", "has":
			return tok
		}
		// "works" VBZ -> "work" VB
		if strings.HasSuffix(origTok.Text, "s") {
			tok.Tag = "VB"
			tok.Text = strings.TrimSuffix(origTok.Text, "s")
			explainConversion(&origTok, &tok)
			return tok
		}

	}
	return tok
}

func aggregate(meaningfulTokens []prose.Token) {
	var wordCount = map[string]map[string]int{}
	for _, origTok := range meaningfulTokens {
		tok := manipulateToken(origTok)
		lowerText := strings.ToLower(tok.Text)
		_, ok := wordCount[lowerText]
		if !ok {
			wordCount[lowerText] = make(map[string]int, 100)
		}
		wordCount[lowerText][tok.Tag]++
	}

	type sameFreqGroup []string
	var frequency []sameFreqGroup = make([]sameFreqGroup, 10000)
	for text, tagCount := range wordCount {
		for tag, cnt := range tagCount {
			frequency[cnt] = append(frequency[cnt], text + "\t" + tag)
		}
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