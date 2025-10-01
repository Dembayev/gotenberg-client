// Package gotenberg provides a client for the Gotenberg service.
// It offers a convenient API for converting HTML and URLs to PDF documents.
package gotenberg

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	httpclient "github.com/nativebpm/http-client"
	"github.com/nativebpm/http-client/request"
)

// Client is a Gotenberg HTTP client that wraps the base HTTP client
// with Gotenberg-specific functionality for document conversion.
type Client struct {
	*httpclient.Client
}

// Request represents a Gotenberg conversion request builder.
// It wraps the underlying multipart request and provides Gotenberg-specific methods.
type Request struct {
	req *request.Multipart
	wh  map[string]string
}

// Response represents a Gotenberg conversion response.
// It wraps the HTTP response and provides access to the Gotenberg trace header.
type Response struct {
	*http.Response
	GotenbergTrace string
}

// NewClient creates a new Gotenberg client with the given HTTP client and base URL.
// Returns an error if the base URL is invalid.
func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
	client, err := httpclient.NewClient(httpClient, baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: client,
	}, nil
}

// ConvertHTML creates a request to convert HTML content to PDF.
// The html parameter should contain the HTML content to be converted.
func (c *Client) ConvertHTML(ctx context.Context, html io.Reader) *Request {
	r := &Request{}
	r.req = c.MultipartPOST(ctx, ConvertHTML).File(FieldFiles, FileIndexHTML, html)
	return r
}

// ConvertURL creates a request to convert a web page at the given URL to PDF.
func (c *Client) ConvertURL(ctx context.Context, url string) *Request {
	r := &Request{}
	r.req = c.MultipartPOST(ctx, ConvertURL).Param(FieldURL, url)
	return r
}

// Send executes the conversion request and returns the response.
// Returns an error if the request fails or the conversion cannot be completed.
func (r *Request) Send() (*Response, error) {
	resp, err := r.req.Send()
	if err != nil {
		return nil, err
	}
	return &Response{
		Response:       resp,
		GotenbergTrace: resp.Header.Get(HeaderGotenbergTrace),
	}, nil
}

// Header adds a header to the conversion request.
func (r *Request) Header(key, value string) *Request {
	r.req.Header(key, value)
	return r
}

// Param adds a form parameter to the conversion request.
func (r *Request) Param(key, value string) *Request {
	r.req.Param(key, value)
	return r
}

// Bool adds a boolean form parameter to the conversion request.
func (r *Request) Bool(fieldName string, value bool) *Request {
	r.req.Bool(fieldName, value)
	return r
}

// Float adds a float64 form parameter to the conversion request.
func (r *Request) Float(fieldName string, value float64) *Request {
	r.req.Float(fieldName, value)
	return r
}

// File adds a file to the conversion request.
func (r *Request) File(key, filename string, content io.Reader) *Request {
	r.req.File(key, filename, content)
	return r
}

// WebhookURL sets the webhook URL and HTTP method for successful conversions.
func (r *Request) WebhookURL(url, method string) *Request {
	r.req.Header(HeaderWebhookURL, url).
		Header(HeaderWebhookMethod, method)
	return r
}

// WebhookErrorURL sets the webhook URL and HTTP method for failed conversions.
func (r *Request) WebhookErrorURL(url, method string) *Request {
	r.req.Header(HeaderWebhookErrorURL, url).
		Header(HeaderWebhookErrorMethod, method)
	return r
}

// WebhookHeader adds a custom header to be sent with webhook requests.
// Multiple headers can be added by calling this method multiple times.
func (r *Request) WebhookHeader(key, value string) *Request {
	if r.wh == nil {
		r.wh = make(map[string]string)
	}

	r.wh[key] = value
	webhookHeaders, _ := json.Marshal(r.wh)
	r.req.Header(HeaderWebhookExtraHTTPHeaders, string(webhookHeaders))
	return r
}

// OutputFilename sets the output filename for the generated PDF.
func (r *Request) OutputFilename(filename string) *Request {
	r.req.Header(HeaderOutputFilename, filename)
	return r
}

// PaperSize sets the paper size for the PDF using width and height in inches.
func (r *Request) PaperSize(width, height float64) *Request {
	r.req.Float(FieldPaperWidth, width)
	r.req.Float(FieldPaperHeight, height)
	return r
}

// PaperSizeA4 sets the paper size to A4 format.
func (r *Request) PaperSizeA4() *Request {
	return r.PaperSize(PaperSizeA4[0], PaperSizeA4[1])
}

// PaperSizeA6 sets the paper size to A6 format.
func (r *Request) PaperSizeA6() *Request {
	return r.PaperSize(PaperSizeA6[0], PaperSizeA6[1])
}

// PaperSizeLetter sets the paper size to Letter format.
func (r *Request) PaperSizeLetter() *Request {
	return r.PaperSize(PaperSizeLetter[0], PaperSizeLetter[1])
}

// Margins sets the page margins for the PDF in inches.
// Parameters are in order: top, right, bottom, left.
func (r *Request) Margins(top, right, bottom, left float64) *Request {
	r.req.Float(FieldMarginTop, top)
	r.req.Float(FieldMarginRight, right)
	r.req.Float(FieldMarginBottom, bottom)
	r.req.Float(FieldMarginLeft, left)
	return r
}
