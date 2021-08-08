<template>
  <page-layout>
    <input-and-list
      :name.sync="name"
      :description.sync="description"
      :add-meta="addMeta"
      :list-meta="listMeta"
      :headers="headers"
      :items="items"
      :edit-action="false"
      @submit="onSubmit"
      @delete="onDelete"
      ref="form"
    >
    </input-and-list>
  </page-layout>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'Permission',
    mixins: [axiosMixin],
    data() {
      return {
        name: '',
        description: '',
        items: [],
      }
    },
    created() {
      this.getData()
    },
    methods: {
      getData() {
        this.items = []
        this.get(this.$urls.api.v1.permissions)
          .then((result) => {
            if (result.status == 200) {
              if (result.data.permissions) {
                this.items = result.data.permissions
              }
            }
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onSubmit() {
        this.post(this.$urls.api.v1.permissions, {
          permissions: [
            {
              name: this.name,
              description: this.description,
            },
          ],
        })
          .then(() => {
            this.getData()
            this.$refs.form.reset()
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onDelete(e) {
        const id = e.id
        this.delete(`${this.$urls.api.v1.permissions}/${id}`)
          .then((result) => {
            console.log(result)
            this.getData()
          })
          .catch((err) => {
            console.log(err)
          })
      },
    },
    computed: {
      headers() {
        return [
          {
            text: this.$t('permission.entity.name'),
            align: 'start',
            value: 'name',
          },
          {
            text: this.$t('permission.entity.description'),
            align: 'start',
            value: 'description',
          },
        ]
      },
      addMeta() {
        return {
          title: this.$t('permission.form.add.title'),
          caption: this.$t('permission.form.add.caption'),
        }
      },
      listMeta() {
        return {
          title: this.$t('permission.form.list.title'),
          caption: this.$t('permission.form.list.caption'),
        }
      },
    },
  }
</script>

<style></style>
