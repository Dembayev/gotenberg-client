# gotenberg-client

[![Go Reference](https://pkg.go.dev/badge/github.com/nativebpm/gotenberg-client.svg)](https://pkg.go.dev/github.com/nativebpm/gotenberg-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/nativebpm/gotenberg-client)](https://goreportcard.com/report/github.com/nativebpm/gotenberg-client)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A clean, performant and idiomatic Go client for the [Gotenberg](https://gotenberg.dev/) HTTP API.

## Features

- 🚀 **High Performance**: Optimized buffer pools and minimal allocations
- 🎯 **Fluent API**: Modern builder pattern for readable configuration  
- 🔄 **Streaming Support**: Handle large PDFs without memory buffering
- 📦 **Zero Dependencies**: Uses only Go standard library
- 🔗 **Webhook Support**: Full async processing with webhooks
- ⚡ **Context Support**: Proper timeout and cancellation handling
- 🛡️ **Type Safe**: Compile-time validation of configurations

## Installation

```bash
go get github.com/nativebpm/gotenberg-client
```

**Requirements**: Go 1.21 or later

## Quick Start

```go
package main

import (
    "context"
    "net/http"
    "time"
    
    "github.com/nativebpm/gotenberg-client"
)

func main() {
    client := gotenberg.NewClient(
        &http.Client{Timeout: 30 * time.Second},
        "http://localhost:3000",
    )
    
    resp, err := client.ConvertURLToPDF(context.Background(), "https://example.com")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    // Save or stream the PDF...
}
```

## Usage Examples

### Simple URL to PDF conversion

```go
ctx := context.Background()
httpClient := &http.Client{Timeout: 30 * time.Second}
cli := gotenberg.NewClient(httpClient, "http://localhost:3000")

resp, err := cli.ConvertURLToPDF(ctx, "https://example.com")
if err != nil {
    // handle error
}
defer resp.Body.Close()

// stream or save the PDF
// out, _ := os.Create("out.pdf")
// io.Copy(out, resp.Body)
```

### HTML to PDF with Builder Pattern (Recommended)

```go
htmlContent := "<html><body><h1>Hello World</h1></body></html>"
cssContent := "h1 { color: blue; }"

// Use fluent builder pattern for clean, readable configuration
options := gotenberg.NewOptionsBuilder().
    File("styles.css", bytes.NewReader([]byte(cssContent))).
    PaperSizeA4().
    Margins(1.0, 1.0, 1.0, 1.0).
    PrintBackground(true).
    Scale(0.9).
    OutputFilename("document.pdf").
    Build()

resp, err := cli.ConvertHTMLToPDF(ctx, bytes.NewReader([]byte(htmlContent)), options)
if err != nil {
    // handle error
}
defer resp.Body.Close()
```



### Webhook Configuration with Builder

```go
options := gotenberg.NewOptionsBuilder().
    PaperSizeLetter().
    WebhookSuccess("https://your-domain.com/webhook/success", "POST").
    WebhookError("https://your-domain.com/webhook/error", "POST").
    WebhookExtraHeader("Authorization", "Bearer your-token").
    OutputFilename("async-document.pdf").
    Build()

resp, err := cli.ConvertURLToPDF(ctx, "https://example.com", options)
```

## API Reference

### Client Methods

#### `ConvertURLToPDF`
Converts a web page to PDF.

```go
func (c *Client) ConvertURLToPDF(ctx context.Context, url string, opts ...ClientOptions) (*http.Response, error)
```

#### `ConvertHTMLToPDF`
Converts HTML content with optional assets to PDF.

```go
func (c *Client) ConvertHTMLToPDF(ctx context.Context, indexHTML io.Reader, opts ...ClientOptions) (*http.Response, error)
```

### Configuration Options

#### Builder Pattern

```go
options := gotenberg.NewOptionsBuilder().
    PaperSizeA4().                           // Set paper size
    Margins(1.0, 1.0, 1.0, 1.0).           // top, right, bottom, left (inches)
    PrintBackground(true).                   // Include CSS backgrounds
    Landscape(false).                        // Portrait orientation
    Scale(0.8).                             // Scale factor (0.1-2.0)
    SinglePage(false).                       // Generate single page
    OutputFilename("document.pdf").          // Custom filename
    File("style.css", cssReader).           // Add CSS file
    WebhookSuccess("https://...", "POST").   // Success webhook
    WebhookError("https://...", "POST").     // Error webhook
    WebhookExtraHeader("Auth", "Bearer ..."). // Custom headers
    Build()
```



### Paper Sizes

Pre-defined paper sizes available:
- `PaperSizeA4()`, `PaperSizeA3()`, `PaperSizeA5()`, etc.
- `PaperSizeLetter()`, `PaperSizeLegal()`, `PaperSizeTabloid()`
- `PaperSize(width, height)` for custom sizes

### Webhook Configuration

For async processing:

```go
options := gotenberg.NewOptionsBuilder().
    WebhookSuccess("https://your-api.com/webhook/success", "POST").
    WebhookError("https://your-api.com/webhook/error", "POST").
    WebhookExtraHeader("Authorization", "Bearer your-token").
    WebhookExtraHeader("X-Request-ID", "unique-id").
    Build()
```

## Examples

See the [examples](./example/) directory for complete working examples:

- [`html2pdf`](./example/cmd/html2pdf/) - HTML to PDF with CSS styling
- [`html2pdf_minimal`](./example/cmd/html2pdf_minimal/) - Minimal HTML conversion
- [`url2pdf`](./example/cmd/url2pdf/) - URL to PDF conversion
- [`html2pdf_webhook`](./example/cmd/html2pdf_webhook/) - Async processing with webhooks
- [`builder_demo`](./example/cmd/builder_demo/) - Builder pattern showcase

## Performance

The client is optimized for performance:

- **Buffer pooling** reduces memory allocations
- **Streaming responses** handle large PDFs efficiently  
- **Context support** enables proper timeout handling
- **5-15KB/op** memory usage (varies by complexity)
- **48-192 allocs/op** for complete request lifecycle

Benchmark results on modern hardware:
```
BenchmarkConvertURLToPDF-12      291750    3810 ns/op    5539 B/op    48 allocs/op
BenchmarkConvertHTMLToPDF-12     229570    4679 ns/op    8071 B/op    55 allocs/op
BenchmarkOptionsBuilder-12        91573   11741 ns/op   14810 B/op   192 allocs/op
```

## Error Handling

```go
resp, err := client.ConvertHTMLToPDF(ctx, htmlReader, options)
if err != nil {
    return fmt.Errorf("conversion failed: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
    body, _ := io.ReadAll(resp.Body)
    return fmt.Errorf("gotenberg error %d: %s", resp.StatusCode, body)
}

file, err := os.Create("output.pdf")
if err != nil {
    return err
}
defer file.Close()

_, err = io.Copy(file, resp.Body)
return err
```

## Best Practices

- ✅ **Always close `resp.Body`** to prevent resource leaks
- ✅ **Use `context.Context`** for timeouts and cancellations
- ✅ **Check response status codes** before processing
- ✅ **Use builder pattern** for clean, readable configurations
- ✅ **Inspect response headers** (e.g., `Gotenberg-Trace`) for debugging
- ✅ **Handle large files with streaming** instead of loading into memory
- ✅ **Set appropriate HTTP client timeouts** for your use case

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
