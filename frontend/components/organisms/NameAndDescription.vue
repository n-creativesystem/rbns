<template>
  <v-container>
    <v-row>
      <div :class="nameClass">
        <slot name="before-name"></slot>
        <required-text
          name="name"
          :label="$t('inputs.Name')"
          id="name"
          :counter="255"
          v-model="innerName"
        ></required-text>
        <slot name="after-name"></slot>
      </div>
      <div :class="descriptionClass">
        <slot name="before-description"></slot>
        <required-text
          name="description"
          :label="$t('inputs.Description')"
          id="description"
          :counter="255"
          v-model="innerDescription"
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
      nameClass: {
        type: String,
        default: 'col-12 col-md-4',
      },
      description: String,
      descriptionClass: {
        type: String,
        default: 'col-12 col-md-4',
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
    },
  }
</script>

<style></style>
