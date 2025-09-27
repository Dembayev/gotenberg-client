package gotenberg

import (
	"context"
	"fmt"
	"net/http"
)

type URLConversionBuilder struct {
	client *Client
	url    string
	config *clientOptions
}

func (ucb *URLConversionBuilder) PaperSize(width, height float64) *URLConversionBuilder {
	ucb.config.Page.PaperWidth = &width
	ucb.config.Page.PaperHeight = &height
	return ucb
}

func (ucb *URLConversionBuilder) PaperSizeA4() *URLConversionBuilder {
	return ucb.PaperSize(PaperSizeA4[0], PaperSizeA4[1])
}

func (ucb *URLConversionBuilder) PaperSizeA3() *URLConversionBuilder {
	return ucb.PaperSize(PaperSizeA3[0], PaperSizeA3[1])
}

func (ucb *URLConversionBuilder) PaperSizeA5() *URLConversionBuilder {
	return ucb.PaperSize(PaperSizeA5[0], PaperSizeA5[1])
}

func (ucb *URLConversionBuilder) PaperSizeLetter() *URLConversionBuilder {
	return ucb.PaperSize(PaperSizeLetter[0], PaperSizeLetter[1])
}

func (ucb *URLConversionBuilder) PaperSizeLegal() *URLConversionBuilder {
	return ucb.PaperSize(PaperSizeLegal[0], PaperSizeLegal[1])
}

func (ucb *URLConversionBuilder) PaperSizeTabloid() *URLConversionBuilder {
	return ucb.PaperSize(PaperSizeTabloid[0], PaperSizeTabloid[1])
}

func (ucb *URLConversionBuilder) Margins(top, right, bottom, left float64) *URLConversionBuilder {
	ucb.config.Page.MarginTop = &top
	ucb.config.Page.MarginRight = &right
	ucb.config.Page.MarginBottom = &bottom
	ucb.config.Page.MarginLeft = &left
	return ucb
}

func (ucb *URLConversionBuilder) SinglePage(enabled bool) *URLConversionBuilder {
	ucb.config.Page.SinglePage = &enabled
	return ucb
}

func (ucb *URLConversionBuilder) Landscape(enabled bool) *URLConversionBuilder {
	ucb.config.Page.Landscape = &enabled
	return ucb
}

func (ucb *URLConversionBuilder) PrintBackground(enabled bool) *URLConversionBuilder {
	ucb.config.Page.PrintBackground = &enabled
	return ucb
}

func (ucb *URLConversionBuilder) Scale(scale float64) *URLConversionBuilder {
	ucb.config.Page.Scale = &scale
	return ucb
}

func (ucb *URLConversionBuilder) OutputFilename(filename string) *URLConversionBuilder {
	ucb.config.OutputFilename = &filename
	return ucb
}

func (ucb *URLConversionBuilder) WebhookSuccess(url, method string) *URLConversionBuilder {
	ucb.config.Webhook.URL = &url
	ucb.config.Webhook.Method = &method
	return ucb
}

func (ucb *URLConversionBuilder) WebhookError(errorURL, errorMethod string) *URLConversionBuilder {
	ucb.config.Webhook.ErrorURL = &errorURL
	ucb.config.Webhook.ErrorMethod = &errorMethod
	return ucb
}

func (ucb *URLConversionBuilder) WebhookExtraHeader(name, value string) *URLConversionBuilder {
	if ucb.config.Webhook.ExtraHeaders == nil {
		ucb.config.Webhook.ExtraHeaders = make(map[string]string)
	}
	ucb.config.Webhook.ExtraHeaders[name] = value
	return ucb
}

func (ucb *URLConversionBuilder) PreferCSSPageSize(enabled bool) *URLConversionBuilder {
	ucb.config.Page.PreferCSSPageSize = &enabled
	return ucb
}

func (ucb *URLConversionBuilder) GenerateDocumentOutline(enabled bool) *URLConversionBuilder {
	ucb.config.Page.GenerateDocumentOutline = &enabled
	return ucb
}

func (ucb *URLConversionBuilder) GenerateTaggedPDF(enabled bool) *URLConversionBuilder {
	ucb.config.Page.GenerateTaggedPDF = &enabled
	return ucb
}

func (ucb *URLConversionBuilder) OmitBackground(enabled bool) *URLConversionBuilder {
	ucb.config.Page.OmitBackground = &enabled
	return ucb
}

func (ucb *URLConversionBuilder) PageRanges(ranges string) *URLConversionBuilder {
	ucb.config.Page.NativePageRanges = &ranges
	return ucb
}

func (ucb *URLConversionBuilder) Execute(ctx context.Context, url string) (*http.Response, error) {
	if ucb.url == "" {
		return nil, fmt.Errorf("URL is required")
	}
	config := *ucb.config
	options := func(c *clientOptions) {
		*c = config
	}
	return ucb.client.ConvertURLToPDF(ctx, ucb.url, options)
}
