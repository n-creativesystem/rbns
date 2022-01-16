<template>
  <v-menu
    bottom
    left
    max-height="500"
    offset-y
    content-class="rounded"
    transition="slide-y-transition"
    v-bind="$attrs"
    v-on="$listeners"
  >
    <template #activator="props">
      <slot name="activator" v-bind="{ ...props }" />
    </template>

    <n-sheet :outlined="false">
      <slot v-if="$slots.default" />
      <v-list v-else flat>
        <v-subheader>{{ header }}</v-subheader>
        <v-list-item-group color="info">
          <v-list-item v-for="(item, idx) in items" :key="idx">
            <v-list-item-content>
              <v-list-item-title>{{ item.title }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </n-sheet>
  </v-menu>
</template>
<script>
  export default {
    name: 'NMenu',
    props: {
      header: String,
      items: {
        type: Array,
        default: () => [],
      },
    },
  }
</script>
<style lang=""></style>
