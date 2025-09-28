package gotenberg

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func BenchmarkNewClient(b *testing.B) {
	httpClient := &http.Client{}
	baseURL := "http://localhost:3000"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewClient(httpClient, baseURL)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkClient_MethodGet(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodGet(ctx, "/test")
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_MethodPost(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test")
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_Header(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodGet(ctx, "/test").Header("X-Test", "value")
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_Headers(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	headers := map[string]string{
		"X-Test-1":     "value1",
		"X-Test-2":     "value2",
		"Content-Type": "application/json",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodGet(ctx, "/test").Headers(headers)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_QueryParam(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodGet(ctx, "/test").QueryParam("param", "value")
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_QueryParams(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
		"param3": "value3",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodGet(ctx, "/test").QueryParams(params)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_JSONBody(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").JSONBody(data)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_BytesBody(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	body := bytes.Repeat([]byte("test"), 1000) // 4KB body

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").BytesBody(body)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_StringBody(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	body := strings.Repeat("test", 1000) // 4KB body

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").StringBody(body)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_Multipart(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").Multipart()
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_FormField(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().MethodPost(ctx, "/test").FormField("field", "value")
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_File(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	content := strings.NewReader("file content")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		content.Reset("file content")
		result := client.Reset().MethodPost(ctx, "/test").File("files", "test.txt", content)
		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_Reset(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	// Pre-setup client
	client.MethodPost(ctx, "/test").
		Header("X-Test", "value").
		FormField("field", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Reset()
	}
}

func BenchmarkClient_ComplexChain(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	headers := map[string]string{
		"X-Custom": "value",
		"X-Auth":   "token",
	}
	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().
			MethodPost(ctx, "/test").
			Headers(headers).
			QueryParams(params).
			FormField("field1", "value1").
			FormField("field2", "value2")

		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_Send_GET(b *testing.B) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.Reset().MethodGet(ctx, "/test").Send()
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkClient_Send_POST_JSON(b *testing.B) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.Reset().MethodPost(ctx, "/test").JSONBody(data).Send()
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkClient_Send_POST_Multipart(b *testing.B) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	fileContent := strings.NewReader("file content data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileContent.Reset("file content data")
		resp, err := client.Reset().
			MethodPost(ctx, "/test").
			FormField("field", "value").
			File("files", "test.txt", fileContent).
			Send()
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// Memory allocation benchmarks
func BenchmarkClient_MemoryAllocation_SimpleChain(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := client.Reset().
			MethodPost(ctx, "/test").
			Header("X-Test", "value").
			QueryParam("param", "value")

		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

func BenchmarkClient_MemoryAllocation_MultipartForm(b *testing.B) {
	client := createBenchClient(b)
	ctx := context.Background()
	content := strings.NewReader("file content")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		content.Reset("file content")
		result := client.Reset().
			MethodPost(ctx, "/test").
			FormField("field1", "value1").
			FormField("field2", "value2").
			File("files", "test.txt", content)

		if result.err != nil {
			b.Fatal(result.err)
		}
	}
}

// Helper function to create benchmark client
func createBenchClient(b *testing.B) *Client {
	client, err := NewClient(&http.Client{}, "http://localhost:3000")
	if err != nil {
		b.Fatal(err)
	}
	return client
}
