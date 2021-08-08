<template>
  <page-layout>
    <list-form
      :title="$t('role.form.title')"
      :caption="$t('role.form.caption')"
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
          >{{ $t('role.form.add') }}</btn-tfn
        >
      </template>
      <template v-slot:[`item.name`]="{ item }">
        <router-link :to="`/roles/${item.id}`">
          {{ item.name }}
        </router-link>
      </template>
    </list-form>
    <dialog-name-and-desc
      v-if="createDialog"
      v-model="createDialog"
      :title="$t('role.dialog.create.title')"
      :name="name"
      :description="description"
      @submit="onCreateClick"
    ></dialog-name-and-desc>
  </page-layout>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'Roles',
    mixins: [axiosMixin],
    data() {
      return {
        items: [],
        createDialog: false,
        name: '',
        description: '',
      }
    },
    created() {
      this.getData()
    },
    methods: {
      getData() {
        this.get(this.$urls.api.v1.roles)
          .then((result) => {
            if (result.status == 200) {
              if (result.data.roles) {
                this.items = result.data.roles
              }
            }
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onCreateClick(e) {
        this.createDialog = false
        this.post(this.$urls.api.v1.roles, {
          roles: [
            {
              name: e.name,
              description: e.description,
            },
          ],
        })
          .then(() => {
            return this.getData()
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onDelete(item) {
        this.delete(`${this.$urls.api.v1.roles}/${item.id}`)
          .then(() => {
            return this.getData()
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
            text: this.$t('role.entity.name'),
            align: 'start',
            value: 'name',
          },
          {
            text: this.$t('role.entity.description'),
            align: 'start',
            value: 'description',
          },
        ]
      },
    },
  }
</script>

<style></style>
