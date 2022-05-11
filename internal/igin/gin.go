package igin

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Use(recovery())
	return e
}

var (
	server *http.Server

	handler *gin.Engine

	handlerOnce sync.Once
)

func S() *gin.Engine {
	handlerOnce.Do(func() {
		handler = New()
	})
	return handler
}

func Run(addr string) error {
	server = &http.Server{Addr: addr, Handler: handler}
	log.Info().Msgf("[http] listening on %s", addr)
	return server.ListenAndServe()
}

func Close() {
	if server == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Error().Msgf("[http] failed to shutdown server: %s", err)
		return
	}

	log.Info().Msg("[http] server shutdown")
}
