package gotenberg

import "io"

type ConvOption func(*convConfig)

func WithFile(name string, r io.Reader) ConvOption {
	return func(c *convConfig) {
		if c.Files == nil {
			c.Files = make(map[string]io.Reader)
		}
		c.Files[name] = r
	}
}

func WithPaperSize(width, height float64) ConvOption {
	return func(c *convConfig) {
		c.PaperWidth = &width
		c.PaperHeight = &height
	}
}

func WithMargins(top, right, bottom, left float64) ConvOption {
	return func(c *convConfig) {
		c.MarginTop = &top
		c.MarginRight = &right
		c.MarginBottom = &bottom
		c.MarginLeft = &left
	}
}

func WithSinglePage(enabled bool) ConvOption {
	return func(c *convConfig) {
		c.SinglePage = &enabled
	}
}

func WithLandscape(enabled bool) ConvOption {
	return func(c *convConfig) {
		c.Landscape = &enabled
	}
}

func WithPrintBackground(enabled bool) ConvOption {
	return func(c *convConfig) {
		c.PrintBackground = &enabled
	}
}

func WithScale(scale float64) ConvOption {
	return func(c *convConfig) {
		c.Scale = &scale
	}
}

func WithOutputFilename(filename string) ConvOption {
	return func(c *convConfig) {
		c.OutputFilename = &filename
	}
}

func WithWebhookSuccess(url, method string) ConvOption {
	return func(c *convConfig) {
		c.WebhookURL = &url
		c.WebhookMethod = &method
	}
}

func WithWebhookError(errorURL, errorMethod string) ConvOption {
	return func(c *convConfig) {
		c.WebhookErrorURL = &errorURL
		c.WebhookErrorMethod = &errorMethod
	}
}

func WithWebhookExtraHeader(name, value string) ConvOption {
	return func(c *convConfig) {
		if c.WebhookExtraHeaders == nil {
			c.WebhookExtraHeaders = make(map[string]string)
		}
		c.WebhookExtraHeaders[name] = value
	}
}

func WithPaperSizeA4() ConvOption {
	return WithPaperSize(PaperSizeA4[0], PaperSizeA4[1])
}

func WithPaperSizeLetter() ConvOption {
	return WithPaperSize(PaperSizeLetter[0], PaperSizeLetter[1])
}
