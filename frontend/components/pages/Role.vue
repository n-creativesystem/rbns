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
            :organization-id="organizationId"
            @submit="onSubmit"
          ></role-tab-settings>
        </v-container>
      </v-tab-item>
      <v-tab-item value="permissions">
        <v-container v-if="tabs === 'permissions'">
          <role-tab-permissions
            ref="permissions"
            :id="id"
            :organization-id="organizationId"
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
  export default {
    name: 'Role',
    data() {
      return {
        organizationId: '',
        tabs: '',
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
      this.tabs = 'settings'
      this.organizationId = this.$route.params.orgId
      this.id = this.$route.params.id
      this.$store
        .dispatch('roles/findById', {
          organizationId: this.organizationId,
          roleId: this.id,
        })
        .then((result) => {
          this.role = result
        })
        .catch((err) => {
          console.log(err)
        })
    },
    methods: {
      onSubmit() {
        this.$store
          .dispatch('roles/put', {
            organizationId: this.organizationId,
            roleId: this.id,
            data: {
              name: this.role.name,
              description: this.role.description,
            },
          })
          .catch((err) => {
            console.log(err)
          })
      },
    },
  }
</script>

<style></style>
