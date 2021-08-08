<template>
  <div v-if="$vuetify.breakpoint.mobile" class="col-12 dense-4">
    <v-alert :outlined="outlined" :type="type" v-model="innerValue">
      {{ message }}
    </v-alert>
    <btn-tfn block @click.stop="onClick" :color="type">{{ btnText }}</btn-tfn>
  </div>
  <v-alert
    v-else
    :dense="dense"
    :outlined="outlined"
    :type="type"
    v-model="innerValue"
  >
    <div class="d-flex justify-space-between align-center">
      <div>
        {{ message }}
      </div>
      <btn-tfn @click.stop="onClick" :color="type">{{ btnText }}</btn-tfn>
    </div>
  </v-alert>
</template>

<script>
  export default {
    name: 'NAlert',
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
      type: String,
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
