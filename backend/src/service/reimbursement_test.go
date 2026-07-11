package service

import (
	"testing"
	"time"

	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
)

type fakeReimbursementRepo struct {
	claims []models.Reimbursement
}

func (f *fakeReimbursementRepo) Create(claim models.Reimbursement) (models.Reimbursement, error) {
	f.claims = append(f.claims, claim)
	return claim, nil
}

func (f *fakeReimbursementRepo) FindAll() ([]models.Reimbursement, error) {
	return f.claims, nil
}

func (f *fakeReimbursementRepo) FindByID(id uint) (models.Reimbursement, error) {
	for _, claim := range f.claims {
		if claim.ID == id {
			return claim, nil
		}
	}
	return models.Reimbursement{}, nil
}

func (f *fakeReimbursementRepo) FindByTripID(tripID uint) ([]models.Reimbursement, error) {
	var claims []models.Reimbursement
	for _, claim := range f.claims {
		if claim.TripID == tripID {
			claims = append(claims, claim)
		}
	}
	return claims, nil
}

func (f *fakeReimbursementRepo) Update(claim models.Reimbursement) (models.Reimbursement, error) {
	for i, existing := range f.claims {
		if existing.ID == claim.ID {
			f.claims[i] = claim
			return claim, nil
		}
	}
	f.claims = append(f.claims, claim)
	return claim, nil
}

func (f *fakeReimbursementRepo) VerifyOwnership(claimID uint, userID uint) (bool, error) {
	return true, nil
}

func (f *fakeReimbursementRepo) Delete(id uint) error {
	for i, claim := range f.claims {
		if claim.ID == id {
			f.claims = append(f.claims[:i], f.claims[i+1:]...)
			return nil
		}
	}
	return nil
}

func TestSubmitClaimStoresReimbursement(t *testing.T) {
	repo := &fakeReimbursementRepo{}
	svc := NewReimbursementService(repo)

	tripID := uint(7)
	claim, err := svc.SubmitClaim(3, dto.SubmitClaimInput{TripID: tripID, Title: "Taxi", Description: "Airport transfer", Amount: 25000, TransactionDate: "2026-07-08"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claim.TripID != tripID || claim.Amount != 25000 {
		t.Fatalf("unexpected reimbursement payload: %+v", claim)
	}
}

func TestReviewClaimRejectsWithReason(t *testing.T) {
	repo := &fakeReimbursementRepo{}
	svc := NewReimbursementService(repo)

	claim, err := svc.SubmitClaim(3, dto.SubmitClaimInput{TripID: 8, Title: "Hotel", Description: "Lodging", Amount: 500000, TransactionDate: "2026-07-09"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	reviewed, err := svc.ReviewClaim(claim.ID, dto.ReviewClaimInput{Status: "REJECTED", RejectedReason: "Invoice missing"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if reviewed.Status != "REJECTED" || reviewed.RejectedReason != "Invoice missing" {
		t.Fatalf("expected rejected claim with reason, got %+v", reviewed)
	}
}

func TestReviewClaimAcceptsApprovedStatus(t *testing.T) {
	repo := &fakeReimbursementRepo{}
	svc := NewReimbursementService(repo)

	claim, err := svc.SubmitClaim(3, dto.SubmitClaimInput{TripID: 9, Title: "Meals", Description: "Per diem", Amount: 150000, TransactionDate: "2026-07-09"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	reviewed, err := svc.ReviewClaim(claim.ID, dto.ReviewClaimInput{Status: "APPROVED"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if reviewed.Status != "APPROVED" || reviewed.RejectedReason != "" {
		t.Fatalf("expected approved claim without rejection reason, got %+v", reviewed)
	}
}

func TestParseTransactionDate(t *testing.T) {
	parsed, err := time.Parse("2006-01-02", "2026-07-08")
	if err != nil {
		t.Fatalf("expected date to parse, got %v", err)
	}
	if parsed.Format("2006-01-02") != "2026-07-08" {
		t.Fatalf("unexpected parsed date %s", parsed.Format("2006-01-02"))
	}
}
