package main

import (
	"go-perjalanan-dinas/config"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/routes"
	"go-perjalanan-dinas/src/handler"
	"go-perjalanan-dinas/src/repository"
	"go-perjalanan-dinas/src/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, &models.BusinessTrip{})

	r := gin.Default()

	// 1. Inisialisasi Auth Layers
	userRepo := repository.NewUserRepository(config.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService) // <-- Pastikan ini nama variabelnya

	// 2. Inisialisasi Trip Layers
	tripRepo := repository.NewTripRepository(config.DB)
	tripService := service.NewTripService(tripRepo)
	tripHandler := handler.NewTripHandler(tripService) // <-- Pastikan ini nama variabelnya

	// 3. Masukkan ke SetupRouter dengan urutan yang sesuai
	routes.SetupRouter(r, authHandler, tripHandler)

	r.Run(":8080")
}