//go:build cgo

package flowdata

import (
	"bytes"
	"image"

	turbojpeg "github.com/pixiv/go-libjpeg/jpeg"
)

var DefaultQualityJPEG = 95

func EncodeJPEG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := turbojpeg.Encode(&buf, img, &turbojpeg.EncoderOptions{Quality: DefaultQualityJPEG})

	return buf.Bytes(), err
}

func DecodeJPEG(data []byte) (image.Image, error) {
	return turbojpeg.Decode(bytes.NewReader(data), &turbojpeg.DecoderOptions{})
}

func DecodeJPEGInfo(data []byte) (image.Config, error) {
	return turbojpeg.DecodeConfig(bytes.NewReader(data))
}
