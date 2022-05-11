package iroute

import (
	"net/http"
	"time"

	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/ierr"
	"github.com/starudream/creative-apartment/internal/itask"
	"github.com/starudream/creative-apartment/internal/ivalidator"
)

type GetHouseDataReq struct {
	Phone     string `json:"phone" validate:"required,min=1"`
	StartDate string `json:"startDate" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"endDate" validate:"required,datetime=2006-01-02,gtecsfield=StartDate"`
}

type GetHouseDataResp struct {
	// 年-月-日
	Dates []string `json:"dates"`
	// 0用电量 1电费
	// 2用水量 3水费
	// 4用气量 5气费
	Datasets [6][]float64 `json:"datasets"`
}

func GetHouseData(c *Context) {
	req := GetHouseDataReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(ierr.New(0, ivalidator.Translate(err)).Param())
		return
	}

	sd, _ := time.ParseInLocation(config.DateFormat, req.StartDate, time.Local)
	ed, _ := time.ParseInLocation(config.DateFormat, req.EndDate, time.Local)
	ed = ed.AddDate(0, 0, 1)

	resp := GetHouseDataResp{}

	dates := map[string]bool{}
	for t := sd; t.Before(ed); t = t.AddDate(0, 0, 1) {
		s := t.Format(config.DateFormat)
		resp.Dates = append(resp.Dates, s)
		dates[s] = true
	}

	task := itask.New()
	task.Add(ierr.WrapToFunc(ibolt.View(func(tx *ibolt.Tx) error {
		// cur := tx.Bucket([]byte("house_type1_" + req.Phone)).Cursor()
		// for k, v := cur.Seek([]byte(req.StartDate)); k != nil && dates[string(k)]; k, v = cur.Next() {
		// 	fmt.Println(string(v))
		// }
		return nil
	})))
	err = task.Run(c, time.Minute)
	if err != nil {
		c.AbortWithStatusJSON(ierr.New(0, err.Error()).Internal())
		return
	}

	c.JSON(http.StatusOK, resp)
}
