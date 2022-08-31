<template>
  <section class="main-container">
    <div class="page-container">
      <div v-if="siteIndexHtml" v-html="siteIndexHtml"></div>
      <div class="widget">
        <div class="widget-header">社区精华帖</div>
        <div class="widget-content columns">
          <div class="is-6 column">
            <topic-list :topics="topics.even" :in-home="true" />
          </div>
          <div class="is-6 column">
            <topic-list :topics="topics.odd" :in-home="true" />
          </div>
        </div>
        <div class="widget-footer">
          <a href="/topics/excellent">查看更多精华帖...</a>
        </div>
      </div>
      <index-sections :sections="sections" />
    </div>
  </section>
</template>

<script>
import TopicList from '~/components/TopicList'
import IndexSections from '~/components/IndexSections'

export default {
  components: {
    TopicList,
    IndexSections
  },
  async asyncData({ $axios, params }) {
    try {
      const [sections, topics] = await Promise.all([
        $axios.get('/api/sections'),
        $axios.get('/api/topics/excellent')
      ])
      return { sections, topics }
    } catch (e) {
      console.error(e)
    }
  },
  computed: {
    siteIndexHtml() {
      return this.$store.state.config.setting.siteIndexHtml
    }
  },
  head() {
    return {
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
