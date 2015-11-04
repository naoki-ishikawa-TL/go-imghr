package wikipedia

import (
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xpath"
	"io/ioutil"
	"net/http"
	"errors"
)

func GetSummary(query string) (string, error) {
	resp, _ := http.Get(GenerateJaWikipediaURL(query))
	if resp.StatusCode != 200 {
		return "", errors.New("page not found")
	}
	page, _ := ioutil.ReadAll(resp.Body)

	doc, _ := gokogiri.ParseHtml(page)
	defer doc.Free()

	xps := xpath.Compile("//*[@id=\"mw-content-text\"]/p[1]")
	ss, _ := doc.Root().Search(xps)

	content := ""
	for _, s := range ss {
		content += s.Content()
	}

	return content, nil
}
