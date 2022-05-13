package igin

import (
	"github.com/gin-gonic/gin"

	"github.com/starudream/creative-apartment/internal/ierr"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(ierr.NoAuth())
			return
		}

		if auth != secret {
			c.AbortWithStatusJSON(ierr.Forbidden())
			return
		}

		c.Next()
	}
}
