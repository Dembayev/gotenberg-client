package main

import (
	"bytes"
	"context"
	"log"
	"log/slog"
	"net/http"
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

	client := gotenberg.NewClient(httpClient, gotenbergURL)

	data := model.InvoiceData
	html := bytes.NewBuffer(nil)
	invoice.Template.Execute(html, data)

	logo := image.LogoPNG()

	resp, err := client.ConvertHTMLToPDF(context.Background(), html,
		gotenberg.WithPrintBackground(true),
		gotenberg.WithOutputFilename("invoice_async.pdf"),
		gotenberg.WithWebhookSuccess(
			"https://your-webhook-url.com/success",
			"POST",
		),
		gotenberg.WithWebhookError(
			"https://your-webhook-url.com/error",
			"POST",
		),
		gotenberg.WithWebhookExtraHeader("Authorization", "Bearer your-token"),
		gotenberg.WithWebhookExtraHeader("X-Custom-Header", "custom-value"),
		gotenberg.WithFile("logo.png", bytes.NewReader(logo)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	slog.Info("Async HTML to PDF conversion started",
		"trace", resp.Header.Get("Gotenberg-Trace"))
}
