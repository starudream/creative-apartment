package igin

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/starudream/creative-apartment/internal/ierr"
)

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if ev := recover(); ev != nil {
				stack := debug.Stack()
				log.Error().Msgf("[http] panic, %v\n%s", ev, stack)
				c.AbortWithStatusJSON(ierr.New().Internal())
			}
		}()
		c.Next()
	}
}
