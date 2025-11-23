package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSConfig はCORS設定
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// CORS はCORSミドルウェアを返す
func CORS(config CORSConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     config.AllowedOrigins,
		AllowMethods:     config.AllowedMethods,
		AllowHeaders:     config.AllowedHeaders,
		ExposeHeaders:    config.ExposedHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           time.Duration(config.MaxAge) * time.Second,
	})
}
