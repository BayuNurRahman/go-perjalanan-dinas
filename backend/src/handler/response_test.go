package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-perjalanan-dinas/dto"

	"github.com/gin-gonic/gin"
)

func TestWriteSuccessUsesWebResponseShape(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	writeSuccess(c, http.StatusOK, "Operasi berhasil", map[string]string{"id": "1"})

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response dto.WebResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid JSON response, got %v", err)
	}

	if !response.Success {
		t.Fatalf("expected success=true, got %v", response.Success)
	}
	if response.Message != "Operasi berhasil" {
		t.Fatalf("expected message to be preserved, got %q", response.Message)
	}
	if response.Data == nil {
		t.Fatal("expected data payload to be present")
	}
}
