<template>
  <section class="main-container">
    <div class="container">
      <div class="error">
        <div class="description">
          <span v-if="error.statusCode === 404">{{
            error.message ? error.message : '页面没找到'
          }}</span>
          <span v-if="error.statusCode === 403">{{
            error.message ? error.message : 'forbidden'
          }}</span>
          <span v-else>{{ error.statusCode }} 页面异常</span>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  props: {
    error: {
      type: Object,
      default: null
    }
  }
}
</script>

<style lang="scss" scoped>
.error {
  background: #fff;
  text-align: center;
  vertical-align: center;
  padding: 100px 0;

  .description {
    margin-top: 30px;
    span {
      font-size: 18px;
      font-weight: bold;
      line-height: 22px;
      color: rgb(230, 76, 76);
    }
  }
}
</style>
