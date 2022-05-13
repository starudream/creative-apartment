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
	"github.com/shopspring/decimal"
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
	"github.com/starudream/creative-apartment/internal/iu"
	"github.com/starudream/creative-apartment/internal/json"
	"github.com/starudream/creative-apartment/route"
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
	am := igin.Auth(viper.GetString("secret"))

	igin.S().GET("/version", func(c *gin.Context) { c.String(http.StatusOK, config.FULL_VERSION) })

	g := igin.S().Group("/api/v1").Use(igin.Logger(), am)
	{
		g.POST("customers", route.ListCustomers)
		g.POST("house/data", route.GetHouseData)
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
	log.Info().Msgf("[cron] done")
}

func storeHouseInfo(customer *config.Customer, info *api.HouseInfoResp) {
	_ = ibolt.Update(func(tx *ibolt.Tx) error {
		for _, house := range info.Content {
			for _, data := range house.List {
				t0, err := time.ParseInLocation(config.DateTimeFormat, data.LastReadTime, time.Local)
				if !ilog.WrapError(err) {
					continue
				}
				surplus, err := decimal.NewFromString(data.Surplus)
				if !ilog.WrapError(err) {
					continue
				}
				surplusAmount, err := decimal.NewFromString(data.SurplusAmount)
				if !ilog.WrapError(err) {
					continue
				}

				s0 := customer.Phone + "_house_data_" + t0.Format("2006") + "_" + strconv.Itoa(data.EquipmentType)
				bucket0, err := tx.CreateBucketIfNotExists([]byte(s0))
				if !ilog.WrapError(err) {
					continue
				}

				v0 := api.SimpleEquipmentInfo{}
				v0.Surplus, _ = surplus.Float64()
				v0.SurplusAmount, _ = surplusAmount.Float64()
				v0.UnitPrice = data.UnitPrice
				v0.LastReadTime = t0

				if !ilog.WrapError(bucket0.Put([]byte(t0.Format("0102")), json.MustMarshal(v0)), "store") {
					continue
				}

				log.Info().Str("phone", customer.Phone).Int("type", data.EquipmentType).Time("time", t0).Msgf("store house info success")

				t1 := t0.AddDate(0, 0, -1)

				bs := bucket0.Get([]byte(t1.Format("0102")))
				if len(bs) == 0 {
					continue
				}

				v1, err := json.UnmarshalTo[api.SimpleEquipmentInfo](bs)
				if !ilog.WrapError(err) {
					continue
				}

				s1 := customer.Phone + "_house_stats_" + t1.Format("2006") + "_" + strconv.Itoa(data.EquipmentType)
				bucket1, err := tx.CreateBucketIfNotExists([]byte(s1))
				if !ilog.WrapError(err) {
					continue
				}

				a := v1.Surplus - v0.Surplus
				for a < 0 {
					a += config.RechargeAmount * v1.UnitPrice
				}
				if !ilog.WrapError(bucket1.Put([]byte(t0.Format("0102")+"_a"), []byte(decimal.NewFromFloat(a).StringFixed(2))), "store") {
					continue
				}

				b := v1.SurplusAmount - v0.SurplusAmount
				for b < 0 {
					b += config.RechargeAmount
				}
				if !ilog.WrapError(bucket1.Put([]byte(t0.Format("0102")+"_b"), []byte(decimal.NewFromFloat(b).StringFixed(2))), "store") {
					continue
				}
			}
		}

		ilog.WrapError(tx.Bucket([]byte("customer")).Put([]byte(customer.Phone), nil), "store")
		return nil
	})
}
