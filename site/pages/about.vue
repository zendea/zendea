<template>
  <section class="main">
    <div class="container">
      <div class="main-container columns">
        <div class="column is-12">
          <div class="widget">
            <div class="widget-header">
              <a href="/">Zendea</a>
              › About <i class="el-icon-edit"></i>
            </div>
            <div class="widget-content content">
              <h2>About Zendea</h2>
              <p>
                Zendea discussion community, ask questions and share your ideas.
              </p>
              <p>
                Zendea 是创意工作者们的社区。
                这里目前汇聚了主要来自于互联网、游戏和媒体等行业的创意工作者们。
                Zendea 希望能够成为创意工作者们生活和事业的一部分。
              </p>
              <p>
                希望大家能够多多分享自己正在做的有趣事物、交流想法，在这里找到朋友甚至新的机会。
                并且，最重要的是，在这一切的过程中，保持对他人的友善。
              </p>

              <h2>Explore</h2>
              <p>
                为了获得访问 Zendea 的最佳体验，我们推荐使用
                <a href="http://chrome.google.com/" target="_blank"
                  >Google Chrome
                </a>
                或
                <a href="http://www.getfirefox.com/" target="_blank"
                  >Mozilla Firefox
                </a>
                浏览器
              </p>

              <h2>Tips</h2>
              <p>
                1. Zendea 不反对文章的原作者自己全文转载自己写的原创文章<br />
                2. “路过”、“沙发”之类的 0
                信息量回复会被自动规则阻挡或被管理员删除。与讨论主题完全无关的回复，尤其是在技术类主题下出现的话，会被管理员删除。
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  head() {
    return {
      title: this.$siteTitle('关于'),
    }
  },
}
</script>

<style lang="scss" scoped></style>
