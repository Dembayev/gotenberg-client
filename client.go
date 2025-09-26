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
)

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	u, err := url.Parse(strings.TrimSuffix(baseURL, "/"))
	if err != nil {
		panic(fmt.Sprintf("invalid base URL: %s", err))
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    u,
	}
}

func (c Client) ConvertURLToPDF(ctx context.Context, url string, opts ...ConvOption) (*http.Response, error) {
	if url == "" {
		return nil, fmt.Errorf("URL is required")
	}

	config := &convConfig{}
	for _, opt := range opts {
		opt(config)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if err := writer.WriteField("url", url); err != nil {
		return nil, fmt.Errorf("failed to write url field: %w", err)
	}

	if err := c.addPageProperties(writer, pageProperties{
		SinglePage:              config.SinglePage,
		PaperWidth:              config.PaperWidth,
		PaperHeight:             config.PaperHeight,
		MarginTop:               config.MarginTop,
		MarginBottom:            config.MarginBottom,
		MarginLeft:              config.MarginLeft,
		MarginRight:             config.MarginRight,
		PreferCSSPageSize:       config.PreferCSSPageSize,
		GenerateDocumentOutline: config.GenerateDocumentOutline,
		GenerateTaggedPDF:       config.GenerateTaggedPDF,
		PrintBackground:         config.PrintBackground,
		OmitBackground:          config.OmitBackground,
		Landscape:               config.Landscape,
		Scale:                   config.Scale,
		NativePageRanges:        config.NativePageRanges,
	}); err != nil {
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

	c.addWebhookHeaders(req, webhookOptions{
		URL:          config.WebhookURL,
		ErrorURL:     config.WebhookErrorURL,
		Method:       config.WebhookMethod,
		ErrorMethod:  config.WebhookErrorMethod,
		ExtraHeaders: config.WebhookExtraHeaders,
	})

	if config.OutputFilename != nil {
		req.Header.Set("Gotenberg-Output-Filename", *config.OutputFilename)
	}

	return c.httpClient.Do(req)
}

func (c Client) ConvertHTMLToPDF(ctx context.Context, indexHTML io.Reader, opts ...ConvOption) (*http.Response, error) {
	if indexHTML == nil {
		return nil, fmt.Errorf("indexHTML is required")
	}

	config := &convConfig{}
	for _, opt := range opts {
		opt(config)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if config.Files == nil {
		config.Files = make(map[string]io.Reader)
	}
	config.Files["index.html"] = indexHTML
	for filename, content := range config.Files {
		if err := c.addFileField(writer, "files", filename, content); err != nil {
			return nil, fmt.Errorf("failed to add file %s: %w", filename, err)
		}
	}

	if err := c.addPageProperties(writer, pageProperties{
		SinglePage:              config.SinglePage,
		PaperWidth:              config.PaperWidth,
		PaperHeight:             config.PaperHeight,
		MarginTop:               config.MarginTop,
		MarginBottom:            config.MarginBottom,
		MarginLeft:              config.MarginLeft,
		MarginRight:             config.MarginRight,
		PreferCSSPageSize:       config.PreferCSSPageSize,
		GenerateDocumentOutline: config.GenerateDocumentOutline,
		GenerateTaggedPDF:       config.GenerateTaggedPDF,
		PrintBackground:         config.PrintBackground,
		OmitBackground:          config.OmitBackground,
		Landscape:               config.Landscape,
		Scale:                   config.Scale,
		NativePageRanges:        config.NativePageRanges,
	}); err != nil {
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

	c.addWebhookHeaders(req, webhookOptions{
		URL:          config.WebhookURL,
		ErrorURL:     config.WebhookErrorURL,
		Method:       config.WebhookMethod,
		ErrorMethod:  config.WebhookErrorMethod,
		ExtraHeaders: config.WebhookExtraHeaders,
	})

	if config.OutputFilename != nil {
		req.Header.Set("Gotenberg-Output-Filename", *config.OutputFilename)
	}

	return c.httpClient.Do(req)
}
