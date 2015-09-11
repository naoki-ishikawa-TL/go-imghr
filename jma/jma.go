package jma

import (
	"../image"
	"crypto/sha1"
	"fmt"
	"github.com/gographics/imagick/imagick"
	"log"
	"strconv"
)

const (
	MIN_WIDTH  int = 35
	MAX_WIDTH  int = 36
	MIN_HEIGHT int = 29
	MAX_HEIGHT int = 30
)

type JmaImageGenerator struct {
	RequestChan  chan string
	ResponseChan chan string
}

func NewJmaImageGenerator() *JmaImageGenerator {
	requestChan := make(chan string)
	responseChan := make(chan string)

	this := &JmaImageGenerator{RequestChan: requestChan, ResponseChan: responseChan}

	go func() {
		for {
			select {
			case date := <-requestChan:
				digest := fmt.Sprintf("%x", sha1.Sum([]byte(date+"jma")))
				var imgPath string
				if image.FileIsExist("public/data/" + digest + ".png") {
					imgPath = "data/" + digest + ".png"
				} else {
					imgPath = this.GenerateImage(date, digest)
				}

				responseChan <- imgPath
			}
		}
	}()

	return this
}

func (this *JmaImageGenerator) GenerateImage(date string, digest string) string {
	mapImage := image.ReadImageFromAsset("data/jma_map.png")
	defer mapImage.Destroy()
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
			tmp := image.GetImageFromUrl("http://www.jma.go.jp/jp/highresorad/highresorad_tile/HRKSNC/" + date + "/" + date + "/zoom6/" + strconv.Itoa(w) + "_" + strconv.Itoa(h) + ".png")
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

	mapImage.WriteImage("public/data/" + digest + ".png")
	return "data/" + digest + ".png"
}

func (this *JmaImageGenerator) Generate(date string) string {
	this.RequestChan <- date
	return <-this.ResponseChan
}
