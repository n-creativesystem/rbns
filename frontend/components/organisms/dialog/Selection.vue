<template>
  <v-card>
    <v-card-text>
      <v-card outlined tile class="pa-4">
        <n-selection v-model="innerSelection" :items="items"></n-selection>
      </v-card>
    </v-card-text>
    <v-card-actions class="justify-end">
      <btn-tfn @click="onClick">{{ text }}</btn-tfn>
    </v-card-actions>
  </v-card>
</template>

<script>
  import loadingMixin from '@mixin/loading'
  export default {
    name: 'DialogSelection',
    mixins: [loadingMixin],
    props: {
      text: String,
      items: Array,
      selection: Array,
    },
    computed: {
      innerItems: {
        get() {
          return this.items
        },
        set(val) {
          this.$emit('update:items', val)
        },
      },
      innerSelection: {
        get() {
          return this.selection
        },
        set(val) {
          this.$emit('update:selection', val)
        },
      },
    },
    methods: {
      onClick() {
        const items = this.innerItems.filter((_, idx) =>
          this.innerSelection.some((index) => idx === index)
        )
        this.$emit('click', items)
      },
    },
  }
</script>

<style></style>
