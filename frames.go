package flowdata

import (
	"image"
	"time"
)

// Frame contains image data and time when that image was received
type Frame struct {
	Time time.Time
	Img  image.Image
}

// GetImage returns image.Image
func (f Frame) GetImage() image.Image {
	return f.Img
}

// GetTime returns time when image received
func (f Frame) GetTime() time.Time {
	return f.Time
}

// GetSize returns image resolution
func (f Frame) GetSize() (w, h int) {
	b := f.Img.Bounds()
	return b.Max.X, b.Max.Y
}
