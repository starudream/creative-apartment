<template>
  <el-container v-if="done">
    <el-main>
      <Nuxt/>
    </el-main>
    <el-footer>
      <Version/>
      <div v-show="this.$store.state.isLogin" class="user">
        <i class="el-icon-switch-button" @click.prevent="logout"></i>
      </div>
    </el-footer>
    <GitHub/>
  </el-container>
</template>

<script>
import { key } from "@/consts"

export default {
  data() {
    return {
      done: false,
    }
  },
  async mounted() {
    const notLogin = !this.$route.path.startsWith("/login")

    let baseURL = localStorage.getItem(key.baseURL)
    if (!baseURL) {
      baseURL = location.origin
    }
    this.$store.commit("setBaseURL", baseURL)

    const secret = localStorage.getItem(key.secret)
    if (!secret) {
      if (notLogin) {
        await this.$router.replace("/login")
      }
    } else {
      this.$store.commit("setSecret", secret)

      await this.$axios.$get(baseURL + "/version", {
        responseType: "json",
      }).then((x) => {
        if (x.code === 200) {
          this.$store.commit("setLogin", true)
          this.$store.commit("setVersion", x.metadata.version)
          this.$axios.setToken(secret)
          if (!notLogin) {
            location.href = process.env.path
          }
        }
      }).catch(() => {
        if (notLogin) {
          this.$router.replace("/login")
        }
      })
    }

    this.done = true
  },
  methods: {
    logout() {
      this.$confirm("", "确认退出吗", {
        customClass: "logout-confirm",
        center:      true,
      }).then(() => {
        localStorage.clear()
        this.$store.commit("clear")
        location.href = process.env.path
      })
    },
  },
}
</script>

<style>
body {
  padding: 0;
  margin: 0 0 100px;
}

.user {
  font-size: 30px;
  position: fixed;
  left: 0;
  bottom: 0;
  width: 100px;
  height: 100px;
  text-align: center;
  line-height: 100px;
  cursor: pointer;
  overflow: auto;
}

.logout-confirm .el-message-box__btns {
  display: flex;
  justify-content: center;
  gap: 10px;
  flex-direction: row-reverse;
}
</style>
