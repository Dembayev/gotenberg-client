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

// BenchmarkConcurrentSafetyBasic тестирует базовую потокобезопасность
func BenchmarkConcurrentSafetyBasic(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient(&http.Client{}, server.URL)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.MethodGet(ctx, "/test").Send()
			if err != nil {
				b.Error(err)
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
		}
	})
}

// BenchmarkConcurrentHTMLConversion тестирует concurrent HTML конверсии
func BenchmarkConcurrentHTMLConversion(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("mock pdf content"))
	}))
	defer server.Close()

	client, _ := NewClient(&http.Client{}, server.URL)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			html := strings.NewReader("<html><body>Test Content</body></html>")
			resp, err := client.ConvertHTML(ctx, html).PaperSizeA4().Margins(1.0, 1.0, 1.0, 1.0).Send()
			if err != nil {
				b.Error(err)
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
		}
	})
}

// BenchmarkConcurrentURLConversion тестирует concurrent URL конверсии
func BenchmarkConcurrentURLConversion(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("mock pdf content"))
	}))
	defer server.Close()

	client, _ := NewClient(&http.Client{}, server.URL)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.ConvertURL(ctx, "https://example.com").PaperSizeA4().Send()
			if err != nil {
				b.Error(err)
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
		}
	})
}

// BenchmarkConcurrentRequestIsolation тестирует изоляцию между запросами
func BenchmarkConcurrentRequestIsolation(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient(&http.Client{}, server.URL)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		requestID := 0
		for pb.Next() {
			// Каждый запрос имеет уникальные параметры
			resp, err := client.MethodPost(ctx, "/test").
				Header("X-Request-ID", strconv.Itoa(requestID)).
				QueryParam("id", strconv.Itoa(requestID)).
				StringBody("Request body " + strconv.Itoa(requestID)).
				Send()

			if err != nil {
				b.Error(err)
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
			requestID++
		}
	})
}

// BenchmarkConcurrentStressTest стресс-тест для высокой нагрузки
func BenchmarkConcurrentStressTest(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Имитация более реальной обработки
		time.Sleep(2 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient(&http.Client{}, server.URL)
	ctx := context.Background()

	const goroutines = 50
	var wg sync.WaitGroup

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(goroutines)

		for j := 0; j < goroutines; j++ {
			go func(requestID int) {
				defer wg.Done()

				resp, err := client.MethodGet(ctx, "/stress-test").
					Header("X-Request-ID", strconv.Itoa(requestID)).
					QueryParam("worker", strconv.Itoa(requestID%10)).
					Send()

				if err != nil {
					b.Error(err)
				}
				if resp != nil {
					_ = resp.Body.Close()
				}
			}(j)
		}

		wg.Wait()
	}
}
