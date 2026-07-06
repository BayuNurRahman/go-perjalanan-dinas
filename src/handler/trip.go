package handler

import (
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/src/service"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TripHandler struct {
	tripService service.TripService
}

func NewTripHandler(tripService service.TripService) *TripHandler {
	return &TripHandler{tripService}
}

func (h *TripHandler) CreateTrip(c *gin.Context) {
	var input dto.CreateTripInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mengambil userID dari middleware
	userID := c.MustGet("userID").(uint)

	trip, err := h.tripService.CreateTrip(input, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengajukan dinas"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": trip})
}

func (h *TripHandler) GetAllTrips(c *gin.Context) {
	trips, err := h.tripService.GetAllTrips()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data monitoring",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   trips,
	})
}

func (h *TripHandler) UpdateStatus(c *gin.Context) {
	// 1. Ambil ID dari URL
	idParam := c.Param("id")
	
	// 2. Konversi string ke uint agar bisa diproses
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// 3. Ambil status dari JSON body
	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Panggil service (di sini variabel 'id' digunakan!)
	trip, err := h.tripService.UpdateStatus(uint(id), input.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trip})
}