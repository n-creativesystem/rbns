<template>
  <v-container>
    <n-alert-information
      :value="true"
      :message="$t('user.info.role.message')"
      :btn-text="$t('user.info.role.text')"
      @click="dialog = true"
    ></n-alert-information>
    <dialog-roles
      v-model="dialog"
      @click="onDialogClick"
      :ommit-ids="ommitIds"
      :selection.sync="selection"
    ></dialog-roles>
    <n-data-table
      :actions="true"
      :delete-action="true"
      @delete="onDelete"
      :headers="headers"
      :items="user.roles"
    ></n-data-table>
  </v-container>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'UserTabRoles',
    mixins: [axiosMixin],
    data() {
      return {
        user: {
          key: '',
          organizationId: '',
          roles: [],
          permissions: [],
        },
        dialog: false,
        ids: [],
        selection: [],
      }
    },
    props: {
      id: String,
      userKey: String,
    },
    methods: {
      getUrl() {
        return this.$urls.api.v1.users.format(this.id, this.userKey)
      },
      getData() {
        this.get(this.getUrl())
          .then((result) => {
            if (result.status == 200) {
              const data = result.data
              if (data) {
                this.user = data
              }
            }
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onDialogClick(items) {
        this.put(`${this.getUrl()}/roles`, {
          roles: items.map((item) => {
            return {
              id: item.id,
            }
          }),
        })
          .then(() => {
            this.getData()
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onDelete(item) {
        const url = `${this.getUrl()}/roles/${item.id}`
        this.delete(url)
          .then(() => {
            this.getData()
          })
          .catch((err) => {
            console.log(err)
          })
      },
      created() {
        this.getData()
      },
    },
    computed: {
      ommitIds() {
        return this.user.roles.map((role) => role.id)
      },
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
    created() {
      this.created()
    },
  }
</script>

<style></style>
