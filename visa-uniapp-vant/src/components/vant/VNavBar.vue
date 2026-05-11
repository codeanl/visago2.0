<template>
  <view class="v-nav-shell">
    <view class="v-page-width v-nav-bar">
      <view class="v-nav-side v-nav-side--left">
        <view v-if="showBack" class="v-nav-btn" @tap="goBack">
          <VIcon class="v-nav-icon" name="arrow-left" />
        </view>
        <image v-else-if="home" class="v-nav-logo" :src="logoSrc" mode="heightFix" />
      </view>

      <view class="v-nav-title-wrap">
        <text class="v-nav-title">{{ title }}</text>
      </view>

      <view class="v-nav-side v-nav-side--right">
        <slot name="right" />
      </view>
    </view>
  </view>
</template>

<script>
import VIcon from './VIcon.vue'

export default {
  name: 'VNavBar',
  components: { VIcon },
  props: {
    title: {
      type: String,
      default: '',
    },
    home: {
      type: Boolean,
      default: false,
    },
    showBack: {
      type: Boolean,
      default: false,
    },
    logoSrc: {
      type: String,
      default: '/static/header-logo.png',
    },
  },
  methods: {
    goBack() {
      const pages = typeof getCurrentPages === 'function' ? getCurrentPages() : []
      if (pages.length > 1) {
        uni.navigateBack({ delta: 1 })
        return
      }
      uni.reLaunch({ url: '/pages/home/index' })
    },
  },
}
</script>

<style scoped>
.v-nav-shell {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 90;
  padding-top: calc(var(--vant-safe-top) + 2px);
  background: rgba(255, 255, 255, 0.98);
  border-bottom: 1px solid rgba(15, 24, 36, 0.06);
  box-shadow: 0 2px 12px rgba(16, 24, 40, 0.04);
}

.v-nav-bar {
  position: relative;
  height: 74px;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding: 0 16px 12px;
  box-sizing: border-box;
}

.v-nav-side {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.v-nav-title-wrap {
  position: absolute;
  left: 50%;
  bottom: 13px;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 120px;
  max-width: calc(100% - 120px);
}

.v-nav-title {
  font-size: 18px;
  line-height: 1.2;
  font-weight: 800;
  color: #1f2329;
  text-align: center;
}

.v-nav-btn {
  width: 40px;
  height: 40px;
  border-radius: 999px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f7f8fa;
}

.v-nav-icon {
  font-size: 22px;
  color: #607089;
}

.v-nav-logo {
  height: 32px;
  width: auto;
  max-width: 220px;
}
</style>
