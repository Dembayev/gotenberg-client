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
		buf = p.([]byte)
	} else {
		buf = make([]byte, defaultBufferSize)
	}
	defer c.bufPool.Put(buf)

	_, err = io.CopyBuffer(part, content, buf)
	return err
}

func (c *Client) addPageProperties(writer *multipart.Writer, props pageProperties) error {
	if props.SinglePage != nil {
		if err := writer.WriteField(FieldSinglePage, strconv.FormatBool(*props.SinglePage)); err != nil {
			return err
		}
	}

	if props.PaperWidth != nil {
		if err := writer.WriteField(FieldPaperWidth, strconv.FormatFloat(*props.PaperWidth, 'f', -1, 64)); err != nil {
			return err
		}
	}

	if props.PaperHeight != nil {
		if err := writer.WriteField(FieldPaperHeight, strconv.FormatFloat(*props.PaperHeight, 'f', -1, 64)); err != nil {
			return err
		}
	}

	if props.MarginTop != nil {
		if err := writer.WriteField(FieldMarginTop, strconv.FormatFloat(*props.MarginTop, 'f', -1, 64)); err != nil {
			return err
		}
	}

	if props.MarginBottom != nil {
		if err := writer.WriteField(FieldMarginBottom, strconv.FormatFloat(*props.MarginBottom, 'f', -1, 64)); err != nil {
			return err
		}
	}

	if props.MarginLeft != nil {
		if err := writer.WriteField(FieldMarginLeft, strconv.FormatFloat(*props.MarginLeft, 'f', -1, 64)); err != nil {
			return err
		}
	}

	if props.MarginRight != nil {
		if err := writer.WriteField(FieldMarginRight, strconv.FormatFloat(*props.MarginRight, 'f', -1, 64)); err != nil {
			return err
		}
	}

	if props.PreferCSSPageSize != nil {
		if err := writer.WriteField(FieldPreferCSSPageSize, strconv.FormatBool(*props.PreferCSSPageSize)); err != nil {
			return err
		}
	}

	if props.GenerateDocumentOutline != nil {
		if err := writer.WriteField(FieldGenerateDocumentOutline, strconv.FormatBool(*props.GenerateDocumentOutline)); err != nil {
			return err
		}
	}

	if props.GenerateTaggedPDF != nil {
		if err := writer.WriteField(FieldGenerateTaggedPDF, strconv.FormatBool(*props.GenerateTaggedPDF)); err != nil {
			return err
		}
	}

	if props.PrintBackground != nil {
		if err := writer.WriteField(FieldPrintBackground, strconv.FormatBool(*props.PrintBackground)); err != nil {
			return err
		}
	}

	if props.OmitBackground != nil {
		if err := writer.WriteField(FieldOmitBackground, strconv.FormatBool(*props.OmitBackground)); err != nil {
			return err
		}
	}

	if props.Landscape != nil {
		if err := writer.WriteField(FieldLandscape, strconv.FormatBool(*props.Landscape)); err != nil {
			return err
		}
	}

	if props.Scale != nil {
		if err := writer.WriteField(FieldScale, strconv.FormatFloat(*props.Scale, 'f', -1, 64)); err != nil {
			return err
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
