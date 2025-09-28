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
	"strings"
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
	ApplicationJSON = "application/json"
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
)

const (
	HeaderWebhookURL              = "Gotenberg-Webhook-Url"
	HeaderWebhookErrorURL         = "Gotenberg-Webhook-Error-Url"
	HeaderWebhookMethod           = "Gotenberg-Webhook-Method"
	HeaderWebhookErrorMethod      = "Gotenberg-Webhook-Error-Method"
	HeaderWebhookExtraHTTPHeaders = "Gotenberg-Webhook-Extra-Http-Headers"
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
	bufferSize = 1 << 12 // 4096 bytes (4 KB)
)

const (
	ConvertHTML = "html"
	ConvertURL  = "url"
)

var (
	routesURI = map[string]string{
		ConvertHTML: "/forms/chromium/convert/html",
		ConvertURL:  "/forms/chromium/convert/url",
	}

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

type Header struct {
	Name  string
	Value string
}

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	buffer     *bytes.Buffer
	writer     *multipart.Writer
	bufPool    sync.Pool
	err        error
}

type RequestOptions func(r *http.Request) error

func WithHost() RequestOptions {
	return func(r *http.Request) error {
		return WithHeader("Host", r.URL.Host)(r)
	}
}

func WithQueryParams(values url.Values) RequestOptions {
	return func(r *http.Request) error {
		q := r.URL.Query()
		for k := range values {
			v := values.Get(k)
			if v == "" {
				q.Del(k)
				continue
			}
			q.Add(k, values.Get(k))
		}
		r.URL.RawQuery = q.Encode()
		return nil
	}
}

func WithHeader(key, value string) RequestOptions {
	return func(r *http.Request) error {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set(key, value)

		return nil
	}
}

func WebhookErrorURL(url, method string) RequestOptions {
	WithHeader(HeaderWebhookErrorURL, url)
	return WithHeader(HeaderWebhookErrorMethod, method)
}

func WebhookExtraHTTPHeaders(headers map[string]string) RequestOptions {
	h, _ := json.Marshal(headers)
	return WithHeader(HeaderWebhookExtraHTTPHeaders, string(h))
}

func WithJSONBody(body interface{}) RequestOptions {
	return func(r *http.Request) error {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		err = WithBytesBody(b)(r)
		if err != nil {
			return err
		}
		return WithHeader(ContentType, ApplicationJSON)(r)
	}
}

func WithBytesBody(body []byte) RequestOptions {
	return func(r *http.Request) error {
		r.Body = io.NopCloser(bytes.NewReader(body))
		return WithHeader(ContentLength, strconv.Itoa(len(body)))(r)
	}
}

func WithReaderBody(body io.ReadCloser) RequestOptions {
	return func(r *http.Request) error {
		r.Body = body
		return nil
	}
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	u, err := url.Parse(strings.TrimSuffix(baseURL, "/"))
	if err != nil {
		return &Client{
			err: fmt.Errorf("invalid base URL: %w", err),
		}
	}

	buf := &bytes.Buffer{}
	buf.Grow(bufferSize)

	c := &Client{
		httpClient: httpClient,
		baseURL:    u,
		buffer:     buf,
		writer:     multipart.NewWriter(buf),
	}

	c.bufPool = sync.Pool{
		New: func() any {
			buf := make([]byte, 0, bufferSize)
			return &buf
		},
	}

	return c
}

func (c *Client) File(filename string, content io.Reader) *Client {
	if c.err != nil {
		return c
	}

	if c.writer == nil {
		c.err = fmt.Errorf("client not properly initialized")
		return c
	}

	part, err := c.writer.CreateFormFile(FieldFiles, filename)
	if err != nil {
		c.err = fmt.Errorf("failed to create form file: %w", err)
		return c
	}

	var buf []byte
	if p := c.bufPool.Get(); p != nil {
		buf = (*p.(*[]byte))[:bufferSize]
	} else {
		buf = make([]byte, bufferSize)
	}
	defer func() {
		buf = buf[:0]
		c.bufPool.Put(&buf)
	}()

	_, err = io.CopyBuffer(part, content, buf)
	if err != nil {
		c.err = fmt.Errorf("failed to copy buffer: %w", err)
		return c
	}
	return c
}

func (c *Client) String(field string, value string) *Client {
	if c.err != nil {
		return c
	}

	if c.writer == nil {
		c.err = fmt.Errorf("client not properly initialized")
		return c
	}

	err := c.writer.WriteField(field, value)
	if err != nil {
		c.err = fmt.Errorf("failed to write field %q: %w", field, err)
		return c
	}
	return c
}

func (c *Client) URL(url string) *Client {
	return c.String(FieldURL, url)
}

func (c *Client) IndexHTML(html io.Reader) *Client {
	return c.File(FileIndexHTML, html)
}

func (c *Client) FooterHTML(html io.Reader) *Client {
	return c.File(FileFooterHTML, html)
}

func (c *Client) HeaderHTML(html io.Reader) *Client {
	return c.File(FileHeaderHTML, html)
}

func (c *Client) StylesCSS(css io.Reader) *Client {
	return c.File(FileStylesCSS, css)
}

func (c *Client) Bool(field string, value bool) *Client {
	return c.String(field, fmt.Sprintf("%t", value))
}

func (c *Client) Float(field string, value float64) *Client {
	return c.String(field, fmt.Sprintf("%g", value))
}

func (c *Client) paperSize(wh [2]float64) *Client {
	return c.Float(FieldPaperWidth, wh[0]).
		Float(FieldPaperHeight, wh[1])
}

func (c *Client) PaperSize(width, height float64) *Client {
	return c.paperSize([2]float64{width, height})
}

func (c *Client) PaperSizeA4() *Client {
	return c.paperSize(PaperSizeA4)
}

func (c *Client) PaperSizeLetter() *Client {
	return c.paperSize(PaperSizeLetter)
}

func (c *Client) Margins(top, right, bottom, left float64) *Client {
	return c.Float(FieldMarginTop, top).
		Float(FieldMarginRight, right).
		Float(FieldMarginBottom, bottom).
		Float(FieldMarginLeft, left)
}

func (c *Client) Err() error {
	return c.err
}

func (c *Client) ResetClient() *Client {
	c.buffer.Reset()
	c.writer = multipart.NewWriter(c.buffer)
	c.err = nil
	return c
}

func (c *Client) Execute(ctx context.Context, route string, opts ...RequestOptions) (*http.Response, error) {
	defer c.ResetClient()

	if c.err != nil {
		return nil, c.err
	}

	if c.buffer == nil || c.writer == nil {
		return nil, fmt.Errorf("client not properly initialized")
	}

	if err := c.writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	uri, ok := routesURI[route]

	if !ok {
		return nil, fmt.Errorf("unknown route: %s", route)
	}

	opts = append(opts,
		WithReaderBody(io.NopCloser(c.buffer)),
		WithHeader(ContentType, c.writer.FormDataContentType()),
		WithHeader(ContentLength, strconv.Itoa(c.buffer.Len())))

	return c.sendRequest(ctx, "POST", c.baseURL.JoinPath(uri).String(), opts...)
}

func (c *Client) sendRequest(ctx context.Context, method, url string, opts ...RequestOptions) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if len(opts) > 0 {
		for _, v := range opts {
			if err = v(req); err != nil {
				return nil, err
			}
		}
	}

	return c.httpClient.Do(req)
}

func (c *Client) ConvertHTML(ctx context.Context) (*http.Response, error) {
	return c.Execute(ctx, ConvertHTML)
}

func (c *Client) AsyncConvertHTML(ctx context.Context, webhookURL, webhookMethod, webhookErrorURL, webhookErrorMethod string, webHookHeaders map[string]string) (*http.Response, error) {
	return c.Execute(ctx, ConvertHTML,
		WithHeader(HeaderWebhookURL, webhookURL),
		WithHeader(HeaderWebhookMethod, webhookMethod),
		WithHeader(HeaderWebhookErrorURL, webhookErrorURL),
		WithHeader(HeaderWebhookErrorMethod, webhookErrorMethod),
		WebhookExtraHTTPHeaders(webHookHeaders),
	)
}

func (c *Client) ConvertURL(ctx context.Context) (*http.Response, error) {
	return c.Execute(ctx, ConvertURL)
}
