package igin

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/starudream/creative-apartment/internal/ierr"
)

func init() {
	log.Logger = log.Output(&zerolog.ConsoleWriter{Out: colorable.NewColorableStdout(), TimeFormat: "2006-01-02T15:04:05.000Z07:00"})
}

func TestMiddleware(t *testing.T) {
	t.Run("logger", func(t *testing.T) {
		h(http.MethodPost, "/v1", nil, Logger(), func(c *gin.Context) { c.JSON(ierr.OK()) })
		h(http.MethodPost, "/v2", nil)
	})

	t.Run("recovery", func(t *testing.T) {
		h(http.MethodPost, "/v2", nil, Logger(), func(c *gin.Context) { panic("no way") })
	})
}

func h(method, target string, body io.Reader, handlers ...gin.HandlerFunc) *httptest.ResponseRecorder {
	e := New()
	e.Handle(method, target, handlers...)
	r := httptest.NewRequest(method, target, body)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)
	return w
}
