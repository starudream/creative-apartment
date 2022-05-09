package api

type CommonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (v CommonResp) GetCode() int {
	return v.Code
}

func (v CommonResp) GetMessage() string {
	return v.Message
}

type xCommonResp interface {
	GetCode() int
	GetMessage() string
}
