<template>
  <view class="top-shell">
    <view class="visago-page-width top-bar">
      <view class="left-zone">
        <view v-if="showBack" class="back-button" hover-class="back-button--hover" @tap="handleBack">
          <text class="material-symbols-outlined back-icon">arrow_back_ios_new</text>
        </view>
        <image v-else-if="home" class="brand-logo" :src="logoSrc" mode="heightFix" />
      </view>

      <view v-if="!home" class="title-center">
        <text class="menu-title">{{ pageName }}</text>
      </view>

      <view v-if="showNotice" class="notice-button" hover-class="notice-button--hover" @tap="openNoticePage">
        <text class="material-symbols-outlined notice-icon">notifications</text>
        <view v-if="showNoticeDot" class="notice-dot" />
      </view>
      <view v-else class="notice-placeholder" />
    </view>
  </view>
</template>

<script>
export default {
  name: 'VisagoTopBar',
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
  methods: {
    handleBack() {
      if (this.customBack) {
        this.$emit('back')
        return
      }
      this.goBack()
    },
    goBack() {
      const pages = getCurrentPages()
      if (pages.length > 1) {
        uni.navigateBack({ delta: 1 })
      } else {
        uni.reLaunch({ url: '/pages/home/index' })
      }
    },
    openNoticePage() {
      return
    },
  },
}
</script>

<style scoped>
.top-shell {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: var(--visago-surface);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--visago-line);
}

.top-bar {
  position: relative;
  height: 74px;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding: 0 16px 12px;
  box-sizing: border-box;
}

.left-zone {
  min-width: 150px;
  min-height: 40px;
  display: flex;
  align-items: center;
  justify-content: flex-start;
}

.back-button {
  width: 40px;
  height: 40px;
  border-radius: 9999px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-button--hover {
  background: var(--visago-surface-soft);
}

.back-icon {
  font-size: 22px;
  color: var(--visago-text-soft);
}

.title-center {
  position: absolute;
  left: 50%;
  bottom: 12px;
  transform: translateX(-50%);
  min-height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.brand-logo {
  height: 32px;
  width: auto;
  max-width: 240px;
}

.menu-title {
  font-size: 34rpx;
  font-weight: 700;
  color: var(--visago-text);
  letter-spacing: 0.02em;
}

.notice-button,
.notice-placeholder {
  width: 40px;
  height: 40px;
}

.notice-button {
  border-radius: 9999px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.notice-button--hover {
  background: var(--visago-surface-soft);
}

.notice-icon {
  font-size: 24px;
  color: var(--visago-text-soft);
}

.notice-dot {
  width: 8px;
  height: 8px;
  border-radius: 9999px;
  position: absolute;
  right: 9px;
  top: 9px;
  background: var(--visago-danger);
  border: 1px solid var(--visago-surface);
}
</style>

