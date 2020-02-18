package imageutil

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/jpeg" // required for jpeg formatted images
	_ "image/png"  // required for png formatted images
	"io"
	"strings"

	"github.com/disintegration/imaging"
)

// Widths are the 'preselected' sizes of the thumbnails
var Widths = []int{32, 64}

var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

// IsImage checks if a byte array is image
func IsImage(incipit []byte) bool {
	incipitStr := string(incipit)
	for magic := range magicTable {
		if strings.HasPrefix(incipitStr, magic) {
			return true
		}
	}

	return false
}

// GetImageSize returns the size of the image
func GetImageSize(r io.Reader) (int, int, error) {
	img, _, err := image.DecodeConfig(r)
	if err != nil {
		return -1, -1, err
	}
	return img.Width, img.Height, nil
}

// ResizeImage resizes image to expected sizes
func ResizeImage(r io.Reader, widthSize int) ([]byte, error) {
	srcImg, _, err := image.Decode(r)
	dstImgSize := imaging.Resize(srcImg, widthSize, 0, imaging.Lanczos) // whatever interpolation method
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, dstImgSize, nil) // make this jpeg easier for life
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
