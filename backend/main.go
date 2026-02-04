package main

import (
	"factory-api/controllers"
	"factory-api/database"
	"factory-api/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 1. Inisialisasi Database & Seeder
	database.ConnectDatabase()
	database.SeedAdmin()

	// 2. Route Public (Tanpa Token)
	r.POST("/login", controllers.Login)

	// 3. GROUP ADMIN (Hanya Role ADMIN)
	// Bisa: Tambah user, Lihat Log, Buat SPK/Work Order
	admin := r.Group("/admin")
	admin.Use(middleware.AuthAndRoleMiddleware("ADMIN"))
	{
		admin.POST("/add-operator", controllers.CreateOperator)
		admin.GET("/audit-logs", controllers.GetAuditLogs)    // Pastikan func ini ada di admin_controller.go
		admin.POST("/work-order", controllers.CreateWorkOrder) // Endpoint baru
	}

	// 4. GROUP PRODUKSI UMUM (Bisa diakses Admin & Semua Operator)
	// Bisa: Lihat daftar pekerjaan (Read Only)
	prodView := r.Group("/production")
	prodView.Use(middleware.AuthAndRoleMiddleware("ADMIN", "OPERATOR_CUTTING", "OPERATOR_PRESSING"))
	{
		prodView.GET("/work-orders", controllers.GetAllWorkOrders)
		prodView.GET("/cutting/status", controllers.GetCuttingStatus)
		prodView.GET("/pressing/status", controllers.GetPressingStatus)
	}

	// 5. GROUP AKSI OPERATOR (Update Status & Jalankan Mesin)
	// Menggunakan ActionMiddleware agar write access terkontrol
	prodAction := r.Group("/production")
	prodAction.Use(middleware.AuthAndRoleMiddleware("ADMIN", "OPERATOR_CUTTING", "OPERATOR_PRESSING"))
	{
		// Operator bisa update status WO (Misal: Cutting selesai, lanjut Pressing)
		prodAction.PATCH("/work-order/:id/status", controllers.UpdateWOStatus)
		
		// Endpoint spesifik mesin (Butuh Role Spesifik di ActionMiddleware)
		prodAction.POST("/cutting/start", middleware.ActionMiddleware("OPERATOR_CUTTING"), controllers.StartCutting)
		prodAction.POST("/pressing/start", middleware.ActionMiddleware("OPERATOR_PRESSING"), controllers.StartPressing)
	}

	r.Run(":8080")
}