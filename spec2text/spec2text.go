package spec2text

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

func GetTextFromHTML(html io.Reader) string {
	gdoc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		panic(err)
	}
	mainTag := gdoc.Find("main")
	mainTag.Find("pre").Remove()
	mainTag.Find("code").Remove()
	text := mainTag.Text()
	return text
}
