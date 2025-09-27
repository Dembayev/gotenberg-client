package gotenberg

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func BenchmarkNewClient(b *testing.B) {
	httpClient := &http.Client{Timeout: 30 * time.Second}
	baseURL := "http://localhost:3000"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewClient(httpClient, baseURL)
	}
}

func BenchmarkClientChaining(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := NewClient(httpClient, "http://localhost:3000")
		html := strings.NewReader("<html><body>Benchmark Test</body></html>")
		css := strings.NewReader("body { margin: 0; padding: 0; }")

		client.
			IndexHTML(html).
			StylesCSS(css).
			PaperSizeA4().
			Margins(1.0, 1.0, 1.0, 1.0).
			Bool(FieldPrintBackground, true).
			Bool(FieldLandscape, false)
	}
}

func BenchmarkClientFile(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		content := strings.NewReader("test file content for benchmarking")
		client.File("test.txt", content)
	}
}

func BenchmarkClientString(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.String("testField", "testValue")
	}
}

func BenchmarkClientBool(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Bool(FieldPrintBackground, true)
	}
}

func BenchmarkClientFloat(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Float(FieldPaperWidth, 8.5)
	}
}

func BenchmarkClientPaperSize(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.PaperSize(8.5, 11.0)
	}
}

func BenchmarkClientMargins(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Margins(1.0, 1.0, 1.0, 1.0)
	}
}

func BenchmarkClientConvertHTML(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := NewClient(httpClient, "http://localhost:3000")
		html := strings.NewReader("<html><body>Benchmark HTML conversion</body></html>")
		client.IndexHTML(html)

		b.StartTimer()
		resp, err := client.ConvertHTML(ctx)
		b.StopTimer()

		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		resp.Body.Close()
	}
}

func BenchmarkClientConvertURL(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := NewClient(httpClient, "http://localhost:3000")
		client.URL("https://example.com")

		b.StartTimer()
		resp, err := client.ConvertURL(ctx)
		b.StopTimer()

		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		resp.Body.Close()
	}
}

func BenchmarkClientReuse(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		html := strings.NewReader("<html><body>Reuse test</body></html>")
		client.IndexHTML(html).PaperSizeA4()

		resp, err := client.ConvertHTML(ctx)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		resp.Body.Close()
	}
}

func BenchmarkBufferPool(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := client.bufPool.Get()
		client.bufPool.Put(buf)
	}
}

func BenchmarkClientFullWorkflow(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := NewClient(httpClient, "http://localhost:3000")

		html := strings.NewReader(`
			<html>
				<head>
					<link rel="stylesheet" href="styles.css">
				</head>
				<body>
					<h1>Full Workflow Benchmark</h1>
					<p>This is a comprehensive benchmark test.</p>
					<img src="logo.png" alt="Logo">
				</body>
			</html>
		`)

		css := strings.NewReader(`
			body { 
				font-family: Arial, sans-serif; 
				margin: 2cm; 
				line-height: 1.6; 
			}
			h1 { 
				color: #333; 
				border-bottom: 2px solid #007bff; 
			}
		`)

		logo := strings.NewReader("fake-png-data")

		client.
			IndexHTML(html).
			StylesCSS(css).
			File("logo.png", logo).
			PaperSizeA4().
			Margins(1.0, 1.0, 1.0, 1.0).
			Bool(FieldPrintBackground, true).
			Bool(FieldLandscape, false).
			Float(FieldScale, 1.0)

		resp, err := client.ConvertHTML(ctx)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}

		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkClientMemoryAllocation(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client := NewClient(httpClient, "http://localhost:3000")
		html := strings.NewReader("<html><body>Memory allocation test</body></html>")

		client.
			IndexHTML(html).
			PaperSizeA4().
			Bool(FieldPrintBackground, true)
	}
}

func BenchmarkPaperSizeOperations(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	client := NewClient(httpClient, "http://localhost:3000")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch i % 4 {
		case 0:
			client.PaperSizeA4()
		case 1:
			client.PaperSizeLetter()
		case 2:
			client.PaperSize(PaperSizeA3[0], PaperSizeA3[1])
		case 3:
			client.PaperSize(11.7, 16.54)
		}
	}
}

func BenchmarkClientConcurrent(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	ctx := context.Background()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client := NewClient(httpClient, "http://localhost:3000")
			html := strings.NewReader("<html><body>Concurrent test</body></html>")

			client.IndexHTML(html).PaperSizeA4()

			resp, err := client.ConvertHTML(ctx)
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
			resp.Body.Close()
		}
	})
}

func BenchmarkLargeHTMLContent(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	ctx := context.Background()

	largeHTML := strings.Repeat(`
		<div>
			<h2>Section Title</h2>
			<p>This is a paragraph with some content that repeats many times to simulate a large HTML document.</p>
			<ul>
				<li>List item 1</li>
				<li>List item 2</li>
				<li>List item 3</li>
			</ul>
		</div>
	`, 1000)

	fullHTML := "<html><head><title>Large Document</title></head><body>" + largeHTML + "</body></html>"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := NewClient(httpClient, "http://localhost:3000")
		html := strings.NewReader(fullHTML)

		client.IndexHTML(html).PaperSizeA4()

		resp, err := client.ConvertHTML(ctx)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		resp.Body.Close()
	}
}

func BenchmarkMultipleFiles(b *testing.B) {
	httpClient := newMockHTTPClient(200, "PDF content")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := NewClient(httpClient, "http://localhost:3000")

		html := strings.NewReader("<html><head><link rel='stylesheet' href='main.css'><link rel='stylesheet' href='theme.css'></head><body><img src='logo.png'><img src='banner.jpg'></body></html>")
		mainCSS := strings.NewReader("body { font-family: Arial; }")
		themeCSS := strings.NewReader("h1 { color: blue; }")
		logo := strings.NewReader("fake-png-data")
		banner := strings.NewReader("fake-jpg-data")

		client.
			IndexHTML(html).
			File("main.css", mainCSS).
			File("theme.css", themeCSS).
			File("logo.png", logo).
			File("banner.jpg", banner).
			PaperSizeA4()

		resp, err := client.ConvertHTML(ctx)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		resp.Body.Close()
	}
}
