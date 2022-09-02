<template>
  <div class="comment-component main-content">
    <load-more
      ref="commentsLoadMore"
      v-if="commentsPage"
      v-slot="{ results }"
      :init-data="commentsPage"
      :params="{ entityType: entityType, entityId: entityId }"
      url="/api/comments"
    >
      <ul class="comments">
        <li
          v-for="(comment, index) in results"
          :key="comment.commentId"
          class="comment"
          itemprop="comment"
          itemscope
          itemtype="http://schema.org/Comment"
        >
          <div class="comment-avatar">
            avatar
          </div>
          <div class="comment-meta">
            <span
              class="comment-username"
              itemprop="creator"
              itemscope
              itemtype="http://schema.org/Person"
            >
              <a :href="'/user/' + comment.user.id" itemprop="name">
                {{ comment.user.username }}
              </a>
            </span>
            <span class="comment-floor">#{{ index + 1 }}</span>
            <span class="comment-time">
              <time
                :datetime="
                  comment.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')
                "
                itemprop="datePublished"
                >{{ comment.createTime | prettyDate }}</time
              >
            </span>
            <span class="comment-reply">
              <a @click="reply(comment)">
                <i class="iconfont icon-reply"></i> 回复
              </a>
            </span>
          </div>
          <div class="comment-content content">
            <blockquote v-if="comment.quote" class="comment-quote">
              <div class="comment-quote-user">
                <img
                  v-lazy="comment.quote.user.avatar"
                  class="avatar size-20"
                />
                <a class="quote-username">{{ comment.quote.user.username }}</a>
                <span class="quote-time">
                  {{ comment.quote.createTime | prettyDate }}
                </span>
              </div>
              <div
                v-html="comment.quote.content"
                v-lazy-container="{ selector: 'img' }"
                itemprop="text"
              />
            </blockquote>
            <p
              v-html="comment.content"
              v-lazy-container="{ selector: 'img' }"
            />
          </div>
        </li>
      </ul>
    </load-more>
    <div ref="commentTitle" class="comment-title">
      发起回帖<span v-if="commentCount > 0">({{ commentCount }})</span>
    </div>
    <div class="comment-form">
      <div v-if="isLogin" class="comment-create">
        <div ref="commentEditor" class="comment-input-wrapper">
          <div v-if="quote" class="comment-quote-info">
            <i class="iconfont icon-reply"></i> 回复：
            <label v-text="quote.user.username" />
            <i
              @click="cancelReply"
              class="iconfont icon-close"
              title="取消回复"
            />
          </div>
          <markdown-editor
            ref="mdEditor"
            v-model="content"
            @submit="ctrlEnterCreate"
            editor-id="createEditor"
            height="150px"
            placeholder="请发表你的观点..."
          />
        </div>
        <div class="comment-button-wrapper">
          <button @click="create" v-text="btnName" class="button is-light" />
        </div>
      </div>
      <div v-else class="comment-not-login">
        <div class="comment-login-div">
          请
          <a @click="toLogin" style="font-weight: 700;">登录</a>后发表观点
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import utils from '~/common/utils'
import LoadMore from '~/components/LoadMore'
import MarkdownEditor from '~/components/MarkdownEditor'
export default {
  name: 'Comment',
  components: {
    LoadMore,
    MarkdownEditor
  },
  props: {
    entityType: {
      type: String,
      default: '',
      required: true
    },
    entityId: {
      type: Number,
      default: 0,
      required: true
    },
    commentsPage: {
      type: Object,
      default() {
        return {}
      }
    },
    commentCount: {
      type: Number,
      default: 0
    }
  },
  data() {
    return {
      content: '', // 内容
      sending: false, // 发送中
      quote: null // 引用的对象
    }
  },
  computed: {
    btnName() {
      return this.sending ? '正在发表...' : '回复 (Ctrl/Cmd + Enter)'
    },
    user() {
      return this.$store.state.auth.currentUser
    },
    isLogin() {
      return this.$store.state.auth.currentUser != null
    }
  },
  methods: {
    ctrlEnterCreate() {
      this.create()
    },
    async create() {
      if (!this.content || this.content.length < 3) {
        this.$toast.error('请输入评论内容')
        return
      }
      if (this.sending) {
        console.log('正在发送中，请不要重复提交...')
        return
      }
      this.sending = true
      try {
        const data = await this.$axios.post('/api/comments', {
          entityType: this.entityType,
          entityId: this.entityId,
          content: this.content,
          quoteId: this.quote ? this.quote.commentId : ''
        })
        this.$refs.commentsLoadMore.pushResults(data)
        this.content = ''
        this.$refs.mdEditor.clear()
        this.quote = null
      } catch (e) {
        console.error(e)
        this.$toast.error('评论失败：' + (e.message || e))
      } finally {
        this.sending = false
      }
    },
    reply(quote) {
      if (!this.isLogin) {
        utils.toSignin()
      }
      this.quote = quote
      this.$refs.commentTitle.scrollIntoView({
        block: 'start',
        behavior: 'smooth'
      })
    },
    cancelReply() {
      this.quote = null
    },
    autoHeight() {
      const elem = this.$refs.commentEditor
      elem.style.height = 'auto'
      elem.scrollTop = 0 // 防抖动
      elem.style.height = elem.scrollHeight + 'px'
    },
    toLogin() {
      utils.toSignin()
    }
  }
}
</script>

<style lang="scss" scoped>
.comment-component {
  padding-top: 10px;

  .comment-title {
    font-size: 16px;
    font-weight: 700;
    border-top: 1px solid #f0f0f0;
    // margin: 8px 0;
    padding: 10px 0 5px 0;
  }

  .comment-form {
    .comment-create {
      /*border: 1px solid #f0f0f0;*/
      border-radius: 4px;
      margin-bottom: 10px;
      overflow: hidden;
      position: relative;
      padding: 0;
      box-sizing: border-box;

      .comment-quote-info {
        font-size: 14px;
        color: #000;
        margin-bottom: 5px;

        i.icon-close {
          font-size: 14px !important;
          line-height: 14px;
          color: #448ef6;
          cursor: pointer;
        }

        i.icon-close:hover {
          color: red;
        }
      }

      .comment-input-wrapper {
        margin-bottom: 10px;
      }

      .comment-button-wrapper {
        button {
          float: left;
        }
      }
    }

    .comment-not-login {
      border: 1px solid #f0f0f0;
      border-radius: 0px;
      margin-bottom: 10px;
      overflow: hidden;
      position: relative;
      padding: 10px;
      box-sizing: border-box;

      .comment-login-div {
        color: #d5d5d5;
        cursor: pointer;
        border-radius: 3px;
        padding: 0 10px;

        a {
          margin-left: 10px;
          margin-right: 10px;
        }
      }
    }
  }

  .comments {
    .comment {
      padding: 8px 0;
      overflow: hidden;

      &:not(:last-child) {
        border-bottom: 1px dashed #d1d1d1;
      }

      .comment-avatar {
        float: left;
        padding: 3px;
        margin-right: 5px;

        .avatar {
          min-width: 36px;
          min-height: 36px;
          width: 36px;
          height: 36px;
        }
      }

      .comment-meta {
        position: relative;
        height: 36px;

        .comment-username {
          position: relative;
          font-size: 14px;
          font-weight: 800;
          margin-right: 5px;
          cursor: pointer;
          color: #1abc9c;
          text-decoration: none;
          display: inline-block;
        }
        .comment-floor {
          font-size: 14px;
          color: #7aa87a;
        }
        .comment-time {
          font-size: 12px;
          color: #999999;
          line-height: 1;
          display: inline-block;
          position: relative;
        }

        .comment-reply {
          float: right;
          font-size: 12px;
        }
      }

      .comment-content {
        word-wrap: break-word;
        word-break: break-all;
        text-align: justify;
        color: #4a4a4a;
        font-size: 14px;
        line-height: 1.6;
        position: relative;
        padding-left: 45px;
        margin-top: -5px;
      }

      .comment-quote {
        font-size: 12px;
        padding: 10px 10px;
        border-left: 2px solid #5978f3;

        &::after {
          content: '\201D';
          font-family: Georgia, serif;
          font-size: 60px;
          font-weight: bold;
          color: #aaa;
          position: absolute;
          right: 0px;
          top: -18px;
        }

        .comment-quote-user {
          display: flex;

          .quote-username {
            line-height: 20px;
            font-weight: 700;
            margin-left: 5px;
          }

          .quote-time {
            line-height: 20px;
            margin-left: 5px;
            color: rgba(134, 142, 150, 0.8);
          }
        }
      }
    }
  }
}
</style>
