package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/alexisgalvan/onepiece-api/internal/domain"
	"github.com/alexisgalvan/onepiece-api/internal/usecase"
)

// CharacterHandler handles HTTP requests for characters.
type CharacterHandler struct {
	uc *usecase.CharacterUseCase
}

// NewCharacterHandler creates a new CharacterHandler.
func NewCharacterHandler(uc *usecase.CharacterUseCase) *CharacterHandler {
	return &CharacterHandler{uc: uc}
}

// List godoc
// @Summary     List characters
// @Description Returns a paginated list of One Piece characters with optional filters
// @Tags        characters
// @Produce     json
// @Param       crew_id query int false "Filter by crew ID"
// @Param       status  query string false "Filter by status (alive|deceased|unknown)"
// @Param       limit   query int false "Number of results (default: 20, max: 100)"
// @Param       offset  query int false "Pagination offset"
// @Success     200 {object} domain.Response
// @Failure     500 {object} domain.Response
// @Router      /api/v1/characters [get]
func (h *CharacterHandler) List(c *gin.Context) {
	filters := domain.CharacterFilters{
		Limit:  parseIntQuery(c, "limit", 20),
		Offset: parseIntQuery(c, "offset", 0),
	}

	if crewID := parseIntQuery(c, "crew_id", 0); crewID > 0 {
		v := int64(crewID)
		filters.CrewID = &v
	}

	if status := c.Query("status"); status != "" {
		s := domain.CharacterStatus(status)
		filters.Status = &s
	}

	result, err := h.uc.List(c.Request.Context(), filters)
	if err != nil {
		errMsg := "failed to retrieve characters"
		c.JSON(http.StatusInternalServerError, domain.Response{
			Success: false,
			Data:    nil,
			Error:   &errMsg,
		})
		return
	}

	items := result.Items
	if items == nil {
		items = []domain.Character{}
	}

	c.JSON(http.StatusOK, domain.Response{
		Success: true,
		Data:    items,
		Error:   nil,
		Meta: &domain.PaginationMeta{
			Total:   int(result.Total),
			Limit:   filters.Limit,
			Offset:  filters.Offset,
			HasNext: int64(filters.Offset+filters.Limit) < result.Total,
		},
	})
}

// GetByID godoc
// @Summary     Get a character by ID
// @Description Returns a single character with crew and devil fruit info
// @Tags        characters
// @Produce     json
// @Param       id path int true "Character ID"
// @Success     200 {object} domain.Response
// @Failure     400 {object} domain.Response
// @Failure     404 {object} domain.Response
// @Failure     500 {object} domain.Response
// @Router      /api/v1/characters/{id} [get]
func (h *CharacterHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		errMsg := "invalid character id"
		c.JSON(http.StatusBadRequest, domain.Response{
			Success: false,
			Data:    nil,
			Error:   &errMsg,
		})
		return
	}

	character, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		errMsg := "failed to retrieve character"
		c.JSON(http.StatusInternalServerError, domain.Response{
			Success: false,
			Data:    nil,
			Error:   &errMsg,
		})
		return
	}

	if character == nil {
		errMsg := "character not found"
		c.JSON(http.StatusNotFound, domain.Response{
			Success: false,
			Data:    nil,
			Error:   &errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Success: true,
		Data:    character,
		Error:   nil,
	})
}

// parseIntQuery parses an integer query param with a default fallback.
func parseIntQuery(c *gin.Context, key string, defaultVal int) int {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil || n < 0 {
		return defaultVal
	}
	return n
}
