package gotenberg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBufferSize = 1 << 12 // 4KB - better for typical file operations
)

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	bufPool    sync.Pool
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	u, err := url.Parse(strings.TrimSuffix(baseURL, "/"))
	if err != nil {
		panic(fmt.Sprintf("invalid base URL: %s", err))
	}

	c := &Client{
		httpClient: httpClient,
		baseURL:    u,
	}

	c.bufPool = sync.Pool{
		New: func() any {
			buf := make([]byte, 0, defaultBufferSize)
			return &buf
		},
	}

	return c
}

func (c *Client) ConvertURLToPDF(ctx context.Context, url string, opts ...ClientOptions) (*http.Response, error) {
	if len(url) == 0 {
		return nil, fmt.Errorf("URL is required")
	}

	config := &clientOptions{}
	for _, opt := range opts {
		opt(config)
	}

	var buf bytes.Buffer
	buf.Grow(2048)
	writer := multipart.NewWriter(&buf)

	if err := writer.WriteField("url", url); err != nil {
		return nil, fmt.Errorf("failed to write url field: %w", err)
	}

	if err := c.addPageProperties(writer, config.Page); err != nil {
		return nil, fmt.Errorf("failed to add page properties: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL.JoinPath("/forms/chromium/convert/url").String(), &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ContentLength = int64(buf.Len())

	c.addWebhookHeaders(req, config.Webhook)

	if config.OutputFilename != nil {
		req.Header.Set("Gotenberg-Output-Filename", *config.OutputFilename)
	}

	return c.httpClient.Do(req)
}

func (c *Client) ConvertHTMLToPDF(ctx context.Context, indexHTML io.Reader, opts ...ClientOptions) (*http.Response, error) {
	if indexHTML == nil {
		return nil, fmt.Errorf("indexHTML is required")
	}

	config := &clientOptions{}
	for _, opt := range opts {
		opt(config)
	}

	var buf bytes.Buffer
	buf.Grow(4096)
	writer := multipart.NewWriter(&buf)

	if err := c.addFileField(writer, "files", "index.html", indexHTML); err != nil {
		return nil, fmt.Errorf("failed to add file %s: %w", "index.html", err)
	}

	for filename, content := range config.Files {
		if err := c.addFileField(writer, "files", filename, content); err != nil {
			return nil, fmt.Errorf("failed to add file %s: %w", filename, err)
		}
	}

	if err := c.addPageProperties(writer, config.Page); err != nil {
		return nil, fmt.Errorf("failed to add page properties: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL.JoinPath("/forms/chromium/convert/html").String(), &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ContentLength = int64(buf.Len())

	c.addWebhookHeaders(req, config.Webhook)

	if config.OutputFilename != nil {
		req.Header.Set("Gotenberg-Output-Filename", *config.OutputFilename)
	}

	return c.httpClient.Do(req)
}
