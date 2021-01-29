package main

import (
	"fmt"
	"os"
	"github.com/PuerkitoBio/goquery"
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
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		panic(err)
	}
	doc.Find("pre").Remove()
	noCodeText := doc.Text()
	fmt.Print(noCodeText)
}
