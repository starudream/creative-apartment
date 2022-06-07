export const state = () => ({
  isLogin: false,
  baseURL: "",
  secret:  "",
  version: "",
})

export const mutations = {
  clear(state) {
    state.isLogin = false
    state.baseURL = ""
    state.secret = ""
    state.version = ""
  },
  setLogin(state, isLogin) {
    state.isLogin = isLogin
  },
  setBaseURL(state, text) {
    state.baseURL = text
  },
  setSecret(state, text) {
    state.secret = text
  },
  setVersion(state, text) {
    state.version = text
  },
}
