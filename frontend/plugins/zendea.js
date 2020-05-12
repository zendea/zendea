import Vue from 'vue'

Vue.use({
  install(Vue, options) {
    Vue.prototype.$siteTitle = function(subTitle) {
      const siteTitle = this.$store.getters['config/siteTitle'] || ''
      if (subTitle) {
        return subTitle + (siteTitle ? ' | ' + siteTitle : '')
      }
      return siteTitle
    }

    Vue.prototype.$siteDescription = function() {
      return this.$store.getters['config/siteDescription']
    }

    Vue.prototype.$siteKeywords = function() {
      return this.$store.getters['config/siteKeywords']
    }

    Vue.prototype.$statUserCount = function() {
      return this.$store.getters['config/statUserCount']
    }

    Vue.prototype.$statTopicCount = function() {
      return this.$store.getters['config/statTopicCount']
    }

    Vue.prototype.$statCommentCount = function() {
      return this.$store.getters['config/statCommentCount']
    }
  }
})
