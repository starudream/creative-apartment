package api

import (
	"net/http"
)

type GetLatestAPKResp struct {
	CommonResp
	Content struct {
		Source        int    `json:"source"`
		IsToUpdate    int    `json:"isToUpdate"`
		VersionNumber int    `json:"versionNumber"`
		ApkAddress    string `json:"apkAdress"`
		CreateTime    string `json:"createTime"`
	} `json:"content"`
}

func GetLatestAPK(apples ...bool) *GetLatestAPKResp {
	source := "2"
	if len(apples) > 0 && apples[0] {
		source = "1"
	}
	result, _ := Execute[*GetLatestAPKResp](http.MethodGet, "/customer/sysVersionControl/single?source="+source, nil, nil)
	return result
}
