<template>
  <div class="zendea">
    <nav-bar />

    <nuxt />

    <my-footer />
    <back-to-top />
  </div>
</template>

<script>
import NavBar from '~/components/NavBar'
import MyFooter from '~/components/MyFooter'
import BackToTop from '~/components/BackToTop'

export default {
  components: {
    NavBar,
    MyFooter,
    BackToTop
  }
}
</script>
