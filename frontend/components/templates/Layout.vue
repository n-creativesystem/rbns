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

      <template v-if="$rbns.isTenant">
        <v-divider class="mx-4" inset vertical></v-divider>
        <n-menu>
          <template v-slot:activator="{ on, attrs }">
            <v-btn class="tfn" tile text v-bind="attrs" v-on="on"
              >{{ selectedTenant || $t('tenant.name')
              }}<v-icon>mdi-menu-down</v-icon></v-btn
            >
          </template>
          <template>
            <v-list flat>
              <v-subheader>{{ $t('tenant.name') }}</v-subheader>
              <v-list-item-group v-model="tenant" color="info">
                <v-list-item v-for="(item, idx) in tenants" :key="idx">
                  <v-list-item-content>
                    <v-list-item-title>{{ item.Name }}</v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </template>
        </n-menu>
      </template>

      <v-spacer></v-spacer>

      <n-menu close-delay="100" open-delay="60" open-on-hover>
        <template v-slot:activator="{ on, attrs }">
          <v-btn class="tfn" tile text v-bind="attrs" v-on="on"
            ><v-icon>mdi-translate</v-icon><v-icon>mdi-menu-down</v-icon></v-btn
          >
        </template>
        <template>
          <v-list flat>
            <v-subheader>{{ $t('translations') }}</v-subheader>
            <v-list-item-group v-model="localSelect" color="info">
              <v-list-item v-for="(i, idx) in locales" :key="idx">
                <v-list-item-content>
                  <v-list-item-title>{{ i.title }}</v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list-item-group>
          </v-list>
        </template>
      </n-menu>
    </v-app-bar>
    <v-main>
      <v-container fluid>
        <router-view></router-view>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
  import Menu from '../../plugins/menu'
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
      tenant: undefined,
    }),
    created() {
      const tenantId = this.$rbns.user.Tenant
      if (tenantId) {
        this.tenant = this.tenants.indexOf(
          this.tenants.find((t) => t.ID === tenantId)
        )
      }
    },
    watch: {
      localSelect(val) {
        const locale = this.locales.find((_, idx) => idx === val)
        this.$vuetify.lang.current = locale.alternate || locale.locale
        this.$i18n.locale = locale.alternate || locale.locale
      },
    },
    computed: {
      tenants() {
        return this.$rbns.tenants
      },
      selectedTenant() {
        return this.tenant >= 0 ? this.tenants[this.tenant].Name : undefined
      },
    },
  }
</script>
