import Axios from './axios'
import Components from '../components'
import Urls from './urls'
import i18n from './i18n'
import './utils'
import LoginMixin from './login'

const plugins = [
  Axios,
  Components,
  Urls,
  i18n
]

export default {
  install: (Vue) => {
    Vue.mixin(LoginMixin)
    plugins.forEach(plugin => {
      Vue.use(plugin)
    })
  }
}