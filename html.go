package main

import (
	"github.com/PuerkitoBio/goquery"
	"os"
)

func extractText(specFile string) string {
	f, err := os.Open(specFile)
	if err != nil {
		panic(err)
	}
	gdoc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		panic(err)
	}
	mainTag := gdoc.Find("main")
	mainTag.Find("pre").Remove()
	text := mainTag.Text()
	return text
}
