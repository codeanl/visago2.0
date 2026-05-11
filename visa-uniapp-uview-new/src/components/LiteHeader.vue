<template>
  <VuPageHeader :title="title" :show-back="showBack" />
</template>

<script>
import VuPageHeader from './uview/VuPageHeader.vue'

export default {
  name: 'LiteHeader',
  components: { VuPageHeader },
  props: {
    title: {
      type: String,
      default: '',
    },
    subtitle: {
      type: String,
      default: '',
    },
    showBack: {
      type: Boolean,
      default: false,
    },
  },
}
</script>
