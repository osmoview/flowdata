package flowdata

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/oliamb/cutter"
	"golang.org/x/image/draw"
)

var DefaultFrameTimeLayout = "2006-01-02_15:04:05"
var HeadersFrameTimeLayout = "2006-01-02 15:04:05.999999999 -0700 MST"

var (
	ErrNoFrameData  = errors.New("no frame data")
	ErrUnknownImage = errors.New("unknown image format")
)

var (
	jpegStartSignature = []byte{0xff, 0xd8}
	jpegEndSignature   = []byte{0xff, 0xd9}

	pngStartSignature = []byte{0x89, 0x50, 0x4E, 0x47}
)

const (
	ImageFormatJPEG = "jpeg"
	ImageFormatPNG  = "png"
	ImageFormatRAW  = "raw"
)

const (
	ContentTypeImageJPEG = "image/jpeg"
	ContentTypeImagePNG  = "image/png"
	ContentTypeImageRAW  = "application/octet-stream"
)

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

	// datasource contains information about the datasource that released this frame
	DataSource SourceInfo `json:"datasource,omitzero"`

	// source where frame incoming, image or video filename, rtsp stream etc
	Source string `json:"source,omitempty"`

	// video position in ms
	Position int `json:"position"`

	// original image, exist for crop images
	Original *Frame `json:"-"`

	// if image based on the some bbox, it required to restore objects bboxes from the crop imaged to original
	BBox BBox `json:"bbox,omitzero"`

	// mutex required to avoid double reencoding with simultaneous requests
	mutex *sync.Mutex
}

type ImageInfo struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`

	colorModel color.Model

	// image format, default: jpeg
	Format string
}

// SourceInfo contins some basic information about datasource, flow and other upper entity
type SourceInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (si SourceInfo) IsZero() bool {
	return si.ID == ""
}

func NewFrameBytes(data []byte) (*Frame, error) {
	info, err := resolveImageFormat(data)
	return &Frame{
		ID:    uuid.NewString(),
		Time:  time.Now(),
		Bytes: data,
		Info:  info,
		mutex: &sync.Mutex{},
	}, err
}

func resolveImageFormat(data []byte) (info ImageInfo, err error) {
	var imgcfg image.Config
	// jpeg
	if bytes.HasPrefix(data, jpegStartSignature) && bytes.HasSuffix(data, jpegEndSignature) {
		info.Format = ImageFormatJPEG

		imgcfg, err = DecodeJPEGInfo(data)
		if err != nil {
			return
		}
	} else if bytes.HasPrefix(data, pngStartSignature) {
		info.Format = ImageFormatPNG

		imgcfg, err = png.DecodeConfig(bytes.NewReader(data))
		if err != nil {
			return
		}
	} else {
		return info, ErrUnknownImage
	}

	info.Width = imgcfg.Width
	info.Height = imgcfg.Height
	info.colorModel = imgcfg.ColorModel

	return
}

func NewFrameImage(img image.Image) (*Frame, error) {
	b := img.Bounds().Max

	return &Frame{
		ID:    uuid.NewString(),
		Time:  time.Now(),
		Image: img,
		Info:  ImageInfo{Width: b.X, Height: b.Y, colorModel: img.ColorModel(), Format: ImageFormatJPEG},
		mutex: &sync.Mutex{},
	}, nil
}

func NewFrameRAW(data []byte, width, height int) *Frame {
	return &Frame{
		ID:    uuid.NewString(),
		Time:  time.Now(),
		Bytes: data,
		Info:  ImageInfo{Width: width, Height: height, Format: ImageFormatRAW},
		mutex: &sync.Mutex{},
	}
}

func NewFrameRGB(data []byte, width, height int) (*Frame, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	if pxcap := width * height * 3; pxcap != len(data) {
		return nil, fmt.Errorf("invalid resolution or data size %d != %d", pxcap, len(data))
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := y*width + x
			img.Pix[i] = data[i]     // R
			img.Pix[i+1] = data[i+1] // G
			img.Pix[i+2] = data[i+2] // B
			img.Pix[i+3] = 255       // A
		}
	}

	return &Frame{
		ID:    uuid.NewString(),
		Time:  time.Now(),
		Image: img,
		Info:  ImageInfo{Width: width, Height: height, colorModel: img.ColorModel()},
		mutex: &sync.Mutex{},
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
	return fmt.Sprintf("id: %s, %d bytes, image: %v", f.ID, len(f.Bytes), f.Image != nil)
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
		max := f.Image.Bounds().Max
		min := f.Image.Bounds().Min
		return max.X - min.X, max.Y - min.Y
	}

	return f.Info.Width, f.Info.Height
}

// GetBytes return slice from `bytes` field, if it exists, overwise encode image to the jpeg
func (f *Frame) GetBytes() (imgdata []byte, err error) {
	if f.mutex != nil {
		f.mutex.Lock()
		defer func() {
			f.mutex.Unlock()
		}()
	}

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

	img, err := f.decodeBytes()
	if err != nil {
		return color.Black
	}

	// assign image to the frame to avoid decoding on each pixel
	f.Image = img

	return f.Image.At(x, y)
}

func (f *Frame) decodeBytes() (img image.Image, err error) {
	if len(f.Bytes) == 0 {
		return f, errors.New("frame does not contain bytes")
	}

	// use mutex only for format required decoding
	if f.Info.Format == ImageFormatJPEG || f.Info.Format == ImageFormatPNG {
		if f.mutex != nil {
			f.mutex.Lock()
			defer f.mutex.Unlock()
		}
	}

	// if bytes exists, decode it as jpeg and return color from obtained image.Image

	switch f.Info.Format {
	case ImageFormatJPEG:
		return DecodeJPEG(f.Bytes)
	case ImageFormatPNG:
		return png.Decode(bytes.NewReader(f.Bytes))
	case ImageFormatRAW:
		// impossible to decode RAW images
		return f, nil
	default:
		// log.Println("image format not specified")
		// use jpeg decoder if format not specified
		return DecodeJPEG(f.Bytes)
	}
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
	defer func() {
		if frame != nil {
			frame.Time = f.Time
			frame.Source = f.Source
			frame.Position = f.Position
			frame.DataSource = f.DataSource
		}
	}()

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
		// Options: cutter.Copy,
	})
	if err != nil {
		return nil, err
	}

	return NewFrameImage(bboxImage)
}

// resize from with save ratio to required width
// NearestNeighbor - very ugly result
// ApproxBiLinear - same ugly, and very slow
// BiLinear - good result and fast
// BenchmarkResizeNearestNeighbor-32    	     457	   2601221 ns/op
// BenchmarkResizeApproxBiLinear-32     	     126	   9820116 ns/op
// BenchmarkResizeBiLinear-32           	     303	   4033383 ns/op
// BenchmarkResizeCatmullRom-32         	     217	   5418807 ns/op
func (f *Frame) Resize(width int) (*Frame, error) {
	// if no image, then check bytes is exists and it can be encoded to the image
	// sometimes it is not possible, like JPEG images without Huffman table
	if f.Image == nil {
		if _, err := f.decodeBytes(); err != nil {
			// if not possible to resize, returns original frame
			return f, err
		}
	}

	height := int(float64(width) / (float64(f.Info.Width) / float64(f.Info.Height)))

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.BiLinear.Scale(img, img.Rect, f, f.Bounds(), draw.Over, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		return nil, err
	}

	return NewFrameImage(img)
}

type RespContext interface {
	Set(key, val string)
	Send([]byte) error
}

type FrameWriteOptions struct {
	Resize int
}

func (f *Frame) WriteResponse(c RespContext, opt FrameWriteOptions) (err error) {
	if opt.Resize > 0 && f.Info.Width > opt.Resize {
		if resizedFrame, err := f.Resize(opt.Resize); err == nil {
			f = resizedFrame
		}
	}

	imgBytes, err := f.GetBytes()
	if err != nil {
		return err
	}

	c.Set("Content-Type", ContentTypeImageJPEG)
	c.Set("Content-Length", strconv.Itoa(len(imgBytes)))

	c.Set("X-DateTime", f.Time.Format(HeadersFrameTimeLayout))
	c.Set("X-Timestamp_ms", strconv.Itoa(int(f.Time.UnixMilli())))
	c.Set("X-Position", strconv.Itoa(f.Position))
	c.Set("X-Source", f.Source)
	c.Set("X-Size", fmt.Sprintf("%dx%d", f.Info.Width, f.Info.Height))
	c.Set("X-ID", f.ID)

	return c.Send(imgBytes)
}

func (f *Frame) FormDataFile() (filedata []byte, filename, contentType string, err error) {
	filedata, err = f.GetBytes()
	return filedata, f.ID, ContentTypeImageJPEG, err
}

func (f *Frame) MarshalJSON() ([]byte, error) {
	if _, err := f.GetBytes(); err != nil {
		return nil, err
	}

	return json.Marshal(struct{ Frame }{*f})
}

func (f *Frame) UnmarshalJSON(b []byte) error {
	if len(b) <= 3 {
		return fmt.Errorf("invalid incoming data")
	}

	// json
	if b[0] == '{' && b[len(b)-1] == '}' {
		type frame Frame
		var frameAlias struct{ frame }
		if err := json.Unmarshal(b, &frameAlias); err != nil {
			return err
		}

		*f = Frame(frameAlias.frame)
		return nil
	}

	// trim quotes if it exists
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = bytes.Trim(b, `"`)
	}

	// try to decode bytes from base64 string
	dec := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(b))
	if raw, err := io.ReadAll(dec); err == nil {
		b = raw
	}

	frame, err := NewFrameBytes(b)
	if err != nil {
		return err
	}

	*f = *frame
	return nil
}

func (f *Frame) ContentType() string {
	switch f.Info.Format {
	case ImageFormatJPEG:
		return ContentTypeImageJPEG
	case ImageFormatPNG:
		return ContentTypeImagePNG
	case ImageFormatRAW:
		return ContentTypeImageRAW

	// WARN: if format unknown returns as jpeg
	default:
		return ContentTypeImageJPEG
	}
}
