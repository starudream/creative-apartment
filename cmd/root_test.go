package cmd

import (
	"math/rand"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/starudream/creative-apartment/api"
	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/internal/itest"
	"github.com/starudream/creative-apartment/internal/itime"
)

func TestStoreHouseInfo(t *testing.T) {
	itest.Init()

	customer := &config.Customer{Phone: "13312341234"}

	var (
		e1 = float64(500)
		e2 = float64(500)
		st = itime.Now().BeginningOfMonth()
	)

	for i := 0; i < 30; i++ {
		info := &api.HouseInfoResp{
			Content: []*api.HouseInfoContent{
				{
					// CustomerPhone: customer.Phone,
					List: []*api.EquipmentInfo{
						{
							EquipmentType: 1,
							Surplus:       decimal.NewFromFloat(e1 / config.E1UnitPrice).StringFixed(2),
							SurplusAmount: decimal.NewFromFloat(e1).StringFixed(2),
							UnitPrice:     config.E1UnitPrice,
							LastReadTime:  st.Format(config.DateTimeFormat),
						},
						{
							EquipmentType: 2,
							Surplus:       decimal.NewFromFloat(e2 / config.E2UnitPrice).StringFixed(2),
							SurplusAmount: decimal.NewFromFloat(e2).StringFixed(2),
							UnitPrice:     config.E2UnitPrice,
							LastReadTime:  st.Format(config.DateTimeFormat),
						},
					},
				},
			},
		}
		storeHouseInfo(customer, info)

		e1 -= 20 * rand.Float64()
		e2 -= 20 * rand.Float64()
		st = st.AddDate(0, 0, 1)
	}
}
