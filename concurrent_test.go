package gotenberg

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestClient_ConcurrentSafety(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // Simulate some processing time
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		t.Fatal(err)
	}

	const numGoroutines = 100
	const requestsPerGoroutine = 10

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*requestsPerGoroutine)

	ctx := context.Background()

	// Launch multiple goroutines making concurrent requests
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < requestsPerGoroutine; j++ {
				// Each goroutine creates its own request chain
				resp, err := client.MethodGet(ctx, "/test").
					Header("X-Goroutine-ID", strconv.Itoa(goroutineID)).
					QueryParam("request", strconv.Itoa(j)).
					Send()

				if err != nil {
					errors <- err
					continue
				}

				_ = resp.Body.Close()
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent request failed: %v", err)
	}
}

func TestClient_ConcurrentHTMLConversion(t *testing.T) {
	// Create test server that handles multipart requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Error("Expected multipart/form-data content type")
		}

		// Simulate processing time
		time.Sleep(5 * time.Millisecond)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("PDF content"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		t.Fatal(err)
	}

	const numGoroutines = 50
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	ctx := context.Background()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			html := strings.NewReader("<html><body><h1>Test PDF " + string(rune(id)) + "</h1></body></html>")

			resp, err := client.ConvertHTML(ctx, html).
				Bool(FieldPrintBackground, true).
				PaperSizeA4().
				Margins(1.0, 1.0, 1.0, 1.0).
				Send()

			if err != nil {
				errors <- err
				return
			}

			_ = resp.Body.Close()
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent HTML conversion failed: %v", err)
	}
}

func TestClient_ConcurrentURLConversion(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("PDF content"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		t.Fatal(err)
	}

	const numGoroutines = 30
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	ctx := context.Background()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := client.ConvertURL(ctx, "https://example.com").
				Bool(FieldLandscape, true).
				PaperSizeLetter().
				Send()

			if err != nil {
				errors <- err
				return
			}

			_ = resp.Body.Close()
		}()
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent URL conversion failed: %v", err)
	}
}

func TestClient_ConcurrentWithWebhooks(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for webhook headers
		webhookURL := r.Header.Get("Gotenberg-Webhook-Url")
		if webhookURL == "" {
			t.Error("Expected webhook URL header")
		}

		time.Sleep(5 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("PDF content"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		t.Fatal(err)
	}

	const numGoroutines = 25
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	ctx := context.Background()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			html := strings.NewReader("<html><body>Test " + string(rune(id)) + "</body></html>")

			resp, err := client.ConvertHTML(ctx, html).
				WebhookURLMethodPost("http://localhost:8080/success").
				WebhookErrorURLMethodPost("http://localhost:8080/error").
				WebhookExtraHeaders(map[string]string{
					"X-Request-ID": "req-" + string(rune(id)),
				}).
				Send()

			if err != nil {
				errors <- err
				return
			}

			_ = resp.Body.Close()
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent webhook request failed: %v", err)
	}
}
