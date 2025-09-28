package gotenberg

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type mockRoundTripper struct {
	response *http.Response
	err      error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func newMockHTTPClient(statusCode int, body string) *http.Client {
	response := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
	response.Header.Set("Content-Type", "application/pdf")
	response.Header.Set("Gotenberg-Trace", "test-trace-id")

	return &http.Client{
		Transport: &mockRoundTripper{response: response},
	}
}

func TestNewClient(t *testing.T) {
	httpClient := &http.Client{Timeout: 30 * time.Second}
	baseURL := "http://localhost:3000"

	client := NewClient(httpClient, baseURL)

	if client.httpClient != httpClient {
		t.Error("httpClient not set correctly")
	}

	expectedURL := "http://localhost:3000"
	if client.baseURL.String() != expectedURL {
		t.Errorf("baseURL = %s, expected %s", client.baseURL.String(), expectedURL)
	}

	if client.buffer == nil {
		t.Error("buffer not initialized")
	}

	if client.writer == nil {
		t.Error("writer not initialized")
	}
}

func TestNewClientWithTrailingSlash(t *testing.T) {
	httpClient := &http.Client{Timeout: 30 * time.Second}
	baseURL := "http://localhost:3000/"

	client := NewClient(httpClient, baseURL)

	expectedURL := "http://localhost:3000"
	if client.baseURL.String() != expectedURL {
		t.Errorf("baseURL = %s, expected %s", client.baseURL.String(), expectedURL)
	}
}

func TestNewClientInvalidURL(t *testing.T) {
	httpClient := &http.Client{Timeout: 30 * time.Second}
	client := NewClient(httpClient, "://invalid-url")

	if client.err == nil {
		t.Error("expected error for invalid URL")
	}

	if !strings.Contains(client.err.Error(), "invalid base URL") {
		t.Errorf("expected 'invalid base URL' error, got %v", client.err)
	}
}

func TestClientFile(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	content := strings.NewReader("test content")
	result := client.File("test.txt", content)

	if result != client {
		t.Error("File should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientString(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	result := client.String("testField", "testValue")

	if result != client {
		t.Error("String should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientURL(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	result := client.URL("https://example.com")

	if result != client {
		t.Error("URL should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientIndexHTML(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	html := strings.NewReader("<html><body>Test</body></html>")
	result := client.IndexHTML(html)

	if result != client {
		t.Error("IndexHTML should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientFooterHTML(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	html := strings.NewReader("<div>Footer</div>")
	result := client.FooterHTML(html)

	if result != client {
		t.Error("FooterHTML should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientHeaderHTML(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	html := strings.NewReader("<div>Header</div>")
	result := client.HeaderHTML(html)

	if result != client {
		t.Error("HeaderHTML should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientStylesCSS(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	css := strings.NewReader("body { margin: 0; }")
	result := client.StylesCSS(css)

	if result != client {
		t.Error("StylesCSS should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientBool(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	result := client.Bool(FieldPrintBackground, true)

	if result != client {
		t.Error("Bool should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientFloat(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	result := client.Float(FieldPaperWidth, 8.5)

	if result != client {
		t.Error("Float should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientPaperSize(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	result := client.PaperSize(8.5, 11.0)

	if result != client {
		t.Error("PaperSize should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientPaperSizeA4(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	result := client.PaperSizeA4()

	if result != client {
		t.Error("PaperSizeA4 should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientPaperSizeLetter(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	result := client.PaperSizeLetter()

	if result != client {
		t.Error("PaperSizeLetter should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientMargins(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	result := client.Margins(1.0, 1.0, 1.0, 1.0)

	if result != client {
		t.Error("Margins should return client for chaining")
	}

	if client.err != nil {
		t.Errorf("unexpected error: %v", client.err)
	}
}

func TestClientConvertHTML(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	html := strings.NewReader("<html><body>Test</body></html>")
	client.IndexHTML(html)

	ctx := context.Background()
	resp, err := client.ConvertHTML(ctx)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/pdf" {
		t.Errorf("expected PDF content type")
	}
}

func TestClientConvertURL(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	client.URL("https://example.com")

	ctx := context.Background()
	resp, err := client.ConvertURL(ctx)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestClientExecuteUnknownRoute(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	ctx := context.Background()
	_, err := client.Execute(ctx, "unknown")

	if err == nil {
		t.Error("expected error for unknown route")
	}

	expectedError := "unknown route: unknown"
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestClientWithPreviousError(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	client.err = io.ErrUnexpectedEOF

	result := client.String("test", "value")
	if result != client {
		t.Error("should return client even with error")
	}

	if client.err != io.ErrUnexpectedEOF {
		t.Error("should preserve original error")
	}
}

func TestClientExecuteWithError(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	client.err = io.ErrUnexpectedEOF

	ctx := context.Background()
	_, err := client.ConvertHTML(ctx)

	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected preserved error, got %v", err)
	}
}

func TestClientChaining(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	html := strings.NewReader("<html><body>Test</body></html>")
	css := strings.NewReader("body { margin: 0; }")

	result := client.
		IndexHTML(html).
		StylesCSS(css).
		PaperSizeA4().
		Margins(1.0, 1.0, 1.0, 1.0).
		Bool(FieldPrintBackground, true).
		Bool(FieldLandscape, false)

	if result != client {
		t.Error("chaining should return same client instance")
	}

	if client.err != nil {
		t.Errorf("unexpected error in chain: %v", client.err)
	}

	ctx := context.Background()
	resp, err := client.ConvertHTML(ctx)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestPaperSizeConstants(t *testing.T) {
	testCases := []struct {
		name     string
		size     [2]float64
		expected [2]float64
	}{
		{"Letter", PaperSizeLetter, [2]float64{8.5, 11}},
		{"Legal", PaperSizeLegal, [2]float64{8.5, 14}},
		{"Tabloid", PaperSizeTabloid, [2]float64{11, 17}},
		{"Ledger", PaperSizeLedger, [2]float64{17, 11}},
		{"A0", PaperSizeA0, [2]float64{33.1, 46.8}},
		{"A1", PaperSizeA1, [2]float64{23.4, 33.1}},
		{"A2", PaperSizeA2, [2]float64{16.54, 23.4}},
		{"A3", PaperSizeA3, [2]float64{11.7, 16.54}},
		{"A4", PaperSizeA4, [2]float64{8.27, 11.7}},
		{"A5", PaperSizeA5, [2]float64{5.83, 8.27}},
		{"A6", PaperSizeA6, [2]float64{4.13, 5.83}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.size[0] != tc.expected[0] || tc.size[1] != tc.expected[1] {
				t.Errorf("paper size %s = [%v, %v], expected [%v, %v]",
					tc.name, tc.size[0], tc.size[1], tc.expected[0], tc.expected[1])
			}
		})
	}
}

func TestClientReset(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	html := strings.NewReader("<html><body>Test</body></html>")
	client.IndexHTML(html)
	client.Bool(FieldPrintBackground, true)

	ctx := context.Background()
	_, err := client.ConvertHTML(ctx)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if client.buffer.Len() != 0 {
		t.Error("buffer should be reset after execute")
	}

	if client.err != nil {
		t.Error("error should be reset after execute")
	}
}

func TestServerIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("expected multipart/form-data content type")
		}

		if r.URL.Path != "/forms/chromium/convert/html" {
			t.Errorf("expected HTML conversion path, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Gotenberg-Trace", "integration-test")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PDF content"))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 30 * time.Second}
	client := NewClient(httpClient, server.URL)

	html := strings.NewReader("<html><body>Integration Test</body></html>")
	client.IndexHTML(html).PaperSizeA4()

	ctx := context.Background()
	resp, err := client.ConvertHTML(ctx)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if resp.Header.Get("Gotenberg-Trace") != "integration-test" {
		t.Error("expected Gotenberg-Trace header")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("error reading response body: %v", err)
	}

	if string(body) != "PDF content" {
		t.Errorf("expected 'PDF content', got %s", string(body))
	}
}

func TestClientContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 30 * time.Second}
	client := NewClient(httpClient, server.URL)

	html := strings.NewReader("<html><body>Test</body></html>")
	client.IndexHTML(html)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.ConvertHTML(ctx)

	if err == nil {
		t.Error("expected timeout error")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("expected context deadline exceeded error, got %v", err)
	}
}

func TestClientBufferPool(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	if client.bufPool.New == nil {
		t.Error("buffer pool not initialized")
	}

	buf := client.bufPool.Get()
	if buf == nil {
		t.Error("buffer pool should return a buffer")
	}

	client.bufPool.Put(buf)

	buf2 := client.bufPool.Get()
	if buf2 != buf {
		t.Error("buffer pool should reuse buffers")
	}
}

func TestClientError(t *testing.T) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	if client.Err() != nil {
		t.Error("new client should have no error")
	}

	client.err = io.ErrUnexpectedEOF

	if client.Err() != io.ErrUnexpectedEOF {
		t.Error("ClientError should return the set error")
	}
}

func TestInvalidBaseURL(t *testing.T) {
	httpClient := &http.Client{}
	client := NewClient(httpClient, "ht!tp://invalid")

	if client.err == nil {
		t.Error("expected error for invalid base URL")
	}
}
