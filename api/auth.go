package api

import (
	"net/http"
)

type LoginReq struct {
	Code            string `json:"code" validate:"required,min=1"`
	Username        string `json:"username" validate:"required,min=1"`
	Password        string `json:"password" validate:"required,min=1"`
	RegistrationIds string `json:"registrationIds"`
}

type LoginResp struct {
	CommonResp
	Content struct {
		Id           string `json:"id"`
		Type         int    `json:"type"`
		IsRepairman  int    `json:"isRepairman"`
		UserName     string `json:"userName"`
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	} `json:"content"`
}

func Login(req LoginReq) *LoginResp {
	result, _ := Execute[*LoginResp](http.MethodPost, "/auth/authentication/customer/phone/app", nil, req)
	return result
}

type SendCodeReq struct {
	Type     int    `json:"type"`
	Username string `json:"username" validate:"required,min=1"`
}

type SendCodeResp struct {
	CommonResp
	Content string `json:"content"`
}

func SendCode(req SendCodeReq) *SendCodeResp {
	req.Type = 1
	result, _ := Execute[*SendCodeResp](http.MethodPost, "/auth/auth/sendCode", nil, req)
	return result
}

type LogoutReq struct {
	AccessToken string `json:"access_token"`
}

type LogoutResp struct {
	CommonResp
	Content string `json:"content"`
}

func Logout(accessToken string) *LogoutResp {
	result, _ := Execute[*LogoutResp](http.MethodPost, "/auth/auth/exit", HAuth(accessToken), &LogoutReq{AccessToken: accessToken})
	return result
}
