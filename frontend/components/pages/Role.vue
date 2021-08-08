<template>
  <page-layout>
    <v-tabs v-model="tabs">
      <v-tab class="tfn-important" href="#settings">
        {{ $t('role.tabs.settings') }}
      </v-tab>
      <v-tab class="tfn-important" href="#permissions">
        {{ $t('role.tabs.permissions') }}
      </v-tab>
      <v-tab class="tfn-important" href="#users">
        {{ $t('role.tabs.users') }}
      </v-tab>
    </v-tabs>
    <v-tabs-items v-model="tabs">
      <v-tab-item value="settings">
        <v-container v-if="tabs === 'settings'">
          <role-tab-settings
            :name.sync="role.name"
            :description.sync="role.description"
            @submit="onSubmit"
          ></role-tab-settings>
        </v-container>
      </v-tab-item>
      <v-tab-item value="permissions">
        <v-container v-if="tabs === 'permissions'">
          <role-tab-permissions
            ref="permissions"
            :id="id"
          ></role-tab-permissions>
        </v-container>
      </v-tab-item>
      <v-tab-item value="users">
        <v-container v-if="tabs === 'users'">
          <role-tab-users :items="role.organizationUsers"></role-tab-users>
        </v-container>
      </v-tab-item>
    </v-tabs-items>
  </page-layout>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'Role',
    mixins: [axiosMixin],
    data() {
      return {
        tabs: 'settings',
        id: '',
        role: {
          name: '',
          description: '',
          permissions: [],
          organizationUsers: [],
        },
        snackbar: false,
        message: '',
      }
    },
    created() {
      this.id = this.$route.params.id
      this.getData()
    },
    methods: {
      getRoleUrl() {
        return `${this.$urls.api.v1.roles}/${this.id}`
      },
      getData() {
        const url = this.getRoleUrl()
        this.get(url)
          .then((result) => {
            if (result.status == 200) {
              const data = result.data
              if (data) {
                this.role = {
                  name: data.name || '',
                  description: data.description || '',
                  permissions: data.permissions || [],
                  organizationUsers: data.organizationUsers || [],
                }
              }
            }
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onSubmit() {
        const url = this.getRoleUrl()
        this.put(url, {
          name: this.role.name,
          description: this.role.description,
        }).catch((err) => {
          console.log(err)
        })
      },
    },
  }
</script>

<style></style>
