<template>
  <view class="v-tabbar-shell">
    <view class="v-page-width v-tabbar-wrap">
      <view
        v-for="item in items"
        :key="item.key"
        class="v-tabbar-item"
        :class="{ 'v-tabbar-item--active': active === item.key }"
        @tap="switchTo(item)"
      >
        <VIcon class="v-tabbar-icon" :name="item.icon" />
        <text class="v-tabbar-label">{{ item.label }}</text>
      </view>
    </view>
  </view>
</template>

<script>
import VIcon from './VIcon.vue'
import { MAIN_TABS } from '../../config/tabs'

export default {
  name: 'VTabbar',
  components: { VIcon },
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
.v-tabbar-shell {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 60;
  padding: 12px 16px calc(12px + var(--vant-safe-bottom));
  pointer-events: none;
}

.v-tabbar-wrap {
  display: flex;
  padding: 10px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.98);
  border: 1px solid rgba(20, 37, 63, 0.08);
  box-shadow: 0 -6px 24px rgba(23, 30, 48, 0.08);
  pointer-events: auto;
}

.v-tabbar-item {
  flex: 1;
  height: 56px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #7b8797;
}

.v-tabbar-item + .v-tabbar-item {
  margin-left: 8px;
}

.v-tabbar-item--active {
  color: #fff;
  background: #1677ff;
}

.v-tabbar-icon {
  font-size: 16px;
  margin-right: 8px;
}

.v-tabbar-label {
  font-size: 14px;
  font-weight: 800;
}
</style>
