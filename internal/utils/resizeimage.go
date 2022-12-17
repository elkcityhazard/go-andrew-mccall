package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"

	"golang.org/x/image/draw"
)

func (m *Tools) ResizeImage(currentImage, pathToNewFile string) error {

	input, _ := os.Open(currentImage)
	defer input.Close()

	output, _ := os.Create(pathToNewFile)
	defer output.Close()

	// Decode the image (from PNG to image.Image):
	src, _ := jpeg.Decode(input)

	// Set the expected size that you want:
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/2, src.Bounds().Max.Y/2))

	// Resize:
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	// Encode to `output`:
	jpeg.Encode(bytes.NewReader(src), dst)

	return nil
}
