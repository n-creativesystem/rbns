export const error = (state, show) => {
  state.error = show
}

export const info = (state, show) => {
  state.info = show
}

export const success = (state, show) => {
  state.success = show
}

export const set = (state, message) => {
  state.message = message
}