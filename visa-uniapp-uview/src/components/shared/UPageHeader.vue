<template>
  <view class="u-header-shell">
    <view class="app-page-width u-header-bar">
      <view class="u-header-side">
        <view v-if="showBack" class="u-header-btn" @tap="goBack">
          <u-icon name="arrow-left" size="20" color="#4e5969" />
        </view>
        <image v-else-if="home" class="u-header-logo" :src="logoSrc" mode="heightFix" />
      </view>

      <view class="u-header-title-wrap">
        <text class="u-header-title">{{ title }}</text>
      </view>

      <view class="u-header-side u-header-side--right">
        <slot v-if="!showMockCapsule" name="right" />
        <view v-else class="u-header-capsule">
          <view class="u-header-capsule__dots">
            <text class="u-header-capsule__dot" />
            <text class="u-header-capsule__dot" />
            <text class="u-header-capsule__dot" />
          </view>
          <view class="u-header-capsule__divider" />
          <view class="u-header-capsule__circle" />
        </view>
      </view>
    </view>
  </view>
</template>

<script>
export default {
  name: 'UPageHeader',
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
  data() {
    return {
      showMockCapsule: false,
    }
  },
  created() {
    try {
      const info = typeof uni !== 'undefined' && typeof uni.getSystemInfoSync === 'function'
        ? uni.getSystemInfoSync()
        : {}
      const platform = String(info.uniPlatform || info.platform || '').toLowerCase()
      this.showMockCapsule = platform !== 'mp-weixin'
    } catch (error) {
      this.showMockCapsule = true
    }
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
.u-header-shell {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 90;
  padding-top: calc(var(--uview-safe-top) + 2px);
  background: rgba(255, 255, 255, 0.98);
  border-bottom: 1px solid rgba(15, 24, 36, 0.05);
  box-shadow: 0 1px 8px rgba(16, 24, 40, 0.03);
}

.u-header-bar {
  position: relative;
  height: 74px;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding: 0 16px 12px;
}

.u-header-side {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.u-header-side--right {
  width: 96px;
  justify-content: flex-end;
}

.u-header-title-wrap {
  position: absolute;
  left: 50%;
  bottom: 13px;
  transform: translateX(-50%);
  min-width: 120px;
  max-width: calc(100% - 120px);
  display: flex;
  align-items: center;
  justify-content: center;
}

.u-header-title {
  font-size: 18px;
  line-height: 1.2;
  font-weight: 800;
  color: var(--uview-text);
  text-align: center;
}

.u-header-btn {
  width: 40px;
  height: 40px;
  border-radius: 999px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f7f8fa;
  border: 1px solid var(--uview-border);
}

.u-header-logo {
  height: 32px;
  width: auto;
  max-width: 220px;
}

.u-header-capsule {
  width: 88px;
  height: 34px;
  border-radius: 999px;
  background: #ffffff;
  border: 1px solid rgba(20, 37, 63, 0.14);
  box-shadow: 0 2px 10px rgba(16, 24, 40, 0.04);
  display: flex;
  align-items: center;
  justify-content: center;
}

.u-header-capsule__dots {
  width: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.u-header-capsule__dot {
  width: 4px;
  height: 4px;
  border-radius: 999px;
  background: #1f2329;
  margin: 0 2px;
}

.u-header-capsule__divider {
  width: 1px;
  height: 14px;
  background: rgba(20, 37, 63, 0.12);
}

.u-header-capsule__circle {
  width: 18px;
  height: 18px;
  border-radius: 999px;
  border: 2px solid #1f2329;
  margin-left: 12px;
  position: relative;
}

.u-header-capsule__circle::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 4px;
  height: 4px;
  border-radius: 999px;
  background: #1f2329;
  transform: translate(-50%, -50%);
}
</style>
