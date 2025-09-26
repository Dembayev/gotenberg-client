package gotenberg

import "io"

type clientOptions struct {
	Files          map[string]io.Reader
	OutputFilename *string
	Page           pageProperties
	Webhook        webhookOptions
}

type pageProperties struct {
	SinglePage              *bool
	PaperWidth              *float64
	PaperHeight             *float64
	MarginTop               *float64
	MarginBottom            *float64
	MarginLeft              *float64
	MarginRight             *float64
	PreferCSSPageSize       *bool
	GenerateDocumentOutline *bool
	GenerateTaggedPDF       *bool
	PrintBackground         *bool
	OmitBackground          *bool
	Landscape               *bool
	Scale                   *float64
	NativePageRanges        *string
}

type webhookOptions struct {
	URL          *string
	ErrorURL     *string
	Method       *string
	ErrorMethod  *string
	ExtraHeaders map[string]string
}
