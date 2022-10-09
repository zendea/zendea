<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="main-content">
          <topics-nav />
          <topic-list :topics="topicsPage.results" :show-ad="true" />
          <pagination :page="topicsPage.page" url-prefix="/topics?p=" />
        </div>
      </div>
      <topic-side :score-rank="scoreRank" :links="links" :stat="stat" />
    </div>
  </section>
</template>

<script>
import TopicSide from '~/components/TopicSide'
import TopicsNav from '~/components/TopicsNav'
import TopicList from '~/components/TopicList'
import Pagination from '~/components/Pagination'

export default {
  components: {
    TopicSide,
    TopicsNav,
    TopicList,
    Pagination,
  },
  async asyncData({ $axios, query }) {
    try {
      const [topicsPage, scoreRank, links, stat] = await Promise.all([
        $axios.get('/api/topics', {
          params: {
            page: query.p || 1,
          },
        }),
        $axios.get('/api/user/score/rank'),
        $axios.get('/api/links/top'),
        $axios.get('/api/stat'),
      ])
      return { topicsPage, scoreRank, links, stat }
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
    },
  },
  head() {
    return {
      title: this.$siteTitle('社区'),
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: this.$siteDescription(),
        },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() },
      ],
    }
  },
}
</script>

<style lang="scss" scoped></style>
