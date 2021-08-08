import Vue from 'vue';
import Vuetify from 'vuetify/lib/framework';
import { i18n } from './i18n'

Vue.use(Vuetify);

export default new Vuetify({
  lang: {
    current: 'ja',
    t: (key, ...params) => i18n.t(key, params)
  }
});
