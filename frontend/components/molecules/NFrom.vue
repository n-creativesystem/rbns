<template>
  <v-form @submit.prevent="onSubmit" ref="form">
    <slot></slot>
  </v-form>
</template>

<script>
  export default {
    name: 'NForm',
    methods: {
      onSubmit(e) {
        if (this.$refs.form.validate()) {
          this.$emit('submit', e)
        }
      },
      reset() {
        this.$refs.form.reset()
      },
    },
  }
</script>

<style></style>
