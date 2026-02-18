package main

import (
	"factory-api/controllers"
	"factory-api/database"
	"factory-api/middleware"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"time"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Sesuaikan dengan port frontend Svelte-mu
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 1. Inisialisasi Database & Seeder
	database.ConnectDatabase()
	database.SeedUsers()
	
	// 2. Route Public (Tanpa Token)
	r.POST("/login", controllers.Login)

	// 3. GROUP ADMIN (Hanya Role ADMIN)
	// Bisa: Tambah user, Lihat Log, Buat SPK/Work Order
admin := r.Group("/admin")
	admin.Use(middleware.AuthAndRoleMiddleware("ADMIN"))
	{
		admin.POST("/add-operator", controllers.CreateOperator) // Create
		admin.GET("/users", controllers.GetAllUsers)            // Read
		admin.PUT("/users/:id", controllers.UpdateUser)         // Update
		admin.DELETE("/users/:id", controllers.DeleteUser)      // Delete
		
		admin.GET("/audit-logs", controllers.GetAuditLogs)
		admin.POST("/work-order", controllers.CreateWorkOrder)
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
		prodAction.POST("/pressing/cycle", middleware.ActionMiddleware("OPERATOR_PRESSING"), controllers.RecordPressingCycle)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthAndRoleMiddleware(
		"ADMIN", "admin", 
		"OPERATOR_CUTTING", "cutting", 
		"OPERATOR_PRESSING", "pressing",
		"MANAGER", "manager",  
		"LEADER", "leader",
	)) 
	{
		api.GET("/users/profile", controllers.GetUserProfile)
		api.GET("/dashboard/stats", controllers.GetDashboardStats)
		api.GET("/pressing/today", controllers.GetPressingDashboard)
		api.POST("/scan-machine", controllers.ScanMachine)
		api.POST("/lwp", controllers.CreateLWP)
	}

	// -----------------------------------------------------------
    // 6. GROUP CHART DASHBOARD (Drill Down System)
    // -----------------------------------------------------------
    chartApi := r.Group("/api/chart")
    
    // Middleware: Izinkan ADMIN (untuk testing), MANAGER, dan LEADER
    // Pastikan role "MANAGER" dan "LEADER" nanti ada di tabel users jika ingin dipakai real
    chartApi.Use(middleware.AuthAndRoleMiddleware("ADMIN", "MANAGER", "LEADER")) 
    {
        // Level 1: Manager melihat Overview semua Proses
        // Usage: GET /api/chart/manager?tanggal=2026-02-01
        chartApi.GET("/manager", controllers.GetManagerOverview)

        // Level 2: Klik Proses -> Lihat Overview Mesin
        // Usage: GET /api/chart/process?tanggal=2026-02-01&proses=PRS
        chartApi.GET("/process", controllers.GetLeaderProcessView)

    // Level 3a: Klik Mesin -> Lihat Summary/Overview Mesin
    // Usage: GET /api/chart/machine?tanggal=2026-02-01&no_mc=04A
    chartApi.GET("/machine", controllers.GetMachineDetail)
    
    // Level 3b: Detail Per Jam dengan Item Info untuk Chart & Tabel
    // Usage: GET /api/chart/machine/hourly?tanggal=2026-02-01&no_mc=04A&shift=1
    chartApi.GET("/machine/hourly", controllers.GetChartDataByMachine)
    }

	r.Run(":8080")
}