package jma

import (
    "log"
    "strconv"
    "github.com/gographics/imagick/imagick"
    "../image"
)

const (
    MIN_WIDTH int = 35
    MAX_WIDTH int = 36
    MIN_HEIGHT int = 29
    MAX_HEIGHT int = 30
)

func GenerateJmaImage(date string) *imagick.MagickWand {
    mapImage := image.ReadImageFromAsset("data/jma_map.png")
    maskImage := image.ReadImageFromAsset("data/jma_mask.png")
    defer maskImage.Destroy()
    manucipalityImage := image.ReadImageFromAsset("data/jma_manucipality.png")
    defer manucipalityImage.Destroy()
    ihrImage := image.ReadImageFromAsset("data/ihr.png")
    defer ihrImage.Destroy()

    rainImage := imagick.NewMagickWand()
    defer rainImage.Destroy()
    for w := MIN_WIDTH; w <= MAX_WIDTH; w++ {
        tmpImage := imagick.NewMagickWand()
        defer tmpImage.Destroy()
        for h := MIN_HEIGHT; h <= MAX_HEIGHT; h++ {
            tmp := image.GetImageFromUrl("http://www.jma.go.jp/jp/highresorad/highresorad_tile/HRKSNC/"+date+"/"+date+"/zoom6/"+strconv.Itoa(w)+"_"+strconv.Itoa(h)+".png")
            defer tmp.Destroy()
            tmpImage.AddImage(tmp)
            tmpImage.SetLastIterator()
        }
        tmpImage.SetFirstIterator()
        tmpImage = tmpImage.AppendImages(true)
        rainImage.AddImage(tmpImage)
        rainImage.SetLastIterator()
    }
    rainImage.SetFirstIterator()
    rainImage = rainImage.AppendImages(false)
    rainImage.AdaptiveResizeImage(2048, 2048)

    err := mapImage.CompositeImage(rainImage, imagick.COMPOSITE_OP_OVER, 0, 0)
    if err != nil {
        log.Println(err)
    }
    err = mapImage.CompositeImage(maskImage, imagick.COMPOSITE_OP_OVER, 0, 0)
    if err != nil {
        log.Println(err)
    }
    err = mapImage.CompositeImage(manucipalityImage, imagick.COMPOSITE_OP_OVER, 0, 0)
    if err != nil {
        log.Println(err)
    }
    err = mapImage.CropImage(1000, 500, 580, 750)
    if err != nil {
        log.Println(err)
    }
    err = mapImage.CompositeImage(ihrImage, imagick.COMPOSITE_OP_OVER, 750, 300)
    if err != nil {
        log.Println(err)
    }

    return mapImage
}
