package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"wacast/core/config"
	"wacast/core/database"
	"wacast/core/utils"
)

type ConfigHandler struct {
	db     *database.Database
	config *config.Config
}

func NewConfigHandler(db *database.Database, cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{
		db:     db,
		config: cfg,
	}
}

func (h *ConfigHandler) RegisterRoutes(r *gin.RouterGroup) {
	configGroup := r.Group("/config")
	{
		configGroup.GET("/db", h.GetDBStatus)
		// Disabled for SaaS version for security (only configurable via .env)
		// configGroup.POST("/db/test", h.TestDBConnection)
		// configGroup.POST("/db/save", h.SaveDBConnection)
	}
}

func (h *ConfigHandler) GetDBStatus(c *gin.Context) {
	status := "connected"
	if !h.db.HealthCheck() {
		status = "disconnected"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  status,
		"config": gin.H{
			"host":     h.config.Database.Host,
			"port":     h.config.Database.Port,
			"user":     h.config.Database.User,
			"database": h.config.Database.DBName,
			"ssl_mode": h.config.Database.SSLMode,
		},
	})
}

func (h *ConfigHandler) TestDBConnection(c *gin.Context) {
	var req config.DatabaseConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	// Ensure defaults for optional fields if not provided
	if req.Port == 0 {
		req.Port = 5432
	}
	if req.ConnectionTimeout == 0 {
		req.ConnectionTimeout = 5
	}

	// Try to open and ping
	testDb, err := database.InitDatabase(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Connection failed: " + err.Error(),
		})
		return
	}
	defer testDb.Close()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Database connection test successful",
	})
}

func (h *ConfigHandler) SaveDBConnection(c *gin.Context) {
	var req config.DatabaseConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	// 1. Save to .env
	if err := config.SaveDatabaseConfigToEnv(&req); err != nil {
		utils.Error("Failed to save DB config to .env", zap.Error(err))
		internalError(c, err)
		return
	}

	// 2. Hot-swap connection
	if err := h.db.UpdateConnection(&req); err != nil {
		utils.Error("Failed to hot-swap DB connection", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Config saved to .env but failed to update live connection: " + err.Error(),
		})
		return
	}

	// 3. Update in-memory config for future use
	h.config.Database = &req

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Database configuration saved and updated successfully",
	})
}
