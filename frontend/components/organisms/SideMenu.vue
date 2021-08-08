<template>
  <v-navigation-drawer :clipped="clipped" v-model="innerValue" app>
    <v-list>
      <v-list-item> </v-list-item>
      <v-list-item v-for="link in links" :key="link.key" :to="link.href">
        <v-list-item-icon v-if="link.icon">
          <v-icon>{{ link.icon }}</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>{{ $t(`menu.${link.title}`) }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
  export default {
    name: 'SideMenu',
    data() {
      return {}
    },
    props: {
      value: Boolean,
      links: Array,
      clipped: Boolean,
    },
    computed: {
      innerValue: {
        get() {
          return this.value
        },
        set(value) {
          this.$emit('input', value)
        },
      },
    },
  }
</script>

<style></style>
