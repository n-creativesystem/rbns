<template>
  <v-snackbar
    v-model="innerValue"
    :color="color"
    :top="position.top"
    :right="position.right"
    :bottom="position.bottom"
    :left="position.left"
    v-bind="$attrs"
  >
    <div class="d-flex justify-space-between align-center">
      <span class="text-body-1">
        {{ message }}
      </span>
      <v-btn class="justify-end" text color="white" @click="onClose">{{
        $t('close')
      }}</v-btn>
    </div>
  </v-snackbar>
</template>

<script>
  export default {
    name: 'NSnackbar',
    props: {
      value: Boolean,
      message: String,
      color: String,
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
    methods: {
      onClose() {
        this.innerValue = false
      },
    },
  }
</script>

<style></style>
