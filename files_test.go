package flowdata

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileFromBytes(t *testing.T) {
	data := []byte("hello world")
	f := NewFileFromBytes(data)

	if f.ID == "" {
		t.Error("empty ID")
	}
	if f.Size != uint64(len(data)) {
		t.Errorf("invalid size: got %d, want %d", f.Size, len(data))
	}
	if string(f.Bytes) != string(data) {
		t.Error("bytes mismatch")
	}
	if f.Type == "" {
		t.Error("empty type")
	}
	if f.ReadAt.IsZero() {
		t.Error("ReadAt is zero")
	}
}

func TestNewFileFromPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	data := []byte("file contents")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	f, err := NewFileFromPath(path)
	if err != nil {
		t.Fatal(err)
	}
	if f.Name != "sample.txt" {
		t.Errorf("invalid name: %q", f.Name)
	}
	if f.Size != uint64(len(data)) {
		t.Errorf("invalid size: %d", f.Size)
	}
	if string(f.Bytes) != string(data) {
		t.Error("bytes mismatch")
	}
}

func TestFileFormDataFile(t *testing.T) {
	f := NewFileFromBytes([]byte("abc"))
	f.Name = "x.bin"

	data, name, ct, err := f.FormDataFile()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "abc" {
		t.Error("data mismatch")
	}
	if name != "x.bin" {
		t.Errorf("name: %q", name)
	}
	if ct != f.Type {
		t.Errorf("content type: %q", ct)
	}
}

func TestFile(t *testing.T) {
	f, err := NewFileFromPath("testdata/logo.jpg")
	if err != nil {
		t.Fatal(err)
	}

	if f.Name != "logo.jpg" {
		t.Errorf("invalid name: %q", f.Name)
	}
	if f.Type != "image/jpeg" {
		t.Errorf("invalid type: %q", f.Type)
	}
	if f.Ext != ".jpg" {
		t.Errorf("invalid ext: %q", f.Ext)
	}
	if f.Size == 0 || f.Size != uint64(len(f.Bytes)) {
		t.Errorf("invalid size: %d", f.Size)
	}
	if f.ID == "" {
		t.Error("empty ID")
	}
}
