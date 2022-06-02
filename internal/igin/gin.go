package igin

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"github.com/starudream/creative-apartment/internal/ivalidator"
)

func New() *gin.Engine {
	{
		validate, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			*validate = *(ivalidator.V())
		}
	}
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Use(recovery())
	return e
}

var (
	running int32

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
	atomic.StoreInt32(&running, 1)
	return server.ListenAndServe()
}

func Close() {
	if atomic.LoadInt32(&running) != 1 {
		return
	}

	log.Info().Msg("[http] server gracefully closing, will be closed in 5s")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Error().Msgf("[http] failed to shutdown server: %s", err)
		return
	}

	log.Info().Msg("[http] server shutdown")
}
