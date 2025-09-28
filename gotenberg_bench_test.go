package gotenberg

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func BenchmarkConvertHTML(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	html := strings.NewReader("<html><body>Test</body></html>")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		html.Reset("<html><body>Test</body></html>")
		result := client.Reset().ConvertHTML(ctx, html)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkConvertURL(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	url := "https://example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().ConvertURL(ctx, url)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkWebhookURLMethodPost(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	webhookURL := "http://example.com/webhook"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").WebhookURLMethodPost(webhookURL)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkWebhookErrorURLMethodPost(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	errorURL := "http://example.com/error"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").WebhookErrorURLMethodPost(errorURL)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkWebhookExtraHeaders(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	headers := map[string]string{
		"X-Custom-1": "value1",
		"X-Custom-2": "value2",
		"X-Auth":     "token123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").WebhookExtraHeaders(headers)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkBool(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").Bool("testField", true)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkFloat(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").Float("testField", 123.456)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkPaperSize(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").PaperSize(8.5, 11.0)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkPaperSizeA4(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").PaperSizeA4()
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkPaperSizeLetter(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").PaperSizeLetter()
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkMargins(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").Margins(1.0, 2.0, 3.0, 4.0)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkCompleteHTMLConversion_Chain(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	html := strings.NewReader("<html><body>Test</body></html>")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		html.Reset("<html><body>Test</body></html>")
		result := client.Reset().
			ConvertHTML(ctx, html).
			Bool(FieldPrintBackground, true).
			PaperSizeA4().
			Margins(1.0, 1.0, 1.0, 1.0).
			WebhookURLMethodPost("http://example.com/webhook")

		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkCompleteURLConversion_Chain(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	url := "https://example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().
			ConvertURL(ctx, url).
			Bool(FieldLandscape, true).
			PaperSizeLetter().
			Float(FieldScale, 0.8).
			WebhookErrorURLMethodPost("http://example.com/error")

		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkHTMLConversion_EndToEnd(b *testing.B) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PDF content"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	html := strings.NewReader("<html><body><h1>Test PDF</h1></body></html>")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		html.Reset("<html><body><h1>Test PDF</h1></body></html>")
		resp, err := client.Reset().
			ConvertHTML(ctx, html).
			Bool(FieldPrintBackground, true).
			PaperSizeA4().
			Margins(1.0, 1.0, 1.0, 1.0).
			Send()

		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkURLConversion_EndToEnd(b *testing.B) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PDF content"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	testURL := "https://example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.Reset().
			ConvertURL(ctx, testURL).
			Bool(FieldLandscape, true).
			PaperSizeLetter().
			Send()

		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkWebhookConfiguration(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	headers := map[string]string{
		"X-Custom-Header": "MyValue",
		"Authorization":   "Bearer token123",
		"X-Request-ID":    "req-456",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().
			MethodPost(ctx, "/test").
			WebhookURLMethodPost("http://host.docker.internal:28080/success").
			WebhookErrorURLMethodPost("http://host.docker.internal:28080/error").
			WebhookExtraHeaders(headers)

		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

// Memory allocation benchmarks
func BenchmarkHTMLConversion_MemoryAllocation(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	html := strings.NewReader("<html><body>Test</body></html>")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		html.Reset("<html><body>Test</body></html>")
		result := client.Reset().
			ConvertHTML(ctx, html).
			Bool(FieldPrintBackground, true).
			PaperSizeA4()

		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkURLConversion_MemoryAllocation(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	url := "https://example.com"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().
			ConvertURL(ctx, url).
			Bool(FieldLandscape, true).
			PaperSizeLetter()

		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkWebhookSetup_MemoryAllocation(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	headers := map[string]string{
		"X-Custom": "value",
		"X-Auth":   "token",
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().
			MethodPost(ctx, "/test").
			WebhookURLMethodPost("http://example.com/webhook").
			WebhookExtraHeaders(headers)

		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

// Comparative benchmarks
func BenchmarkPaperSize_Manual_vs_A4(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.Run("Manual", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := client.Reset().MethodPost(ctx, "/test").PaperSize(8.27, 11.7)
			if result.err != nil {
				b.Fatal(result.err)
			}
		}
	})

	b.Run("A4", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := client.Reset().MethodPost(ctx, "/test").PaperSizeA4()
			if result.err != nil {
				b.Fatal(result.err)
			}
		}
	})
}

func BenchmarkPaperSize_Manual_vs_Letter(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.Run("Manual", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := client.Reset().MethodPost(ctx, "/test").PaperSize(8.5, 11)
			if result.err != nil {
				b.Fatal(result.err)
			}
		}
	})

	b.Run("Letter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := client.Reset().MethodPost(ctx, "/test").PaperSizeLetter()
			if result.err != nil {
				b.Fatal(result.err)
			}
		}
	})
}

// Real-world usage simulation
func BenchmarkRealWorldHTMLUsage(b *testing.B) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write(bytes.Repeat([]byte("P"), 1024)) // Simulate 1KB PDF
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 10 * time.Second}, server.URL)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	htmlTemplate := `
		<html>
		<head><title>Invoice</title></head>
		<body>
			<h1>Invoice #12345</h1>
			<p>Customer: John Doe</p>
			<table>
				<tr><td>Item</td><td>Price</td></tr>
				<tr><td>Service A</td><td>$100</td></tr>
				<tr><td>Service B</td><td>$200</td></tr>
			</table>
			<p>Total: $300</p>
		</body>
		</html>
	`
	html := strings.NewReader(htmlTemplate)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		html.Reset(htmlTemplate)
		resp, err := client.Reset().
			ConvertHTML(ctx, html).
			Bool(FieldPrintBackground, true).
			Bool(FieldGenerateDocumentOutline, true).
			PaperSizeA4().
			Margins(0.5, 0.5, 0.5, 0.5).
			WebhookURLMethodPost("http://localhost:8080/success").
			WebhookErrorURLMethodPost("http://localhost:8080/error").
			WebhookExtraHeaders(map[string]string{
				"X-Request-ID": "req-12345",
				"X-User-ID":    "user-67890",
			}).
			Send()

		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
