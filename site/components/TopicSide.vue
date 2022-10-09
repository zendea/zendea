<template>
  <div class="right-container">
    <post-btn-sidebar :current-node-id="currentNodeId" />
    <site-notice />
    <site-tip />
    <div v-if="scoreRank && scoreRank.length" class="widget">
      <div class="widget-header">
        <span class="widget-title">积分排行</span>
      </div>
      <div class="widget-content">
        <ul class="score-rank">
          <li v-for="user in scoreRank" :key="user.id">
            <a :href="'/user/' + user.id" class="score-user-avatar">avatar</a>
            <div class="score-user-info">
              <a :href="'/user/' + user.id">{{ user.username }}</a>
              <p>{{ user.topicCount }} 帖子 • {{ user.commentCount }} 评论</p>
            </div>
            <div class="score-rank-info">
              <span class="score-user-score">
                <i class="iconfont icon-dollar" /><span>{{ user.score }}</span>
              </span>
            </div>
          </li>
        </ul>
      </div>
    </div>
    <div v-if="links && links.length" class="widget">
      <div class="widget-header">
        <span>友情链接</span>
        <span class="slot"><a href="/links">查看更多&gt;&gt;</a></span>
      </div>
      <div class="widget-content links">
        <ul class="list-group">
          <li v-for="link in links" :key="link.linkId" class="list-group-item">
            <a
              :href="link.url"
              :title="link.title"
              class="link-title"
              rel="nofollow"
              target="_blank"
            >
              <p v-if="!link.logo" class="link-title">{{ link.title }}</p>
              <p v-if="!link.logo" class="link-summary">
                {{ link.summary }}
              </p>
              <p v-if="link.logo">
                <img v-if="link.logo" :src="link.logo" />
              </p>
            </a>
          </li>
        </ul>
      </div>
    </div>
    <site-stat :stat="stat" />
  </div>
</template>

<script>
import PostBtnSidebar from '~/components/PostBtnSidebar'
import SiteNotice from '~/components/SiteNotice'
import SiteTip from '~/components/SiteTip'
import SiteStat from '~/components/SiteStat'

export default {
  components: {
    PostBtnSidebar,
    SiteNotice,
    SiteTip,
    SiteStat,
  },
  props: {
    currentNodeId: {
      type: Number,
      default: 0,
    },
    links: {
      type: Array,
      default() {
        return null
      },
    },
    stat: {
      type: Object,
      default() {
        return {}
      },
    },
    scoreRank: {
      type: Array,
      default() {
        return null
      },
    },
  },
}
</script>

<style lang="scss" scoped>
.links {
  img {
    width: 130px;
  }
  .link-title {
    font-size: 14px;
  }
  .link-summary {
    font-size: 12px;
  }
}
.score-rank {
  li {
    display: flex;
    list-style: none;
    margin: 8px 0;
    font-size: 13px;
    position: relative;

    &:not(:last-child) {
      border-bottom: 1px solid #f7f7f7;
    }

    .score-user-avatar {
      min-width: 30px;
      .avatar {
        width: 30px;
        height: 30px;
      }
    }

    .score-user-info {
      width: 100%;
      margin-left: 5px;
      line-height: 1.4;
      font-size: 12px;
      a {
        font-weight: 700;
        &:hover {
          text-decoration: underline;
        }
      }
    }

    .score-rank-info {
      width: 120px;
      .score-user-score {
        float: right;
        border-radius: 12px;
        color: #778087;
        height: 21px;
        line-height: 21px;
        padding: 0 6px;
        text-shadow: 0 0 1px #fff;
        background-color: #f5f5f5;
        font-size: 0.75rem;
        align-items: center;
        i {
          margin-right: 3px;
        }
      }
    }
  }
}
</style>
