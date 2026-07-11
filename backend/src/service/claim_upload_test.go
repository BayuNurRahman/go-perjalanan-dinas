package service

import (
	"mime/multipart"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveClaimAttachmentsStoresFilesInClaimDirectory(t *testing.T) {
	fileHeader := createClaimMultipartFileHeader(t, "claim-proof.pdf", "application/pdf", []byte("%PDF-1.4\n%fake pdf"))

	paths, err := saveClaimAttachments([]*multipart.FileHeader{fileHeader}, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("expected 1 saved path, got %d", len(paths))
	}
	if !strings.Contains(paths[0], filepath.Join("uploads", "claims", "claim-3")) {
		t.Fatalf("expected saved path to be under uploads/claims/claim-3, got %s", paths[0])
	}
}

func createClaimMultipartFileHeader(t *testing.T, filename string, contentType string, content []byte) *multipart.FileHeader {
	t.Helper()

	var body strings.Builder
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("files", filename)
	if err != nil {
		t.Fatalf("unable to create multipart form file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("unable to write multipart content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("unable to close multipart writer: %v", err)
	}

	req := httptest.NewRequest("POST", "/", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_ = contentType
	if err := req.ParseMultipartForm(5 << 20); err != nil {
		t.Fatalf("unable to parse multipart form: %v", err)
	}

	return req.MultipartForm.File["files"][0]
}
