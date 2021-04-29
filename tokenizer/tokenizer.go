package tokenizer

import (
	"fmt"
	jprose "github.com/jdkato/prose/v2"
	"os"
)

func Tokenize(textFile string) []jprose.Token {
	text, err := os.ReadFile(textFile)
	if err != nil {
		panic(err)
	}
	// NLP
	doc, err := jprose.NewDocument(string(text))
	if err != nil {
		panic(err)
	}

	return doc.Tokens()
}

func String(tok *jprose.Token) string {
	return fmt.Sprintf("%s\t%s\t%s", tok.Text, tok.Tag, tok.Label)
}
