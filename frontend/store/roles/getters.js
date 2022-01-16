/*
export function someGetter (state) {
}
*/

export const list = (state) => {
  return state.roles || []
}

export const listPermissions = (state) => {
  return state.rolePermissions || []
}