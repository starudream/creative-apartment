<template>
  <div class="version">{{ this.$store.state.version }}</div>
</template>

<style>
.version {
  font-size: 14px;
  position: fixed;
  right: 0;
  bottom: 0;
  width: 200px;
  height: 100px;
  text-align: center;
  line-height: 100px;
}
</style>
