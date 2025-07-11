package flowdata

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/go-querystring/query"
)

type Request struct {
	Method      string
	URL         string
	Query       map[string]any
	Headers     map[string]string
	Body        any
	ContentType string
}

func NewRequest() *Request {
	return &Request{}
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
	r.Headers = headers
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
		q, err := query.Values(r.Query)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = q.Encode()
	}

	// set headers
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	// do request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// handle error
	if resp.StatusCode >= 400 {
		err := &HTTPError{Code: resp.StatusCode, Message: resp.Status}
		if e := json.Unmarshal(body, err); e != nil {
			err.Message = string(body)
		}
		return nil, err
	}

	return &Response{resp, body}, nil
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
