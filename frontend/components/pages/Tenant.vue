<template>
  <form-parts class="input-and-list">
    <input-form
      :title="title"
      :caption="caption"
      :name="name"
      @update:name="toLower"
      :description.sync="description"
      :name-attrs="nameAttrs"
      :description-attrs="descriptionAttrs"
      @submit="onSubmit"
    ></input-form>
  </form-parts>
</template>
<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'Tenant',
    mixins: [axiosMixin],
    data: () => ({
      name: '',
      description: '',
    }),
    methods: {
      onSubmit() {
        this.post(this.$urls.api.v1.tenants, {
          name: this.name,
          description: this.description,
        })
          .then((result) => {
            console.log(result)
            location.href = '/permissions'
          })
          .catch((err) => {
            console.log(err)
          })
      },
      possibleCharactersValidate(value) {
        return (
          /^[a-zA-Z0-9_+-]+(.[a-zA-Z0-9_+-]+)*$/.test(value) ||
          this.$t('tenant.possibleCharactersValidate')
        )
      },
      toLower(value) {
        this.name = value.toLowerCase()
      },
    },
    computed: {
      title() {
        return this.$t('tenant.title')
      },
      caption() {
        return this.$t('tenant.caption')
      },
      tenantRules() {
        return [this.possibleCharactersValidate]
      },
      nameAttrs() {
        return {
          rules: this.tenantRules,
          counter: 30,
        }
      },
      descriptionAttrs() {
        return {}
      },
    },
  }
</script>
<style lang=""></style>
