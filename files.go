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
	ID        string `json:"id"` // uuid
	Filename  string `json:"filename"`
	Type      string `json:"type"`
	Extension string `json:"extension"`
	Size      uint64 `json:"size"` // size in bytes
	Data      []byte `json:"data,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func (f *File) GetData() []byte {
	return f.Data
}

func (f *File) GetFileName() string {
	return f.Filename
}

func NewFileFromData(data []byte) *File {
	mt := mimetype.Detect(data)
	id := uuid.New()

	return &File{
		ID:        id.String(),
		Filename:  id.String() + mt.Extension(),
		Type:      mt.String(),
		Extension: mt.Extension(),
		Size:      uint64(len(data)),
		Data:      data,
		CreatedAt: time.Now().UTC(),
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
	file := NewFileFromData(data)
	file.Filename = filepath.Base(path)

	return file, nil
}
