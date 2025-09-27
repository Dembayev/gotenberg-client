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
		WithHTMLReader(html).
		WithFile("logo.png", bytes.NewReader(logo)).
		Execute(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	outFile, err := os.Create("./invoice_minimal.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	n, err := io.Copy(outFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Minimal HTML to PDF conversion completed",
		"pdf_size_bytes", n,
		"trace", resp.Header.Get("Gotenberg-Trace"))
}
