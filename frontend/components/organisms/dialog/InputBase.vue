<template>
  <v-card>
    <n-form @submit="onClick" ref="form">
      <v-card-text>
        <slot></slot>
      </v-card-text>
      <v-card-actions class="justify-end">
        <btn-tfn type="submit">{{ text }}</btn-tfn>
      </v-card-actions>
    </n-form>
  </v-card>
</template>

<script>
  export default {
    name: 'dialogInputBase',
    data() {
      return {
        text: '作成',
      }
    },
    props: {
      activator: String,
    },
    methods: {
      onClick() {
        this.$emit('submit')
      },
      reset() {
        this.$refs.form.reset()
      },
    },
  }
</script>

<style></style>
