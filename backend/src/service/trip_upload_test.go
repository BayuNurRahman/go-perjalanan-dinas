package service

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveTripAttachmentsRequiresAtLeastOneFile(t *testing.T) {
	_, err := saveTripAttachments(nil, 1)
	if err == nil {
		t.Fatalf("expected error when no files are provided")
	}
}

func TestSaveTripAttachmentsStoresFiles(t *testing.T) {
	fileHeader := createMultipartFileHeader(t, "trip-attachment.pdf", "application/pdf", []byte("%PDF-1.4\n%fake pdf"))
	paths, err := saveTripAttachments([]*multipart.FileHeader{fileHeader}, 7)
	if err != nil {
		t.Fatalf("expected attachments to be saved, got error: %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("expected 1 saved path, got %d", len(paths))
	}
	if _, err := os.Stat(paths[0]); err != nil {
		t.Fatalf("expected saved file to exist, got error: %v", err)
	}
}

func TestSaveTripAttachmentsAcceptsPdfWithAlternativeMimeType(t *testing.T) {
	fileHeader := createMultipartFileHeader(t, "trip-attachment.pdf", "application/x-pdf", []byte("not a real pdf"))
	paths, err := saveTripAttachments([]*multipart.FileHeader{fileHeader}, 7)
	if err != nil {
		t.Fatalf("expected PDF attachments to be accepted, got error: %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("expected 1 saved path, got %d", len(paths))
	}
}

func TestGetTripUploadDirUsesProjectRoot(t *testing.T) {
	dir := getTripUploadDir(7)
	if dir == "" {
		t.Fatalf("expected upload directory to be resolved")
	}
	if !strings.HasSuffix(dir, filepath.Join("uploads", "trips", "user-7")) {
		t.Fatalf("expected upload directory to end with uploads/trips/user-7, got %s", dir)
	}
}

func createMultipartFileHeader(t *testing.T, filename, contentType string, content []byte) *multipart.FileHeader {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", "form-data; name=\"files\"; filename=\""+filename+"\"")
	header.Set("Content-Type", contentType)

	part, err := writer.CreatePart(header)
	if err != nil {
		t.Fatalf("unable to create multipart part: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("unable to write multipart content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("unable to close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(5 << 20); err != nil {
		t.Fatalf("unable to parse multipart form: %v", err)
	}

	return req.MultipartForm.File["files"][0]
}
