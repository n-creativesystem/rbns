<template>
  <page-layout>
    <v-tabs v-model="tabs">
      <v-tab class="tfn-important" href="#role">
        {{ $t('user.tabs.roles') }}
      </v-tab>
    </v-tabs>
    <v-tabs-items v-model="tabs">
      <v-tab-item value="role">
        <user-tab-roles
          v-if="tabs === 'role'"
          ref="role"
          :id="id"
          :user-key="userKey"
        ></user-tab-roles>
      </v-tab-item>
    </v-tabs-items>
  </page-layout>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'User',
    mixins: [axiosMixin],
    data() {
      return {
        tabs: 'role',
        id: '',
        userKey: '',
        user: {
          key: '',
          organizationId: '',
          roles: [],
          permissions: [],
        },
      }
    },
    created() {
      this.id = this.$route.params.id
      this.userKey = this.$route.params.userKey
    },
    methods: {
      onSubmit() {
        this.put(this.getUrl(), {
          name: this.organization.name,
          description: this.organization.description,
        })
          .then((result) => {
            console.log(result)
          })
          .catch((err) => {
            console.log(err)
          })
      },
    },
  }
</script>

<style></style>
