import Vue from 'vue'
import Vuex from 'vuex'
import loading from './loading'
import notification from './notification'
import organizations from './organizations'
import roles from './roles'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    loading,
    notification,
    organizations,
    roles
  },

  strict: process.env.DEV
})
