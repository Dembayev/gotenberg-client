package gotenberg

import "io"

type convConfig struct {
	Files          map[string]io.Reader
	OutputFilename *string

	PaperWidth              *float64
	PaperHeight             *float64
	MarginTop               *float64
	MarginBottom            *float64
	MarginLeft              *float64
	MarginRight             *float64
	Scale                   *float64
	SinglePage              *bool
	PreferCSSPageSize       *bool
	GenerateDocumentOutline *bool
	GenerateTaggedPDF       *bool
	PrintBackground         *bool
	OmitBackground          *bool
	Landscape               *bool

	NativePageRanges    *string
	WebhookURL          *string
	WebhookErrorURL     *string
	WebhookMethod       *string
	WebhookErrorMethod  *string
	WebhookExtraHeaders map[string]string
}
