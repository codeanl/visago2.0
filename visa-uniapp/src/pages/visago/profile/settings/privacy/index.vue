<template>
  <view class="settings-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="隐私与安全" />

    <scroll-view scroll-y class="settings-scroll">
      <view class="settings-content visago-page-width">
        <view v-for="group in groups" :key="group.key" class="group">
          <text class="group-title">{{ group.title }}</text>
          <view class="group-card">
            <view v-for="(item, idx) in group.items" :key="item.key" class="row" @tap="handleRow(item)">
              <view class="row-left">
                <view class="row-icon" :style="{ color: item.color, background: item.tint }">
                  <text class="material-symbols-outlined">{{ item.icon }}</text>
                </view>
                <text class="row-title" :class="{ 'row-title--danger': item.danger }">{{ item.title }}</text>
              </view>

              <view class="row-right">
                <text v-if="item.value" class="row-value">{{ item.value }}</text>
                <view v-if="item.toggle" class="mini-switch" :class="{ 'mini-switch--on': item.on }">
                  <view class="mini-thumb" />
                </view>
                <text v-else class="material-symbols-outlined row-arrow">chevron_right</text>
              </view>

              <view v-if="idx < group.items.length - 1" class="divider" />
            </view>
          </view>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../../components/VisagoTopBar.vue'
import { applyTheme, getStoredTheme } from '../../../../../utils/theme'

const STORAGE_KEY = 'visago_privacy_preferences'

function buildGroups() {
  return [
    {
      key: 'account',
      title: '账户与登录',
      items: [
        { key: 'login-alert', title: '新设备登录提醒', icon: 'devices', toggle: true, on: true, color: '#0891b2', tint: 'rgba(8,145,178,0.12)' },
        { key: 'sensitive-confirm', title: '敏感操作二次确认', icon: 'shield', toggle: true, on: true, color: '#16a34a', tint: 'rgba(22,163,74,0.12)' },
      ],
    },
    {
      key: 'data',
      title: '数据管理',
      items: [
        { key: 'policy', title: '查看隐私政策', icon: 'policy', color: '#0f65d8', tint: 'rgba(15,101,216,0.12)' },
        { key: 'delete', title: '注销账户', icon: 'delete_forever', danger: true, color: '#ef4444', tint: 'rgba(239,68,68,0.12)' },
      ],
    },
  ]
}

export default {
  components: { VisagoTopBar },
  data() {
    return {
      groups: buildGroups(),
    }
  },
  onLoad() {
    this.restorePreferences()
  },
  onShow() {
    applyTheme(getStoredTheme())
  },
  methods: {
    restorePreferences() {
      const stored = uni.getStorageSync(STORAGE_KEY)
      if (!stored) return

      const parsed = typeof stored === 'string' ? JSON.parse(stored) : stored
      if (!parsed || typeof parsed !== 'object') return

      this.groups.forEach((group) => {
        group.items.forEach((item) => {
          if (!item.toggle) return
          if (Object.prototype.hasOwnProperty.call(parsed, item.key)) {
            item.on = Boolean(parsed[item.key])
          }
        })
      })
    },
    persistPreferences() {
      const payload = {}
      this.groups.forEach((group) => {
        group.items.forEach((item) => {
          if (item.toggle) payload[item.key] = item.on
        })
      })
      uni.setStorageSync(STORAGE_KEY, payload)
    },
    handleRow(item) {
      if (item.toggle) {
        item.on = !item.on
        this.persistPreferences()
        return
      }

      if (item.key === 'policy') {
        uni.navigateTo({
          url: '/pages/visago/profile/settings/about/privacy-policy/index',
          fail: () => {
            uni.showToast({ title: '打开失败', icon: 'none' })
          },
        })
        return
      }
      if (item.key === 'delete') {
        uni.navigateTo({
          url: '/pages/visago/profile/settings/privacy/delete-account/index',
          fail: () => {
            uni.showToast({ title: '打开失败', icon: 'none' })
          },
        })
        return
      }

      uni.showToast({
        title: `${item.title}功能演示`,
        icon: 'none',
      })
    },
  },
}
</script>

<style scoped>
.settings-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.settings-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.settings-content {
  padding: 18px 16px 34px;
  box-sizing: border-box;
}

.group + .group {
  margin-top: 18px;
}

.group-title,
.row-title,
.row-value {
  display: block;
}

.group-title {
  margin: 0 0 8px 4px;
  font-size: 12px;
  font-weight: 900;
  color: var(--visago-text-soft);
  letter-spacing: 0.12em;
}

.group-card {
  overflow: hidden;
  border-radius: 18px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
}

.row {
  position: relative;
  min-height: 62px;
  padding: 0 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.row-left,
.row-right {
  display: flex;
  align-items: center;
}

.row-left {
  gap: 12px;
  min-width: 0;
  flex: 1;
}

.row-right {
  gap: 6px;
  flex-shrink: 0;
}

.row-icon {
  width: 38px;
  height: 38px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.row-icon .material-symbols-outlined {
  font-size: 20px;
}

.row-title {
  font-size: 15px;
  color: var(--visago-text);
  font-weight: 700;
}

.row-title--danger {
  color: #ef4444;
}

.row-value {
  font-size: 12px;
  font-weight: 700;
  color: var(--visago-text-soft);
}

.row-arrow {
  color: var(--visago-text-soft);
}

.divider {
  position: absolute;
  left: 64px;
  right: 0;
  bottom: 0;
  height: 1px;
  background: var(--visago-line);
}

.mini-switch {
  width: 42px;
  height: 24px;
  border-radius: 999px;
  padding: 2px;
  background: var(--visago-surface-soft);
  box-sizing: border-box;
  display: flex;
}

.mini-switch--on {
  justify-content: flex-end;
  background: var(--visago-primary);
}

.mini-thumb {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #fff;
}
</style>
