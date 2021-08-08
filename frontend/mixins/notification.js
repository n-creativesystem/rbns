import Vue from 'vue'

export default Vue.extend({
  methods: {
    __set_notification(show, type, message) {
      if (show) {
        this.$store.dispatch(`notification/show`, { type: type, message: message })
      } else {
        this.$store.dispatch('notification/hide', { type: type })
      }
    },
    showErrorNotification(message) {
      this.__set_notification(true, 'error', message)
    },
    showInfoNotification(message) {
      this.__set_notification(true, 'info', message)
    },
    showSuccessNotification(message) {
      this.__set_notification(true, 'success', message)
    },
    hideErrorNotification() {
      this.__set_notification(false, 'error')
    },
    hideInfoNotification() {
      this.__set_notification(false, 'info')
    },
    hideSuccessNotification() {
      this.__set_notification(false, 'success')
    },
  },
  computed: {
    errorNotification: {
      get() {
        return this.$store.getters['notification/error']
      },
      set() {
        this.hideErrorNotification()
      }
    },
    infoNotification: {
      get() {
        return this.$store.getters['notification/info']
      },
      set() {
        this.hideInfoNotification()
      }
    },
    successNotification: {
      get() {
        return this.$store.getters['notification/success']
      },
      set() {
        this.hideSuccessNotification()
      }
    },
    notificationMessage() {
      return this.$store.getters['notification/message']
    }
  }
})