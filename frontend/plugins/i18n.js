import Vue from 'vue'
import VueI18n from 'vue-i18n'
import * as locales from '../locale'


Vue.use(VueI18n)

export const i18n = new VueI18n({
  locale: 'ja',
  messages: locales
})

export default {
  install: (Vue) => {
    Vue.prototype.$t = function (key, ...values) {
      const t = i18n.t(key, values)
      if (typeof t === 'string') {
        return t
      }
      throw new Error('not found key.')
    }
  }
}