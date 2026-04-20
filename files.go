package flowdata

import (
	"io"
	"os"
	"path/filepath"
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

func (f *File) GetBytes() []byte {
	return f.Bytes
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) FormDataFile() (filedata []byte, filename, contentType string, err error) {
	return f.GetBytes(), f.GetName(), f.Type, nil
}

func NewFileFromBytes(bytes []byte) *File {
	mt := mimetype.Detect(bytes)
	id := uuid.New()

	return &File{
		ID:     id.String(),
		Type:   mt.String(),
		Ext:    mt.Extension(),
		Size:   uint64(len(bytes)),
		Bytes:  bytes,
		ReadAt: time.Now().UTC(),
	}
}

// NewFileFromPath return file data on path
func NewFileFromPath(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	file := NewFileFromBytes(data)
	file.Name = filepath.Base(path)

	return file, nil
}
