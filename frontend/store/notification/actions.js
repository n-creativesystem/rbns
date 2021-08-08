export const show = ({ commit }, { type, message }) => {
  const typ = type.toLowerCase()
  switch (typ) {
    case 'error':
    case 'info':
    case 'success':
      commit(typ, true)
      break;
    default:
      break;
  }
  commit('set', message)
}

export const hide = ({ commit }, { type }) => {
  const typ = type.toLowerCase()
  switch (typ) {
    case 'error':
    case 'info':
    case 'success':
      commit(typ, false)
      break;
    default:
      break;
  }
  commit('set', '')
}
