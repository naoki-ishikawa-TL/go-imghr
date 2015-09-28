package tenki

import (
    "net/http"
    "io/ioutil"
    "github.com/moovweb/gokogiri"
    "github.com/moovweb/gokogiri/xpath"
)

func GetTemperature() string {
    resp, _ := http.Get("http://www.jma.go.jp/jp/amedas_h/today-44132.html?areaCode=000&groupCode=30")
    page, _ := ioutil.ReadAll(resp.Body)

    doc, _ := gokogiri.ParseHtml(page)
    defer doc.Free()

    xps := xpath.Compile("//*[@id=\"tbl_list\"]/tr/td[2]")
    ss, _ := doc.Root().Search(xps)

    var temperature string
    for _, s := range ss {
        if len(s.InnerHtml()) > 2 {
            temperature = s.InnerHtml()
        }
    }

    return temperature
}
