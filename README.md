# gotenberg — Go client for Gotenberg

A clean, dependency-free and idiomatic Go client for the Gotenberg HTTP API.

- Request bodies are constructed with `mime/multipart`.
- Adds optional webhook headers (`Gotenberg-Webhook-*`) and `Gotenberg-Output-Filename` when provided via options.
- Returning `*http.Response` allows streaming large PDFs without buffering them entirely in memory.
- No third-party dependencies: the client uses only the Go standard library (net/http, mime/multipart, context, etc.).

## Usage Examples

### Simple URL to PDF conversion

```go
ctx := context.Background()
httpClient := &http.Client{}
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

### Using Traditional Functional Options

```go
resp, err := cli.ConvertHTMLToPDF(ctx, htmlReader,
    gotenberg.WithPaperSizeA4(),
    gotenberg.WithMargins(1.0, 1.0, 1.0, 1.0),
    gotenberg.WithPrintBackground(true),
    gotenberg.WithOutputFilename("document.pdf"),
)
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

Recommendations

- Always close `resp.Body`.
- Use `context.Context` for timeouts and cancellations.
- Inspect response headers (e.g. `Gotenberg-Trace`) when needed.
