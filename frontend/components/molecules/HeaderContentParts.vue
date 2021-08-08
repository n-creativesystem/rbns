<template>
  <div class="row">
    <div class="col-12">
      <slot name="header"></slot>
      <slot name="content"></slot>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'headerContentParts',
  }
</script>

<style></style>
