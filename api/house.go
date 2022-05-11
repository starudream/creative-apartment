package api

import (
	"net/http"
	"time"
)

type HouseInfoResp struct {
	CommonResp
	Content []*HouseInfoContent `json:"content"`
}

type HouseInfoContent struct {
	Id                    string           `json:"id"`
	OrderContractId       string           `json:"orderContractId"`
	OrderContractDetailId string           `json:"orderContractDetailId"`
	HouseName             string           `json:"houseName"`
	HouseCode             string           `json:"houseCode"`
	CustomerInfoId        string           `json:"customerInfoId"`
	CustomerName          string           `json:"customerName"`
	CustomerPhone         string           `json:"customerPhone"`
	List                  []*EquipmentInfo `json:"list"`
}

type EquipmentInfo struct {
	Id               string  `json:"id"`
	EquipmentType    int     `json:"equipmentType"`
	HouseCode        string  `json:"houseCode"`
	MeterOnOff       int     `json:"meterOnOff"`
	NoPowerOff       int     `json:"noPowerOff"`
	MeterAddr        string  `json:"meterAddr"`
	MeterInstallAddr string  `json:"meterInstallAddr"`
	Surplus          string  `json:"surplus"`      // 余量
	SurplusAmount    string  `json:"surplusAmout"` // 余额
	UnitPrice        float64 `json:"unitPrice"`    // 单价
	OperateTime      string  `json:"operateTime"`
	LastReadTime     string  `json:"lastReadTime"` // 上次读表时间
}

type SimpleEquipmentInfo struct {
	Surplus       float64   `json:"surplus"`       // 余量
	SurplusAmount float64   `json:"surplusAmount"` // 余额
	UnitPrice     float64   `json:"unitPrice"`     // 单价
	LastReadTime  time.Time `json:"lastReadTime"`  // 上次读表时间
}

func GetHouseInfo(accessToken string) *HouseInfoResp {
	resp, _ := Execute[*HouseInfoResp](http.MethodGet, "/customer/hydropower/houseInfoByPage", HAuth(accessToken), nil)
	return resp
}
