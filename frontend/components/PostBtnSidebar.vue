<template>
  <div class="widget">
    <div class="widget-content content">
      <a :href="createTopicUrl" class="button is-link btn-block">
        <span>发布新话题</span>
      </a>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    currentNodeId: {
      type: Number,
      default: 0
    }
  },
  computed: {
    createTopicUrl() {
      let url = '/topic/create'
      if (this.currentNodeId) {
        url += '?nodeId=' + this.currentNodeId
      }
      return url
    }
  }
}
</script>

<style lang="scss" scoped>
.btn-block {
  display: 100%;
  width: 100%;
}
</style>
