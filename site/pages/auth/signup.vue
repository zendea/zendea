<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="widget">
          <div class="widget-header">
            注册
          </div>
          <div class="widget-content">
            <div class="field">
              <label class="label">用户名</label>
              <div class="control has-icons-left">
                <input
                  v-model="username"
                  @keyup.enter="signup"
                  class="input is-success"
                  type="text"
                  placeholder="请输入用户名"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-user" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">邮箱</label>
              <div class="control has-icons-left">
                <input
                  v-model="email"
                  @keyup.enter="signup"
                  class="input is-success"
                  type="text"
                  placeholder="请输入邮箱"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-email" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">密码</label>
              <div class="control has-icons-left">
                <input
                  v-model="password"
                  @keyup.enter="signup"
                  class="input"
                  type="password"
                  placeholder="请输入密码"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">确认密码</label>
              <div class="control has-icons-left">
                <input
                  v-model="rePassword"
                  @keyup.enter="signup"
                  class="input"
                  type="password"
                  placeholder="请再次输入密码"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">验证码</label>
              <div class="control has-icons-left">
                <div class="field is-horizontal">
                  <div class="field" style="width:100%;">
                    <input
                      v-model="captchaCode"
                      @keyup.enter="signup"
                      class="input"
                      type="text"
                      placeholder="验证码"
                    />
                    <span class="icon is-small is-left"
                      ><i class="iconfont icon-captcha"
                    /></span>
                  </div>
                  <div v-if="captchaUrl" class="field">
                    <a @click="showCaptcha"
                      ><img :src="captchaUrl" style="height: 40px;"
                    /></a>
                  </div>
                </div>
              </div>
            </div>

            <div class="field">
              <div class="control">
                <button @click="signup" class="button is-link">
                  注册
                </button>
                <nuxt-link class="button is-text" to="/auth/signin">
                  已有账号，前往登录&gt;&gt;
                </nuxt-link>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="right-container">
        <div class="widget">
          <div class="widget-header">
            使用其他平台帐号登录
          </div>
          <div class="widget-content">
            <ul class="list-group">
              <li class="list-group-item">
                <github-login :ref-url="ref" />
              </li>
              <li class="list-group-ltem">
                <gitee-login :ref-url="ref" />
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import utils from '~/common/utils'
import GithubLogin from '~/components/GithubLogin'
import GiteeLogin from '~/components/GiteeLogin'
// import QqLogin from '~/components/QqLogin'
export default {
  components: {
    GithubLogin,
    GiteeLogin
    // QqLogin
  },
  asyncData({ params, query }) {
    return {
      ref: query.ref
    }
  },
  data() {
    return {
      username: '',
      email: '',
      password: '',
      rePassword: '',
      captchaId: '',
      captchaUrl: '',
      captchaCode: ''
    }
  },
  mounted() {
    this.showCaptcha()
  },
  methods: {
    async signup() {
      const me = this
      if (!me.username) {
        this.$toast.error('用户名不能为空')
        return
      }
      if (!me.email) {
        this.$toast.error('邮箱不能为空')
        return
      }
      try {
        await this.$store.dispatch('auth/signup', {
          captchaId: me.captchaId,
          captchaCode: me.captchaCode,
          username: me.username,
          email: me.email,
          password: me.password,
          rePassword: me.rePassword,
          ref: me.ref
        })

        utils.linkTo('/auth/signin')
      } catch (err) {
        this.$toast.error(err.message || err)
        await this.showCaptcha()
      }
    },
    async showCaptcha() {
      try {
        const ret = await this.$axios.get('/api/captcha/request')
        this.captchaId = ret.captchaId
        this.captchaUrl = ret.captchaUrl
        this.captchaCode = ''
      } catch (e) {
        this.$toast.error(e.message || e)
      }
    }
  },
  head() {
    return {
      title: this.$siteTitle('注册')
    }
  }
}
</script>

<style lang="scss" scoped></style>
