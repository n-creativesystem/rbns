<template>
  <v-container>
    <n-alert-information
      :value="true"
      :message="$t('organization.info.user.message')"
      :btn-text="$t('organization.info.user.text')"
      @click="dialog = true"
    ></n-alert-information>
    <dialog-users
      v-model="dialog"
      :user-id.sync="userId"
      :selection.sync="selection"
      @submit="onSubmit"
    ></dialog-users>
    <n-data-table
      :headers="headers"
      :items="items"
      actions
      delete-action
      @delete="onDelete"
    >
      <template v-slot:[`item.key`]="{ item }">
        <router-link :to="`/organizations/${id}/users/${item.key}`">
          {{ item.key }}
        </router-link>
      </template>
    </n-data-table>
  </v-container>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'OrganizationTabUsers',
    mixins: [axiosMixin],
    data() {
      return {
        items: [],
        dialog: false,
        userId: '',
        selection: [],
      }
    },
    props: {
      id: String,
    },
    created() {
      this.getData()
    },
    methods: {
      getData() {
        this.get(`${this.$urls.api.v1.organizations}/${this.id}`)
          .then((res) => {
            this.items = res.data.users
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onSubmit(e) {
        const data = {
          key: e.userId,
          organization_id: this.id,
          roles: e.selection.map((s) => {
            return { id: s.id }
          }),
        }
        this.post(`${this.$urls.api.v1.organizations}/${this.id}/users`, data)
          .then(() => {
            this.getData()
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onDialogClick() {},
      onDelete(item) {
        console.log(item)
        this.delete(
          `${this.$urls.api.v1.organizations}/${this.id}/users/${item.key}`
        )
          .then(() => {
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
            text: this.$t('user.entity.userId'),
            name: 'userId',
            align: 'start',
            value: 'key',
          },
        ]
      },
    },
  }
</script>

<style></style>
