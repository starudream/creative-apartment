export default {
  // Target: https://go.nuxtjs.dev/config-target
  target: "static",

  // Env: https://nuxtjs.org/docs/configuration-glossary/configuration-env
  env: {
    path: "/",
  },

  // Router: https://nuxtjs.org/docs/configuration-glossary/configuration-router
  router: {
    base: "/",
  },

  // Global page headers: https://go.nuxtjs.dev/config-head
  head: {
    title:     "城投宽庭",
    htmlAttrs: {
      lang: "zh-CN",
    },
    meta:      [
      {charset: "utf-8"},
      {name: "viewport", content: "width=device-width, initial-scale=1"},
      {hid: "description", name: "description", content: ""},
      {name: "format-detection", content: "telephone=no"},
    ],
    link:      [
      {rel: "icon", type: "image/x-icon", href: "/favicon.ico"},
    ],
  },

  // Global CSS: https://go.nuxtjs.dev/config-css
  css: [
    "element-ui/lib/theme-chalk/index.css",
  ],

  // Plugins to run before rendering page: https://go.nuxtjs.dev/config-plugins
  plugins: [
    // https://element.eleme.io
    {src: "@/plugins/element-ui", ssr: true},
    // https://github.com/ecomfe/vue-echarts
    {src: "@/plugins/vue-echarts", ssr: false},
  ],

  // Auto import components: https://go.nuxtjs.dev/config-components
  components: true,

  // Modules for dev and build (recommended): https://go.nuxtjs.dev/config-modules
  buildModules: [
    // https://github.com/ecomfe/vue-echarts
    "@nuxtjs/composition-api/module",
  ],

  // Modules: https://go.nuxtjs.dev/config-modules
  modules: [
    // https://go.nuxtjs.dev/axios
    "@nuxtjs/axios",
    // https://github.com/nuxt-community/dayjs-module
    "@nuxtjs/dayjs",
  ],

  // Axios module configuration: https://go.nuxtjs.dev/config-axios
  axios: {
    // Workaround to avoid enforcing hard-coded localhost:3000: https://github.com/nuxt-community/axios-module/issues/308
    baseURL:  "/",
    progress: true,
  },

  // Day.js: https://github.com/nuxt-community/dayjs-module
  dayjs: {
    defaultLocale:   "zh",
    defaultTimeZone: "Asia/Shanghai",
    plugins:         [
      "utc",
      "duration",
      "timezone",
      "relativeTime",
    ],
  },

  // Loading: https://nuxtjs.org/docs/features/loading/
  loading: {
    color: "blue",
  },

  // Build Configuration: https://go.nuxtjs.dev/config-build
  build: {
    transpile: [/^element-ui/],
  },
}
