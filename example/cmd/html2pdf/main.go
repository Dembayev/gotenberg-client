package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/nativebpm/gotenberg-client"
	"github.com/nativebpm/gotenberg-client/example/model"
	"github.com/nativebpm/gotenberg-client/example/pkg/image"
	"github.com/nativebpm/gotenberg-client/example/pkg/templates/invoice"
)

func main() {
	gotenbergURL := `http://localhost:3000`

	httpClient := &http.Client{
		Timeout: 90 * time.Second,
	}

	clientBuilder := gotenberg.NewClientBuilder(httpClient, gotenbergURL)
	data := model.InvoiceData

	html := bytes.NewBuffer(nil)
	invoice.Template.Execute(html, data)

	logo := image.LogoPNG()

	resp, err := clientBuilder.ConvertHTML().
		WithFile("logo.png", bytes.NewReader(logo)).
		PrintBackground(true).
		Landscape(false).
		Scale(1.0).
		OutputFilename("invoice.pdf").
		PaperSizeA4().
		Margins(1.0, 1.0, 1.0, 1.0).
		Execute(context.Background(), html)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	outFile, err := os.Create("./invoice_new.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	n, err := io.Copy(outFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("HTML to PDF conversion completed",
		"pdf_size_bytes", n,
		"content_type", resp.Header.Get("Content-Type"),
		"trace", resp.Header.Get("Gotenberg-Trace"))
}
