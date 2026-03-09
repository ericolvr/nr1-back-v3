package storage

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
)

type StorageClient struct {
	client     *storage.Client
	bucketName string
}

func NewStorageClient(ctx context.Context, bucketName string) (*StorageClient, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	return &StorageClient{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// UploadFile faz upload de um arquivo para o Cloud Storage
func (s *StorageClient) UploadFile(ctx context.Context, objectName string, data []byte, contentType string) (string, error) {
	bucket := s.client.Bucket(s.bucketName)
	obj := bucket.Object(objectName)

	// Create a writer
	wc := obj.NewWriter(ctx)
	wc.ContentType = contentType
	wc.CacheControl = "public, max-age=86400" // Cache por 1 dia

	// Write data
	if _, err := wc.Write(data); err != nil {
		return "", fmt.Errorf("failed to write data: %w", err)
	}

	// Close writer
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// Retorna a URL pública do arquivo
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucketName, objectName)
	log.Printf("✅ File uploaded successfully: %s", publicURL)

	return publicURL, nil
}

// DeleteFile deleta um arquivo do Cloud Storage
func (s *StorageClient) DeleteFile(ctx context.Context, objectName string) error {
	bucket := s.client.Bucket(s.bucketName)
	obj := bucket.Object(objectName)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	log.Printf("🗑️ File deleted successfully: %s", objectName)
	return nil
}

// GetFileURL retorna a URL pública de um arquivo
func (s *StorageClient) GetFileURL(objectName string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucketName, objectName)
}

// Close fecha o cliente de storage
func (s *StorageClient) Close() error {
	return s.client.Close()
}
