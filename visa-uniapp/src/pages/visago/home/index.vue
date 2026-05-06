<template>
  <view class="home-page">
    <VisagoTopBar :home="true" :show-notice="false" logo-src="/static/header-logo.png" />

    <view class="home-content visago-page-width">
      <view class="hero-section">
        <text class="hero-title">想去哪里探索世界？</text>
        <view class="search-box">
          <text class="material-symbols-outlined search-icon" @tap="searchVisaCountries">search</text>
          <input
            v-model.trim="keyword"
            class="search-input"
            placeholder="搜索目的地、签证类型..."
            confirm-type="search"
            @confirm="searchVisaCountries"
          />
        </view>
      </view>


      <view class="visa-type-wrap" :class="themeMode === 'dark' ? 'visa-type-wrap--dark' : 'visa-type-wrap--light'">
        <view class="visa-type-row">
          <view
            v-for="item in visaTypeMenu"
            :key="item.key"
            class="visa-type-item"
            @tap="onVisaTypeTap(item)"
          >
            <view class="visa-type-icon-box" :class="{ 'visa-type-icon-box--active': activeVisaType === item.key }">
              <text class="material-symbols-outlined visa-type-icon" :style="{ color: item.color }">{{ item.icon }}</text>
            </view>
            <text class="visa-type-label" :class="{ 'visa-type-label--active': activeVisaType === item.key }">{{ item.label }}</text>
          </view>
        </view>
      </view>

      <view class="insight-grid">
        <view class="insight-primary insight-primary--clickable visago-card" @tap="openVisaFreePage">
          <text class="insight-primary-tip">免签国家</text>
          <text class="insight-primary-count">{{ freeCountryCountText }}</text>
          <text class="insight-primary-desc">覆盖免签、落地签与常见材料清单</text>
          <view class="insight-primary-row">
            <text class="material-symbols-outlined globe">public</text>
          </view>
        </view>

        <view class="insight-right">
          <view class="ai-predict-card ai-predict-card--static visago-card">
            <view class="ai-head">
              <view class="ai-head-top">
                <view class="ai-badge">AI</view>
                <text class="ai-status">内测中</text>
              </view>
              <text class="ai-title">AI出签预测</text>
            </view>
          </view>

          <view class="insight-light visago-card">
            <view class="insight-light-head">
              <text class="insight-light-tip">升级经验</text>
              <text class="insight-light-tag">即将开放</text>
            </view>
            <text class="insight-light-title">准备好下一次冒险了吗？</text>
            <view class="insight-score">
              <text class="score">0 / 3,000</text>
              <text class="goal">🎆</text>
            </view>
            <view class="progress-track">
              <view class="progress-fill" />
            </view>
          </view>
        </view>
      </view>

      <view class="section-header">
        <text class="section-title">我的计划进度</text>
        <text class="section-link" @tap="openPlanPage">{{ planList.length ? '全部' : '去创建' }}</text>
      </view>

      <view v-if="loadingPlans" class="plan-empty plan-empty--loading visago-card">
        <text class="material-symbols-outlined plan-empty-icon">hourglass_top</text>
        <view class="plan-empty-copy">
          <text class="plan-empty-title">正在同步你的目标进度</text>
          <text class="plan-empty-desc">我们会根据实际计划、步骤和任务状态更新首页进度。</text>
        </view>
      </view>

      <swiper
        v-else-if="planList.length"
        class="plan-swiper"
        :circular="planList.length > 1"
        :autoplay="planList.length > 1"
        :interval="3600"
        :duration="420"
      >
        <swiper-item v-for="plan in planList" :key="plan.id">
          <view class="plan-card visago-card" @tap="openPlanDetail(plan.id)">
            <view class="plan-flag">{{ plan.flag }}</view>
            <view class="plan-main">
              <view class="plan-top">
                <text class="plan-name">{{ plan.name }}</text>
                <text class="plan-step">步骤 {{ plan.step }}/{{ plan.totalStep }}</text>
              </view>
              <text class="plan-sub">{{ plan.statusText }}</text>
              <view class="plan-track">
                <view class="plan-track-fill" :style="{ width: `${plan.percent}%` }" />
              </view>
            </view>
            <view class="plan-circle" :style="progressCircleStyle(plan.percent)">
              <view class="plan-circle-inner">{{ plan.percent }}%</view>
            </view>
          </view>
        </swiper-item>
      </swiper>

      <view v-else class="plan-empty plan-empty--idle visago-card" @tap="openVisaPage()">
        <view class="plan-empty-badge">NEW GOAL</view>
        <text class="material-symbols-outlined plan-empty-icon">track_changes</text>
        <view class="plan-empty-copy">
          <text class="plan-empty-title">还没有签证目标</text>
          <text class="plan-empty-desc">创建第一个出行计划后，这里会自动显示真实进度、当前步骤和完成比例。</text>
        </view>
        <view class="plan-empty-btn">去创建计划</view>
      </view>

      <view class="section-header">
        <text class="section-title">热门目的地</text>
        <text class="section-link">更多</text>
      </view>

      <scroll-view class="destination-scroll" scroll-x>
        <view class="destination-row">
          <view v-for="item in filteredDestinations" :key="item.id" class="destination-card visago-card">
            <view class="destination-image-wrap">
              <image class="destination-image" :src="item.image" mode="aspectFill" />
              <view v-if="item.hot" class="hot-tag">热门</view>
            </view>
            <view class="destination-body">
              <view class="destination-top">
                <text class="destination-name">{{ item.name }}</text>
                <text class="destination-flag">{{ item.flag }}</text>
              </view>
              <view class="destination-meta">
                <view class="destination-meta-left">
                  <text class="material-symbols-outlined tiny-icon">verified</text>
                  <text>{{ item.type }} | {{ item.time }}</text>
                </view>
                <view class="destination-meta-right">
                  <text class="destination-price">{{ item.price }}</text>
                  <text class="price-suffix">起</text>
                </view>
              </view>
            </view>
          </view>

          <view v-if="loadingDestinations" class="empty-tip visago-card">
            正在加载热门目的地...
          </view>

          <view v-else-if="!filteredDestinations.length" class="empty-tip visago-card">
            没有匹配目的地，试试输入“日本”或“澳大利亚”。
          </view>
        </view>
      </scroll-view>

      <view class="tool-group">
        <text class="section-title">快捷工具</text>
        <view class="visago-card tool-card">
          <view class="tool-item tool-item--clickable" hover-class="tool-item--active" @tap.stop="openExchangeTool">
            <view class="tool-left">
              <text class="material-symbols-outlined tool-icon">{{ homeToolGroups.quick[0].icon }}</text>
              <text class="tool-title">{{ homeToolGroups.quick[0].title }}</text>
            </view>
            <text class="material-symbols-outlined tool-arrow">chevron_right</text>
          </view>
          <view
            v-for="item in homeToolGroups.quick.slice(1)"
            :key="item.key"
            class="tool-item tool-item--clickable"
            hover-class="tool-item--active"
            @tap.stop="onQuickToolTap(item)"
          >
            <view class="tool-left">
              <text class="material-symbols-outlined tool-icon">{{ item.icon }}</text>
              <text class="tool-title">{{ item.title }}</text>
            </view>
            <text class="material-symbols-outlined tool-arrow">chevron_right</text>
          </view>
        </view>
      </view>

      <view v-if="homeToolGroups.common.length" class="tool-group">
        <text class="section-title">常用工具</text>
        <view class="visago-card tool-card">
          <view
            v-for="item in homeToolGroups.common"
            :key="item.key"
            class="tool-item tool-item--clickable"
            hover-class="tool-item--active"
            @tap.stop="onCommonToolTap(item)"
          >
            <view class="tool-left">
              <text class="material-symbols-outlined tool-icon tool-icon-secondary">{{ item.icon }}</text>
              <text class="tool-title">{{ item.title }}</text>
            </view>
            <text class="material-symbols-outlined tool-arrow">chevron_right</text>
          </view>
        </view>
      </view>
    </view>

    <VisagoBottomNav active-tab="home" />
  </view>
</template>

<script>
import VisagoTopBar from '../../../components/VisagoTopBar.vue'
import VisagoBottomNav from '../../../components/VisagoBottomNav.vue'
import { api } from '../../../utils/api'
import { applyTheme, getStoredTheme } from '../../../utils/theme'

const HOME_TOOL_GROUPS = {
  quick: [
    { key: 'fx', icon: 'currency_exchange', title: '汇率换算' },
    // { key: 'photo-check', icon: 'portrait', title: '证件照标准检查' },
    // { key: 'scan', icon: 'document_scanner', title: '文件扫描件' },
  ],
  common: [
    // { key: 'translate', icon: 'g_translate', title: '翻译助手' },
    // { key: 'insurance', icon: 'health_and_safety', title: '保险购买' },
    // { key: 'embassy', icon: 'apartment', title: '驻华使领馆' },
  ],
}

const HOME_PLAN_STEP_FALLBACK = Object.freeze([
  { stepKey: 'apply', title: '申请' },
  { stepKey: 'docs', title: '材料' },
  { stepKey: 'book', title: '预约' },
  { stepKey: 'result', title: '结果' },
])

function clampPercent(value) {
  const n = Number(value || 0)
  if (!Number.isFinite(n)) return 0
  return Math.max(0, Math.min(100, Math.round(n)))
}

function getPlanSteps(detail, activeStepKey) {
  const detailSteps = Array.isArray(detail && detail.steps) ? detail.steps.filter(Boolean) : []
  if (detailSteps.length) return detailSteps
  return HOME_PLAN_STEP_FALLBACK.map((item) => ({
    ...item,
    status: item.stepKey === activeStepKey ? 'active' : 'todo',
  }))
}

function findActiveStepIndex(steps, activeStepKey) {
  const activeIndex = steps.findIndex((step) => step.status === 'active')
  if (activeIndex >= 0) return activeIndex
  const keyIndex = steps.findIndex((step) => step.stepKey === activeStepKey)
  if (keyIndex >= 0) return keyIndex
  const doneCount = steps.filter((step) => step.status === 'done').length
  return Math.min(doneCount, Math.max(steps.length - 1, 0))
}

function buildPlanStatusText(plan, steps) {
  const activeStep = steps.find((step) => step.status === 'active') || steps[findActiveStepIndex(steps, plan.activeStepKey)]
  if (activeStep && activeStep.title) {
    return `当前进行：${activeStep.title}`
  }
  if (clampPercent(plan.progress) >= 100) {
    return '当前计划已完成，等待处理结果'
  }
  return '计划进行中'
}

function buildHomePlanCard(plan, detail) {
  const steps = getPlanSteps(detail, plan.activeStepKey)
  const currentIndex = findActiveStepIndex(steps, plan.activeStepKey)
  return {
    id: plan.id,
    flag: plan.countryFlag || plan.countryName?.slice(0, 1) || '签',
    name: plan.visaTitle || `${plan.countryName || '签证'}计划`,
    step: Math.min(currentIndex + 1, steps.length || 1),
    totalStep: steps.length || 1,
    percent: clampPercent(plan.progress),
    statusText: buildPlanStatusText(plan, steps),
  }
}

function normalizeHotDestination(item = {}) {
  const countryName = String(item.name || '').trim()
  const visaName = String(item.visaName || '').trim()
  const visaType = String(item.type || '').trim()
  const processingTime = String(item.time || '').trim()
  const price = String(item.price || '').trim()
  const note = String(item.note || '').trim()
  return {
    id: item.visaId || item.countryId || `${countryName}-${visaName}`,
    countryId: Number(item.countryId || 0),
    visaId: Number(item.visaId || 0),
    name: countryName || '热门目的地',
    flag: item.flag || '签',
    price: price || '待更新',
    type: visaType || visaName || '签证',
    time: processingTime || '待更新',
    hot: Boolean(item.hot),
    image: item.image || '',
    note,
    keywords: [countryName, visaName, visaType, processingTime, price, note].filter(Boolean),
  }
}

export default {
  components: {
    VisagoTopBar,
    VisagoBottomNav,
  },
  data() {
    return {
      keyword: '',
      activeVisaType: 'tourism',
      themeMode: 'light',
      loadingFreeCountries: false,
      freeCountryCount: 0,
      loadingPlans: false,
      loadingDestinations: false,
      homeToolGroups: HOME_TOOL_GROUPS,
      destinations: [],
      planList: [],
      visaTypeMenu: [
        { key: 'tourism', label: '旅游签', icon: 'flight', color: '#2f89ff' },
        { key: 'business', label: '商务签', icon: 'business_center', color: '#9333ea' },
        { key: 'study', label: '学生签', icon: 'school', color: '#f59e0b' },
        { key: 'work', label: '工作签', icon: 'description', color: '#22c55e' },
      ],
    }
  },
  computed: {
    freeCountryCountText() {
      if (this.loadingFreeCountries) {
        return '加载中...'
      }
      return `${Number(this.freeCountryCount || 0)} 个地区`
    },
    filteredDestinations() {
      const word = this.keyword.trim().toLowerCase()
      if (!word) {
        return this.destinations
      }
      return this.destinations.filter((item) => item.keywords.some((token) => token.toLowerCase().includes(word)))
    },
  },
  onShow() {
    this.themeMode = applyTheme(getStoredTheme())
    this.loadFreeCountryCount()
    this.loadHomePlans()
    this.loadHotDestinations()
  },
  methods: {
    async loadFreeCountryCount() {
      if (this.loadingFreeCountries) return
      this.loadingFreeCountries = true
      try {
        const items = (await api.listFreeCountries({ enabled: 1 })) || []
        this.freeCountryCount = Array.isArray(items) ? items.length : 0
      } catch (error) {
        this.freeCountryCount = 0
      } finally {
        this.loadingFreeCountries = false
      }
    },
    async loadHotDestinations() {
      if (this.loadingDestinations) return
      this.loadingDestinations = true
      try {
        const items = (await api.listHotDestinations({ limit: 8 })) || []
        this.destinations = items.map((item) => normalizeHotDestination(item))
      } catch (error) {
        this.destinations = []
      } finally {
        this.loadingDestinations = false
      }
    },
    async loadHomePlans() {
      if (this.loadingPlans) return
      this.loadingPlans = true
      try {
        const plans = (await api.listPlans()) || []
        const activePlans = plans.filter((item) => String(item.resultStatus || 'pending') === 'pending')
        if (!activePlans.length) {
          this.planList = []
          return
        }
        const details = await Promise.all(
          activePlans.slice(0, 5).map((item) =>
            api
              .getPlanDetail(item.id)
              .then((detail) => detail || null)
              .catch(() => null)
          )
        )
        this.planList = activePlans.slice(0, 5).map((plan, index) => buildHomePlanCard(plan, details[index]))
      } catch (error) {
        this.planList = []
      } finally {
        this.loadingPlans = false
      }
    },
    openVisaPage(keyword = '') {
      const search = String(keyword || '').trim()
      const url = search ? `/pages/visago/visa/index?keyword=${encodeURIComponent(search)}` : '/pages/visago/visa/index'
      uni.navigateTo({
        url,
        fail: (error) => {
          const msg = (error && error.errMsg) || '打开失败'
          uni.showToast({
            title: msg.slice(0, 18),
            icon: 'none',
          })
        },
      })
    },
    searchVisaCountries() {
      this.openVisaPage(this.keyword)
    },
    openPlanPage() {
      if (this.planList.length) {
        uni.navigateTo({ url: '/pages/visago/plan/index' })
        return
      }
      this.openVisaPage()
    },
    openPlanDetail(planId) {
      uni.navigateTo({
        url: `/pages/visago/plan/index?planId=${planId}`,
        fail: (error) => {
          const msg = (error && error.errMsg) || '打开失败'
          uni.showToast({
            title: msg.slice(0, 18),
            icon: 'none',
          })
        },
      })
    },
    onVisaTypeTap(item) {
      this.activeVisaType = item.key
      this.openVisaPage()
    },
    progressCircleStyle(percent) {
      const safe = Math.max(0, Math.min(100, Number(percent) || 0))
      return {
        background: `conic-gradient(var(--visago-primary) ${safe}%, rgba(15, 101, 216, 0.18) ${safe}% 100%)`,
      }
    },
    openVisaFreePage() {
      uni.navigateTo({
        url: '/pages/visago/visa-free/index',
        fail: () => {
          uni.showToast({
            title: '打开失败',
            icon: 'none',
          })
        },
      })
    },
    openExchangeTool() {
      uni.navigateTo({
        url: '/pages/visago/tools/exchange/index',
        fail: (error) => {
          const msg = (error && error.errMsg) || '打开失败'
          uni.showToast({
            title: msg.slice(0, 18),
            icon: 'none',
          })
        },
      })
    },
    onQuickToolTap(item) {
      if (item.key === 'fx') {
        this.openExchangeTool()
      }
    },
    onCommonToolTap(item) {
      if (item.key === 'embassy') {
        uni.navigateTo({
          url: '/pages/visago/tools/embassy/index',
          fail: (error) => {
            const msg = (error && error.errMsg) || '打开失败'
            uni.showToast({
              title: msg.slice(0, 18),
              icon: 'none',
            })
          },
        })
        return
      }
    },
  },
}
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background: var(--visago-bg);
}

.home-content {
  box-sizing: border-box;
  padding: 96px 16px calc(120px + var(--visago-safe-bottom));
}

.hero-section {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.hero-title {
  font-size: 52rpx;
  font-weight: 700;
  color: var(--visago-text);
}

.search-box {
  height: 52px;
  border-radius: 14px;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  padding: 0 14px;
  gap: 8px;
}

.search-icon {
  color: var(--visago-text-soft);
  font-size: 24px;
}

.search-input {
  flex: 1;
  height: 52px;
  font-size: 18px;
}

.ai-predict-card {
  padding: 10px 12px;
  min-height: 64px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  border: 1px solid color-mix(in srgb, var(--visago-primary) 22%, transparent);
  box-shadow: 0 10px 24px rgba(15, 101, 216, 0.12);
  background: linear-gradient(
    140deg,
    color-mix(in srgb, var(--visago-primary) 14%, var(--visago-surface) 86%) 0%,
    color-mix(in srgb, #ffffff 30%, var(--visago-surface) 70%) 36%,
    var(--visago-surface) 68%
  );
}

.ai-head {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 6px;
}

.ai-head-top {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.ai-predict-card--static {
  cursor: default;
}

.ai-badge {
  min-width: 28px;
  height: 20px;
  padding: 0 6px;
  border-radius: 999px;
  background: var(--visago-primary);
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ai-title {
  width: 100%;
  font-size: 15px;
  font-weight: 700;
  color: var(--visago-text);
  line-height: 1.2;
}

.ai-status {
  flex-shrink: 0;
  padding: 3px 8px;
  border-radius: 999px;
  background: rgba(15, 101, 216, 0.1);
  color: var(--visago-primary);
  font-size: 11px;
  font-weight: 700;
}

.ai-arrow {
  color: var(--visago-primary);
  font-size: 18px;
}

.ai-sub {
  display: block;
  margin-top: 6px;
  font-size: 11px;
  color: var(--visago-text-soft);
  line-height: 1.35;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
}

.visa-type-wrap {
  margin-top: 22px;
  padding: 0;
}

.visa-type-wrap--dark {
  background: transparent;
  border: none;
}

.visa-type-wrap--light {
  background: transparent;
  border: none;
}

.visa-type-row {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.visa-type-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.visa-type-icon-box {
  width: 100%;
  aspect-ratio: 1 / 1;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid transparent;
}

.visa-type-wrap--dark .visa-type-icon-box {
  background: #171b24;
}

.visa-type-wrap--light .visa-type-icon-box {
  background: #ffffff;
  border-color: #e7edf8;
}

.visa-type-icon-box--active {
  transform: translateY(-1px);
}

.visa-type-wrap--dark .visa-type-icon-box--active {
  background: #1f2633;
  box-shadow: inset 0 0 0 1px rgba(180, 190, 210, 0.24);
}

.visa-type-wrap--light .visa-type-icon-box--active {
  background: #f8fbff;
  border-color: rgba(15, 101, 216, 0.28);
}

.visa-type-icon {
  font-size: 26px;
}

.visa-type-label {
  font-size: 13px;
  line-height: 1;
  font-weight: 500;
}

.visa-type-wrap--dark .visa-type-label {
  color: #9aa3b8;
}

.visa-type-wrap--light .visa-type-label {
  color: #65708a;
}

.visa-type-label--active {
  font-weight: 600;
}

.visa-type-wrap--dark .visa-type-label--active {
  color: #d7deec;
}

.visa-type-wrap--light .visa-type-label--active {
  color: #1d4fa8;
}

.insight-grid {
  margin-top: 22px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.insight-right {
  display: grid;
  grid-template-rows: 1fr 1fr;
  gap: 10px;
}

.insight-primary {
  background: linear-gradient(160deg, #0f65d8 0%, #0759c3 100%);
  color: #fff;
  padding: 14px;
  min-height: 150px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 12px 26px rgba(15, 101, 216, 0.28);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.insight-primary--clickable:active {
  transform: translateY(1px);
}

.insight-primary-tip {
  font-size: 12px;
  letter-spacing: 0.3px;
  opacity: 0.88;
}

.insight-primary-count {
  margin-top: 6px;
  font-size: 42rpx;
  font-weight: 700;
  line-height: 1.2;
}

.insight-primary-desc {
  margin-top: 6px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.84);
  line-height: 1.4;
}

.insight-primary-row {
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.globe {
  font-size: 58px;
  color: rgba(255, 255, 255, 0.2);
}

.insight-light {
  padding: 10px 12px;
  min-height: 64px;
  box-sizing: border-box;
  border: 1px solid color-mix(in srgb, var(--visago-primary) 14%, var(--visago-line));
  box-shadow: 0 8px 22px rgba(17, 24, 39, 0.08);
  background: linear-gradient(
    180deg,
    color-mix(in srgb, var(--visago-primary) 5%, var(--visago-surface) 95%) 0%,
    var(--visago-surface) 62%
  );
  display: flex;
  flex-direction: column;
}

.insight-light-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.insight-light-tip {
  color: var(--visago-text-soft);
  font-size: 13px;
}

.insight-light-tag {
  flex-shrink: 0;
  padding: 3px 8px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--visago-primary) 10%, var(--visago-surface) 90%);
  color: var(--visago-primary);
  font-size: 10px;
  font-weight: 700;
}

.insight-light-title {
  margin-top: 4px;
  font-size: 14px;
  line-height: 1.25;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.insight-score {
  margin-top: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.score {
  color: var(--visago-primary);
  font-weight: 700;
  font-size: 14px;
}

.goal {
  font-size: 14px;
}

.progress-track {
  margin-top: 5px;
  height: 8px;
  width: 100%;
  border-radius: 9999px;
  background: var(--visago-surface-soft);
  overflow: hidden;
}

.progress-fill {
  width: 81%;
  height: 100%;
  border-radius: 9999px;
  background: var(--visago-primary);
}

.insight-light-sub {
  margin-top: 4px;
  font-size: 11px;
  color: var(--visago-text-soft);
  text-align: right;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.section-header {
  margin-top: 28px;
  margin-bottom: 10px;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.section-title {
  font-size: 36rpx;
  font-weight: 700;
}

.section-link {
  font-size: 14px;
  color: var(--visago-primary);
}

.plan-swiper {
  height: 122px;
}

.plan-empty {
  min-height: 136px;
  padding: 16px 18px;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  gap: 14px;
}

.plan-empty--loading {
  border: 1px solid color-mix(in srgb, var(--visago-primary) 10%, var(--visago-line));
  background: linear-gradient(
    135deg,
    color-mix(in srgb, var(--visago-primary) 7%, var(--visago-surface) 93%) 0%,
    var(--visago-surface) 100%
  );
}

.plan-empty--idle {
  position: relative;
  overflow: hidden;
  border: 1px solid color-mix(in srgb, var(--visago-primary) 18%, var(--visago-line));
  background:
    radial-gradient(circle at top right, color-mix(in srgb, var(--visago-primary) 16%, transparent) 0%, transparent 38%),
    linear-gradient(145deg, var(--visago-surface) 0%, color-mix(in srgb, var(--visago-primary) 6%, var(--visago-surface) 94%) 100%);
}

.plan-empty-badge {
  position: absolute;
  top: 14px;
  right: 14px;
  padding: 4px 8px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--visago-primary) 12%, var(--visago-surface) 88%);
  color: var(--visago-primary);
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.08em;
}

.plan-empty-icon {
  width: 48px;
  height: 48px;
  border-radius: 16px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: color-mix(in srgb, var(--visago-primary) 12%, var(--visago-surface) 88%);
  color: var(--visago-primary);
  font-size: 24px;
}

.plan-empty-copy {
  min-width: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.plan-empty-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--visago-text);
}

.plan-empty-desc {
  font-size: 12px;
  line-height: 1.55;
  color: var(--visago-text-soft);
}

.plan-empty-btn {
  flex-shrink: 0;
  padding: 10px 14px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--visago-primary) 0%, var(--visago-primary-strong) 100%);
  color: #fff;
  font-size: 13px;
  font-weight: 700;
}

.plan-card {
  height: 122px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  box-sizing: border-box;
}

.plan-card:active {
  transform: translateY(1px);
}

.plan-flag {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 30rpx;
}

.plan-main {
  flex: 1;
  min-width: 0;
}

.plan-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.plan-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--visago-text);
}

.plan-step {
  font-size: 11px;
  color: var(--visago-text-soft);
  flex-shrink: 0;
}

.plan-sub {
  margin-top: 4px;
  display: block;
  font-size: 12px;
  color: var(--visago-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.plan-track {
  margin-top: 8px;
  height: 7px;
  border-radius: 9999px;
  background: var(--visago-surface-soft);
  overflow: hidden;
}

.plan-track-fill {
  height: 100%;
  border-radius: 9999px;
  background: var(--visago-primary);
}

.plan-circle {
  width: 54px;
  height: 54px;
  border-radius: 9999px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.plan-circle-inner {
  width: 42px;
  height: 42px;
  border-radius: 9999px;
  background: var(--visago-surface);
  color: var(--visago-text);
  font-size: 11px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.destination-scroll {
  width: 100%;
  white-space: nowrap;
}

.destination-row {
  display: inline-flex;
  gap: 14px;
  padding-right: 8px;
}

.destination-card {
  width: 272px;
  overflow: hidden;
}

.destination-image-wrap {
  position: relative;
  height: 130px;
}

.destination-image {
  width: 100%;
  height: 100%;
}

.hot-tag {
  position: absolute;
  right: 10px;
  top: 8px;
  font-size: 11px;
  color: #fff;
  background: rgba(0, 0, 0, 0.55);
  padding: 4px 9px;
  border-radius: 9999px;
}

.destination-body {
  padding: 12px 14px;
}

.destination-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.destination-name {
  font-size: 18px;
  font-weight: 600;
}

.destination-flag {
  font-size: 28rpx;
  line-height: 1;
}

.destination-price {
  color: var(--visago-primary);
  font-size: 30rpx;
  font-weight: 700;
}

.price-suffix {
  margin-left: 2px;
  color: var(--visago-text-soft);
  font-size: 11px;
}

.destination-meta {
  margin-top: 6px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: var(--visago-text-muted);
  font-size: 13px;
}

.tiny-icon {
  font-size: 14px;
}

.destination-meta-left,
.destination-meta-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.empty-tip {
  width: 272px;
  min-height: 86px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  font-size: 14px;
  color: var(--visago-text-muted);
  text-align: center;
}

.tool-group {
  margin-top: 28px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tool-card {
  overflow: hidden;
}

.tool-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid var(--visago-line);
}

.tool-item--clickable {
  cursor: pointer;
}

.tool-item--active,
.tool-item--clickable:active {
  background: var(--visago-surface-soft);
}

.tool-item:last-child {
  border-bottom: none;
}

.tool-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.tool-icon {
  color: var(--visago-primary);
  font-size: 22px;
}

.tool-icon-secondary {
  color: var(--visago-secondary);
}

.tool-title {
  font-size: 16px;
}

.tool-arrow {
  color: var(--visago-text-soft);
  font-size: 19px;
}
</style>

