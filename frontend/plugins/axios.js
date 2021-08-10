import Axios from 'axios'

export default {
  install: (Vue) => {
    const axios = Axios.create({
      baseURL: window.VUE_APP_RBNS_API_ENDPOINT || '/'
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
