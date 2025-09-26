package gotenberg

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

func (c *Client) addFileField(writer *multipart.Writer, fieldName, filename string, content io.Reader) error {
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return err
	}

	var buf []byte
	if p := c.bufPool.Get(); p != nil {
		buf = *p.(*[]byte)
	} else {
		buf = make([]byte, defaultBufferSize)
	}
	defer func() { c.bufPool.Put(&buf) }()

	_, err = io.CopyBuffer(part, content, buf)
	return err
}

func (c *Client) addPageProperties(writer *multipart.Writer, props pageProperties) error {
	type fieldWriter struct {
		value  any
		field  string
		format func(any) (string, bool)
	}

	boolFormatter := func(v any) (string, bool) {
		if ptr, ok := v.(*bool); ok && ptr != nil {
			return strconv.FormatBool(*ptr), true
		}
		return "", false
	}

	floatFormatter := func(v any) (string, bool) {
		if ptr, ok := v.(*float64); ok && ptr != nil {
			return strconv.FormatFloat(*ptr, 'f', -1, 64), true
		}
		return "", false
	}

	stringFormatter := func(v any) (string, bool) {
		if ptr, ok := v.(*string); ok && ptr != nil {
			return *ptr, true
		}
		return "", false
	}

	fields := []fieldWriter{
		{props.SinglePage, FieldSinglePage, boolFormatter},
		{props.PaperWidth, FieldPaperWidth, floatFormatter},
		{props.PaperHeight, FieldPaperHeight, floatFormatter},
		{props.MarginTop, FieldMarginTop, floatFormatter},
		{props.MarginBottom, FieldMarginBottom, floatFormatter},
		{props.MarginLeft, FieldMarginLeft, floatFormatter},
		{props.MarginRight, FieldMarginRight, floatFormatter},
		{props.PreferCSSPageSize, FieldPreferCSSPageSize, boolFormatter},
		{props.GenerateDocumentOutline, FieldGenerateDocumentOutline, boolFormatter},
		{props.GenerateTaggedPDF, FieldGenerateTaggedPDF, boolFormatter},
		{props.PrintBackground, FieldPrintBackground, boolFormatter},
		{props.OmitBackground, FieldOmitBackground, boolFormatter},
		{props.Landscape, FieldLandscape, boolFormatter},
		{props.Scale, FieldScale, floatFormatter},
		{props.NativePageRanges, FieldNativePageRanges, stringFormatter},
	}

	for _, fw := range fields {
		if val, ok := fw.format(fw.value); ok {
			if err := writer.WriteField(fw.field, val); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) addWebhookHeaders(req *http.Request, opts webhookOptions) {
	if opts.URL != nil {
		req.Header.Set(HeaderWebhookURL, *opts.URL)
	}

	if opts.ErrorURL != nil {
		req.Header.Set(HeaderWebhookErrorURL, *opts.ErrorURL)
	}

	if opts.Method != nil {
		req.Header.Set(HeaderWebhookMethod, strings.ToUpper(*opts.Method))
	}

	if opts.ErrorMethod != nil {
		req.Header.Set(HeaderWebhookErrorMethod, strings.ToUpper(*opts.ErrorMethod))
	}

	if len(opts.ExtraHeaders) > 0 {
		extraHeaders, err := json.Marshal(opts.ExtraHeaders)
		if err == nil {
			req.Header.Set(HeaderWebhookExtraHTTPHeaders, string(extraHeaders))
		}
	}
}
