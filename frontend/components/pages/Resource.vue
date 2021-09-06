<template>
  <page-layout>
    <v-tabs v-model="tabs">
      <v-tab class="tfn-important" href="#settings">
        {{ $t('resource.tabs.settings') }}
      </v-tab>
      <v-tab class="tfn-important" href="#permissions">
        {{ $t('resource.tabs.permissions') }}
      </v-tab>
    </v-tabs>
    <v-tabs-items v-model="tabs">
      <v-tab-item value="settings">
        <v-container v-if="tabs === 'settings'">
          <resource-tab-settings
            :id="resources.id"
            :description.sync="resources.description"
            @submit="onSubmit"
          ></resource-tab-settings>
        </v-container>
      </v-tab-item>
      <v-tab-item value="permissions">
        <v-container v-if="tabs === 'permissions'">
          <resource-tab-permissions
            ref="permissions"
            :id="id"
          ></resource-tab-permissions>
        </v-container>
      </v-tab-item>
    </v-tabs-items>
  </page-layout>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'Resource',
    mixins: [axiosMixin],
    data() {
      return {
        tabs: 'settings',
        id: '',
        resources: {
          id: '',
          description: '',
          permissions: [],
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
      getResourceUrl() {
        return `${this.$urls.api.v1.resources}/${this.id}`
      },
      getData() {
        const url = this.getResourceUrl()
        this.get(url)
          .then((result) => {
            if (result.status == 200) {
              const data = result.data
              if (data) {
                this.resources = {
                  id: data.id || '',
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
        const url = this.getResourceUrl()
        this.put(url, {
          description: this.resources.description,
        }).catch((err) => {
          console.log(err)
        })
      },
    },
  }
</script>

<style></style>
