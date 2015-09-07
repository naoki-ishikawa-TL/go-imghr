package image

import (
    "github.com/gographics/imagick/imagick"
    "crypto/sha1"
    "fmt"
)

func GenerateImageForBot(date string, imageType string, generateFunc func(string) *imagick.MagickWand) string {
    digest := fmt.Sprintf("%x", sha1.Sum([]byte(date+imageType)))
    if FileIsExist("public/data/"+digest+".png") {
        return "data/"+digest+".png"
    }

    rainImage := generateFunc(date)
    defer rainImage.Destroy()
    rainImage.WriteImage("public/data/"+digest+".png")

    return "data/"+digest+".png"
}
