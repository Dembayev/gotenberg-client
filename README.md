# gotenberg-client

[![Go Reference](https://pkg.go.dev/badge/github.com/nativebpm/gotenberg-client.svg)](https://pkg.go.dev/github.com/nativebpm/gotenberg-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/nativebpm/gotenberg-client)](https://goreportcard.com/report/github.com/nativebpm/gotenberg-client)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)


A high-performance Go client for the [Gotenberg](https://gotenberg.dev/) HTTP API with a fluent interface. Built using only the Go standard library (via [http-client](https://github.com/nativebpm/http-client)).

**Features:**
- Minimal dependencies (only stdlib + [http-client](https://github.com/nativebpm/http-client))
- Fluent API for building requests
- Support for HTML/URL to PDF conversion
- Webhook support for async jobs (with example)
- Easy file/multipart uploads
- Paper size, margins, and advanced PDF options


## Quick Start: Synchronous HTML to PDF


```go
package main

import (
    "context"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
    "github.com/nativebpm/gotenberg-client"
)

func main() {
    client, err := gotenberg.NewClient(&http.Client{}, "http://localhost:3000")
    if err != nil {
        log.Fatal(err)
    }

    html := strings.NewReader("<html><body><h1>Hello World!</h1></body></html>")

    resp, err := client.
        ConvertHTML(context.Background(), html).
        PaperSizeA4().
        Send()
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    f, _ := os.Create("out.pdf")
    io.Copy(f, resp.Body)
    f.Close()
}
```

## Example: Async PDF Generation with Webhook

See [`example/cmd/webhook`](example/cmd/webhook) for a full async webhook demo (HTML invoice to PDF, with logo, webhook server, and error handling):

```sh
go run ./example/cmd/webhook
```

This will:
- Start a local webhook server
- Generate an invoice PDF using HTML template and logo
- Receive the PDF via webhook callback from Gotenberg


## Installation


```bash
go get github.com/nativebpm/gotenberg-client
```


## Testing

Run all tests and benchmarks:

```sh
go test -v -bench=. ./...
```

## Project Structure

- `gotenberg.go` — main client implementation
- `example/` — real-world usage: invoice template, logo, webhook server
- `example/cmd/webhook` — async webhook demo
- `example/model` — invoice data structs
- `example/pkg/templates/invoice` — HTML template for invoice
- `example/pkg/image` — logo generator

## Dependencies

- Go standard library
- [`github.com/nativebpm/http-client`](https://github.com/nativebpm/http-client)

No third-party dependencies are required, ensuring minimal bloat and maximum compatibility.


## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
