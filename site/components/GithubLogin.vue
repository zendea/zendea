<template>
  <a
    :class="{ button: isButton }"
    @click="githubLogin"
    class="is-default is-block is-primary"
  >
    <i class="iconfont icon-github" />&nbsp;
    <strong>{{ title }}</strong>
  </a>
</template>

<script>
export default {
  name: 'GithubLogin',
  props: {
    title: {
      type: String,
      default: 'Github'
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
    async githubLogin() {
      try {
        if (!this.refUrlValue && process.client) {
          // 如果没配置refUrl，那么取当前地址
          this.refUrlValue = window.location.pathname
        }
        const ret = await this.$axios.get('/api/oauth/github/authorize', {
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
