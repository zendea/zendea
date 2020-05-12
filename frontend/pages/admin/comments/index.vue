<template>
  <section v-loading="listLoading" class="page-container">
    <div class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.id" placeholder="编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.userId" placeholder="用户编号"></el-input>
        </el-form-item>
        <el-form-item>
          <el-select
            v-model="filters.entityType"
            @change="list"
            clearable
            placeholder="请选择类型"
          >
            <el-option label="话题" value="topic"></el-option>
            <el-option label="文章" value="article"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="filters.entityId" placeholder="上级编号">
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-select
            v-model="filters.status"
            @change="list"
            clearable
            placeholder="请选择状态"
          >
            <el-option label="正常" value="0"></el-option>
            <el-option label="删除" value="1"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button v-on:click="list" type="primary">查询</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="main-content">
      <ul class="comments">
        <li v-for="item in results" :key="item.id">
          <div class="comment-item">
            <div
              :style="{ backgroundImage: 'url(' + item.user.avatar + ')' }"
              class="avatar"
            ></div>
            <div class="content">
              <div class="meta">
                <span class="nickname"
                  ><a :href="'/user/' + item.user.id" target="_blank">{{
                    item.user.nickname
                  }}</a></span
                >
                <span class="create-time"
                  >@{{ item.createTime | formatDate }}</span
                >
                <span v-if="item.entityType === 'article'">
                  <a :href="'/article/' + item.entityId" target="_blank"
                    >文章：{{ item.entityId }}</a
                  >
                </span>

                <span v-if="item.entityType === 'topic'">
                  <a :href="'/topic/' + item.entityId" target="_blank"
                    >话题：{{ item.entityId }}</a
                  >
                </span>
              </div>
              <div v-html="item.content" class="summary"></div>
              <div class="tools">
                <span v-if="item.status === 1" class="item info">已删除</span>
                <a @click="handleDelete(item)" class="item">删除</a>
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>

    <div class="pagebar">
      <el-pagination
        :page-sizes="[20, 50, 100, 300]"
        @current-change="handlePageChange"
        @size-change="handleLimitChange"
        :current-page="page.page"
        :page-size="page.limit"
        :total="page.total"
        layout="total, sizes, prev, pager, next, jumper"
      ></el-pagination>
    </div>
  </section>
</template>

<script>
export default {
  layout: 'admin',
  data() {
    return {
      results: [],
      listLoading: false,
      page: {},
      filters: {
        id: '',
        userId: '',
        entityType: '',
        entityId: '',
        status: ''
      },
      selectedRows: [],

      addForm: {
        userId: '',
        entityType: '',
        entityId: '',
        content: '',
        quoteId: '',
        status: '',
        createTime: ''
      },
      addFormVisible: false,
      addFormRules: {},
      addLoading: false,

      editForm: {
        id: '',
        userId: '',
        entityType: '',
        entityId: '',
        content: '',
        quoteId: '',
        status: '',
        createTime: ''
      },
      editFormVisible: false,
      editFormRules: {},
      editLoading: false
    }
  },
  mounted() {
    this.list()
  },
  methods: {
    list() {
      const me = this
      me.listLoading = true
      const params = {}
      params.params = Object.assign(me.filters, {
        page: me.page.page,
        limit: me.page.limit
      })
      this.$axios
        .get('/api/admin/comments', params)
        .then((data) => {
          me.results = data.results
          me.page = data.page
        })
        .finally(() => {
          me.listLoading = false
        })
    },
    handlePageChange(val) {
      this.page.page = val
      this.list()
    },
    handleLimitChange(val) {
      this.page.limit = val
      this.list()
    },
    handleSelectionChange(val) {
      this.selectedRows = val
    },
    handleDelete(row) {
      const me = this
      this.$axios
        .delete(`/api/admin/comments/${row.id}`)
        .then((data) => {
          me.$message.success('删除成功')
          me.list()
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    }
  }
}
</script>

<style scoped lang="scss">
.comments {
  list-style: none;
  padding: 0px;

  li {
    border-bottom: 1px solid #f2f2f2;
    padding-top: 10px;
    padding-bottom: 10px;

    .comment-item {
      display: flex;

      .avatar {
        min-width: 40px;
        min-height: 40px;
        width: 40px;
        height: 40px;
        border-radius: 50%;
        margin-right: 10px;
        background-repeat: no-repeat;
        background-size: contain;
        background-position: center;
      }

      .content {
        width: 100%;
        .meta {
          span {
            &:not(:last-child) {
              margin-right: 5px;
            }

            font-size: 13px;
            color: #999;
            font-weight: bold;

            &.nickname {
              color: #1a1a1a;
              font-size: 14px;
              font-weight: bold;
            }

            &.create-time {
              color: #999;
              font-size: 13px;
              font-weight: normal;
            }
          }
        }

        .summary {
          font-size: 15px;
          color: #555;
        }

        .tools {
          float: right;
          .item {
            color: blue;
            cursor: pointer;
            &:not(:last-child) {
              margin-right: 10px;
            }

            &.info {
              color: red;
            }
          }
        }
      }
    }
  }
}
</style>
