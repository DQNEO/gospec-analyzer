package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jdkato/prose/v2"
	"github.com/jinzhu/inflection"
	"github.com/surgebase/porter2"
)

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

	var modeDump, modeCount, modeUniq bool
	switch arg {
	case "dump":
		modeDump = true
	case "count":
		modeCount = true
	case "uniq":
		modeUniq = true
	default:
		showUsage()
		return
	}
	tsv, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(tsv), "\n")
	var tokens []prose.Token
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

	var meaningfulTokens []prose.Token
	// exclude meaningless tokens
	for _, tok := range tokens {
		if len(tok.Text) == 1 { // exclue one letter
			continue
		}

		// Exclude tokens with punctuations
		first := tok.Text[0]
		if !(('a' <= first && first <= 'z') || ('A' <= first && first <= 'Z')) {
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
			"CD",   // cardinal number
			"DT",   // determiner
			"IN",   // conjunction, subordinating or preposition
			"CC",   // conjunction, coordinating
			"PRP",  // pronoun, personal
			"PRP$", // pronoun, possessive
			"TO",   // infinitival to
			"WDT",  // wh-determiner
			"WP",   // wh-pronoun, personal
			"WP$",  // wh-pronoun, possessive
			"WRB",  // wh-adverb
			"MD":   // verb, modal auxiliary
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

	var importantTokens []prose.Token
	for _, tok := range meaningfulTokens {
		switch tok.Text {
		case
		"is", "are", "be","has",
		"not","same", "only",
		"use",
		"name",
		"first", "one","zero":
			continue
		default:
			importantTokens = append(importantTokens, tok)
		}
	}
	if modeCount {
		countByTags(meaningfulTokens)
		return
	}
	if modeUniq {
		countByWord(meaningfulTokens)
		return
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
