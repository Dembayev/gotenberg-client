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
)

const (
	ConvertHTML = "/forms/chromium/convert/html"
	ConvertURL  = "/forms/chromium/convert/url"
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
	bufferSize = 1 << 12 // 4096 bytes (4 KB)
)

const (
	FieldURL       = "url"
	FieldFiles     = "files"
	FileIndexHTML  = "index.html"
	FileFooterHTML = "footer.html"
	FileHeaderHTML = "header.html"
	FileStylesCSS  = "styles.css"
)

const (
	ApplicationJSON = "application/json"
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
)

type Form struct {
	multipart bool
	writer    *multipart.Writer
	buffer    *bytes.Buffer
	bufPool   sync.Pool
}

type Request struct {
	client  *http.Client
	request *http.Request
	err     error
	Form
}

func NewReqest(ctx context.Context, client *http.Client, url string) *Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return &Request{
			err: fmt.Errorf("failed to create request: %w", err),
		}
	}

	return &Request{
		client:  client,
		request: req,
	}
}

func (r *Request) MethodGet(ctx context.Context, client *http.Client, url string) *Request {
	if r.err != nil {
		return r
	}

	r.request.Method = http.MethodGet
	return r
}

func (r *Request) MethodPost(ctx context.Context, client *http.Client, url string) *Request {
	if r.err != nil {
		return r
	}

	r.request.Method = http.MethodPost
	return r
}

func (r *Request) MethodPut(ctx context.Context, client *http.Client, url string) *Request {
	if r.err != nil {
		return r
	}

	r.request.Method = http.MethodPut
	return r
}

func (r *Request) MethodPatch(ctx context.Context, client *http.Client, url string) *Request {
	if r.err != nil {
		return r
	}

	r.request.Method = http.MethodPatch
	return r
}

func (r *Request) MethodDelete(ctx context.Context, client *http.Client, url string) *Request {
	if r.err != nil {
		return r
	}

	r.request.Method = http.MethodDelete
	return r
}

func (r *Request) Multipart() *Request {
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

func (r *Request) Header(key, value string) *Request {
	if r.err != nil {
		return r
	}

	if r.request.Header == nil {
		r.request.Header = make(http.Header)
	}
	r.request.Header.Set(key, value)

	return r
}

func (r *Request) Headers(headers map[string]string) *Request {
	if r.err != nil {
		return r
	}

	for key, value := range headers {
		r = r.Header(key, value)
	}

	return r
}

func (r *Request) ContentType(contentType string) *Request {
	return r.Header(ContentType, contentType)
}

func (r *Request) JSONContentType() *Request {
	return r.ContentType(ApplicationJSON)
}

func (r *Request) QueryParam(key, value string) *Request {
	if r.err != nil {
		return r
	}

	q := r.request.URL.Query()
	q.Set(key, value)
	r.request.URL.RawQuery = q.Encode()

	return r
}

func (r *Request) QueryParams(params map[string]string) *Request {
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

func (r *Request) QueryValues(values url.Values) *Request {
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

func (r *Request) Body(body io.ReadCloser) *Request {
	if r.err != nil {
		return r
	}

	r.request.Body = body
	return r
}

func (r *Request) BytesBody(body []byte) *Request {
	if r.err != nil {
		return r
	}

	r.request.Body = io.NopCloser(bytes.NewReader(body))
	r = r.Header(ContentLength, strconv.Itoa(len(body)))

	return r
}

func (r *Request) StringBody(body string) *Request {
	return r.BytesBody([]byte(body))
}

func (r *Request) JSONBody(body any) *Request {
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

func (r *Request) WebhookURLMethodPost(url string) *Request {
	return r.Header(HeaderWebhookURL, url).Header(HeaderWebhookMethod, http.MethodPost)
}

func (r *Request) WebhookErrorURLMethodPost(url string) *Request {
	return r.Header(HeaderWebhookErrorURL, url).Header(HeaderWebhookErrorMethod, http.MethodPost)
}

func (r *Request) WebhookExtraHeaders(headers map[string]string) *Request {
	if r.err != nil {
		return r
	}

	jsonHeaders, err := json.Marshal(headers)
	if err != nil {
		r.err = fmt.Errorf("failed to marshal webhook headers: %w", err)
		return r
	}

	return r.Header(HeaderWebhookExtraHTTPHeaders, string(jsonHeaders))
}

func (r *Request) File(fieldName, filename string, content io.Reader) *Request {
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

func (r *Request) FormField(fieldName, value string) *Request {
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

func (r *Request) IndexHTML(html io.Reader) *Request {
	return r.File(FieldFiles, FileIndexHTML, html)
}

func (r *Request) FooterHTML(html io.Reader) *Request {
	return r.File(FieldFiles, FileFooterHTML, html)
}

func (r *Request) HeaderHTML(html io.Reader) *Request {
	return r.File(FieldFiles, FileHeaderHTML, html)
}

func (r *Request) StylesCSS(css io.Reader) *Request {
	return r.File(FieldFiles, FileStylesCSS, css)
}

func (r *Request) URL(url string) *Request {
	return r.FormField(FieldURL, url)
}

func (r *Request) Bool(fieldName string, value bool) *Request {
	return r.FormField(fieldName, fmt.Sprintf("%t", value))
}

func (r *Request) Float(fieldName string, value float64) *Request {
	return r.FormField(fieldName, fmt.Sprintf("%g", value))
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

func (r *Request) Err() error {
	return r.err
}

func (r *Request) Send() (*http.Response, error) {
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
