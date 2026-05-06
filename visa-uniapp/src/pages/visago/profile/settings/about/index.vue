<template>
  <view class="about-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="关于 Visago" />

    <scroll-view scroll-y class="about-scroll">
      <view class="about-content visago-page-width">
        <view class="brand-row">
          <image class="brand-logo" src="/static/logo.png" mode="aspectFit" />
          <view class="brand-copy">
            <text class="brand-name">Visago</text>
            <text class="brand-version">版本 1.0.1</text>
          </view>
        </view>

        <view class="list-card">
          <view class="list-row">
            <view class="row-left">
              <text class="material-symbols-outlined row-icon">info</text>
              <text class="row-title">当前版本</text>
            </view>
            <text class="row-value">1.0.1</text>
          </view>

          <view class="list-row" @tap="onContactTap(contactItems[0])">
            <view class="row-left">
              <text class="material-symbols-outlined row-icon">mail</text>
              <text class="row-title">联系邮箱</text>
            </view>
            <text class="row-value">jiaanliao23@gmail.com</text>
          </view>

          <view class="list-row" @tap="onContactTap(contactItems[1])">
            <view class="row-left">
              <text class="material-symbols-outlined row-icon">public</text>
              <text class="row-title">官方网站</text>
            </view>
            <text class="row-value">visago.nova2026.top</text>
          </view>

          <view class="list-row" @tap="onLegalTap('agreement')">
            <view class="row-left">
              <text class="material-symbols-outlined row-icon">gavel</text>
              <text class="row-title">用户协议</text>
            </view>
            <text class="material-symbols-outlined row-arrow">chevron_right</text>
          </view>

          <view class="list-row" @tap="onLegalTap('privacy')">
            <view class="row-left">
              <text class="material-symbols-outlined row-icon">verified_user</text>
              <text class="row-title">隐私政策</text>
            </view>
            <text class="material-symbols-outlined row-arrow">chevron_right</text>
          </view>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../../components/VisagoTopBar.vue'
import { applyTheme, getStoredTheme } from '../../../../../utils/theme'

export default {
  components: { VisagoTopBar },
  data() {
    return {
      contactItems: [
        { key: 'email', value: 'jiaanliao23@gmail.com' },
        { key: 'site', value: 'visago.nova2026.top' },
      ],
    }
  },
  onShow() {
    applyTheme(getStoredTheme())
  },
  methods: {
    onContactTap(item) {
      if (!item || !item.value) return
      uni.setClipboardData({
        data: item.value,
        fail: () => {
          uni.showToast({ title: '复制失败', icon: 'none' })
        },
      })
    },
    onLegalTap(key) {
      const routeMap = {
        agreement: '/pages/visago/profile/settings/about/agreement/index',
        privacy: '/pages/visago/profile/settings/about/privacy-policy/index',
      }
      const url = routeMap[key]
      if (!url) return
      uni.navigateTo({
        url,
        fail: () => {
          uni.showToast({ title: '打开失败', icon: 'none' })
        },
      })
    },
  },
}
</script>

<style scoped>
.about-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.about-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.about-content {
  padding: 18px 16px 34px;
  box-sizing: border-box;
}

.brand-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.brand-logo {
  width: 64px;
  height: 64px;
  border-radius: 18px;
  background: #fff;
  border: 1px solid var(--visago-line);
}

.brand-name,
.brand-version,
.row-title,
.row-value {
  display: block;
}

.brand-name {
  font-size: 20px;
  font-weight: 900;
  color: var(--visago-text);
}

.brand-version {
  margin-top: 4px;
  font-size: 12px;
  color: var(--visago-text-muted);
}

.list-card {
  overflow: hidden;
  border-radius: 18px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
}

.list-row {
  min-height: 56px;
  padding: 0 13px;
  border-bottom: 1px solid var(--visago-line);
  display: flex;
  align-items: center;
  gap: 10px;
}

.list-row:last-child {
  border-bottom: none;
}

.row-left {
  min-width: 0;
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
}

.row-icon {
  color: var(--visago-primary);
  font-size: 20px;
}

.row-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--visago-text);
}

.row-value {
  max-width: 42%;
  text-align: right;
  font-size: 12px;
  color: var(--visago-text-muted);
}

.row-arrow {
  color: var(--visago-text-soft);
}
</style>
