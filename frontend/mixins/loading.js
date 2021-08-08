import Vue from 'vue'

export default Vue.extend({
  data() {
    return {
      internalLoadingCounter: 0
    }
  },
  methods: {
    onloading() {
      this.internalLoadingCounter++
      this.$store.dispatch('loading/loading')
      return this.internalLoadingCounter
    },
    unloading(counter) {
      if (counter && counter === this.internalLoadingCounter) {
        this.$store.dispatch('loading/unloading')
      }
    }
  },
  computed: {
    loading() {
      return this.$store.getters['loading/loading']
    }
  }
})