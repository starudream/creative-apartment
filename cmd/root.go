package cmd

import (
	"context"
	"fmt"
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
	"github.com/starudream/creative-apartment/internal/app"
	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/ibot"
	"github.com/starudream/creative-apartment/internal/icron"
	"github.com/starudream/creative-apartment/internal/ierr"
	"github.com/starudream/creative-apartment/internal/igin"
	"github.com/starudream/creative-apartment/internal/ilog"
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
	secret := viper.GetString("secret")

	igin.S().GET("/version", func(c *gin.Context) { c.String(http.StatusOK, config.FULL_VERSION) })

	g := igin.S().Group("/api/v1").Use(igin.Logger(), igin.Auth(secret))
	{
		g.POST("customers", route.ListCustomers)
		g.POST("house/data", route.GetHouseData)
	}

	igin.S().Use(igin.Serve("/", igin.LocalFileSystem("dist"), func(c *gin.Context) { c.AbortWithStatusJSON(ierr.NotFound()) }))
	igin.S().NoMethod(func(c *gin.Context) { c.AbortWithStatusJSON(ierr.NotAllowed()) })

	return igin.Run(":" + viper.GetString("port"))
}

func runCron(context.Context) error {
	time.Sleep(time.Second)

	if len(config.GetCustomers()) == 0 {
		log.Error().Msg("config not contain customers")
	} else {
		c := icron.New()
		icron.WrapError(c.AddFunc("0 0 06 * * *", runCronCustomers))
		icron.WrapError(c.AddFunc("0 0 14 * * *", runCronCustomers))
		icron.WrapError(c.AddFunc("0 0 22 * * *", runCronCustomers))
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
			var vs0, vs1, datetime = make([]api.SimpleEquipmentInfo, 3), make([]api.SimpleEquipmentInfo, 3), ""
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

				if len(bucket0.Get([]byte(t0.Format("0102")))) > 0 {
					log.Debug().Str("phone", customer.Phone).Int("type", data.EquipmentType).Time("time", t0).Msgf("already stored house info")
					continue
				}

				v0 := api.SimpleEquipmentInfo{}
				v0.Surplus, _ = surplus.Float64()
				v0.SurplusAmount, _ = surplusAmount.Float64()
				v0.UnitPrice = data.UnitPrice
				v0.LastReadTime = t0

				vs0[data.EquipmentType-1] = v0

				if !ilog.WrapError(bucket0.Put([]byte(t0.Format("0102")), json.MustMarshal(v0)), "store") {
					continue
				}

				log.Debug().Str("phone", customer.Phone).Int("type", data.EquipmentType).Time("time", t0).Msgf("store house info success")

				t1 := t0.AddDate(0, 0, -1)

				if datetime == "" {
					datetime = t1.Format(config.DateFormat)
				}

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

				vs1[data.EquipmentType-1] = api.SimpleEquipmentInfo{Surplus: a, SurplusAmount: b, UnitPrice: v1.UnitPrice}
			}
			if house.CustomerPhone != "" && viper.GetString("dingtalk.token") != "" {
				sendHouseInfoMessage(house.CustomerPhone, datetime, vs0, vs1)
			}
		}
		ilog.WrapError(tx.Bucket([]byte("customer")).Put([]byte(customer.Phone), nil), "store")
		return nil
	})
}

func sendHouseInfoMessage(phone, datetime string, vs0, vs1 []api.SimpleEquipmentInfo) {
	msg := strings.Builder{}
	for i := 0; i < 3; i++ {
		tag := ""
		switch i {
		case 0:
			tag = "电"
		case 1:
			tag = "水"
		case 2:
			tag = "气"
		}
		x := false
		if v := vs1[i]; v.UnitPrice > 0 {
			msg.WriteString(tag)
			msg.WriteString("费")
			msg.WriteString(decimal.NewFromFloat(v.SurplusAmount).StringFixed(2))
			x = true
		}
		if v := vs0[i]; v.UnitPrice > 0 {
			if x {
				msg.WriteString("，剩余")
			} else {
				msg.WriteString("剩余")
				msg.WriteString(tag)
				msg.WriteString("费")
			}
			msg.WriteString(decimal.NewFromFloat(v.SurplusAmount).StringFixed(2))
			x = true
		}
		if x {
			msg.WriteString("\n")
		}
	}
	if msg.Len() > 0 {
		s := fmt.Sprintf("【%s】【%s】\n%s", phone, datetime, msg.String())
		ierr.CheckErr(ibot.Dingtalk.SendMessage(s))
	}
}
