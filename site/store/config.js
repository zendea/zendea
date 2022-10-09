export const state = () => ({
  setting: {},
  appinfo: {},
})

export const mutations = {
  setConfig(state, config) {
    state.setting = config.setting
    state.appinfo = config.appinfo
  },
}

export const actions = {
  // 加载配置
  async loadConfig(context) {
    const ret = await this.$axios.get('/api/config/configs')
    context.commit('setConfig', ret)
    return ret
  },
}

export const getters = {
  siteTitle(state) {
    return state.setting.siteTitle || ''
  },
  siteDescription(state) {
    return state.setting.siteDescription || ''
  },
  siteKeywords(state) {
    return state.setting.siteKeywords || ''
  },
}
