package main

import (
	"encoding/json"
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
	filter1: exclude punctuations
	filter2: exclude meaningless tokens
	filter3: exclude basic words
	filter4: exclude technical terms
	normalize: normalize word variations
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

	switch arg {
	case "filter1":
		filter(loadTokens(os.Stdin), isPunct)
	case "filter2":
		filter(loadTokens(os.Stdin), isMeaningless)
	case "filter3":
		filter(loadTokens(os.Stdin), isBasicWord)
	case "filter4":
		filter(loadTokens(os.Stdin), isTechnicalTerm)
	case "normalize":
		normalize(loadTokens(os.Stdin))
	case "normalizejson":
		normalizeJson(loadTokens(os.Stdin))
	case "count":
		tokens := loadTokens(os.Stdin)
		countByTags(tokens)
	case "uniq":
		tokens := loadTokens(os.Stdin)
		countByWord(tokens)
	default:
		showUsage()
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

// exclude unneeded things
func filter(tokens []prose.Token, fn func(*prose.Token) bool) {
	for _, tok := range tokens {
		if fn(&tok) {
			println("skip: " + tokenizer.String(&tok))
		} else {
			fmt.Println(tokenizer.String(&tok))
		}
	}
}

func isPunct(tok *prose.Token) bool {

	if len(tok.Text) == 1 { // exclude one letter token
		return true
	}

	// Exclude tokens with punctuations
	first := tok.Text[0]
	if !(('a' <= first && first <= 'z') || ('A' <= first && first <= 'Z')) {
		return true
	}
	if strings.Contains(tok.Text, "[") {
		return true
	}
	if strings.Contains(tok.Text, "(") {
		return true
	}
	if strings.Contains(tok.Text, "+") {
		return true
	}

	return false
}

func isMeaningless(tok *prose.Token) bool {
	// Exclude tokens of DT (a,an,the,..)
	return meaninglessTokens[strings.ToUpper(tok.Tag)]
}

func isBasicWord(tok *prose.Token) bool {
	return basicWords[tok.Text]
}

func isTechnicalTerm(tok *prose.Token) bool {
	return technicalTerm[tok.Text]
}

func normalize(tokens []prose.Token) {
	for _, tok := range tokens {
		newTok := manipulateToken(tok)
		fmt.Printf(
			"%-15s\t%4s => %-15s\t%4s\t:%s\n",
			tok.Text, tok.Tag,
			newTok.Text, newTok.Tag,
			tok.Label,
		)
	}
}

func normalizeJson(tokens []prose.Token) {
	var assoc map[string]string = make(map[string]string, len(tokens))
	for _, tok := range tokens {
		newTok := manipulateToken(tok)
		assoc[tok.Text] = newTok.Text
	}
	bytes , err := json.MarshalIndent(assoc, "", "  ")
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(bytes)
}

func explainConversion(old *prose.Token, new *
	prose.Token) {
	fmt.Fprintf(os.Stderr, "Converting [%4s] %20s => [%4s] %20s\n",
		old.Tag, old.Text, new.Tag, new.Text)
}

func singulifyToken(origTok *prose.Token) *prose.Token {
	if strings.HasSuffix(origTok.Text, "s") {
		tok := &prose.Token{}
		tok.Tag = "nn"
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
	case "nns", "nnps":
		// dogs NNS -> dog NN
		return *singulifyToken(&origTok)
	case "vbd", "vbg", "vbn", "vbp", "vbz":
		tok.Text = porter2.Stem(origTok.Text)
		tok.Tag = "vb"
		return tok
	case "jjs":
		if strings.HasSuffix(origTok.Text, "est") {
			tok.Text = strings.TrimSuffix(origTok.Text, "est")
			tok.Tag = "jj"
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
