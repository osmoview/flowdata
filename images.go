package flowdata

import (
	"bytes"
	"image"
	"image/jpeg"
)

var DefaultQualityJPEG = 95

func EncodeJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: DefaultQualityJPEG})

	return buf.Bytes(), err
}

func DecodeJPEG(data []byte) (image.Image, error) {
	return jpeg.Decode(bytes.NewReader(data))
}
