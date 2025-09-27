package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/nativebpm/gotenberg-client"
)

func main() {
	gotenbergURL := `http://localhost:3000`

	httpClient := &http.Client{
		Timeout: 90 * time.Second,
	}

	clientBuilder := gotenberg.NewClientBuilder(httpClient, gotenbergURL)

	// Using builder pattern for URL to PDF conversion
	resp, err := clientBuilder.ConvertURL().
		PrintBackground(true).
		OutputFilename("example.pdf").
		PaperSizeLetter().
		Execute(context.Background(), "https://example.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	outFile, err := os.Create("./example.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	n, err := io.Copy(outFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("URL to PDF conversion completed",
		"pdf_size_bytes", n,
		"trace", resp.Header.Get("Gotenberg-Trace"))
}
