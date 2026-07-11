package service

import (
	"testing"

	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"

	"gorm.io/gorm"
)

type stubTripRepo struct {
	trips []models.BusinessTrip
}

func (s *stubTripRepo) Save(trip models.BusinessTrip) (models.BusinessTrip, error) {
	s.trips = append(s.trips, trip)
	return trip, nil
}

func (s *stubTripRepo) FindByID(id uint) (models.BusinessTrip, error) {
	for _, trip := range s.trips {
		if trip.ID == id {
			return trip, nil
		}
	}
	return models.BusinessTrip{}, gorm.ErrRecordNotFound
}

func (s *stubTripRepo) Update(trip models.BusinessTrip) (models.BusinessTrip, error) {
	for i, existing := range s.trips {
		if existing.ID == trip.ID {
			s.trips[i] = trip
			return trip, nil
		}
	}
	return trip, nil
}

func (s *stubTripRepo) FindAll() ([]models.BusinessTrip, error) {
	return s.trips, nil
}

func (s *stubTripRepo) Delete(id uint) error {
	for i, trip := range s.trips {
		if trip.ID == id {
			s.trips = append(s.trips[:i], s.trips[i+1:]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (s *stubTripRepo) FindByUserID(userID uint) ([]models.BusinessTrip, error) {
	var result []models.BusinessTrip
	for _, trip := range s.trips {
		if trip.UserID == userID {
			result = append(result, trip)
		}
	}
	return result, nil
}

func (s *stubTripRepo) FindByDepartmentID(departmentID uint) ([]models.BusinessTrip, error) {
	var result []models.BusinessTrip
	for _, trip := range s.trips {
		if trip.User.DepartmentID != nil && *trip.User.DepartmentID == departmentID {
			result = append(result, trip)
		}
	}
	return result, nil
}

func (s *stubTripRepo) FindUsersByDepartmentID(departmentID uint) ([]models.User, error) {
	var users []models.User
	seen := map[uint]bool{}
	for _, trip := range s.trips {
		if trip.User.DepartmentID != nil && *trip.User.DepartmentID == departmentID {
			if !seen[trip.UserID] {
				seen[trip.UserID] = true
				users = append(users, trip.User)
			}
		}
	}
	return users, nil
}

func (s *stubTripRepo) FindAllUsers() ([]models.User, error) {
	var users []models.User
	seen := map[uint]bool{}
	for _, trip := range s.trips {
		if !seen[trip.UserID] {
			seen[trip.UserID] = true
			users = append(users, trip.User)
		}
	}
	return users, nil
}

func TestGetManagerDashboardCountsStatuses(t *testing.T) {
	repo := &stubTripRepo{trips: []models.BusinessTrip{
		{ID: 1, Status: "PENDING"},
		{ID: 2, Status: "APPROVED"},
		{ID: 3, Status: "REJECTED"},
		{ID: 4, Status: "REVISION_REQUESTED"},
		{ID: 5, Status: "PENDING"},
	}}

	svc := &tripService{repo: repo}
	dashboard, err := svc.GetManagerDashboard()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	summary, ok := dashboard["summary"].(map[string]int)
	if !ok {
		t.Fatalf("expected summary map, got %T", dashboard["summary"])
	}

	if summary["pending"] != 2 || summary["approved"] != 1 || summary["rejected"] != 1 || summary["revision_requested"] != 1 {
		t.Fatalf("unexpected summary values: %#v", summary)
	}
}

func TestGetIncomingApplicationsIncludesPendingAndRevision(t *testing.T) {
	repo := &stubTripRepo{trips: []models.BusinessTrip{
		{ID: 1, Status: "PENDING"},
		{ID: 2, Status: "APPROVED"},
		{ID: 3, Status: "REVISION_REQUESTED"},
	}}

	svc := &tripService{repo: repo}
	applications, err := svc.GetIncomingApplications()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(applications) != 2 {
		t.Fatalf("expected 2 applications, got %d", len(applications))
	}
}

func TestGetIncomingApplicationsFiltersByManagerDepartment(t *testing.T) {
	deptA := uint(10)
	deptB := uint(20)
	repo := &stubTripRepo{trips: []models.BusinessTrip{
		{ID: 1, Status: "PENDING", User: models.User{DepartmentID: &deptA}},
		{ID: 2, Status: "REVISION_REQUESTED", User: models.User{DepartmentID: &deptB}},
		{ID: 3, Status: "PENDING", User: models.User{DepartmentID: &deptA}},
	}}

	svc := &tripService{repo: repo}
	applications, err := svc.GetIncomingApplications(deptA)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(applications) != 2 {
		t.Fatalf("expected 2 applications for department %d, got %d", deptA, len(applications))
	}

	for _, trip := range applications {
		if trip.User.DepartmentID == nil || *trip.User.DepartmentID != deptA {
			t.Fatalf("expected trip to belong to department %d, got %#v", deptA, trip.User.DepartmentID)
		}
	}
}

func TestUpdateTripUpdatesOwnedTrip(t *testing.T) {
	repo := &stubTripRepo{trips: []models.BusinessTrip{{ID: 7, UserID: 12, Destination: "Old", Description: "Old desc", Status: "PENDING"}}}
	svc := &tripService{repo: repo}

	updated, err := svc.UpdateTrip(7, 12, dto.UpdateTripInput{Destination: "New", Description: "Updated desc"}, nil)
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}

	if updated.Destination != "New" || updated.Description != "Updated desc" {
		t.Fatalf("expected updated trip fields, got %#v", updated)
	}
}

func TestDeleteTripDeletesOwnedTrip(t *testing.T) {
	repo := &stubTripRepo{trips: []models.BusinessTrip{{ID: 8, UserID: 15, Status: "PENDING"}}}
	svc := &tripService{repo: repo}

	err := svc.DeleteTrip(8, 15, "EMPLOYEE")
	if err != nil {
		t.Fatalf("expected delete to succeed, got %v", err)
	}

	if len(repo.trips) != 0 {
		t.Fatalf("expected trip to be removed, remaining: %d", len(repo.trips))
	}
}

func TestUpdateClaimRejectsCrossUserAccess(t *testing.T) {
	repo := &stubTripRepo{trips: []models.BusinessTrip{{ID: 7, UserID: 12}}}
	svc := &tripService{repo: repo}

	_, err := svc.UpdateClaim(7, 99, dto.UpdateClaimInput{Notes: "tampered"})
	if err == nil {
		t.Fatalf("expected unauthorized access error")
	}
}

var _ TripService = (*tripService)(nil)
var _ repositoryTripStub = (*stubTripRepo)(nil)

type repositoryTripStub interface {
	Save(trip models.BusinessTrip) (models.BusinessTrip, error)
	FindByID(id uint) (models.BusinessTrip, error)
	Update(trip models.BusinessTrip) (models.BusinessTrip, error)
	FindAll() ([]models.BusinessTrip, error)
	FindByUserID(userID uint) ([]models.BusinessTrip, error)
	FindByDepartmentID(departmentID uint) ([]models.BusinessTrip, error)
	FindUsersByDepartmentID(departmentID uint) ([]models.User, error)
	FindAllUsers() ([]models.User, error)
}
