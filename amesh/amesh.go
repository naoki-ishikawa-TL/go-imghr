package amesh

import (
    "../image"
    "github.com/gographics/imagick/imagick"
    "crypto/sha1"
    "fmt"
)

type AmeshImageGenerator struct {
    RequestChan chan string
    ResponseChan chan string
}

func NewAmeshImageGenerator() *AmeshImageGenerator {
    requestChan := make(chan string)
    responseChan := make(chan string)

    this := &AmeshImageGenerator{RequestChan: requestChan, ResponseChan: responseChan}

    go func() {
        for {
            select {
            case date := <-requestChan:
                digest := fmt.Sprintf("%x", sha1.Sum([]byte(date+"amesh")))
                var imgPath string
                if image.FileIsExist("public/data/"+digest+".png") {
                    imgPath = "data/"+digest+".png"
                } else {
                    imgPath = this.ImageGenerate(date, digest)
                }

                responseChan <- imgPath
            }
        }
    }()

    return this
}

func (this *AmeshImageGenerator) ImageGenerate(date string, digest string) string {
    mapImage := image.ReadImageFromAsset("data/amesh_map.jpg")
    defer mapImage.Destroy()
    maskImage := image.ReadImageFromAsset("data/amesh_mask.png")
    defer maskImage.Destroy()
    rainImage := image.GetImageFromUrl("http://tokyo-ame.jwa.or.jp/mesh/100/"+date+".gif")
    defer rainImage.Destroy()

    err := mapImage.CompositeImage(rainImage, imagick.COMPOSITE_OP_OVER, 0, 0)
    if err != nil {
        fmt.Println(err)
    }
    err = mapImage.CompositeImage(maskImage, imagick.COMPOSITE_OP_OVER, 0, 0)
    if err != nil {
        fmt.Println(err)
    }
    err = mapImage.CropImage(1300, 600, 900, 650)
    if err != nil {
        fmt.Println(err)
    }
    ihrImage := image.ReadImageFromAsset("data/ihr.png")
    defer ihrImage.Destroy()

    err = mapImage.CompositeImage(ihrImage, imagick.COMPOSITE_OP_OVER, 970, 320)
    if err != nil {
        fmt.Println(err)
    }

    mapImage.WriteImage("public/data/"+digest+".png")
    return "data/"+digest+".png"
}

func (this *AmeshImageGenerator) Generate(date string) string {
    this.RequestChan <- date
    return <- this.ResponseChan
}

