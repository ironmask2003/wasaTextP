package imageFunctions

import (
	"encoding/base64"
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	"os"
)

// Function return the path of the user photo profile
func SetDefaultPhoto(userId int) string {
	return fmt.Sprintf("./storage/%d/user_photo.png", userId)
}

// Function convert an image in base64
func ImageToBase64(filename string) (string, error) {
	imageFile, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() { err = imageFile.Close() }()

	imageData, err := io.ReadAll(imageFile)
	if err != nil {
		return "", err
	}

	base64 := base64.StdEncoding.EncodeToString(imageData)
	return base64, err
}

// Function used to save and crop an image
func SaveAndCrop(filename string, w uint, h uint) error {
	file, err := os.Open(filename) // Opena the image file
	if err != nil {
		return err
	}
	defer func() { err = file.Close() }()

	// Decode the image
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	resizedImg := resize.Resize(w, h, img, resize.NearestNeighbor)
	// Save cropped image
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { err = out.Close() }()
	if err := jpeg.Encode(out, resizedImg, nil); err != nil {
		return err
	}

	return err
}
