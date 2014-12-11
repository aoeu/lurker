/*
The thumb package wraps the excellent smartcrop and gift libraries for
generating thumbnail images from comic panels.
*/
package thumb

import (
	_ "fmt"
	"image"
	"errors"
	_ "github.com/muesli/smartcrop"
	_ "github.com/disintegration/gift"
)

// VisCrop takes a pointer to an image and uses the smartcrop library to
// crop and generate a new thumbnail image.
func VisCrop(src *image.Image) (image.Image, error) {
	return nil, errors.New("oops")
}

// BlurMeDown takes a pointer to an image and uses gift to blur it artfully
// and return it as a new image.
func BlurMeDown(src *image.Image) (image.Image, error) {
	return nil, errors.New("oops")
}

