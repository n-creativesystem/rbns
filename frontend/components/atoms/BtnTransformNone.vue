<template>
  <v-btn
    :loading="loading"
    :color="color"
    class="tfn"
    v-bind="$attrs"
    v-on="$listeners"
  >
    <slot></slot>
  </v-btn>
</template>

<script>
  import loadingMixin from '@mixin/loading'
  export default {
    name: 'BtnTfn',
    mixins: [loadingMixin],
    props: {
      color: {
        type: String,
        default: 'info',
      },
    },
  }
</script>

<style></style>
