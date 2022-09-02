<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="main-content">
          <article
            class="topic-detail topic-wrap"
            itemscope
            itemtype="http://schema.org/BlogPosting"
          >
            <div class="topic-header">
              <div class="topic-header-left">
                <a
                  :href="'/user/' + topic.user.id"
                  :title="topic.user.username"
                >
                  avatar
                </a>
              </div>
              <div class="topic-header-center">
                <h1 class="title" itemprop="headline">
                  {{ topic.title }}
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
                  <span class="meta-item">
                    <time
                      :datetime="
                        topic.lastCommentTime
                          | formatDate('yyyy-MM-ddTHH:mm:ss')
                      "
                      itemprop="datePublished"
                      >{{ topic.lastCommentTime | prettyDate }}</time
                    >
                  </span>
                  <span class="meta-item">
                    <a
                      v-if="topic.node"
                      :href="'/topics/' + topic.node.nodeId"
                      class="node"
                      >{{ topic.node.name }}</a
                    >
                  </span>
                  <span class="meta-item">
                    <span
                      v-for="tag in topic.tags"
                      :key="tag.tagId"
                      class="tag"
                    >
                      <a :href="'/topics/tag/' + tag.tagId">{{
                        tag.tagName
                      }}</a>
                    </span>
                  </span>
                  <span v-if="isOwner" class="meta-item act">
                    <a @click="deleteTopic(topic.topicId)">
                      <i class="iconfont icon-delete" />&nbsp;删除
                    </a>
                  </span>
                  <!-- 话题类型为普通时才可以修改 -->
                  <span
                    v-if="isOwner && topic.type === 0"
                    class="meta-item act"
                  >
                    <a :href="'/topic/edit/' + topic.topicId">
                      <i class="iconfont icon-edit" />&nbsp;修改
                    </a>
                  </span>
                </div>
              </div>
              <div class="topic-header-right">
                <div class="like">
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
                <span class="count"
                  >{{ topic.commentCount }}&nbsp;/&nbsp;{{
                    topic.viewCount
                  }}</span
                >
              </div>
            </div>

            <div class="content" itemprop="articleBody">
              <div
                v-html="topic.content"
                v-lazy-container="{ selector: 'img' }"
              ></div>
              <div v-if="topic.imageList">
                <figure v-for="image in topic.imageList" :key="image">
                  <img v-lazy="image" />
                </figure>
              </div>
            </div>

            <div class="topic-actions">
              <div
                :class="{ active: favorited }"
                @click="addFavorite(topic.topicId)"
                class="action favorite"
                title="收藏"
              >
                <i class="iconfont icon-favorite" />
              </div>
              <span class="split"></span>
              <div
                :class="{ active: topic.liked }"
                @click="like(topic)"
                class="action like"
                title="点赞"
              >
                <i class="iconfont icon-heart" />
              </div>
              <div v-for="user in likeUsers" :key="user.id">
                <a
                  :href="'/user/' + user.id"
                  :alt="user.username"
                  target="_blank"
                >
                  avatar
                </a>
              </div>
            </div>
          </article>
        </div>

        <!-- 评论 -->
        <comment
          :entity-id="topic.topicId"
          :comments-page="commentsPage"
          :comment-count="topic.commentCount"
          entity-type="topic"
        />
      </div>
      <div class="right-container">
        <div class="user-info">
          <div class="base">
            <div>
              <a :href="'/user/' + topic.user.id" :alt="topic.user.username">
                avatar
              </a>
            </div>
            <div class="username">
              <a :href="'/user/' + topic.user.id" :alt="topic.user.username">
                {{ topic.user.username }}
                <span :class="'level-' + topic.user.level">
                  ({{ topic.user.levelName }})
                </span>
              </a>
            </div>
            <div class="description">
              {{ topic.user.description }}
            </div>
          </div>
          <div class="extra">
            <ul class="basic">
              <li>
                <span>UID</span><br />
                <b>{{ topic.user.id }}</b>
              </li>
              <li>
                <span>积分</span><br />
                <b>{{ topic.user.score }}</b>
              </li>
            </ul>
            <div class="other">
              注册时间:
              <span>
                {{ topic.user.createTime | formatDate('yyyy-MM-dd') }}
              </span>
            </div>
          </div>
        </div>

        <div ref="toc" v-if="topic.toc" class="widget no-bg toc">
          <div class="widget-header">目录</div>
          <div v-html="topic.toc" class="widget-content" />
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import utils from '~/common/utils'
import Comment from '~/components/Comment'

export default {
  components: {
    Comment,
  },
  async asyncData({ $axios, params, error }) {
    let topic
    try {
      topic = await $axios.get('/api/topic/' + params.id)
    } catch (e) {
      error({
        statusCode: 404,
        message: '话题不存在',
      })
      return
    }

    const [favorited, commentsPage, likeUsers] = await Promise.all([
      $axios.get('/api/favorites/favorited', {
        params: {
          entityType: 'topic',
          entityId: params.id,
        },
      }),
      $axios.get('/api/comments', {
        params: {
          entityType: 'topic',
          entityId: params.id,
        },
      }),
      $axios.get('/api/topic/' + params.id + '/recentlikes'),
    ])

    return {
      topic,
      commentsPage,
      favorited: favorited.favorited,
      likeUsers,
    }
  },
  computed: {
    isOwner() {
      return (
        this.$store.state.auth.currentUser &&
        this.topic &&
        this.$store.state.auth.currentUser.id === this.topic.user.id
      )
    },
  },
  mounted() {
    utils.handleToc(this.$refs.toc)
  },
  methods: {
    async addFavorite(topicId) {
      try {
        if (this.favorited) {
          await this.$axios.delete('/api/favorite/delete', {
            params: {
              entityType: 'topic',
              entityId: topicId,
            },
          })
          this.favorited = false
          this.$toast.success('已取消收藏！')
        } else {
          await this.$axios.post('/api/topic/' + topicId + '/favorite')
          this.favorited = true
          this.$toast.success('收藏成功')
        }
      } catch (e) {
        console.error(e)
        this.$toast.error('收藏失败：' + (e.message || e))
      }
    },
    async deleteTopic(topicId) {
      if (process.client && !window.confirm('是否确认删除该话题？')) {
        return
      }
      try {
        await this.$axios.post('/api/topic/delete/' + topicId)
        this.$toast.success('删除成功', {
          duration: 2000,
          onComplete() {
            utils.linkTo('/topics')
          },
        })
      } catch (e) {
        console.error(e)
        this.$toast.error('删除失败：' + (e.message || e))
      }
    },
    async like(topic) {
      try {
        await this.$axios.post('/api/topic/' + topic.topicId + '/like')
        topic.liked = true
        topic.likeCount++
        this.likeUsers = this.likeUsers || []
        this.likeUsers.unshift(this.$store.state.auth.currentUser)
      } catch (e) {
        if (e.code === 1) {
          this.$toast.info('请登录后点赞！！！', {
            action: {
              text: '去登录',
              onClick: (e, toastObject) => {
                utils.toSignin()
              },
            },
          })
        } else {
          topic.liked = true
          this.$toast.error(e.message || e)
        }
      }
    },
  },
  head() {
    return {
      title: this.$siteTitle(this.topic.title),
    }
  },
}
</script>

<style lang="scss" scoped>
.user-info {
  background: #fff;
  padding: 0;
  margin: 0 0 10px 0;

  .base {
    padding: 10px;
    text-align: center;

    .avatar {
      min-width: 80px;
      min-height: 80px;
      width: 80px;
      height: 80px;
    }
    .vue-avatar--wrapper {
      margin-left: 80px;
    }

    .username {
      font-size: 15px;
      font-weight: 700;
      a:hover {
        text-decoration: underline;
      }
      span.level-0 {
        font-size: 13px;
      }
      span.level-10 {
        font-size: 13px;
        color: red;
      }
    }

    .description {
      font-size: 13px;
      margin-top: 5px;
      overflow: hidden;
      word-break: break-all;
      text-overflow: ellipsis;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      display: -webkit-box;
    }
  }

  .extra {
    padding: 0 10px;
    background: rgba(0, 0, 0, 0.01);
    border-top: 1px solid #f5f5f5;
    ul.basic {
      display: flex;
      li {
        width: 100%;
        text-align: center;
        span {
          font-size: 13px;
          font-weight: 400;
          color: #868e96;
        }
      }
    }
    .other {
      border-top: 1px solid #f5f5f5;
      line-height: 30px;
      text-align: center;
      color: #868e96;
      font-size: 14px;
      span {
        margin-left: 5px;
        color: #000;
      }
    }
  }
}
</style>
