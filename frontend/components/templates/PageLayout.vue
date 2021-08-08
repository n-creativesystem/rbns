<template>
  <v-container>
    <slot></slot>
    <n-snackbar-error
      v-model="errorNotification"
      :message="notificationMessage"
    ></n-snackbar-error>
    <n-snackbar-information
      v-model="infoNotification"
      :message="notificationMessage"
    ></n-snackbar-information>
    <n-snackbar-success
      v-model="successNotification"
      :message="notificationMessage"
    ></n-snackbar-success>
  </v-container>
</template>

<script>
  import NotificationMixin from '@mixin/notification'
  export default {
    name: 'PageLayout',
    mixins: [NotificationMixin],
    data() {
      return {}
    },
  }
</script>

<style></style>
