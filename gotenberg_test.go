package gotenberg

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockRoundTripper captures the last request and body for inspection.
type mockRoundTripper struct {
	lastReq  *http.Request
	lastBody []byte
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.lastReq = req
	if req.Body != nil {
		b, _ := io.ReadAll(req.Body)
		m.lastBody = b
	} else {
		m.lastBody = nil
	}

	// Return a minimal successful response
	return &http.Response{
		StatusCode:    200,
		Body:          io.NopCloser(strings.NewReader("ok")),
		ContentLength: 2,
		Header:        make(http.Header),
		Request:       req,
	}, nil
}

func TestConvertURLToPDF_Errors(t *testing.T) {
	c := NewClient(&http.Client{}, "http://example.com")
	_, err := c.ConvertURLToPDF(context.Background(), "")
	if err == nil {
		t.Fatalf("expected error when URL is empty")
	}
}

func TestConvertHTMLToPDF_Errors(t *testing.T) {
	c := NewClient(&http.Client{}, "http://example.com")
	_, err := c.ConvertHTMLToPDF(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected error when indexHTML is empty")
	}
}

func TestConvertURLToPDF_Success(t *testing.T) {
	mrt := &mockRoundTripper{}
	client := NewClient(&http.Client{Transport: mrt}, "http://example.com")

	resp, err := client.ConvertURLToPDF(context.Background(), "https://golang.org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response, got nil")
	}

	if mrt.lastReq == nil {
		t.Fatalf("request was not captured by transport")
	}
	if mrt.lastReq.Method != "POST" {
		t.Fatalf("expected POST method, got %s", mrt.lastReq.Method)
	}
	if !strings.Contains(mrt.lastReq.URL.Path, "/forms/chromium/convert/url") {
		t.Fatalf("unexpected request path: %s", mrt.lastReq.URL.Path)
	}
	if ct := mrt.lastReq.Header.Get("Content-Type"); !strings.Contains(ct, "multipart/form-data") {
		t.Fatalf("expected multipart content type, got %s", ct)
	}
	if !bytes.Contains(mrt.lastBody, []byte("https://golang.org")) {
		t.Fatalf("request body does not contain the provided URL")
	}
}

func TestConvertHTMLToPDF_Success(t *testing.T) {
	mrt := &mockRoundTripper{}
	client := NewClient(&http.Client{Transport: mrt}, "http://example.com")

	html := []byte("<html><body>Hello</body></html>")
	resp, err := client.ConvertHTMLToPDF(context.Background(), bytes.NewReader(html))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response, got nil")
	}

	if mrt.lastReq == nil {
		t.Fatalf("request was not captured by transport")
	}
	if !strings.Contains(mrt.lastReq.URL.Path, "/forms/chromium/convert/html") {
		t.Fatalf("unexpected request path: %s", mrt.lastReq.URL.Path)
	}
	if !bytes.Contains(mrt.lastBody, []byte("filename=\"index.html\"")) {
		t.Fatalf("multipart body missing index.html filename")
	}
	if !bytes.Contains(mrt.lastBody, []byte("Hello")) {
		t.Fatalf("body does not contain HTML content")
	}
}

// Benchmarks: minimal loops calling the three conversion methods.
func BenchmarkConvertURLToPDF(b *testing.B) {
	mrt := &mockRoundTripper{}
	client := NewClient(&http.Client{Transport: mrt}, "http://example.com")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.ConvertURLToPDF(context.Background(), "https://example.com")
	}
}

func BenchmarkConvertHTMLToPDF(b *testing.B) {
	mrt := &mockRoundTripper{}
	client := NewClient(&http.Client{Transport: mrt}, "http://example.com")
	html := []byte("<html><body>Benchmark</body></html>")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.ConvertHTMLToPDF(context.Background(), bytes.NewReader(html))
	}
}

func BenchmarkOptionsBuilder(b *testing.B) {
	mrt := &mockRoundTripper{}
	clientBuilder := NewClientBuilder(&http.Client{Transport: mrt}, "http://example.com")
	html := "<html><body>Benchmark</body></html>"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = clientBuilder.ConvertHTML().
			WithHTML(html).
			PaperSizeA4().
			Margins(1.0, 1.0, 1.0, 1.0).
			Landscape(false).
			PrintBackground(true).
			Execute(context.Background())
	}
}

func TestOptionsBuilder_Success(t *testing.T) {
	mrt := &mockRoundTripper{}
	clientBuilder := NewClientBuilder(&http.Client{Transport: mrt}, "http://example.com")

	html := "<html><body>Hello Builder</body></html>"
	cssFile := bytes.NewReader([]byte("body { color: red; }"))

	// Test builder with complex configuration
	resp, err := clientBuilder.ConvertHTML().
		WithHTML(html).
		WithFile("styles.css", cssFile).
		PaperSizeA4().
		Margins(1.5, 1.0, 1.5, 1.0).
		PrintBackground(true).
		Landscape(false).
		Scale(0.8).
		OutputFilename("test-builder.pdf").
		WebhookSuccess("https://webhook.example.com/success", "POST").
		WebhookError("https://webhook.example.com/error", "POST").
		WebhookExtraHeader("Authorization", "Bearer token123").
		WebhookExtraHeader("X-Custom", "test-value").
		Execute(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response, got nil")
	}

	// Verify request was made properly
	if mrt.lastReq == nil {
		t.Fatalf("request was not captured by transport")
	}
	if !strings.Contains(mrt.lastReq.URL.Path, "/forms/chromium/convert/html") {
		t.Fatalf("unexpected request path: %s", mrt.lastReq.URL.Path)
	}

	// Check headers set by builder
	if filename := mrt.lastReq.Header.Get("Gotenberg-Output-Filename"); filename != "test-builder.pdf" {
		t.Fatalf("expected output filename test-builder.pdf, got %s", filename)
	}
	if webhookURL := mrt.lastReq.Header.Get("Gotenberg-Webhook-Url"); webhookURL != "https://webhook.example.com/success" {
		t.Fatalf("expected webhook URL, got %s", webhookURL)
	}

	// Verify multipart body contains our CSS file
	if !bytes.Contains(mrt.lastBody, []byte("filename=\"styles.css\"")) {
		t.Fatalf("multipart body missing styles.css filename")
	}
	if !bytes.Contains(mrt.lastBody, []byte("color: red;")) {
		t.Fatalf("body does not contain CSS content")
	}
}
