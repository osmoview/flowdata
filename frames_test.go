package flowdata

import (
	"bytes"
	"encoding/json"
	"image"
	"log"
	"os"
	"testing"

	"golang.org/x/image/draw"
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

	b := frame.Bounds()

	dst := image.NewRGBA(image.Rect(0, 0, b.Max.X, b.Max.Y))

	// draw.BiLinear.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	draw.Draw(dst, dst.Bounds(), frame, b.Min, draw.Src)

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

func TestNewFrameBytesRGB(t *testing.T) {
	frame, err := NewFrameRGB([]byte{255, 255, 255}, 1, 1)
	if err != nil {
		log.Fatal(err)
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

	crop, err := frame.Crop(BBox{0.35, 0.8, 0.3, 0.2})
	if err != nil {
		t.Fatal(err)
	}

	cropBytes, err := crop.GetBytes()
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/crop.jpg", cropBytes, 0666)
}

func TestFrameCrop2(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	frame, err := NewFrameImage(imageImage)
	if err != nil {
		t.Fatal(err)
	}

	crop, err := frame.Crop(BBox{0.35, 0.8, 0.3, 0.2})
	if err != nil {
		t.Fatal(err)
	}

	cropBytes, err := crop.GetBytes()
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/crop.jpg", cropBytes, 0666)

	// check crop from crop
	crop, err = crop.Crop(BBox{0.2, 0.2, 0.2, 0.2})
	if err != nil {
		t.Fatal(err)
	}

	cropBytes, err = crop.GetBytes()
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile("testdata/crop2.jpg", cropBytes, 0666)
}

func TestFrameMarshalJSON(t *testing.T) {
	frame, err := NewFrameImage(imageImage)
	if err != nil {
		t.Fatal(err)
	}

	framejson, err := json.Marshal(frame)
	if err != nil {
		t.Fatal(err)
	}

	var f Frame
	if err := json.Unmarshal(framejson, &f); err != nil {
		t.Fatal(err)
	}

	if f.ID != frame.ID {
		t.Fatal("id not equal")
	}

	if len(f.Bytes) == 0 || !bytes.Equal(f.Bytes, frame.Bytes) {
		t.Fatal("frame bytes is empty or unequal")
	}
}

func TestFrameUnmarshalJSONBytes(t *testing.T) {
	b, err := json.Marshal(map[string]any{"frame": imageBytes})
	if err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Frame *Frame
	}

	if err := json.Unmarshal(b, &resp); err != nil {
		t.Fatal(err)
	}

	if resp.Frame == nil {
		t.Fatal("frame is nil")
	}
	if resp.Frame.ID == "" {
		t.Fatal("frame id is empty")
	}
}

func BenchmarkResizeNearestNeighbor(b *testing.B) {
	rect := imageImage.Bounds()
	height := int(float64(720) / (float64(rect.Max.X) / float64(rect.Max.Y)))
	img := image.NewRGBA(image.Rect(0, 0, 720, height))

	for i := 0; i < b.N; i++ {
		draw.NearestNeighbor.Scale(img, img.Rect, imageImage, rect, draw.Over, nil)
	}
}

func BenchmarkResizeApproxBiLinear(b *testing.B) {
	rect := imageImage.Bounds()
	height := int(float64(720) / (float64(rect.Max.X) / float64(rect.Max.Y)))
	img := image.NewRGBA(image.Rect(0, 0, 720, height))

	for i := 0; i < b.N; i++ {
		draw.ApproxBiLinear.Scale(img, img.Rect, imageImage, rect, draw.Over, nil)
	}
}

func BenchmarkResizeBiLinear(b *testing.B) {
	rect := imageImage.Bounds()
	height := int(float64(720) / (float64(rect.Max.X) / float64(rect.Max.Y)))
	img := image.NewRGBA(image.Rect(0, 0, 720, height))

	for i := 0; i < b.N; i++ {
		draw.BiLinear.Scale(img, img.Rect, imageImage, rect, draw.Over, nil)
	}
}

func BenchmarkResizeCatmullRom(b *testing.B) {
	rect := imageImage.Bounds()
	height := int(float64(720) / (float64(rect.Max.X) / float64(rect.Max.Y)))
	img := image.NewRGBA(image.Rect(0, 0, 720, height))

	for i := 0; i < b.N; i++ {
		draw.CatmullRom.Scale(img, img.Rect, imageImage, rect, draw.Over, nil)
	}
}
