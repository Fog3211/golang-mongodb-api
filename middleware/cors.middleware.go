package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// validate origin is in allowed origins
		allowedOrigins := getAllowedOrigins()
		isAllowed := false

		for _, allowed := range allowedOrigins {
			if allowed == origin {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}

		c.Next()
	}
}

func getAllowedOrigins() []string {
	return []string{
		"http://localhost:3000",
		"http://localhost:8000",
		"https://your-frontend.vercel.app",
	}
}
