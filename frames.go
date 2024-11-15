package flowdata

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/oliamb/cutter"
)

var DefaultFrameTimeLayout = "2006-01-02_15:04:05"
var HeadersFrameTimeLayout = "2006-01-02 15:04:05.999999999 -0700 MST"

var ErrNoFrameData = errors.New("no frame data")

const ContentTypeImageJPEG = "image/jpeg"

// Frame contains image data and time when that image was received.
// it can contain only `bytes` or `image` or both of it.
// required to be a careful with reassigning an image data
type Frame struct {
	// uniquey frame ID
	ID string `json:"id"`

	// time when frame received
	Time time.Time `json:"time"`

	// presentation of an image in the image.Image interface
	Image image.Image `json:"-"`

	// raw bytes of the image(jpeg usually)
	Bytes []byte `json:"bytes"`

	// image info
	Info ImageInfo `json:"info"`

	// video position in ms
	Position int `json:"position"`

	// original image, exist for crop images
	Original *Frame `json:"-"`

	// if image based on the some bbox, it required to restore objects bboxes from the crop imaged to original
	BBox BBox `json:"-"`
}

type ImageInfo struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`

	colorModel color.Model

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
		ID:    uuid.NewString(),
		Time:  time.Now(),
		Bytes: data,
		Info:  ImageInfo{Width: info.Width, Height: info.Height, colorModel: info.ColorModel},
	}, nil
}

func NewFrameImage(img image.Image) (*Frame, error) {
	b := img.Bounds().Max

	return &Frame{
		ID:    uuid.NewString(),
		Time:  time.Now(),
		Image: img,
		Info:  ImageInfo{Width: b.X, Height: b.Y, colorModel: img.ColorModel()},
	}, nil
}

// GetImage returns image.Image
func (f *Frame) GetImage() image.Image {
	if f.Image != nil {
		return f.Image
	}

	return f
}

func (f *Frame) String() string {
	return fmt.Sprintf("{id: %s, size: %d}", f.ID, len(f.Bytes))
}

// GetID returns the ID. generated when frame was created
func (f *Frame) GetID() string {
	return f.ID
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
func (f *Frame) GetBytes() (imgdata []byte, err error) {
	// if bytes exists then just return it
	if len(f.Bytes) > 0 {
		return f.Bytes, nil
	}

	// if image exists encode it to the jpeg
	if f.Image != nil {
		f.Bytes, err = EncodeJPEG(f.Image)
		return f.Bytes, err
	}

	// if nothing exists, then do not return anything
	return imgdata, ErrNoFrameData
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

	return f.Info.colorModel
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func (f *Frame) Crop(bbox BBox) (frame *Frame, err error) {
	width, height := f.GetSize()

	x := int(bbox[0] * float32(width))
	y := int(bbox[1] * float32(height))
	w := int(bbox[2] * float32(width))
	h := int(bbox[3] * float32(height))

	defer func() {
		if frame != nil {
			frame.Original = f
			frame.BBox = bbox
		}
	}()

	rect := image.Rect(x, y, x+w, y+h)

	frameImage := f.GetImage()

	// if image compatible with SubImager interface
	sub, ok := frameImage.(SubImager)
	if ok {
		return NewFrameImage(sub.SubImage(rect))
	}

	bboxImage, err := cutter.Crop(frameImage, cutter.Config{
		Anchor: image.Point{X: x, Y: y},
		Width:  w,
		Height: h,
	})
	if err != nil {
		return nil, err
	}

	return NewFrameImage(bboxImage)
}

type RespContext interface {
	Set(key, val string)
	Send([]byte) error
}

func (f *Frame) WriteResponse(c RespContext) error {
	bytes, err := f.GetBytes()
	if err != nil {
		return err
	}

	c.Set("Content-Type", ContentTypeImageJPEG)
	c.Set("Content-Length", strconv.Itoa(len(bytes)))

	c.Set("X-DateTime", f.Time.Format(HeadersFrameTimeLayout))
	c.Set("X-Timestamp_ms", strconv.Itoa(int(f.Time.UnixMilli())))
	c.Set("X-Position", strconv.Itoa(f.Position))
	c.Set("X-Size", fmt.Sprintf("%dx%d", f.Info.Width, f.Info.Height))
	c.Set("X-ID", f.ID)

	return c.Send(bytes)
}

func (f *Frame) FormDataFile() (filedata []byte, filename, contentType string, err error) {
	filedata, err = f.GetBytes()
	return filedata, f.ID, ContentTypeImageJPEG, err
}
