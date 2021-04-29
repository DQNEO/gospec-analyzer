package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jdkato/prose/v2"
	"github.com/jinzhu/inflection"
	"github.com/surgebase/porter2"
	"github.com/DQNEO/gospec-analyzer/tokenizer"
)

func showUsage() {
	help := `
Usage:
	filter1: exclude meaningful tokens
	count: show statistics
	uniq: show statistics uniq by word
`
	fmt.Print(help)
}

func main() {
	if len(os.Args) == 1 {
		showUsage()
		return
	}
	arg := os.Args[1]

	var modeCount, modeUniq, modeFilter1 bool
	switch arg {
	case "filter1":
		modeFilter1 = true
	case "count":
		modeCount = true
	case "uniq":
		modeUniq = true
	default:
		showUsage()
		return
	}

	tokens := loadTokens(os.Stdin)
	if modeFilter1 {
		filter1(tokens)
		return
	}
	if modeCount {
		countByTags(tokens)
		return
	}
	if modeUniq {
		countByWord(tokens)
		return
	}
}

func filter1(tokens []prose.Token) {
	var meaningfulTokens []prose.Token
	// exclude meaningless tokens
	for _, tok := range tokens {
		if isMeaningful(&tok) {
			meaningfulTokens = append(meaningfulTokens, tok)
		} else {
			println("skip: " + tok.Text)
		}
	}

	// exclude basic words
	for _, tok := range meaningfulTokens {
		if basicWords[tok.Text] {
			// log
		} else {
			fmt.Println(tokenizer.String(&tok))
		}
	}
}

func loadTokens(r io.Reader) []prose.Token {
	tsv, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(tsv), "\n")
	var tokens []prose.Token = make([]prose.Token, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 { // EOF
			break
		}
		fields := strings.Split(line, "\t")
		tk := prose.Token{
			Text:   fields[0],
			Tag:  fields[1],
			Label: fields[2],
		}
		tokens = append(tokens, tk)
	}
	return tokens
}

func isMeaningful(tok *prose.Token) bool {

	if len(tok.Text) == 1 { // exclude one letter token
		return false
	}

	// Exclude tokens with punctuations
	first := tok.Text[0]
	if !(('a' <= first && first <= 'z') || ('A' <= first && first <= 'Z')) {
		return false
	}
	if strings.Contains(tok.Text, "[") {
		return false
	}
	if strings.Contains(tok.Text, "(") {
		return false
	}
	if strings.Contains(tok.Text, "+") {
		return false
	}

	// Exclude tokens of DT (a,an,the,..)
	return !meaninglessTokens[tok.Tag]
}

func explainConversion(old *prose.Token, new *prose.Token) {
	fmt.Fprintf(os.Stderr, "Converting [%4s] %20s => [%4s] %20s\n",
		old.Tag, old.Text, new.Tag, new.Text)
}

func singulifyToken(origTok *prose.Token) *prose.Token {
	if strings.HasSuffix(origTok.Text, "s") {
		tok := &prose.Token{}
		tok.Tag = "NN"
		tok.Text = inflection.Singular(origTok.Text)
		explainConversion(origTok, tok)
		return tok
	} else {
		return origTok
	}
}

func manipulateToken(origTok prose.Token) prose.Token {
	// Manipulate token
	tok := origTok
	switch origTok.Tag {
	case "NNS", "NNPS":
		// dogs NNS -> dog NN
		return *singulifyToken(&origTok)
	case "VBD", "VBG", "VBN", "VBP", "VBZ":
		tok.Text = porter2.Stem(origTok.Text)
		tok.Tag = "VB"
		return tok
	case "JJS":
		if strings.HasSuffix(origTok.Text, "est") {
			tok.Text = strings.TrimSuffix(origTok.Text, "est")
			tok.Tag = "JJ"
			return tok
		}
	}
	return tok
}

func countByTags(meaningfulTokens []prose.Token) {
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
			frequency[cnt] = append(frequency[cnt], text+"\t"+tag)
		}
	}

	for n, grp := range frequency {
		for _, w := range grp {
			fmt.Printf("%4d\t%s\n", n, w)
		}
	}
}

func countByWord(meaningfulTokens []prose.Token) {
	var wordCount = map[string]int{}
	for _, origTok := range meaningfulTokens {
		tok := manipulateToken(origTok)
		lowerText := strings.ToLower(tok.Text)
		cnt := wordCount[lowerText]
		wordCount[lowerText] = cnt + 1
	}

	frequency := make([][]string, 10000)
	for text, cnt := range wordCount {
		frequency[cnt] = append(frequency[cnt], text)
	}

	for cnt, grp := range frequency {
		if grp == nil {
			continue
		}
		for _, w := range grp {
			fmt.Printf("%4d\t%s\n", cnt, w)
		}
	}
}
