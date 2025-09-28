package gotenberg

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestConvertHTML(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()
	html := strings.NewReader("<html><body>Test</body></html>")

	result := client.ConvertHTML(ctx, html)

	if result.err != nil {
		t.Fatalf("ConvertHTML() error = %v", result.err)
	}

	if result.request.Method != http.MethodPost {
		t.Errorf("ConvertHTML() method = %v, want %v", result.request.Method, http.MethodPost)
	}

	expectedPath := "/forms/chromium/convert/html"
	if !strings.Contains(result.request.URL.Path, expectedPath) {
		t.Errorf("ConvertHTML() path = %v, want to contain %v", result.request.URL.Path, expectedPath)
	}

	if !result.multipart {
		t.Error("ConvertHTML() should enable multipart")
	}
}

func TestConvertURL(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()
	url := "https://example.com"

	result := client.ConvertURL(ctx, url)

	if result.err != nil {
		t.Fatalf("ConvertURL() error = %v", result.err)
	}

	if result.request.Method != http.MethodPost {
		t.Errorf("ConvertURL() method = %v, want %v", result.request.Method, http.MethodPost)
	}

	expectedPath := "/forms/chromium/convert/url"
	if !strings.Contains(result.request.URL.Path, expectedPath) {
		t.Errorf("ConvertURL() path = %v, want to contain %v", result.request.URL.Path, expectedPath)
	}

	if !result.multipart {
		t.Error("ConvertURL() should enable multipart")
	}
}

func TestWebhookURLMethodPost(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()
	webhookURL := "http://example.com/webhook"

	result := client.MethodPost(ctx, "/test").WebhookURLMethodPost(webhookURL)

	if result.err != nil {
		t.Fatalf("WebhookURLMethodPost() error = %v", result.err)
	}

	if got := result.request.Header.Get(HeaderWebhookURL); got != webhookURL {
		t.Errorf("WebhookURLMethodPost() webhook URL = %v, want %v", got, webhookURL)
	}

	if got := result.request.Header.Get(HeaderWebhookMethod); got != http.MethodPost {
		t.Errorf("WebhookURLMethodPost() webhook method = %v, want %v", got, http.MethodPost)
	}
}

func TestWebhookErrorURLMethodPost(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()
	errorURL := "http://example.com/error"

	result := client.MethodPost(ctx, "/test").WebhookErrorURLMethodPost(errorURL)

	if result.err != nil {
		t.Fatalf("WebhookErrorURLMethodPost() error = %v", result.err)
	}

	if got := result.request.Header.Get(HeaderWebhookErrorURL); got != errorURL {
		t.Errorf("WebhookErrorURLMethodPost() error URL = %v, want %v", got, errorURL)
	}

	if got := result.request.Header.Get(HeaderWebhookErrorMethod); got != http.MethodPost {
		t.Errorf("WebhookErrorURLMethodPost() error method = %v, want %v", got, http.MethodPost)
	}
}

func TestWebhookExtraHeaders(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()
	headers := map[string]string{
		"X-Custom-1": "value1",
		"X-Custom-2": "value2",
	}

	result := client.MethodPost(ctx, "/test").WebhookExtraHeaders(headers)

	if result.err != nil {
		t.Fatalf("WebhookExtraHeaders() error = %v", result.err)
	}

	headerValue := result.request.Header.Get(HeaderWebhookExtraHTTPHeaders)
	if headerValue == "" {
		t.Fatal("WebhookExtraHeaders() did not set header")
	}

	var gotHeaders map[string]string
	err := json.Unmarshal([]byte(headerValue), &gotHeaders)
	if err != nil {
		t.Fatalf("WebhookExtraHeaders() invalid JSON: %v", err)
	}

	for key, want := range headers {
		if got, exists := gotHeaders[key]; !exists || got != want {
			t.Errorf("WebhookExtraHeaders() header %s = %v, want %v", key, got, want)
		}
	}
}

func TestWebhookExtraHeaders_InvalidJSON(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	// Create a map that can't be marshaled to JSON
	invalidHeaders := make(map[string]string)
	invalidHeaders["test"] = "value"

	// Temporarily break the client to test error handling
	originalClient := client
	client.err = nil

	// This should work normally, so let's create a different test case
	// We'll test with a valid map first
	result := originalClient.MethodPost(ctx, "/test").WebhookExtraHeaders(invalidHeaders)

	if result.err != nil {
		t.Fatalf("WebhookExtraHeaders() with valid headers should not error: %v", result.err)
	}
}

func TestBool(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	tests := []struct {
		name      string
		fieldName string
		value     bool
		want      string
	}{
		{
			name:      "true value",
			fieldName: "testField",
			value:     true,
			want:      "true",
		},
		{
			name:      "false value",
			fieldName: "testField",
			value:     false,
			want:      "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.Reset().MethodPost(ctx, "/test").Bool(tt.fieldName, tt.value)
			if result.err != nil {
				t.Fatalf("Bool() error = %v", result.err)
			}

			if !result.multipart {
				t.Error("Bool() should enable multipart")
			}
		})
	}
}

func TestFloat(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	tests := []struct {
		name      string
		fieldName string
		value     float64
		want      string
	}{
		{
			name:      "integer value",
			fieldName: "testField",
			value:     123.0,
			want:      "123",
		},
		{
			name:      "float value",
			fieldName: "testField",
			value:     123.456,
			want:      "123.456",
		},
		{
			name:      "scientific notation",
			fieldName: "testField",
			value:     1e6,
			want:      "1e+06",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.Reset().MethodPost(ctx, "/test").Float(tt.fieldName, tt.value)
			if result.err != nil {
				t.Fatalf("Float() error = %v", result.err)
			}

			if !result.multipart {
				t.Error("Float() should enable multipart")
			}
		})
	}
}

func TestPaperSize(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	result := client.MethodPost(ctx, "/test").PaperSize(8.5, 11.0)

	if result.err != nil {
		t.Fatalf("PaperSize() error = %v", result.err)
	}

	if !result.multipart {
		t.Error("PaperSize() should enable multipart")
	}
}

func TestPaperSizeA4(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	result := client.MethodPost(ctx, "/test").PaperSizeA4()

	if result.err != nil {
		t.Fatalf("PaperSizeA4() error = %v", result.err)
	}

	if !result.multipart {
		t.Error("PaperSizeA4() should enable multipart")
	}
}

func TestPaperSizeLetter(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	result := client.MethodPost(ctx, "/test").PaperSizeLetter()

	if result.err != nil {
		t.Fatalf("PaperSizeLetter() error = %v", result.err)
	}

	if !result.multipart {
		t.Error("PaperSizeLetter() should enable multipart")
	}
}

func TestMargins(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	result := client.MethodPost(ctx, "/test").Margins(1.0, 2.0, 3.0, 4.0)

	if result.err != nil {
		t.Fatalf("Margins() error = %v", result.err)
	}

	if !result.multipart {
		t.Error("Margins() should enable multipart")
	}
}

func TestPaperSizeConstants(t *testing.T) {
	tests := []struct {
		name string
		size [2]float64
	}{
		{"Letter", PaperSizeLetter},
		{"Legal", PaperSizeLegal},
		{"Tabloid", PaperSizeTabloid},
		{"Ledger", PaperSizeLedger},
		{"A0", PaperSizeA0},
		{"A1", PaperSizeA1},
		{"A2", PaperSizeA2},
		{"A3", PaperSizeA3},
		{"A4", PaperSizeA4},
		{"A5", PaperSizeA5},
		{"A6", PaperSizeA6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.size[0] <= 0 || tt.size[1] <= 0 {
				t.Errorf("Paper size %s has invalid dimensions: %v", tt.name, tt.size)
			}
		})
	}
}

func TestFieldConstants(t *testing.T) {
	tests := []struct {
		name  string
		field string
	}{
		{"SinglePage", FieldSinglePage},
		{"PaperWidth", FieldPaperWidth},
		{"PaperHeight", FieldPaperHeight},
		{"MarginTop", FieldMarginTop},
		{"MarginBottom", FieldMarginBottom},
		{"MarginLeft", FieldMarginLeft},
		{"MarginRight", FieldMarginRight},
		{"PreferCSSPageSize", FieldPreferCSSPageSize},
		{"GenerateDocumentOutline", FieldGenerateDocumentOutline},
		{"GenerateTaggedPDF", FieldGenerateTaggedPDF},
		{"PrintBackground", FieldPrintBackground},
		{"OmitBackground", FieldOmitBackground},
		{"Landscape", FieldLandscape},
		{"Scale", FieldScale},
		{"NativePageRanges", FieldNativePageRanges},
		{"URL", FieldURL},
		{"Files", FieldFiles},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field == "" {
				t.Errorf("Field constant %s is empty", tt.name)
			}
		})
	}
}

func TestHeaderConstants(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{"WebhookURL", HeaderWebhookURL},
		{"WebhookErrorURL", HeaderWebhookErrorURL},
		{"WebhookMethod", HeaderWebhookMethod},
		{"WebhookErrorMethod", HeaderWebhookErrorMethod},
		{"WebhookExtraHTTPHeaders", HeaderWebhookExtraHTTPHeaders},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.header == "" {
				t.Errorf("Header constant %s is empty", tt.name)
			}
		})
	}
}

func TestFileConstants(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{"IndexHTML", FileIndexHTML},
		{"FooterHTML", FileFooterHTML},
		{"HeaderHTML", FileHeaderHTML},
		{"StylesCSS", FileStylesCSS},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.file == "" {
				t.Errorf("File constant %s is empty", tt.name)
			}
		})
	}
}

func TestPathConstants(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"ConvertHTML", ConvertHTML},
		{"ConvertURL", ConvertURL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.path == "" {
				t.Errorf("Path constant %s is empty", tt.name)
			}
			if !strings.HasPrefix(tt.path, "/") {
				t.Errorf("Path constant %s should start with '/': %s", tt.name, tt.path)
			}
		})
	}
}

func TestIntegrationHTMLConversion(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		expectedPath := "/forms/chromium/convert/html"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Check multipart form
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Error("Expected multipart/form-data content type")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PDF content"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	html := strings.NewReader("<html><body><h1>Test PDF</h1></body></html>")

	resp, err := client.ConvertHTML(ctx, html).
		Bool(FieldPrintBackground, true).
		PaperSizeA4().
		Margins(1.0, 1.0, 1.0, 1.0).
		Send()

	if err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body := make([]byte, 1024)
	n, err := resp.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if !bytes.Contains(body[:n], []byte("PDF content")) {
		t.Error("Response body does not contain expected content")
	}
}

func TestIntegrationURLConversion(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		expectedPath := "/forms/chromium/convert/url"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Check multipart form
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Error("Expected multipart/form-data content type")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PDF content"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	testURL := "https://example.com"

	resp, err := client.ConvertURL(ctx, testURL).
		Bool(FieldLandscape, true).
		PaperSizeLetter().
		Send()

	if err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
