package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alexisgalvan/onepiece-api/internal/domain"
	"github.com/alexisgalvan/onepiece-api/internal/usecase"
)

// APIKeyHandler handles HTTP requests for API key management.
type APIKeyHandler struct {
	uc *usecase.APIKeyUseCase
}

// NewAPIKeyHandler creates a new APIKeyHandler.
func NewAPIKeyHandler(uc *usecase.APIKeyUseCase) *APIKeyHandler {
	return &APIKeyHandler{uc: uc}
}

// GenerateRequest is the request body for generating an API key.
type GenerateRequest struct {
	Label     string `json:"label"`
	RateLimit int    `json:"rate_limit_per_min"`
}

// Generate godoc
// @Summary     Generate an API key
// @Description Creates a new API key. The plain-text key is returned ONCE and not stored.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body GenerateRequest false "API key config"
// @Success     201 {object} domain.Response
// @Failure     500 {object} domain.Response
// @Router      /api/v1/auth/api-key [post]
func (h *APIKeyHandler) Generate(c *gin.Context) {
	var req GenerateRequest
	// Ignore binding error — defaults are fine
	_ = c.ShouldBindJSON(&req)

	result, err := h.uc.Generate(c.Request.Context(), req.Label, req.RateLimit)
	if err != nil {
		errMsg := "failed to generate API key"
		c.JSON(http.StatusInternalServerError, domain.Response{
			Success: false,
			Data:    nil,
			Error:   &errMsg,
		})
		return
	}

	c.JSON(http.StatusCreated, domain.Response{
		Success: true,
		Data:    result,
		Error:   nil,
	})
}
