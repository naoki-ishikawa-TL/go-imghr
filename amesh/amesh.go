package amesh

import (
    "../image"
    "github.com/gographics/imagick/imagick"
    "fmt"
)

func GenerateAmeshImage(date string) *imagick.MagickWand {
    mapImage := image.ReadImageFromAsset("data/jma_map.jpg")
    maskImage := image.ReadImageFromAsset("data/jma_mask.png")
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

    return mapImage
}
