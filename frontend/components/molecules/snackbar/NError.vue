<template>
  <n-snackbar
    v-model="innerValue"
    :message="message"
    :position="position"
    color="error"
    v-bind="$attrs"
  ></n-snackbar>
</template>

<script>
  export default {
    name: 'NSnackbarError',
    props: {
      value: Boolean,
      message: String,
      position: {
        type: Object,
        default: () => {
          return {
            top: true,
            right: true,
            bottom: false,
            left: false,
          }
        },
      },
    },
    computed: {
      innerValue: {
        get() {
          return this.value
        },
        set(val) {
          this.$emit('input', val)
        },
      },
    },
  }
</script>

<style></style>
