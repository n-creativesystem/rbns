<template>
  <name-and-description-form
    :name.sync="id"
    name-class="col-12 col-md-8"
    :name-attrs="{ disabled: true }"
    :description.sync="innerDescription"
    description-class="col-12 col-md-8"
    :btn-text="$t('save')"
    submit-class="col-12 col-md-8"
    @submit="onSubmit"
  >
  </name-and-description-form>
</template>

<script>
  export default {
    name: 'ResourceTabSettings',
    data() {
      return {}
    },
    props: {
      id: String,
      description: String,
    },
    computed: {
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
    },
  }
</script>

<style></style>
