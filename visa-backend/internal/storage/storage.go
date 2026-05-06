package storage

import (
	"context"
	"errors"
	"io"
)

var ErrNotConfigured = errors.New("qiniu object storage is not configured")

type ObjectInfo struct {
	Key string
	URL string
}

type Uploader interface {
	Upload(ctx context.Context, folder, prefix, originalName string, src io.Reader) (ObjectInfo, error)
	Configured() bool
	Provider() string
}
