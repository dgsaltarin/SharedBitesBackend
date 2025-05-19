package ports

import (
	"context"
	"io"
)

type FileStore interface {
	UploadFile(ctx context.Context, file io.Reader, destinationPath, contentType string) (string, error)
	GetFileURL(ctx context.Context, storagePath string) (string, error)
	DeleteFile(ctx context.Context, storagePath string) error
}
