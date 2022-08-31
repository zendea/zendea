<template>
  <el-dialog
    :visible.sync="isShowNodes"
    v-if="isShowNodes"
    width="80%"
    title="节点列表"
  >
    <el-table
      :data="results"
      v-loading="listLoading"
      highlight-current-row
      border
    >
      <el-table-column prop="id" label="编号"></el-table-column>
      <el-table-column prop="name" label="名称"></el-table-column>
      <el-table-column prop="description" label="描述"></el-table-column>
      <el-table-column prop="topicCount" label="话题数"></el-table-column>
      <el-table-column prop="sortNo" label="排序"></el-table-column>
      <el-table-column prop="createTime" label="创建时间">
        <template slot-scope="scope">{{
          scope.row.createTime | formatDate
        }}</template>
      </el-table-column>
    </el-table>

    <div class="pagebar">
      <el-pagination
        :page-sizes="[20, 50, 100, 300]"
        @current-change="handlePageChange"
        @size-change="handleLimitChange"
        :current-page="page.page"
        :page-size="page.limit"
        :total="page.total"
        layout="total, sizes, prev, pager, next, jumper"
      >
      </el-pagination>
    </div>
  </el-dialog>
</template>

<script>
export default {
  data() {
    return {
      isShowNodes: false,
      sectionId: 0,
      results: [],
      listLoading: false,
      page: {},
      filters: {}
    }
  },
  mounted() {},
  methods: {
    async showNodes(sectionId) {
      this.sectionId = sectionId
      this.isShowNodes = true
      await this.list()
    },
    async list() {
      const me = this
      me.listLoading = true
      const params = {}
      params.params = Object.assign(me.filters, {
        page: me.page.page,
        limit: me.page.limit
      })

      try {
        const data = await this.$axios.get(
          '/api/admin/sections/' + me.sectionId + '/nodes'
        )
        this.results = data.results
        this.page = data.page
      } catch (err) {
        this.$notify.error({ title: '错误', message: err.message || err })
      } finally {
        this.listLoading = false
      }
    },
    handlePageChange(val) {
      this.page.page = val
      this.list()
    },
    handleLimitChange(val) {
      this.page.limit = val
      this.list()
    }
  }
}
</script>

<style lang="scss" scoped></style>
