<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <load-more
          v-slot="{ results }"
          :init-data="articlesPage"
          :params="{ tagId: tag.tagId }"
          url="/api/articles/tag/{ tag.tagId }"
        >
          <article-list :articles="results" :show-ad="true" />
        </load-more>
      </div>
      <div class="right-container">
        <site-notice />
      </div>
    </div>
  </section>
</template>

<script>
import ArticleList from '~/components/ArticleList'
import LoadMore from '~/components/LoadMore'
import SiteNotice from '~/components/SiteNotice'

export default {
  components: { ArticleList, LoadMore, SiteNotice },
  async asyncData({ $axios, params }) {
    try {
      const [tag, articlesPage] = await Promise.all([
        $axios.get('/api/tag/' + params.tagId),
        $axios.get('/api/articles/tag/' + params.tagId)
      ])
      return {
        tag,
        articlesPage
      }
    } catch (e) {
      console.error(e)
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.tag.tagName + ' - 文章'),
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
