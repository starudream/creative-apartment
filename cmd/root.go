package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
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
	"github.com/starudream/creative-apartment/internal/iscript"
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
		app.Add(iscript.FixHouseStatsOffset)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initRouter(context.Context) error {
	auth := igin.Auth(viper.GetString("secret"), authLimiter())

	igin.S().Use(igin.CORS())

	igin.S().GET("/version", version)
	igin.S().GET("/verifySecret", auth)

	g := igin.S().Group("/api/v1").Use(igin.Logger(), auth)
	{
		g.POST("customers", route.ListCustomers)
		g.POST("house/data", route.GetHouseData)
	}

	igin.S().Use(igin.Static("/", igin.LocalFileSystem("dist"), func(c *gin.Context) { c.AbortWithStatusJSON(ierr.NotFound()) }))
	igin.S().NoMethod(func(c *gin.Context) { c.AbortWithStatusJSON(ierr.NotAllowed()) })

	return igin.Run(":" + viper.GetString("port"))
}

func authLimiter() igin.AuthLimitFunc {
	max := 5
	cc := cache.New(12*time.Hour, time.Minute)
	return func(c *gin.Context) bool {
		k := c.ClientIP()
		x, _ := cc.IncrementInt(k, 1)
		if x <= 0 {
			x = 1
			cc.SetDefault(k, x)
		}
		_, d, _ := cc.GetWithExpiration(k)
		c.Header("X-RateLimit-Limit", strconv.Itoa(max))
		if max-x < 0 {
			c.Header("X-RateLimit-Remaining", "0")
		} else {
			c.Header("X-RateLimit-Remaining", strconv.Itoa(max-x))
		}
		c.Header("X-RateLimit-Reset", strconv.Itoa(int(d.Unix())))
		if x > max {
			c.AbortWithStatusJSON(ierr.Frequent())
		}
		return x <= max
	}
}

func version(c *gin.Context) {
	md := ierr.MD{"version": config.VERSION, "bidtime": config.BIDTIME}
	c.JSON(ierr.New().SetMetadata(md).OK())
}

func runCron(context.Context) error {
	time.Sleep(time.Second)

	if len(config.GetCustomers()) == 0 {
		log.Error().Msg("config not contain customers")
	} else {
		go func() {
			if viper.GetBool("startup") {
				// 启动时运行一次
				runCronCustomers()
			}
		}()
		go func() {
			c := icron.New()
			// 抄表之前，如果有充值记录下充值后的数据
			icron.WrapError(c.AddFunc("0 0 03 * * *", runCronCustomers))
			// 抄表之后，记录下前一天的消耗量
			icron.WrapError(c.AddFunc("0 0 04 * * *", runCronCustomers))
			c.Run()
		}()
	}
	return nil
}

func runCronCustomers() {
	for i := 0; i < len(config.GetCustomers()); i++ {
		customer := config.GetCustomers()[i]
		info := api.GetHouseInfo(customer.GetToken())
		if info == nil {
			log.Error().Str("phone", customer.Phone).Msgf("get house info error")
			sendMessage("获取信息失败，请检查配置")
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

				v0 := api.SimpleEquipmentInfo{}
				v0.Surplus, _ = surplus.Float64()
				v0.SurplusAmount, _ = surplusAmount.Float64()
				v0.UnitPrice = data.UnitPrice
				v0.LastReadTime = t0

				vs0[data.EquipmentType-1] = v0

				if bs0 := bucket0.Get([]byte(t0.Format("0102"))); len(bs0) > 0 {
					vt, err := json.UnmarshalTo[api.SimpleEquipmentInfo](bs0)
					if !ilog.WrapError(err) {
						continue
					}
					if vt.Surplus != v0.Surplus || vt.SurplusAmount != v0.SurplusAmount {
						v0.LastRecord = &vt
					} else {
						log.Debug().Str("phone", customer.Phone).Int("type", data.EquipmentType).Time("time", t0).Msgf("already stored house info and no change")
						continue
					}
				}

				if !ilog.WrapError(bucket0.Put([]byte(t0.Format("0102")), json.MustMarshal(v0)), "store") {
					continue
				}

				log.Debug().Str("phone", customer.Phone).Int("type", data.EquipmentType).Time("time", t0).Msgf("store house info success")

				if v0.LastRecord != nil {
					continue
				}

				t1 := t0.AddDate(0, 0, -1)

				if datetime == "" {
					datetime = t1.Format(config.DateFormat)
				}

				bs1 := bucket0.Get([]byte(t1.Format("0102")))
				if len(bs1) == 0 {
					continue
				}

				v1, err := json.UnmarshalTo[api.SimpleEquipmentInfo](bs1)
				if !ilog.WrapError(err) {
					continue
				}

				s1 := customer.Phone + "_house_stats_" + t1.Format("2006") + "_" + strconv.Itoa(data.EquipmentType)
				bucket1, err := tx.CreateBucketIfNotExists([]byte(s1))
				if !ilog.WrapError(err) {
					continue
				}

				a := v1.Surplus - v0.Surplus
				if !ilog.WrapError(bucket1.Put([]byte(t1.Format("0102")+"_a"), []byte(decimal.NewFromFloat(a).StringFixed(2))), "store") {
					continue
				}

				b := v1.SurplusAmount - v0.SurplusAmount
				if !ilog.WrapError(bucket1.Put([]byte(t1.Format("0102")+"_b"), []byte(decimal.NewFromFloat(b).StringFixed(2))), "store") {
					continue
				}

				vs1[data.EquipmentType-1] = api.SimpleEquipmentInfo{Surplus: a, SurplusAmount: b, UnitPrice: v1.UnitPrice}
			}
			sendHouseInfoMessage(house.CustomerPhone, datetime, vs0, vs1)
		}
		ilog.WrapError(tx.Bucket([]byte("customer")).Put([]byte(customer.Phone), nil), "store")
		return nil
	})
}

func sendMessage(s string) {
	if viper.GetString("dingtalk.token") == "" {
		return
	}
	ierr.CheckErr(ibot.Dingtalk.SendMessage(s))
}

func sendHouseInfoMessage(phone, datetime string, vs0, vs1 []api.SimpleEquipmentInfo) {
	if phone == "" || datetime == "" || viper.GetString("dingtalk.token") == "" {
		return
	}
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
