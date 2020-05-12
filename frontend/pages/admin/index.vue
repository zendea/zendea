<template>
  <div class="page-container">
    <h1 class="index-title">控制台</h1>
    <div class="widget">
      <div class="widget-header">构建信息</div>
      <div class="widget-content">
        <div>App name: {{ systeminfo.appName }}</div>
        <div>App version: {{ systeminfo.appVersion }}</div>
        <div>Build time: {{ systeminfo.buildTime }}</div>
        <div>Build commit: {{ systeminfo.buildCommit }}</div>
      </div>
    </div>
    <div class="widget">
      <div class="widget-header">运行信息</div>
      <div class="widget-content">
        <div>Runtime: {{ systeminfo.upTime }}</div>
        <div>OS: {{ systeminfo.os }}</div>
        <div>Arch: {{ systeminfo.arch }}</div>
        <div>CPU: {{ systeminfo.numCpu }}</div>
        <div>Go version: {{ systeminfo.goversion }}</div>
      </div>
    </div>
    <div class="widget">
      <div class="widget-header">Zendea官方信息</div>
      <div class="widget-content">
        <div>
          Official Site:
          <a href="http://zendea.com/" target="_blank">http://zendea.com</a>
        </div>
        <div>
          Github:
          <a href="https://github.com/zendea/zendea" target="_blank"
            >https://github.com/zendea/zendea</a
          >
        </div>
        <div>
          Gitee:
          <a href="https://gitee.com/zendea/zendea" target="_blank"
            >https://gitee.com/zendea/zendea</a
          >
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  layout: 'admin',
  async asyncData({ $axios }) {
    const systeminfo = await $axios.get('/api/admin/dashboard/systeminfo')
    return {
      systeminfo
    }
  }
}
</script>

<style lang="scss" scoped></style>
