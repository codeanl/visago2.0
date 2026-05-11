<template>
  <view class="vu-tabbar-shell" :style="shellStyle">
    <view class="page-width vu-tabbar-card">
      <view
        v-for="item in items"
        :key="item.key"
        class="vu-tabbar-item"
        :class="{ 'vu-tabbar-item--active': active === item.key }"
        hover-class="vu-tabbar-item--hover"
        @tap="switchTo(item)"
      >
        <u-icon
          :name="active === item.key ? item.activeIcon : item.icon"
          :size="18"
          :color="active === item.key ? '#ffffff' : '#7b8495'"
        />
        <text class="vu-tabbar-label">{{ item.label }}</text>
      </view>
    </view>
  </view>
</template>

<script>
import { MAIN_TABS } from '../../config/tabs'
import { getLayoutMetrics } from '../../utils/layout'

export default {
  name: 'VuTabbar',
  props: {
    active: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      metrics: getLayoutMetrics(),
    }
  },
  computed: {
    items() {
      return MAIN_TABS
    },
    shellStyle() {
      return `padding-bottom:${this.metrics.safeBottom + 10}px;`
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
.vu-tabbar-shell {
  position: fixed;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 80;
  padding-top: 10px;
  padding-left: 16px;
  padding-right: 16px;
  background: rgba(247, 248, 251, 0.96);
}

.vu-tabbar-card {
  height: 68px;
  padding: 8px;
  border-radius: 22px;
  background: #ffffff;
  border: 1px solid #e8edf5;
  box-shadow: 0 -4px 18px rgba(31, 35, 41, 0.07);
  display: flex;
}

.vu-tabbar-item {
  flex: 1;
  height: 52px;
  border-radius: 16px;
  color: #7b8495;
  display: flex;
  align-items: center;
  justify-content: center;
}

.vu-tabbar-item + .vu-tabbar-item {
  margin-left: 8px;
}

.vu-tabbar-item--active {
  color: #ffffff;
  background: #1677ff;
}

.vu-tabbar-item--hover {
  background: #eef5ff;
}

.vu-tabbar-item--active.vu-tabbar-item--hover {
  background: #1677ff;
}

.vu-tabbar-label {
  margin-left: 8px;
  font-size: 14px;
  line-height: 1;
  font-weight: 800;
}
</style>
