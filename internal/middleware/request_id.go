package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sake-kasu/sake-hack-backend/internal/logger"
)

// RequestID はリクエストごとにユニークなIDを生成し、コンテキストに設定するミドルウェア
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// X-Request-ID ヘッダーがあればそれを使用、なければ生成
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// レスポンスヘッダーに設定
		c.Header("X-Request-ID", requestID)

		// コンテキストに設定
		ctx := context.WithValue(c.Request.Context(), logger.RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
