# gotenberg — Go client for Gotenberg

A clean, dependency-free and idiomatic Go client for the Gotenberg HTTP API.

- Request bodies are constructed with `mime/multipart`.
- Adds optional webhook headers (`Gotenberg-Webhook-*`) and `Gotenberg-Output-Filename` when provided via options.
- Returning `*http.Response` allows streaming large PDFs without buffering them entirely in memory.
- No third-party dependencies: the client uses only the Go standard library (net/http, mime/multipart, context, etc.).

Usage example

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

Recommendations

- Always close `resp.Body`.
- Use `context.Context` for timeouts and cancellations.
- Inspect response headers (e.g. `Gotenberg-Trace`) when needed.
