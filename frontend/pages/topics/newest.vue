<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <!--
        <div class="main-content" style="padding: 0;">
          <post-twitter @created="twitterCreated" />
        </div>
      -->
        <div class="main-content">
          <topics-nav :current-tab="0" />
          <topic-list :topics="topicsPage.results" :in-home="false" />
          <pagination :page="topicsPage.page" url-prefix="/topics/newest?p=" />
        </div>
      </div>
      <topic-side :score-rank="scoreRank" :links="links" />
    </div>
  </section>
</template>

<script>
// import PostTwitter from '~/components/PostTwitter'
import TopicSide from '~/components/TopicSide'
import TopicsNav from '~/components/TopicsNav'
import TopicList from '~/components/TopicList'
import Pagination from '~/components/Pagination'

export default {
  components: {
    // PostTwitter,
    TopicSide,
    TopicsNav,
    TopicList,
    Pagination
  },
  async asyncData({ $axios, query }) {
    try {
      const [topicsPage, scoreRank, links] = await Promise.all([
        $axios.get('/api/topics', {
          params: {
            page: query.p || 1
          }
        }),
        $axios.get('/api/user/score/rank'),
        $axios.get('/api/links/top')
      ])
      return { topicsPage, scoreRank, links }
    } catch (e) {
      console.error(e)
    }
  },
  methods: {
    twitterCreated(data) {
      if (this.topicsPage) {
        if (this.topicsPage.results) {
          this.topicsPage.results.unshift(data)
        } else {
          this.topicsPage.results = [data]
        }
      }
    }
  },
  head() {
    return {
      title: this.$siteTitle('默认 - 话题'),
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: this.$siteDescription()
        },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() }
      ]
    }
  }
}
</script>

<style lang="scss" scoped></style>
