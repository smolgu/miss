// Copyright 2018 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package img

import (
	"image"
	"net/http"
	// jpeg and png used for decoding remote images
	_ "image/jpeg"
	_ "image/png"

	"github.com/cenkalti/dominantcolor"
	"github.com/corona10/goimagehash"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
)

// DominantColor retruns hex primary color for image
func DominantColor(img image.Image) string {
	return dominantcolor.Hex(dominantcolor.Find(img))
}

// PerceptiveHash returns perceptive hash for image
func PerceptiveHash(img image.Image) (string, error) {
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return "", err
	}
	return hash.ToString(), nil
}

// HasFaces check is image contains face
func HasFaces(img image.Image) bool {
	panic("not implemented")
}

// Crop smart crop
func Crop(img image.Image, w, h int) image.Image {
	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	topCrop, _ := analyzer.FindBestCrop(img, 250, 250)
	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	if croppedimg, ok := img.(SubImager); ok {
		return croppedimg.SubImage(topCrop)
	}
	return img
}

// FromURL load remote image
func FromURL(url string) (img image.Image, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	img, _, err = image.Decode(resp.Body)
	if err != nil {
		return
	}
	return
}
