<template>
  <ul class="article-list">
    <li v-for="article in articles" :key="article.articleId">
      <article
        class="article-item"
        itemscope
        itemtype="http://schema.org/BlogPosting"
      >
        <h1 class="article-title" itemprop="headline">
          <a :href="'/article/' + article.articleId">{{ article.title }}</a>
        </h1>

        <div class="article-summary" itemprop="description">
          {{ article.summary }}
        </div>

        <div class="article-meta">
          <span class="article-meta-item">
            由
            <a
              :href="'/user/' + article.user.id"
              class="article-author"
              itemprop="author"
              itemscope
              itemtype="http://schema.org/Person"
              ><span itemprop="name">{{ article.user.username }}</span></a
            >发布于
            <time
              :datetime="article.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')"
              itemprop="datePublished"
              >{{ article.createTime | prettyDate }}</time
            >
          </span>

          <span
            v-if="article.tags && article.tags.length > 0"
            class="article-meta-item"
          >
            <span
              v-for="tag in article.tags"
              :key="tag.tagId"
              class="article-tag tag"
            >
              <a :href="'/articles/tag/' + tag.tagId" class>
                {{ tag.tagName }}
              </a>
            </span>
          </span>
        </div>
      </article>
    </li>
  </ul>
</template>

<script>
export default {
  props: {
    articles: {
      type: Array,
      default() {
        return []
      },
      required: false
    }
  }
}
</script>

<style lang="scss" scoped></style>
