package route

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/ierr"
	"github.com/starudream/creative-apartment/internal/itask"
	"github.com/starudream/creative-apartment/internal/itime"
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
	ed = itime.New(ed.AddDate(0, 0, 1)).EndOfDay()

	resp := GetHouseDataResp{}

	dates := map[string]int{}
	for i, t := 0, sd; t.Before(ed.AddDate(0, 0, -1)); i, t = i+1, t.AddDate(0, 0, 1) {
		s := t.Format(config.DateFormat)
		resp.Dates = append(resp.Dates, s)
		dates[s] = i
	}

	for i := range resp.Datasets {
		resp.Datasets[i] = make([]float64, len(resp.Dates))
	}

	for t := sd; t.Before(ed); {
		// nolint
		x, m := itime.New(t).EndOfYear(), map[string][]float64{}
		if ed.After(x) {
			m, err = GetHouseStatsByYear(c, req.Phone, t.Format("2006"), t.Format("0102"), x.Format("0102"))
		} else {
			m, err = GetHouseStatsByYear(c, req.Phone, t.Format("2006"), t.Format("0102"), ed.Format("0102"))
		}
		if err != nil {
			c.AbortWithStatusJSON(ierr.New(0, err.Error()).Internal())
			return
		}
		for k, vs := range m {
			for i, v := range vs {
				resp.Datasets[i][dates[k]] = v
			}
		}
		t = x.Add(time.Nanosecond)
	}

	c.JSON(http.StatusOK, resp)
}

func GetHouseStatsByYear(ctx context.Context, phone, year, smd, emd string) (map[string][]float64, error) {
	data, mu := map[string][]float64{}, sync.Mutex{}
	task := itask.New(ctx, time.Minute)
	for _, suffix := range []int{1, 2, 3} {
		suffix := suffix
		task.Add(func(ctx context.Context) error {
			return ibolt.View(func(tx *ibolt.Tx) error {
				bucket := tx.Bucket([]byte(phone + "_house_stats_" + year + "_" + strconv.Itoa(suffix)))
				if bucket == nil {
					return nil
				}
				cur := bucket.Cursor()
				for k, v := cur.Seek([]byte(smd)); len(k) > 0 && bytes.Compare(k, []byte(emd)) <= 0; k, v = cur.Next() {
					sk := year + "-" + string(k[0:2]) + "-" + string(k[2:4])
					vd, err := decimal.NewFromString(string(v))
					if err != nil {
						return err
					}
					mu.Lock()
					if _, exist := data[sk]; !exist {
						data[sk] = make([]float64, 6)
					}
					switch string(k[5:6]) {
					case "a":
						data[sk][2*(suffix-1)], _ = vd.Float64()
					case "b":
						data[sk][2*suffix-1], _ = vd.Float64()
					}
					mu.Unlock()
				}
				return nil
			})
		})
	}
	return data, task.Run(3)
}
