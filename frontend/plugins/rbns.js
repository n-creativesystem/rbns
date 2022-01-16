const tenants = window.rbnsBootstrapData.tenants || []
export default {
  install: (Vue) => {
    Vue.prototype.$rbns = {
      bootstrap: window.rbnsBootstrapData,
      user: window.rbnsBootstrapData.user || {},
      tenants: tenants,
      isTenant: tenants.length > 0
    }
  }
}