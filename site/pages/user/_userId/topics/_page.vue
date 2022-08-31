<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <nav class="breadcrumb my-breadcrumb">
          <ul>
            <li><a href="article">首页</a></li>
            <li>
              <a :href="'/user/' + user.id + '?tab=topics'">{{
                user.username
              }}</a>
            </li>
            <li class="is-active">
              <a href="#" aria-current="page">话题列表</a>
            </li>
          </ul>
        </nav>

        <topic-list :topics="topicsPage.results" />
        <pagination
          :page="topicsPage.page"
          :url-prefix="'/user/' + user.id + '/topics/'"
        />
      </div>
      <div class="right-container">
        <user-center-sidebar :user="user" />
      </div>
    </div>
  </section>
</template>

<script>
import TopicList from '~/components/TopicList'
import Pagination from '~/components/Pagination'
import UserCenterSidebar from '~/components/UserCenterSidebar'
export default {
  components: {
    TopicList,
    Pagination,
    UserCenterSidebar
  },
  async asyncData({ $axios, params, error }) {
    let user
    try {
      user = await $axios.get('/api/profile/' + params.userId)
    } catch (err) {
      error({
        statusCode: 404,
        message: err.message || '系统错误'
      })
      return
    }

    const [topicsPage] = await Promise.all([
      $axios.get('/api/user/topics/' + params.userId, {
        params: {
          userId: params.userId,
          page: params.page
        }
      })
    ])

    return {
      user,
      topicsPage
    }
  },
  computed: {
    currentUser() {
      return this.$store.state.auth.currentUser
    },
    // 是否是主人态
    isOwner() {
      const current = this.$store.state.auth.currentUser
      return this.user && current && this.user.id === current.id
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.user.username + ' - 话题')
    }
  }
}
</script>

<style lang="scss" scoped></style>
