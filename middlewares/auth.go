package middlewares

import (
	"github.com/gin-gonic/gin"
)

func NoAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Code for middlewares
		c.Next()
	}
}
