export const state = () => ({
  baseURL: "",
  secret:  "",
  version: "",
})

export const mutations = {
  clear(state) {
    state.baseURL = ""
    state.secret = ""
    state.version = ""
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
