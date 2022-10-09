<template>
  <div class="widget">
    <div class="widget-header">站点统计</div>
    <div class="widget-content extra-info">
      <ul>
        <li>
          <span>主题数</span><br />
          <b>{{ stat.topicCount }}</b>
        </li>
        <li>
          <span>回帖数</span><br />
          <b>{{ stat.commentCount }}</b>
        </li>
        <li>
          <span>用户数</span><br />
          <b>{{ stat.userCount }}</b>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    stat: {
      type: Object,
      required: true,
    },
  },
}
</script>

<style lang="scss" scoped>
.extra-info {
  ul {
    display: flex;
    li {
      width: 100%;
      text-align: center;
      span {
        font-size: 13px;
        font-weight: 400;
        color: var(--text-color3);
      }
      a,
      b {
        color: var(--text-color);
      }
    }
  }
}
</style>
