package gotenberg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

const (
	bufferSize = 1 << 12 // 4096 bytes (4 KB)
)

const (
	ApplicationJSON = "application/json"
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
)

type Request struct {
	request   *http.Request
	multipart bool
	writer    *multipart.Writer
	buffer    *bytes.Buffer
	bufPool   sync.Pool
}

type Client struct {
	baseURL *url.URL
	client  *http.Client
	err     error
	Request
}

func NewClient(client *http.Client, baseURL string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	return &Client{
		client:  client,
		baseURL: u,
	}, nil
}

func (r *Client) MethodGet(ctx context.Context, path string) *Client {
	if r.err != nil {
		return r
	}

	r.request, r.err = http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL.JoinPath(path).String(), nil)
	return r
}

func (r *Client) MethodPost(ctx context.Context, path string) *Client {
	if r.err != nil {
		return r
	}

	r.request, r.err = http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL.JoinPath(path).String(), nil)
	return r
}

func (r *Client) MethodPut(ctx context.Context, path string) *Client {
	if r.err != nil {
		return r
	}

	r.request, r.err = http.NewRequestWithContext(ctx, http.MethodPut, r.baseURL.JoinPath(path).String(), nil)
	return r
}

func (r *Client) MethodPatch(ctx context.Context, path string) *Client {
	if r.err != nil {
		return r
	}

	r.request, r.err = http.NewRequestWithContext(ctx, http.MethodPatch, r.baseURL.JoinPath(path).String(), nil)
	return r
}

func (r *Client) MethodDelete(ctx context.Context, path string) *Client {
	if r.err != nil {
		return r
	}

	r.request, r.err = http.NewRequestWithContext(ctx, http.MethodDelete, r.baseURL.JoinPath(path).String(), nil)
	return r
}

func (r *Client) Multipart() *Client {
	if r.err != nil {
		return r
	}

	if r.multipart {
		return r
	}

	r.bufPool = sync.Pool{
		New: func() any {
			buf := make([]byte, 0, bufferSize)
			return &buf
		},
	}
	r.buffer = &bytes.Buffer{}
	r.buffer.Grow(bufferSize)
	r.writer = multipart.NewWriter(r.buffer)
	r.multipart = true

	return r
}

func (r *Client) Header(key, value string) *Client {
	if r.err != nil {
		return r
	}

	if r.request.Header == nil {
		r.request.Header = make(http.Header)
	}
	r.request.Header.Set(key, value)

	return r
}

func (r *Client) Headers(headers map[string]string) *Client {
	if r.err != nil {
		return r
	}

	for key, value := range headers {
		r = r.Header(key, value)
	}

	return r
}

func (r *Client) ContentType(contentType string) *Client {
	return r.Header(ContentType, contentType)
}

func (r *Client) JSONContentType() *Client {
	return r.ContentType(ApplicationJSON)
}

func (r *Client) QueryParam(key, value string) *Client {
	if r.err != nil {
		return r
	}

	q := r.request.URL.Query()
	q.Set(key, value)
	r.request.URL.RawQuery = q.Encode()

	return r
}

func (r *Client) QueryParams(params map[string]string) *Client {
	if r.err != nil {
		return r
	}

	q := r.request.URL.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	r.request.URL.RawQuery = q.Encode()

	return r
}

func (r *Client) QueryValues(values url.Values) *Client {
	if r.err != nil {
		return r
	}

	q := r.request.URL.Query()
	for k := range values {
		v := values.Get(k)
		if v == "" {
			q.Del(k)
			continue
		}
		q.Add(k, values.Get(k))
	}
	r.request.URL.RawQuery = q.Encode()

	return r
}

func (r *Client) Body(body io.ReadCloser) *Client {
	if r.err != nil {
		return r
	}

	r.request.Body = body
	return r
}

func (r *Client) BytesBody(body []byte) *Client {
	if r.err != nil {
		return r
	}

	r.request.Body = io.NopCloser(bytes.NewReader(body))
	r = r.Header(ContentLength, strconv.Itoa(len(body)))

	return r
}

func (r *Client) StringBody(body string) *Client {
	return r.BytesBody([]byte(body))
}

func (r *Client) JSONBody(body any) *Client {
	if r.err != nil {
		return r
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		r.err = fmt.Errorf("failed to marshal JSON: %w", err)
		return r
	}

	r = r.BytesBody(jsonData)
	r = r.JSONContentType()

	return r
}

func (r *Client) File(fieldName, filename string, content io.Reader) *Client {
	if r.err != nil {
		return r
	}

	if !r.multipart {
		r = r.Multipart()
		if r.err != nil {
			return r
		}
	}

	part, err := r.writer.CreateFormFile(fieldName, filename)
	if err != nil {
		r.err = fmt.Errorf("failed to create form file: %w", err)
		return r
	}

	var buf []byte
	if p := r.bufPool.Get(); p != nil {
		buf = (*p.(*[]byte))[:bufferSize]
	} else {
		buf = make([]byte, bufferSize)
	}
	defer func() {
		buf = buf[:0]
		r.bufPool.Put(&buf)
	}()

	_, err = io.CopyBuffer(part, content, buf)
	if err != nil {
		r.err = fmt.Errorf("failed to copy file content: %w", err)
		return r
	}

	return r
}

func (r *Client) FormField(fieldName, value string) *Client {
	if r.err != nil {
		return r
	}

	if !r.multipart {
		r = r.Multipart()
		if r.err != nil {
			return r
		}
	}

	err := r.writer.WriteField(fieldName, value)
	if err != nil {
		r.err = fmt.Errorf("failed to write form field %q: %w", fieldName, err)
		return r
	}

	return r
}

func (r *Client) Err() error {
	return r.err
}

func (r *Client) Reset() *Client {
	r.request = nil
	r.err = nil
	r.multipart = false
	r.writer = nil
	if r.buffer != nil {
		r.buffer.Reset()
	}
	return r
}

func (r *Client) Send() (*http.Response, error) {
	defer r.Reset()

	if r.err != nil {
		return nil, r.err
	}

	if r.multipart && r.writer != nil && r.buffer != nil {
		if err := r.writer.Close(); err != nil {
			return nil, fmt.Errorf("failed to close multipart writer: %w", err)
		}

		r.request.Body = io.NopCloser(r.buffer)
		r = r.Header(ContentType, r.writer.FormDataContentType())
		r = r.Header(ContentLength, strconv.Itoa(r.buffer.Len()))
	}

	return r.client.Do(r.request)
}
