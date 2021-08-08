<template>
  <header-content-parts>
    <template v-slot:header>
      <title-and-caption :title="`${title}`" :caption="caption">
        <template v-slot:title>
          <slot name="title"></slot>
        </template>
        <template v-slot:caption>
          <slot name="caption"></slot>
        </template>
      </title-and-caption>
    </template>
    <template v-slot:content>
      <n-data-table
        class="justify-center"
        :headers="headers"
        :items="items"
        :loading="loading"
        actions
        :edit-action="editAction"
        :delete-action="deleteAction"
        @edit="onEdit"
        @delete="onDelete"
      >
        <template
          v-for="(slotContent, slotName) of $scopedSlots"
          #[slotName]="data"
        >
          <slot
            v-if="slotName.startsWith('item.')"
            :name="slotName"
            v-bind="data"
          ></slot>
        </template>
      </n-data-table>
    </template>
  </header-content-parts>
</template>

<script>
  import loadingMixin from '@mixin/loading'
  export default {
    name: 'listForm',
    mixins: [loadingMixin],
    data() {
      return {}
    },
    props: {
      title: String,
      caption: String,
      headers: Array,
      items: Array,
      editAction: {
        type: Boolean,
        default: true,
      },
      deleteAction: {
        type: Boolean,
        default: true,
      },
    },
    methods: {
      onEdit(e) {
        this.$emit('edit', e)
      },
      onDelete(e) {
        this.$emit('delete', e)
      },
    },
  }
</script>

<style></style>
