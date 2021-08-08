import Vue from 'vue'
import loadingMixin from './loading'
import NotificationMixin from './notification'

export default Vue.extend({
  mixins: [loadingMixin, NotificationMixin],
  methods: {
    __axios_call(config) {
      const counter = this.onloading()
      return new Promise((resolve, reject) => {
        this.$axios(config)
          .then((response) => {
            resolve(response)
          }).catch((err) => {
            if (err.response && err.response.data && err.response.data.description) {
              this.showErrorNotification(err.response.data.description)
              reject(err)
            } else {
              this.showErrorNotification(err.toString())
              reject(err)
            }
          })
          .finally(() => this.unloading(counter))
      })
    },
    get(url, config) {
      if (!config) {
        config = {}
      }
      return this.__axios_call({
        method: 'get',
        url: url,
        ...config
      })
    },
    post(url, data, config) {
      if (!config) {
        config = {}
      }
      return this.__axios_call({
        method: 'post',
        url: url,
        data: data,
        ...config
      })
    },
    put(url, data, config) {
      if (!config) {
        config = {}
      }
      return this.__axios_call({
        method: 'put',
        url: url,
        data: data,
        ...config
      })
    },
    delete(url, config) {
      if (!config) {
        config = {}
      }
      return this.__axios_call({
        method: 'delete',
        url: url,
        ...config
      })
    },
  }
})