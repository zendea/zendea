<template>
  <div>
    <section class="main">
      <div class="container main-container left-main">
        <div class="left-container">
          <div class="main-content">
            <div class="node-header">
              <div class="title">
                {{ node.name }}
                <span class="total">
                  共有 {{ node.topicCount }} 个讨论话题
                </span>
              </div>
              <div class="summary">
                <p>{{ node.description }}</p>
              </div>
            </div>
            <topic-list :topics="topicsPage.results" :in-home="false" />
            <pagination
              :page="topicsPage.page"
              :url-prefix="'/topics/' + node.nodeId + '?p='"
            />
            <div v-if="!topicsPage.page.total" class="summary">
              本节点暂无话题。
            </div>
          </div>
        </div>
        <topic-side
          :current-node-id="node.nodeId"
          :score-rank="scoreRank"
          :links="links"
        />
      </div>
    </section>
  </div>
</template>

<script>
// import PostTwitter from '~/components/PostTwitter'
import TopicSide from '~/components/TopicSide'
// import TopicsNav from '~/components/TopicsNav'
import TopicList from '~/components/TopicList'
import Pagination from '~/components/Pagination'

export default {
  components: {
    // PostTwitter,
    TopicSide,
    // TopicsNav,
    TopicList,
    Pagination
  },
  async asyncData({ $axios, params, query }) {
    const [node, topicsPage, scoreRank, links] = await Promise.all([
      $axios.get('/api/node/' + params.nodeId),
      $axios.get('/api/topics/node', {
        params: {
          nodeId: params.nodeId,
          page: query.p || 1
        }
      }),
      $axios.get('/api/user/score/rank'),
      $axios.get('/api/links/top')
    ])
    return {
      node,
      topicsPage,
      scoreRank,
      links
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
      title: this.$siteTitle(this.node.name + ' - 话题'),
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
.node-header {
  margin-bottom: 5px;
  border-bottom: 1px solid #f2f2f2;
  padding: 5px 10px;

  .container {
    padding: 15px 20px;
  }
  .title {
    font-size: 24px;
    color: #333;
    margin-bottom: 8px;
    .total {
      color: #999;
      font-size: 14px;
      margin-left: 10px;
    }
  }
  .summary {
    font-size: 14px;
  }
}
</style>
