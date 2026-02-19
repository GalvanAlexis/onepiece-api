package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alexisgalvan/onepiece-api/internal/config"
	"github.com/alexisgalvan/onepiece-api/internal/domain"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	version string
	env     string
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		version: cfg.App.Version,
		env:     cfg.App.Env,
	}
}

// HealthData represents the data payload of a health response.
type HealthData struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Env     string `json:"env"`
}

// Check godoc
// @Summary     Health check
// @Description Returns API status, version and environment
// @Tags        utility
// @Produce     json
// @Success     200 {object} domain.Response
// @Router      /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, domain.Response{
		Success: true,
		Data: HealthData{
			Status:  "ok",
			Version: h.version,
			Env:     h.env,
		},
		Error: nil,
	})
}
