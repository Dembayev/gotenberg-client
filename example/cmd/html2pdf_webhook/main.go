package main

import (
	"bytes"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nativebpm/gotenberg-client"
	"github.com/nativebpm/gotenberg-client/example/model"
	"github.com/nativebpm/gotenberg-client/example/pkg/image"
	"github.com/nativebpm/gotenberg-client/example/pkg/templates/invoice"
)

func main() {
	srv := StartServer(":28080")

	gotenbergURL := `http://localhost:3000`

	httpClient := &http.Client{
		Timeout: 90 * time.Second,
	}

	clientBuilder := gotenberg.NewClientBuilder(httpClient, gotenbergURL)

	data := model.InvoiceData
	html := bytes.NewBuffer(nil)
	invoice.Template.Execute(html, data)

	logo := image.LogoPNG()

	// Using builder pattern with webhook configuration
	resp, err := clientBuilder.ConvertHTML().
		WithFile("logo.png", bytes.NewReader(logo)).
		PrintBackground(true).
		OutputFilename("invoice_async.pdf").
		WebhookSuccess("http://host.docker.internal:28080/success", "POST").
		WebhookError("http://host.docker.internal:28080/error", "POST").
		WebhookExtraHeader("Authorization", "Bearer your-token").
		WebhookExtraHeader("X-Custom-Header", "custom-value").
		Execute(context.Background(), html)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	slog.Info("Async HTML to PDF conversion started",
		"trace", resp.Header.Get("Gotenberg-Trace"))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	slog.Info("waiting for webhook callbacks; press Ctrl+C to exit")
	<-sigCh

	slog.Info("shutting down webhook server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Info("error shutting down server", "err", err)
	}
}
