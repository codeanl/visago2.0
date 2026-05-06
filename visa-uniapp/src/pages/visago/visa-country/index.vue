<template>
  <view class="country-visa-page" :class="{ 'country-visa-page--dark': themeMode === 'dark' }">
    <view class="hero">
      <image class="hero-image" :src="country.image || defaultHero" mode="aspectFill" />
      <view class="hero-mask" />
      <view class="hero-top">
        <view class="hero-icon-btn" @tap="goBack">
          <text class="material-symbols-outlined">arrow_back</text>
        </view>
      </view>
      <view class="hero-bottom">
        <text class="hero-title">{{ country.name || '国家签证' }}</text>
        <text class="hero-sub">{{ country.region || '全球' }} · 可选签证类型</text>
      </view>
    </view>

    <view class="page-body visago-page-width">
      <view class="section-head">
        <text class="section-title">可选签证类型</text>
      </view>

      <view class="visa-stack">
        <view v-if="loading" class="visa-card visago-card">
          <view class="visa-title-row">
            <text class="visa-title">正在加载签证列表</text>
          </view>
          <text class="visa-desc">请稍等，正在从数据库读取该国家下的签证信息。</text>
        </view>

        <view v-for="visa in visaOptions" :key="visa.id" class="visa-card visago-card" @tap="openVisaDetail(visa)">
          <view class="visa-title-row">
            <text class="visa-title">{{ visa.name }}</text>
            <view class="visa-tag-row">
              <text v-if="visa.hot" class="visa-hot">热门</text>
              <text v-if="visa.visaFree" class="visa-free">免签</text>
            </view>
          </view>
          <text class="visa-desc">{{ visa.description || visa.longIntro }}</text>
          <view class="visa-meta-grid">
            <view class="meta-item">
              <text class="meta-label">办理时长</text>
              <text class="meta-value">{{ visa.processingTime || '-' }}</text>
            </view>
            <view class="meta-item">
              <text class="meta-label">有效期</text>
              <text class="meta-value">{{ visa.validity || '-' }}</text>
            </view>
            <view class="meta-item">
              <text class="meta-label">预估费用</text>
              <text class="meta-value meta-value--price">{{ visa.fee || '-' }}</text>
            </view>
          </view>
        </view>

        <view v-if="!loading && !visaOptions.length" class="visa-card visago-card">
          <view class="visa-title-row">
            <text class="visa-title">还没有签证配置</text>
          </view>
          <text class="visa-desc">请先到后台管理端为这个国家新增签证和详细步骤。</text>
        </view>
      </view>
    </view>

    <view v-if="activeVisa" class="detail-mask" @tap="closeVisaDetail">
      <view class="detail-sheet" @tap.stop>
        <view class="detail-head">
          <text class="detail-title">{{ activeVisa.name }}</text>
          <view class="hero-icon-btn hero-icon-btn--detail" @tap="closeVisaDetail">
            <text class="material-symbols-outlined">close</text>
          </view>
        </view>

        <view v-if="detailLoading" class="detail-loading">
          <text>正在加载签证详情...</text>
        </view>

        <template v-else>
          <view class="badge-row">
            <text class="badge">{{ activeVisa.visaType || '签证类型' }}</text>
            <text v-if="activeVisa.visaFree" class="badge badge--free">免签</text>
            <text class="badge">{{ activeVisa.entries || '入境次数待定' }}</text>
          </view>

          <text class="detail-intro">{{ activeVisa.longIntro || activeVisa.description }}</text>

          <view class="detail-grid">
            <view class="detail-box">
              <text class="material-symbols-outlined detail-box-icon">calendar_month</text>
              <text class="detail-box-label">办理周期</text>
              <text class="detail-box-value">{{ activeVisa.processingTime || '-' }}</text>
            </view>
            <view class="detail-box">
              <text class="material-symbols-outlined detail-box-icon">payments</text>
              <text class="detail-box-label">预估费用</text>
              <text class="detail-box-value">{{ activeVisa.fee || '-' }}</text>
            </view>
            <view class="detail-box">
              <text class="material-symbols-outlined detail-box-icon">assignment_ind</text>
              <text class="detail-box-label">护照要求</text>
              <text class="detail-box-value">{{ passportRequirement }}</text>
            </view>
            <view class="detail-box">
              <text class="material-symbols-outlined detail-box-icon">account_balance_wallet</text>
              <text class="detail-box-label">材料重点</text>
              <text class="detail-box-value">{{ materialHighlight }}</text>
            </view>
          </view>

          <view class="join-btn" :class="{ 'join-btn--disabled': creatingPlan }" @tap="addToPlan">
            <text>{{ creatingPlan ? '正在加入计划...' : '加入我的计划' }}</text>
          </view>
        </template>
      </view>
    </view>
  </view>
</template>

<script>
import { api } from '../../../utils/api'
import { applyTheme, getStoredTheme, THEME_CHANGE_EVENT } from '../../../utils/theme'

const defaultHero = 'https://images.pexels.com/photos/918275/pexels-photo-918275.jpeg?auto=compress&cs=tinysrgb&w=1200'

export default {
  data() {
    return {
      defaultHero,
      themeMode: 'light',
      countryId: 0,
      loading: false,
      detailLoading: false,
      creatingPlan: false,
      country: {
        id: 0,
        name: '',
        region: '',
        image: '',
      },
      visaOptions: [],
      activeVisa: null,
    }
  },
  computed: {
    passportRequirement() {
      return '回国后至少 6 个月有效'
    },
    materialHighlight() {
      return '近 3-6 个月基础申请材料'
    },
  },
  onLoad(query) {
    this.countryId = Number(query.countryId || 0)
    if (typeof uni !== 'undefined' && uni.$on) {
      uni.$on(THEME_CHANGE_EVENT, this.onThemeChange)
    }
    this.syncTheme()
    this.loadCountryAndVisas()
  },
  onShow() {
    this.syncTheme()
  },
  onUnload() {
    if (typeof uni !== 'undefined' && uni.$off) {
      uni.$off(THEME_CHANGE_EVENT, this.onThemeChange)
    }
  },
  methods: {
    syncTheme() {
      this.themeMode = applyTheme(getStoredTheme())
    },
    onThemeChange(theme) {
      this.themeMode = theme === 'dark' ? 'dark' : 'light'
      applyTheme(this.themeMode)
    },
    async loadCountryAndVisas() {
      if (!this.countryId) return
      this.loading = true
      try {
        const [countries, visas] = await Promise.all([api.listCountries(), api.listVisasByCountry(this.countryId)])
        this.country = (countries || []).find((item) => Number(item.id) === this.countryId) || this.country
        this.visaOptions = visas || []
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '加载失败',
          icon: 'none',
        })
      } finally {
        this.loading = false
      }
    },
    goBack() {
      uni.navigateBack({ delta: 1 })
    },
    async openVisaDetail(visa) {
      this.activeVisa = { ...visa, steps: [] }
      this.detailLoading = true
      try {
        this.activeVisa = await api.getVisaDetail(visa.id)
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '详情加载失败',
          icon: 'none',
        })
        this.activeVisa = null
      } finally {
        this.detailLoading = false
      }
    },
    closeVisaDetail() {
      this.activeVisa = null
      this.detailLoading = false
    },
    async addToPlan() {
      if (!this.activeVisa || this.creatingPlan) return
      this.creatingPlan = true
      try {
        const currentPlans = (await api.listPlans()) || []
        const existed = currentPlans.find((item) => Number(item.visaId) === Number(this.activeVisa.id) && item.status === 'active')
        if (existed) {
          uni.showToast({
            title: '该签证已有进行中的申请',
            icon: 'none',
          })
          uni.navigateTo({
            url: `/pages/visago/plan/index?planId=${existed.id}`,
          })
          return
        }

        const plan = await api.createPlan({
          countryId: this.country.id,
          countryName: this.country.name,
          visaId: this.activeVisa.id,
          visaTitle: this.activeVisa.name,
          source: 'visa',
        })

        uni.navigateTo({
          url: `/pages/visago/plan/index?planId=${plan.id}`,
        })
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '加入计划失败',
          icon: 'none',
        })
      } finally {
        this.creatingPlan = false
      }
    },
  },
}
</script>

<style scoped>
.country-visa-page {
  --cv-bg: #f3f2f8;
  --cv-text: #101827;
  --cv-text-soft: #63708d;
  --cv-text-muted: #8b96ad;
  --cv-surface: #ffffff;
  --cv-surface-alt: #f7f9fc;
  --cv-line: #dde3ef;
  --cv-hero-mask-top: rgba(14, 29, 53, 0.2);
  --cv-hero-mask-bottom: rgba(14, 29, 53, 0.62);
  --cv-btn-bg: rgba(255, 255, 255, 0.85);
  --cv-btn-text: #2b3c5a;
  --cv-primary: #0f65d8;
  --cv-hot: #d97706;
  min-height: 100vh;
  background: var(--cv-bg);
  color: var(--cv-text);
}

.hero {
  position: relative;
  height: 276px;
  overflow: hidden;
}

.hero-image {
  width: 100%;
  height: 100%;
}

.hero-mask {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, var(--cv-hero-mask-top) 0%, var(--cv-hero-mask-bottom) 100%);
}

.hero-top {
  position: absolute;
  top: 12px;
  left: 12px;
  right: 12px;
  display: flex;
  justify-content: flex-start;
}

.hero-icon-btn {
  width: 38px;
  height: 38px;
  border-radius: 9999px;
  background: var(--cv-btn-bg);
  color: var(--cv-btn-text);
  display: flex;
  align-items: center;
  justify-content: center;
}

.hero-icon-btn--detail {
  background: var(--cv-surface-alt);
}

.hero-bottom {
  position: absolute;
  left: 16px;
  right: 16px;
  bottom: 16px;
}

.hero-title {
  display: block;
  font-size: 33px;
  font-weight: 700;
  color: #fff;
}

.hero-sub {
  display: block;
  margin-top: 4px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.9);
}

.page-body {
  margin-top: -14px;
  position: relative;
  z-index: 2;
  border-radius: 18px 18px 0 0;
  background: var(--cv-bg);
  padding: 18px 16px 24px;
  box-sizing: border-box;
}

.section-head {
  margin-bottom: 12px;
}

.section-title {
  font-size: 22px;
  font-weight: 700;
}

.visa-stack {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding-bottom: 16px;
}

.visa-card {
  background: var(--cv-surface);
  border: 1px solid var(--cv-line);
  border-radius: 14px;
  padding: 12px;
}

.visa-title-row {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  align-items: center;
}

.visa-tag-row {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.visa-title {
  font-size: 16px;
  font-weight: 700;
}

.visa-hot {
  padding: 1px 7px;
  border-radius: 6px;
  font-size: 10px;
  color: var(--cv-hot);
  border: 1px solid rgba(217, 119, 6, 0.28);
  background: rgba(217, 119, 6, 0.12);
}

.visa-free {
  padding: 1px 7px;
  border-radius: 6px;
  font-size: 10px;
  color: #15803d;
  border: 1px solid rgba(21, 128, 61, 0.24);
  background: rgba(34, 197, 94, 0.12);
}

.visa-desc {
  margin-top: 6px;
  display: block;
  font-size: 12px;
  color: var(--cv-text-soft);
  line-height: 1.45;
}

.visa-meta-grid {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid var(--cv-line);
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 6px;
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.meta-label {
  font-size: 10px;
  color: var(--cv-text-muted);
}

.meta-value {
  font-size: 12px;
  color: var(--cv-text);
  font-weight: 600;
}

.meta-value--price {
  color: var(--cv-primary);
}

.detail-mask {
  position: fixed;
  inset: 0;
  z-index: 99;
  background: rgba(0, 0, 0, 0.42);
  display: flex;
  align-items: flex-end;
}

.detail-sheet {
  width: 100%;
  max-height: 88vh;
  background: var(--cv-bg);
  border-radius: 18px 18px 0 0;
  padding: 14px 16px calc(16px + var(--visago-safe-bottom));
  box-sizing: border-box;
  overflow-y: auto;
}

.detail-head {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  align-items: center;
}

.detail-title {
  font-size: 20px;
  font-weight: 700;
}

.detail-loading {
  margin-top: 18px;
  padding: 20px 0;
  text-align: center;
  color: var(--cv-text-soft);
  font-size: 13px;
}

.badge-row {
  margin-top: 10px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.badge {
  padding: 3px 10px;
  border-radius: 8px;
  font-size: 11px;
  color: var(--cv-primary);
  background: rgba(15, 101, 216, 0.12);
  border: 1px solid rgba(15, 101, 216, 0.2);
}

.badge--free {
  color: #15803d;
  background: rgba(34, 197, 94, 0.12);
  border-color: rgba(21, 128, 61, 0.22);
}

.detail-intro {
  display: block;
  margin-top: 10px;
  font-size: 13px;
  line-height: 1.5;
  color: var(--cv-text-soft);
}

.detail-grid {
  margin-top: 12px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.detail-box {
  background: var(--cv-surface);
  border: 1px solid var(--cv-line);
  border-radius: 12px;
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.detail-box-icon {
  font-size: 20px;
  color: var(--cv-primary);
}

.detail-box-label {
  font-size: 11px;
  color: var(--cv-text-muted);
}

.detail-box-value {
  font-size: 12px;
  color: var(--cv-text);
  font-weight: 600;
  line-height: 1.35;
}

.process-wrap {
  margin-top: 16px;
}

.process-title {
  font-size: 18px;
  font-weight: 700;
}

.process-item {
  margin-top: 10px;
  display: flex;
  gap: 10px;
}

.process-dot {
  width: 22px;
  height: 22px;
  border-radius: 9999px;
  background: rgba(15, 101, 216, 0.12);
  border: 1px solid rgba(15, 101, 216, 0.24);
  color: var(--cv-primary);
  font-size: 11px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.process-main {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.process-item-title {
  font-size: 14px;
  color: var(--cv-text);
  font-weight: 600;
}

.process-item-desc {
  font-size: 12px;
  line-height: 1.45;
  color: var(--cv-text-soft);
}


.join-btn {
  margin-top: 16px;
  height: 46px;
  border-radius: 9999px;
  background: var(--cv-primary);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
}

.join-btn--disabled {
  opacity: 0.72;
}

.country-visa-page--dark,
:global(html.theme-dark) .country-visa-page {
  --cv-bg: #000000;
  --cv-text: #f8fafc;
  --cv-text-soft: #cbd5e1;
  --cv-text-muted: #94a3b8;
  --cv-surface: #171717;
  --cv-surface-alt: #1e293b;
  --cv-line: rgba(255, 255, 255, 0.08);
  --cv-hero-mask-top: rgba(0, 0, 0, 0.32);
  --cv-hero-mask-bottom: rgba(0, 0, 0, 0.76);
  --cv-btn-bg: rgba(0, 0, 0, 0.45);
  --cv-btn-text: rgba(255, 255, 255, 0.92);
  --cv-primary: #60a5fa;
  --cv-hot: #ff9f0a;
}
</style>
