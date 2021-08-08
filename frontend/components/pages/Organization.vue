<template>
  <page-layout>
    <v-tabs v-model="tabs">
      <v-tab class="tfn-important" href="#settings">
        {{ $t('organization.tabs.settings') }}
      </v-tab>
      <v-tab class="tfn-important" href="#users">
        {{ $t('organization.tabs.users') }}
      </v-tab>
    </v-tabs>
    <v-tabs-items v-model="tabs">
      <v-tab-item value="settings">
        <v-container v-if="tabs === 'settings'">
          <organization-tab-settings
            :name.sync="organization.name"
            :description.sync="organization.description"
            @submit="onSubmit"
          ></organization-tab-settings>
        </v-container>
      </v-tab-item>
      <v-tab-item value="users">
        <v-container v-if="tabs === 'users'">
          <organization-tab-users :id="id"></organization-tab-users>
        </v-container>
      </v-tab-item>
    </v-tabs-items>
  </page-layout>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'Organization',
    mixins: [axiosMixin],
    data() {
      return {
        tabs: 'settings',
        id: '',
        organization: {
          name: '',
          description: '',
          permissions: [],
        },
      }
    },
    created() {
      this.id = this.$route.params.id
      this.getData()
    },
    methods: {
      getOrganizationUrl() {
        return `${this.$urls.api.v1.organizations}/${this.id}`
      },
      getData() {
        const url = this.getOrganizationUrl()
        this.get(url)
          .then((result) => {
            if (result.status == 200) {
              const data = result.data
              if (data) {
                this.organization = {
                  name: data.name || '',
                  description: data.description || '',
                  permissions: data.permissions || [],
                }
              }
            }
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onSubmit() {
        const url = this.getOrganizationUrl()
        this.put(url, {
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
