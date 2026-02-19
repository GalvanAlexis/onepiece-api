package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alexisgalvan/onepiece-api/internal/domain"
	"github.com/alexisgalvan/onepiece-api/pkg/logger"
)

// Recovery returns a Gin middleware that recovers from panics,
// logs the error with Zap, and responds with a structured 500 JSON.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Error("panic recovered",
					zap.Any("error", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
				)

				errMsg := "internal server error"
				c.AbortWithStatusJSON(http.StatusInternalServerError, domain.Response{
					Success: false,
					Data:    nil,
					Error:   &errMsg,
				})
			}
		}()
		c.Next()
	}
}
