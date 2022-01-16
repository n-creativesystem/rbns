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
      <template v-slot:content-header>
        <div class="pt-10">
          <v-row>
            <v-col md="6" cols="12">
              <v-select
                :items="organizations"
                item-text="name"
                item-value="id"
                v-model="organizationValue"
                :label="$t('organization.name')"
                @change="onChange"
              ></v-select>
            </v-col>
          </v-row>
        </div>
      </template>
      <template v-slot:[`item.name`]="{ item }">
        <router-link :to="`/roles/${organizationValue}/${item.id}`">
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
  import loadingMixin from '@mixin/loading'
  export default {
    name: 'Roles',
    mixins: [loadingMixin],
    data() {
      return {
        createDialog: false,
        name: '',
        description: '',
        organizationValue: undefined,
      }
    },
    created() {
      const lodingCounter = this.onloading()
      this.$store.commit('roles/setRoles', [])
      this.$store
        .dispatch('organizations/findAll')
        .then(() => {
          if (this.organizations.length == 1) {
            this.organizationValue = this.organizations[0].id
            this.onChange()
          }
        })
        .finally(() => this.unloading(lodingCounter))
    },
    methods: {
      onCreateClick(e) {
        this.createDialog = false
        const counter = this.onloading()
        this.$store
          .dispatch('roles/add', {
            organizationId: this.organizationValue,
            roles: [
              {
                name: e.name,
                description: e.description,
              },
            ],
          })
          .then(() => {
            this.onChange()
          })
          .catch((err) => {
            console.log(err)
          })
          .finally(() => this.unloading(counter))
      },
      onDelete(item) {
        const counter = this.onloading()
        this.$store
          .dispatch('roles/remove', {
            organizationId: this.organizationValue,
            roleId: item.id,
          })
          .then(() => {
            this.onChange()
          })
          .catch((err) => {
            console.log(err)
          })
          .finally(() => this.unloading(counter))
      },
      onChange() {
        this.$store.dispatch('roles/findAll', {
          organizationId: this.organizationValue,
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
      organizations() {
        return this.$store.getters['organizations/list']
      },
      items() {
        return this.$store.getters['roles/list']
      },
    },
  }
</script>

<style></style>
