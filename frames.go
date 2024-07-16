package flowdata

import (
	"image"
	"image/color"
	"time"
)

var DefaultFrameTimeLayout = "2006-01-02_15:04:05"

// Frame contains image data and time when that image was received.
// it can contain only `bytes` or `image` or both of it.
// required to be a careful with reassigning an image data
type Frame struct {
	// time when frame received
	Time time.Time

	// presentation of an image in the image.Image interface
	Image image.Image

	// raw bytes of the image(jpeg usually)
	Bytes []byte

	// image info
	Info image.Config

	// video position in ms
	Position int

	// TODO: support different image types
	// image format, default: jpeg
	// Format string
}

func NewFrameBytes(data []byte) (*Frame, error) {
	// TODO: support different image types
	info, err := DecodeJPEGInfo(data)
	if err != nil {
		return nil, err
	}

	return &Frame{
		Time:  time.Now(),
		Bytes: data,
		Info:  info,
	}, nil
}

func NewFrameImage(img image.Image) (*Frame, error) {
	b := img.Bounds().Max

	return &Frame{
		Time:  time.Now(),
		Image: img,
		Info:  image.Config{ColorModel: img.ColorModel(), Width: b.X, Height: b.Y},
	}, nil
}

// GetImage returns image.Image
func (f *Frame) GetImage() image.Image {
	if f.Image != nil {
		return f.Image
	}

	return f
}

// GetTime returns time when image received
func (f *Frame) GetTime() time.Time {
	return f.Time
}

func (f *Frame) FormatTime() string {
	return f.Time.Format(DefaultFrameTimeLayout)
}

func (f *Frame) Timestamp() int64 {
	return f.Time.Unix()
}

func (f *Frame) TimestampNano() int64 {
	return f.Time.UnixNano()
}

func (f *Frame) TimestampMilli() int64 {
	return f.Time.UnixMilli()
}

// GetSize returns image resolution
func (f *Frame) GetSize() (w, h int) {
	if f.Image != nil {
		b := f.Image.Bounds().Max
		return b.X, b.Y
	}

	return f.Info.Width, f.Info.Height
}

// GetBytes return slice from `bytes` field, if it exists, overwise encode image to the jpeg
func (f *Frame) GetBytes() []byte {
	// if bytes exists then just return it
	if len(f.Bytes) > 0 {
		return f.Bytes
	}

	// if image exists encode it to the jpeg
	if f.Image != nil {
		f.Bytes, _ = EncodeJPEG(f.Image)
		return f.Bytes
	}

	// if nothing exists, then do not return anything
	return nil
}

// Bounds of image received from jpeg.DecodeConfig
func (f *Frame) Bounds() image.Rectangle {
	if f.Image != nil {
		return f.Image.Bounds()
	}

	return image.Rect(0, 0, f.Info.Width, f.Info.Height)
}

// At returns color by specified pixel
func (f *Frame) At(x, y int) color.Color {
	// if image exist returns color from it
	if f.Image != nil {
		return f.Image.At(x, y)
	}

	// if bytes exists, decode it as jpeg and return color from obtained image.Image
	if len(f.Bytes) > 0 {
		f.Image, _ = DecodeJPEG(f.Bytes)
		if f.Image != nil {
			return f.Image.At(x, y)
		}
	}

	return color.Black
}

// ColorModel resolved from the image config
func (f *Frame) ColorModel() color.Model {
	if f.Image != nil {
		return f.Image.ColorModel()
	}

	return f.Info.ColorModel
}
