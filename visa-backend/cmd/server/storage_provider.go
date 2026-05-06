package main

import (
	"context"
	"errors"
	"io"
	"net/http"

	qstorage "visa-backend/internal/storage"
)

func newObjectStorage(cfg appConfig) qstorage.Uploader {
	return qstorage.NewQiniuUploader(qstorage.QiniuConfig{
		AccessKey: cfg.QiniuAccessKey,
		SecretKey: cfg.QiniuSecretKey,
		Bucket:    cfg.QiniuBucket,
		Domain:    cfg.QiniuDomain,
		UploadURL: cfg.QiniuUploadURL,
	})
}

func (s *appServer) uploadToObjectStorage(ctx context.Context, folder, prefix, originalName string, src io.Reader) (qstorage.ObjectInfo, error) {
	if s.storage == nil || !s.storage.Configured() {
		return qstorage.ObjectInfo{}, qstorage.ErrNotConfigured
	}
	return s.storage.Upload(ctx, folder, prefix, originalName, src)
}

func mapUploadError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	if errors.Is(err, qstorage.ErrNotConfigured) {
		writeError(w, http.StatusServiceUnavailable, errors.New("七牛对象存储尚未配置"))
		return
	}
	writeError(w, http.StatusInternalServerError, err)
}
