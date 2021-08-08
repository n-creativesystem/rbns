import Axios from 'axios'

const baseNode = document.querySelector('base')
const baseUrl = baseNode ? baseNode.href : '/'

export default {
  install: (Vue) => {
    const axios = Axios.create({
      baseURL: baseUrl
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
