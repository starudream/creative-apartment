package api

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/internal/ihttp"
	"github.com/starudream/creative-apartment/internal/ilog"
	"github.com/starudream/creative-apartment/internal/ivalidator"
)

type H map[string]string

func HAuth(accessToken string) H {
	return H{config.HAuthKey: config.HAuthValuePrefix + accessToken}
}

func Execute[T any](method, path string, headers H, body any) (T, bool) {
	var emptyResult, realResult T
	if body != nil && !ilog.WrapError(ivalidator.Struct(body), "validator") {
		return emptyResult, false
	}
	if v, exist := headers[config.HAuthKey]; exist && !ilog.WrapError(ivalidator.Var(strings.TrimPrefix(v, config.HAuthValuePrefix), "AccessToken", "min=1"), "validator") {
		return emptyResult, false
	}
	resp, err := ihttp.R().SetHeaders(headers).SetBody(body).SetResult(&realResult).SetError(&CommonResp{}).Execute(method, config.ApiAddr+path)
	if err != nil {
		log.Error().CallerSkipFrame(1).Msgf("send request error: %v", err)
		return emptyResult, false
	}
	if resp.IsError() {
		log.Error().CallerSkipFrame(1).Msgf("response status code: %d, body: %s", resp.StatusCode(), resp.String())
		return emptyResult, false
	}
	if v, ok := any(realResult).(xCommonResp); ok {
		if v.GetCode() != http.StatusOK {
			log.Error().CallerSkipFrame(1).Msgf("api error code: %d, message: %s", v.GetCode(), v.GetMessage())
			return emptyResult, false
		}
	}
	return realResult, true
}
