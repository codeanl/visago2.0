<template>
  <view class="nav-shell">
    <view class="visago-page-width nav-inner">
      <view v-for="item in navItems" :key="item.key" class="nav-slot" :class="{ 'nav-slot--plan': item.key === 'plan' }" @tap="goTab(item.key)">
        <template v-if="item.key === 'plan'">
          <view class="plan-bubble" :class="{ 'plan-bubble--active': activeTab === 'plan' }">
            <text class="material-symbols-outlined plan-icon">task_alt</text>
          </view>
        </template>
        <template v-else>
          <text class="material-symbols-outlined nav-icon" :class="{ 'nav-icon--active': activeTab === item.key }">{{ item.icon }}</text>
          <text class="nav-label" :class="{ 'nav-label--active': activeTab === item.key }">{{ item.label }}</text>
        </template>
      </view>
    </view>
  </view>
</template>

<script>
const TAB_ROUTES = {
  home: '/pages/visago/home/index',
  visa: '/pages/visago/visa/index',
  plan: '/pages/visago/plan/index',
  community: '/pages/visago/community/index',
  profile: '/pages/visago/profile/index',
}

export default {
  name: 'VisagoBottomNav',
  props: {
    activeTab: {
      type: String,
      required: true,
      validator(value) {
        return ['home', 'visa', 'plan', 'community', 'profile'].includes(value)
      },
    },
  },
  data() {
    return {
      navigating: false,
      navItems: [
        { key: 'home', label: '首页', icon: 'home' },
        { key: 'visa', label: '签证', icon: 'description' },
        { key: 'plan', label: '计划', icon: 'task_alt' },
        { key: 'community', label: '社区', icon: 'group' },
        { key: 'profile', label: '我的', icon: 'person' },
      ],
    }
  },
  methods: {
    goTab(key) {
      if (this.navigating || key === this.activeTab) {
        return
      }
      const url = TAB_ROUTES[key]
      if (!url) {
        return
      }

      this.navigating = true
      uni.reLaunch({
        url,
        complete: () => {
          setTimeout(() => {
            this.navigating = false
          }, 240)
        },
      })
    },
  },
}
</script>

<style scoped>
.nav-shell {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  background: var(--visago-surface);
  border-top: 1px solid var(--visago-line);
  box-shadow: 0 -6px 24px rgba(23, 30, 48, 0.06);
}

.nav-inner {
  height: calc(82px + var(--visago-safe-bottom));
  padding: 2px 8px calc(var(--visago-safe-bottom) + 2px);
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-slot {
  flex: 1;
  min-height: 56px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.nav-slot--plan {
  min-height: 62px;
}

.nav-icon {
  font-size: 27px;
  color: #8a95ad;
}

.nav-icon--active {
  color: var(--visago-primary);
  font-variation-settings: 'FILL' 1;
}

.nav-label {
  font-size: 12px;
  color: #8a95ad;
  line-height: 1;
}

.nav-label--active {
  color: var(--visago-primary);
}

.plan-bubble {
  width: 54px;
  height: 54px;
  border-radius: 9999px;
  margin-bottom: 0;
  background: var(--visago-primary);
  border: 3px solid var(--visago-surface);
  box-shadow: var(--visago-shadow-fab);
  display: flex;
  align-items: center;
  justify-content: center;
}

.plan-bubble--active {
  background: linear-gradient(165deg, #258cff 0%, #0f65d8 72%);
}

.plan-icon {
  color: #fff;
  font-size: 28px;
  font-variation-settings: 'FILL' 1;
}

</style>

