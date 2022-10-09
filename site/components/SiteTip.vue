<template>
  <div v-if="randomSiteTip" class="widget tip">
    <div v-if="randomSiteTip.title" class="widget-header">
      {{ randomSiteTip.title }}
    </div>
    <div class="widget-content content tip-content">
      <p v-html="randomSiteTip.content"></p>
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    randomSiteTip() {
      const siteTips = this.$store.state.config.setting.siteTips
      if (!siteTips || siteTips.length === 0) {
        return null
      }
      return siteTips[Math.floor(Math.random() * siteTips.length)]
    },
  },
}
</script>

<style lang="scss" scoped>
.tip {
  .content {
    font-size: 80%;
  }
}
</style>
