<template>
  <ul class="topic-list topic-wrap">
    <li v-for="topic in topics" :key="topic.topicId" class="topic-item">
      <article itemscope itemtype="http://schema.org/BlogPosting">
        <div class="topic-header">
          <div class="topic-header-left">
            <a :href="'/user/' + topic.user.id" :title="topic.user.username">
              <img
                v-if="topic.user.avatar"
                v-lazy="topic.user.avatar"
                class="avatar"
              />
              <avatar v-else :username="topic.user.username" :size="50" />
            </a>
          </div>
          <div class="topic-header-center">
            <h1 itemprop="headline">
              <span v-if="topic.node && !hasNodeId" class="topic-node">
                <a :href="'/topics/' + topic.node.nodeId" class="node">{{
                  topic.node.name
                }}</a>
              </span>
              <span class="topic-title">
                <a :href="'/topic/' + topic.topicId" :title="topic.title">{{
                  topic.title
                }}</a>
              </span>
            </h1>

            <div class="topic-meta">
              <span
                class="meta-item"
                itemprop="author"
                itemscope
                itemtype="http://schema.org/Person"
              >
                <a :href="'/user/' + topic.user.id" itemprop="name">{{
                  topic.user.username
                }}</a>
              </span>
              <span
                v-if="topic.lastCommentUser.id"
                class="meta-item"
                itemprop="author"
                itemscope
                itemtype="http://schema.org/Person"
              >
                最后由
                <a :href="'/user/' + topic.lastCommentUser.id">{{
                  topic.lastCommentUser.username
                }}</a>
                回复于
              </span>
              <span class="meta-item">
                <time
                  :datetime="
                    topic.lastCommentTime | formatDate('yyyy-MM-ddTHH:mm:ss')
                  "
                  itemprop="datePublished"
                  >{{ topic.lastCommentTime | prettyDate }}</time
                >
              </span>
              <span class="meta-item">
                <span v-for="tag in topic.tags" :key="tag.tagId" class="tag">
                  <a :href="'/topics/tag/' + tag.tagId">{{ tag.tagName }}</a>
                </span>
              </span>
            </div>
          </div>
          <div v-if="!inHome" class="topic-header-right">
            <div v-if="topic.likeCount" class="like">
              <span
                :class="{ liked: topic.liked }"
                @click="like(topic)"
                class="like-btn"
              >
                <i class="iconfont icon-heart" />
              </span>
              <span v-if="topic.likeCount" class="like-count">{{
                topic.likeCount
              }}</span>
            </div>
            <span v-if="topic.commentCount" class="count">{{
              topic.commentCount
            }}</span>
          </div>
        </div>
        <ul
          v-if="topic.imageList && topic.imageList.length > 0"
          class="topic-images"
        >
          <li v-for="image in topic.imageList" :key="image">
            <a
              :href="'/topic/' + topic.topicId"
              :title="topic.title"
              class="topic-image-item"
            >
              <img v-lazy="image" />
            </a>
          </li>
        </ul>
      </article>
    </li>
  </ul>
</template>

<script>
import utils from '~/common/utils'
export default {
  props: {
    topics: {
      type: Array,
      default() {
        return []
      },
      required: false
    },
    inHome: {
      type: Boolean,
      default: false
    },
    hasNodeId: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    async like(topic) {
      try {
        await this.$axios.post('/api/topic/like/' + topic.topicId)
        topic.liked = true
        topic.likeCount++
      } catch (e) {
        if (e.code === 1) {
          this.$toast.info('请登录后点赞！！！', {
            action: {
              text: '去登录',
              onClick: (e, toastObject) => {
                utils.toSignin()
              }
            }
          })
        } else {
          this.$toast.error(e.message || e)
        }
      }
    }
  }
}
</script>

<style lang="scss" scoped></style>
