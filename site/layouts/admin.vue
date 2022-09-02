<template>
  <div class="bbs-admin">
    <nav
      class="navbar is-link admin-navbar"
      role="navigation"
      aria-label="main navigation"
    >
      <div class="container">
        <div class="navbar-brand">
          <div class="navbar-item">
            <i class="iconfont icon-zendea brand"></i>
          </div>

          <a
            :class="{ 'is-active': navbarActive }"
            @click="toggleNav"
            role="button"
            class="navbar-burger burger"
            aria-label="menu"
            aria-expanded="false"
            data-target="navbarBasic"
          >
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
          </a>
        </div>

        <div
          id="navbarBasic"
          :class="{ 'is-active': navbarActive }"
          class="navbar-menu"
        >
          <div class="navbar-start">
            <a
              :class="{ 'is-active': $route.path == '/admin' }"
              class="navbar-item"
              href="/admin"
            >
              <i class="iconfont icon-dashboard"></i>
              <span>控制台</span>
            </a>
            <a
              :class="{ 'is-active': $route.path == '/admin/nodes' }"
              class="navbar-item"
              href="/admin/nodes"
            >
              <i class="iconfont icon-node"></i>
              <span>节点</span>
            </a>
            <a
              :class="{ 'is-active': $route.path == '/admin/topics' }"
              class="navbar-item"
              href="/admin/topics"
            >
              <i class="iconfont icon-topic"></i>
              <span>话题</span>
            </a>

            <a
              :class="{ 'is-active': $route.path == '/admin/articles' }"
              class="navbar-item"
              href="/admin/articles"
            >
              <i class="iconfont icon-article"></i>
              <span>文章</span>
            </a>

            <a
              :class="{ 'is-active': $route.path == '/admin/comments' }"
              class="navbar-item"
              href="/admin/comments"
            >
              <i class="iconfont icon-comment"></i>
              <span>评论</span>
            </a>
            <a
              :class="{ 'is-active': $route.path == '/admin/users' }"
              class="navbar-item"
              href="/admin/users"
            >
              <i class="iconfont icon-user"></i>
              <span>用户</span>
            </a>
            <a
              :class="{ 'is-active': $route.path == '/admin/users/score' }"
              class="navbar-item"
              href="/admin/users/score"
            >
              <i class="iconfont icon-dollar"></i>
              <span>积分</span>
            </a>
            <a
              :class="{ 'is-active': $route.path == '/admin/links' }"
              class="navbar-item"
              href="/admin/links"
            >
              <i class="iconfont icon-link"></i>
              <span>友链</span>
            </a>
            <a
              :class="{ 'is-active': $route.path == '/admin/settings' }"
              class="navbar-item"
              href="/admin/settings"
            >
              <i class="iconfont icon-setting"></i>
              <span>设置</span>
            </a>
          </div>

          <div class="navbar-end">
            <div class="navbar-item">
              <div class="buttons">
                <a class="button is-default" href="/">
                  <strong>回到前台</strong>
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>
    <div class="container">
      <nuxt />
    </div>
  </div>
</template>

<script>
import '~/plugins/element-ui'

export default {
  middleware: ['authenticated'],
  data() {
    return {
      navbarActive: false,
    }
  },
  methods: {
    toggleNav() {
      this.navbarActive = !this.navbarActive
    },
  },
  head() {
    return {
      title: '后台管理',
    }
  },
}
</script>

<style lang="scss">
@import '@/assets/styles/variables';
.bbs-admin {
  .admin-navbar {
    .navbar-brand {
      .brand {
        font-weight: 400;
        font-size: 40px;
        line-height: 40px;
        color: #fff;
      }
    }
    .navbar-item {
      font-weight: 700;
      .iconfont {
        margin-right: 3px;
      }
    }
  }
  .page-container {
    margin: 10px 0 !important;
    padding: 10px !important;
    background: #fff;

    .toolbar {
      border-bottom: 1px solid #e9ecf0;
      margin: 0 0 5px 0;
      padding: 10px 10px 0 10px;

      .el-form-item {
        margin-bottom: 10px;
      }
    }

    .pagebar {
      text-align: right;
      margin-top: 5px;
    }
  }
}
</style>
