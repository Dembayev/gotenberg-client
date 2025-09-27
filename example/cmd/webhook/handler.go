package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func webhookHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		outFile, err := os.Create("./invoice_async.pdf")
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
			"x-custom-header", r.Header.Get("X-Custom-Header"),
			"body lenth", n,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}
}
