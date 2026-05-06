<template>
  <view class="exchange-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="汇率换算" />

    <view class="exchange-content visago-page-width">
      <view class="converter-card visago-card">
        <view v-if="countriesLoading" class="loading-block">
          <text class="loading-title">正在加载国家和币种数据...</text>
          <text class="loading-desc">首次进入会拉取完整国家列表，请稍等一下。</text>
        </view>

        <view v-else-if="countriesError" class="loading-block">
          <text class="loading-title">国家数据加载失败</text>
          <text class="loading-desc">{{ countriesError }}</text>
          <view class="retry-btn" @tap="loadCountries">重新加载</view>
        </view>

        <template v-else>
          <view class="field-block">
            <text class="field-label">FROM</text>
            <view class="field-row">
              <view class="field-currency" @tap="openCountryPicker('from')">
                <text class="flag">{{ countryFlag(fromCountryCode) }}</text>
                <text class="code">{{ fromCountry ? fromCountry.currencyCode : '---' }}</text>
                <text class="material-symbols-outlined drop-icon">expand_more</text>
              </view>
              <input v-model="fromValue" type="digit" class="field-input" placeholder="0" />
            </view>
            <text class="field-caption">{{ fromCountry ? fromCountry.countryName : '选择国家' }}</text>
          </view>

          <view class="swap-wrap">
            <view class="swap-btn" @tap="swapCountries">
              <text class="material-symbols-outlined">swap_vert</text>
            </view>
          </view>

          <view class="field-block field-block--to">
            <text class="field-label">TO</text>
            <view class="field-row">
              <view class="field-currency" @tap="openCountryPicker('to')">
                <text class="flag">{{ countryFlag(toCountryCode) }}</text>
                <text class="code">{{ toCountry ? toCountry.currencyCode : '---' }}</text>
                <text class="material-symbols-outlined drop-icon">expand_more</text>
              </view>
              <text class="field-output">{{ converted }}</text>
            </view>
            <text class="field-caption field-caption--align-right">{{ toCountry ? toCountry.countryName : '选择国家' }}</text>
          </view>

          <view class="rate-row">
            <view class="rate-left">
              <text class="material-symbols-outlined info-icon">info</text>
              <text>{{ rateSummary }}</text>
            </view>
            <text class="rate-update">{{ rateUpdateText }}</text>
          </view>
          <text v-if="quoteError" class="rate-error">{{ quoteError }}</text>
        </template>
      </view>

      <view class="section-head">
        <text class="section-title">市场行情</text>
        <view class="all-link" @tap="openCountryPicker('to')">全部国家</view>
      </view>

      <view class="market-grid">
        <view class="market-card market-card--clickable visago-card" @tap="openTrendDrawer('primary')">
          <view class="market-top">
            <text class="pair">{{ primaryPairLabel }}</text>
            <text :class="primaryPairStatusClass">{{ primaryPairStatus }}</text>
          </view>
          <text class="market-value">{{ primaryPairValue }}</text>
          <view class="market-graph" :class="primaryGraphClass">
            <canvas canvas-id="exchange-trend-primary" class="market-graph-canvas" width="320" height="88" />
            <view v-if="trendLoading" class="market-graph__overlay">加载中</view>
            <view v-else-if="!primaryTrendPoints.length" class="market-graph__overlay">暂无趋势</view>
          </view>
          <text class="market-more">点击查看详情</text>
        </view>

        <view class="market-card market-card--clickable visago-card" @tap="openTrendDrawer('reverse')">
          <view class="market-top">
            <text class="pair">{{ reversePairLabel }}</text>
            <text :class="reversePairStatusClass">{{ reversePairStatus }}</text>
          </view>
          <text class="market-value">{{ reversePairValue }}</text>
          <view class="market-graph" :class="reverseGraphClass">
            <canvas canvas-id="exchange-trend-reverse" class="market-graph-canvas" width="320" height="88" />
            <view v-if="trendLoading" class="market-graph__overlay">加载中</view>
            <view v-else-if="!reverseTrendPoints.length" class="market-graph__overlay">暂无趋势</view>
          </view>
          <text class="market-more">点击查看详情</text>
        </view>
      </view>

      <view class="section-head section-head-gap">
        <text class="section-title">最近记录</text>
      </view>

      <view class="history-card visago-card">
        <template v-if="recentPairs.length">
          <view v-for="(item, idx) in recentPairs" :key="`${item.fromCountryCode}-${item.toCountryCode}-${idx}`" class="history-item" @tap="applyRecentPair(item)">
            <view class="history-left">
              <view class="history-icon">
                <text class="material-symbols-outlined">history</text>
              </view>
              <view>
                <text class="history-title">{{ item.fromCurrency }} to {{ item.toCurrency }}</text>
                <text class="history-time">{{ item.fromCountryName }} → {{ item.toCountryName }} · {{ formatHistoryTime(item.savedAt) }}</text>
              </view>
            </view>
            <view class="history-right">
              <text class="history-amount">{{ item.convertedAmountText }} {{ item.toCurrency }}</text>
              <text class="history-rate">汇率: {{ item.rateText }}</text>
            </view>
            <view v-if="idx < recentPairs.length - 1" class="history-divider" />
          </view>
        </template>
        <view v-else class="history-empty">
          <text class="history-empty__title">暂无换算记录</text>
          <text class="history-empty__desc">选择国家并完成一次换算后，会在这里保留最近记录。</text>
        </view>
      </view>
    </view>

    <view v-if="pickerVisible" class="sheet-mask" @tap="closeCountryPicker">
      <view class="sheet-panel" @tap.stop>
        <view class="sheet-head">
          <text class="sheet-title">{{ pickerMode === 'from' ? '选择付款国家' : '选择到账国家' }}</text>
          <view class="sheet-close" @tap="closeCountryPicker">
            <text class="material-symbols-outlined">close</text>
          </view>
        </view>

        <view class="search-box">
          <text class="material-symbols-outlined search-icon">search</text>
          <input
            v-model.trim="pickerKeyword"
            class="search-input"
            placeholder="搜索国家、英文名或币种代码"
            confirm-type="search"
          />
        </view>

        <scroll-view scroll-y class="sheet-list" show-scrollbar="false">
          <view
            v-for="item in filteredCountries"
            :key="item.countryCode"
            class="sheet-item"
            :class="{ 'sheet-item--active': isSelectedCountry(item) }"
            @tap="selectCountry(item)"
          >
            <view class="sheet-item__main">
              <text class="flag flag--sheet">{{ countryFlag(item.countryCode) }}</text>
              <view class="sheet-item__copy">
                <text class="sheet-item__name">{{ item.countryName }}</text>
                <text class="sheet-item__meta">{{ item.englishName }} · {{ item.currencyCode }} · {{ item.currencyName }}</text>
              </view>
            </view>
            <view class="sheet-item__side">
              <text v-if="!item.supported" class="unsupported-tag">暂不支持</text>
              <text v-else-if="isSelectedCountry(item)" class="material-symbols-outlined selected-icon">check</text>
            </view>
          </view>

          <view v-if="!filteredCountries.length" class="sheet-empty">
            <text class="sheet-empty__title">没有找到匹配的国家</text>
            <text class="sheet-empty__desc">试试输入中文名、英文名或币种代码，例如 USD、JPY。</text>
          </view>
        </scroll-view>
      </view>
    </view>

    <view v-if="trendDrawerVisible" class="trend-mask" @tap="closeTrendDrawer">
      <view class="trend-panel" @tap.stop>
        <view class="trend-panel__head">
          <view>
            <text class="trend-panel__title">{{ activeTrendLabel }}</text>
            <text class="trend-panel__subtitle">{{ activeTrendCountryLabel }}</text>
          </view>
          <view class="sheet-close" @tap="closeTrendDrawer">
            <text class="material-symbols-outlined">close</text>
          </view>
        </view>

        <view class="trend-panel__summary">
          <view>
            <text class="trend-panel__value">{{ activeTrendValue }}</text>
            <text :class="activeTrendStatusClass">{{ activeTrendStatus }}</text>
          </view>
          <text class="trend-panel__date">{{ quoteDate || '最近 7 日' }}</text>
        </view>

        <view class="trend-panel__chart" :class="activeTrendGraphClass">
          <canvas canvas-id="exchange-trend-detail" class="trend-panel__canvas" width="720" height="240" />
          <view v-if="trendLoading" class="market-graph__overlay">加载中</view>
          <view v-else-if="!activeTrendPoints.length" class="market-graph__overlay">暂无趋势</view>
        </view>

        <view class="trend-stats">
          <view class="trend-stat">
            <text class="trend-stat__label">最新</text>
            <text class="trend-stat__value">{{ activeTrendLatest }}</text>
          </view>
          <view class="trend-stat">
            <text class="trend-stat__label">最高</text>
            <text class="trend-stat__value">{{ activeTrendHigh }}</text>
          </view>
          <view class="trend-stat">
            <text class="trend-stat__label">最低</text>
            <text class="trend-stat__value">{{ activeTrendLow }}</text>
          </view>
        </view>

        <scroll-view scroll-y class="trend-list" show-scrollbar="false">
          <view v-for="(item, idx) in activeTrendPoints" :key="`${activeTrendKey}-${item.date}`" class="trend-list__item">
            <text class="trend-list__date">{{ item.date }}</text>
            <text class="trend-list__rate">{{ formatRate(item.rate) }}</text>
            <view v-if="idx < activeTrendPoints.length - 1" class="trend-list__divider" />
          </view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { api } from '../../../../utils/api'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'

const RECENT_PAIR_KEY = 'visago_exchange_recent_pairs'

export default {
  components: { VisagoTopBar },
  data() {
    return {
      countries: [],
      countriesLoading: false,
      countriesError: '',
      fromCountryCode: 'CN',
      toCountryCode: 'US',
      fromValue: '1000',
      quoteLoading: false,
      quoteError: '',
      quoteDate: '',
      rate: 0,
      pickerVisible: false,
      pickerMode: 'from',
      pickerKeyword: '',
      recentPairs: [],
      latestQuoteToken: '',
      trendLoading: false,
      primaryTrendPoints: [],
      reverseTrendPoints: [],
      trendDrawerVisible: false,
      activeTrendKey: 'primary',
    }
  },
  computed: {
    countryMap() {
      const map = {}
      this.countries.forEach((item) => {
        map[item.countryCode] = item
      })
      return map
    },
    fromCountry() {
      return this.countryMap[this.fromCountryCode] || null
    },
    toCountry() {
      return this.countryMap[this.toCountryCode] || null
    },
    filteredCountries() {
      const keyword = String(this.pickerKeyword || '').trim().toLowerCase()
      if (!keyword) {
        return this.countries
      }
      return this.countries.filter((item) => {
        const values = [item.countryName, item.englishName, item.currencyCode, item.currencyName]
        return values.some((value) => String(value || '').toLowerCase().includes(keyword))
      })
    },
    amountNumber() {
      const raw = Number(this.fromValue)
      return Number.isFinite(raw) ? raw : 0
    },
    converted() {
      if (this.quoteLoading || this.quoteError || !this.rate) {
        return '--'
      }
      return this.formatNumber(this.amountNumber * this.rate)
    },
    rateSummary() {
      if (this.quoteLoading) {
        return '正在获取实时汇率...'
      }
      if (this.quoteError) {
        return '当前汇率暂不可用'
      }
      if (!this.fromCountry || !this.toCountry || !this.rate) {
        return '选择国家后自动获取最新参考汇率'
      }
      return `1 ${this.fromCountry.currencyCode} ≈ ${this.formatRate(this.rate)} ${this.toCountry.currencyCode}`
    },
    rateUpdateText() {
      if (this.quoteLoading) {
        return '更新中'
      }
      if (this.quoteDate) {
        return this.quoteDate
      }
      return '等待查询'
    },
    primaryPairLabel() {
      return `${this.fromCountry ? this.fromCountry.currencyCode : '--'}/${this.toCountry ? this.toCountry.currencyCode : '--'}`
    },
    primaryPairStatus() {
      if (this.quoteError) {
        return '异常'
      }
      if (this.primaryTrendPoints.length < 2) {
        return '实时'
      }
      return this.formatTrendStatus(this.primaryTrendDelta)
    },
    primaryPairStatusClass() {
      if (this.primaryTrendPoints.length < 2) {
        return 'up'
      }
      return this.primaryTrendDelta >= 0 ? 'up' : 'down'
    },
    primaryPairValue() {
      return this.rate ? this.formatRate(this.rate) : '--'
    },
    primaryGraphClass() {
      if (this.primaryTrendPoints.length < 2) {
        return 'market-graph-up'
      }
      return this.primaryTrendDelta >= 0 ? 'market-graph-up' : 'market-graph-down'
    },
    primaryTrendDelta() {
      return this.getTrendDelta(this.primaryTrendPoints)
    },
    reversePairLabel() {
      return `${this.toCountry ? this.toCountry.currencyCode : '--'}/${this.fromCountry ? this.fromCountry.currencyCode : '--'}`
    },
    reversePairStatus() {
      if (this.quoteError) {
        return '异常'
      }
      if (this.reverseTrendPoints.length < 2) {
        return '反向'
      }
      return this.formatTrendStatus(this.reverseTrendDelta)
    },
    reversePairStatusClass() {
      if (this.reverseTrendPoints.length < 2) {
        return 'down'
      }
      return this.reverseTrendDelta >= 0 ? 'up' : 'down'
    },
    reversePairValue() {
      return this.rate ? this.formatRate(1 / this.rate) : '--'
    },
    reverseGraphClass() {
      if (this.reverseTrendPoints.length < 2) {
        return 'market-graph-down'
      }
      return this.reverseTrendDelta >= 0 ? 'market-graph-up' : 'market-graph-down'
    },
    reverseTrendDelta() {
      return this.getTrendDelta(this.reverseTrendPoints)
    },
    activeTrendLabel() {
      return this.activeTrendKey === 'reverse' ? this.reversePairLabel : this.primaryPairLabel
    },
    activeTrendCountryLabel() {
      if (this.activeTrendKey === 'reverse') {
        return `${this.toCountry ? this.toCountry.countryName : '--'} → ${this.fromCountry ? this.fromCountry.countryName : '--'}`
      }
      return `${this.fromCountry ? this.fromCountry.countryName : '--'} → ${this.toCountry ? this.toCountry.countryName : '--'}`
    },
    activeTrendPoints() {
      return this.activeTrendKey === 'reverse' ? this.reverseTrendPoints : this.primaryTrendPoints
    },
    activeTrendDelta() {
      return this.activeTrendKey === 'reverse' ? this.reverseTrendDelta : this.primaryTrendDelta
    },
    activeTrendStatus() {
      if (this.quoteError) {
        return '异常'
      }
      if (this.activeTrendPoints.length < 2) {
        return '最近 7 日'
      }
      return this.formatTrendStatus(this.activeTrendDelta)
    },
    activeTrendStatusClass() {
      return this.activeTrendDelta >= 0 ? 'up' : 'down'
    },
    activeTrendGraphClass() {
      return this.activeTrendDelta >= 0 ? 'market-graph-up' : 'market-graph-down'
    },
    activeTrendValue() {
      return this.activeTrendKey === 'reverse' ? this.reversePairValue : this.primaryPairValue
    },
    activeTrendLatest() {
      if (!this.activeTrendPoints.length) {
        return '--'
      }
      return this.formatRate(this.activeTrendPoints[this.activeTrendPoints.length - 1].rate)
    },
    activeTrendHigh() {
      if (!this.activeTrendPoints.length) {
        return '--'
      }
      return this.formatRate(Math.max(...this.activeTrendPoints.map((item) => Number(item.rate || 0))))
    },
    activeTrendLow() {
      if (!this.activeTrendPoints.length) {
        return '--'
      }
      return this.formatRate(Math.min(...this.activeTrendPoints.map((item) => Number(item.rate || 0))))
    },
  },
  onLoad() {
    this.loadCountries()
  },
  onShow() {
    applyTheme(getStoredTheme())
    if (this.primaryTrendPoints.length || this.reverseTrendPoints.length) {
      this.$nextTick(() => {
        this.renderTrendCharts()
      })
    }
  },
  methods: {
    async loadCountries() {
      this.countriesLoading = true
      this.countriesError = ''
      try {
        const items = await api.listExchangeCountries()
        this.countries = Array.isArray(items) ? items : []
        if (!this.countries.length) {
          this.countriesError = '暂时没有获取到国家列表。'
          return
        }
        this.ensureDefaultCountries()
        this.restoreRecentPairs()
        await this.refreshQuote()
      } catch (error) {
        this.countriesError = error && error.message ? error.message : '国家列表加载失败'
      } finally {
        this.countriesLoading = false
      }
    },
    ensureDefaultCountries() {
      if (!this.countryMap[this.fromCountryCode]) {
        this.fromCountryCode = this.pickDefaultCountryCode(['CN', 'US', 'JP'])
      }
      if (!this.countryMap[this.toCountryCode] || this.toCountryCode === this.fromCountryCode) {
        this.toCountryCode = this.pickDefaultCountryCode(['US', 'JP', 'DE'], this.fromCountryCode)
      }
    },
    pickDefaultCountryCode(candidates = [], excludeCode = '') {
      for (const code of candidates) {
        if (this.countryMap[code] && code !== excludeCode) {
          return code
        }
      }
      const fallback = this.countries.find((item) => item.countryCode !== excludeCode)
      return fallback ? fallback.countryCode : ''
    },
    openCountryPicker(mode) {
      this.pickerMode = mode
      this.pickerKeyword = ''
      this.pickerVisible = true
    },
    closeCountryPicker() {
      this.pickerVisible = false
      this.pickerKeyword = ''
    },
    isSelectedCountry(item) {
      if (!item) return false
      return this.pickerMode === 'from' ? item.countryCode === this.fromCountryCode : item.countryCode === this.toCountryCode
    },
    async selectCountry(item) {
      if (!item) return
      if (this.pickerMode === 'from') {
        this.fromCountryCode = item.countryCode
        if (this.fromCountryCode === this.toCountryCode) {
          this.toCountryCode = this.pickDefaultCountryCode(['US', 'JP', 'DE'], this.fromCountryCode)
        }
      } else {
        this.toCountryCode = item.countryCode
        if (this.toCountryCode === this.fromCountryCode) {
          this.fromCountryCode = this.pickDefaultCountryCode(['CN', 'US', 'JP'], this.toCountryCode)
        }
      }
      this.closeCountryPicker()
      await this.refreshQuote()
    },
    async swapCountries() {
      if (!this.fromCountryCode || !this.toCountryCode) return
      const previousFrom = this.fromCountryCode
      this.fromCountryCode = this.toCountryCode
      this.toCountryCode = previousFrom
      await this.refreshQuote()
    },
    async applyRecentPair(item) {
      if (!item) return
      if (!this.countryMap[item.fromCountryCode] || !this.countryMap[item.toCountryCode]) {
        return
      }
      this.fromCountryCode = item.fromCountryCode
      this.toCountryCode = item.toCountryCode
      if (item.amountText) {
        this.fromValue = item.amountText
      }
      await this.refreshQuote()
    },
    openTrendDrawer(key) {
      this.activeTrendKey = key === 'reverse' ? 'reverse' : 'primary'
      this.trendDrawerVisible = true
      this.$nextTick(() => {
        this.renderTrendCharts()
      })
    },
    closeTrendDrawer() {
      this.trendDrawerVisible = false
    },
    async refreshQuote() {
      if (!this.fromCountry || !this.toCountry) {
        return
      }
      const fromCurrency = this.fromCountry.currencyCode
      const toCurrency = this.toCountry.currencyCode
      if (!fromCurrency || !toCurrency) {
        return
      }
      const token = `${fromCurrency}-${toCurrency}-${Date.now()}`
      this.latestQuoteToken = token
      if (!this.fromCountry.supported) {
        this.quoteLoading = false
        this.rate = 0
        this.quoteDate = ''
        this.quoteError = `${this.fromCountry.countryName}使用的 ${fromCurrency} 暂不支持实时汇率查询`
        this.clearTrend()
        return
      }
      if (!this.toCountry.supported) {
        this.quoteLoading = false
        this.rate = 0
        this.quoteDate = ''
        this.quoteError = `${this.toCountry.countryName}使用的 ${toCurrency} 暂不支持实时汇率查询`
        this.clearTrend()
        return
      }
      this.quoteLoading = true
      this.quoteError = ''
      try {
        const result = await api.getExchangeQuote({ from: fromCurrency, to: toCurrency })
        if (this.latestQuoteToken !== token) {
          return
        }
        this.rate = Number(result && result.rate ? result.rate : 0)
        this.quoteDate = result && result.date ? result.date : ''
        if (!this.rate) {
          this.quoteError = '暂未获取到有效汇率'
          this.clearTrend()
          return
        }
        await this.refreshTrend(fromCurrency, toCurrency)
        this.saveRecentPair()
      } catch (error) {
        if (this.latestQuoteToken !== token) {
          return
        }
        this.rate = 0
        this.quoteDate = ''
        this.quoteError = error && error.message ? error.message : '汇率获取失败'
        this.clearTrend()
      } finally {
        if (this.latestQuoteToken === token) {
          this.quoteLoading = false
        }
      }
    },
    async refreshTrend(fromCurrency, toCurrency) {
      this.trendLoading = true
      try {
        const result = await api.getExchangeTrend({ from: fromCurrency, to: toCurrency, days: 7 })
        const primaryPoints = Array.isArray(result && result.points)
          ? result.points
              .map((item) => ({
                date: item.date || '',
                rate: Number(item.rate || 0),
              }))
              .filter((item) => item.date && item.rate > 0)
          : []
        this.primaryTrendPoints = primaryPoints
        this.reverseTrendPoints = primaryPoints.map((item) => ({
          date: item.date,
          rate: item.rate > 0 ? 1 / item.rate : 0,
        }))
        this.$nextTick(() => {
          this.renderTrendCharts()
        })
      } catch (error) {
        this.clearTrend()
      } finally {
        this.trendLoading = false
      }
    },
    clearTrend() {
      this.trendLoading = false
      this.primaryTrendPoints = []
      this.reverseTrendPoints = []
      this.$nextTick(() => {
        this.renderTrendCharts()
      })
    },
    renderTrendCharts() {
      this.drawTrendChart('exchange-trend-primary', this.primaryTrendPoints, {
        line: '#14893a',
        fillTop: 'rgba(86, 157, 108, 0.30)',
        fillBottom: 'rgba(86, 157, 108, 0.06)',
      }, { width: 320, height: 88, lineWidth: 3, pointRadius: 4, paddingX: 10, paddingY: 8 })
      this.drawTrendChart('exchange-trend-reverse', this.reverseTrendPoints, {
        line: '#c45c5c',
        fillTop: 'rgba(196, 92, 92, 0.26)',
        fillBottom: 'rgba(196, 92, 92, 0.06)',
      }, { width: 320, height: 88, lineWidth: 3, pointRadius: 4, paddingX: 10, paddingY: 8 })
      if (this.trendDrawerVisible) {
        this.drawTrendChart('exchange-trend-detail', this.activeTrendPoints, {
          line: this.activeTrendDelta >= 0 ? '#14893a' : '#c45c5c',
          fillTop: this.activeTrendDelta >= 0 ? 'rgba(86, 157, 108, 0.30)' : 'rgba(196, 92, 92, 0.26)',
          fillBottom: this.activeTrendDelta >= 0 ? 'rgba(86, 157, 108, 0.06)' : 'rgba(196, 92, 92, 0.06)',
        }, { width: 720, height: 240, lineWidth: 4, pointRadius: 5, paddingX: 18, paddingY: 16 })
      }
    },
    drawTrendChart(canvasId, points, palette, options = {}) {
      const ctx = uni.createCanvasContext(canvasId, this)
      const width = Number(options.width || 320)
      const height = Number(options.height || 96)
      ctx.clearRect(0, 0, width, height)
      if (!Array.isArray(points) || points.length < 2) {
        ctx.draw()
        return
      }

      const rates = points.map((item) => Number(item.rate || 0))
      let min = Math.min(...rates)
      let max = Math.max(...rates)
      if (max === min) {
        const padding = max === 0 ? 1 : Math.abs(max) * 0.02
        max += padding
        min -= padding
      }

      const paddingX = Number(options.paddingX || 10)
      const paddingY = Number(options.paddingY || 8)
      const chartWidth = width - paddingX * 2
      const chartHeight = height - paddingY * 2
      const stepX = chartWidth / Math.max(points.length - 1, 1)
      const coords = rates.map((rate, index) => ({
        x: paddingX + stepX * index,
        y: height - paddingY - ((rate - min) / (max - min)) * chartHeight,
      }))

      const fillGradient = ctx.createLinearGradient(0, 0, 0, height)
      fillGradient.addColorStop(0, palette.fillTop)
      fillGradient.addColorStop(1, palette.fillBottom)

      ctx.beginPath()
      ctx.moveTo(coords[0].x, height - paddingY)
      coords.forEach((point) => {
        ctx.lineTo(point.x, point.y)
      })
      ctx.lineTo(coords[coords.length - 1].x, height - paddingY)
      ctx.closePath()
      ctx.setFillStyle(fillGradient)
      ctx.fill()

      ctx.beginPath()
      ctx.setStrokeStyle(palette.line)
      ctx.setLineWidth(Number(options.lineWidth || 3))
      ctx.setLineCap('round')
      ctx.setLineJoin('round')
      ctx.moveTo(coords[0].x, coords[0].y)
      coords.slice(1).forEach((point) => {
        ctx.lineTo(point.x, point.y)
      })
      ctx.stroke()

      const lastPoint = coords[coords.length - 1]
      ctx.beginPath()
      ctx.setFillStyle(palette.line)
      ctx.arc(lastPoint.x, lastPoint.y, Number(options.pointRadius || 4), 0, Math.PI * 2)
      ctx.fill()
      ctx.draw()
    },
    saveRecentPair() {
      if (!this.fromCountry || !this.toCountry || !this.rate) {
        return
      }
      const convertedAmount = this.amountNumber * this.rate
      const current = {
        fromCountryCode: this.fromCountry.countryCode,
        fromCountryName: this.fromCountry.countryName,
        fromCurrency: this.fromCountry.currencyCode,
        toCountryCode: this.toCountry.countryCode,
        toCountryName: this.toCountry.countryName,
        toCurrency: this.toCountry.currencyCode,
        amountText: this.fromValue,
        convertedAmountText: this.formatNumber(convertedAmount),
        rateText: this.formatRate(this.rate),
        savedAt: Date.now(),
      }
      const next = [current, ...this.recentPairs.filter((item) => !(item.fromCountryCode === current.fromCountryCode && item.toCountryCode === current.toCountryCode))].slice(0, 6)
      this.recentPairs = next
      uni.setStorageSync(RECENT_PAIR_KEY, next)
    },
    restoreRecentPairs() {
      const saved = uni.getStorageSync(RECENT_PAIR_KEY)
      const list = Array.isArray(saved) ? saved : []
      this.recentPairs = list.filter((item) => this.countryMap[item.fromCountryCode] && this.countryMap[item.toCountryCode]).slice(0, 6)
    },
    countryFlag(code) {
      const value = String(code || '').trim().toUpperCase()
      if (!/^[A-Z]{2}$/.test(value)) {
        return '🏳️'
      }
      return String.fromCodePoint(...value.split('').map((char) => 127397 + char.charCodeAt(0)))
    },
    getTrendDelta(points) {
      if (!Array.isArray(points) || points.length < 2) {
        return 0
      }
      const first = Number(points[0].rate || 0)
      const last = Number(points[points.length - 1].rate || 0)
      if (!first || !last) {
        return 0
      }
      return ((last - first) / first) * 100
    },
    formatTrendStatus(delta) {
      const prefix = delta >= 0 ? '近7日 +' : '近7日 '
      return `${prefix}${this.trimNumber(delta.toFixed(2))}%`
    },
    formatRate(value) {
      if (!Number.isFinite(value)) {
        return '--'
      }
      if (value >= 100) {
        return this.trimNumber(value.toFixed(2))
      }
      if (value >= 1) {
        return this.trimNumber(value.toFixed(4))
      }
      return this.trimNumber(value.toFixed(6))
    },
    formatNumber(value) {
      if (!Number.isFinite(value)) {
        return '--'
      }
      const digits = Math.abs(value) < 1 ? 4 : 2
      return this.addThousands(this.trimNumber(value.toFixed(digits)))
    },
    trimNumber(value) {
      return String(value).replace(/(\.\d*?[1-9])0+$/u, '$1').replace(/\.0+$/u, '')
    },
    addThousands(value) {
      const [integerPart, decimalPart] = String(value).split('.')
      const formattedInteger = integerPart.replace(/\B(?=(\d{3})+(?!\d))/gu, ',')
      return decimalPart ? `${formattedInteger}.${decimalPart}` : formattedInteger
    },
    formatHistoryTime(value) {
      const date = new Date(Number(value || 0))
      if (Number.isNaN(date.getTime())) {
        return '刚刚'
      }
      const year = date.getFullYear()
      const month = `${date.getMonth() + 1}`.padStart(2, '0')
      const day = `${date.getDate()}`.padStart(2, '0')
      const hours = `${date.getHours()}`.padStart(2, '0')
      const minutes = `${date.getMinutes()}`.padStart(2, '0')
      return `${year}-${month}-${day} ${hours}:${minutes}`
    },
  },
}
</script>

<style scoped>
.exchange-page {
  min-height: 100vh;
  background: var(--visago-bg);
}

.exchange-content {
  box-sizing: border-box;
  padding: 92px 16px 32px;
}

.converter-card {
  padding: 16px 14px 14px;
  border-radius: 16px;
}

.loading-block {
  padding: 8px 2px;
}

.loading-title {
  display: block;
  font-size: 16px;
  font-weight: 600;
  color: var(--visago-text);
}

.loading-desc {
  display: block;
  margin-top: 6px;
  color: var(--visago-text-muted);
  font-size: 13px;
  line-height: 1.6;
}

.retry-btn {
  margin-top: 12px;
  width: 96px;
  height: 36px;
  border-radius: 999px;
  background: var(--visago-primary);
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
}

.field-block + .field-block {
  margin-top: 2px;
}

.field-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--visago-text-soft);
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.field-row {
  margin-top: 6px;
  height: 52px;
  border-radius: 12px;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
}

.field-currency {
  display: flex;
  align-items: center;
  gap: 8px;
  border-right: 1px solid #b9c0d0;
  padding-right: 10px;
}

.field-caption {
  display: block;
  margin-top: 6px;
  color: var(--visago-text-soft);
  font-size: 11px;
}

.field-caption--align-right {
  text-align: right;
}

.flag {
  font-size: 19px;
}

.code {
  font-size: 20px;
  font-weight: 600;
}

.drop-icon {
  color: var(--visago-text-soft);
}

.field-input,
.field-output {
  flex: 1;
  text-align: right;
  font-size: 28px;
  font-weight: 600;
}

.field-input {
  color: var(--visago-primary);
}

.field-block--to {
  margin-top: 6px;
}

.swap-wrap {
  display: flex;
  justify-content: center;
  margin: -8px 0 -4px;
  position: relative;
  z-index: 1;
}

.swap-btn {
  width: 40px;
  height: 40px;
  border-radius: 9999px;
  background: var(--visago-primary);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--visago-shadow-fab);
}

.rate-row {
  margin-top: 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.rate-left {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--visago-text-muted);
  font-size: 13px;
}

.info-icon {
  font-size: 16px;
}

.rate-update {
  color: #14893a;
  font-size: 13px;
  font-weight: 600;
}

.rate-error {
  display: block;
  margin-top: 8px;
  color: #dc2626;
  font-size: 12px;
  line-height: 1.6;
}

.section-head {
  margin-top: 22px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-head-gap {
  margin-top: 26px;
}

.section-title {
  font-size: 22px;
  font-weight: 600;
}

.all-link {
  color: var(--visago-primary);
  font-size: 13px;
  font-weight: 500;
}

.market-grid {
  margin-top: 10px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.market-card {
  padding: 12px;
  border-radius: 14px;
}

.market-card--clickable {
  position: relative;
}

.market-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 4px;
}

.pair {
  font-size: 12px;
  font-weight: 500;
}

.up {
  color: #14893a;
  font-size: 11px;
}

.down {
  color: #dc2626;
  font-size: 11px;
}

.market-value {
  margin-top: 8px;
  display: block;
  font-size: 20px;
  font-weight: 600;
}

.market-graph {
  margin-top: 8px;
  height: 88px;
  border-radius: 2px;
  position: relative;
  overflow: hidden;
}

.market-graph-up {
  background: linear-gradient(145deg, rgba(168, 196, 168, 0.9), rgba(200, 224, 200, 0.9));
}

.market-graph-down {
  background: linear-gradient(145deg, rgba(177, 161, 161, 0.9), rgba(214, 165, 165, 0.9));
}

.market-graph-canvas {
  width: 100%;
  height: 88px;
}

.market-graph__overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(17, 24, 39, 0.52);
  font-size: 11px;
  font-weight: 600;
}

.market-more {
  display: block;
  margin-top: 8px;
  color: var(--visago-text-soft);
  font-size: 11px;
  text-align: right;
}

.history-card {
  margin-top: 12px;
  overflow: hidden;
  border-radius: 14px;
}

.history-item {
  position: relative;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  gap: 8px;
}

.history-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.history-icon {
  width: 40px;
  height: 40px;
  border-radius: 9999px;
  background: var(--visago-surface-soft);
  color: var(--visago-primary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.history-title {
  display: block;
  font-size: 13px;
  font-weight: 500;
}

.history-time {
  display: block;
  margin-top: 2px;
  color: var(--visago-text-soft);
  font-size: 11px;
  line-height: 1.5;
}

.history-right {
  text-align: right;
}

.history-amount {
  display: block;
  font-size: 13px;
  font-weight: 600;
}

.history-rate {
  display: block;
  margin-top: 2px;
  color: #14893a;
  font-size: 11px;
}

.history-divider {
  position: absolute;
  left: 64px;
  right: 0;
  bottom: 0;
  height: 1px;
  background: var(--visago-line);
}

.history-empty {
  padding: 24px 18px;
  text-align: center;
}

.history-empty__title {
  display: block;
  font-size: 15px;
  font-weight: 600;
  color: var(--visago-text);
}

.history-empty__desc {
  display: block;
  margin-top: 8px;
  color: var(--visago-text-soft);
  font-size: 12px;
  line-height: 1.6;
}

.sheet-mask {
  position: fixed;
  inset: 0;
  z-index: 80;
  background: rgba(12, 18, 32, 0.36);
  display: flex;
  align-items: flex-end;
}

.trend-mask {
  position: fixed;
  inset: 0;
  z-index: 90;
  background: rgba(12, 18, 32, 0.36);
  display: flex;
  align-items: flex-end;
}

.sheet-panel {
  width: 100%;
  max-height: 82vh;
  border-radius: 24px 24px 0 0;
  background: var(--visago-surface);
  padding: 18px 16px 22px;
  box-sizing: border-box;
}

.trend-panel {
  width: 100%;
  height: 82vh;
  border-radius: 24px 24px 0 0;
  background: var(--visago-surface);
  padding: 18px 16px 22px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.sheet-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.trend-panel__head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  flex-shrink: 0;
}

.sheet-title {
  color: var(--visago-text);
  font-size: 20px;
  font-weight: 700;
}

.trend-panel__title {
  display: block;
  color: var(--visago-text);
  font-size: 22px;
  font-weight: 700;
}

.trend-panel__subtitle {
  display: block;
  margin-top: 6px;
  color: var(--visago-text-muted);
  font-size: 13px;
}

.trend-panel__summary {
  margin-top: 16px;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 12px;
  flex-shrink: 0;
}

.trend-panel__value {
  display: block;
  color: var(--visago-text);
  font-size: 30px;
  font-weight: 700;
}

.trend-panel__date {
  color: var(--visago-text-soft);
  font-size: 12px;
}

.trend-panel__chart {
  margin-top: 16px;
  height: 240px;
  border-radius: 16px;
  position: relative;
  overflow: hidden;
  flex-shrink: 0;
}

.trend-panel__canvas {
  width: 100%;
  height: 240px;
  display: block;
}

.trend-stats {
  margin-top: 16px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  flex-shrink: 0;
}

.trend-stat {
  padding: 12px;
  border-radius: 14px;
  background: var(--visago-surface-soft);
}

.trend-stat__label {
  display: block;
  color: var(--visago-text-soft);
  font-size: 11px;
}

.trend-stat__value {
  display: block;
  margin-top: 6px;
  color: var(--visago-text);
  font-size: 16px;
  font-weight: 700;
}

.trend-list {
  margin-top: 16px;
  flex: 1;
  min-height: 0;
}

.trend-list__item {
  position: relative;
  min-height: 50px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.trend-list__date {
  color: var(--visago-text-muted);
  font-size: 13px;
}

.trend-list__rate {
  color: var(--visago-text);
  font-size: 14px;
  font-weight: 600;
}

.trend-list__divider {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 1px;
  background: var(--visago-line);
}

.sheet-close {
  width: 36px;
  height: 36px;
  border-radius: 999px;
  background: var(--visago-surface-soft);
  color: var(--visago-text-soft);
  display: flex;
  align-items: center;
  justify-content: center;
}

.search-box {
  margin-top: 16px;
  height: 46px;
  border-radius: 16px;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 14px;
}

.search-icon {
  color: var(--visago-text-soft);
  font-size: 20px;
}

.search-input {
  flex: 1;
  min-width: 0;
  color: var(--visago-text);
  font-size: 14px;
}

.sheet-list {
  margin-top: 14px;
  max-height: calc(82vh - 122px);
}

.sheet-item {
  min-height: 72px;
  padding: 14px 2px;
  border-bottom: 1px solid var(--visago-line);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.sheet-item--active {
  background: rgba(99, 102, 241, 0.05);
}

.sheet-item__main {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.flag--sheet {
  font-size: 22px;
}

.sheet-item__copy {
  min-width: 0;
}

.sheet-item__name {
  display: block;
  color: var(--visago-text);
  font-size: 15px;
  font-weight: 700;
}

.sheet-item__meta {
  display: block;
  margin-top: 4px;
  color: var(--visago-text-muted);
  font-size: 12px;
  line-height: 1.6;
}

.sheet-item__side {
  flex-shrink: 0;
}

.unsupported-tag {
  padding: 4px 8px;
  border-radius: 999px;
  background: rgba(245, 158, 11, 0.14);
  color: #d97706;
  font-size: 11px;
  font-weight: 700;
}

.selected-icon {
  color: var(--visago-primary);
  font-size: 20px;
}

.sheet-empty {
  padding: 36px 12px 16px;
  text-align: center;
}

.sheet-empty__title {
  display: block;
  color: var(--visago-text);
  font-size: 16px;
  font-weight: 700;
}

.sheet-empty__desc {
  display: block;
  margin-top: 8px;
  color: var(--visago-text-muted);
  font-size: 13px;
  line-height: 1.7;
}

:global(html.theme-dark) .field-currency {
  border-right-color: #4e576b;
}

:global(html.theme-dark) .market-graph__overlay {
  color: rgba(255, 255, 255, 0.65);
}

:global(html.theme-dark) .trend-mask,
:global(html.theme-dark) .sheet-mask {
  background: rgba(0, 0, 0, 0.54);
}
</style>

