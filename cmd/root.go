package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/starudream/creative-apartment/api"
	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/dist"
	"github.com/starudream/creative-apartment/internal/app"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/icron"
	"github.com/starudream/creative-apartment/internal/igin"
	"github.com/starudream/creative-apartment/internal/ilog"
	"github.com/starudream/creative-apartment/internal/json"
)

var rootCmd = &cobra.Command{
	Use:     config.AppName,
	Short:   config.AppName,
	Version: fmt.Sprintf("%s (%s)", config.VERSION, config.BIDTIME),
	Run: func(cmd *cobra.Command, args []string) {
		app.Add(initRouter)
		app.Add(runCron)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initRouter(context.Context) error {
	igin.S().StaticFS("/", http.FS(dist.FS))
	return igin.Run(":" + viper.GetString("port"))
}

func runCron(context.Context) error {
	if len(config.GetCustomers()) == 0 {
		log.Error().Msg("config not contain customers")
	} else {
		c := icron.New()
		c.AddFunc("0 0 6 * * *", runCronCustomers)
		// c.AddFunc("* * * * * *", runCronCustomers)
		c.Run()
	}
	return nil
}

func runCronCustomers() {
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
