package controllers

import (
	"besq-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ActivityLog defines the structure for a single activity entry.
type ActivityLog struct {
	User   string `json:"user"`
	Action string `json:"action"`
	Time   string `json:"time"`   // e.g., "10:45 AM"
	Status string `json:"status"` // e.g., "Completed"
}

// DashboardStats defines the structure for dashboard statistics.
type DashboardStats struct {
	TotalEmployees  int64         `json:"totalEmployees"`
	TotalOutput     int           `json:"totalOutput"`
	RejectRate      string        `json:"rejectRate"`
	ActiveShift     string        `json:"activeShift"`
	CuttingOutput   int           `json:"cuttingOutput"`
	PressingOutput  int           `json:"pressingOutput"`
	FinishingOutput int           `json:"finishingOutput"`
	Activities      []ActivityLog `json:"activities"`
}

// DashboardController handles dashboard-related requests.
type DashboardController struct {
	db *gorm.DB
}

// NewDashboardController creates a new instance of DashboardController.
func NewDashboardController(db *gorm.DB) *DashboardController {
	return &DashboardController{db: db}
}

// GetDashboardStats retrieves and returns dashboard statistics.
// It counts the total number of users (employees) and returns it along with placeholder stats.
func (dc *DashboardController) GetDashboardStats(c *gin.Context) {
	var count int64

	// Use GORM to count the rows in the 'users' table.
	if err := dc.db.Model(&models.User{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user count"})
		return
	}

	// Create mock activities as requested.
	mockActivities := []ActivityLog{
		{User: "Budi", Action: "Cutting Batch A", Time: "10:45 AM", Status: "Completed"},
		{User: "Ani", Action: "Pressing Batch C", Time: "11:02 AM", Status: "In Progress"},
		{User: "Rina", Action: "Finishing Batch B", Time: "11:15 AM", Status: "Pending"},
	}

	// Create stats object with a mix of real and placeholder data as requested.
	stats := DashboardStats{
		TotalEmployees:  count,
		TotalOutput:     0, // Placeholder
		RejectRate:      "0%", // Placeholder
		ActiveShift:     "Shift 1", // Placeholder
		CuttingOutput:   0, // New mock field
		PressingOutput:  0, // New mock field
		FinishingOutput: 0, // New mock field
		Activities:      mockActivities, // New mock field
	}

	c.JSON(http.StatusOK, stats)
}
