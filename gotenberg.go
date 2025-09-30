package gotenberg

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	httpclient "github.com/nativebpm/http-client"
	"github.com/nativebpm/http-client/request"
)

type Client struct {
	*httpclient.Client
}

type Request struct {
	req *request.Multipart
	wh  map[string]string
	err error
}

type Response struct {
	*http.Response
	GotenbergTrace string
}

func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
	client, err := httpclient.NewClient(httpClient, baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: client,
	}, nil
}

func (c *Client) ConvertHTML(ctx context.Context, html io.Reader) *Request {
	r := &Request{}
	r.req = c.MultipartPOST(ctx, ConvertHTML).File(FieldFiles, FileIndexHTML, html)
	r.err = r.req.Err()
	return r
}

func (c *Client) ConvertURL(ctx context.Context, url string) *Request {
	r := &Request{}
	r.req = c.MultipartPOST(ctx, ConvertURL).Param(FieldURL, url)
	r.err = r.req.Err()
	return r
}

func (r *Request) Err() error {
	return r.err
}

func (r *Request) Send() (*Response, error) {
	if r.err != nil {
		return nil, r.err
	}

	resp, err := r.req.Send()
	if err != nil {
		return nil, err
	}
	return &Response{
		Response:       resp,
		GotenbergTrace: resp.Header.Get(HeaderGotenbergTrace),
	}, nil
}

func (r *Request) Header(key, value string) *Request {
	if r.err != nil {
		return r
	}

	r.req.Header(key, value)
	return r
}

func (r *Request) Param(key, value string) *Request {
	if r.err != nil {
		return r
	}

	r.req.Param(key, value)
	return r
}

func (r *Request) Bool(fieldName string, value bool) *Request {
	if r.err != nil {
		return r
	}

	r.req.Bool(fieldName, value)
	return r
}

func (r *Request) Float(fieldName string, value float64) *Request {
	if r.err != nil {
		return r
	}

	r.req.Float(fieldName, value)
	return r
}

func (r *Request) File(key, filename string, content io.Reader) *Request {
	if r.err != nil {
		return r
	}

	r.req.File(key, filename, content)
	return r
}

func (r *Request) WebhookURL(url, method string) *Request {
	if r.err != nil {
		return r
	}

	r.req.Header(HeaderWebhookURL, url).Header(HeaderWebhookMethod, method)
	r.err = r.req.Err()
	return r
}

func (r *Request) OutputFilename(filename string) *Request {
	if r.err != nil {
		return r
	}

	r.req.Header(HeaderOutputFilename, filename)
	r.err = r.req.Err()
	return r
}

func (r *Request) WebhookErrorURL(url, method string) *Request {
	if r.err != nil {
		return r
	}

	r.req.Header(HeaderWebhookErrorURL, url).Header(HeaderWebhookErrorMethod, method)
	r.err = r.req.Err()
	return r
}

func (r *Request) WebhookHeader(key, value string) *Request {
	if r.err != nil {
		return r
	}

	if r.wh == nil {
		r.wh = make(map[string]string)
	}

	r.wh[key] = value
	webhookHeaders, err := json.Marshal(r.wh)
	if err != nil {
		r.err = err
		return r
	}

	r.req.Header(HeaderWebhookExtraHTTPHeaders, string(webhookHeaders))
	r.err = r.req.Err()
	return r
}

func (r *Request) PaperSize(width, height float64) *Request {
	if r.err != nil {
		return r
	}

	r.req.Float(FieldPaperWidth, width)
	r.req.Float(FieldPaperHeight, height)
	return r
}

func (r *Request) PaperSizeA4() *Request {
	return r.PaperSize(PaperSizeA4[0], PaperSizeA4[1])
}

func (r *Request) PaperSizeA6() *Request {
	return r.PaperSize(PaperSizeA6[0], PaperSizeA6[1])
}

func (r *Request) PaperSizeLetter() *Request {
	return r.PaperSize(PaperSizeLetter[0], PaperSizeLetter[1])
}

func (r *Request) Margins(top, right, bottom, left float64) *Request {
	if r.err != nil {
		return r
	}

	r.req.Float(FieldMarginTop, top)
	r.req.Float(FieldMarginRight, right)
	r.req.Float(FieldMarginBottom, bottom)
	r.req.Float(FieldMarginLeft, left)
	return r
}
