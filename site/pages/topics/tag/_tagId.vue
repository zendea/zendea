<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <div class="tag-header">
          <div class="title">
            标签:
            <span class="name">{{ tag.tagName }}</span>
          </div>
        </div>
        <topic-list :topics="topicsPage.results" />
        <pagination
          :page="topicsPage.page"
          :url-prefix="'/topics/' + tag.tagId + '?p='"
        />
      </div>
      <topic-side :score-rank="scoreRank" :links="links" />
    </div>
  </section>
</template>

<script>
import TopicSide from '~/components/TopicSide'
// import TopicsNav from '~/components/TopicsNav'
import TopicList from '~/components/TopicList'
import Pagination from '~/components/Pagination'

export default {
  components: {
    TopicSide,
    // TopicsNav,
    TopicList,
    Pagination
  },
  async asyncData({ $axios, params, query }) {
    const [tag, topicsPage, scoreRank, links] = await Promise.all([
      $axios.get('/api/tag/' + params.tagId),
      $axios.get('/api/topics/tag', {
        params: {
          tagId: params.tagId,
          page: query.p || 1
        }
      }),
      $axios.get('/api/user/score/rank'),
      $axios.get('/api/links/top')
    ])
    return {
      tag,
      topicsPage,
      scoreRank,
      links
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.tag.tagName + ' - 话题'),
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

<style lang="scss" scoped>
.tag-header {
  margin-bottom: 5px;
  border-bottom: 1px solid #f2f2f2;
  padding: 5px 10px;

  .title {
    font-size: 14px;
    color: #999;
    margin-bottom: 8px;
    .name {
      color: #333;
      font-size: 24px;
      margin-left: 10px;
    }
  }
}
</style>
