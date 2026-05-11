<template>
  <VuTabbar :active="active" />
</template>

<script>
import VuTabbar from './uview/VuTabbar.vue'

export default {
  name: 'BottomDock',
  components: { VuTabbar },
  props: {
    active: {
      type: String,
      default: 'home',
    },
  },
}
</script>
