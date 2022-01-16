import Axios from './axios'
import Components from '../components'
import Urls from './urls'
import i18n from './i18n'
import './utils'
import LoginMixin from './login'
import rbns from './rbns'

const plugins = [
  Axios,
  Components,
  Urls,
  i18n,
  rbns
]

export default {
  install: (Vue) => {
    Vue.mixin(LoginMixin)
    plugins.forEach(plugin => {
      Vue.use(plugin)
    })
  }
}