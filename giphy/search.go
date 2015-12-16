package giphy

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const ApiKey = "dc6zaTOxFJmzC"

type GiphyResponse struct {
	Meta struct {
		Status int
		Msg    string
	}
}

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
	GiphyResponse
	Data []SearchResponseImage
}

type RandomResponse struct {
	GiphyResponse
	Data struct {
		Type                         string
		Id                           string
		Url                          string
		ImageOriginalUrl             string `json:"image_original_url"`
		ImageUrl                     string `json:"image_url"`
		ImageMp4Url                  string `json:"image_mp4_url"`
		ImageFrames                  string `json:"image_frames"`
		ImageWidth                   string `json:"image_width"`
		ImageHeight                  string `json:"image_height"`
		FixedHeightDownsampledUrl    string `json:"fixed_height_downsampled_url"`
		FixedHeightDownsampledWidth  string `json:"fixed_height_downsampled_width"`
		FixedHeightDownsampledHeight string `json:"fixed_height_downsampled_height"`
		FixedWidthDownsampledUrl     string `json:"fixed_width_downsampled_url"`
		FixedWidthDownsampledWidth   string `json:"fixed_width_downsampled_width"`
		FixedWidthDownsampledHeight  string `json:"fixed_width_downsampled_height"`
		FixedHeightSmallUrl          string `json:"fixed_height_small_url"`
		FixedHeightSmallStillUrl     string `json:"fixed_height_small_still_url"`
		FixedHeightSmallWidth        string `json:"fixed_height_small_width"`
		FixedHeightSmallHeight       string `json:"fixed_height_small_height"`
		FixedWidthSmallUrl           string `json:"fixed_width_small_url"`
		FixedWidthSmallStillUrl      string `json:"fixed_Width_small_still_url"`
		FixedWidthSmallWidth         string `json:"fixed_width_small_width"`
		FixedWidthSmallHeight        string `json:"fixed_width_small_height"`
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
