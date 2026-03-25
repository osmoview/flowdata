package flowdata

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/hetiansu5/urlquery"
)

const (
	ContentTypeJSON      = "application/json"
	ContentTypeMultipart = "multipart/form-data"
)

type Request struct {
	Method      string
	URL         string
	Query       map[string]any
	Headers     map[string]string
	Body        any
	ContentType string
	Stream      bool
}

func NewRequest() *Request {
	return &Request{
		Query:   map[string]any{},
		Headers: map[string]string{},
	}
}

func (r *Request) SetBody(body any) *Request {
	r.Body = body
	return r
}

func (r *Request) SetQuery(query map[string]any) *Request {
	r.Query = query
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.Headers[k] = v
	}
	return r
}

func (r *Request) GET(url string) (*Response, error) {
	r.Method = http.MethodGet
	r.URL = url
	return r.Do()
}

func (r *Request) POST(url string, body any) (*Response, error) {
	r.Method = http.MethodPost
	r.URL = url
	r.ContentType = ContentTypeJSON
	if body != nil {
		r.Body = body
	}
	return r.Do()
}

func (r *Request) Upload(url string, body any) (*Response, error) {
	r.Method = http.MethodPost
	r.URL = url
	r.ContentType = ContentTypeMultipart
	if body != nil {
		r.Body = body
	}
	return r.Do()
}

func (r *Request) Do() (*Response, error) {
	// prepare body
	var buf bytes.Buffer
	if r.Body != nil {
		if err := r.prepareBody(&buf); err != nil {
			return nil, err
		}
	}

	// create new request
	req, err := http.NewRequest(r.Method, r.URL, &buf)
	if err != nil {
		return nil, err
	}

	// set query parameters if it exists
	if len(r.Query) > 0 {
		q, err := urlquery.Marshal(r.Query)
		if err != nil {
			return nil, err
		}

		if len(req.URL.RawQuery) > 0 && req.URL.RawQuery[len(req.URL.RawQuery)-1] != '&' {
			req.URL.RawQuery += "&"
		}

		req.URL.RawQuery += string(q)
	}

	// set headers
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	// do request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &Response{Response: resp}, err
	}

	// do not read and decode body for stream responses
	if r.Stream {
		if resp.StatusCode >= 400 {
			err = &HTTPError{Code: resp.StatusCode, Message: resp.Status}
		}
		return &Response{Response: resp}, err
	}

	// read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{resp, body}, err
	}

	// handle error
	if resp.StatusCode >= 400 {
		err := &HTTPError{Code: resp.StatusCode, Message: resp.Status}
		if e := json.Unmarshal(body, err); e != nil {
			err.Message = string(body)
		}
		return nil, err
	}

	return &Response{resp, body}, err
}

func (r *Request) prepareBody(w io.Writer) error {
	switch v := r.Body.(type) {
	case []byte:
		_, err := w.Write(v)
		return err
	case string:
		_, err := w.Write([]byte(v))
		return err
	}

	// prepare body as multipart form data
	if strings.HasPrefix(r.ContentType, ContentTypeMultipart) {
		e := NewMultipartEncoder(w)
		if err := e.Encode(r.Body); err != nil {
			return err
		}
		e.Close()

		r.SetHeaders(map[string]string{"Content-Type": e.FormDataContentType()})
		return nil
	}

	// prepare body as json
	r.SetHeaders(map[string]string{"Content-Type": ContentTypeJSON})
	return json.NewEncoder(w).Encode(r.Body)
}

//
// response
//

type Response struct {
	*http.Response
	Body []byte
}

func (r *Response) JSON(x any) error {
	return json.Unmarshal(r.Body, x)
}

func (r *Response) String() string {
	return string(r.Body)
}
