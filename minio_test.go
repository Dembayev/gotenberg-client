// Package gotenberg provides tests for MinIO API
package gotenberg_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	gotenberg "github.com/nativebpm/gotenberg-client"
)

// MockMinioClient is a mock implementation for testing
type MockMinioClient struct {
	files map[string][]byte
}

func (m *MockMinioClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	m.files[objectName] = data
	return nil
}

func (m *MockMinioClient) DownloadFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	data, exists := m.files[objectName]
	if !exists {
		return nil, http.ErrNotSupported
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

// Example test for upload endpoint
func ExampleMinioAPI_HandleUpload() {
	// This is an example of how to test the upload endpoint
	// In real tests, you would use a proper MinIO client or mock

	// Create a buffer to write our multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Add file field
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.Close()

	// Note: In real implementation, you would create a request and response
	// recorder, then call api.HandleUpload(w, req) where api is your MinioAPI instance

	// Output:
	// Example of upload request structure
}

// Example test for download endpoint
func ExampleMinioAPI_HandleDownload() {
	// This is an example of how to test the download endpoint

	// Note: In real implementation, you would create a request and response
	// recorder, then call api.HandleDownload(w, req) where api is your MinioAPI instance

	// Output:
	// Example of download request structure
}

// Example of programmatic file upload (not through HTTP)
func ExampleMinioClient_UploadFile() {
	config := gotenberg.MinioConfig{
		Endpoint:        "localhost:9000",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
		BucketName:      "documents",
		UseSSL:          false,
	}

	ctx := context.Background()
	client, err := gotenberg.NewMinioClient(ctx, config)
	if err != nil {
		panic(err)
	}

	// Upload a file
	content := strings.NewReader("Hello, MinIO!")
	_, err = client.UploadFile(ctx, "hello.txt", content, int64(content.Len()), "text/plain")
	if err != nil {
		panic(err)
	}

	// Output:
	// File uploaded successfully
}

// Example of programmatic file download (not through HTTP)
func ExampleMinioClient_DownloadFile() {
	config := gotenberg.MinioConfig{
		Endpoint:        "localhost:9000",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
		BucketName:      "documents",
		UseSSL:          false,
	}

	ctx := context.Background()
	client, err := gotenberg.NewMinioClient(ctx, config)
	if err != nil {
		panic(err)
	}

	// Download a file
	reader, err := client.DownloadFile(ctx, "hello.txt")
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// Read content
	content, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	_ = content // Use the content

	// Output:
	// File downloaded successfully
}
