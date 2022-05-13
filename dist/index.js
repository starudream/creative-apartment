// noinspection JSUnresolvedVariable

let _loading = $("#loading")
let _login = $("#login")
let _chart = $("#chart")
let chart

let options = {
  locale:      "zh-CN",
  responsive:  true,
  interaction: {
    intersect: false,
    mode:      "nearest",
    axis:      "x",
  },
  plugins:     {
    legend:  {
      display:  true,
      position: "top",
      align:    "center",
    },
    tooltip: {
      callbacks: {
        footer: function (items) {
          let sum = 0
          items.forEach(function (item) {
            if (item.dataset.label && item.dataset.label.indexOf("费") > -1) {
              sum += item.parsed.y
            }
          })
          return "费用总和 " + sum.toFixed(2) + " 元"
        },
      },
    },
    title:   {
      display: true,
      text:    "费用使用情况",
    },
  },
}

let params = {type: "line", options: options}

let Update = function (data) {
  console.log(data)
  chart.data = data
  chart.data.datasets.forEach(function (dataset) {
    if (!dataset.label) return
    if (dataset.label) {
      if (dataset.label.indexOf("电量") > -1) {
        dataset.borderColor = "rgb(217, 136, 128)"
      } else if (dataset.label.indexOf("电费") > -1) {
        dataset.borderColor = "rgb(203, 67, 53)"
      } else if (dataset.label.indexOf("水量") > -1) {
        dataset.borderColor = "rgb(127, 179, 213)"
      } else if (dataset.label.indexOf("水费") > -1) {
        dataset.borderColor = "rgb(46, 134, 193)"
      } else if (dataset.label.indexOf("气量") > -1) {
        dataset.borderColor = "rgb(118, 215, 196)"
      } else if (dataset.label.indexOf("气费") > -1) {
        dataset.borderColor = "rgb(19, 141, 117)"
      }
    }
  })
  chart.update()
}

let POST = function (url, data, callback) {
  $.ajax({
    type:        "POST",
    url:         url,
    data:        JSON.stringify(data),
    contentType: "application/json",
    timeout:     60 * 1000,
    complete:    function (xhr) {
      const resp = JSON.parse(xhr.responseText)
      if (xhr.status === 200) {
        callback(resp)
      } else {
        Swal.fire({title: "Error " + xhr.status, text: resp.message, icon: "error"})
      }
    },
  })
}

let ListCustomers = function () {
  POST("/api/v1/customers", {}, function (resp) {
    if (resp.customers && resp.customers.length > 0) {
      // TODO: TEST
      // GetHouseData("13312341234", dayjs().add(-10, "day"), dayjs())
      GetHouseData(resp.customers[0].phone, dayjs().add(-30, "day"), dayjs())
    } else {
      Swal.fire({title: "Error", text: "没有用户信息", icon: "error"})
    }
  })
}

let GetHouseData = function (phone, startTime, endTime) {
  const data = {
    phone:     phone,
    startDate: startTime.format("YYYY-MM-DD"),
    endDate:   endTime.format("YYYY-MM-DD"),
  }
  POST("/api/v1/house/data", data, function (resp) {
    Update({
      labels:   resp.dates,
      datasets: [
        {label: "用电量（度）", data: resp.datasets[0]},
        {label: "电费（元）", data: resp.datasets[1]},
        {label: "用水量（立方米）", data: resp.datasets[2]},
        {label: "水费（元）", data: resp.datasets[3]},
        {label: "用气量（立方米）", data: resp.datasets[4]},
        {label: "气费（元）", data: resp.datasets[5]},
      ],
    })
  })
}

$(function () {
  chart = new Chart(_chart, params)
  _loading.hide()
  _login.show()
  // _chart.show()
})
