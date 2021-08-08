<template>
  <n-text-field
    required
    outlined
    :rules="innerRules"
    :loading="loading"
    v-bind="$attrs"
    v-model="innverValue"
    v-on="$listeners"
  >
    <template v-slot:label>
      {{ $attrs.label }} <span class="red--text">*</span>
    </template>
    <template
      v-for="(slotContent, slotName) of $scopedSlots"
      #[slotName]="data"
    >
      <slot :name="slotName" v-bind="data"></slot>
    </template>
  </n-text-field>
</template>

<script>
  import loadingMixin from '@mixin/loading'

  export default {
    name: 'requiredText',
    mixins: [loadingMixin],
    props: {
      value: {
        required: true,
        default: undefined,
      },
      rules: {
        type: Array,
        default: () => [],
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
      innerRules() {
        if (this.rules) {
          return [this.required].concat(this.rules.filter((r) => r))
        } else {
          return [this.required]
        }
      },
    },
    methods: {
      required(value) {
        return !!value || this.$t('required fields')
      },
    },
  }
</script>

<style></style>
