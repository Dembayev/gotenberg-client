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
		buf = (*p.(*[]byte))[:defaultBufferSize]
	} else {
		buf = make([]byte, defaultBufferSize)
	}
	defer func() {
		buf = buf[:0]
		c.bufPool.Put(&buf)
	}()

	_, err = io.CopyBuffer(part, content, buf)
	return err
}

func (c *Client) addPageProperties(writer *multipart.Writer, props pageProperties) error {
	boolFields := []struct {
		value *bool
		field string
	}{
		{props.SinglePage, FieldSinglePage},
		{props.PreferCSSPageSize, FieldPreferCSSPageSize},
		{props.GenerateDocumentOutline, FieldGenerateDocumentOutline},
		{props.GenerateTaggedPDF, FieldGenerateTaggedPDF},
		{props.PrintBackground, FieldPrintBackground},
		{props.OmitBackground, FieldOmitBackground},
		{props.Landscape, FieldLandscape},
	}

	for _, field := range boolFields {
		if field.value != nil {
			if err := writer.WriteField(field.field, strconv.FormatBool(*field.value)); err != nil {
				return err
			}
		}
	}

	floatFields := []struct {
		value *float64
		field string
	}{
		{props.PaperWidth, FieldPaperWidth},
		{props.PaperHeight, FieldPaperHeight},
		{props.MarginTop, FieldMarginTop},
		{props.MarginBottom, FieldMarginBottom},
		{props.MarginLeft, FieldMarginLeft},
		{props.MarginRight, FieldMarginRight},
		{props.Scale, FieldScale},
	}

	for _, field := range floatFields {
		if field.value != nil {
			if err := writer.WriteField(field.field, strconv.FormatFloat(*field.value, 'f', -1, 64)); err != nil {
				return err
			}
		}
	}

	if props.NativePageRanges != nil {
		if err := writer.WriteField(FieldNativePageRanges, *props.NativePageRanges); err != nil {
			return err
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
		// Pre-uppercase common methods to avoid allocation
		method := *opts.Method
		switch method {
		case "post", "POST":
			req.Header.Set(HeaderWebhookMethod, "POST")
		case "get", "GET":
			req.Header.Set(HeaderWebhookMethod, "GET")
		case "put", "PUT":
			req.Header.Set(HeaderWebhookMethod, "PUT")
		case "patch", "PATCH":
			req.Header.Set(HeaderWebhookMethod, "PATCH")
		case "delete", "DELETE":
			req.Header.Set(HeaderWebhookMethod, "DELETE")
		default:
			req.Header.Set(HeaderWebhookMethod, strings.ToUpper(method))
		}
	}

	if opts.ErrorMethod != nil {
		// Same optimization for error method
		method := *opts.ErrorMethod
		switch method {
		case "post", "POST":
			req.Header.Set(HeaderWebhookErrorMethod, "POST")
		case "get", "GET":
			req.Header.Set(HeaderWebhookErrorMethod, "GET")
		case "put", "PUT":
			req.Header.Set(HeaderWebhookErrorMethod, "PUT")
		case "patch", "PATCH":
			req.Header.Set(HeaderWebhookErrorMethod, "PATCH")
		case "delete", "DELETE":
			req.Header.Set(HeaderWebhookErrorMethod, "DELETE")
		default:
			req.Header.Set(HeaderWebhookErrorMethod, strings.ToUpper(method))
		}
	}

	if len(opts.ExtraHeaders) > 0 {
		extraHeaders, err := json.Marshal(opts.ExtraHeaders)
		if err == nil {
			req.Header.Set(HeaderWebhookExtraHTTPHeaders, string(extraHeaders))
		}
	}
}
