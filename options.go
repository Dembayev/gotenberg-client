package gotenberg

import "io"

type ClientOptions func(*clientOptions)

func WithFile(name string, r io.Reader) ClientOptions {
	return func(c *clientOptions) {
		if c.Files == nil {
			c.Files = make(map[string]io.Reader)
		}
		c.Files[name] = r
	}
}

func WithPaperSize(width, height float64) ClientOptions {
	return func(c *clientOptions) {
		c.Page.PaperWidth = &width
		c.Page.PaperHeight = &height
	}
}

func WithMargins(top, right, bottom, left float64) ClientOptions {
	return func(c *clientOptions) {
		c.Page.MarginTop = &top
		c.Page.MarginRight = &right
		c.Page.MarginBottom = &bottom
		c.Page.MarginLeft = &left
	}
}

func WithSinglePage(enabled bool) ClientOptions {
	return func(c *clientOptions) {
		c.Page.SinglePage = &enabled
	}
}

func WithLandscape(enabled bool) ClientOptions {
	return func(c *clientOptions) {
		c.Page.Landscape = &enabled
	}
}

func WithPrintBackground(enabled bool) ClientOptions {
	return func(c *clientOptions) {
		c.Page.PrintBackground = &enabled
	}
}

func WithScale(scale float64) ClientOptions {
	return func(c *clientOptions) {
		c.Page.Scale = &scale
	}
}

func WithOutputFilename(filename string) ClientOptions {
	return func(c *clientOptions) {
		c.OutputFilename = &filename
	}
}

func WithWebhookSuccess(url, method string) ClientOptions {
	return func(c *clientOptions) {
		c.Webhook.URL = &url
		c.Webhook.Method = &method
	}
}

func WithWebhookError(errorURL, errorMethod string) ClientOptions {
	return func(c *clientOptions) {
		c.Webhook.ErrorURL = &errorURL
		c.Webhook.ErrorMethod = &errorMethod
	}
}

func WithWebhookExtraHeader(name, value string) ClientOptions {
	return func(c *clientOptions) {
		if c.Webhook.ExtraHeaders == nil {
			c.Webhook.ExtraHeaders = make(map[string]string)
		}
		c.Webhook.ExtraHeaders[name] = value
	}
}

func WithPaperSizeA4() ClientOptions {
	return WithPaperSize(PaperSizeA4[0], PaperSizeA4[1])
}

func WithPaperSizeLetter() ClientOptions {
	return WithPaperSize(PaperSizeLetter[0], PaperSizeLetter[1])
}
