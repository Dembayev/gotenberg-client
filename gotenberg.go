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
	"strings"
	"sync"
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
	defaultBufferSize = 1 << 12 // 4KB
)

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	buffer     *bytes.Buffer
	writer     *multipart.Writer
	bufPool    sync.Pool
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	u, err := url.Parse(strings.TrimSuffix(baseURL, "/"))
	if err != nil {
		panic(fmt.Sprintf("invalid base URL: %s", err))
	}

	buf := &bytes.Buffer{}
	buf.Grow(8192)

	c := &Client{
		httpClient: httpClient,
		baseURL:    u,
		buffer:     buf,
		writer:     multipart.NewWriter(buf),
	}

	c.bufPool = sync.Pool{
		New: func() any {
			buf := make([]byte, 0, defaultBufferSize)
			return &buf
		},
	}

	return c
}

func (c *Client) WriteHTML(html io.Reader) error {
	return c.addFileField(c.writer, "files", "index.html", html)
}

func (c *Client) WriteFile(filename string, content io.Reader) error {
	return c.addFileField(c.writer, "files", filename, content)
}

func (c *Client) WriteURL(url string) error {
	return c.writer.WriteField("url", url)
}

func (c *Client) WritePaperSize(width, height float64) error {
	if err := c.writer.WriteField(FieldPaperWidth, fmt.Sprintf("%g", width)); err != nil {
		return err
	}
	return c.writer.WriteField(FieldPaperHeight, fmt.Sprintf("%g", height))
}

func (c *Client) WriteMargins(top, right, bottom, left float64) error {
	fields := []struct {
		name  string
		value float64
	}{
		{FieldMarginTop, top},
		{FieldMarginRight, right},
		{FieldMarginBottom, bottom},
		{FieldMarginLeft, left},
	}

	for _, field := range fields {
		if err := c.writer.WriteField(field.name, fmt.Sprintf("%g", field.value)); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) WriteBoolProperty(field string, value bool) error {
	return c.writer.WriteField(field, fmt.Sprintf("%t", value))
}

func (c *Client) WriteStringProperty(field, value string) error {
	return c.writer.WriteField(field, value)
}

// Webhook поддержка
func (c *Client) SetWebhookSuccess(url, method string) error {
	if err := c.writer.WriteField("webhookURL", url); err != nil {
		return err
	}
	return c.writer.WriteField("webhookMethod", method)
}

func (c *Client) SetWebhookError(url, method string) error {
	if err := c.writer.WriteField("webhookErrorURL", url); err != nil {
		return err
	}
	return c.writer.WriteField("webhookErrorMethod", method)
}

func (c *Client) SetWebhookHeaders(headers map[string]string) error {
	if len(headers) == 0 {
		return nil
	}

	headersJSON, err := json.Marshal(headers)
	if err != nil {
		return err
	}

	return c.writer.WriteField("webhookExtraHTTPHeaders", string(headersJSON))
}

func (c *Client) addFileField(writer *multipart.Writer, fieldName, filename string, content io.Reader) error {
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return err
	}

	var buf []byte
	if p := c.bufPool.Get(); p != nil {
		buf = (*p.(*[]byte))[:defaultBufferSize]
	} else {
		buf = make([]byte, defaultBufferSize)
	}
	defer func() {
		buf = buf[:0]
		c.bufPool.Put(&buf)
	}()

	_, err = io.CopyBuffer(part, content, buf)
	return err
}

func (c *Client) Execute(ctx context.Context) (*http.Response, error) {
	return c.ExecuteHTML(ctx)
}

func (c *Client) ExecuteHTML(ctx context.Context) (*http.Response, error) {
	return c.executeRequest(ctx, "/forms/chromium/convert/html")
}

func (c *Client) ExecuteURL(ctx context.Context) (*http.Response, error) {
	return c.executeRequest(ctx, "/forms/chromium/convert/url")
}

func (c *Client) executeRequest(ctx context.Context, endpoint string) (*http.Response, error) {
	if err := c.writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.baseURL.JoinPath(endpoint).String(),
		c.buffer,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", c.writer.FormDataContentType())
	req.ContentLength = int64(c.buffer.Len())

	return c.httpClient.Do(req)
}

func (c *Client) Reset() {
	c.buffer.Reset()
	c.writer = multipart.NewWriter(c.buffer)
}

func (c *Client) ContentType() string {
	return c.writer.FormDataContentType()
}

func (c *Client) BufferSize() int {
	return c.buffer.Len()
}
