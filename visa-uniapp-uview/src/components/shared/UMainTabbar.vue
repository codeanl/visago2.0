<template>
  <view class="u-tabbar-shell">
    <view class="app-page-width u-tabbar-wrap">
      <view
        v-for="item in items"
        :key="item.key"
        class="u-tabbar-item"
        :class="{ 'u-tabbar-item--active': active === item.key }"
        @tap="switchTo(item)"
      >
        <u-icon
          class="u-tabbar-icon"
          :name="active === item.key ? item.activeIcon : item.icon"
          :size="18"
          :color="active === item.key ? '#ffffff' : '#7b8797'"
        />
        <text class="u-tabbar-label">{{ item.label }}</text>
      </view>
    </view>
  </view>
</template>

<script>
import { MAIN_TABS } from '../../config/tabs'

export default {
  name: 'UMainTabbar',
  props: {
    active: {
      type: String,
      default: '',
    },
  },
  computed: {
    items() {
      return MAIN_TABS
    },
  },
  methods: {
    switchTo(item) {
      if (!item || !item.route || item.key === this.active) return
      uni.reLaunch({ url: item.route })
    },
  },
}
</script>

<style scoped>
.u-tabbar-shell {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 60;
  padding: 10px 16px calc(10px + var(--uview-safe-bottom));
}

.u-tabbar-wrap {
  display: flex;
  padding: 8px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.98);
  border: 1px solid rgba(20, 37, 63, 0.06);
  box-shadow: 0 -4px 16px rgba(23, 30, 48, 0.06);
}

.u-tabbar-item {
  flex: 1;
  height: 52px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #7b8797;
}

.u-tabbar-item + .u-tabbar-item {
  margin-left: 8px;
}

.u-tabbar-item--active {
  color: #fff;
  background: var(--uview-brand);
}

.u-tabbar-icon {
  margin-right: 8px;
}

.u-tabbar-label {
  font-size: 14px;
  font-weight: 800;
}
</style>
