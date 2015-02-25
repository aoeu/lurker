/*
The thumb file wraps the excellent smartcrop and gift libraries for
generating thumbnail images from comic panels.
*/
package lurker

import (
	"errors"
	"github.com/disintegration/gift"
	"github.com/muesli/smartcrop"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

// VisCrop takes a pointer to an image and uses the smartcrop library to
// crop and generate a new thumbnail image.
func VisCrop(src image.Image) (image.Image, error) {
	crop, _ := smartcrop.SmartCrop(&src, 100, 100)
	log.Println("**********")
	log.Printf("%+v\n", crop)

	dstRect := image.Rect(crop.X, crop.Y, crop.Width+crop.X, crop.Height+crop.Y)

	g := gift.New(gift.Crop(dstRect))
	dst := image.NewRGBA(dstRect)

	g.Draw(dst, src)
	return dst, nil
}

// BlurMeDown takes a pointer to an image and uses gift to blur it artfully
// and return it as a new image.
func BlurMeDown(src *image.Image) (image.Image, error) {
	return nil, errors.New("oops")
}
