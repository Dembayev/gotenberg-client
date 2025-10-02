# gotenberg-client

[![Go Reference](https://pkg.go.dev/badge/github.com/nativebpm/gotenberg-client.svg)](https://pkg.go.dev/github.com/nativebpm/gotenberg-client)
[![Go Report Card](https://goreportcard.com/badge/github.com/nativebpm/gotenberg-client)](https://goreportcard.com/report/github.com/nativebpm/gotenberg-client)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)


A high-performance Go client for the [Gotenberg](https://gotenberg.dev/) HTTP API with a fluent interface and MinIO storage integration. Built using only the Go standard library (via [http-client](https://github.com/nativebpm/http-client)).

**Features:**
- Minimal dependencies (only stdlib + [http-client](https://github.com/nativebpm/http-client))
- Fluent API for building requests
- Support for HTML/URL to PDF conversion
- Webhook support
- Easy file/multipart uploads
- Paper size, margins, and advanced PDF options
- **MinIO storage integration with HTTP APIs**
- **File upload/download REST endpoints**


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
		PaperSizeA6().
		Margins(0.5, 0.5, 0.5, 0.5).
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

## MinIO Storage Integration

This package now includes a complete MinIO storage solution with HTTP APIs for file upload and download operations.

### MinIO API Endpoints

The package provides two REST API endpoints:

1. **POST /api/upload** - Upload files to MinIO
2. **GET /api/download** - Download files from MinIO

See [MINIO_API.md](MINIO_API.md) for detailed documentation and [API_EXAMPLES.md](API_EXAMPLES.md) for code examples in multiple languages.

### Quick Start with MinIO

```go
package main

import (
	"context"
	"log"
	"net/http"

	gotenberg "github.com/nativebpm/gotenberg-client"
)

func main() {
	// Configure MinIO
	config := gotenberg.MinioConfig{
		Endpoint:        "localhost:9000",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
		BucketName:      "documents",
		UseSSL:          false,
	}

	// Create MinIO client
	ctx := context.Background()
	minioClient, err := gotenberg.NewMinioClient(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	// Create API handler
	api := gotenberg.NewMinioAPI(minioClient)

	// Setup HTTP server
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	// Start server
	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", mux)
}
```

### Running with Docker Compose

Start MinIO and Gotenberg services:

```bash
docker-compose up -d
```

This starts:
- MinIO on ports 9000 (API) and 9001 (Console)
- Gotenberg on port 3000

Access MinIO Console at: http://localhost:9001

### Using Make Commands

```bash
# Start MinIO
make minio-run

# Install dependencies
make api-deps

# Run API server
make api-run

# Test API endpoints
make api-test

# Start both MinIO and API server
make dev

# Clean up containers
make clean
```

## Example: Async PDF Generation with Webhook

See [`examples/cmd/webhook`](examples/cmd/webhook) for a full async webhook demo (HTML invoice to PDF, with logo, webhook server, and error handling):

```sh
go run ./examples/cmd/webhook
```

This will:
- Start a local webhook server
- Generate an invoice PDF using HTML template and logo
- Receive the PDF via webhook callback from Gotenberg


## Installation


```bash
go get github.com/nativebpm/gotenberg-client
```

For MinIO support, install additional dependency:

```bash
go get github.com/minio/minio-go/v7
```

## Testing

Run all tests and benchmarks:

```sh
go test -v -bench=. ./...
```

## Project Structure

- `gotenberg.go` — main client implementation
- `minio.go` — MinIO client implementation
- `minio_api.go` — HTTP API handlers for MinIO operations
- `examples/` — real-world usage: invoice template, logo, webhook server
- `examples/cmd/webhook` — async webhook demo
- `examples/minio_api_server.go` — MinIO API server example
- `examples/model` — invoice data structs
- `examples/pkg/templates/invoice` — HTML template for invoice
- `examples/pkg/image` — logo generator
- `MINIO_API.md` — MinIO API documentation
- `API_EXAMPLES.md` — API usage examples in multiple languages

## API Examples

The package includes comprehensive examples for:
- cURL commands
- JavaScript/TypeScript (Fetch API, Axios)
- Python (requests library)
- Go (standard library)
- React components
- PHP

See [API_EXAMPLES.md](API_EXAMPLES.md) for details.

## Dependencies

- Go standard library
- [`github.com/nativebpm/http-client`](https://github.com/nativebpm/http-client)
- [`github.com/minio/minio-go/v7`](https://github.com/minio/minio-go) (for MinIO features)

No other third-party dependencies are required, ensuring minimal bloat and maximum compatibility.


## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
