package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductionController handles production-related operations (Cutting & Pressing)
type ProductionController struct {
	db *gorm.DB
}

// NewProductionController creates a new instance of ProductionController
func NewProductionController(db *gorm.DB) *ProductionController {
	return &ProductionController{db: db}
}

// CreateCuttingEntry creates a new cutting production entry
func (pc *ProductionController) CreateCuttingEntry(c *gin.Context) {
	var cuttingEntry struct {
		Quantity int    `json:"quantity" binding:"required"`
		Notes    string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&cuttingEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create cutting entry (to be stored in database)
	// For now, returning success response
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Cutting entry created successfully",
		"quantity": cuttingEntry.Quantity,
		"notes":    cuttingEntry.Notes,
	})
}

// GetCuttingLogs retrieves today's cutting production logs
func (pc *ProductionController) GetCuttingLogs(c *gin.Context) {
	// Retrieve cutting logs for today
	// For now, returning empty list
	c.JSON(http.StatusOK, gin.H{
		"cutting_logs": []interface{}{},
		"total":        0,
	})
}

// CreatePressingEntry creates a new pressing production entry
func (pc *ProductionController) CreatePressingEntry(c *gin.Context) {
	var pressingEntry struct {
		Quantity int    `json:"quantity" binding:"required"`
		Notes    string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&pressingEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create pressing entry (to be stored in database)
	// For now, returning success response
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Pressing entry created successfully",
		"quantity": pressingEntry.Quantity,
		"notes":    pressingEntry.Notes,
	})
}

// GetPressingLogs retrieves today's pressing production logs
func (pc *ProductionController) GetPressingLogs(c *gin.Context) {
	// Retrieve pressing logs for today
	// For now, returning empty list
	c.JSON(http.StatusOK, gin.H{
		"pressing_logs": []interface{}{},
		"total":         0,
	})
}
