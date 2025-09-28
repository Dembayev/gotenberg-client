package gotenberg

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		wantError bool
	}{
		{
			name:      "valid URL",
			baseURL:   "http://localhost:3000",
			wantError: false,
		},
		{
			name:      "invalid URL",
			baseURL:   "://invalid-url",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(&http.Client{}, tt.baseURL)
			if (err != nil) != tt.wantError {
				t.Errorf("NewClient() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestClient_Methods(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	tests := []struct {
		name   string
		method func() *Client
		want   string
	}{
		{
			name:   "GET",
			method: func() *Client { return client.MethodGet(ctx, "/test") },
			want:   http.MethodGet,
		},
		{
			name:   "POST",
			method: func() *Client { return client.MethodPost(ctx, "/test") },
			want:   http.MethodPost,
		},
		{
			name:   "PUT",
			method: func() *Client { return client.MethodPut(ctx, "/test") },
			want:   http.MethodPut,
		},
		{
			name:   "PATCH",
			method: func() *Client { return client.MethodPatch(ctx, "/test") },
			want:   http.MethodPatch,
		},
		{
			name:   "DELETE",
			method: func() *Client { return client.MethodDelete(ctx, "/test") },
			want:   http.MethodDelete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()
			if result.err != nil {
				t.Errorf("Method() error = %v", result.err)
				return
			}
			if result.request.Method != tt.want {
				t.Errorf("Method() = %v, want %v", result.request.Method, tt.want)
			}
		})
	}
}

func TestClient_Header(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	result := client.MethodGet(ctx, "/test").
		Header("X-Test", "value").
		Header("Content-Type", "application/json")

	if result.err != nil {
		t.Fatalf("Header() error = %v", result.err)
	}

	if got := result.request.Header.Get("X-Test"); got != "value" {
		t.Errorf("Header X-Test = %v, want value", got)
	}

	if got := result.request.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("Header Content-Type = %v, want application/json", got)
	}
}

func TestClient_Headers(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	headers := map[string]string{
		"X-Test-1":     "value1",
		"X-Test-2":     "value2",
		"Content-Type": "application/json",
	}

	result := client.MethodGet(ctx, "/test").Headers(headers)

	if result.err != nil {
		t.Fatalf("Headers() error = %v", result.err)
	}

	for key, want := range headers {
		if got := result.request.Header.Get(key); got != want {
			t.Errorf("Header %s = %v, want %v", key, got, want)
		}
	}
}

func TestClient_QueryParam(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	result := client.MethodGet(ctx, "/test").
		QueryParam("param1", "value1").
		QueryParam("param2", "value2")

	if result.err != nil {
		t.Fatalf("QueryParam() error = %v", result.err)
	}

	query := result.request.URL.Query()
	if got := query.Get("param1"); got != "value1" {
		t.Errorf("QueryParam param1 = %v, want value1", got)
	}
	if got := query.Get("param2"); got != "value2" {
		t.Errorf("QueryParam param2 = %v, want value2", got)
	}
}

func TestClient_QueryParams(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	result := client.MethodGet(ctx, "/test").QueryParams(params)

	if result.err != nil {
		t.Fatalf("QueryParams() error = %v", result.err)
	}

	query := result.request.URL.Query()
	for key, want := range params {
		if got := query.Get(key); got != want {
			t.Errorf("QueryParam %s = %v, want %v", key, got, want)
		}
	}
}

func TestClient_QueryValues(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	values := url.Values{}
	values.Add("param1", "value1")
	values.Add("param2", "value2")
	values.Add("empty", "")

	result := client.MethodGet(ctx, "/test").QueryValues(values)

	if result.err != nil {
		t.Fatalf("QueryValues() error = %v", result.err)
	}

	query := result.request.URL.Query()
	if got := query.Get("param1"); got != "value1" {
		t.Errorf("QueryValues param1 = %v, want value1", got)
	}
	if got := query.Get("param2"); got != "value2" {
		t.Errorf("QueryValues param2 = %v, want value2", got)
	}
	if query.Has("empty") {
		t.Error("QueryValues should remove empty parameters")
	}
}

func TestClient_Body(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	body := io.NopCloser(strings.NewReader("test body"))
	result := client.MethodPost(ctx, "/test").Body(body)

	if result.err != nil {
		t.Fatalf("Body() error = %v", result.err)
	}

	if result.request.Body != body {
		t.Error("Body() did not set request body correctly")
	}
}

func TestClient_BytesBody(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	bodyBytes := []byte("test body")
	result := client.MethodPost(ctx, "/test").BytesBody(bodyBytes)

	if result.err != nil {
		t.Fatalf("BytesBody() error = %v", result.err)
	}

	if result.request.Body == nil {
		t.Fatal("BytesBody() did not set request body")
	}

	got, err := io.ReadAll(result.request.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	if !bytes.Equal(got, bodyBytes) {
		t.Errorf("BytesBody() = %v, want %v", got, bodyBytes)
	}
}

func TestClient_StringBody(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	bodyString := "test body"
	result := client.MethodPost(ctx, "/test").StringBody(bodyString)

	if result.err != nil {
		t.Fatalf("StringBody() error = %v", result.err)
	}

	got, err := io.ReadAll(result.request.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	if string(got) != bodyString {
		t.Errorf("StringBody() = %v, want %v", string(got), bodyString)
	}
}

func TestClient_JSONBody(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	data := map[string]string{"key": "value"}
	result := client.MethodPost(ctx, "/test").JSONBody(data)

	if result.err != nil {
		t.Fatalf("JSONBody() error = %v", result.err)
	}

	contentType := result.request.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("JSONBody() Content-Type = %v, want application/json", contentType)
	}

	got, err := io.ReadAll(result.request.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	expected := `{"key":"value"}`
	if string(got) != expected {
		t.Errorf("JSONBody() = %v, want %v", string(got), expected)
	}
}

func TestClient_Multipart(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	result := client.MethodPost(ctx, "/test").Multipart()

	if result.err != nil {
		t.Fatalf("Multipart() error = %v", result.err)
	}

	if !result.multipart {
		t.Error("Multipart() did not set multipart flag")
	}

	if result.writer == nil {
		t.Error("Multipart() did not create writer")
	}

	if result.buffer == nil {
		t.Error("Multipart() did not create buffer")
	}
}

func TestClient_FormField(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	result := client.MethodPost(ctx, "/test").
		FormField("field1", "value1").
		FormField("field2", "value2")

	if result.err != nil {
		t.Fatalf("FormField() error = %v", result.err)
	}

	if !result.multipart {
		t.Error("FormField() should enable multipart")
	}
}

func TestClient_File(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	fileContent := strings.NewReader("file content")
	result := client.MethodPost(ctx, "/test").
		File("files", "test.txt", fileContent)

	if result.err != nil {
		t.Fatalf("File() error = %v", result.err)
	}

	if !result.multipart {
		t.Error("File() should enable multipart")
	}
}

func TestClient_Reset(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	// Setup client with some state
	result := client.MethodPost(ctx, "/test").
		Header("X-Test", "value").
		FormField("field", "value")

	if result.err != nil {
		t.Fatalf("Setup error = %v", result.err)
	}

	// Reset client
	reset := result.Reset()

	if reset.request != nil {
		t.Error("Reset() should clear request")
	}

	if reset.err != nil {
		t.Error("Reset() should clear error")
	}

	if reset.multipart {
		t.Error("Reset() should clear multipart flag")
	}

	if reset.writer != nil {
		t.Error("Reset() should clear writer")
	}
}

func TestClient_Send(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := context.Background()
	resp, err := client.MethodGet(ctx, "/test").Send()

	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Send() status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != "OK" {
		t.Errorf("Send() body = %v, want OK", string(body))
	}
}

func TestClient_SendWithMultipart(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Error("Expected multipart/form-data content type")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(&http.Client{Timeout: 5 * time.Second}, server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := context.Background()
	resp, err := client.MethodPost(ctx, "/test").
		FormField("field", "value").
		Send()

	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Send() status = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	// Test error propagation
	result := client.MethodGet(ctx, "/test")
	result.err = http.ErrMissingFile // Set an error

	// All subsequent calls should return the same client with error
	result2 := result.Header("X-Test", "value")
	if result2.err != http.ErrMissingFile {
		t.Error("Error should be propagated through method calls")
	}

	result3 := result2.QueryParam("param", "value")
	if result3.err != http.ErrMissingFile {
		t.Error("Error should be propagated through method calls")
	}

	// Test Err() method
	if result3.Err() != http.ErrMissingFile {
		t.Error("Err() should return the stored error")
	}
}

func TestClient_JSONBody_MarshalError(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	// Create a type that can't be marshaled to JSON
	invalidData := make(chan int)

	result := client.MethodPost(ctx, "/test").JSONBody(invalidData)

	if result.err == nil {
		t.Error("JSONBody() should return error for invalid data")
	}
}

func TestClient_Send_WithError(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	// Set an error
	result := client.MethodGet(ctx, "/test")
	result.err = http.ErrMissingFile

	resp, err := result.Send()

	if err != http.ErrMissingFile {
		t.Errorf("Send() should return stored error, got %v", err)
	}

	if resp != nil {
		t.Error("Send() should return nil response when error exists")
	}
}

func TestWebhookExtraHeaders_MarshalError(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	// This won't actually cause a marshal error since map[string]string is always valid JSON
	// But we can test the normal flow
	headers := map[string]string{
		"X-Test": "value",
	}

	result := client.MethodPost(ctx, "/test").WebhookExtraHeaders(headers)

	if result.err != nil {
		t.Errorf("WebhookExtraHeaders() should not error with valid headers: %v", result.err)
	}
}

// Helper function to create test client
func createTestClient(t *testing.T) *Client {
	client, err := NewClient(&http.Client{}, "http://localhost:3000")
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	return client
}
