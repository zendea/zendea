<template>
  <section class="main">
    <div class="container">
      <div class="main-container columns">
        <div class="column is-12">
          <div class="widget">
            <div class="widget-header">
              <span>标签</span>
            </div>
            <div class="widget-content">
              <div class="tags are-medium">
                <span
                  v-for="tag in tagsPage.results"
                  :key="tag.tagId"
                  class="tag is-normal"
                >
                  <a :href="'/topics/tag/' + tag.tagId" :title="tag.tagName">
                    {{ tag.tagName }}
                  </a>
                </span>
              </div>
            </div>
            <pagination :page="tagsPage.page" url-prefix="/tags/" />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import Pagination from '~/components/Pagination'
// import SiteNotice from '~/components/SiteNotice'

export default {
  components: {
    Pagination
    // SiteNotice
  },
  async asyncData({ $axios, params }) {
    const [tagsPage] = await Promise.all([
      $axios.get('/api/tags', {
        params: {
          page: params.page
        }
      })
    ])
    return {
      tagsPage
    }
  },
  head() {
    return {
      title: this.$siteTitle('标签')
    }
  }
}
</script>

<style lang="scss" scoped></style>
