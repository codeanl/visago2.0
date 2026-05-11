<template>
  <VuPageHeader
    :title="pageName"
    :home="home"
    :show-back="showBack"
    :custom-back="customBack"
    :show-notice="showNotice"
    :show-notice-dot="showNoticeDot"
    :logo-src="logoSrc"
    @back="$emit('back')"
  />
</template>

<script>
import VuPageHeader from './uview/VuPageHeader.vue'

export default {
  name: 'VisagoTopBar',
  components: { VuPageHeader },
  emits: ['back'],
  props: {
    home: {
      type: Boolean,
      default: false,
    },
    pageName: {
      type: String,
      default: '',
    },
    showNoticeDot: {
      type: Boolean,
      default: false,
    },
    showNotice: {
      type: Boolean,
      default: false,
    },
    showBack: {
      type: Boolean,
      default: false,
    },
    customBack: {
      type: Boolean,
      default: false,
    },
    logoSrc: {
      type: String,
      default: '/static/header-logo.png',
    },
  },
}
</script>
