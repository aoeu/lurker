/*
thumb_test.go runs tests for the thumb package.
*/
package lurker

import (
	_ "fmt"
	"image"
	"image/png"
	"os"
	"strings"
	"testing"
)

var (
	testImages = [4]string{
		"test_data/douglas_engelbart_1925_2013.png",
		"test_data/small_moon.png",
		"test_data/20131209.png",
		"test_data/20141210.png"}
)

func TestVisCrop(t *testing.T) {
	for _, path := range testImages {
		filename := strings.Split(strings.Split(path, "/")[1], ".")[0]

		srcFile, _ := os.Open(path)
		defer srcFile.Close()
		src, _, err := image.Decode(srcFile)
		if err != nil {
			t.Error(err)
		}
		
		dst, err := VisCrop(src)
		if err != nil {
			t.Error(err)
		}
		// fmt.Printf("%v\n", dst)
		dstFile, _ := os.Create(filename + "_thumb.png")
		defer dstFile.Close()
		png.Encode(dstFile, dst)
	}
}
