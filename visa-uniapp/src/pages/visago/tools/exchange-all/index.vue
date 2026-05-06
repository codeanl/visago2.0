<template>
  <view class="exchange-all-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="全部汇率" />

    <view class="all-content visago-page-width">
      <view class="search-box">
        <text class="material-symbols-outlined search-icon">search</text>
        <input v-model.trim="keyword" class="search-input" placeholder="搜索币种或国家..." />
      </view>

      <view class="filter-row">
        <view
          v-for="item in filters"
          :key="item"
          class="filter-chip"
          :class="{ 'filter-chip--active': activeFilter === item }"
          @tap="activeFilter = item"
        >
          {{ item }}
        </view>
        <view class="sort-btn">
          <text class="material-symbols-outlined">sort</text>
          <text>Sort</text>
        </view>
      </view>

      <view class="base-card visago-card">
        <view>
          <text class="base-label">基准货币</text>
          <view class="base-main">
            <text class="flag">🇨🇳</text>
            <text class="base-value">1 CNY</text>
          </view>
        </view>
        <view class="base-right">
          <text class="base-label">最后更新</text>
          <text class="base-time">刚刚</text>
        </view>
      </view>

      <view class="rate-list">
        <view v-for="item in filteredRates" :key="item.code" class="rate-card visago-card">
          <view class="rate-left">
            <view class="rate-avatar">{{ item.flag }}</view>
            <view>
              <text class="rate-code">{{ item.code }}</text>
              <text class="rate-name">{{ item.name }}</text>
            </view>
          </view>
          <view class="rate-right">
            <text class="rate-num">{{ item.value }}</text>
            <text class="rate-change" :class="item.changeClass">{{ item.changeText }}</text>
          </view>
          <text class="material-symbols-outlined star">star</text>
        </view>
      </view>

      <view class="tips">汇率仅供参考。交易可能受银行手续费及市场波动影响</view>
    </view>

    <view class="float-btn">
      <text class="material-symbols-outlined">currency_exchange</text>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'

export default {
  components: { VisagoTopBar },
  data() {
    return {
      keyword: '',
      activeFilter: '全部',
      filters: ['全部', '前10', '收藏'],
      rates: [
        { code: 'USD', name: '美元', value: '0.1384', change: 0.12, flag: '🇺🇸' },
        { code: 'EUR', name: '欧元', value: '0.1272', change: -0.04, flag: '🇪🇺' },
        { code: 'JPY', name: '日元', value: '20.84', change: 0.85, flag: '🇯🇵' },
        { code: 'GBP', name: '英镑', value: '0.1082', change: 0.0, flag: '🇬🇧' },
        { code: 'AUD', name: '澳元', value: '0.2075', change: -0.21, flag: '🇦🇺' },
        { code: 'HKD', name: '港币', value: '1.082', change: 0.05, flag: '🇭🇰' },
        { code: 'CAD', name: '加元', value: '0.1891', change: 0.11, flag: '🇨🇦' },
      ],
    }
  },
  computed: {
    filteredRates() {
      const keyword = this.keyword.trim().toLowerCase()
      let list = this.rates
      if (this.activeFilter === '前10') {
        list = this.rates.slice(0, 5)
      }
      if (this.activeFilter === '收藏') {
        list = this.rates.filter((item) => ['USD', 'EUR', 'JPY'].includes(item.code))
      }
      if (!keyword) {
        return list.map(this.formatRateItem)
      }
      return list
        .filter((item) => item.code.toLowerCase().includes(keyword) || item.name.includes(keyword))
        .map(this.formatRateItem)
    },
  },
  onShow() {
    applyTheme(getStoredTheme())
  },
  methods: {
    formatRateItem(item) {
      let changeClass = 'rate-change--flat'
      let sign = ''
      if (item.change > 0) {
        changeClass = 'rate-change--up'
        sign = '+'
      } else if (item.change < 0) {
        changeClass = 'rate-change--down'
        sign = '-'
      }
      return {
        ...item,
        changeClass,
        changeText: `${sign}${Math.abs(item.change).toFixed(2)}%`,
      }
    },
  },
}
</script>

<style scoped>
.exchange-all-page {
  min-height: 100vh;
  background: var(--visago-bg);
}

.all-content {
  box-sizing: border-box;
  padding: 92px 16px 96px;
}

.search-box {
  height: 48px;
  border-radius: 12px;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  padding: 0 14px;
  gap: 8px;
}

.search-icon {
  color: var(--visago-text-soft);
}

.search-input {
  flex: 1;
  height: 48px;
  font-size: 15px;
}

.filter-row {
  margin-top: 12px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.filter-chip {
  padding: 6px 16px;
  border-radius: 9999px;
  background: var(--visago-surface-soft);
  color: var(--visago-text-muted);
  font-size: 13px;
}

.filter-chip--active {
  background: var(--visago-primary);
  color: #fff;
}

.sort-btn {
  margin-left: auto;
  color: var(--visago-primary);
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
}

.base-card {
  margin-top: 12px;
  padding: 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-radius: 16px;
}

.base-label {
  display: block;
  color: var(--visago-text-soft);
  font-size: 13px;
  font-weight: 600;
}

.base-main {
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.base-value {
  font-size: 22px;
  font-weight: 600;
}

.base-right {
  text-align: right;
}

.base-time {
  margin-top: 8px;
  display: block;
  font-size: 13px;
}

.rate-list {
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rate-card {
  position: relative;
  padding: 12px 14px;
  border-radius: 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.rate-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.rate-avatar {
  width: 48px;
  height: 48px;
  border-radius: 9999px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26px;
  background: var(--visago-surface-soft);
}

.rate-code {
  display: block;
  font-size: 17px;
  font-weight: 600;
}

.rate-name {
  margin-top: 2px;
  display: block;
  color: var(--visago-text-muted);
  font-size: 13px;
}

.rate-right {
  margin-right: 18px;
  text-align: right;
}

.rate-num {
  display: block;
  font-size: 17px;
  font-weight: 700;
}

.rate-change {
  margin-top: 2px;
  display: block;
  font-size: 12px;
}

.rate-change--up {
  color: #15803d;
}

.rate-change--down {
  color: #dc2626;
}

.rate-change--flat {
  color: var(--visago-text-muted);
}

.star {
  position: absolute;
  top: 10px;
  right: 10px;
  color: #c5cde1;
}

.tips {
  margin-top: 22px;
  color: var(--visago-text-soft);
  text-align: center;
  font-size: 11px;
}

.float-btn {
  position: fixed;
  right: 20px;
  bottom: 18px;
  width: 56px;
  height: 56px;
  border-radius: 9999px;
  background: var(--visago-primary);
  box-shadow: var(--visago-shadow-fab);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.float-btn .material-symbols-outlined {
  font-size: 28px;
}
</style>

