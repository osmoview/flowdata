package flowdata

import (
	"image"
	"os"
	"testing"
)

var (
	imageFilename = "testdata/logo.jpg"
	imageBytes    []byte
	imageImage    image.Image
)

func init() {
	var err error
	imageBytes, err = os.ReadFile(imageFilename)
	if err != nil {
		panic(err)
	}

	imageImage, err = DecodeJPEG(imageBytes)
	if err != nil {
		panic(err)
	}
}

func TestNewFrameBytes(t *testing.T) {
	frame, err := NewFrameBytes(imageBytes)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := frame.GetBytes()
	if err != nil {
		t.Fatal(err)
	}

	if len(bytes) == 0 {
		t.Error("frame bytes is empty")
	}

	if _, err := DecodeJPEG(bytes); err != nil {
		t.Error(err)
	}

	if frame.ID == "" {
		t.Error("no id")
	}

	if frame.Time.IsZero() {
		t.Error("no time")
	}
}

func TestNewFrameImage(t *testing.T) {
	frame, err := NewFrameImage(imageImage)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := frame.GetBytes()
	if err != nil {
		t.Fatal(err)
	}

	if len(bytes) == 0 {
		t.Error("frame bytes is empty")
	}

	if _, err := DecodeJPEG(bytes); err != nil {
		t.Error(err)
	}

	if frame.ID == "" {
		t.Error("no id")
	}

	if frame.Time.IsZero() {
		t.Error("no time")
	}
}

func TestFrameCrop(t *testing.T) {
	frame, err := NewFrameImage(imageImage)
	if err != nil {
		t.Fatal(err)
	}

	crop, err := frame.Crop(BBox{0.1, 0.1, 0.8, 0.8})
	if err != nil {
		t.Fatal(err)
	}

	cropBytes, err := crop.GetBytes()
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/crop.jpg", cropBytes, 0666)
}
