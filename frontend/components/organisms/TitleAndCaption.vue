<template>
  <span>
    <div class="d-flex justify-space-between align-center">
      <div class="text-h5">
        {{ title }}
      </div>
      <slot name="title"></slot>
    </div>
    <div class="d-flex justify-space-betweeb align-center">
      <div class="text-caption">
        {{ caption }}
      </div>
      <slot name="caption"></slot>
    </div>
  </span>
</template>

<script>
  export default {
    name: 'titleAndCaption',
    props: {
      title: String,
      caption: String,
    },
  }
</script>

<style></style>
