<template>
  <section class="page-container">
    <div class="toolbar">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.name" placeholder="名称"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button v-on:click="list" type="primary">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="handleAdd" type="primary">新增</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      :data="results"
      v-loading="listLoading"
      @selection-change="handleSelectionChange"
      highlight-current-row
      border
      style="width: 100%;"
    >
      <el-table-column type="selection" width="55"></el-table-column>
      <el-table-column prop="id" label="编号"></el-table-column>
      <el-table-column prop="name" label="名称"></el-table-column>
      <el-table-column prop="sortNo" label="排序"></el-table-column>

      <el-table-column prop="createTime" label="创建时间">
        <template slot-scope="scope">{{
          scope.row.createTime | formatDate
        }}</template>
      </el-table-column>

      <el-table-column label="操作" width="250">
        <template slot-scope="scope">
          <el-button
            @click="showNodes(scope.$index, scope.row)"
            type="success"
            size="small"
            >节点</el-button
          >
          <el-button @click="handleEdit(scope.$index, scope.row)" size="small"
            >编辑</el-button
          >
          <el-button
            @click="deleteSubmit(scope.row)"
            type="danger"
            size="small"
          >
            删除
          </el-button>
        </template>
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

    <el-dialog
      :visible.sync="addFormVisible"
      :close-on-click-modal="false"
      title="新增"
    >
      <el-form ref="addForm" :model="addForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="addForm.name"></el-input>
        </el-form-item>
        <el-form-item label="排序">
          <el-input v-model="addForm.sortNo"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addFormVisible = false">取消</el-button>
        <el-button
          @click.native="addSubmit"
          :loading="addLoading"
          type="primary"
          >提交</el-button
        >
      </div>
    </el-dialog>

    <el-dialog
      :visible.sync="editFormVisible"
      :close-on-click-modal="false"
      title="编辑"
    >
      <el-form ref="editForm" :model="editForm" label-width="80px">
        <el-input v-model="editForm.id" type="hidden"></el-input>
        <el-form-item label="名称">
          <el-input v-model="editForm.name"></el-input>
        </el-form-item>
        <el-form-item label="排序">
          <el-input v-model="editForm.sortNo"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="editFormVisible = false">取消</el-button>
        <el-button
          @click.native="editSubmit"
          :loading="editLoading"
          type="primary"
          >提交</el-button
        >
      </div>
    </el-dialog>
    <section-nodes ref="sectionNodes" />
  </section>
</template>

<script>
import SectionNodes from './section-nodes'
export default {
  layout: 'admin',
  components: { SectionNodes },
  data() {
    return {
      results: [],
      sectionNodes: [],
      listLoading: false,
      page: {},
      filters: {},
      selectedRows: [],

      addForm: {
        name: '',
        sortNo: '',
        createTime: ''
      },
      addFormVisible: false,
      addLoading: false,

      editForm: {
        id: '',
        name: '',
        sortNo: '',
        createTime: ''
      },
      editFormVisible: false,
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
        .get('/api/admin/sections', params)
        .then((data) => {
          me.results = data.results
          me.page = data.page
        })
        .finally(() => {
          me.listLoading = false
        })
    },
    showNodes(index, row) {
      this.$refs.sectionNodes.showNodes(row.id)
    },
    handlePageChange(val) {
      this.page.page = val
      this.list()
    },
    handleLimitChange(val) {
      this.page.limit = val
      this.list()
    },
    handleAdd() {
      this.addForm = {
        name: ''
      }
      this.addFormVisible = true
    },
    addSubmit() {
      const me = this
      this.$axios
        .post('/api/admin/sections', this.addForm)
        .then((data) => {
          me.$message({ message: '提交成功', type: 'success' })
          me.addFormVisible = false
          me.list()
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },
    handleEdit(index, row) {
      const me = this
      this.$axios
        .get('/api/admin/sections/' + row.id)
        .then((data) => {
          me.editForm = Object.assign({}, data)
          me.editFormVisible = true
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },
    editSubmit() {
      const me = this
      this.$axios
        .put('/api/admin/sections/' + me.editForm.id, me.editForm)
        .then((data) => {
          me.$message({ message: '编辑成功', type: 'success' })
          me.list()
          me.editFormVisible = false
        })
        .catch((rsp) => {
          me.$notify.error({ title: '错误', message: rsp.message })
        })
    },
    async deleteSubmit(row) {
      await this.$confirm('是否确认删除该分类?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then((confirm) => {
          try {
            this.$axios.delete('/api/admin/sections/' + row.id)
            this.$message({ message: '删除成功', type: 'success' })
            this.list()
          } catch (err) {
            this.$notify.error({ title: '错误', message: err.message || err })
          }
        })
        .catch((cancel) => {
          console.log('cancel')
        })
    },
    handleSelectionChange(val) {
      this.selectedRows = val
    }
  }
}
</script>

<style lang="scss" scoped></style>
