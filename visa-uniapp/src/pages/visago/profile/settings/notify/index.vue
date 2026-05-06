<template>
  <view class="settings-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="通知设置" />

    <scroll-view scroll-y class="settings-scroll">
      <view class="settings-content visago-page-width">
        <view v-for="group in groups" :key="group.key" class="group">
          <text class="group-title">{{ group.title }}</text>
          <view class="group-card">
            <view v-for="(item, idx) in group.items" :key="item.key" class="row" @tap="toggleItem(item)">
              <view class="row-left">
                <view class="row-icon" :style="{ color: item.color, background: item.tint }">
                  <text class="material-symbols-outlined">{{ item.icon }}</text>
                </view>
                <view class="row-copy">
                  <text class="row-title">{{ item.title }}</text>
                  <text class="row-desc">{{ item.desc }}</text>
                </view>
              </view>

              <view class="switch" :class="{ 'switch--on': item.on }">
                <view class="switch-thumb" />
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

const STORAGE_KEY = 'visago_notify_preferences'

function buildGroups() {
  return [
    {
      key: 'visa',
      title: '签证相关',
      items: [
        {
          key: 'status',
          title: '状态变更',
          desc: '受理、审核、出签等关键节点提醒。',
          icon: 'task_alt',
          on: true,
          color: '#0f65d8',
          tint: 'rgba(15,101,216,0.12)',
        },
        {
          key: 'material',
          title: '材料提醒',
          desc: '缺件、过期、重新上传等提醒。',
          icon: 'inventory_2',
          on: true,
          color: '#f59e0b',
          tint: 'rgba(245,158,11,0.13)',
        },
        {
          key: 'appointment',
          title: '预约提醒',
          desc: '递签、录指纹、取件前提醒。',
          icon: 'event_available',
          on: true,
          color: '#22c55e',
          tint: 'rgba(34,197,94,0.12)',
        },
      ],
    },
    {
      key: 'account',
      title: '账户与社区',
      items: [
        {
          key: 'account',
          title: '账户安全',
          desc: '异地登录、密码变更通知。',
          icon: 'shield',
          on: true,
          color: '#8b5cf6',
          tint: 'rgba(139,92,246,0.12)',
        },
        {
          key: 'community',
          title: '社区互动',
          desc: '点赞、评论、收藏提醒。',
          icon: 'forum',
          on: false,
          color: '#ec4899',
          tint: 'rgba(236,72,153,0.12)',
        },
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
          payload[item.key] = item.on
        })
      })
      uni.setStorageSync(STORAGE_KEY, payload)
    },
    toggleItem(item) {
      item.on = !item.on
      this.persistPreferences()
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
.row-desc {
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
  min-height: 68px;
  padding: 0 13px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.row-left {
  min-width: 0;
  flex: 1;
  display: flex;
  align-items: center;
  gap: 11px;
}

.row-copy {
  min-width: 0;
  flex: 1;
}

.row-icon {
  width: 36px;
  height: 36px;
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
  font-weight: 800;
  color: var(--visago-text);
}

.row-desc {
  margin-top: 3px;
  font-size: 11px;
  line-height: 1.45;
  color: var(--visago-text-muted);
}

.switch {
  width: 44px;
  height: 26px;
  border-radius: 999px;
  padding: 2px;
  background: var(--visago-surface-soft);
  box-sizing: border-box;
  display: flex;
  flex-shrink: 0;
}

.switch--on {
  justify-content: flex-end;
  background: var(--visago-primary);
}

.switch-thumb {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #fff;
}

.divider {
  position: absolute;
  left: 60px;
  right: 0;
  bottom: 0;
  height: 1px;
  background: var(--visago-line);
}
</style>
