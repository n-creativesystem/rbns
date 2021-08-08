<template>
  <v-data-table
    :headers="innerHeaders"
    :items="loading ? [] : items"
    :loading="loading"
    :search="search"
    v-bind="$attrs"
  >
    <template v-slot:top>
      <div class="d-flex flex-row-reverse">
        <div class="pa-2">
          <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            :label="$t('search')"
            single-line
            hide-details
          ></v-text-field>
        </div>
      </div>
    </template>
    <template v-slot:no-data>
      {{ $t('no data') }}
    </template>
    <template v-slot:no-results>
      {{ $t('no record') }}
    </template>
    <template v-slot:loading>
      Loading
      <v-progress-circular indeterminate color="purple"></v-progress-circular>
    </template>
    <template v-if="nemeToChip" v-slot:[`item.name`]="{ item }">
      <v-chip label>{{ item.name }}</v-chip>
    </template>
    <template v-slot:[`item.actions`]="{ item }">
      <v-icon v-if="editAction" class="mr-2" small @click="onEdit(item)"
        >mdi-pencil</v-icon
      >
      <v-icon v-if="deleteAction" small @click="onDelete(item)"
        >mdi-delete</v-icon
      >
    </template>
    <template
      v-for="(slotContent, slotName) of $scopedSlots"
      #[slotName]="data"
    >
      <slot :name="slotName" v-bind="data"></slot>
    </template>
  </v-data-table>
</template>

<script>
  import loadingMixin from '@mixin/loading'
  export default {
    name: 'nDataTable',
    mixins: [loadingMixin],
    data() {
      return {
        search: null,
      }
    },
    props: {
      items: {
        type: Array,
        required: true,
        default: () => [],
      },
      headers: {
        type: Array,
        required: true,
        default: () => [],
      },
      nemeToChip: {
        type: Boolean,
        default: true,
      },
      actions: Boolean,
      editAction: Boolean,
      deleteAction: Boolean,
    },
    computed: {
      innerHeaders() {
        return this.actions
          ? this.headers.concat([
              {
                text: '',
                value: 'actions',
              },
            ])
          : this.headers
      },
    },
    methods: {
      onEdit(item) {
        this.$emit('edit', item)
      },
      onDelete(item) {
        this.$emit('delete', item)
      },
    },
  }
</script>

<style></style>
