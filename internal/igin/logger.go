package igin

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/starudream/creative-apartment/internal/iio"
	"github.com/starudream/creative-apartment/internal/iseq"
	"github.com/starudream/creative-apartment/internal/iu"
	"github.com/starudream/creative-apartment/internal/json"
)

func Logger() gin.HandlerFunc {
	var (
		gHeader = func(c *gin.Context, keys ...string) string {
			if len(keys) == 0 {
				return ""
			}
			for i := 0; i < len(keys); i++ {
				if v := c.GetHeader(keys[i]); v != "" {
					return v
				}
			}
			return ""
		}
		gPath = func(c *gin.Context) string {
			path := c.Request.URL.Path
			raw := c.Request.URL.RawQuery
			if raw != "" {
				path = path + "?" + raw
			}
			return path
		}
		gReqId = func(c *gin.Context) string {
			v := gHeader(c, HRequestId)
			if v == "" {
				v = iseq.UUID()
			}
			c.Set(HRequestId, v)
			c.Writer.Header().Set(HRequestId, v)
			return v
		}
		gLevel = func(c *gin.Context) zerolog.Level {
			switch {
			case c.Writer.Status() >= http.StatusInternalServerError:
				return zerolog.ErrorLevel
			case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
				return zerolog.WarnLevel
			default:
				return zerolog.InfoLevel
			}
		}
	)
	return func(c *gin.Context) {
		c.Writer = &responseWriter{ResponseWriter: c.Writer, bs: &bytes.Buffer{}}

		start := time.Now()

		l := log.With().
			Str("method", c.Request.Method).
			Str("path", gPath(c)).
			Str("tid", gReqId(c)).
			Str("ip", c.ClientIP()).
			Logger()

		var body []byte
		c.Request.Body, body = iio.ReadBody(c.Request.Body)
		l.Info().Str("span", "req").Str("type", gHeader(c, HContentType)).Msg(iu.Ternary[string](len(body) > 0, string(body), "-"))

		c.Next()

		latency := time.Now().Sub(start)

		l = l.With().Dur("latency", latency).Int("code", c.Writer.Status()).Logger()

		msg := ""
		if len(c.Errors.Errors()) > 0 {
			msg = json.MustMarshalString(c.Errors.JSON())
		} else {
			if w, ok := c.Writer.(*responseWriter); ok {
				msg = w.bs.String()
			}
		}

		l.WithLevel(gLevel(c)).Str("span", "resp").Msg(iu.Ternary[string](msg != "", msg, "-"))
	}
}

type responseWriter struct {
	gin.ResponseWriter

	bs *bytes.Buffer
}

var _ io.Writer = (*responseWriter)(nil)

func (w responseWriter) Write(p []byte) (n int, err error) {
	if w.bs == nil {
		w.bs = &bytes.Buffer{}
	}
	w.bs.Write(p)
	return w.ResponseWriter.Write(p)
}
