<template>
  <v-container>
    <v-row>
      <div :class="nameClass">
        <slot name="before-name"></slot>
        <required-text
          name="name"
          :label="$t('inputs.Name')"
          id="name"
          v-model="innerName"
          v-bind="innerNameAttrs"
        ></required-text>
        <slot name="after-name"></slot>
      </div>
      <div :class="descriptionClass">
        <slot name="before-description"></slot>
        <required-text
          name="description"
          :label="$t('inputs.Description')"
          id="description"
          v-model="innerDescription"
          v-bind="innerDescriptionAttrs"
        ></required-text>
        <slot name="after-description"></slot>
      </div>
      <slot></slot>
    </v-row>
  </v-container>
</template>

<script>
  export default {
    name: 'nameAndDescription',
    props: {
      name: String,
      description: String,
      nameClass: {
        type: String,
        default: 'col-12 col-md-4',
      },
      descriptionClass: {
        type: String,
        default: 'col-12 col-md-4',
      },
      nameAttrs: {
        type: Object,
        default: () => ({}),
      },
      descriptionAttrs: {
        type: Object,
        default: () => ({}),
      },
    },
    computed: {
      innerName: {
        get() {
          return this.name
        },
        set(val) {
          this.$emit('update:name', val)
        },
      },
      innerDescription: {
        get() {
          return this.description
        },
        set(val) {
          this.$emit('update:description', val)
        },
      },
      innerNameAttrs() {
        return Object.assign(
          {
            counter: 255,
          },
          this.nameAttrs
        )
      },
      innerDescriptionAttrs() {
        return Object.assign(
          {
            counter: 255,
          },
          this.descriptionAttrs
        )
      },
    },
  }
</script>

<style></style>
