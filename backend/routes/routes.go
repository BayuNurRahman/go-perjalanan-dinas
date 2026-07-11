package routes

import (
	"go-perjalanan-dinas/middleware"
	"go-perjalanan-dinas/src/handler"
	"go-perjalanan-dinas/src/repository"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, authHandler *handler.AuthHandler, departmentHandler *handler.DepartmentHandler, userHandler *handler.UserHandler, roleHandler *handler.RoleHandler, tripHandler *handler.TripHandler, reimbursementHandler *handler.ReimbursementHandler, blacklistRepo repository.BlacklistedTokenRepository) {
	api := r.Group("/api/v1")

	// Auth Routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", middleware.AuthMiddleware(blacklistRepo), middleware.RoleBlockMiddleware("SUPER_ADMIN"), authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", middleware.AuthMiddleware(blacklistRepo), authHandler.Logout)
	}

	// Department Routes
	departments := api.Group("/departments")
	departments.Use(middleware.AuthMiddleware(blacklistRepo))
	{
		departments.GET("/", departmentHandler.GetDepartments)
		departments.GET("", departmentHandler.GetDepartments)
		departments.GET("/:id", departmentHandler.GetDepartmentByID)
		departments.POST("/", middleware.RoleBlockMiddleware("SUPER_ADMIN"), departmentHandler.CreateDepartment)
		departments.POST("", middleware.RoleBlockMiddleware("SUPER_ADMIN"), departmentHandler.CreateDepartment)
		departments.PUT("/:id", middleware.RoleBlockMiddleware("SUPER_ADMIN"), departmentHandler.UpdateDepartment)
		departments.DELETE("/:id", middleware.RoleBlockMiddleware("SUPER_ADMIN"), departmentHandler.DeleteDepartment)
	}

	// User Routes
	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware(blacklistRepo))
	{
		users.GET("/", middleware.RoleBlockMiddleware("SUPER_ADMIN"), userHandler.GetUsers)
		users.GET("", middleware.RoleBlockMiddleware("SUPER_ADMIN"), userHandler.GetUsers)
		users.GET("/:id", userHandler.GetUserByID)
		users.PUT("/:id", middleware.RoleBlockMiddleware("SUPER_ADMIN"), userHandler.UpdateUser)
		users.DELETE("/:id", middleware.RoleBlockMiddleware("SUPER_ADMIN"), userHandler.DeleteUser)
	}

	// Role Routes
	roles := api.Group("/roles")
	roles.Use(middleware.AuthMiddleware(blacklistRepo))
	{
		roles.GET("/", roleHandler.GetRoles)
		roles.GET("", roleHandler.GetRoles)
		roles.GET("/:id", roleHandler.GetRoleByID)
		roles.POST("/", middleware.RoleBlockMiddleware("SUPER_ADMIN"), roleHandler.CreateRole)
		roles.POST("", middleware.RoleBlockMiddleware("SUPER_ADMIN"), roleHandler.CreateRole)
		roles.PUT("/:id", middleware.RoleBlockMiddleware("SUPER_ADMIN"), roleHandler.UpdateRole)
		roles.DELETE("/:id", middleware.RoleBlockMiddleware("SUPER_ADMIN"), roleHandler.DeleteRole)
	}

	// Trip Routes
	trip := api.Group("/trips")
	trip.Use(middleware.AuthMiddleware(blacklistRepo))
	{
		// Employee features
		trip.POST("/", middleware.RoleBlockMiddleware("EMPLOYEE"), tripHandler.CreateTrip)
		trip.POST("", middleware.RoleBlockMiddleware("EMPLOYEE"), tripHandler.CreateTrip)
		trip.POST("/request", middleware.RoleBlockMiddleware("EMPLOYEE"), tripHandler.CreateTrip)
		trip.GET("/dashboard", middleware.RoleBlockMiddleware("EMPLOYEE", "MANAGER", "SUPER_ADMIN"), tripHandler.GetEmployeeDashboard)
		trip.GET("/:id", middleware.RoleBlockMiddleware("EMPLOYEE", "MANAGER", "SUPER_ADMIN", "ADMIN_FIN"), tripHandler.GetTripByID)
		trip.GET("/:id/files/:filename", middleware.RoleBlockMiddleware("EMPLOYEE", "MANAGER", "SUPER_ADMIN", "ADMIN_FIN"), tripHandler.DownloadAttachment)
		trip.GET("/me", middleware.RoleBlockMiddleware("EMPLOYEE"), tripHandler.GetMyTrips)
		trip.PUT("/:id", middleware.RoleBlockMiddleware("EMPLOYEE"), tripHandler.UpdateTrip)
		trip.DELETE("/:id", middleware.RoleBlockMiddleware("EMPLOYEE", "MANAGER"), tripHandler.DeleteTrip)
		trip.PATCH("/:id/claim", middleware.RoleBlockMiddleware("EMPLOYEE"), tripHandler.UpdateClaim)

		// Manager/Admin monitoring
		trip.GET("/", middleware.RoleBlockMiddleware("MANAGER", "SUPER_ADMIN", "ADMIN_FIN"), tripHandler.GetAllTrips)
		trip.GET("", middleware.RoleBlockMiddleware("MANAGER", "SUPER_ADMIN", "ADMIN_FIN"), tripHandler.GetAllTrips)
		trip.GET("/manager/dashboard", middleware.RoleBlockMiddleware("MANAGER", "SUPER_ADMIN"), tripHandler.GetManagerDashboard)
		trip.GET("/manager/applications", middleware.RoleBlockMiddleware("MANAGER", "SUPER_ADMIN"), tripHandler.GetIncomingApplications)
		trip.GET("/manager/team-distribution", middleware.RoleBlockMiddleware("MANAGER", "SUPER_ADMIN"), tripHandler.GetTeamDistribution)
		trip.PATCH("/:id/status", middleware.RoleBlockMiddleware("MANAGER", "SUPER_ADMIN", "EMPLOYEE"), tripHandler.UpdateStatus)
		trip.PATCH("/:id/review-financial", middleware.RoleBlockMiddleware("ADMIN_FIN", "MANAGER"), tripHandler.ReviewFinancial)
		trip.PATCH("/:id/disburse", middleware.RoleBlockMiddleware("ADMIN_FIN", "MANAGER"), tripHandler.DisburseFunds)
	}

	claims := api.Group("/claims")
	claims.Use(middleware.AuthMiddleware(blacklistRepo))
	{
		claims.POST("/", middleware.RoleBlockMiddleware("EMPLOYEE"), reimbursementHandler.SubmitClaim)
		claims.POST("", middleware.RoleBlockMiddleware("EMPLOYEE"), reimbursementHandler.SubmitClaim)
		claims.GET("/trip/:trip_id", middleware.RoleBlockMiddleware("EMPLOYEE", "ADMIN_FIN", "MANAGER"), reimbursementHandler.GetClaimsByTripID)
		claims.GET("/:id/files/:filename", middleware.RoleBlockMiddleware("EMPLOYEE", "ADMIN_FIN", "MANAGER"), reimbursementHandler.DownloadClaimAttachment)
		claims.PATCH("/:id/review", middleware.RoleBlockMiddleware("ADMIN_FIN", "MANAGER"), reimbursementHandler.ReviewClaim)
		claims.GET("/:id", reimbursementHandler.GetClaimByID)
		claims.PUT("/:id", middleware.RoleBlockMiddleware("EMPLOYEE"), reimbursementHandler.UpdateClaim)
		claims.DELETE("/:id", middleware.RoleBlockMiddleware("EMPLOYEE"), reimbursementHandler.DeleteClaim)
	}
}
