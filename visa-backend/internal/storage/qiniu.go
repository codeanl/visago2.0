package storage

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	qiniustorage "github.com/qiniu/go-sdk/v7/storage"
)

type QiniuConfig struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Domain    string
	UploadURL string
	TokenTTL  time.Duration
	Client    *http.Client
}

type qiniuUploader struct {
	accessKey string
	secretKey string
	bucket    string
	domain    string
	tokenTTL  time.Duration
	config    *qiniustorage.Config
}

func NewQiniuUploader(cfg QiniuConfig) Uploader {
	tokenTTL := cfg.TokenTTL
	if tokenTTL <= 0 {
		tokenTTL = time.Hour
	}

	uploadCfg := &qiniustorage.Config{
		UseHTTPS: true,
	}

	// Keep a custom upload host hook point for special regions or private gateways.
	if uploadURL := strings.TrimSpace(cfg.UploadURL); uploadURL != "" {
		uploadCfg.UseCdnDomains = false
		uploadCfg.RsHost = uploadURL
	}

	return &qiniuUploader{
		accessKey: strings.TrimSpace(cfg.AccessKey),
		secretKey: strings.TrimSpace(cfg.SecretKey),
		bucket:    strings.TrimSpace(cfg.Bucket),
		domain:    normalizeDomain(cfg.Domain),
		tokenTTL:  tokenTTL,
		config:    uploadCfg,
	}
}

func (q *qiniuUploader) Configured() bool {
	return q.accessKey != "" && q.secretKey != "" && q.bucket != "" && q.domain != ""
}

func (q *qiniuUploader) Provider() string {
	return "Qiniu"
}

func (q *qiniuUploader) Upload(ctx context.Context, folder, prefix, originalName string, src io.Reader) (ObjectInfo, error) {
	if !q.Configured() {
		return ObjectInfo{}, ErrNotConfigured
	}

	objectKey := buildObjectKey(folder, prefix, originalName)
	putPolicy := qiniustorage.PutPolicy{
		Scope:   q.bucket,
		Expires: uint64(q.tokenTTL.Seconds()),
	}
	mac := qbox.NewMac(q.accessKey, q.secretKey)
	upToken := putPolicy.UploadToken(mac)

	data, err := io.ReadAll(src)
	if err != nil {
		return ObjectInfo{}, err
	}

	formUploader := qiniustorage.NewFormUploader(q.config)
	var ret struct {
		Key  string `json:"key"`
		Hash string `json:"hash"`
	}

	putExtra := qiniustorage.PutExtra{}
	if err := formUploader.Put(ctx, &ret, upToken, objectKey, bytes.NewReader(data), int64(len(data)), &putExtra); err != nil {
		return ObjectInfo{}, err
	}

	if strings.TrimSpace(ret.Key) == "" {
		ret.Key = objectKey
	}

	return ObjectInfo{
		Key: ret.Key,
		URL: q.domain + "/" + escapeObjectKey(ret.Key),
	}, nil
}

func buildObjectKey(folder, prefix, originalName string) string {
	folder = strings.Trim(strings.ReplaceAll(strings.TrimSpace(folder), "\\", "/"), "/")
	if folder == "" {
		folder = "misc"
	}
	prefix = strings.Trim(strings.ReplaceAll(strings.TrimSpace(prefix), "\\", "-"), "/")
	if prefix == "" {
		prefix = "file"
	}
	ext := normalizeExt(originalName)
	filename := prefix + "-" + time.Now().Format("20060102150405.000000000") + ext
	return folder + "/" + filename
}

func normalizeExt(originalName string) string {
	ext := strings.ToLower(strings.TrimSpace(path.Ext(originalName)))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return ext
	default:
		return ".jpg"
	}
}

func normalizeDomain(raw string) string {
	value := strings.TrimSpace(raw)
	if value == "" {
		return ""
	}
	if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		value = "https://" + value
	}
	return strings.TrimRight(value, "/")
}

func escapeObjectKey(key string) string {
	parts := strings.Split(strings.TrimLeft(key, "/"), "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}
