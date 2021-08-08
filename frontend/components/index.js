import * as Atoms from './atoms'
import * as Molecules from './molecules'
import * as Organisms from './organisms'
import * as Templates from './templates'
import * as Pages from './pages'

const components = {
  ...Atoms,
  ...Molecules,
  ...Organisms,
  ...Pages,
  ...Templates
}

function install(Vue, components) {
  components && Object.keys(components).forEach(key => {
    const c = components[key]
    if (typeof c === 'function') {
      Vue.component(c.options.name, c)
    } else {
      const comp = Vue.extend(c)
      try {
        Vue.component(comp.options.name, comp)
      } catch (error) {
        install(Vue, c)
      }
    }
  })
}

export default {
  install: function (Vue) {
    install(Vue, components)
  }
}
