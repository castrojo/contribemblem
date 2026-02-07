//go:build ignore
// +build ignore

package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func main() {
	// Create a 474x96 image (Destiny 2 emblem dimensions)
	img := image.NewRGBA(image.Rect(0, 0, 474, 96))

	// Fill with a simple gradient for visual interest
	for y := 0; y < 96; y++ {
		for x := 0; x < 474; x++ {
			// Create a blue-to-purple gradient
			r := uint8(50 + (x * 100 / 474))
			g := uint8(50 + (y * 100 / 96))
			b := uint8(150 + (x * 100 / 474))
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// Save as JPEG
	f, err := os.Create("test_emblem.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 90}); err != nil {
		panic(err)
	}
}
