<template>
  <v-container v-bind="containerAttrs">
    <v-row>
      <div
        v-for="input in inputs"
        :key="input.name"
        :class="input.class || 'col-12 col-md-4'"
      >
        <slot :name="`before-${input.name}`" v-bind="input"></slot>
        <required-text
          :name="input.name"
          :label="input.label"
          :id="input.id || input.name"
          v-bind="input.attrs"
          v-model="input.value"
          @input="onInput"
        ></required-text>
        <slot :name="`after-${input.name}`" v-bind="input"></slot>
      </div>
    </v-row>
  </v-container>
</template>

<script>
  export default {
    name: 'inputsContainer',
    data() {
      return {
        inputs: this.value,
      }
    },
    props: {
      containerAttrs: Object,
      value: {
        type: Array,
        required: true,
      },
      btnText: {
        type: String,
        default: '保存',
      },
      submitClass: {
        type: String,
        default: 'col-12 col-md-4',
      },
    },
    methods: {
      submit(e) {
        this.$emit('submit', {
          event: e,
          inputs: this.inputs,
        })
      },
      onInput() {
        this.$emit('input', this.inputs)
      },
    },
  }
</script>

<style></style>
