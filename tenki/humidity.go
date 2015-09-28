package tenki

import (
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xpath"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetHumidity() string {
	resp, _ := http.Get("http://www.jma.go.jp/jp/amedas_h/today-44132.html?areaCode=000&groupCode=30")
	page, _ := ioutil.ReadAll(resp.Body)

	doc, _ := gokogiri.ParseHtml(page)
	defer doc.Free()

	xps := xpath.Compile("//*[@id=\"tbl_list\"]/tr/td[7]")
	ss, _ := doc.Root().Search(xps)

	var humidity string
	for _, s := range ss {
		if _, err := strconv.Atoi(s.InnerHtml()); err == nil {
			humidity = s.InnerHtml()
		}
	}

	return humidity
}
