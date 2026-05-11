<template>
  <view class="page-screen visas-screen">
    <VuPageHeader :title="country.name || countryName || '签证列表'" :show-back="true" />
    <view class="page-top-spacer" :style="topSpacerStyle" />

    <view class="hero">
      <image class="hero-image" :src="country.image || defaultHero" mode="aspectFill" />
      <view class="hero-mask" />
      <view class="hero-copy">
        <text class="hero-title">{{ country.name || countryName || '国家签证' }}</text>
        <text class="hero-sub">{{ country.region || '全球' }} · 可选签证类型</text>
      </view>
    </view>

    <view class="page-width page-body visas-body">
      <text class="app-section-title">可选签证类型</text>

      <view class="visa-list">
        <view v-if="loading" class="app-card app-empty visa-empty">
          <u-loading-icon mode="circle" size="22" color="#1677ff" />
          <text class="app-empty-title">正在加载签证列表</text>
          <text class="app-empty-desc">请稍等，正在读取该国家下的签证信息。</text>
        </view>

        <view
          v-for="visa in visaOptions"
          :key="visa.id"
          class="app-card visa-card"
          hover-class="visa-card--hover"
          @tap="openVisaDetail(visa)"
        >
          <view class="visa-title-row">
            <view class="visa-title-main">
              <text class="visa-title">{{ visa.name }}</text>
              <u-tag v-if="visa.visaFree" text="免签" type="warning" size="mini" />
            </view>
            <u-tag v-if="visa.hot" text="热门" type="primary" size="mini" />
          </view>

          <text class="visa-desc">{{ visa.description || visa.longIntro || '查看签证详情与办理信息' }}</text>

          <view class="visa-meta-row">
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

        <view v-if="!loading && !visaOptions.length" class="app-card app-empty visa-empty">
          <u-icon name="empty-list" size="30" color="#9aa4b5" />
          <text class="app-empty-title">还没有签证配置</text>
          <text class="app-empty-desc">这个国家暂时还没有可展示的签证信息。</text>
        </view>
      </view>
    </view>

    <view v-if="activeVisa" class="detail-mask" @tap="closeVisaDetail">
      <view class="detail-sheet" @tap.stop>
        <view class="detail-handle" />
        <view class="detail-head">
          <view class="detail-head-icon">
            <u-icon :name="activeVisa.visaFree ? 'map' : 'bookmark'" size="20" color="#1677ff" />
          </view>
          <view class="detail-title-wrap">
            <text class="detail-title">{{ activeVisa.name }}</text>
            <text class="detail-sub">{{ activeVisa.description || activeVisa.longIntro || '查看签证详情与办理信息' }}</text>
          </view>
          <view class="detail-close" @tap="closeVisaDetail">
            <u-icon name="close" size="16" color="#7b8495" />
          </view>
        </view>

        <view v-if="detailLoading" class="detail-loading">
          <u-loading-icon mode="circle" size="22" color="#1677ff" />
          <text class="detail-loading-text">正在加载签证详情...</text>
        </view>

        <scroll-view v-else scroll-y class="detail-scroll">
          <view class="badge-row">
            <view class="badge-cell">
              <u-tag :text="activeVisa.visaType || '签证类型'" type="primary" size="mini" />
            </view>
            <view v-if="activeVisa.visaFree" class="badge-cell">
              <u-tag text="免签" type="warning" size="mini" />
            </view>
            <view class="badge-cell">
              <u-tag :text="activeVisa.entries || '入境次数待定'" type="info" size="mini" />
            </view>
          </view>

          <view class="intro-card">
            <u-icon name="file-text" size="16" color="#1677ff" />
            <text class="intro-text">{{ activeVisa.longIntro || activeVisa.description || '暂无详情说明' }}</text>
          </view>

          <view class="detail-grid">
            <view class="detail-box">
              <u-icon name="clock" size="17" color="#1677ff" />
              <text class="detail-box-label">办理周期</text>
              <text class="detail-box-value">{{ activeVisa.processingTime || '-' }}</text>
            </view>
            <view class="detail-box">
              <u-icon name="rmb" size="17" color="#1677ff" />
              <text class="detail-box-label">预估费用</text>
              <text class="detail-box-value">{{ activeVisa.fee || '-' }}</text>
            </view>
            <view class="detail-box">
              <u-icon name="bookmark" size="17" color="#1677ff" />
              <text class="detail-box-label">有效期</text>
              <text class="detail-box-value">{{ activeVisa.validity || '-' }}</text>
            </view>
            <view class="detail-box">
              <u-icon name="calendar" size="17" color="#1677ff" />
              <text class="detail-box-label">入境次数</text>
              <text class="detail-box-value">{{ activeVisa.entries || '-' }}</text>
            </view>
          </view>

          <view v-if="activeVisa.steps && activeVisa.steps.length" class="process-block">
            <text class="process-title">办理步骤</text>
            <view
              v-for="(step, idx) in activeVisa.steps"
              :key="step.id || step.stepKey || idx"
              class="process-item"
            >
              <view class="process-index">{{ idx + 1 }}</view>
              <view class="process-main">
                <text class="process-item-title">{{ step.title }}</text>
                <text class="process-desc">{{ step.desc || (step.strategies || []).join('；') || '查看详情后再准备下一步。' }}</text>
              </view>
            </view>
          </view>

          <u-button
            :text="creatingGoal ? '正在加入目标...' : goalExists ? '去查看我的目标' : '加入我的目标'"
            type="primary"
            shape="circle"
            :disabled="creatingGoal"
            customStyle="margin-top: 16px;"
            @click="addToGoal"
          />
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script>
import VuPageHeader from '../../components/uview/VuPageHeader.vue'
import { api } from '../../utils/api'
import { getGoalByVisaId, upsertGoal } from '../../utils/goals'
import { getTopSpacerStyle } from '../../utils/layout'

const defaultHero = 'https://images.pexels.com/photos/918275/pexels-photo-918275.jpeg?auto=compress&cs=tinysrgb&w=1200'

export default {
  components: {
    VuPageHeader,
  },
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
      topSpacerStyle: 'height:96px;',
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
  background: #f7f8fb;
}

.page-top-spacer {
  height: 96px;
}

.hero {
  position: relative;
  height: 218px;
  overflow: hidden;
}

.hero-image {
  width: 100%;
  height: 218px;
}

.hero-mask {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  background: linear-gradient(180deg, rgba(14, 29, 53, 0.18) 0%, rgba(14, 29, 53, 0.64) 100%);
}

.hero-copy {
  position: absolute;
  right: 16px;
  bottom: 16px;
  left: 16px;
}

.hero-title {
  display: block;
  color: #ffffff;
  font-size: 30px;
  line-height: 1.2;
  font-weight: 900;
}

.hero-sub {
  display: block;
  margin-top: 5px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 13px;
  line-height: 1.4;
}

.visas-body {
  position: relative;
  z-index: 2;
  margin-top: -14px;
  padding-top: 18px;
  padding-bottom: 30px;
  border-radius: 18px 18px 0 0;
  background: #f7f8fb;
}

.visa-list {
  margin-top: 12px;
}

.visa-card {
  margin-bottom: 12px;
  padding: 14px;
}

.visa-card--hover {
  background: #f8fbff;
}

.visa-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.visa-title-main {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
}

.visa-title-main :deep(.u-tag) {
  margin-left: 8px;
  flex-shrink: 0;
}

.visa-title {
  display: block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  color: #1f2329;
  font-size: 16px;
  line-height: 1.3;
  font-weight: 800;
}

.visa-desc {
  display: block;
  margin-top: 8px;
  color: #7b8495;
  font-size: 12px;
  line-height: 1.55;
}

.visa-meta-row {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #eef2f7;
  display: flex;
}

.meta-item {
  width: 33.3333%;
  padding-right: 8px;
}

.meta-label {
  display: block;
  color: #8a94a6;
  font-size: 11px;
  line-height: 1.2;
}

.meta-value {
  display: block;
  margin-top: 6px;
  color: #1f2329;
  font-size: 13px;
  line-height: 1.35;
  font-weight: 800;
}

.meta-value--price {
  color: #1677ff;
}

.visa-empty {
  margin-bottom: 12px;
}

.detail-mask {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 120;
  background: rgba(18, 24, 38, 0.42);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.detail-sheet {
  width: 100%;
  max-width: 430px;
  height: 72vh;
  min-height: 430px;
  max-height: 680px;
  padding: 10px 16px 18px;
  border-radius: 22px 22px 0 0;
  background: #ffffff;
  display: flex;
  flex-direction: column;
}

.detail-handle {
  width: 42px;
  height: 4px;
  margin: 0 auto 16px;
  border-radius: 2px;
  background: #e4e9f2;
  flex-shrink: 0;
}

.detail-head {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.detail-head-icon {
  width: 44px;
  height: 44px;
  margin-right: 12px;
  border-radius: 14px;
  background: #eef6ff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.detail-title-wrap {
  flex: 1;
  min-width: 0;
}

.detail-title {
  display: block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  color: #1f2329;
  font-size: 17px;
  font-weight: 900;
}

.detail-sub {
  display: block;
  margin-top: 4px;
  color: #7b8495;
  font-size: 11px;
  line-height: 1.45;
}

.detail-close {
  width: 34px;
  height: 34px;
  margin-left: 8px;
  border-radius: 17px;
  background: #f2f4f8;
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
  color: #7b8495;
  font-size: 13px;
}

.detail-scroll {
  flex: 1;
  min-height: 0;
  margin-top: 12px;
}

.badge-row {
  display: flex;
  flex-wrap: wrap;
}

.badge-cell {
  margin-right: 8px;
  margin-bottom: 8px;
}

.intro-card {
  margin-top: 8px;
  padding: 12px;
  border-radius: 16px;
  background: #f4f7fb;
  display: flex;
  align-items: flex-start;
}

.intro-card :deep(.u-icon) {
  margin-right: 8px;
  flex-shrink: 0;
}

.intro-text {
  flex: 1;
  color: #6b7280;
  font-size: 12px;
  line-height: 1.6;
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
  margin-left: 4px;
  margin-right: 4px;
  margin-bottom: 8px;
  padding: 12px;
  border-radius: 16px;
  background: #f4f7fb;
}

.detail-box-label {
  display: block;
  margin-top: 8px;
  color: #8a94a6;
  font-size: 11px;
}

.detail-box-value {
  display: block;
  margin-top: 6px;
  color: #1f2329;
  font-size: 13px;
  line-height: 1.35;
  font-weight: 800;
}

.process-block {
  margin-top: 10px;
}

.process-title {
  display: block;
  font-size: 15px;
  font-weight: 900;
  color: #1f2329;
}

.process-item {
  margin-top: 10px;
  display: flex;
}

.process-index {
  width: 26px;
  height: 26px;
  margin-right: 10px;
  border-radius: 13px;
  background: #eef6ff;
  color: #1677ff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 900;
  flex-shrink: 0;
}

.process-main {
  flex: 1;
  min-width: 0;
}

.process-item-title {
  display: block;
  color: #1f2329;
  font-size: 14px;
  font-weight: 800;
}

.process-desc {
  display: block;
  margin-top: 5px;
  color: #7b8495;
  font-size: 12px;
  line-height: 1.55;
}
</style>
