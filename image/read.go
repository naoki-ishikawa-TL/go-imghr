package image

import (
    "github.com/gographics/imagick/imagick"
    "net/http"
    "io/ioutil"
    "log"
    "../assets"
)

func ReadImageFromAsset(name string) *imagick.MagickWand {
    imageBlob, _ := assets.Asset(name)
    image := imagick.NewMagickWand()
    err := image.ReadImageBlob(imageBlob)
    if err != nil {
        log.Println(err)
    }

    return image
}

func ReadImageFromFile(path string) *imagick.MagickWand {
    image := imagick.NewMagickWand()
    err := image.ReadImage(path)
    if err != nil {
        log.Println(err)
    }

    return image
}

func GetImageFromUrl(url string) *imagick.MagickWand {
    log.Println(url)
    response, err := http.Get(url)
    if err != nil {
        log.Println(err)
    }
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Println(err)
    }

    image := imagick.NewMagickWand()
    err = image.ReadImageBlob(body)
    if err != nil {
        log.Println(err)
    }

    return image
}

