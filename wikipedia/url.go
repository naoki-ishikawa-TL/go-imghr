package wikipedia

import (
	"net/url"
)

func GenerateJaWikipediaURL(query string) string {
	return "https://ja.wikipedia.org/wiki/" + url.QueryEscape(query)
}
