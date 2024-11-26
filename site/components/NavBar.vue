<template>
  <nav ref="nav" class="navbar" role="navigation" aria-label="main navigation">
    <div class="container">
      <div class="navbar-brand">
        <nuxt-link to="/" class="navbar-item">
          <i class="iconfont icon-zendea brand"></i>
        </nuxt-link>
        <a
          :class="{ 'is-active': navbarActive }"
          @click="toggleNav"
          class="navbar-burger burger"
          data-target="navbarBasic"
        >
          <span aria-hidden="true" />
          <span aria-hidden="true" />
          <span aria-hidden="true" />
        </a>
      </div>
      <div :class="{ 'is-active': navbarActive }" class="navbar-menu">
        <div class="navbar-start">
          <a
            v-for="(nav, index) in setting.siteNavs"
            :key="index"
            :href="nav.url"
            :class="{ 'is-active': $route.path == nav.url }"
            class="navbar-item"
            >{{ nav.title }}</a
          >
        </div>

        <div class="navbar-end">
          <div class="navbar-item searchFormDiv">
            <form
              id="searchForm"
              action="https://www.google.com/search"
              target="_blank"
            >
              <div class="control has-icons-right">
                <input name="q" type="hidden" value="site:zendea.com" />
                <input
                  name="q"
                  class="input"
                  type="text"
                  maxlength="30"
                  placeholder="搜索"
                />
                <span class="icon is-medium is-right">
                  <i class="iconfont icon-search" />
                </span>
              </div>
            </form>
          </div>

          <div class="navbar-item has-dropdown is-hoverable">
            <a href="/topic/create" title="发表话题" class="navbar-link">
              <i class="iconfont icon-create"></i>
            </a>
            <div class="navbar-dropdown">
              <a href="/topic/create" class="navbar-item">
                <i class="iconfont icon-topic" />&nbsp;话题
              </a>
              <a href="/article/create" class="navbar-item">
                <i class="iconfont icon-article" />&nbsp;文章
              </a>
            </div>
          </div>

          <notifier v-if="user" />

          <div v-if="isAdmin" class="navbar-item dropdown">
            <a href="/admin" title="后台管理">
              <i class="iconfont icon-wrench"></i>
            </a>
          </div>

          <div v-if="user" class="navbar-item has-dropdown is-hoverable">
            <a :href="'/user/' + user.id" class="navbar-link">
              <strong>{{ user.username }}</strong>
            </a>
            <div class="navbar-dropdown">
              <a :href="'/user/' + user.id" class="navbar-item">
                <i class="iconfont icon-home" />&nbsp;我的首页
              </a>
              <a class="navbar-item" href="/user/settings">
                <i class="iconfont icon-user" />&nbsp;编辑资料
              </a>
              <a class="navbar-item" href="/user/favorites">
                <i class="iconfont icon-favorite" />&nbsp;我的收藏
              </a>
              <a @click="signout" class="navbar-item">
                <i class="iconfont icon-logout" />&nbsp;退出登录
              </a>
            </div>
          </div>
          <div v-else class="navbar-item">
            <div class="buttons">
              <nuxt-link class="button is-small" to="/auth/signup"
                >注册
              </nuxt-link>
              <nuxt-link class="button is-small is-info" to="/auth/signin"
                >登录
              </nuxt-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
import utils from '~/common/utils'
import Notifier from '~/components/Notifier'

export default {
  components: {
    Notifier,
  },
  data() {
    return {
      navbarActive: false,
    }
  },
  computed: {
    user() {
      return this.$store.state.auth.currentUser
    },
    isAdmin() {
      const user = this.$store.state.auth.currentUser
      const LevelUserAdmin = this.$store.state.config.appinfo.user_level_admin
      if (!user) {
        return false
      }
      if (user.level === LevelUserAdmin) {
        return true
      }
      return false
    },
    setting() {
      return this.$store.state.config.setting
    },
  },
  methods: {
    async signout() {
      try {
        await this.$store.dispatch('auth/signout')
        let ref = '/'
        if (ref === '/' && process.client) {
          // 如果没配置refUrl，那么取当前地址
          ref = window.location.pathname
        }
        this.$toast.success('退出成功。', {
          duration: 1000,
          onComplete() {
            utils.linkTo(ref)
          },
        })
      } catch (e) {
        console.error(e)
      }
    },
    toggleNav() {
      this.navbarActive = !this.navbarActive
    },
  },
}
</script>

<style lang="scss" scoped></style>
