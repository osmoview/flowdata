package flowdata

import (
	"bytes"
	"testing"
)

func TestMultipartEncoder(t *testing.T) {
	var buf bytes.Buffer
	encoder := NewMultipartEncoder(&buf)

	frame, err := NewFrameBytes(imageBytes)
	if err != nil {
		t.Fatal(err)
	}

	objects := []*Object{
		{Type: "test", Name: "test", BBox: [4]float32{0.1, 0.1, 0.2, 0.2}},
		{Type: "test", Name: "test2", BBox: [4]float32{0.7, 0.7, 0.2, 0.2}},
	}

	regions := []*Region{
		{Name: "region"},
	}

	data := map[string]any{"frame": frame, "objects": objects, "status": "just a some text", "regions": regions}

	if err := encoder.Encode(data); err != nil {
		t.Fatal(err)
	}

	ctype := encoder.FormDataContentType()
	if ctype == "" {
		t.Error(err)
	}
	t.Log(ctype)
}
