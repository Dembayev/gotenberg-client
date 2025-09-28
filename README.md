# gotenberg-client

[![Go Reference](https://pkg.go.dev/badge/github.com/nativebpm/gotenberg-client.svg)](https://pkg.go.dev/github.com/nativebpm/gotenberg-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/nativebpm/gotenberg-client)](https://goreportcard.com/report/github.com/nativebpm/gotenberg-client)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A high-performance Go client for the [Gotenberg](https://gotenberg.dev/) HTTP API with fluent interface and zero dependencies.

## Quick Start

```go
package main

import (
    "context"
    "log"
    "net/http"
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
}
```

## Installation

```bash
go get github.com/nativebpm/gotenberg-client
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
