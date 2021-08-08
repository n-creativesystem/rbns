<template>
  <form-parts class="input-and-list">
    <input-form
      v-if="addShow"
      :title="addMeta.title"
      :caption="addMeta.caption"
      ref="form"
      :name.sync="innerName"
      :description.sync="innerDescription"
      @submit="onSubmit"
    >
    </input-form>
    <list-form
      v-if="listShow"
      :title="listMeta.title"
      :caption="listMeta.caption"
      :headers="headers"
      :items="items"
      :search="search"
      :edit-action="editAction"
      :delete-action="deleteAction"
      @edit="onEdit"
      @delete="onDelete"
    >
      <template v-slot:title>
        <slot name="list-title"></slot>
      </template>
      <template v-slot:caption>
        <slot name="list-caption"></slot>
      </template>
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
    </list-form>
  </form-parts>
</template>

<script>
  import loadingMixin from '@mixin/loading'
  export default {
    name: 'inputAndList',
    mixins: [loadingMixin],
    data() {
      return {
        search: null,
      }
    },
    props: {
      addShow: {
        type: Boolean,
        default: true,
      },
      listShow: {
        type: Boolean,
        default: true,
      },
      addMeta: {
        type: Object,
        default: () => {
          return {
            title: '',
            caption: '',
          }
        },
      },
      listMeta: {
        type: Object,
        default: () => {
          return {
            title: '',
            caption: '',
          }
        },
      },
      name: String,
      description: String,
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
    computed: {
      innerName: {
        get() {
          return this.name
        },
        set(val) {
          this.$emit('update:name', val)
        },
      },
      innerDescription: {
        get() {
          return this.description
        },
        set(val) {
          this.$emit('update:description', val)
        },
      },
    },
    methods: {
      onSubmit(e) {
        this.$emit('submit', e)
      },
      onEdit(e) {
        this.$emit('edit', e)
      },
      onDelete(e) {
        this.$emit('delete', e)
      },
      reset() {
        this.$refs.form.reset()
      },
    },
  }
</script>

<style></style>
