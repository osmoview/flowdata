package flowdata

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
)

type File struct {
	ID    string `json:"id"` // uuid
	Name  string `json:"name"`
	Type  string `json:"type"`
	Ext   string `json:"ext"`
	Size  uint64 `json:"size"` // size in bytes
	Bytes []byte `json:"bytes,omitempty"`

	ReadAt time.Time `json:"read_at"`
}

// NewFile returns `File` structure with passed bytes and filename
// it tries to resolve the file type(mime) by inspect bytes
func NewFile(bytes []byte, name string) *File {
	// try to resolve mime type from the bytes
	mt := mimetype.Detect(bytes)

	// resolve extenstion from the file name,
	// if failed, resolve from the detected mime type
	ext := strings.ToLower(filepath.Ext(name))
	if ext == "" {
		ext = mt.Extension()
	}

	return &File{
		ID:     uuid.NewString(),
		Name:   name,
		Type:   mt.String(),
		Ext:    ext,
		Size:   uint64(len(bytes)),
		Bytes:  bytes,
		ReadAt: time.Now().UTC(),
	}
}

// ReadFile return the `File` structure with bytes readed by specified file path
func ReadFile(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewFile(data, filepath.Base(path)), nil
}

func (f *File) GetBytes() []byte {
	return f.Bytes
}

func (f *File) GetName() string {
	return f.Name
}

// FromDataFile required to match to FormDataFile interface to write file as part of multipart http request
func (f *File) FormDataFile() (filedata []byte, filename, contentType string, err error) {
	return f.GetBytes(), f.GetName(), f.Type, nil
}
