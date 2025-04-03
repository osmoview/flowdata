package flowdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"reflect"
	"strconv"
	"strings"
)

//
// http error
//

// HTTPError required to handle resonse in the json format, usually it contains message to text field with error description
type HTTPError struct {
	Code    int
	Message string
	Text    string
}

func (e *HTTPError) Error() string {
	if e.Text != "" {
		return e.Text
	}

	return e.Message
}

var ErrUnsupportedVariableType = errors.New("unsupported variable type")

// FormDataFile is interface describe entitiy like *Frame suitable to sending by http as part of multipart/form-data
type FormDataFile interface {
	FormDataFile() (filedata []byte, filename string, contentType string, err error)
}

// MultipartEncoder is just a wrapper aroind *multipart.Writer to more comfortable encode data
type MultipartEncoder struct {
	*multipart.Writer
}

func NewMultipartEncoder(w io.Writer) *MultipartEncoder {
	return &MultipartEncoder{multipart.NewWriter(w)}
}

func (e *MultipartEncoder) Encode(v any) error {
	m, ok := v.(map[string]any)
	if !ok {
		return ErrUnsupportedVariableType
	}

	for k, v := range m {

		// i dont want to import `flowdata` package
		// *flowdata.Frame implement FormDataFile interface, but golang cannot properly handle
		// slice of structs implements some interface, so checks it by reflect, sorry
		if rv, ok := e.isSliceOfFormDataFile(v); ok {
			for i := 0; i < rv.Len(); i++ {
				v := rv.Index(i)
				if err := e.WriteFile(k, v.Interface().(FormDataFile)); err != nil {
					return err
				}
			}
			continue
		}

		switch f := v.(type) {
		// upload multiple frames/files to the remote service
		case []FormDataFile:
			for _, frame := range f {
				if err := e.WriteFile(k, frame); err != nil {
					return err
				}
			}

		// upload one frame/file to the remote service
		case FormDataFile:
			if err := e.WriteFile(k, f); err != nil {
				return err
			}

		// write part as bytes
		case []byte:
			if err := e.WritePart(k, f, "", "application/octet-stream"); err != nil {
				return err
			}

		// write part as string
		case string:
			if err := e.WriteField(k, f); err != nil {
				return err
			}

		// write part as json
		default:
			data, err := json.Marshal(f)
			if err != nil {
				return err
			}
			if err := e.WritePart(k, data, "", "application/json"); err != nil {
				return err
			}
		}

	}

	return e.Close()
}

func (e *MultipartEncoder) GetPartHeader(partname string, data []byte, filename, contentType string) (h textproto.MIMEHeader) {
	h = make(textproto.MIMEHeader)

	if contentType != "" {
		h.Set("Content-Type", contentType)
	}

	contentLength := strconv.Itoa(len(data))
	h.Set("Content-Length", contentLength)

	contentDisposition := []string{"form-data"}
	if partname != "" {
		contentDisposition = append(contentDisposition, fmt.Sprintf(`name="%s"`, partname))
	}
	if filename != "" {
		contentDisposition = append(contentDisposition, fmt.Sprintf(`filename="%s"`, filename))
	}

	h.Set("Content-Disposition", strings.Join(contentDisposition, "; "))

	return h
}

func (e *MultipartEncoder) WritePart(partname string, data []byte, filename, contentType string) error {
	part, err := e.CreatePart(e.GetPartHeader(partname, data, filename, contentType))
	if err != nil {
		return err
	}

	if _, err := part.Write(data); err != nil {
		return err
	}

	return nil
}

func (e *MultipartEncoder) WriteFile(partname string, f FormDataFile) error {
	filedata, filename, contentType, err := f.FormDataFile()
	if err != nil {
		return err
	}

	return e.WritePart(partname, filedata, filename, contentType)
}

func (e *MultipartEncoder) isSliceOfFormDataFile(v any) (reflect.Value, bool) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		return reflect.Value{}, false
	}

	// WARN: maybe required additional checks
	elem := reflect.New(rv.Type().Elem()).Elem()

	_, ok := elem.Interface().(FormDataFile)
	return rv, ok
}
