export const state = () => ({
  currentUser: null,
  userJwt: null
})

export const mutations = {
  setCurrentUser(state, user) {
    state.currentUser = user
  },
  setUserJwt(state, jwt) {
    state.userJwt = jwt
  }
}

export const actions = {
  // 登录成功
  loginSuccess(context, { token, user }) {
    console.log(token)
    this.$cookies.set('jwt', token, { maxAge: 86400 * 7, path: '/' })
    context.commit('setUserJwt', token)
  },

  // 获取当前登录用户
  async getCurrentUser(context) {
    const ret = await this.$axios.get('/api/user/current')
    if (ret.success) {
      const user = ret.data
      context.commit('setCurrentUser', user)
      return user
    }
  },

  // 登录
  async signin(context, { captchaId, captchaCode, username, password }) {
    const ret = await this.$axios.post('/api/auth/login', {
      captchaId,
      captchaCode,
      username,
      password
    })
    if (ret.success === true) {
      context.dispatch('loginSuccess', ret.data)
    } else {
      throw ret
    }
  },

  // github登录
  async signinByGithub(context, { code, state }) {
    const ret = await this.$axios.get('/api/oauth/github/callback', {
      params: {
        code,
        state
      }
    })
    if (ret.success === true) {
      context.dispatch('loginSuccess', ret.data)
    } else {
      throw ret
    }
  },

  // github登录
  async signinByGitee(context, { code, state }) {
    const ret = await this.$axios.get('/api/oauth/gitee/callback', {
      params: {
        code,
        state
      }
    })
    if (ret.success === true) {
      context.dispatch('loginSuccess', ret.data)
    } else {
      throw ret
    }
  },

  // qq登录
  async signinByQQ(context, { code, state }) {
    const ret = await this.$axios.get('/api/oauth/qq/callback', {
      params: {
        code,
        state
      }
    })
    if (ret.success === true) {
      context.dispatch('loginSuccess', ret.data)
    } else {
      throw ret
    }
  },

  // 注册
  async signup(
    context,
    { captchaId, captchaCode, nickname, username, email, password, rePassword }
  ) {
    const ret = await this.$axios.post('/api/auth/signup', {
      captchaId,
      captchaCode,
      nickname,
      username,
      email,
      password,
      rePassword
    })
    context.dispatch('loginSuccess', ret)
    return ret.user
  },

  // 退出登录
  signout(context) {
    context.commit('setCurrentUser', null)
    this.$cookies.remove('jwt')
  }
}
