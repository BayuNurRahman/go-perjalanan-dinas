package service

import (
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/repository"
	"time"
)

type TripService interface {
	CreateTrip(input dto.CreateTripInput, userID uint) (models.BusinessTrip, error)
	GetAllTrips() ([]models.BusinessTrip, error)
	UpdateStatus(id uint, status string) (models.BusinessTrip, error)
}

type tripService struct {
	repo repository.TripRepository
}

func NewTripService(repo repository.TripRepository) TripService {
	return &tripService{repo}
}

func (s *tripService) CreateTrip(input dto.CreateTripInput, userID uint) (models.BusinessTrip, error) {
	// Konversi string ke time.Time
	start, _ := time.Parse("2006-01-02", input.StartDate)
	end, _ := time.Parse("2006-01-02", input.EndDate)

	trip := models.BusinessTrip{
		UserID:      userID,
		Destination: input.Destination,
		StartDate:   start,
		EndDate:     end,
		Status:      "PENDING",
	}

	return s.repo.Save(trip)
}

func (s *tripService) GetAllTrips() ([]models.BusinessTrip, error) {
    return s.repo.FindAll()
}

func (s *tripService) UpdateStatus(id uint, status string) (models.BusinessTrip, error) {
    trip, err := s.repo.FindByID(id)
    if err != nil {
        return trip, err
    }

    trip.Status = status
    return s.repo.Update(trip)
}