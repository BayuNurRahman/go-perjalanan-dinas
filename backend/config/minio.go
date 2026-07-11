package config

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client
var MinioBucket string

func ConnectMinio() {
	LoadEnv()

	endpoint := GetEnv("MINIO_ENDPOINT", "localhost:9000")
	accessKeyID := GetEnv("MINIO_ACCESS_KEY", "minioadmin")
	secretAccessKey := GetEnv("MINIO_SECRET_KEY", "minioadmin")
	useSSLStr := GetEnv("MINIO_USE_SSL", "false")
	MinioBucket = GetEnv("MINIO_BUCKET", "perjalanan-dinas")

	useSSL, err := strconv.ParseBool(useSSLStr)
	if err != nil {
		useSSL = false
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		slog.Error("Gagal inisialisasi MinIO client", "error", err)
		return
	}

	MinioClient = client
	slog.Info("Koneksi MinIO berhasil diinisialisasi", "endpoint", endpoint)

	// Create bucket if it doesn't exist (fail gracefully so tests don't crash without MinIO)
	ctx := context.Background()
	exists, err := MinioClient.BucketExists(ctx, MinioBucket)
	if err != nil {
		slog.Warn("Gagal memeriksa bucket MinIO. Pastikan server MinIO berjalan di "+endpoint, "error", err)
		return
	}

	if !exists {
		err = MinioClient.MakeBucket(ctx, MinioBucket, minio.MakeBucketOptions{})
		if err != nil {
			slog.Error("Gagal membuat bucket MinIO", "bucket", MinioBucket, "error", err)
			return
		}
		slog.Info("Bucket MinIO berhasil dibuat", "bucket", MinioBucket)
	} else {
		slog.Info("Bucket MinIO sudah tersedia", "bucket", MinioBucket)
	}
}

// GetMinioKey parses a local path or key and returns a normalized MinIO object key.
func GetMinioKey(pathOrKey string, category string, id uint) string {
	// If already prefixed with category + "/", it's a MinIO key
	if strings.HasPrefix(pathOrKey, category+"/") {
		return pathOrKey
	}

	// Normalise path separators
	normalized := strings.ReplaceAll(pathOrKey, "\\", "/")
	base := filepath.Base(normalized)

	// Determine prefix category ID name (e.g. category "claims" -> prefix "claim-3")
	var idPrefix string
	switch category {
	case "claims":
		idPrefix = fmt.Sprintf("claim-%d", id)
	case "trips":
		idPrefix = fmt.Sprintf("user-%d", id)
	default:
		idPrefix = fmt.Sprintf("id-%d", id)
	}

	return fmt.Sprintf("%s/%s/%s", category, idPrefix, base)
}
