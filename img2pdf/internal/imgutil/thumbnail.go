package imgutil

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

func GenerateThumbnail(img image.Image, width int) string {
	thumb := imaging.Resize(img, width, 0, imaging.Lanczos)
	var buf bytes.Buffer
	jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: 75})
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func ThumbnailFromFile(path string, width int) (string, error) {
	img, err := LoadImage(path)
	if err != nil {
		return "", err
	}
	return GenerateThumbnail(img, width), nil
}
