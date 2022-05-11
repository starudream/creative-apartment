package cmd

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/starudream/creative-apartment/api"
	"github.com/starudream/creative-apartment/config"
	"github.com/starudream/creative-apartment/dist"
	"github.com/starudream/creative-apartment/internal/app"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/icron"
	"github.com/starudream/creative-apartment/internal/ierr"
	"github.com/starudream/creative-apartment/internal/igin"
	"github.com/starudream/creative-apartment/internal/ilog"
	"github.com/starudream/creative-apartment/internal/iroute"
	"github.com/starudream/creative-apartment/internal/iu"
	"github.com/starudream/creative-apartment/internal/json"
)

var rootCmd = &cobra.Command{
	Use:     config.AppName,
	Short:   config.AppName,
	Version: config.FULL_VERSION,
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
	igin.S().GET("/version", func(c *gin.Context) { c.String(http.StatusOK, config.FULL_VERSION) })
	g := igin.S().Group("/api/v1").Use(igin.Logger())
	{
		g.GET("electricity", iroute.Electricity)
		g.GET("water", iroute.Water)
	}
	igin.S().Use(igin.Serve("/", igin.StaticFile(dist.FS, ".", true)))
	igin.S().NoRoute(func(c *gin.Context) {
		rp := strings.TrimPrefix(c.Request.URL.Path, "/")
		if iu.SliceContains(dist.Files, rp) {
			c.FileFromFS(rp, http.FS(dist.FS))
		} else {
			c.AbortWithStatusJSON(ierr.NotFound())
		}
	})
	igin.S().NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(ierr.NotAllowed())
	})
	log.Info().Msgf("[http] load static files: %s", strings.Join(dist.Files, ", "))
	return igin.Run(":" + viper.GetString("port"))
}

func runCron(context.Context) error {
	if len(config.GetCustomers()) == 0 {
		log.Error().Msg("config not contain customers")
	} else {
		c := icron.New()
		icron.WrapError(c.AddFunc("0 0 6 * * *", runCronCustomers))
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
