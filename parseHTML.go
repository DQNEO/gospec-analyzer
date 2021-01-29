// Usage:
//   go run parseHTML.go
package main
import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jdkato/prose"
	"os"
)

func main() {
	//html, err := ioutil.ReadFile("./spec.html")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Print(string(html))
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
		wordCount[tok.Text]++
	}

	type sameFreqGroup []string
	var frequency []sameFreqGroup = make([]sameFreqGroup, len(wordCount))
	for w, n := range wordCount {
		//fmt.Printf("%3d\t%s\n", n,w)
		frequency[n] = append(frequency[n], w)
	}

	for n, grp := range frequency {
		for _, w := range grp {
			fmt.Printf("%4d\t%s\n", n, w)
		}
	}
	fmt.Printf("%4d\tTotal\n", len(wordCount))
}
