import Urls from '@assets/urls.json'

export default {
  install: (Vue) => {
    const urls = {
      ...Urls
    }
    Vue.prototype.$urls = urls
  }
}