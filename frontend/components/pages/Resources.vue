<template>
  <page-layout>
    <list-form
      :title="$t('resource.form.title')"
      :caption="$t('resource.form.caption')"
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
          >{{ $t('resource.form.add') }}</btn-tfn
        >
      </template>
      <template v-slot:[`item.id`]="{ item }">
        <router-link :to="`/resources/${encodeURI(item.id)}`">
          {{ item.id }}
        </router-link>
      </template>
    </list-form>
    <dialog-name-and-desc
      v-if="createDialog"
      v-model="createDialog"
      :title="$t('resource.dialog.create.title')"
      :name="name"
      :description="description"
      @submit="onCreateClick"
    ></dialog-name-and-desc>
  </page-layout>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'Resources',
    mixins: [axiosMixin],
    data() {
      return {
        items: [],
        createDialog: false,
        name: '',
        description: '',
      }
    },
    created() {
      this.getData()
    },
    methods: {
      getData() {
        this.get(this.$urls.api.v1.resources)
          .then((result) => {
            if (result.status == 200) {
              if (result.data.resources) {
                this.items = result.data.resources
              }
            }
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onCreateClick(e) {
        this.createDialog = false
        this.post(this.$urls.api.v1.resources, {
          roles: [
            {
              name: e.name,
              description: e.description,
            },
          ],
        })
          .then(() => {
            return this.getData()
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onDelete(item) {
        this.delete(`${this.$urls.api.v1.resource}/${item.id}`)
          .then(() => {
            return this.getData()
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
            text: this.$t('resource.entity.id'),
            align: 'start',
            value: 'id',
          },
          {
            text: this.$t('resource.entity.description'),
            align: 'start',
            value: 'description',
          },
        ]
      },
    },
  }
</script>

<style></style>
