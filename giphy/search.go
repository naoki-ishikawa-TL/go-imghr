package giphy

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const ApiKey = "dc6zaTOxFJmzC"

type SearchResponseImage struct {
	Type        string
	Id          string
	Url         string
	BitlyGifUrl string `json:"bitly_gif_url"`
	BitlyUrl    string `json:"bitly_url"`
	EmbedUrl    string `json:"embed_url"`
	Images      struct {
		Original struct {
			Url  string
			Mp4  string
			Webp string
		}
	}
}

type SearchResponse struct {
	Data []SearchResponseImage
	Meta struct {
		Status int
		Msg    string
	}
}

type RandomResponse struct {
	Data struct {
		Type             string
		Id               string
		Url              string
		ImageOriginalUrl string `json:"image_original_url"`
		ImageUrl         string `json:"image_url"`
	}
	Meta struct {
		Status int
		Msg    string
	}
}

func Search(query string) string {
	v := &url.Values{}
	v.Set("api_key", ApiKey)
	v.Set("q", query)
	v.Set("limit", "1")

	response, _ := http.Get("http://api.giphy.com/v1/gifs/search?" + v.Encode())
	var data SearchResponse
	dec := json.NewDecoder(response.Body)
	dec.Decode(&data)

	return data.Data[0].Images.Original.Url
}

func Random() string {
	v := &url.Values{}
	v.Set("api_key", ApiKey)

	response, _ := http.Get("http://api.giphy.com/v1/gifs/random?" + v.Encode())
	var data RandomResponse
	dec := json.NewDecoder(response.Body)
	dec.Decode(&data)

	return data.Data.ImageOriginalUrl
}
