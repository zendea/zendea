<template>
  <div class="">
    <div v-if="show" @click="backToTop" class="back-to-top">
      回顶部
    </div>
  </div>
</template>

<script>
export default {
  name: 'BackToTop',
  props: {
    leaveTop: {
      type: Number,
      default: 300
    }
  },
  data() {
    return {
      show: false,
      timer: '',
      speed: 10
    }
  },
  mounted() {
    window.addEventListener('scroll', () => {
      this.show = this.getScrollTop() > this.leaveTop
    })
  },
  methods: {
    getScrollTop() {
      return document.documentElement.scrollTop || document.body.scrollTop
    },
    backToTop() {
      const currentScroll = this.getScrollTop()
      if (currentScroll > 0) {
        window.scrollTo(0, currentScroll - this.speed)
        // 加速
        this.speed += 2
        this.timer = requestAnimationFrame(() => {
          this.backToTop()
        })
      } else {
        this.speed = 10
        cancelAnimationFrame(this.timer)
      }
    }
  }
}
</script>

<style scoped lang="scss">
@import '@/assets/styles/variables';
.back-to-top {
  position: fixed;
  right: 40px;
  bottom: 80px;
  padding: 10px;
  width: 36px;
  color: #999;
  text-align: center;
  line-height: 1.25;
  border: 1px solid #ddd;
  background: #fff;
  cursor: pointer;
  z-index: 9999;
  &:hover {
    background: #f3f3f3;
  }
}
</style>
