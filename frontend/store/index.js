import Vue from 'vue'
import Vuex from 'vuex'
import loading from './loading'
import notification from './notification'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    loading,
    notification
  },

  strict: process.env.DEV
})
