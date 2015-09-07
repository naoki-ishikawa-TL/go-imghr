package image

import (
    "github.com/gographics/imagick/imagick"
    "crypto/sha1"
    "fmt"
    "sync"
)

type ImageGenerator struct {
    Table map[string] func(string) *imagick.MagickWand
    LockTable map[string] *sync.Mutex
}

func NewImageGenerator() *ImageGenerator {
    table := make(map[string] func(string) *imagick.MagickWand)
    lockTable := make(map[string] *sync.Mutex)
    return &ImageGenerator{Table: table, LockTable: lockTable}
}

func (this *ImageGenerator) AddGenerator(imageType string, generateFunc func(string) *imagick.MagickWand) {
    this.Table[imageType] = generateFunc
    this.LockTable[imageType] = new(sync.Mutex)
}

func (this *ImageGenerator) Generate(imageType string, date string) string {
    m := this.LockTable[imageType]
    m.Lock()
    defer m.Unlock()
    digest := fmt.Sprintf("%x", sha1.Sum([]byte(date+imageType)))
    if FileIsExist("public/data/"+digest+".png") {
        return "data/"+digest+".png"
    }

    generateFunc := this.Table[imageType]
    rainImage := generateFunc(date)
    defer rainImage.Destroy()
    rainImage.WriteImage("public/data/"+digest+".png")

    return "data/"+digest+".png"
}
