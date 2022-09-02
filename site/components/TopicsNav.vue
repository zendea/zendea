<template>
  <div class="topics-nav">
    <ul class="topics-nav-list">
      <li :class="{ active: currentTab === 0 }" class="topics-nav-item">
        <a href="/topics/newest">默认</a>
      </li>
      <li :class="{ active: currentTab === -1 }" class="topics-nav-item">
        <a href="/topics/excellent">
          <i class="iconfont icon-diamond"></i> 精华帖
        </a>
      </li>
      <li :class="{ active: currentTab === -2 }" class="topics-nav-item">
        <a href="/topics/no_reply">无人问津</a>
      </li>
      <li :class="{ active: currentTab === -3 }" class="topics-nav-item">
        <a href="/topics/last">最新发布</a>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  props: {
    currentTab: {
      type: Number,
      default: 0,
    },
  },
}
</script>

<style lang="scss" scoped>
.topics-nav {
  font-size: 16px;
  font-weight: bold;
  border-bottom: 1px solid #e6ecf0 !important;

  .topics-nav-list {
    .topics-nav-item {
      display: inline-block;
      padding: 0 15px;

      a {
        position: relative;
        display: inline-block;
        padding: 14px 0;
        line-height: 22px;
        text-align: center;
        font-weight: 500;
        color: #8590a6;

        &:hover {
          color: #343a40;
        }
      }

      &.active {
        a {
          color: #343a40;
          font-weight: 700;

          &:after {
            position: absolute;
            right: -3px;
            bottom: -1px;
            left: -3px;
            height: 3px;
            background: #3273dc;
            content: '';
          }
        }
      }
      .icon-divider-vertical {
        color: #ccc;
      }
    }
  }
}
</style>
