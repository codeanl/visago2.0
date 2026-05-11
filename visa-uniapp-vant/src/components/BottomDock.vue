<template>
  <view class="dock-shell">
    <view class="page-wrap dock-wrap card">
      <view class="dock-item" :class="{ 'dock-item--active': active === 'home' }" @tap="switchTo('/pages/home/index')">
        <text class="dock-emoji">◧</text>
        <text class="dock-label">首页</text>
      </view>
      <view class="dock-item" :class="{ 'dock-item--active': active === 'goals' }" @tap="switchTo('/pages/goals/index')">
        <text class="dock-emoji">◎</text>
        <text class="dock-label">目标</text>
      </view>
    </view>
  </view>
</template>

<script>
export default {
  props: {
    active: {
      type: String,
      default: 'home',
    },
  },
  methods: {
    switchTo(url) {
      const pages = typeof getCurrentPages === 'function' ? getCurrentPages() : []
      const current = pages.length ? `/${pages[pages.length - 1].route}` : ''
      if (current === url) return
      uni.reLaunch({ url })
    },
  },
}
</script>

<style scoped>
.dock-shell {
  position: fixed;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 40;
  padding: 12px 16px calc(12px + var(--lite-safe-bottom));
}

.dock-wrap {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
  padding: 10px;
  border-radius: 22px;
}

.dock-item {
  height: 56px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: var(--lite-text-muted);
}

.dock-item--active {
  color: #fff;
  background: linear-gradient(135deg, var(--lite-primary) 0%, var(--lite-primary-strong) 100%);
}

.dock-emoji {
  font-size: 16px;
}

.dock-label {
  font-size: 14px;
  font-weight: 800;
}
</style>
