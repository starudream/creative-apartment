<template>
  <div class="login">
    <div class="form">
      <el-form inline size="medium">
        <el-form-item>
          <el-input v-model="secret" autofocus clearable placeholder="请输入秘钥" show-password type="password"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click.prevent="login">登录</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="setting" @click.prevent="showSetting">
      <i class="el-icon-setting"></i>
    </div>
  </div>
</template>

<script>
import { key } from "@/consts"

export default {
  computed: {
    baseURL: {
      get() {
        return this.$store.state.baseURL
      },
      set(v) {
        localStorage.setItem(key.baseURL, v)
        this.$store.commit("setBaseURL", v)
      },
    },
    secret:  {
      get() {
        return this.$store.state.secret
      },
      set(v) {
        localStorage.setItem(key.secret, v)
        this.$store.commit("setSecret", v)
      },
    },
  },
  methods:  {
    login() {
      this.$axios.$get(this.baseURL + "/verifySecret", {
        params:       {secret: this.secret},
        responseType: "json",
      }).then((v) => {
        if (v.code === 200) {
          this.$message.success("秘钥验证成功")
          this.$router.replace("/")
        } else {
          this.$message.error("后端地址设置错误")
        }
      }).catch(() => {
        this.$message.error("秘钥错误")
      })
    },
    showSetting() {
      this.$prompt("", "自定义后端地址", {
        customClass:       "setting-prompt",
        closeOnClickModal: false,
        inputPlaceholder:  this.baseURL,
        inputValue:        this.baseURL,
        inputType:         "text",
        center:            true,
      }).then(({value}) => {
        if (value) {
          if (!value.startsWith("http")) {
            // noinspection HttpUrlsUsage
            value = "http://" + value
          }
          if (value.endsWith("/")) {
            value = value.substring(0, value.length - 1)
          }
          this.baseURL = value
        }
      })
    },
  },
}
</script>

<style>
body {
  background-color: whitesmoke;
}

.form {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.setting {
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

.setting-prompt .el-message-box__btns {
  display: flex;
  justify-content: center;
  gap: 10px;
  flex-direction: row-reverse;
}
</style>
