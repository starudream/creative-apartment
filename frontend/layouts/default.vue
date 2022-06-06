<template>
  <div class="container">
    <Nuxt/>
    <GitHub/>
    <Version/>
  </div>
</template>

<script>
import { key } from "@/consts"

export default {
  mounted() {
    const notLogin = this.$route.path !== "/login"

    let baseURL = localStorage.getItem(key.baseURL)
    if (!baseURL) {
      baseURL = location.origin
    }
    this.$store.commit("setBaseURL", baseURL)

    const secret = localStorage.getItem(key.secret)
    if (!secret) {
      if (notLogin) {
        this.$router.replace("/login")
      }
      return
    }
    this.$store.commit("setSecret", secret)

    this.$axios.$get(baseURL + "/version", {
      responseType: "json",
    }).then((x) => {
      if (x.code === 200) {
        this.$store.commit("setVersion", x.metadata.version)
      }
    }).catch(() => {
      if (notLogin) {
        this.$router.replace("/login")
      }
    })
  },
}
</script>
