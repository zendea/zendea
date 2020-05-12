<template>
  <a
    :class="{ button: isButton }"
    @click="giteeLogin"
    class="is-default is-block is-link"
  >
    <i class="iconfont icon-gitee" />&nbsp;
    <strong>{{ title }}</strong>
  </a>
</template>

<script>
export default {
  name: 'GiteeLogin',
  props: {
    title: {
      type: String,
      default: 'Gitee'
    },
    refUrl: {
      // 登录来源地址，控制登录成功之后要跳到该地址，默认跳转到首页，如果设置为空就跳回登录页
      type: String,
      default: '/'
    },
    isButton: {
      type: Boolean,
      default: true
    }
  },
  data() {
    return {
      refUrlValue: this.refUrl
    }
  },
  methods: {
    async giteeLogin() {
      try {
        if (!this.refUrlValue && process.client) {
          // 如果没配置refUrl，那么取当前地址
          this.refUrlValue = window.location.pathname
        }
        const ret = await this.$axios.get('/api/oauth/gitee/authorize', {
          params: {
            ref: this.refUrlValue
          }
        })
        window.location = ret.url
      } catch (e) {
        console.error(e)
        this.$toast.error('登录失败：' + (e.message || e))
      }
    }
  }
}
</script>

<style lang="scss" scoped></style>
