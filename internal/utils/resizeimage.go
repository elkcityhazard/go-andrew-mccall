package utils

import (
	"fmt"
	"github.com/elkcityhazard/go-andrew-mccall/internal/models"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

// SubImager type is created to use type assertion to cast SubImage to the image

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func (m *Tools) ResizeImage(file *UploadedFile, pathToFile string, user *models.User, w http.ResponseWriter, r *http.Request) error {

	fmt.Println("User email: ", user.Email)

	// Open up the damn file

	input, _ := os.Open(pathToFile)
	//	don't forget to defer closing it to avoid memory leak
	defer input.Close()

	// create a new image buffer

	imageBuff := make([]byte, 512)

	// read the original file into the image buffer so we can check what type of file it is

	_, err := input.Read(imageBuff)

	if err != nil {
		return err
	}

	// get the file type from the buffer

	fileType := http.DetectContentType(imageBuff)

	// create the output file

	output, _ := os.Create(path.Join("./static/uploads/", user.Email, fmt.Sprintf("avatar-%s", filepath.Base(file.NewFileName))))

	// defer closing the output file until we are done writing it and the function exits (avoid memory leak)
	defer output.Close()

	// seek the file back to the beginning or else we won't be able to write the whole file

	input.Seek(0, 0)

	// create a new image variable

	var src image.Image

	// determine if the original file was a png or jpeg before continuing

	if strings.EqualFold(fileType, "image/png") {
		src, _ = png.Decode(input)
	} else {
		// Decode the image (from PNG to image.Image):
		src, _ = jpeg.Decode(input)
	}

	// create a whole new sized image

	// Set the expected size that you want:
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/4, src.Bounds().Max.Y/4))

	// At returns the color of the pixel at (x, y).
	// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
	// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.

	bounds := dst.Bounds()

	//	Get Width

	width := bounds.Dx()

	height := bounds.Dy()

	fmt.Println(width, height)

	// created the cropped size of the image

	cropSize := image.Rect(0, 0, width/2, height/2)

	// dynamically get focal point based on original width

	focalX := math.Floor(float64(width)) * 1.33
	focalY := math.Floor(float64(height))
	//focalY = 0
	//focalX = float64(width)

	//This is the place of the left and top padding of image that you want to crop. In this case we add padding left to width * 1.33 and padding top is the height of the destination image

	cropSize = cropSize.Add(image.Point{int(focalX), int(focalY)})

	//SubImage returns an image representing the portion of the image p visible through r. The returned value shares pixels with the original image.

	croppedImage := src.(SubImager).SubImage(cropSize)

	croppedImageFile, err := os.Create("./static/uploads/cropped.png")

	if err != nil {
		log.Fatalln(err)
	}

	defer croppedImageFile.Close()

	if err := png.Encode(croppedImageFile, croppedImage); err != nil {
		log.Fatalln(err)
		return err
	}

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
