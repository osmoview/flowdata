package flowdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"nodes/gopath/src/github.com/gofiber/fiber/v3/log"
	"reflect"
	"strconv"
	"strings"
)

var ErrUnsupportedVariableType = errors.New("unsupported variable type")

//
// http error
//

// HTTPError required to handle resonse in the json format, usually it contains message to text field with error description
type HTTPError struct {
	Text string `json:"text"`

	// code and message required for fiber compatible
	Code    int    `json:"code"`
	Message string `json:"message"`

	// detail required for compatible with fast api 400+ responses
	Detail []FastAPIDetail `json:"detail"`
}

// FastAPIDetail format,
// example: {"detail":[{"type":"missing","loc":["body","events"],"msg":"Field required","input":{"Events":[]}}]}
type FastAPIDetail struct {
	Type string   `json:"type"`
	Loc  []string `json:"loc"`
	Msg  string   `json:"msg"`
}

func (e *HTTPError) Error() string {
	if e.Text != "" {
		return e.Text
	}

	if len(e.Detail) > 0 {
		msgs := make([]string, 0, len(e.Detail))
		for _, d := range e.Detail {
			msgs = append(msgs, fmt.Sprintf("'%s' %s", strings.Join(d.Loc, "."), strings.ToLower(d.Msg)))
		}

		return strings.Join(msgs, "\n")
	}

	return e.Message
}

//
//
//

// FormDataFile is interface describe entitiy like *Frame suitable to sending by http as part of multipart/form-data
type FormDataFile interface {
	FormDataFile() (filedata []byte, filename string, contentType string, err error)
}

//
//
//

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

		// *flowdata.Frame implement FormDataFile interface, but golang cannot properly handle
		// slice of structs implements some interface, so checks it by reflect, sorry
		if rv, ok := e.isSliceOfFormDataFile(v); ok {
			for i := 0; i < rv.Len(); i++ {
				v := rv.Index(i)

				// if slice elem is match to FormDataFile interface
				if f, ok := v.Interface().(FormDataFile); ok {
					if err := e.WriteFile(k, f); err != nil {
						return err
					}
				}

				// if slice elem is convertable to the flowdata.File
				if f, ok := e.isFileCompatible(v); ok {
					if err := e.WriteFile(k, f); err != nil {
						return err
					}
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
			if file, ok := e.isFileCompatible(f); ok {
				return e.WriteFile(k, file)
			}

			data, err := json.Marshal(f)
			if err != nil {
				return err
			}
			if err := e.WritePart(k, data, "", "application/json"); err != nil {
				return err
			}
		}

	}

	return nil
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

	if _, ok := elem.Interface().(FormDataFile); ok {
		return rv, ok
	}

	if _, ok := e.isFileCompatible(elem); ok {
		return rv, ok
	}

	return reflect.Value{}, false
}

func (e *MultipartEncoder) isFileCompatible(v any) (file *File, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
			log.Error("failed to convert struct or map to the flowdata.File", r)
			return
		}
	}()

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// for structure check fields exists
	if rv.Kind() == reflect.Struct {
		fbytes := rv.FieldByName("Bytes")
		fname := rv.FieldByName("Name")
		if !fbytes.IsValid() || !fname.IsValid() {
			return nil, false
		}

		return NewFile(fbytes.Bytes(), fname.String()), true
	}

	// for the map check map keys exists
	if rv.Kind() == reflect.Map {
		var fbytes reflect.Value
		var fname reflect.Value

		for _, key := range rv.MapKeys() {
			if key.Kind() != reflect.String {
				continue
			}

			if strings.EqualFold(key.String(), "bytes") {
				fbytes = rv.MapIndex(key)
			} else if strings.EqualFold(key.String(), "name") {
				fname = rv.MapIndex(key)
			}
		}

		return NewFile(fbytes.Bytes(), fname.String()), true
	}

	return
}
