<template>
  <v-app>
    <side-menu :links="Menu" v-model="drawer" clipped></side-menu>
    <v-app-bar app clippedLeft>
      <v-app-bar-nav-icon @click="drawer = !drawer">
        <v-icon>{{
          drawer ? 'mdi-dots-vertical' : 'mdi-format-list-bulleted'
        }}</v-icon>
      </v-app-bar-nav-icon>

      <v-toolbar-title>
        <router-link to="/">
          <v-img :src="require('@/assets/logo.svg')"></v-img>
        </router-link>
        <!-- <v-btn class="tfn text-h6" text to="/">Role Based Access Control</v-btn> -->
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <v-menu
        bottom
        offset-y
        open-on-hover
        open-delay="60"
        max-height="500"
        close-delay="100"
        content-class="rounded"
        left
        transition="slide-y-transition"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-btn tile text v-bind="attrs" v-on="on"
            ><v-icon>mdi-translate</v-icon><v-icon>mdi-menu-down</v-icon></v-btn
          >
        </template>
        <v-card class="mx-auto" max-width="300" tile>
          <v-list flat>
            <v-subheader>{{ $t('translations') }}</v-subheader>
            <v-list-item-group v-model="localSelect" color="info">
              <v-list-item v-for="(item, idx) in locales" :key="idx">
                <v-list-item-content>
                  <v-list-item-title>{{ item.title }}</v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list-item-group>
          </v-list>
        </v-card>
      </v-menu>
    </v-app-bar>
    <v-main>
      <v-container fluid>
        <router-view></router-view>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
  import Menu from '@assets/menu.json'
  export default {
    name: 'Layout',
    data: () => ({
      Menu,
      drawer: null,
      localSelect: 0,
      locales: [
        {
          title: '日本語',
          locale: 'ja-JP',
          alternate: 'ja',
        },
        {
          title: 'English',
          locale: 'en',
        },
      ],
    }),
    watch: {
      localSelect(val) {
        const locale = this.locales.find((_, idx) => idx === val)
        this.$vuetify.lang.current = locale.alternate || locale.locale
        this.$i18n.locale = locale.alternate || locale.locale
      },
    },
  }
</script>
