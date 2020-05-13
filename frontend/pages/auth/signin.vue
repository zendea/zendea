<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="widget">
          <div class="widget-header">
            登录
          </div>
          <div class="widget-content">
            <div class="field">
              <label class="label">用户名/邮箱</label>
              <div class="control has-icons-left">
                <input
                  v-model="username"
                  @keyup.enter="submitLogin"
                  class="input is-success"
                  type="text"
                  placeholder="请输入用户名或邮箱"
                />
                <span class="icon is-small is-left"
                  ><i class="iconfont icon-user"
                /></span>
              </div>
            </div>

            <div class="field">
              <label class="label">密码</label>
              <div class="control has-icons-left">
                <input
                  v-model="password"
                  @keyup.enter="submitLogin"
                  class="input"
                  type="password"
                  placeholder="请输入密码"
                />
                <span class="icon is-small is-left"
                  ><i class="iconfont icon-password"
                /></span>
              </div>
            </div>

            <div class="field">
              <label class="label">验证码</label>
              <div class="control has-icons-left">
                <div class="field is-horizontal">
                  <div class="field" style="width:100%;">
                    <input
                      v-model="captchaCode"
                      @keyup.enter="submitLogin"
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
                <button @click="submitLogin" class="button is-link">
                  登录
                </button>
                <nuxt-link class="button is-text" to="/auth/signup">
                  没有账号？点击这里去注册&gt;&gt;
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
              <li class="list-group-item">
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
      password: '',
      captchaId: '',
      captchaUrl: '',
      captchaCode: ''
    }
  },
  computed: {
    currentUser() {
      return this.$store.state.auth.currentUser
    },
    isLogin() {
      return !!this.currentUser
    }
  },
  mounted() {
    if (this.redirectIfLogined()) {
      return
    }
    this.showCaptcha()
  },
  methods: {
    async submitLogin() {
      try {
        if (!this.username) {
          this.$toast.error('请输入用户名或邮箱')
          return
        }
        if (!this.password) {
          this.$toast.error('请输入密码')
          return
        }
        if (!this.captchaCode) {
          this.$toast.error('请输入验证码')
          return
        }
        await this.$store.dispatch('auth/signin', {
          captchaId: this.captchaId,
          captchaCode: this.captchaCode,
          username: this.username,
          password: this.password,
          ref: this.ref
        })
        this.$toast.success('登录成功。', {
          duration: 1000,
          onComplete() {
            if (this.ref) {
              // 跳到登录前
              utils.linkTo(this.ref)
            } else {
              // 跳到个人主页
              utils.linkTo('/')
            }
          }
        })
      } catch (e) {
        this.$toast.error(e.message || e)
        await this.showCaptcha()
      }
    },
    async showCaptcha() {
      try {
        const ret = await this.$axios.get('/api/captcha/request')
        this.captchaId = ret.captchaId
        this.captchaUrl = ret.captchaUrl
      } catch (e) {
        this.$toast.error(e.message || e)
      }
    },
    /**
     * 如果已经登录了，那么直接跳转
     * @returns {boolean}
     */
    redirectIfLogined() {
      if (this.isLogin) {
        const me = this
        this.$toast.success('登录成功！', {
          duration: 3000,
          keepOnHover: false,
          position: 'top-center',
          onComplete() {
            if (me.ref && !utils.isSigninUrl(me.ref)) {
              utils.linkTo(me.ref)
            } else {
              utils.linkTo('/')
            }
          }
        })
        return true
      }
      return false
    }
  },
  head() {
    return {
      title: this.$siteTitle('登录')
    }
  }
}
</script>
