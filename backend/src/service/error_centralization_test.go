package service

import (
	"testing"

	"go-perjalanan-dinas/models"
)

func TestSharedModelErrorsAreDefined(t *testing.T) {
	if models.ErrEmailAlreadyExists == nil {
		t.Fatal("expected ErrEmailAlreadyExists to be defined")
	}
	if models.ErrTripAccessDenied == nil {
		t.Fatal("expected ErrTripAccessDenied to be defined")
	}
	if models.ErrTripAttachmentRequired == nil {
		t.Fatal("expected ErrTripAttachmentRequired to be defined")
	}
}
