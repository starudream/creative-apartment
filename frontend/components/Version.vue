<template>
  <div class="version">{{ this.$store.state.version }}</div>
</template>

<script>
export default {}
</script>

<style>
.version {
  font-size: 14px;
  position: absolute;
  right: 0;
  bottom: 0;
  width: 200px;
  height: 100px;
  text-align: center;
  line-height: 100px;
}
</style>
