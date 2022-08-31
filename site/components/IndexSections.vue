<template>
  <div class="index-sections widget">
    <div class="widget-header">讨论节点分类导航</div>
    <div class="widget-content">
      <div class="row node-list">
        <div v-for="section in sections" class="node media2">
          <label class="media2-left">{{ section.name }}</label>
          <span class="nodes media2-body">
            <span v-for="node in section.nodes" class="name">
              <a :href="'/topics/' + node.nodeId">{{ node.name }}</a>
            </span>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    sections: {
      type: Array,
      default() {
        return []
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.index-sections {
  .node-list {
    display: flex;
    flex-wrap: wrap;
    margin-left: -15px;
    margin-right: -15px;
    .node {
      font-size: 14px;
      margin-bottom: 10px;
      margin-top: 0;
      label {
        display: inline-block;
        margin-bottom: 0.5rem;
        font-weight: normal;
        color: #aaa;
        text-align: right;
      }
      .name {
        margin-bottom: 10px;
        width: 100px;
        display: block;
        float: left;
        text-align: left;
        a {
          color: #333;
          &:hover {
            text-decoration: underline;
          }
        }
      }
    }
    .media2 {
      display: flex;
      align-items: flex-start;
      .media2-left {
        min-width: 130px;
        margin-right: 7.5px !important;
      }
      .media2-body {
        flex: 1;
      }
    }
  }
}
</style>
