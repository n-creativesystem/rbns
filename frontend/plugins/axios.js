import Axios from 'axios'

const baseUrl = document.getElementsByTagName("base")[0]

export default {
  install: (Vue) => {
    const axios = Axios.create({
      baseURL: baseUrl.href || '/'
    })
    axios.interceptors.request.use(config => {
      if (!config.headers) {
        config.headers = {}
      }
      return config
    })
    Vue.prototype.$axios = axios
  }
}
