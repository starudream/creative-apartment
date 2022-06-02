<template>
  <div></div>
</template>

<script>
import { key } from "@/consts"

export default {
  mounted() {
    const secret = localStorage.getItem(key.secret)
    if (!secret) {
      this.$router.replace("/login")
    }
  },
}
</script>
