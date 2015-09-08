package google

import (
    "net/url"
    "net/http"
    "time"
    "encoding/json"
    "math/rand"
)

type ImageSearchApi struct {
    ResponseData struct {
        Results []struct {
            UnescapedUrl string
        }
    }
}

func ImageSearch(query string) string {
    rand.Seed(time.Now().UnixNano())
    v := url.Values{}
    v.Set("v", "1.0")
    v.Set("rsz", "8")
    v.Set("q", query)
    v.Set("safe", "active")
    response, _ := http.Get("http://ajax.googleapis.com/ajax/services/search/images?"+v.Encode())
    dec := json.NewDecoder(response.Body)
    var data ImageSearchApi
    dec.Decode(&data)
    if len(data.ResponseData.Results) == 0 {
        return ""
    }
    i := rand.Intn(len(data.ResponseData.Results))

    return data.ResponseData.Results[i].UnescapedUrl
}
