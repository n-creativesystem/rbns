<template>
  <v-chip-group
    v-model="innerValue"
    :column="column"
    :multiple="multiple"
    :active-class="activeClass"
    v-bind="$attrs"
  >
    <v-progress-circular
      v-if="loading"
      indeterminate
      color="purple"
      class="justify-end"
    ></v-progress-circular>
    <v-chip v-else filter v-for="item in items" :key="item.id">{{
      item.name
    }}</v-chip>
  </v-chip-group>
</template>

<script>
  import loadingMixin from '@mixin/loading'
  export default {
    name: 'NSelection',
    mixins: [loadingMixin],
    props: {
      value: Array,
      column: {
        type: Boolean,
        default: true,
      },
      multiple: {
        type: Boolean,
        default: true,
      },
      activeClass: {
        type: String,
        default: 'deep-purple accent-4 white--text',
      },
      items: Array,
    },
    computed: {
      innerValue: {
        get() {
          return this.value
        },
        set(value) {
          this.$emit('input', value)
        },
      },
    },
  }
</script>

<style></style>
