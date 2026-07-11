package service

import (
	"testing"

	"go-perjalanan-dinas/models"
)

func TestPaginateTripsFiltersBySearchAndReturnsMetadata(t *testing.T) {
	trips := []models.BusinessTrip{
		{Destination: "Bandung", Status: "PENDING"},
		{Destination: "Jakarta", Status: "APPROVED"},
		{Destination: "Yogyakarta", Status: "PENDING"},
	}

	result := paginateTrips(trips, 1, 2, "band")
	items, ok := result["items"].([]models.BusinessTrip)
	if !ok {
		t.Fatalf("expected items to be []models.BusinessTrip, got %T", result["items"])
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item after search, got %d", len(items))
	}

	if items[0].Destination != "Bandung" {
		t.Fatalf("expected Bandung trip, got %s", items[0].Destination)
	}

	meta, ok := result["pagination"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected pagination metadata, got %T", result["pagination"])
	}

	if meta["page"] != 1 || meta["limit"] != 2 || meta["total"] != 1 || meta["total_pages"] != 1 {
		t.Fatalf("unexpected pagination metadata: %#v", meta)
	}
}
