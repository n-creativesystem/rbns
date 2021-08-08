<template>
  <n-dialog
    v-if="innerValue"
    v-model="innerValue"
    :title="innerTitle"
    :title-class="innerTitleClass"
    v-bind="$attrs"
    @opened="onOpened"
  >
    <dialog-input-base @submit="onClick" ref="input">
      <name-and-description
        :name.sync="returnData.name"
        name-class="col-12"
        :description.sync="returnData.description"
        description-class="col-12"
      ></name-and-description>
    </dialog-input-base>
  </n-dialog>
</template>

<script>
  export default {
    name: 'dialogNameAndDesc',
    data() {
      return {
        returnData: {
          name: this.name,
          description: this.description,
        },
      }
    },
    props: {
      value: Boolean,
      title: String,
      titleClass: String,
      name: String,
      description: String,
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
      innerTitle() {
        return this.title || 'title'
      },
      innerTitleClass() {
        return this.titleClass || 'text-h5 grey lighten-2'
      },
    },
    methods: {
      onClick() {
        this.$emit('submit', this.returnData)
      },
      onOpened() {
        if (this.$refs.input) {
          this.$refs.input.reset()
        }
      },
    },
  }
</script>

<style></style>
