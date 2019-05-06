package middlewares

import (
	"github.com/gin-gonic/gin"
)

// Authorized blocks unauthorized requestrs
func Authorized(c *gin.Context) {
	_, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(401, gin.H{
			"status": false,
			"message": "Unauthorization!!",
		})

		return
	}
}
