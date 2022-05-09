package cmd

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/starudream/creative-apartment/api"
	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/icron"
	"github.com/starudream/creative-apartment/internal/ierr"
	"github.com/starudream/creative-apartment/internal/ilog"
	"github.com/starudream/creative-apartment/internal/json"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run daily tasks to store data",
	Run: func(cmd *cobra.Command, args []string) {
		if len(config.GetCustomers()) == 0 {
			ierr.CheckErr("config not contain customers")
		}
		c := icron.New()
		c.AddFunc("0 0 6 * * *", runCron)
		// c.AddFunc("0 * * * * *", runCron)
		c.Run()
	},
}

func runCron() {
	for i := 0; i < len(config.GetCustomers()); i++ {
		customer := config.GetCustomers()[i]
		info := api.GetHouseInfo(customer.GetToken())
		if info == nil {
			log.Error().Str("phone", customer.Phone).Msgf("get house info error")
		} else {
			storeHouseInfo(customer, info)
		}
	}
}

func storeHouseInfo(customer *config.Customer, info *api.HouseInfoResp) {
	_ = ibolt.Update(func(tx *ibolt.Tx) error {
		for _, house := range info.Content {
			for _, data := range house.List {
				t, err := time.ParseInLocation(config.DateTimeFormat, data.LastReadTime, time.Local)
				if !ilog.WrapError(err) {
					continue
				}
				name := "house_type" + strconv.Itoa(data.EquipmentType) + "_" + customer.Phone
				bucket, err := tx.CreateBucketIfNotExists([]byte(name))
				if !ilog.WrapError(err) {
					continue
				}
				if ilog.WrapError(bucket.Put([]byte(t.Format(config.DateFormat)), json.MustMarshal(data)), "store") {
					log.Debug().Str("phone", customer.Phone).Int("type", data.EquipmentType).Str("time", t.Format(config.DateFormat)).Msgf("store house info success")
				}
			}
		}
		return nil
	})
}
