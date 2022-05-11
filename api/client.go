package api

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/internal/ihttp"
	"github.com/starudream/creative-apartment/internal/ivalidator"
)

type H map[string]string

func HAuth(accessToken string) H {
	return H{config.HAuthKey: config.HAuthValuePrefix + accessToken}
}

func Execute[T any](method, path string, headers H, body any) (T, bool) {
	if body != nil {
		if err := ivalidator.Struct(body); err != nil {
			return logError[T](err)
		}
	}
	if v, exist := headers[config.HAuthKey]; exist {
		if err := ivalidator.Var(strings.TrimPrefix(v, config.HAuthValuePrefix), "AccessToken", "min=1"); err != nil {
			return logError[T](err)
		}
	}
	var result T
	resp, err := ihttp.R().SetHeaders(headers).SetBody(body).SetResult(&result).SetError(&CommonResp{}).Execute(method, config.ApiAddr+path)
	if err != nil {
		return logError[T](err)
	}
	if resp.IsError() {
		return logError[T](err)
	}
	if v, ok := any(result).(xCommonResp); ok {
		if v.GetCode() != http.StatusOK {
			return logError[T](err)
		}
	}
	return result, true
}

func logError[T any](err error) (t T, b bool) {
	log.Error().CallerSkipFrame(3).Msgf("send request error: %v", err)
	return
}
