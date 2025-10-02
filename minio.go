// Package gotenberg provides MinIO storage integration
package gotenberg

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient wraps the MinIO client for file operations
type MinioClient struct {
	client     *minio.Client
	bucketName string
}

// MinioConfig contains configuration for MinIO connection
type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
}

// NewMinioClient creates a new MinIO client with the given configuration
func NewMinioClient(ctx context.Context, config MinioConfig) (*MinioClient, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Check if bucket exists, create if not
	exists, err := client.BucketExists(ctx, config.BucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinioClient{
		client:     client,
		bucketName: config.BucketName,
	}, nil
}

// UploadFile uploads a file to MinIO
// objectName - the name of the object in MinIO
// reader - the file content
// size - the size of the file in bytes (-1 for unknown size)
// contentType - the MIME type of the file (e.g., "application/pdf")
func (m *MinioClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (*minio.UploadInfo, error) {
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	info, err := m.client.PutObject(ctx, m.bucketName, objectName, reader, size, opts)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// DownloadFile downloads a file from MinIO
// objectName - the name of the object in MinIO
// Returns an io.ReadCloser that must be closed by the caller
func (m *MinioClient) DownloadFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	object, err := m.client.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	// Verify the object exists by getting its stat
	_, err = object.Stat()
	if err != nil {
		object.Close()
		return nil, err
	}

	return object, nil
}

// GetFileInfo returns information about a file in MinIO
func (m *MinioClient) GetFileInfo(ctx context.Context, objectName string) (minio.ObjectInfo, error) {
	return m.client.StatObject(ctx, m.bucketName, objectName, minio.StatObjectOptions{})
}

// DeleteFile deletes a file from MinIO
func (m *MinioClient) DeleteFile(ctx context.Context, objectName string) error {
	return m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
}

// ListFiles lists all files in the bucket with the given prefix
func (m *MinioClient) ListFiles(ctx context.Context, prefix string) <-chan minio.ObjectInfo {
	return m.client.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
}
