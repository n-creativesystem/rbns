<template>
  <n-dialog
    v-if="innerValue"
    v-model="innerValue"
    v-bind="$attrs"
    @open="onOpen"
  >
    <v-card>
      <n-form @submit="onSubmit">
        <div class="col-12">
          <v-subheader>{{
            $t('organization.dialog.create.user.title')
          }}</v-subheader>
          <required-text
            name="userId"
            :label="$t('user.entity.userId')"
            id="userId"
            :counter="255"
            v-model="innerUserId"
            :disabled="editMode"
          ></required-text>
        </div>
        <div class="col-12">
          <v-subheader>{{
            $t('organization.dialog.create.role.title')
          }}</v-subheader>
          <v-card outlined flat class="pa-4">
            <n-selection :items="items" v-model="innerSelection"> </n-selection>
          </v-card>
        </div>
        <div class="col-12">
          <btn-tfn type="submit">
            <v-icon>mdi-add</v-icon>
            {{ $t('save') }}
          </btn-tfn>
        </div>
      </n-form>
    </v-card>
  </n-dialog>
</template>

<script>
  import axiosMixin from '@mixin/axios'
  export default {
    name: 'DialogUsers',
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
      userId: String,
      editMode: {
        type: Boolean,
        default: false,
      },
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
      innerUserId: {
        get() {
          return this.userId
        },
        set(val) {
          this.$emit('update:userId', val)
        },
      },
      getItems() {
        return (
          (this.items &&
            this.items.filter((_, idx) =>
              this.innerSelection.some((i) => i === idx)
            )) ||
          []
        )
      },
    },
    methods: {
      getData() {
        this.get(this.$urls.api.v1.roles)
          .then((result) => {
            if (result.data.roles) {
              this.items = result.data.roles.filter((i) =>
                this.ommitIds && this.ommitIds.length > 0
                  ? this.ommitIds.some((id) => i.id !== id)
                  : true
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
      onSubmit() {
        this.$emit('submit', {
          userId: this.innerUserId,
          selection: this.getItems,
        })
        this.innerValue = false
        this.innerUserId = ''
        this.innerSelection = []
      },
    },
  }
</script>

<style></style>
