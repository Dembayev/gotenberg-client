package gotenberg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HTMLConversionBuilder struct {
	client    *Client
	indexHTML io.Reader
	config    *clientOptions
}

func (hcb *HTMLConversionBuilder) WithHTML(html string) *HTMLConversionBuilder {
	hcb.indexHTML = strings.NewReader(html)
	return hcb
}

func (hcb *HTMLConversionBuilder) WithHTMLReader(reader io.Reader) *HTMLConversionBuilder {
	hcb.indexHTML = reader
	return hcb
}

func (hcb *HTMLConversionBuilder) WithCSS(css string) *HTMLConversionBuilder {
	return hcb.WithFile("styles.css", strings.NewReader(css))
}

func (hcb *HTMLConversionBuilder) WithFile(filename string, reader io.Reader) *HTMLConversionBuilder {
	hcb.config.Files[filename] = reader
	return hcb
}

func (hcb *HTMLConversionBuilder) PaperSize(width, height float64) *HTMLConversionBuilder {
	hcb.config.Page.PaperWidth = &width
	hcb.config.Page.PaperHeight = &height
	return hcb
}

func (hcb *HTMLConversionBuilder) PaperSizeA4() *HTMLConversionBuilder {
	return hcb.PaperSize(PaperSizeA4[0], PaperSizeA4[1])
}

func (hcb *HTMLConversionBuilder) PaperSizeA3() *HTMLConversionBuilder {
	return hcb.PaperSize(PaperSizeA3[0], PaperSizeA3[1])
}

func (hcb *HTMLConversionBuilder) PaperSizeA5() *HTMLConversionBuilder {
	return hcb.PaperSize(PaperSizeA5[0], PaperSizeA5[1])
}

func (hcb *HTMLConversionBuilder) PaperSizeLetter() *HTMLConversionBuilder {
	return hcb.PaperSize(PaperSizeLetter[0], PaperSizeLetter[1])
}

func (hcb *HTMLConversionBuilder) PaperSizeLegal() *HTMLConversionBuilder {
	return hcb.PaperSize(PaperSizeLegal[0], PaperSizeLegal[1])
}

func (hcb *HTMLConversionBuilder) PaperSizeTabloid() *HTMLConversionBuilder {
	return hcb.PaperSize(PaperSizeTabloid[0], PaperSizeTabloid[1])
}

func (hcb *HTMLConversionBuilder) Margins(top, right, bottom, left float64) *HTMLConversionBuilder {
	hcb.config.Page.MarginTop = &top
	hcb.config.Page.MarginRight = &right
	hcb.config.Page.MarginBottom = &bottom
	hcb.config.Page.MarginLeft = &left
	return hcb
}

func (hcb *HTMLConversionBuilder) SinglePage(enabled bool) *HTMLConversionBuilder {
	hcb.config.Page.SinglePage = &enabled
	return hcb
}

func (hcb *HTMLConversionBuilder) Landscape(enabled bool) *HTMLConversionBuilder {
	hcb.config.Page.Landscape = &enabled
	return hcb
}

func (hcb *HTMLConversionBuilder) PrintBackground(enabled bool) *HTMLConversionBuilder {
	hcb.config.Page.PrintBackground = &enabled
	return hcb
}

func (hcb *HTMLConversionBuilder) Scale(scale float64) *HTMLConversionBuilder {
	hcb.config.Page.Scale = &scale
	return hcb
}

func (hcb *HTMLConversionBuilder) OutputFilename(filename string) *HTMLConversionBuilder {
	hcb.config.OutputFilename = &filename
	return hcb
}

func (hcb *HTMLConversionBuilder) WebhookSuccess(url, method string) *HTMLConversionBuilder {
	hcb.config.Webhook.URL = &url
	hcb.config.Webhook.Method = &method
	return hcb
}

func (hcb *HTMLConversionBuilder) WebhookError(errorURL, errorMethod string) *HTMLConversionBuilder {
	hcb.config.Webhook.ErrorURL = &errorURL
	hcb.config.Webhook.ErrorMethod = &errorMethod
	return hcb
}

func (hcb *HTMLConversionBuilder) WebhookExtraHeader(name, value string) *HTMLConversionBuilder {
	if hcb.config.Webhook.ExtraHeaders == nil {
		hcb.config.Webhook.ExtraHeaders = make(map[string]string)
	}
	hcb.config.Webhook.ExtraHeaders[name] = value
	return hcb
}

func (hcb *HTMLConversionBuilder) PreferCSSPageSize(enabled bool) *HTMLConversionBuilder {
	hcb.config.Page.PreferCSSPageSize = &enabled
	return hcb
}

func (hcb *HTMLConversionBuilder) GenerateDocumentOutline(enabled bool) *HTMLConversionBuilder {
	hcb.config.Page.GenerateDocumentOutline = &enabled
	return hcb
}

func (hcb *HTMLConversionBuilder) GenerateTaggedPDF(enabled bool) *HTMLConversionBuilder {
	hcb.config.Page.GenerateTaggedPDF = &enabled
	return hcb
}

func (hcb *HTMLConversionBuilder) OmitBackground(enabled bool) *HTMLConversionBuilder {
	hcb.config.Page.OmitBackground = &enabled
	return hcb
}

func (hcb *HTMLConversionBuilder) PageRanges(ranges string) *HTMLConversionBuilder {
	hcb.config.Page.NativePageRanges = &ranges
	return hcb
}

func (hcb *HTMLConversionBuilder) Execute(ctx context.Context) (*http.Response, error) {
	if hcb.indexHTML == nil {
		return nil, fmt.Errorf("HTML content is required")
	}
	config := *hcb.config
	options := func(c *clientOptions) {
		*c = config
	}
	return hcb.client.ConvertHTMLToPDF(ctx, hcb.indexHTML, options)
}
