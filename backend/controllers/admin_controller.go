package controllers

import (
	"factory-api/database"
	"factory-api/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func GetAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	
	// Ambil parameter dari URL, contoh: /admin/audit-logs?page=1&limit=10
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	
	database.DB.Order("created_at desc").Limit(limit).Offset(offset).Find(&logs)
	
	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"data":  logs,
	})
}