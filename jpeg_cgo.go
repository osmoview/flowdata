//go:build cgo

package flowdata

import (
	"bytes"
	"fmt"
	"image"

	libjpeg "github.com/sg3des/libjpeg"
)

var DefaultQualityJPEG = 95

func EncodeJPEG(img image.Image) (b []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	var buf bytes.Buffer
	err = libjpeg.Encode(&buf, img, &libjpeg.EncoderOptions{Quality: DefaultQualityJPEG})

	return buf.Bytes(), err
}

func DecodeJPEG(data []byte) (image.Image, error) {
	return libjpeg.Decode(bytes.NewReader(data), &libjpeg.DecoderOptions{})
}

func DecodeJPEGInfo(data []byte) (image.Config, error) {
	return libjpeg.DecodeConfig(bytes.NewReader(data))
}
