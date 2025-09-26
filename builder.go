package gotenberg

import "io"

// OptionsBuilder provides a fluent interface for building client options.
// This reduces allocations when chaining multiple options and provides
// better readability compared to functional options.
//
// Example usage:
//
//	options := gotenberg.NewOptionsBuilder().
//		PaperSizeA4().
//		Margins(1.0, 1.0, 1.0, 1.0).
//		PrintBackground(true).
//		OutputFilename("document.pdf").
//		Build()
//
//	resp, err := client.ConvertHTMLToPDF(ctx, htmlReader, options)
type OptionsBuilder struct {
	config *clientOptions
}

func NewOptionsBuilder() *OptionsBuilder {
	return &OptionsBuilder{
		config: &clientOptions{
			Files: make(map[string]io.Reader),
		},
	}
}

func (b *OptionsBuilder) File(name string, r io.Reader) *OptionsBuilder {
	b.config.Files[name] = r
	return b
}

func (b *OptionsBuilder) PaperSize(width, height float64) *OptionsBuilder {
	b.config.Page.PaperWidth = &width
	b.config.Page.PaperHeight = &height
	return b
}

func (b *OptionsBuilder) Margins(top, right, bottom, left float64) *OptionsBuilder {
	b.config.Page.MarginTop = &top
	b.config.Page.MarginRight = &right
	b.config.Page.MarginBottom = &bottom
	b.config.Page.MarginLeft = &left
	return b
}

func (b *OptionsBuilder) SinglePage(enabled bool) *OptionsBuilder {
	b.config.Page.SinglePage = &enabled
	return b
}

func (b *OptionsBuilder) Landscape(enabled bool) *OptionsBuilder {
	b.config.Page.Landscape = &enabled
	return b
}

func (b *OptionsBuilder) PrintBackground(enabled bool) *OptionsBuilder {
	b.config.Page.PrintBackground = &enabled
	return b
}

func (b *OptionsBuilder) Scale(scale float64) *OptionsBuilder {
	b.config.Page.Scale = &scale
	return b
}

func (b *OptionsBuilder) OutputFilename(filename string) *OptionsBuilder {
	b.config.OutputFilename = &filename
	return b
}

func (b *OptionsBuilder) WebhookSuccess(url, method string) *OptionsBuilder {
	b.config.Webhook.URL = &url
	b.config.Webhook.Method = &method
	return b
}

func (b *OptionsBuilder) WebhookError(errorURL, errorMethod string) *OptionsBuilder {
	b.config.Webhook.ErrorURL = &errorURL
	b.config.Webhook.ErrorMethod = &errorMethod
	return b
}

func (b *OptionsBuilder) WebhookExtraHeader(name, value string) *OptionsBuilder {
	if b.config.Webhook.ExtraHeaders == nil {
		b.config.Webhook.ExtraHeaders = make(map[string]string)
	}
	b.config.Webhook.ExtraHeaders[name] = value
	return b
}

func (b *OptionsBuilder) PaperSizeA4() *OptionsBuilder {
	return b.PaperSize(PaperSizeA4[0], PaperSizeA4[1])
}

func (b *OptionsBuilder) PaperSizeLetter() *OptionsBuilder {
	return b.PaperSize(PaperSizeLetter[0], PaperSizeLetter[1])
}

func (b *OptionsBuilder) Build() ClientOptions {
	config := *b.config
	return func(c *clientOptions) {
		*c = config
	}
}
