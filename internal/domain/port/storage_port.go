package port

import (
	"context"
	"io"
)

type ObjectStorage interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType string) (string, error)
	GetPresignedURL(ctx context.Context, key string) (string, error)
}
