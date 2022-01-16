import Axios from 'axios'

const baseUrl = document.getElementsByTagName("base")[0]
export const axios = Axios.create({
  baseURL: baseUrl.href || '/'
})
axios.interceptors.request.use(config => {
  if (!config.headers) {
    config.headers = {}
  }
  config.headers['X-Request-By'] = location.origin
  return config
})

export default {
  install: (Vue) => {
    Vue.prototype.$axios = axios
  }
}
