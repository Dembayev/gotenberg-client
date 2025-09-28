package gotenberg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nativebpm/gotenberg-client/pkg/httpclient"
)

const (
	ConvertHTML = "/forms/chromium/convert/html"
	ConvertURL  = "/forms/chromium/convert/url"
)

const (
	FieldSinglePage              = "singlePage"
	FieldPaperWidth              = "paperWidth"
	FieldPaperHeight             = "paperHeight"
	FieldMarginTop               = "marginTop"
	FieldMarginBottom            = "marginBottom"
	FieldMarginLeft              = "marginLeft"
	FieldMarginRight             = "marginRight"
	FieldPreferCSSPageSize       = "preferCssPageSize"
	FieldGenerateDocumentOutline = "generateDocumentOutline"
	FieldGenerateTaggedPDF       = "generateTaggedPdf"
	FieldPrintBackground         = "printBackground"
	FieldOmitBackground          = "omitBackground"
	FieldLandscape               = "landscape"
	FieldScale                   = "scale"
	FieldNativePageRanges        = "nativePageRanges"
)

const (
	HeaderWebhookURL              = "Gotenberg-Webhook-Url"
	HeaderWebhookErrorURL         = "Gotenberg-Webhook-Error-Url"
	HeaderWebhookMethod           = "Gotenberg-Webhook-Method"
	HeaderWebhookErrorMethod      = "Gotenberg-Webhook-Error-Method"
	HeaderWebhookExtraHTTPHeaders = "Gotenberg-Webhook-Extra-Http-Headers"
	HeaderOutputFilename          = "Gotenberg-Output-Filename"
	HeaderGotenbergTrace          = "Gotenberg-Trace"
)

var (
	PaperSizeLetter  = [2]float64{8.5, 11}
	PaperSizeLegal   = [2]float64{8.5, 14}
	PaperSizeTabloid = [2]float64{11, 17}
	PaperSizeLedger  = [2]float64{17, 11}
	PaperSizeA0      = [2]float64{33.1, 46.8}
	PaperSizeA1      = [2]float64{23.4, 33.1}
	PaperSizeA2      = [2]float64{16.54, 23.4}
	PaperSizeA3      = [2]float64{11.7, 16.54}
	PaperSizeA4      = [2]float64{8.27, 11.7}
	PaperSizeA5      = [2]float64{5.83, 8.27}
	PaperSizeA6      = [2]float64{4.13, 5.83}
)

const (
	FieldURL       = "url"
	FieldFiles     = "files"
	FileIndexHTML  = "index.html"
	FileFooterHTML = "footer.html"
	FileHeaderHTML = "header.html"
	FileStylesCSS  = "styles.css"
)

type Client struct {
	*httpclient.Client
}

type Request struct {
	*httpclient.Request
	err error
}

type Response struct {
	*http.Response
	trace string
}

func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
	httpClientWrapper, err := httpclient.NewClient(httpClient, baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: httpClientWrapper,
	}, nil
}

func (c *Client) ConvertHTML(ctx context.Context, html io.Reader) *Request {
	req := c.Client.MethodPost(ctx, ConvertHTML).File(FieldFiles, FileIndexHTML, html)
	return &Request{Request: req, err: req.Err()}
}

func (c *Client) ConvertURL(ctx context.Context, url string) *Request {
	req := c.Client.MethodPost(ctx, ConvertURL).FormField(FieldURL, url)
	return &Request{Request: req, err: req.Err()}
}

func (r *Request) WebhookURLMethodPost(url string) *Request {
	req := r.Request.Header(HeaderWebhookURL, url).Header(HeaderWebhookMethod, http.MethodPost)
	return &Request{Request: req, err: req.Err()}
}

func (r *Request) OutputFilename(filename string) *Request {
	req := r.Request.Header(HeaderOutputFilename, filename)
	return &Request{Request: req, err: req.Err()}
}

func (r *Request) WebhookErrorURLMethodPost(url string) *Request {
	req := r.Request.Header(HeaderWebhookErrorURL, url).Header(HeaderWebhookErrorMethod, http.MethodPost)
	return &Request{Request: req, err: req.Err()}
}

func (r *Request) WebhookExtraHeaders(headers map[string]string) *Request {
	jsonHeaders, err := json.Marshal(headers)
	if err != nil {
		r.err = err
		return r
	}

	req := r.Request.Header(HeaderWebhookExtraHTTPHeaders, string(jsonHeaders))
	return &Request{Request: req, err: req.Err()}
}

func (r *Request) Bool(fieldName string, value bool) *Request {
	req := r.Request.FormField(fieldName, fmt.Sprintf("%t", value))
	return &Request{Request: req, err: req.Err()}
}

func (r *Request) Float(fieldName string, value float64) *Request {
	req := r.Request.FormField(fieldName, fmt.Sprintf("%g", value))
	return &Request{Request: req, err: req.Err()}
}

func (r *Request) PaperSize(width, height float64) *Request {
	return r.Float(FieldPaperWidth, width).
		Float(FieldPaperHeight, height)
}

func (r *Request) PaperSizeA4() *Request {
	return r.PaperSize(PaperSizeA4[0], PaperSizeA4[1])
}

func (r *Request) PaperSizeLetter() *Request {
	return r.PaperSize(PaperSizeLetter[0], PaperSizeLetter[1])
}

func (r *Request) Margins(top, right, bottom, left float64) *Request {
	return r.Float(FieldMarginTop, top).
		Float(FieldMarginRight, right).
		Float(FieldMarginBottom, bottom).
		Float(FieldMarginLeft, left)
}

func (r *Request) File(fieldName, filename string, content io.Reader) *Request {
	req := r.Request.File(fieldName, filename, content)
	return &Request{Request: req, err: req.Err()}
}

func (r *Request) Header(key, value string) *Request {
	req := r.Request.Header(key, value)
	return &Request{Request: req, err: req.Err()}
}

func (r *Request) FormField(fieldName, value string) *Request {
	req := r.Request.FormField(fieldName, value)
	return &Request{Request: req, err: req.Err()}
}

func (r *Request) Send() (*Response, error) {
	if r.err != nil {
		return nil, r.err
	}
	resp, err := r.Request.Send()
	if err != nil {
		return nil, err
	}
	return &Response{Response: resp, trace: resp.Header.Get(HeaderGotenbergTrace)}, nil
}
