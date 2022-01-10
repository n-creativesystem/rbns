import Urls from './api'

export default {
  install: (Vue) => {
    const urls = {
      ...Urls
    }
    Vue.prototype.$urls = urls
  }
}