# gotenberg-client

[![Go Reference](https://pkg.go.dev/badge/github.com/nativebpm/gotenberg-client.svg)](https://pkg.go.dev/github.com/nativebpm/gotenberg-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/nativebpm/gotenberg-client)](https://goreportcard.com/report/github.com/nativebpm/gotenberg-client)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A high-performance, streaming Go client for the [Gotenberg](https://gotenberg.dev/) HTTP API, optimized for minimal memory allocations and maximum throughput.

## 🚀 Performance Highlights

Based on comprehensive benchmarks:
- **Ultra-low allocations**: `Reset()` method with 0 allocs/op
- **High throughput**: 5.5M+ operations/sec for client creation
- **Memory efficient**: Optimized buffer pooling with sync.Pool
- **Streaming architecture**: Direct multipart.Writer without intermediate buffers

## ✨ Key Features

- 🎯 **Fluent Builder Pattern**: Intuitive method chaining for clean code
- ♻️ **Client Reuse**: `Reset()` method for processing multiple documents efficiently
- 📦 **Zero Dependencies**: Uses only Go standard library
- 🔄 **Full Webhook Support**: Async processing with comprehensive callback handling
- ⚡ **Context Support**: Proper timeout and cancellation handling
- 🛡️ **Type Safety**: Strongly typed constants and methods
- 📊 **Buffer Pooling**: Memory-optimized with sync.Pool for high-load scenarios

## 📊 Benchmark Results

```
BenchmarkClient_Reset-12                    1000000000    1.223 ns/op     0 B/op    0 allocs/op
BenchmarkNewClient-12                        5547428       217.5 ns/op     256 B/op  2 allocs/op
BenchmarkClient_Send_GET-12                  3638         300905 ns/op    6358 B/op  79 allocs/op
BenchmarkHTMLConversion_EndToEnd-12          2292         530779 ns/op   23504 B/op 248 allocs/op
BenchmarkRealWorldHTMLUsage-12               2070         501488 ns/op   27102 B/op 294 allocs/op
```

## 🔧 Quick Start

```go
package main

import (
    "context"
    "log"
    "net/http"
    "strings"
    "time"
    
    "github.com/nativebpm/gotenberg-client"
)

func main() {
    // Create HTTP client with timeout
    httpClient := &http.Client{
        Timeout: 30 * time.Second,
    }
    
    // Initialize Gotenberg client
    client, err := gotenberg.NewClient(httpClient, "http://localhost:3000")
    if err != nil {
        log.Fatal(err)
    }
    
    // Convert HTML to PDF
    html := strings.NewReader("<html><body><h1>Hello World!</h1></body></html>")
    
    resp, err := client.
        ConvertHTML(context.Background(), html).
        PaperSizeA4().
        Margins(1.0, 1.0, 1.0, 1.0).
        Bool(gotenberg.FieldPrintBackground, true).
        Send()
        
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    
    // Process PDF response...
}
```

## 📚 Advanced Usage

### Client Reuse Pattern
```go
// Efficient client reuse for multiple conversions
client, _ := gotenberg.NewClient(httpClient, baseURL)

for _, document := range documents {
    resp, err := client.Reset(). // Reset for reuse
        ConvertHTML(ctx, document.HTML).
        PaperSizeA4().
        Send()
    // Process response...
}
```

### Webhook Configuration
```go
resp, err := client.
    ConvertHTML(ctx, html).
    WebhookURLMethodPost("http://localhost:8080/success").
    WebhookErrorURLMethodPost("http://localhost:8080/error").
    WebhookExtraHeaders(map[string]string{
        "X-Request-ID": "req-12345",
        "Authorization": "Bearer token",
    }).
    Send()
```

### Paper Size Options
```go
// Predefined sizes
client.PaperSizeA4()           // 8.27 x 11.7 inches
client.PaperSizeLetter()       // 8.5 x 11 inches

// Custom size
client.PaperSize(8.5, 11.0)    // width x height in inches

// Complete configuration
client.ConvertHTML(ctx, html).
    PaperSizeA4().
    Margins(0.5, 0.5, 0.5, 0.5). // top, right, bottom, left
    Bool(gotenberg.FieldLandscape, true).
    Float(gotenberg.FieldScale, 0.8).
    Send()
```

## 🏗️ Architecture

### Core Components

- **Client**: Main HTTP client with fluent interface
- **Form**: Multipart form handling with buffer pooling  
- **Constants**: Type-safe field and header constants
- **Methods**: Chainable builder methods for all Gotenberg options

### Memory Management

- **Buffer Pooling**: Reusable buffers via sync.Pool
- **Streaming**: Direct multipart writing without intermediate allocation
- **Reset Pattern**: Zero-allocation client reuse

## 🧪 Testing & Quality

### Comprehensive Test Suite
```bash
# Run all tests
go test -v

# Run benchmarks
go test -bench=. -benchmem

# Coverage report
go test -cover
```

### Test Coverage
- ✅ Unit tests for all public methods
- ✅ Integration tests with mock servers
- ✅ Benchmark tests for performance validation
- ✅ Memory allocation profiling
- ✅ Error handling scenarios

## � API Reference

### Core Methods
- `NewClient(client *http.Client, baseURL string) (*Client, error)`
- `Reset() *Client` - Reset client for reuse
- `ConvertHTML(ctx context.Context, html io.Reader) *Client`
- `ConvertURL(ctx context.Context, url string) *Client`
- `Send() (*http.Response, error)`

### Configuration Methods
- `PaperSizeA4()`, `PaperSizeLetter()`, `PaperSize(w, h float64)`
- `Margins(top, right, bottom, left float64)`
- `Bool(field string, value bool)`, `Float(field string, value float64)`
- `WebhookURLMethodPost(url string)`, `WebhookExtraHeaders(headers map[string]string)`

### HTTP Methods
- `Header(key, value string)`, `Headers(map[string]string)`
- `QueryParam(key, value string)`, `QueryParams(map[string]string)`
- `JSONBody(data interface{})`, `StringBody(body string)`
- `FormField(name, value string)`, `File(field, filename string, content io.Reader)`

## 🏆 Best Practices

- ✅ **Always close `resp.Body`** to prevent resource leaks
- ✅ **Use `context.Context`** for timeouts and cancellations
- ✅ **Call `Reset()`** when reusing client instances
- ✅ **Check response status codes** before processing
- ✅ **Set appropriate HTTP client timeouts** for your use case
- ✅ **Inspect `Gotenberg-Trace` header** for debugging
- ✅ **Use buffer pooling** for high-throughput scenarios

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `go test -v`
5. Run benchmarks: `go test -bench=.`
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
