package gotenberg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	ConvertHTML = "/forms/chromium/convert/html"
	ConvertURL  = "/forms/chromium/convert/url"
)

const (
	FieldSinglePage              = "singlePage"
	FieldPaperWidth              = "paperWidth"
	FieldPaperHeight             = "paperHeight"
	FieldMarginTop               = "marginTop"
	FieldMarginBottom            = "marginBottom"
	FieldMarginLeft              = "marginLeft"
	FieldMarginRight             = "marginRight"
	FieldPreferCSSPageSize       = "preferCssPageSize"
	FieldGenerateDocumentOutline = "generateDocumentOutline"
	FieldGenerateTaggedPDF       = "generateTaggedPdf"
	FieldPrintBackground         = "printBackground"
	FieldOmitBackground          = "omitBackground"
	FieldLandscape               = "landscape"
	FieldScale                   = "scale"
	FieldNativePageRanges        = "nativePageRanges"
)

const (
	HeaderWebhookURL              = "Gotenberg-Webhook-Url"
	HeaderWebhookErrorURL         = "Gotenberg-Webhook-Error-Url"
	HeaderWebhookMethod           = "Gotenberg-Webhook-Method"
	HeaderWebhookErrorMethod      = "Gotenberg-Webhook-Error-Method"
	HeaderWebhookExtraHTTPHeaders = "Gotenberg-Webhook-Extra-Http-Headers"
)

var (
	PaperSizeLetter  = [2]float64{8.5, 11}
	PaperSizeLegal   = [2]float64{8.5, 14}
	PaperSizeTabloid = [2]float64{11, 17}
	PaperSizeLedger  = [2]float64{17, 11}
	PaperSizeA0      = [2]float64{33.1, 46.8}
	PaperSizeA1      = [2]float64{23.4, 33.1}
	PaperSizeA2      = [2]float64{16.54, 23.4}
	PaperSizeA3      = [2]float64{11.7, 16.54}
	PaperSizeA4      = [2]float64{8.27, 11.7}
	PaperSizeA5      = [2]float64{5.83, 8.27}
	PaperSizeA6      = [2]float64{4.13, 5.83}
)

const (
	FieldURL       = "url"
	FieldFiles     = "files"
	FileIndexHTML  = "index.html"
	FileFooterHTML = "footer.html"
	FileHeaderHTML = "header.html"
	FileStylesCSS  = "styles.css"
)

func (r *Client) ConvertHTML(ctx context.Context, html io.Reader) *Client {
	return r.MethodPost(ctx, ConvertHTML).File(FieldFiles, FileIndexHTML, html)
}

func (r *Client) ConvertURL(ctx context.Context, url string) *Client {
	return r.MethodPost(ctx, ConvertURL).FormField(FieldURL, url)
}

func (r *Client) WebhookURLMethodPost(url string) *Client {
	return r.Header(HeaderWebhookURL, url).Header(HeaderWebhookMethod, http.MethodPost)
}

func (r *Client) WebhookErrorURLMethodPost(url string) *Client {
	return r.Header(HeaderWebhookErrorURL, url).Header(HeaderWebhookErrorMethod, http.MethodPost)
}

func (r *Client) WebhookExtraHeaders(headers map[string]string) *Client {
	if r.err != nil {
		return r
	}

	jsonHeaders, err := json.Marshal(headers)
	if err != nil {
		r.err = fmt.Errorf("failed to marshal webhook headers: %w", err)
		return r
	}

	return r.Header(HeaderWebhookExtraHTTPHeaders, string(jsonHeaders))
}

func (r *Client) Bool(fieldName string, value bool) *Client {
	return r.FormField(fieldName, fmt.Sprintf("%t", value))
}

func (r *Client) Float(fieldName string, value float64) *Client {
	return r.FormField(fieldName, fmt.Sprintf("%g", value))
}

func (r *Client) PaperSize(width, height float64) *Client {
	return r.Float(FieldPaperWidth, width).
		Float(FieldPaperHeight, height)
}

func (r *Client) PaperSizeA4() *Client {
	return r.PaperSize(PaperSizeA4[0], PaperSizeA4[1])
}

func (r *Client) PaperSizeLetter() *Client {
	return r.PaperSize(PaperSizeLetter[0], PaperSizeLetter[1])
}

func (r *Client) Margins(top, right, bottom, left float64) *Client {
	return r.Float(FieldMarginTop, top).
		Float(FieldMarginRight, right).
		Float(FieldMarginBottom, bottom).
		Float(FieldMarginLeft, left)
}
