<template>
  <n-data-table :headers="headers" :items="items">
    <template v-slot:[`item.userKey`]="{ item }">
      <router-link
        :to="
          '/organizations/{0}/users/{1}'.format(
            item.organizationId,
            item.userKey
          )
        "
      >
        {{ item.userKey }}
      </router-link>
    </template>
    <template v-slot:[`item.organizationName`]="{ item }">
      <router-link :to="'/organizations/{0}'.format(item.organizationId)">
        {{ item.organizationName }}
      </router-link>
    </template>
  </n-data-table>
</template>

<script>
  export default {
    name: 'RoleTabUsers',
    data() {
      return {}
    },
    props: {
      items: Array,
    },
    computed: {
      headers() {
        return [
          {
            text: this.$t('user.entity.userId'),
            value: 'userKey',
            align: 'start',
          },
          {
            text: this.$t('organization.entity.name'),
            value: 'organizationName',
            align: 'start',
          },
          {
            text: this.$t('organization.entity.description'),
            value: 'organizationDescription',
            align: 'start',
          },
        ]
      },
    },
  }
</script>

<style></style>
