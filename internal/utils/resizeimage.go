package utils

import (
	"fmt"
	"github.com/elkcityhazard/go-andrew-mccall/internal/models"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

func (m *Tools) ResizeImage(file *UploadedFile, pathToFile string, user *models.User, w http.ResponseWriter, r *http.Request) error {

	fmt.Println("User email: ", user.Email)

	input, _ := os.Open(pathToFile)
	defer input.Close()

	imageBuff := make([]byte, 512)

	_, err := input.Read(imageBuff)

	if err != nil {
		return err
	}

	fileType := http.DetectContentType(imageBuff)

	output, _ := os.Create(path.Join("./static/uploads/", user.Email, fmt.Sprintf("avatar-%s", filepath.Base(file.NewFileName))))
	defer output.Close()

	input.Seek(0, 0)

	var src image.Image

	if strings.EqualFold(fileType, "image/png") {
		src, _ = png.Decode(input)
	} else {
		// Decode the image (from PNG to image.Image):
		src, _ = jpeg.Decode(input)
	}

	fmt.Println(fileType)

	// Set the expected size that you want:
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/5, src.Bounds().Max.Y/5))

	// Resize:
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	// Encode to `output`:

	if strings.EqualFold(fileType, "image/png") {
		png.Encode(output, dst)
	} else {
		jpeg.Encode(output, dst, nil)
	}

	return nil
}
