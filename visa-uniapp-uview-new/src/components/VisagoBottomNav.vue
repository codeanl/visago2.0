<template>
  <VuTabbar :active="mappedActive" />
</template>

<script>
import VuTabbar from './uview/VuTabbar.vue'

export default {
  name: 'VisagoBottomNav',
  components: { VuTabbar },
  props: {
    activeTab: {
      type: String,
      default: 'home',
    },
  },
  computed: {
    mappedActive() {
      if (this.activeTab === 'plan') return 'goals'
      return this.activeTab
    },
  },
}
</script>
