package gotenberg

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockRoundTripper implements http.RoundTripper for testing
type mockRoundTripper struct{}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Read and discard the request body to allow multipart goroutine to finish cleanly
	if req.Body != nil {
		io.Copy(io.Discard, req.Body)
		req.Body.Close()
	}
	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("pdf-bytes")),
	}
	resp.Header.Set("Gotenberg-Trace", "trace-id")
	return resp, nil
}

func newTestClient(t *testing.T) *Client {
	httpCli := &http.Client{Transport: &mockRoundTripper{}}
	cli, err := NewClient(httpCli, "http://localhost")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return cli
}

func TestNewClient(t *testing.T) {
	cli, err := NewClient(&http.Client{}, "http://localhost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cli == nil {
		t.Fatal("expected client, got nil")
	}
}

func TestConvertHTML(t *testing.T) {
	c := newTestClient(t)
	buf := bytes.NewBufferString("<html></html>")
	r := c.ConvertHTML(context.Background(), buf)
	if r == nil || r.req == nil {
		t.Fatal("expected request, got nil")
	}
	// Send to trigger the request and check response
	resp, err := r.Send()
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if resp.GotenbergTrace != "trace-id" {
		t.Errorf("expected trace-id, got %s", resp.GotenbergTrace)
	}
}

func TestConvertURL(t *testing.T) {
	c := newTestClient(t)
	r := c.ConvertURL(context.Background(), "http://example.com")
	if r == nil || r.req == nil {
		t.Fatal("expected request, got nil")
	}
	resp, err := r.Send()
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if resp.GotenbergTrace != "trace-id" {
		t.Errorf("expected trace-id, got %s", resp.GotenbergTrace)
	}
}

func TestRequestSend(t *testing.T) {
	c := newTestClient(t)
	r := c.ConvertURL(context.Background(), "http://example.com")
	resp, err := r.Send()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.GotenbergTrace != "trace-id" {
		t.Errorf("expected trace-id, got %s", resp.GotenbergTrace)
	}
}

func TestRequestHeaderParamBoolFloatFile(t *testing.T) {
	c := newTestClient(t)
	r := c.ConvertHTML(context.Background(), bytes.NewBufferString("<html></html>"))
	r.Header("X-Test", "v1").Param("p", "v2").Bool("b", true).Float("f", 1.23).File("k", "f.txt", strings.NewReader("x"))
	// Just check that Send does not error
	_, err := r.Send()
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

func TestWebhookURL(t *testing.T) {
	c := newTestClient(t)
	r := c.ConvertHTML(context.Background(), bytes.NewBufferString("<html></html>"))
	r.WebhookURL("http://webhook", "POST")
	_, err := r.Send()
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

func TestOutputFilename(t *testing.T) {
	c := newTestClient(t)
	r := c.ConvertHTML(context.Background(), bytes.NewBufferString("<html></html>"))
	r.OutputFilename("out.pdf")
	_, err := r.Send()
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

func TestWebhookErrorURL(t *testing.T) {
	c := newTestClient(t)
	r := c.ConvertHTML(context.Background(), bytes.NewBufferString("<html></html>"))
	r.WebhookErrorURL("http://err", "PUT")
	_, err := r.Send()
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

func TestWebhookHeaders(t *testing.T) {
	c := newTestClient(t)
	r := c.ConvertHTML(context.Background(), bytes.NewBufferString("<html></html>"))
	r.WebhookHeader("X-Test", "v1")
	_, err := r.Send()
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

func TestPaperSize(t *testing.T) {
	c := newTestClient(t)
	r := c.ConvertHTML(context.Background(), bytes.NewBufferString("<html></html>"))
	r.PaperSize(1.1, 2.2)
	_, err := r.Send()
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

func TestPaperSizeA4(t *testing.T) {
	c := newTestClient(t)
	r := c.ConvertHTML(context.Background(), bytes.NewBufferString("<html></html>"))
	r.PaperSizeA4()
	_, err := r.Send()
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

func TestPaperSizeLetter(t *testing.T) {
	c := newTestClient(t)
	r := c.ConvertHTML(context.Background(), bytes.NewBufferString("<html></html>"))
	r.PaperSizeLetter()
	_, err := r.Send()
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

func TestMargins(t *testing.T) {
	c := newTestClient(t)
	r := c.ConvertHTML(context.Background(), bytes.NewBufferString("<html></html>"))
	r.Margins(1, 2, 3, 4)
	_, err := r.Send()
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

// Benchmarks

func BenchmarkConvertHTML(b *testing.B) {
	c := newTestClient(nil)
	buf := bytes.NewBufferString("<html></html>")
	for i := 0; i < b.N; i++ {
		c.ConvertHTML(context.Background(), buf)
	}
}

func BenchmarkConvertURL(b *testing.B) {
	c := newTestClient(nil)
	for i := 0; i < b.N; i++ {
		c.ConvertURL(context.Background(), "http://example.com")
	}
}

func BenchmarkConvertHTMLSend(b *testing.B) {
	c := newTestClient(nil)
	for i := 0; i < b.N; i++ {
		r := c.ConvertHTML(context.Background(), bytes.NewBufferString("<html></html>"))
		_, _ = r.Send()
	}
}

func BenchmarkConvertURLSend(b *testing.B) {
	c := newTestClient(nil)
	for i := 0; i < b.N; i++ {
		r := c.ConvertURL(context.Background(), "http://example.com")
		_, _ = r.Send()
	}
}
