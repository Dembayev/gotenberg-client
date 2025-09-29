package gotenberg

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	httpclient "github.com/nativebpm/http-client"
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
	*httpclient.Multipart
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
	r.Multipart, r.err = c.Client.MultipartPOST(ctx, ConvertHTML).File(FieldFiles, FileIndexHTML, html).GetRequest()
	return r
}

func (c *Client) ConvertURL(ctx context.Context, url string) *Request {
	r := &Request{}
	r.Multipart, r.err = c.Client.MultipartPOST(ctx, ConvertURL).FormField(FieldURL, url).GetRequest()
	return r
}

func (r *Request) WebhookURL(url, method string) *Request {
	if r.err != nil {
		return r
	}
	r.Multipart, r.err = r.Multipart.Header(HeaderWebhookURL, url).Header(HeaderWebhookMethod, method).GetRequest()
	return r
}

func (r *Request) OutputFilename(filename string) *Request {
	if r.err != nil {
		return r
	}
	r.Multipart, r.err = r.Multipart.Header(HeaderOutputFilename, filename).GetRequest()
	return r
}

func (r *Request) WebhookErrorURL(url, method string) *Request {
	if r.err != nil {
		return r
	}
	r.Multipart, r.err = r.Multipart.Header(HeaderWebhookErrorURL, url).Header(HeaderWebhookErrorMethod, method).GetRequest()
	return r
}

func (r *Request) WebhookHeaders(headers map[string]string) *Request {
	jsonHeaders, err := json.Marshal(headers)
	if err != nil {
		r.err = err
		return r
	}
	r.Multipart, r.err = r.Multipart.Header(HeaderWebhookExtraHTTPHeaders, string(jsonHeaders)).GetRequest()
	return r
}

func (r *Request) Bool(fieldName string, value bool) *Request {
	if r.err != nil {
		return r
	}
	r.Multipart, r.err = r.Multipart.FormField(fieldName, strconv.FormatBool(value)).GetRequest()
	return r
}

func (r *Request) Float(fieldName string, value float64) *Request {
	if r.err != nil {
		return r
	}
	r.Multipart, r.err = r.Multipart.FormField(fieldName, strconv.FormatFloat(value, 'f', -1, 64)).GetRequest()
	return r
}

func (r *Request) PaperSize(width, height float64) *Request {
	if r.err != nil {
		return r
	}
	r.Float(FieldPaperWidth, width)
	r.Float(FieldPaperHeight, height)
	return r
}

func (r *Request) PaperSizeA4() *Request {
	return r.PaperSize(PaperSizeA4[0], PaperSizeA4[1])
}

func (r *Request) PaperSizeLetter() *Request {
	return r.PaperSize(PaperSizeLetter[0], PaperSizeLetter[1])
}

func (r *Request) Margins(top, right, bottom, left float64) *Request {
	if r.err != nil {
		return r
	}
	r.Float(FieldMarginTop, top)
	r.Float(FieldMarginRight, right)
	r.Float(FieldMarginBottom, bottom)
	r.Float(FieldMarginLeft, left)
	return r
}

func (r *Request) File(fieldName, filename string, content io.Reader) *Request {
	if r.err != nil {
		return r
	}
	r.Multipart, r.err = r.Multipart.File(fieldName, filename, content).GetRequest()
	return r
}

func (r *Request) Header(key, value string) *Request {
	if r.err != nil {
		return r
	}
	r.Multipart, r.err = r.Multipart.Header(key, value).GetRequest()
	return r
}

func (r *Request) FormField(fieldName, value string) *Request {
	if r.err != nil {
		return r
	}
	r.Multipart, r.err = r.Multipart.FormField(fieldName, value).GetRequest()
	return r
}

func (r *Request) Send() (*Response, error) {
	if r.err != nil {
		return nil, r.err
	}
	resp, err := r.Multipart.Send()
	if err != nil {
		return nil, err
	}
	return &Response{Response: resp, GotenbergTrace: resp.Header.Get(HeaderGotenbergTrace)}, nil
}
