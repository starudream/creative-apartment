<template>
  <el-row>
    <el-col>
      <el-result :title="error.statusCode === 404 ? '页面未找到' : '未知错误'" icon="error">
        <template slot="extra">
          <el-button size="medium" type="primary" @click="$router.go(-1)">返回</el-button>
        </template>
      </el-result>
    </el-col>
  </el-row>
</template>

<script>
export default {
  props: ["error"],
}
</script>
