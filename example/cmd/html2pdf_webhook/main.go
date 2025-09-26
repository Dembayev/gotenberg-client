package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
	gotenbergURL := `http://localhost:3000`

	httpClient := &http.Client{
		Timeout: 90 * time.Second,
	}

	client := gotenberg.NewClient(httpClient, gotenbergURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/success", webhookHandler("success"))
	mux.HandleFunc("/error", webhookHandler("error"))

	srv := &http.Server{Addr: ":28080", Handler: mux}
	go func() {
		slog.Info("starting local webhook server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("webhook server error: %v", err)
		}
	}()

	data := model.InvoiceData
	html := bytes.NewBuffer(nil)
	invoice.Template.Execute(html, data)

	logo := image.LogoPNG()

	resp, err := client.ConvertHTMLToPDF(context.Background(), html,
		gotenberg.WithFile("logo.png", bytes.NewReader(logo)),
		gotenberg.WithPrintBackground(true),
		gotenberg.WithOutputFilename("invoice_async.pdf"),
		gotenberg.WithWebhookSuccess(
			"http://host.docker.internal:28080/success",
			"POST",
		),
		gotenberg.WithWebhookError(
			"http://host.docker.internal:28080/error",
			"POST",
		),
		gotenberg.WithWebhookExtraHeader("Authorization", "Bearer your-token"),
		gotenberg.WithWebhookExtraHeader("X-Custom-Header", "custom-value"),
	)
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

func webhookHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		outFile, err := os.Create("invoice_async.pdf")
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()

		n, err := io.Copy(outFile, r.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Body.Close()

		slog.Info(fmt.Sprintf("webhook %s received", name),
			"method", r.Method,
			"path", r.URL.Path,
			"gotenberg-trace", r.Header.Get("Gotenberg-Trace"),
			"authorization", r.Header.Get("Authorization"),
			"x-custom-header", r.Header.Get("X-Custom-Header"),
			"body lenth", n,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}
}
