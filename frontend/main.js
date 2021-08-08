import Vue from 'vue'
import App from './App.vue'
import plugin from '@plugins'
import vuetify from '@plugins/vuetify'
import router from '@plugins/router'
import './css/rbac.css'
import store from './store'
import { i18n } from '@plugins/i18n'

Vue.use(plugin)

Vue.config.productionTip = false

new Vue({
  vuetify,
  router,
  store,
  i18n,
  render: h => h(App)
}).$mount('#app')
