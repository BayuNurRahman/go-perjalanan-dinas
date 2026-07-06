package routes

import (
	"go-perjalanan-dinas/middleware"
	"go-perjalanan-dinas/src/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, authHandler *handler.AuthHandler, tripHandler *handler.TripHandler) {
	api := r.Group("/api/v1")

	// Auth Routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Trip Routes
	trip := api.Group("/trips")
	trip.Use(middleware.AuthMiddleware())
	{
		// Karyawan (EMPLOYEE) membuat pengajuan baru
		trip.POST("/", middleware.RoleBlockMiddleware("EMPLOYEE"), tripHandler.CreateTrip)

		// MANAGER dan ADMIN dapat melihat monitoring semua pegawai
		trip.GET("/", middleware.RoleBlockMiddleware("MANAGER", "ADMIN"), tripHandler.GetAllTrips)

		// Eksklusif MANAGER untuk Approve / Reject
		trip.PATCH("/:id/status", middleware.RoleBlockMiddleware("MANAGER"), tripHandler.UpdateStatus)

		// Hanya MANAGER dan ADMIN yang bisa melihat monitoring semua pegawai (Global Monitoring)
		trip.GET("/", middleware.RoleBlockMiddleware("MANAGER", "ADMIN"), tripHandler.GetAllTrips)
	}
}