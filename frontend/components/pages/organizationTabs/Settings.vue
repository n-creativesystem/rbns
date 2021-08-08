<template>
  <name-and-description-form
    :name.sync="innerName"
    name-class="col-12 col-md-8"
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
    name: 'OrganizationTabSettings',
    data() {
      return {}
    },
    props: {
      id: String,
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
    },
  }
</script>

<style></style>
