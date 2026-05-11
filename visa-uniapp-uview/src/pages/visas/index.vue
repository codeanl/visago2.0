<template>
  <view class="app-screen visas-screen">
    <UPageHeader :title="country.name || countryName || '签证列表'" :show-back="true" />
    <view class="page-top-spacer" :style="topSpacerStyle" />

    <view class="hero-banner">
      <image class="hero-image" :src="country.image || defaultHero" mode="aspectFill" />
      <view class="hero-mask" />
      <view class="hero-copy">
        <text class="hero-title">{{ country.name || countryName || '国家签证' }}</text>
        <text class="hero-sub">{{ country.region || '全球' }} · 可选签证类型</text>
      </view>
    </view>

    <view class="app-page-width visas-body">
      <text class="app-section-title">可选签证类型</text>

      <view class="visa-stack">
        <view v-if="loading" class="app-card visa-card visa-card--empty">
          <u-loading-icon mode="circle" size="22" color="#1677ff" />
          <text class="visa-empty-title">正在加载签证列表</text>
          <text class="visa-empty-desc">请稍等，正在读取该国家下的签证信息。</text>
        </view>

        <view
          v-for="visa in visaOptions"
          :key="visa.id"
          class="app-card visa-card"
          @tap="openVisaDetail(visa)"
        >
          <view class="visa-title-row">
            <view class="visa-title-wrap">
              <text class="visa-title">{{ visa.name }}</text>
              <u-tag v-if="visa.visaFree" text="免签" type="warning" size="mini" />
            </view>
            <u-tag v-if="visa.hot" text="热门" type="primary" size="mini" />
          </view>

          <text class="visa-desc">{{ visa.description || visa.longIntro || '查看签证详情与办理信息' }}</text>

          <view class="meta-row">
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

        <view v-if="!loading && !visaOptions.length" class="app-card visa-card visa-card--empty">
          <text class="visa-empty-title">还没有签证配置</text>
          <text class="visa-empty-desc">这个国家暂时还没有可展示的签证信息。</text>
        </view>
      </view>
    </view>

    <view v-if="activeVisa" class="drawer-mask" @tap="closeVisaDetail">
      <view class="drawer-sheet" @tap.stop>
        <view class="drawer-handle" />

        <view class="drawer-body">
          <view class="drawer-head">
            <view class="drawer-icon">
              <u-icon :name="activeVisa.visaFree ? 'map' : 'bookmark'" size="18" color="#1677ff" />
            </view>
            <view class="drawer-title-wrap">
              <text class="drawer-title">{{ activeVisa.name }}</text>
              <text class="drawer-sub">{{ activeVisa.description || activeVisa.longIntro || '查看签证详情与办理信息' }}</text>
            </view>
            <view class="drawer-close" @tap="closeVisaDetail">
              <u-icon name="close" size="16" color="#86909c" />
            </view>
          </view>

          <view v-if="detailLoading" class="detail-loading">
            <u-loading-icon mode="circle" size="22" color="#1677ff" />
            <text class="detail-loading-text">正在加载签证详情...</text>
          </view>

          <scroll-view v-else scroll-y class="drawer-scroll">
            <view class="drawer-block">
              <view class="badge-row">
                <u-tag :text="activeVisa.visaType || '签证类型'" type="primary" size="mini" />
                <u-tag v-if="activeVisa.visaFree" text="免签" type="warning" size="mini" />
                <u-tag :text="activeVisa.entries || '入境次数待定'" type="info" size="mini" />
              </view>

              <view class="intro-card">
                <u-icon name="file-text" size="16" color="#1677ff" />
                <text class="intro-text">{{ activeVisa.longIntro || activeVisa.description }}</text>
              </view>

              <view class="detail-grid">
                <view class="detail-box">
                  <u-icon name="clock" size="16" color="#1677ff" />
                  <text class="detail-box-label">办理周期</text>
                  <text class="detail-box-value">{{ activeVisa.processingTime || '-' }}</text>
                </view>
                <view class="detail-box">
                  <u-icon name="rmb" size="16" color="#1677ff" />
                  <text class="detail-box-label">预估费用</text>
                  <text class="detail-box-value">{{ activeVisa.fee || '-' }}</text>
                </view>
                <view class="detail-box">
                  <u-icon name="bookmark" size="16" color="#1677ff" />
                  <text class="detail-box-label">有效期</text>
                  <text class="detail-box-value">{{ activeVisa.validity || '-' }}</text>
                </view>
                <view class="detail-box">
                  <u-icon name="calendar" size="16" color="#1677ff" />
                  <text class="detail-box-label">入境次数</text>
                  <text class="detail-box-value">{{ activeVisa.entries || '-' }}</text>
                </view>
              </view>
            </view>

            <view v-if="activeVisa.steps && activeVisa.steps.length" class="drawer-block">
              <text class="drawer-block-title">办理步骤</text>
              <view
                v-for="(step, idx) in activeVisa.steps"
                :key="step.id || step.stepKey || idx"
                class="step-card"
              >
                <view class="step-top">
                  <view class="step-index">{{ idx + 1 }}</view>
                  <text class="step-title">{{ step.title }}</text>
                </view>
                <text class="step-desc">{{ step.desc || (step.strategies || []).join('；') || '查看详情后再准备下一步。' }}</text>
              </view>
            </view>

            <u-button
              :text="creatingGoal ? '正在加入目标...' : goalExists ? '去查看我的目标' : '加入我的目标'"
              type="primary"
              shape="circle"
              :disabled="creatingGoal"
              customStyle="margin-top: 8px;"
              @click="addToGoal"
            />
          </scroll-view>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import UPageHeader from '../../components/shared/UPageHeader.vue'
import { api } from '../../utils/api'
import { getGoalByVisaId, upsertGoal } from '../../utils/goals'
import { getTopSpacerStyle } from '../../utils/layout'

const defaultHero = 'https://images.pexels.com/photos/918275/pexels-photo-918275.jpeg?auto=compress&cs=tinysrgb&w=1200'

export default {
  components: { UPageHeader },
  data() {
    return {
      defaultHero,
      countryId: 0,
      countryName: '',
      loading: false,
      detailLoading: false,
      creatingGoal: false,
      country: {
        id: 0,
        name: '',
        region: '',
        image: '',
        flag: '',
      },
      visaOptions: [],
      activeVisa: null,
      topSpacerStyle: 'height: 96px;',
    }
  },
  computed: {
    goalExists() {
      return this.activeVisa ? Boolean(getGoalByVisaId(this.activeVisa.id)) : false
    },
  },
  onLoad(query) {
    this.syncLayout()
    this.countryId = Number(query.countryId || 0)
    this.countryName = decodeURIComponent(String(query.countryName || ''))
  },
  onShow() {
    this.syncLayout()
    this.loadCountryAndVisas()
  },
  methods: {
    syncLayout() {
      this.topSpacerStyle = getTopSpacerStyle('sub')
    },
    async loadCountryAndVisas() {
      if (!this.countryId) return
      this.loading = true
      try {
        const [countries, visas] = await Promise.all([
          api.listCountries('', ''),
          api.listVisasByCountry(this.countryId),
        ])
        this.country = (countries || []).find((item) => Number(item.id) === this.countryId) || {
          ...this.country,
          id: this.countryId,
          name: this.countryName,
        }
        this.visaOptions = visas || []
      } catch (error) {
        uni.showToast({ title: (error && error.message) || '加载失败', icon: 'none' })
      } finally {
        this.loading = false
      }
    },
    async openVisaDetail(visa) {
      this.activeVisa = { ...visa, steps: [] }
      this.detailLoading = true
      try {
        this.activeVisa = await api.getVisaDetail(visa.id)
      } catch (error) {
        uni.showToast({ title: (error && error.message) || '详情加载失败', icon: 'none' })
        this.activeVisa = null
      } finally {
        this.detailLoading = false
      }
    },
    closeVisaDetail() {
      this.activeVisa = null
      this.detailLoading = false
    },
    addToGoal() {
      if (!this.activeVisa || this.creatingGoal) return
      if (this.goalExists) {
        this.closeVisaDetail()
        uni.reLaunch({ url: '/pages/goals/index' })
        return
      }
      this.creatingGoal = true
      try {
        upsertGoal({
          visaId: this.activeVisa.id,
          countryId: this.country.id,
          countryName: this.country.name || this.countryName || '',
          countryFlag: this.country.flag || ((this.country.name || this.countryName || '').slice(0, 1) || '签'),
          visaName: this.activeVisa.name,
          visaType: this.activeVisa.visaType,
          processingTime: this.activeVisa.processingTime,
          fee: this.activeVisa.fee,
          validity: this.activeVisa.validity,
          entries: this.activeVisa.entries,
          description: this.activeVisa.description,
          visaSteps: this.activeVisa.steps || [],
        })
        uni.showToast({ title: '已加入目标', icon: 'none' })
        this.closeVisaDetail()
        uni.reLaunch({ url: '/pages/goals/index' })
      } finally {
        this.creatingGoal = false
      }
    },
  },
}
</script>

<style scoped>
.visas-screen {
  background: var(--uview-bg);
}

.page-top-spacer {
  height: 96px;
}

.hero-banner {
  position: relative;
  height: 220px;
  overflow: hidden;
}

.hero-image {
  width: 100%;
  height: 100%;
}

.hero-mask {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(14, 29, 53, 0.2) 0%, rgba(14, 29, 53, 0.62) 100%);
}

.hero-copy {
  position: absolute;
  left: 16px;
  right: 16px;
  bottom: 16px;
}

.hero-title {
  display: block;
  font-size: 30px;
  font-weight: 800;
  color: #fff;
}

.hero-sub {
  display: block;
  margin-top: 4px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.9);
}

.visas-body {
  margin-top: -14px;
  position: relative;
  z-index: 2;
  border-radius: 18px 18px 0 0;
  background: var(--uview-bg);
  padding: 18px 16px 32px;
}

.visa-stack {
  display: block;
  margin-top: 12px;
}

.visa-card {
  padding: 14px;
  margin-bottom: 12px;
}

.visa-card--empty {
  text-align: center;
}

.visa-title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.visa-title-wrap {
  display: flex;
  align-items: center;
}

.visa-title-wrap :deep(.u-tag) {
  margin-left: 8px;
}

.visa-title {
  font-size: 16px;
  font-weight: 800;
  color: var(--uview-text);
}

.visa-desc {
  display: block;
  margin-top: 8px;
  color: var(--uview-text-3);
  font-size: 12px;
  line-height: 1.5;
}

.visa-empty-title {
  display: block;
  margin-top: 10px;
  font-size: 16px;
  font-weight: 700;
}

.visa-empty-desc {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: var(--uview-text-3);
}

.meta-row {
  display: flex;
  flex-wrap: wrap;
  margin: 12px -4px 0;
}

.meta-item {
  width: calc(33.3333% - 8px);
  margin: 0 4px 8px;
}

.meta-label {
  display: block;
  font-size: 11px;
  color: var(--uview-text-3);
}

.meta-value {
  display: block;
  margin-top: 6px;
  font-size: 13px;
  font-weight: 800;
  color: var(--uview-text);
}

.meta-value--price {
  color: var(--uview-brand);
}

.drawer-mask {
  position: fixed;
  inset: 0;
  z-index: 120;
  background: rgba(18, 24, 38, 0.38);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.drawer-sheet {
  width: 100%;
  max-width: 430px;
  height: 72vh;
  min-height: 440px;
  max-height: 680px;
  border-radius: 20px 20px 0 0;
  background: #fff;
  padding: 10px 16px 18px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.drawer-handle {
  width: 42px;
  height: 4px;
  border-radius: 999px;
  background: var(--uview-border);
  margin: 0 auto 16px;
}

.drawer-body {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.drawer-head {
  display: flex;
  align-items: center;
}

.drawer-icon {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  background: rgba(22, 119, 255, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  flex-shrink: 0;
}

.drawer-title-wrap {
  flex: 1;
  min-width: 0;
}

.drawer-title {
  display: block;
  font-size: 17px;
  font-weight: 800;
}

.drawer-sub {
  display: block;
  margin-top: 4px;
  font-size: 11px;
  line-height: 1.45;
  color: var(--uview-text-3);
}

.drawer-close {
  width: 34px;
  height: 34px;
  border-radius: 999px;
  background: var(--uview-surface-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.detail-loading {
  margin-top: 18px;
  padding: 20px 0;
  text-align: center;
}

.detail-loading-text {
  display: block;
  margin-top: 10px;
  color: var(--uview-text-3);
  font-size: 13px;
}

.drawer-scroll {
  flex: 1;
  min-height: 0;
  margin-top: 12px;
}

.drawer-block {
  padding-bottom: 16px;
}

.badge-row {
  display: flex;
  flex-wrap: wrap;
}

.badge-row :deep(.u-tag) {
  margin-right: 8px;
  margin-bottom: 8px;
}

.intro-card {
  margin-top: 12px;
  padding: 11px;
  display: flex;
  align-items: flex-start;
  border-radius: 16px;
  background: var(--uview-surface-soft);
}

.intro-card :deep(.u-icon) {
  margin-right: 8px;
  flex-shrink: 0;
}

.intro-text {
  font-size: 11px;
  line-height: 1.5;
  color: var(--uview-text-3);
}

.detail-grid {
  margin-top: 14px;
  display: flex;
  flex-wrap: wrap;
  margin-left: -4px;
  margin-right: -4px;
}

.detail-box {
  width: calc(50% - 8px);
  margin: 0 4px 8px;
  background: var(--uview-surface-soft);
  border-radius: 14px;
  padding: 10px;
}

.detail-box :deep(.u-icon) {
  margin-bottom: 6px;
}

.detail-box-label {
  display: block;
  font-size: 10px;
  color: var(--uview-text-3);
}

.detail-box-value {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  font-weight: 800;
  line-height: 1.4;
}

.drawer-block-title {
  display: block;
  margin-bottom: 8px;
  font-size: 13px;
  font-weight: 800;
}

.step-card {
  margin-top: 8px;
  padding: 12px;
  border-radius: 16px;
  background: var(--uview-surface-soft);
}

.step-top {
  display: flex;
  align-items: center;
}

.step-index {
  width: 24px;
  height: 24px;
  border-radius: 999px;
  background: rgba(22, 119, 255, 0.14);
  color: var(--uview-brand);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 800;
  margin-right: 8px;
}

.step-title {
  font-size: 14px;
  font-weight: 800;
  color: var(--uview-text);
}

.step-desc {
  display: block;
  margin-top: 8px;
  font-size: 12px;
  line-height: 1.55;
  color: var(--uview-text-3);
}
</style>
