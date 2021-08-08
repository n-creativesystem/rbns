<template>
  <v-text-field
    clearable
    outlined
    v-bind="$attrs"
    v-model="innverValue"
    v-on="$listeners"
  >
    <template
      v-for="(slotContent, slotName) of $scopedSlots"
      #[slotName]="data"
    >
      <slot :name="slotName" v-bind="data"></slot>
    </template>
  </v-text-field>
</template>

<script>
  export default {
    name: 'NTextField',
    props: {
      value: {
        required: true,
        default: undefined,
      },
    },
    computed: {
      innverValue: {
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
