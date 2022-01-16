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
        :omit-ids="innerIds"
        :selection.sync="selection"
      ></dialog-permissions>
      <n-data-table
        :headers="headers"
        :items="items"
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
      }
    },
    props: {
      id: String,
      organizationId: String,
    },
    methods: {
      getData() {
        this.$store
          .dispatch('roles/findPermissions', {
            organizationId: this.organizationId,
            roleId: this.id,
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onDialogClick(items) {
        const permissions = items.map((select) => {
          return {
            id: select.id,
          }
        })
        this.$store
          .dispatch('roles/putRolePermissions', {
            organizationId: this.organizationId,
            roleId: this.id,
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
        this.$store
          .dispatch('roles/removeRolePermissions', {
            organizationId: this.organizationId,
            roleId: this.id,
            permissionId: item.id,
          })
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
      items() {
        return this.$store.getters['roles/listPermissions']
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
