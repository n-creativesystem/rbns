export const loading = ({ commit }) => {
  commit('set', true)
}

export const unloading = ({ commit }) => {
  commit('set', false)
}