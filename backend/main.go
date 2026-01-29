package main

import (
	"besq-backend/controllers"
	"besq-backend/middleware"
	"besq-backend/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func main() {
	// Initialize database
	initDatabase()

	// Run auto-migration
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed default admin user
	seedDefaultAdmin()

	// Initialize controllers
	authController := controllers.NewAuthController(db)
	userController := controllers.NewUserController(db)
	dashboardController := controllers.NewDashboardController(db)

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Basic routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is healthy",
			"status":  "ok",
		})
	})

	// Public authentication route
	r.POST("/login", authController.Login)
	
	// API group
	api := r.Group("/api")
	// Use the authentication middleware for all routes in this group
	api.Use(middleware.AuthMiddleware()) 
	{
		// Welcome message for the API
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to BESQ Backend API",
				"version": "1.0.0",
			})
		})

		// Dashboard statistics route
		api.GET("/dashboard/stats", dashboardController.GetDashboardStats)

		// Profile route accessible to any authenticated user
		api.GET("/users/profile", userController.GetProfile)

		// Admin-only routes for user management
		adminRoutes := api.Group("/users")
		adminRoutes.Use(middleware.RequireRoles("admin"))
		{
			adminRoutes.GET("/", userController.GetUsers)
			adminRoutes.POST("/", userController.CreateUser)
			adminRoutes.GET("/:id", userController.GetUser)
			adminRoutes.PUT("/:id", userController.UpdateUser)
			adminRoutes.DELETE("/:id", userController.DeleteUser)
		}
	}

	log.Println("Server starting on port 8080...")
	r.Run(":8080")
}

// initDatabase initializes the SQLite database connection using a pure Go driver.
func initDatabase() {
	var err error
	// Using "github.com/glebarez/sqlite" for CGO-free operation.
	db, err = gorm.Open(sqlite.Open("besq.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")
}

// seedDefaultAdmin creates a default admin user if no users exist.
func seedDefaultAdmin() {
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)

	if userCount == 0 {
		// Hash the default password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}

		defaultAdmin := models.User{
			NIK:      "admin",
			Name:     "Default Admin",
			Password: string(hashedPassword),
			Role:     "admin",
		}

		if err := db.Create(&defaultAdmin).Error; err != nil {
			log.Fatal("Failed to create default admin user:", err)
		}

		log.Println("Default admin user created successfully")
		log.Println("NIK: admin, Password: admin123")
	} else {
		log.Println("Users table is not empty, skipping seeding")
	}
}