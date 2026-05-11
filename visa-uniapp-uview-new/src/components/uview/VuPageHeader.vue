<template>
  <view class="vu-header-shell" :style="headerStyle">
    <view class="page-width vu-header-bar" :style="barStyle">
      <view class="vu-header-left">
        <view v-if="showBack" class="vu-header-btn" hover-class="vu-header-btn--hover" @tap="handleBack">
          <u-icon name="arrow-left" size="20" color="#4b5563" />
        </view>
        <image v-else-if="home" class="vu-header-logo" :src="logoSrc" mode="heightFix" />
        <slot v-else name="left" />
      </view>

      <view v-if="title" class="vu-header-title-wrap">
        <text class="vu-header-title">{{ title }}</text>
      </view>

      <view class="vu-header-right">
        <slot name="right" />
        <view v-if="showNotice" class="vu-header-btn" hover-class="vu-header-btn--hover" @tap="openNotice">
          <u-icon name="bell" size="20" color="#4b5563" />
          <view v-if="showNoticeDot" class="vu-notice-dot" />
        </view>
        <view v-else-if="showMockCapsule" class="vu-capsule">
          <view class="vu-capsule-dots">
            <view class="vu-capsule-dot" />
            <view class="vu-capsule-dot" />
            <view class="vu-capsule-dot" />
          </view>
          <view class="vu-capsule-divider" />
          <view class="vu-capsule-circle" />
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import { getLayoutMetrics } from '../../utils/layout'

export default {
  name: 'VuPageHeader',
  emits: ['back', 'notice'],
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
    customBack: {
      type: Boolean,
      default: false,
    },
    showNotice: {
      type: Boolean,
      default: false,
    },
    showNoticeDot: {
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
      metrics: getLayoutMetrics(),
    }
  },
  computed: {
    headerStyle() {
      return `height:${this.metrics.headerHeight}px;padding-top:${this.metrics.safeTop}px;`
    },
    barStyle() {
      const height = Math.max(48, this.metrics.headerHeight - this.metrics.safeTop)
      return `height:${height}px;`
    },
    showMockCapsule() {
      return this.metrics.platform !== 'mp-weixin'
    },
  },
  methods: {
    handleBack() {
      if (this.customBack) {
        this.$emit('back')
        return
      }
      const pages = typeof getCurrentPages === 'function' ? getCurrentPages() : []
      if (pages.length > 1) {
        uni.navigateBack({ delta: 1 })
        return
      }
      uni.reLaunch({ url: '/pages/home/index' })
    },
    openNotice() {
      this.$emit('notice')
    },
  },
}
</script>

<style scoped>
.vu-header-shell {
  position: fixed;
  top: 0;
  right: 0;
  left: 0;
  z-index: 90;
  background: #ffffff;
  border-bottom: 1px solid #eef2f7;
  box-shadow: 0 2px 12px rgba(31, 35, 41, 0.04);
}

.vu-header-bar {
  position: relative;
  padding-left: 16px;
  padding-right: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.vu-header-left,
.vu-header-right {
  width: 102px;
  height: 44px;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.vu-header-left {
  justify-content: flex-start;
}

.vu-header-right {
  justify-content: flex-end;
}

.vu-header-title-wrap {
  position: absolute;
  right: 112px;
  left: 112px;
  top: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.vu-header-title {
  max-width: 100%;
  color: #1f2329;
  font-size: 18px;
  line-height: 1.25;
  font-weight: 800;
  text-align: center;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.vu-header-logo {
  height: 32px;
  width: 160px;
}

.vu-header-btn {
  position: relative;
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background: #f4f7fb;
  display: flex;
  align-items: center;
  justify-content: center;
}

.vu-header-btn--hover {
  background: #e9eef7;
}

.vu-notice-dot {
  position: absolute;
  right: 9px;
  top: 8px;
  width: 7px;
  height: 7px;
  border-radius: 4px;
  background: #ee0a24;
  border: 1px solid #ffffff;
}

.vu-capsule {
  width: 88px;
  height: 34px;
  border-radius: 17px;
  background: #ffffff;
  border: 1px solid rgba(20, 37, 63, 0.14);
  box-shadow: 0 2px 10px rgba(31, 35, 41, 0.04);
  display: flex;
  align-items: center;
  justify-content: center;
}

.vu-capsule-dots {
  width: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.vu-capsule-dot {
  width: 4px;
  height: 4px;
  border-radius: 2px;
  background: #1f2329;
  margin-left: 2px;
  margin-right: 2px;
}

.vu-capsule-divider {
  width: 1px;
  height: 14px;
  background: rgba(20, 37, 63, 0.14);
}

.vu-capsule-circle {
  width: 18px;
  height: 18px;
  border-radius: 9px;
  border: 2px solid #1f2329;
  margin-left: 12px;
}
</style>
