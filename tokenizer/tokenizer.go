package tokenizer

import (
	"fmt"
	"github.com/jdkato/prose/v2"
	"os"
)

func Tokenize(textFile string) []prose.Token {
	text, err := os.ReadFile(textFile)
	if err != nil {
		panic(err)
	}
	// NLP
	doc, err := prose.NewDocument(string(text))
	if err != nil {
		panic(err)
	}

	return doc.Tokens()
}

func String(tok *prose.Token) string {
	return fmt.Sprintf("%s\t%s\t%s", tok.Text, tok.Tag, tok.Label)
}
