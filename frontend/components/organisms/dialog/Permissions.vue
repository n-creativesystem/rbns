<template>
  <n-dialog
    v-if="innerValue"
    v-model="innerValue"
    v-bind="$attrs"
    :title="$t('permission.dialog.title')"
    @open="onOpen"
  >
    <dialog-selection
      :text="$t('permission.dialog.add')"
      :items.sync="items"
      :selection.sync="innerSelection"
      @click="onClick"
    ></dialog-selection>
  </n-dialog>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'dialogPermissions',
    mixins: [axiosMixin],
    data() {
      return {
        items: [],
      }
    },
    props: {
      value: Boolean,
      ommitIds: Array,
      selection: Array,
    },
    computed: {
      innerValue: {
        get() {
          return this.value
        },
        set(val) {
          this.$emit('input', val)
        },
      },
      innerSelection: {
        get() {
          return this.selection
        },
        set(val) {
          this.$emit('update:selection', val)
        },
      },
    },
    methods: {
      onClick(items) {
        this.$emit('click', items)
        this.innerValue = false
      },
      getData() {
        this.get(this.$urls.api.v1.permissions)
          .then((result) => {
            if (result.data.permissions) {
              this.items = result.data.permissions.filter(
                (i) => !(this.ommitIds || []).some((id) => i.id === id)
              )
            }
          })
          .catch((err) => {
            console.log(err)
          })
      },
      onOpen() {
        this.getData()
      },
    },
  }
</script>

<style></style>
