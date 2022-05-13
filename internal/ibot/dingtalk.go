package ibot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/starudream/creative-apartment/internal/icrypto"
	"github.com/starudream/creative-apartment/internal/ihttp"
)

// API Document: https://open.dingtalk.com/document/group/custom-robot-access

var Dingtalk = dingtalk{}

type dingtalk struct {
	Token  string
	Secret string
}

var _ Interface = (*dingtalk)(nil)

type dingtalkReq struct {
	MsgType string        `json:"msgtype"`
	Text    *dingtalkText `json:"text"`
}

type dingtalkText struct {
	Content string `json:"content"`
}

type dingtalkResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (m dingtalk) Init(token string, secret string) dingtalk {
	m.Token = token
	m.Secret = secret
	return m
}

func (m dingtalk) SendMessage(text string) error {
	if m.Token == "" {
		m.Token = viper.GetString("dingtalk.token")
		m.Secret = viper.GetString("dingtalk.secret")
	}
	if m.Token == "" {
		return fmt.Errorf("[dingtalk] token is empty")
	}
	addr := "https://oapi.dingtalk.com/robot/send?access_token=" + m.Token
	if m.Secret != "" {
		milli := strconv.FormatInt(time.Now().UnixMilli(), 10)
		sign := icrypto.ECB.SHA256(milli+"\n"+m.Secret, m.Secret).Base64Std()
		addr += "&timestamp=" + milli + "&sign=" + sign
	}
	if text == "" {
		text = "这是一条测试消息"
	}
	req := &dingtalkReq{MsgType: "text", Text: &dingtalkText{Content: text}}
	resp, err := ihttp.R().SetHeader("content-type", "application/json").SetBody(req).SetResult(&dingtalkResp{}).Post(addr)
	if err != nil {
		return fmt.Errorf("[dingtalk] send message error: %v", err)
	}
	if e, ok := resp.Result().(*dingtalkResp); ok && e.ErrCode != 0 {
		return fmt.Errorf("[dingtalk] send message error: %d, %s", e.ErrCode, e.ErrMsg)
	}
	return nil
}
