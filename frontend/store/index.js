export const state = () => ({
  baseURL: "",
  secret:  "",
})

export const mutations = {
  setBaseURL(state, text) {
    state.baseURL = text
  },
  setSecret(state, text) {
    state.secret = text
  },
}
