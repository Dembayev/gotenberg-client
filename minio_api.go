// Package gotenberg provides HTTP API handlers for MinIO operations
package gotenberg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

// MinioAPI provides HTTP handlers for MinIO operations
type MinioAPI struct {
	minioClient *MinioClient
}

// NewMinioAPI creates a new MinIO API handler
func NewMinioAPI(minioClient *MinioClient) *MinioAPI {
	return &MinioAPI{
		minioClient: minioClient,
	}
}

// UploadRequest represents the upload response
type UploadResponse struct {
	Success    bool   `json:"success"`
	ObjectName string `json:"object_name"`
	Size       int64  `json:"size"`
	ETag       string `json:"etag"`
	Message    string `json:"message,omitempty"`
}

// DownloadErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// HandleUpload handles file upload to MinIO
// POST /api/upload
// Expects multipart/form-data with a file field named "file"
// Optional query parameter: objectName (if not provided, uses the original filename)
func (api *MinioAPI) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse multipart form (max 100MB in memory)
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to parse multipart form: "+err.Error())
		return
	}

	// Get file from request
	file, header, err := r.FormFile("file")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to get file from request: "+err.Error())
		return
	}
	defer file.Close()

	// Get object name from query parameter or use original filename
	objectName := r.URL.Query().Get("objectName")
	if objectName == "" {
		objectName = header.Filename
	}

	// Get content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to MinIO
	ctx := r.Context()
	uploadInfo, err := api.minioClient.UploadFile(ctx, objectName, file, header.Size, contentType)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to upload file: "+err.Error())
		return
	}

	// Send success response
	response := UploadResponse{
		Success:    true,
		ObjectName: uploadInfo.Key,
		Size:       uploadInfo.Size,
		ETag:       uploadInfo.ETag,
		Message:    "File uploaded successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleDownload handles file download from MinIO
// GET /api/download?objectName=filename.pdf
// Query parameter: objectName (required) - the name of the file in MinIO
func (api *MinioAPI) HandleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get object name from query parameter
	objectName := r.URL.Query().Get("objectName")
	if objectName == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Missing objectName parameter")
		return
	}

	ctx := r.Context()

	// Get file info first to set proper headers
	fileInfo, err := api.minioClient.GetFileInfo(ctx, objectName)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "File not found: "+err.Error())
		return
	}

	// Download from MinIO
	object, err := api.minioClient.DownloadFile(ctx, objectName)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to download file: "+err.Error())
		return
	}
	defer object.Close()

	// Set headers
	w.Header().Set("Content-Type", fileInfo.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size, 10))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(objectName)))
	w.Header().Set("ETag", fileInfo.ETag)

	// Stream file to response
	_, err = io.Copy(w, object)
	if err != nil {
		// Can't write error response here as headers are already sent
		// Log the error instead
		fmt.Printf("Error streaming file: %v\n", err)
	}
}

// writeErrorResponse writes an error response in JSON format
func writeErrorResponse(w http.ResponseWriter, statusCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Error:   errorMsg,
	})
}

// RegisterRoutes registers the MinIO API routes on the provided mux
func (api *MinioAPI) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/upload", api.HandleUpload)
	mux.HandleFunc("/api/download", api.HandleDownload)
}
