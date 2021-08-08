<template>
  <div class="row">
    <div class="col-12">
      <n-alert-information
        :value="true"
        :message="$t('role.info.permission.message')"
        :btn-text="$t('role.info.permission.text')"
        @click="dialog = true"
      ></n-alert-information>
      <dialog-permissions
        v-model="dialog"
        @click="onDialogClick"
        :ommit-ids="innerIds"
        :selection.sync="selection"
      ></dialog-permissions>
      <n-data-table
        :headers="headers"
        :items="loading ? [] : items"
        actions
        delete-action
        @delete="onDelete"
      >
      </n-data-table>
    </div>
  </div>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'RoleTabPermissions',
    mixins: [axiosMixin],
    data() {
      return {
        dialog: false,
        selection: [],
        items: [],
      }
    },
    props: {
      id: String,
    },
    methods: {
      getData() {
        this.get(`${this.$urls.api.v1.roles}/${this.id}/permissions`)
          .then((result) => {
            this.items = result.data.permissions
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onDialogClick(items) {
        const permissions = items.map((select) => {
          return {
            Id: select.id,
          }
        })
        this.put(`${this.$urls.api.v1.roles}/${this.id}/permissions`, {
          permissions: permissions,
        })
          .then(() => {
            this.getData()
          })
          .catch((err) => {
            this.$emit('save-error', err)
          })
      },
      onDelete(item) {
        this.delete(
          `${this.$urls.api.v1.roles}/${this.id}/permissions/${item.id}`
        )
          .then(() => {
            this.getData()
          })
          .catch((err) => {
            this.$emit('save-error', err)
          })
      },
      created() {
        this.getData()
      },
    },
    computed: {
      innerIds() {
        return this.items && this.items.length
          ? this.items.map((i) => i.id)
          : []
      },
      headers() {
        return [
          {
            text: this.$t('permission.entity.name'),
            value: 'name',
            align: 'start',
          },
          {
            text: this.$t('permission.entity.description'),
            value: 'description',
            align: 'start',
          },
        ]
      },
    },
    watch: {
      dialog: {
        handler(val) {
          if (val) {
            this.selection = []
          }
        },
      },
    },
    created() {
      this.created()
    },
  }
</script>

<style></style>
