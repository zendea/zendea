<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <div class="tabs">
          <ul>
            <li :class="{ 'is-active': activeTab === 'topics' }">
              <a :href="'/user/' + user.id + '?tab=topics'">
                <span class="icon is-small">
                  <i class="iconfont icon-topic" aria-hidden="true" />
                </span>
                <span>话题</span>
              </a>
            </li>
            <li :class="{ 'is-active': activeTab === 'articles' }">
              <a :href="'/user/' + user.id + '?tab=articles'">
                <span class="icon is-small">
                  <i class="iconfont icon-article" aria-hidden="true" />
                </span>
                <span>文章</span>
              </a>
            </li>
          </ul>
        </div>

        <div v-if="activeTab === 'topics'">
          <div v-if="recentTopics && recentTopics.length">
            <topic-list :topics="recentTopics" />
            <div class="more">
              <a :href="'/user/' + user.id + '/topics'">查看更多&gt;&gt;</a>
            </div>
          </div>
          <div v-else class="notification is-primary" style="margin-top: 10px">
            暂无话题
          </div>
        </div>

        <div v-if="activeTab === 'articles'">
          <div v-if="recentArticles && recentArticles.length">
            <article-list :articles="recentArticles" />
            <div class="more">
              <a :href="'/user/' + user.id + '/articles'">查看更多&gt;&gt;</a>
            </div>
          </div>
          <div v-else class="notification is-primary" style="margin-top: 10px">
            暂无文章
          </div>
        </div>
      </div>
      <div class="right-container">
        <user-center-sidebar :user="user" />
        <div class="widget">
          <div class="widget-header">关注</div>
          <div class="widget-content watch-actions">
            <div
              v-if="!isOwner"
              :class="{ active: user.watched }"
              @click="watch(user)"
              class="action watch"
              title="关注"
            >
              <i class="iconfont icon-eye" />
            </div>
            <span v-if="!isOwner" class="split"></span>
            <div v-for="user in userWatchers" :key="user.id">
              <a
                :href="'/user/' + user.id"
                :alt="user.username"
                :title="user.username"
              >
                avatar
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import utils from '~/common/utils'
import TopicList from '~/components/TopicList'
import ArticleList from '~/components/ArticleList'
import UserCenterSidebar from '~/components/UserCenterSidebar'

const defaultTab = 'topics'

export default {
  components: {
    TopicList,
    ArticleList,
    UserCenterSidebar,
  },
  async asyncData({ $axios, params, query, error }) {
    let user
    try {
      user = await $axios.get('/api/profile/' + params.userId)
    } catch (err) {
      error({
        statusCode: 404,
        message: err.message || '系统错误',
      })
      return
    }

    const [watched, userWatchers] = await Promise.all([
      $axios.get('/api/watch/watched', {
        params: {
          userId: params.userId,
        },
      }),
      $axios.get('/api/users/' + params.userId + '/recentwatchers'),
    ])

    const activeTab = query.tab || defaultTab
    let recentTopics = null
    let recentArticles = null
    if (activeTab === 'topics') {
      recentTopics = await $axios.get(
        '/api/topics/user/recent/' + params.userId + ''
      )
    } else if (activeTab === 'articles') {
      recentArticles = await $axios.get(
        '/api/articles/user/recent/' + params.userId + ''
      )
    }
    return {
      activeTab,
      user,
      watched: watched.watched,
      userWatchers,
      recentTopics,
      recentArticles,
    }
  },
  data() {
    return {}
  },
  computed: {
    currentUser() {
      return this.$store.state.auth.currentUser
    },
    isOwner() {
      const current = this.$store.state.auth.currentUser
      return this.user && current && this.user.id === current.id
    },
  },
  methods: {
    async watch(user) {
      try {
        if (this.watched) {
          await this.$axios.delete('/api/watch/delete', {
            params: {
              userId: user.id,
            },
          })
          this.watched = false
          user.watchCount--
          this.$toast.success('已取消关注！')
          this.userWatchers = await this.$axios.get(
            '/api/users/' + user.id + '/recentwatchers'
          )
        } else {
          const ret = await this.$axios.post('/api/users/' + user.id + '/watch')
          if (ret != null && ret.success === false) {
            throw ret
          } else {
            user.watched = true
            user.watchCount++
            this.userWatchers = this.userWatchers || []
            this.userWatchers.unshift(this.$store.state.auth.currentUser)
            this.$toast.success('关注成功')
          }
        }
      } catch (e) {
        if (e.success === false) {
          this.$toast.info('请登录后再关注', {
            action: {
              text: '去登录',
              onClick: (e, toastObject) => {
                utils.toSignin()
              },
            },
          })
        } else {
          user.watched = true
          this.$toast.error(e.message || e)
        }
      }
    },
  },
  head() {
    return {
      title: this.$siteTitle(this.user.username),
    }
  },
}
</script>

<style lang="scss" scoped>
.tabs {
  margin-bottom: 5px;
}
.more {
  text-align: right;
}
.watch-actions {
  display: flex;
  height: 42px;
  width: max-content;
  padding: 5px 10px 5px 5px;

  background: #fff;
  // border: 1px solid #dae0e4;
  border-radius: 30px;
  vertical-align: middle;

  .avatar {
    margin-left: 5px;
    border: solid 1px #e8ecee;
  }

  .split {
    display: inline-block;
    margin: auto 7px auto 10px;
    height: 16px;
    width: 1px;
    opacity: 0.4;
    background: #dae0e4;
    vertical-align: middle;
  }

  .action {
    margin-left: 5px;
    height: 30px;
    width: 30px;
    line-height: 30px;
    text-align: center;
    color: #e7672e;
    background-color: rgba(126, 107, 1, 0.08);
    border-radius: 50%;
    transition: all 0.5s;
    cursor: pointer;

    i {
      font-size: 16px;
    }

    &:hover,
    &.active {
      color: #fff;
      background-color: #e7672e;
    }
  }
}
</style>
