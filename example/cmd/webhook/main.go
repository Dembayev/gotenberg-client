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

	client := gotenberg.NewClient(httpClient, gotenbergURL)

	data := model.InvoiceData
	html := bytes.NewBuffer(nil)
	invoice.Template.Execute(html, data)

	logo := image.LogoPNG()

	resp, err := client.IndexHTML(html).
		File("logo.png", bytes.NewReader(logo)).
		Bool(gotenberg.FieldPrintBackground, true).
		ConvertHTML(context.Background())

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// if err := client.SetWebhookSuccess("http://host.docker.internal:28080/success", "POST"); err != nil {
	// 	log.Fatal("Error setting webhook success:", err)
	// }

	// if err := client.SetWebhookError("http://host.docker.internal:28080/error", "POST"); err != nil {
	// 	log.Fatal("Error setting webhook error:", err)
	// }

	// webhookHeaders := map[string]string{
	// 	"Authorization":   "Bearer your-token",
	// 	"X-Custom-Header": "custom-value",
	// }

	// if err := client.SetWebhookHeaders(webhookHeaders); err != nil {
	// 	log.Fatal("Error setting webhook headers:", err)
	// }

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
