// Example of using MinIO API endpoints
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	gotenberg "github.com/nativebpm/gotenberg-client"
)

func main() {
	// MinIO configuration
	config := gotenberg.MinioConfig{
		Endpoint:        "localhost:9000",           // MinIO server endpoint
		AccessKeyID:     "minioadmin",               // MinIO access key
		SecretAccessKey: "minioadmin",               // MinIO secret key
		BucketName:      "documents",                // Bucket name
		UseSSL:          false,                      // Use HTTPS (false for local development)
	}

	// Create MinIO client
	ctx := context.Background()
	minioClient, err := gotenberg.NewMinioClient(ctx, config)
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	// Create API handler
	api := gotenberg.NewMinioAPI(minioClient)

	// Create HTTP server
	mux := http.NewServeMux()
	
	// Register MinIO routes
	api.RegisterRoutes(mux)

	// Add a simple health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Println("Available endpoints:")
	fmt.Println("  POST /api/upload     - Upload file to MinIO")
	fmt.Println("  GET  /api/download   - Download file from MinIO")
	fmt.Println("  GET  /health         - Health check")
	
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
