package igin

import (
	"github.com/gin-gonic/gin"

	"github.com/starudream/creative-apartment/internal/ierr"
)

type AuthLimitFunc func(c *gin.Context) bool

func Auth(secret string, limit AuthLimitFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(ierr.NoAuth())
			return
		}

		if auth != secret {
			if limit != nil {
				if !limit(c) {
					return
				}
			}
			c.AbortWithStatusJSON(ierr.Forbidden())
			return
		}

		c.Next()
	}
}
