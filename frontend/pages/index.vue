<template>
  <div>
    <div class="user">
      <i class="el-icon-switch-button" @click.prevent="logout"></i>
    </div>
  </div>
</template>

<script>
import { key } from "@/consts"

export default {
  methods: {
    logout() {
      this.$confirm("", "确认退出吗", {
        customClass: "logout-confirm",
        center:      true,
      }).then(() => {
        localStorage.clear()
        this.$store.commit("clear")
        this.$router.replace("/")
      })
    },
  },
}
</script>

<style>
.user {
  font-size: 30px;
  position: absolute;
  left: 0;
  bottom: 0;
  width: 100px;
  height: 100px;
  text-align: center;
  line-height: 100px;
  cursor: pointer;
}

.logout-confirm .el-message-box__btns {
  display: flex;
  justify-content: center;
  gap: 10px;
  flex-direction: row-reverse;
}
</style>
