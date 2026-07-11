package main

import (
	"go-perjalanan-dinas/config"
	"go-perjalanan-dinas/docs"
	"go-perjalanan-dinas/middleware"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/routes"
	"go-perjalanan-dinas/src/handler"
	"go-perjalanan-dinas/src/repository"
	"go-perjalanan-dinas/src/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Travel Dinas API
// @version 1.0
// @description API untuk sistem monitoring perjalanan dinas dengan autentikasi JWT, RBAC, upload dokumen, dan review keuangan.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadEnv()
	config.InitLogger()
	config.ConnectDatabase()
	config.ConnectMinio()
	config.DB.AutoMigrate(&models.Department{}, &models.User{}, &models.BusinessTrip{}, &models.Role{}, &models.BlacklistedToken{}, &models.Reimbursement{})

	r := gin.New()
	r.Use(gin.Recovery())

	// Apply CORS middleware FIRST (before any other middleware)
	r.Use(middleware.CORSMiddleware())

	// Apply structured logging middleware
	r.Use(middleware.LoggerMiddleware())

	// Add a health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	if docs.SwaggerInfo.Title != "" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 1. Inisialisasi Auth Layers
	userRepo := repository.NewUserRepository(config.DB)
	blacklistRepo := repository.NewBlacklistedTokenRepository(config.DB)
	authService := service.NewAuthServiceWithBlacklist(userRepo, blacklistRepo)
	authHandler := handler.NewAuthHandler(authService)

	// 2. Inisialisasi Department Layers
	departmentRepo := repository.NewDepartmentRepository(config.DB)
	departmentService := service.NewDepartmentService(departmentRepo)
	departmentHandler := handler.NewDepartmentHandler(departmentService)

	// 3. Inisialisasi User Layers
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 4. Inisialisasi Role Layers
	roleRepo := repository.NewRoleRepository(config.DB)
	roleService := service.NewRoleService(roleRepo)
	roleHandler := handler.NewRoleHandler(roleService)

	// 5. Inisialisasi Trip Layers
	tripRepo := repository.NewTripRepository(config.DB)
	tripService := service.NewTripService(tripRepo)
	tripHandler := handler.NewTripHandler(tripService)

	// 6. Inisialisasi Reimbursement Layers
	reimbursementRepo := repository.NewReimbursementRepository(config.DB)
	reimbursementService := service.NewReimbursementService(reimbursementRepo)
	reimbursementHandler := handler.NewReimbursementHandler(reimbursementService)

	// 7. Masukkan ke SetupRouter dengan urutan yang sesuai
	routes.SetupRouter(r, authHandler, departmentHandler, userHandler, roleHandler, tripHandler, reimbursementHandler, blacklistRepo)

	host := config.GetEnv("APP_HOST", "0.0.0.0")
	port := config.GetEnv("APP_PORT", "8080")
	r.Run(host + ":" + port)
}
