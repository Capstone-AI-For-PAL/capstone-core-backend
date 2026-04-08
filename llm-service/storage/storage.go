package storage

import "context"

// Storage abstracts file upload so implementations can be swapped
// (e.g. S3, GCS, local filesystem).
type Storage interface {
	Upload(ctx context.Context, path string, data []byte, contentType string) (url string, err error)
}
