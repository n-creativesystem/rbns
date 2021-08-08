<template>
  <header-content-parts>
    <template v-slot:header>
      <title-and-caption :title="`${title}`" :caption="caption">
        <template v-slot:title>
          <slot name="title"></slot>
        </template>
        <template v-slot:caption>
          <slot name="caption"></slot>
        </template>
      </title-and-caption>
    </template>
    <template v-slot:content>
      <name-and-description-form
        ref="form"
        :name.sync="innerName"
        :description.sync="innerDescription"
        @submit="onSubmit"
      ></name-and-description-form>
    </template>
  </header-content-parts>
</template>

<script>
  export default {
    name: 'inputForm',
    props: {
      title: String,
      caption: String,
      name: String,
      description: String,
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
    methods: {
      onSubmit(e) {
        this.$emit('submit', e)
      },
      reset() {
        this.$refs.form.reset()
      },
    },
  }
</script>

<style></style>
