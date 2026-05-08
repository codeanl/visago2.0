<template>
  <view class="header-shell">
    <view class="page-wrap header-wrap">
      <view class="header-row">
        <view v-if="showBack" class="back-btn" @tap="goBack">&lt;</view>
        <view v-else class="back-placeholder" />
        <view class="header-copy">
          <text class="header-title">{{ title }}</text>
          <text v-if="subtitle" class="header-subtitle">{{ subtitle }}</text>
        </view>
        <view class="back-placeholder" />
      </view>
    </view>
  </view>
</template>

<script>
export default {
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
  methods: {
    goBack() {
      const pages = typeof getCurrentPages === 'function' ? getCurrentPages() : []
      if (pages.length > 1) {
        uni.navigateBack()
        return
      }
      uni.reLaunch({ url: '/pages/home/index' })
    },
  },
}
</script>

<style scoped>
.header-shell {
  position: sticky;
  top: 0;
  z-index: 20;
  margin: 0 -16px;
  padding: 14px 16px 8px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98) 0%, rgba(255, 255, 255, 0.82) 100%);
  backdrop-filter: blur(16px);
}

.header-wrap {
  width: 100%;
}

.header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.back-btn,
.back-placeholder {
  width: 40px;
  height: 40px;
  flex: 0 0 40px;
}

.back-btn {
  border-radius: 12px;
  background: #fff;
  border: 1px solid var(--lite-line);
  color: var(--lite-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26px;
  line-height: 1;
  font-weight: 800;
}

.header-copy {
  flex: 1;
  min-width: 0;
  text-align: center;
}

.header-title {
  display: block;
  font-size: 22px;
  font-weight: 900;
  letter-spacing: 0.5px;
}

.header-subtitle {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.5;
  color: var(--lite-text-muted);
}
</style>

