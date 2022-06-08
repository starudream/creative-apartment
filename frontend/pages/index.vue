<template>
  <div>
    <el-row style="margin-top: 50px;">
      <el-col :md="{span:8,offset:8}">
        <div style="max-width: 560px; margin: 0 auto;">
          <div style="display: inline-block; width: 200px; margin: 0 auto 5px;">
            <el-select v-model="phone" filterable placeholder="请选择">
              <el-option v-for="customer in customers" :key="customer.phone" :label="customer.phone" :value="customer.phone"/>
            </el-select>
          </div>
          <div style="display: inline-block; width: 350px; margin: 0 auto;">
            <el-date-picker v-model="date" :picker-options="dateOption" clearable end-placeholder="结束日期" range-separator="至" size="large" start-placeholder="开始日期" type="daterange" @blur="changeDate"/>
          </div>
        </div>
      </el-col>
    </el-row>
    <el-row>
      <el-col :lg="{span:16,offset:4}" style="margin-top: 50px;">
        <div style="height: 600px;">
          <client-only>
            <v-chart :option="chartOption"/>
          </client-only>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  data() {
    return {
      customers:   [],
      phone:       "",
      yesterday:   this.$dayjs().startOf("day").add(-1, "day"),
      day1:        this.$dayjs().startOf("day"),
      week1:       this.$dayjs().startOf("week").add(1, "day"),
      month1:      this.$dayjs().startOf("month"),
      date:        [],
      dateOption:  {
        firstDayOfWeek: 1,
        disabledDate:   (date) => {
          return date >= this.day1.toDate()
        },
        shortcuts:      [
          {
            text:    "本周",
            onClick: (picker) => {
              picker.$emit("pick", [this.week1.toDate(), this.week1.add(6, "day").toDate()])
            },
          },
          {
            text:    "上周",
            onClick: (picker) => {
              picker.$emit("pick", [this.week1.add(-7, "day").toDate(), this.week1.add(-1, "day").toDate()])
            },
          },
          {
            text:    "近一周",
            onClick: (picker) => {
              picker.$emit("pick", [this.yesterday.add(-6, "day"), this.yesterday])
            },
          },
          {
            text:    "本月",
            onClick: (picker) => {
              picker.$emit("pick", [this.month1.toDate(), this.month1.add(1, "month").add(-1, "day").toDate()])
            },
          },
          {
            text:    "上月",
            onClick: (picker) => {
              picker.$emit("pick", [this.month1.add(-1, "month").toDate(), this.month1.add(-1, "day").toDate()])
            },
          },
          {
            text:    "近30天",
            onClick: (picker) => {
              picker.$emit("pick", [this.yesterday.add(-30, "day").toDate(), this.yesterday])
            },
          },
          {
            text:    "近90天",
            onClick: (picker) => {
              picker.$emit("pick", [this.yesterday.add(-90, "day").toDate(), this.yesterday])
            },
          },
        ],
      },
      chartOption: {
        grid:    {
          left:         "5%",
          right:        "5%",
          top:          "120",
          bottom:       "0",
          containLabel: true,
        },
        title:   {
          text: "用量统计",
          left: "center",
          top:  "0",
        },
        legend:  {
          data: ["电量", "电费", "水量", "水费"],
          left: "center",
          top:  "35",
        },
        toolbox: {
          left:    "center",
          top:     "70",
          feature: {
            dataView:    {
              readOnly: false,
            },
            magicType:   {
              type: ["line", "bar"],
            },
            saveAsImage: {
              pixelRatio: 2,
            },
            restore:     {},
          },
        },
        tooltip: {
          trigger: "axis",
        },
        xAxis:   {
          type:        "category",
          boundaryGap: false,
          data:        [],
        },
        yAxis:   {
          type: "value",
        },
        series:  [
          {
            type: "line",
            name: "电量",
            data: [],
          },
          {
            type: "line",
            name: "电费",
            data: [],
          },
          {
            type: "line",
            name: "水量",
            data: [],
          },
          {
            type: "line",
            name: "水费",
            data: [],
          },
        ],
      },
    }
  },
  methods: {
    initDate() {
      this.date = [
        this.yesterday.add(-6, "day").toDate(),
        this.yesterday.toDate(),
      ]
    },
    changeDate() {
      if (this.date.length !== 2) {
        this.initDate()
      }
      if (!this.phone) {
        this.$message.error("没有选择租户")
        return
      }
      this.getHouseData()
    },
    async listCustomers() {
      this.customers = []
      await this.$axios.$post(this.$store.state.baseURL + "/api/v1/customers", {}, {
        responseType: "json",
      }).then((v) => {
        if (v.customers && v.customers.length > 0) {
          v.customers.forEach((customer, index) => {
            if (index === 0) {
              this.phone = customer.phone
            }
            this.customers.push({phone: customer.phone})
          })
        }
      }).catch(() => {
        this.$message.error("接口错误")
      })
    },
    async getHouseData() {
      this.$axios.$post(this.$store.state.baseURL + "/api/v1/house/data", {
        phone:     this.phone,
        startDate: this.$dayjs(this.date[0]).format("YYYY-MM-DD"),
        endDate:   this.$dayjs(this.date[1]).format("YYYY-MM-DD"),
      }, {
        responseType: "json",
      }).then((v) => {
        if (v.dates) {
          this.chartOption.xAxis.data = v.dates
          for (let i = 0; i < this.chartOption.series.length; i++) {
            this.chartOption.series[i].data = v.datasets[i]
          }
        }
      }).catch(() => {
        this.$message.error("接口错误")
      })
    },
  },
  async mounted() {
    this.initDate()
    await this.listCustomers()
    if (this.phone) {
      await this.getHouseData()
    }
  },
}
</script>
