<template>
  <n-form @submit="submit" ref="form">
    <name-and-description
      :name.sync="innerName"
      :name-class="nameClass"
      :description.sync="innerDescription"
      :description-class="descriptionClass"
    >
      <div :class="submitClass">
        <slot name="before-submit"></slot>
        <btn-tfn type="submit">
          <v-icon>mdi-add</v-icon>
          {{ innerBtnText }}
        </btn-tfn>
        <slot name="after-submit"></slot>
      </div>
    </name-and-description>
  </n-form>
</template>

<script>
  export default {
    name: 'nameAndDescriptionForm',
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
      btnText: String,
      submitClass: {
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
      innerBtnText() {
        return this.btnText || this.$t('inputs.Add')
      },
    },
    methods: {
      submit(e) {
        this.$emit('submit', e)
      },
      reset() {
        this.$refs.form.reset()
      },
    },
  }
</script>

<style></style>
