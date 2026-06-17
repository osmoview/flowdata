package flowdata

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
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
	Path  string `json:"path"`
	Type  string `json:"type"`
	Ext   string `json:"ext"`
	Size  uint64 `json:"size"` // size in bytes
	Bytes []byte `json:"bytes,omitempty"`

	ReadAt time.Time `json:"read_at"`
}

// NewFile returns `File` structure with passed bytes and filename
// it tries to resolve the file type(mime) by inspect bytes
func NewFile(bytes []byte, fp string) *File {
	file := &File{
		ID:     uuid.NewString(),
		Name:   filepath.Base(fp),
		Path:   fp,
		Ext:    strings.ToLower(filepath.Ext(fp)),
		Size:   uint64(len(bytes)),
		Bytes:  bytes,
		ReadAt: time.Now().UTC(),
	}

	// if extension specified, resolve mime type from it
	if file.Ext != "" {
		file.Type = mime.TypeByExtension(file.Ext)
	}

	// if failed to resolve file type from extension,
	// try to obtain it from bytes
	if file.Type == "" {
		mt := mimetype.Detect(bytes)
		file.Type = mt.String()
		if file.Ext == "" {
			file.Ext = mt.Extension()
		}
	}

	return file
}

// ReadFile return the `File` structure with bytes readed by specified file path
func ReadFile(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewFile(data, path), nil
}

// FromDataFile required to match to FormDataFile interface to write file as part of multipart http request
func (f *File) FormDataFile() (filedata []byte, filename, contentType string, err error) {
	return f.Bytes, f.Name, f.Type, nil
}

func (f *File) UnmarshalJSON(b []byte) error {
	if len(b) <= 3 {
		return fmt.Errorf("invalid incoming data")
	}

	// json
	if b[0] == '{' && b[len(b)-1] == '}' {
		type file File
		var fileAlias struct{ file }
		if err := json.Unmarshal(b, &fileAlias); err != nil {
			return err
		}

		*f = File(fileAlias.file)
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

	file := NewFile(b, "")

	*f = *file
	return nil
}
