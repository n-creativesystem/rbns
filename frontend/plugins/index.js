import Axios from './axios'
import Components from '../components'
import Urls from './urls'
import i18n from './i18n'
import './utils'

const plugins = [
  Axios,
  Components,
  Urls,
  i18n
]

export default {
  install: (Vue) => {
    plugins.forEach(plugin => {
      Vue.use(plugin)
    })
  }
}