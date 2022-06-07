import Vue from "vue"
import ECharts from "vue-echarts"
import { use } from "echarts/core"

import { CanvasRenderer } from "echarts/renderers"
import { BarChart, LineChart } from "echarts/charts"
import { GridComponent, TitleComponent, LegendComponent, TooltipComponent, ToolboxComponent } from "echarts/components"

use([
  CanvasRenderer,

  BarChart,
  LineChart,

  GridComponent,
  TitleComponent,
  LegendComponent,
  TooltipComponent,
  ToolboxComponent,
])

Vue.component("v-chart", ECharts)
