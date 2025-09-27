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

	client := gotenberg.NewClient(httpClient, gotenbergURL)

	if err := client.WriteURL("https://example.com"); err != nil {
		log.Fatal("Error writing URL:", err)
	}

	if err := client.WritePaperSize(gotenberg.PaperSizeLetter[0], gotenberg.PaperSizeLetter[1]); err != nil {
		log.Fatal("Error setting paper size:", err)
	}

	if err := client.WriteBoolProperty(gotenberg.FieldPrintBackground, true); err != nil {
		log.Fatal("Error setting print background:", err)
	}

	resp, err := client.ExecuteURL(context.Background())
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
