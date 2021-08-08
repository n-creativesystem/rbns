<template>
  <v-dialog
    v-model="innerValue"
    :fullscreen="$vuetify.breakpoint.mobile"
    scrollable
    max-width="500px"
    v-bind="$attrs"
  >
    <v-card>
      <v-system-bar light window>
        <v-spacer></v-spacer>
        <v-btn icon @click="onClose">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-system-bar>
      <v-card-title v-if="title" :class="innerTitleClass">
        {{ title }}
      </v-card-title>
      <slot></slot>
    </v-card>
  </v-dialog>
</template>

<script>
  export default {
    name: 'nDialog',
    data() {
      return {}
    },
    props: {
      value: Boolean,
      activator: String,
      title: String,
      titleClass: String,
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
      innerTitleClass() {
        return this.titleClass || 'text-h5 grey lighten-2'
      },
    },
    mounted() {
      this.$nextTick(() => {
        if (this.innerValue) {
          this.$emit('open')
        }
      })
    },
    methods: {
      onClose() {
        this.innerValue = false
      },
    },
    watch: {
      innerValue: {
        handler(val) {
          if (!val) {
            this.$emit('close')
          } else {
            this.$emit('open')
          }
        },
      },
    },
  }
</script>

<style></style>
