<template>
  <page-layout>
    <list-form
      :title="$t('organization.form.title')"
      :caption="$t('organization.form.caption')"
      :headers="headers"
      :items="items"
      :edit-action="false"
      @delete="onDelete"
    >
      <template v-slot:title>
        <btn-tfn
          class="justify-end"
          color="info"
          @click.stop="createDialog = true"
          >{{ $t('organization.form.add') }}</btn-tfn
        >
      </template>
      <template v-slot:[`item.name`]="{ item }">
        <router-link :to="`/organizations/${item.id}`">
          {{ item.name }}
        </router-link>
      </template>
    </list-form>
    <dialog-name-and-desc
      v-if="createDialog"
      v-model="createDialog"
      :title="$t('organization.dialog.create.title')"
      :name="name"
      :description="description"
      @submit="onCreateClick"
    ></dialog-name-and-desc>
  </page-layout>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'Organizations',
    mixins: [axiosMixin],
    data() {
      return {
        name: '',
        description: '',
        items: [],
        createDialog: false,
      }
    },
    methods: {
      getData() {
        this.items = []
        this.get(this.$urls.api.v1.organizations)
          .then((result) => {
            if (result.status == 200) {
              if (result.data.organizations) {
                this.items = result.data.organizations
              }
            }
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onCreateClick(e) {
        this.createDialog = false
        this.post(this.$urls.api.v1.organizations, {
          name: e.name,
          description: e.description,
        })
          .then(() => {
            this.getData()
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onDelete(e) {
        const id = e.id
        this.delete(`${this.$urls.api.v1.organizations}/${id}`)
          .then(() => {
            this.getData()
          })
          .catch((err) => {
            console.log(err)
          })
      },
    },
    created() {
      this.getData()
    },
    computed: {
      headers() {
        return [
          {
            text: this.$t('organization.entity.name'),
            align: 'start',
            value: 'name',
          },
          {
            text: this.$t('organization.entity.description'),
            align: 'start',
            value: 'description',
          },
        ]
      },
    },
  }
</script>

<style></style>
