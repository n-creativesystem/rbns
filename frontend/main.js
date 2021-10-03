import Vue from 'vue'
import App from './App.vue'
import plugin from '@plugins'
import vuetify from '@plugins/vuetify'
import router from '@plugins/router'
import './css/rbac.css'
import store from './store'
import { i18n } from '@plugins/i18n'

(async () => {
  try {
    const res = await fetch('settings.json')
    const data = await res.json()
    window.VUE_APP_RBNS_API_ENDPOINT = data.base_url
  } catch (_) {
    console.log();
  }
  Vue.use(plugin)

  Vue.config.productionTip = false

  new Vue({
    vuetify,
    router,
    store,
    i18n,
    render: h => h(App)
  }).$mount("#app")
})()
