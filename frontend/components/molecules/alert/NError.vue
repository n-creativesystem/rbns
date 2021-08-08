<template>
  <n-alert
    v-model="innerValue"
    :message="message"
    :btn-text="btnText"
    :outline="outline"
    :dense="dense"
    type="error"
    @click="onClick"
  ></n-alert>
</template>

<script>
  export default {
    name: 'NAlertError',
    props: {
      value: Boolean,
      message: String,
      btnText: String,
      outlined: {
        type: Boolean,
        default: true,
      },
      dense: {
        type: Boolean,
        default: true,
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
    methods: {
      onClick(e) {
        this.$emit('click', e)
      },
    },
  }
</script>

<style></style>
