<template>
  <view class="visa-free-page" :class="{ 'visa-free-page--dark': themeMode === 'dark' }">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="免签/落地签国家" />

    <view class="visa-free-content visago-page-width">
      <view class="hero-panel visago-card">
        <view class="hero-panel__search">
          <view class="search-box search-box--hero">
            <text class="material-symbols-outlined search-icon">search</text>
            <input v-model.trim="keyword" class="search-input" placeholder="搜索国家、城市、政策类型" />
            <text v-if="keyword" class="material-symbols-outlined clear-icon" @tap="keyword = ''">close</text>
          </view>
          <view class="hero-panel__hint">
            <text class="material-symbols-outlined hero-panel__hint-icon">travel_explore</text>
            <text>{{ displayCountries.length }} 个结果，支持地图联动与快速筛选</text>
          </view>
        </view>

        <view class="stats-row">
          <view v-for="item in summaryCards" :key="item.key" class="stats-card">
            <view class="stats-icon-wrap" :class="`stats-icon-wrap--${item.key}`">
              <text class="material-symbols-outlined stats-icon">{{ item.icon }}</text>
            </view>
            <text class="stats-label">{{ item.label }}</text>
            <text class="stats-value">{{ item.value }}</text>
          </view>
        </view>
      </view>

      <view class="filters-panel visago-card">
        <scroll-view class="policy-scroll" scroll-x :show-scrollbar="false">
          <view class="policy-row">
            <view
              v-for="item in policyTabs"
              :key="item"
              class="policy-chip"
              :class="{ 'policy-chip--active': activePolicy === item }"
              @tap="activePolicy = item"
            >
              {{ item }}
            </view>
          </view>
        </scroll-view>

        <scroll-view class="continent-scroll" scroll-x :show-scrollbar="false">
          <view class="policy-row">
            <view
              v-for="item in continentTabs"
              :key="item"
              class="policy-chip"
              :class="{ 'policy-chip--active': activeContinent === item }"
              @tap="activeContinent = item"
            >
              {{ item }}
            </view>
          </view>
        </scroll-view>
      </view>

      <view class="section-head">
        <text class="section-title">全球免签地图</text>
        <text class="section-sub">地图和列表联动，筛选后只高亮当前结果，支持拖拽和缩放查看。</text>
      </view>

      <view id="visa-free-map" class="map-card visago-card">
        <view class="map-legend">
          <view class="legend-item"><span class="legend-dot legend-dot--free" />免签证</view>
          <view class="legend-item"><span class="legend-dot legend-dot--voa" />落地签</view>
          <view class="legend-item"><span class="legend-dot legend-dot--eta" />ETA</view>
          <view class="legend-item"><span class="legend-dot legend-dot--evisa" />电子签证</view>
        </view>

        <movable-area class="map-area">
          <movable-view class="map-view" direction="all" :scale="true" :scale-value="mapScale" :scale-min="1" :scale-max="3">
            <image
              class="map-image"
              :class="{ 'map-image--dark': themeMode === 'dark' }"
              src="https://upload.wikimedia.org/wikipedia/commons/thumb/8/80/World_map_-_low_resolution.svg/1280px-World_map_-_low_resolution.svg.png"
              mode="aspectFill"
            />
            <view class="map-mask" />

            <view
              v-for="country in mapCountries"
              :key="country.id"
              class="map-marker"
              :class="[`map-marker--${mapMarkerType(country.policyType)}`, { 'map-marker--active': activeMapCountryId === country.id }]"
              :style="{ left: `${country.mapX}%`, top: `${country.mapY}%` }"
            />

            <view
              v-if="activeMapCountry"
              class="map-label"
              :style="{ left: `${activeMapCountry.mapX}%`, top: `${activeMapCountry.mapY}%` }"
            >
              <text>{{ activeMapCountry.name }} · {{ activeMapCountry.city }}</text>
            </view>
          </movable-view>
        </movable-area>

        <view class="map-zoom">
          <view class="zoom-btn" @tap="changeMapScale(-0.25)">
            <text class="material-symbols-outlined">remove</text>
          </view>
          <view class="zoom-btn" @tap="changeMapScale(0.25)">
            <text class="material-symbols-outlined">add</text>
          </view>
        </view>
      </view>

      <view class="section-head">
        <text class="section-title">国家列表</text>
        <text class="section-sub">{{ displayCountries.length }} 个结果，点定位图标可在地图高亮当前国家标注</text>
      </view>

      <view class="country-list">
        <view
          v-for="country in displayCountries"
          :key="country.id"
          class="country-item visago-card"
          :class="{ 'country-item--active': activeDetailCountry && activeDetailCountry.id === country.id }"
          @tap="openCountryDetail(country)"
        >
          <view class="country-main">
            <view class="country-title-row">
              <text class="country-name">{{ country.name }}</text>
              <text class="country-city">{{ country.city }}</text>
            </view>
            <view class="country-badges">
              <text class="country-badge">{{ country.policyType }}</text>
              <text class="country-badge country-badge--muted">{{ country.stay }}</text>
              <text v-if="country.supportedVisaId" class="country-badge country-badge--success">支持加入目标</text>
            </view>
            <text class="country-note">{{ country.note }}</text>
          </view>
          <view class="country-actions">
            <view class="loc-btn" @tap.stop="focusOnMap(country)">
              <text class="material-symbols-outlined">location_on</text>
            </view>
            <text class="material-symbols-outlined country-arrow">chevron_right</text>
          </view>
        </view>

        <view v-if="!displayCountries.length" class="empty-card visago-card">
          <text class="material-symbols-outlined">travel_explore</text>
          <text>没有匹配的免签国家</text>
        </view>
      </view>
    </view>

    <view class="drawer-mask" :class="{ 'drawer-mask--show': drawerVisible }" @tap="closeCountryDetail">
      <view class="detail-drawer" :class="{ 'detail-drawer--show': drawerVisible }" @tap.stop>
        <view v-if="activeDetailCountry">
          <view class="drawer-head">
            <view class="drawer-title-wrap">
              <text class="drawer-title">{{ activeDetailCountry.name }} {{ activeDetailCountry.policyType }}</text>
              <text class="drawer-sub">{{ activeDetailCountry.stay || '详情信息' }}</text>
            </view>
            <view class="drawer-close" @tap="closeCountryDetail">
              <text class="material-symbols-outlined">close</text>
            </view>
          </view>

          <view class="badge-row">
            <text class="badge">{{ activeDetailCountry.policyType }}</text>
            <text class="badge">{{ activeDetailCountry.stay }}</text>
            <text v-if="activeDetailCountry.supportedVisaId" class="badge badge--success">支持加入目标</text>
          </view>

          <text class="intro-text">{{ activeDetailCountry.note }}</text>

          <view class="support-card" :class="{ 'support-card--disabled': !activeDetailCountry.supportedVisaId }">
            <view>
              <text class="support-title">{{ activeDetailCountry.supportedVisaId ? '已维护，可加入目标' : '暂未维护到签证目标' }}</text>
              <text class="support-sub">
                {{ activeDetailCountry.supportedVisaId ? '这条政策已经关联到签证目标，可直接加入计划继续跟进。' : '当前只提供政策展示，还不能直接加入目标计划。' }}
              </text>
            </view>
            <text class="material-symbols-outlined support-icon">{{ activeDetailCountry.supportedVisaId ? 'verified' : 'info' }}</text>
          </view>

          <view class="info-grid">
            <view class="info-card">
              <text class="material-symbols-outlined info-icon">calendar_month</text>
              <text class="info-label">停留时长</text>
              <text class="info-value">{{ activeDetailCountry.stay || '-' }}</text>
            </view>
            <view class="info-card">
              <text class="material-symbols-outlined info-icon">public</text>
              <text class="info-label">政策类型</text>
              <text class="info-value">{{ activeDetailCountry.policyType || '-' }}</text>
            </view>
            <view class="info-card">
              <text class="material-symbols-outlined info-icon">location_city</text>
              <text class="info-label">代表城市</text>
              <text class="info-value">{{ activeDetailCountry.city || '-' }}</text>
            </view>
            <view class="info-card">
              <text class="material-symbols-outlined info-icon">fact_check</text>
              <text class="info-label">目标支持</text>
              <text class="info-value">{{ activeDetailCountry.supportedVisaId ? `${activeDetailCountry.supportedCountryName || ''} ${activeDetailCountry.supportedVisaName || ''}`.trim() : '暂未维护' }}</text>
            </view>
          </view>

          <view class="process-wrap">
            <text class="process-title">出行前确认</text>
            <view v-for="(item, idx) in currentTips" :key="item" class="process-item">
              <view class="process-dot">{{ idx + 1 }}</view>
              <view class="process-main">
                <text class="process-item-title">{{ item.title }}</text>
                <text class="process-item-desc">{{ item.desc }}</text>
              </view>
            </view>
          </view>

          <view class="drawer-actions">
            <view class="ghost-btn" @tap="focusOnMap(activeDetailCountry)">地图定位</view>
            <view class="primary-btn" :class="{ 'primary-btn--disabled': !activeDetailCountry.supportedVisaId || joiningPlan }" @tap="joinPlan(activeDetailCountry)">
              {{ !activeDetailCountry.supportedVisaId ? '暂未支持' : joiningPlan ? '加入中...' : '加入目标' }}
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../components/VisagoTopBar.vue'
import { api } from '../../../utils/api'
import { applyTheme, getStoredTheme, THEME_CHANGE_EVENT } from '../../../utils/theme'

export default {
  components: { VisagoTopBar },
  data() {
    return {
      themeMode: 'light',
      countries: [],
      keyword: '',
      activePolicy: '全部',
      activeContinent: '全部',
      activeMapCountryId: 0,
      activeDetailCountry: null,
      drawerVisible: false,
      joiningPlan: false,
      mapScale: 1,
    }
  },
  computed: {
    policyTabs() {
      return ['全部', '免签', '落地签', 'ETA', '电子签']
    },
    continentTabs() {
      const list = Array.from(new Set(this.countries.map((item) => item.region).filter(Boolean)))
      return ['全部', ...list]
    },
    displayCountries() {
      const keyword = String(this.keyword || '').trim().toLowerCase()
      return this.countries.filter((item) => {
        const passPolicy = this.activePolicy === '全部' || item.policyType === this.activePolicy
        if (!passPolicy) return false
        const passContinent = this.activeContinent === '全部' || item.region === this.activeContinent
        if (!passContinent) return false
        if (!keyword) return true
        const pool = [item.name, item.code, item.city, item.region, item.policyType, item.stay, item.note, ...(item.keywords || [])]
        return pool.some((entry) => String(entry || '').toLowerCase().includes(keyword))
      })
    },
    activeMapCountry() {
      return this.displayCountries.find((item) => item.id === this.activeMapCountryId) || this.displayCountries[0] || null
    },
    mapCountries() {
      const seen = new Set()
      return this.displayCountries.filter((item) => {
        const key = `${item.name}-${item.city}`
        if (seen.has(key)) return false
        seen.add(key)
        return true
      })
    },
    summaryCards() {
      const freeCount = this.countries.filter((item) => item.policyType.includes('免签')).length
      const arrivalCount = this.countries.filter((item) => item.policyType.includes('落地签')).length
      const etaCount = this.countries.filter((item) => item.policyType.includes('ETA')).length
      const eVisaCount = this.countries.filter((item) => item.policyType.includes('电子签')).length
      return [
        { key: 'free', label: '免签', value: freeCount, icon: 'verified' },
        { key: 'arrival', label: '落地签', value: arrivalCount, icon: 'flight_land' },
        { key: 'eta', label: 'ETA', value: etaCount, icon: 'approval_delegation' },
        { key: 'evisa', label: '电子签', value: eVisaCount, icon: 'travel' },
      ]
    },
    currentTips() {
      const policyType = this.activeDetailCountry?.policyType || ''
      const common = [
        { title: '确认最新政策', desc: '出行前请再次核对航司、使领馆或官方入境政策页面，避免政策更新影响行程。' },
        { title: '准备基础材料', desc: '建议提前准备护照、返程机票、酒店订单和资金证明，以便入境查验。' },
      ]
      if (policyType.includes('落地签')) {
        return [
          { title: '预留办理时间', desc: '落地签通常需要在抵达口岸现场办理，建议预留充足的转场和排队时间。' },
          { title: '准备现金或支付方式', desc: '部分口岸仍要求现场缴费，建议提前准备当地要求的支付方式。' },
          common[1],
        ]
      }
      if (policyType.includes('ETA')) {
        return [
          { title: '提前申请 ETA', desc: 'ETA 通常需要在出发前完成在线授权，建议保存好审批结果与回执。' },
          common[0],
          common[1],
        ]
      }
      if (policyType.includes('电子签')) {
        return [
          { title: '提前完成电子签申请', desc: '电子签证通常需要在线提交材料，请预留审批时间并保留签发文件。' },
          common[0],
          common[1],
        ]
      }
      return common
    },
  },
  watch: {
    keyword() {
      this.ensureActiveCountry()
    },
    activePolicy() {
      this.ensureActiveCountry()
    },
    activeContinent() {
      this.ensureActiveCountry()
    },
  },
  onLoad() {
    if (typeof uni !== 'undefined' && uni.$on) {
      uni.$on(THEME_CHANGE_EVENT, this.onThemeChange)
    }
    this.syncTheme()
  },
  onShow() {
    this.syncTheme()
    this.loadFreeCountries()
  },
  onUnload() {
    if (typeof uni !== 'undefined' && uni.$off) {
      uni.$off(THEME_CHANGE_EVENT, this.onThemeChange)
    }
  },
  methods: {
    async loadFreeCountries() {
      try {
        const countries = (await api.listFreeCountries({ enabled: 1 })) || []
        this.countries = countries
        this.ensureActiveCountry()
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '加载失败',
          icon: 'none',
        })
      }
    },
    syncTheme() {
      this.themeMode = applyTheme(getStoredTheme())
    },
    onThemeChange(theme) {
      this.themeMode = theme === 'dark' ? 'dark' : 'light'
      applyTheme(this.themeMode)
    },
    focusOnMap(country) {
      this.activeMapCountryId = country.id
      uni.pageScrollTo({
        selector: '#visa-free-map',
        duration: 260,
      })
    },
    openCountryDetail(country) {
      this.activeDetailCountry = country
      this.activeMapCountryId = country.id
      this.drawerVisible = true
    },
    changeMapScale(delta) {
      const next = Math.max(1, Math.min(3, Number(this.mapScale || 1) + delta))
      this.mapScale = Number(next.toFixed(2))
    },
    mapMarkerType(policyType) {
      if (policyType === '落地签') return 'voa'
      if (policyType === 'ETA') return 'eta'
      if (policyType === '电子签') return 'evisa'
      return 'free'
    },
    ensureActiveCountry() {
      const exists = this.displayCountries.some((item) => item.id === this.activeMapCountryId)
      if (!exists) {
        this.activeMapCountryId = this.displayCountries[0]?.id || 0
      }
    },
    closeCountryDetail() {
      this.drawerVisible = false
      setTimeout(() => {
        this.activeDetailCountry = null
      }, 180)
    },
    joinPlan(country) {
      if (!country || this.joiningPlan) return
      if (!country.supportedVisaId) {
        uni.showToast({
          title: '当前免签国家暂未关联可加入的目标',
          icon: 'none',
        })
        return
      }
      this.createFreeVisaPlan(country)
    },
    async createFreeVisaPlan(country) {
      this.joiningPlan = true
      try {
        const currentPlans = await api.listPlans()
        const freeVisa = {
          id: country.supportedVisaId,
          name: country.supportedVisaName,
        }
        if (!freeVisa) {
          uni.showToast({
            title: '该国家未配置免签目标',
            icon: 'none',
          })
          return
        }
        const existed = (currentPlans || []).find((item) => Number(item.visaId) === Number(freeVisa.id) && item.status === 'active')
        if (existed) {
          this.closeCountryDetail()
          uni.navigateTo({
            url: `/pages/visago/plan/index?planId=${existed.id}`,
          })
          return
        }
        const plan = await api.createPlan({
          countryId: country.supportedCountryId,
          countryName: country.supportedCountryName || country.name,
          visaId: freeVisa.id,
          visaTitle: freeVisa.name || country.supportedVisaName,
          source: 'visa-free',
        })
        this.closeCountryDetail()
        uni.navigateTo({
          url: `/pages/visago/plan/index?planId=${plan.id}`,
        })
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '加入目标失败',
          icon: 'none',
        })
      } finally {
        this.joiningPlan = false
      }
    },
  },
}
</script>

<style scoped>
.visa-free-page {
  --vf-bg: #f3f2f8;
  --vf-surface: #ffffff;
  --vf-line: #dce2ee;
  --vf-text: #101827;
  --vf-text-soft: #69758f;
  --vf-text-muted: #8a95ac;
  --vf-map-mask: rgba(10, 27, 50, 0.12);
  --vf-marker: #0f65d8;
  --vf-popup-mask: rgba(7, 15, 30, 0.4);
  min-height: 100vh;
  background: var(--vf-bg);
}

.visa-free-content {
  box-sizing: border-box;
  padding: 96px 16px 20px;
}

.hero-panel {
  position: relative;
  overflow: hidden;
  padding: 14px 14px 14px;
  border: 1px solid color-mix(in srgb, var(--vf-marker) 16%, var(--vf-line));
  background:
    radial-gradient(circle at top right, rgba(15, 101, 216, 0.16) 0%, rgba(15, 101, 216, 0) 34%),
    linear-gradient(180deg, color-mix(in srgb, var(--vf-surface) 92%, #ffffff 8%) 0%, var(--vf-surface) 100%);
}

.hero-panel__search,
.hero-panel__hint,
.stats-head {
  position: relative;
  z-index: 1;
}

.hero-panel__hint,
.stats-label,
.stats-value,
.stats-sub {
  display: block;
}

.search-box--hero {
  border-color: transparent;
  background: rgba(255, 255, 255, 0.78);
  box-shadow: inset 0 0 0 1px rgba(15, 101, 216, 0.08);
}

.hero-panel__hint {
  margin-top: 10px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--vf-text-soft);
}

.hero-panel__hint-icon {
  font-size: 16px;
  color: var(--vf-marker);
}

.search-box {
  height: 48px;
  border-radius: 14px;
  padding: 0 14px;
  background: var(--vf-surface);
  border: 1px solid var(--vf-line);
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-icon,
.clear-icon {
  color: var(--vf-text-soft);
  font-size: 18px;
}

.search-input {
  flex: 1;
  height: 48px;
  font-size: 14px;
  color: var(--vf-text);
}

.stats-row {
  margin-top: 14px;
  display: flex;
  gap: 10px;
}

.stats-card {
  flex: 1;
  min-width: 0;
  padding: 14px 12px;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.42);
  background: rgba(255, 255, 255, 0.74);
  backdrop-filter: blur(10px);
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.stats-label,
.stats-value,
.country-title-row,
.country-badges,
.support-title,
.support-sub {
  display: block;
}

.stats-label {
  margin-top: 8px;
  font-size: 12px;
  color: var(--vf-text-soft);
}

.stats-icon-wrap {
  width: 28px;
  height: 28px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stats-icon-wrap--free {
  background: rgba(34, 197, 94, 0.14);
  color: #16a34a;
}

.stats-icon-wrap--arrival {
  background: rgba(59, 130, 246, 0.14);
  color: #2563eb;
}

.stats-icon-wrap--eta {
  background: rgba(168, 85, 247, 0.14);
  color: #9333ea;
}

.stats-icon-wrap--evisa {
  background: rgba(245, 158, 11, 0.16);
  color: #d97706;
}

.stats-icon {
  font-size: 16px;
}

.stats-value {
  margin-top: 10px;
  font-size: 24px;
  font-weight: 800;
  color: var(--vf-text);
}

.filters-panel {
  margin-top: 14px;
  padding: 12px;
  display: grid;
  gap: 12px;
  border: 1px solid var(--vf-line);
  background: var(--vf-surface);
}

.policy-scroll {
  white-space: nowrap;
  width: 100%;
}

.continent-scroll {
  white-space: nowrap;
  width: 100%;
}

.policy-row {
  display: inline-flex;
  gap: 8px;
  padding-right: 10px;
}

.policy-chip {
  padding: 8px 14px;
  border-radius: 9999px;
  background: var(--vf-surface);
  border: 1px solid var(--vf-line);
  color: var(--vf-text-muted);
  font-size: 13px;
}

.policy-chip--active {
  background: var(--vf-marker);
  color: #fff;
  border-color: var(--vf-marker);
}

.map-card {
  margin-top: 14px;
  position: relative;
  height: 280px;
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid var(--vf-line);
  background: var(--vf-surface);
}

.map-legend {
  position: absolute;
  left: 10px;
  right: 10px;
  top: 10px;
  z-index: 3;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 12px;
  background: rgba(15, 23, 42, 0.78);
  color: #fff;
  font-size: 11px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 9999px;
  display: inline-block;
}

.legend-dot--free {
  background: #22c55e;
}

.legend-dot--voa {
  background: #3b82f6;
}

.legend-dot--eta {
  background: #a855f7;
}

.legend-dot--evisa {
  background: #f59e0b;
}

.map-area {
  width: 100%;
  height: 100%;
}

.map-view {
  width: 100%;
  height: 100%;
}

.map-image {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.map-image--dark {
  filter: brightness(0.96) contrast(1.16) saturate(0.85);
}

.map-mask {
  position: absolute;
  inset: 0;
  background: var(--vf-map-mask);
}

.map-marker {
  position: absolute;
  width: 9px;
  height: 9px;
  border-radius: 9999px;
  background: var(--vf-marker);
  border: 2px solid rgba(255, 255, 255, 0.95);
  transform: translate(-50%, -50%);
  box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.18);
}

.map-marker--free {
  background: #22c55e;
}

.map-marker--voa {
  background: #3b82f6;
}

.map-marker--eta {
  background: #a855f7;
}

.map-marker--evisa {
  background: #f59e0b;
}

.map-marker--active {
  width: 12px;
  height: 12px;
  box-shadow: 0 0 0 5px rgba(15, 101, 216, 0.28);
}

.map-label {
  position: absolute;
  transform: translate(-50%, calc(-100% - 12px));
  max-width: 180px;
  border-radius: 9999px;
  background: rgba(15, 101, 216, 0.94);
  color: #fff;
  font-size: 11px;
  line-height: 1.2;
  padding: 5px 10px;
  white-space: nowrap;
  z-index: 2;
}

.map-zoom {
  position: absolute;
  right: 10px;
  bottom: 10px;
  z-index: 3;
  display: grid;
  gap: 8px;
}

.zoom-btn {
  width: 34px;
  height: 34px;
  border-radius: 9999px;
  background: rgba(15, 23, 42, 0.78);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.section-head {
  margin-top: 16px;
}

.section-title {
  display: block;
  font-size: 18px;
  font-weight: 700;
  color: var(--vf-text);
}

.section-sub {
  display: block;
  margin-top: 3px;
  font-size: 12px;
  color: var(--vf-text-soft);
}

.country-list {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding-bottom: calc(8px + var(--visago-safe-bottom));
}

.country-item {
  border: 1px solid var(--vf-line);
  background: var(--vf-surface);
  padding: 12px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.country-item--active {
  border-color: rgba(15, 101, 216, 0.4);
  box-shadow: inset 0 0 0 1px rgba(15, 101, 216, 0.18);
}

.country-main {
  flex: 1;
  min-width: 0;
}

.country-name {
  display: block;
  font-size: 16px;
  font-weight: 700;
  color: var(--vf-text);
}

.country-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.country-city {
  font-size: 12px;
  color: var(--vf-text-soft);
  flex-shrink: 0;
}

.country-meta {
  display: block;
  margin-top: 3px;
  font-size: 12px;
  color: var(--vf-marker);
}

.country-badges {
  margin-top: 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.country-badge {
  padding: 3px 8px;
  border-radius: 9999px;
  background: rgba(15, 101, 216, 0.12);
  color: var(--vf-marker);
  font-size: 11px;
  border: 1px solid rgba(15, 101, 216, 0.16);
}

.country-badge--muted {
  background: var(--vf-bg);
  color: var(--vf-text-soft);
  border-color: var(--vf-line);
}

.country-badge--success {
  color: #15803d;
  background: rgba(34, 197, 94, 0.12);
  border-color: rgba(21, 128, 61, 0.2);
}

.country-note {
  display: block;
  margin-top: 4px;
  font-size: 11px;
  line-height: 1.4;
  color: var(--vf-text-soft);
}

.country-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.loc-btn {
  width: 30px;
  height: 30px;
  border-radius: 9999px;
  background: rgba(15, 101, 216, 0.13);
  color: var(--vf-marker);
  display: flex;
  align-items: center;
  justify-content: center;
}

.loc-btn .material-symbols-outlined {
  font-size: 18px;
}

.country-arrow {
  color: var(--vf-text-muted);
  font-size: 20px;
}

.drawer-mask {
  position: fixed;
  inset: 0;
  z-index: 99;
  background: transparent;
  pointer-events: none;
  transition: background 0.22s ease;
}

.drawer-mask--show {
  background: var(--vf-popup-mask);
  pointer-events: auto;
}

.detail-drawer {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 20px 20px 0 0;
  background: var(--vf-surface);
  border-top: 1px solid var(--vf-line);
  transform: translateY(105%);
  transition: transform 0.22s ease;
  max-height: 84vh;
  overflow-y: auto;
  padding: 14px 16px calc(16px + var(--visago-safe-bottom));
  box-sizing: border-box;
}

.detail-drawer--show {
  transform: translateY(0);
}

.drawer-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.drawer-title-wrap {
  min-width: 0;
}

.drawer-title {
  display: block;
  font-size: 18px;
  font-weight: 700;
  color: var(--vf-text);
}

.drawer-sub {
  margin-top: 3px;
  display: block;
  font-size: 11px;
  color: var(--vf-text-soft);
}

.drawer-close {
  width: 30px;
  height: 30px;
  border-radius: 9999px;
  background: var(--vf-bg);
  color: var(--vf-text-soft);
  display: flex;
  align-items: center;
  justify-content: center;
}

.badge-row {
  margin-top: 10px;
  display: flex;
  gap: 8px;
}

.badge {
  padding: 3px 10px;
  border-radius: 8px;
  font-size: 11px;
  color: var(--vf-marker);
  background: rgba(15, 101, 216, 0.12);
  border: 1px solid rgba(15, 101, 216, 0.22);
}

.badge--success {
  color: #15803d;
  background: rgba(34, 197, 94, 0.12);
  border-color: rgba(21, 128, 61, 0.2);
}

.support-card {
  margin-top: 12px;
  padding: 12px;
  border-radius: 14px;
  border: 1px solid rgba(15, 101, 216, 0.18);
  background: rgba(15, 101, 216, 0.08);
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.support-card--disabled {
  border-color: var(--vf-line);
  background: var(--vf-bg);
}

.support-title {
  font-size: 14px;
  font-weight: 700;
  color: var(--vf-text);
}

.support-sub {
  margin-top: 5px;
  font-size: 12px;
  line-height: 1.5;
  color: var(--vf-text-soft);
}

.support-icon {
  font-size: 20px;
  color: var(--vf-marker);
}

.intro-text {
  margin-top: 10px;
  display: block;
  font-size: 13px;
  line-height: 1.5;
  color: var(--vf-text-soft);
}

.info-grid {
  margin-top: 12px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.info-card {
  border-radius: 12px;
  border: 1px solid var(--vf-line);
  background: var(--vf-bg);
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-icon {
  font-size: 20px;
  color: var(--vf-marker);
}

.info-label {
  font-size: 11px;
  color: var(--vf-text-muted);
}

.info-value {
  font-size: 12px;
  color: var(--vf-text);
  font-weight: 600;
  line-height: 1.35;
}

.process-wrap {
  margin-top: 16px;
}

.process-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--vf-text);
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
  color: var(--vf-marker);
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
  color: var(--vf-text);
  font-weight: 600;
}

.process-item-desc {
  font-size: 12px;
  line-height: 1.45;
  color: var(--vf-text-soft);
}

.drawer-actions {
  margin-top: 14px;
  display: flex;
  gap: 8px;
}

.ghost-btn,
.primary-btn {
  flex: 1;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
}

.ghost-btn {
  background: rgba(15, 101, 216, 0.1);
  color: var(--vf-marker);
}

.primary-btn {
  background: var(--vf-marker);
  color: #fff;
}

.primary-btn--disabled {
  opacity: 0.5;
}

.empty-card {
  min-height: 88px;
  padding: 16px;
  color: var(--vf-text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  text-align: center;
}

.visa-free-page--dark,
:global(html.theme-dark) .visa-free-page {
  --vf-bg: #0f1218;
  --vf-surface: #171b24;
  --vf-line: #4d5669;
  --vf-text: #eef1f8;
  --vf-text-soft: #b6bfd3;
  --vf-text-muted: #8993ac;
  --vf-map-mask: rgba(108, 153, 220, 0.08);
  --vf-marker: #60a5fa;
  --vf-popup-mask: rgba(0, 0, 0, 0.56);
}

@media (max-width: 420px) {
  .stats-card {
    padding: 12px 10px;
  }

  .stats-label {
    font-size: 11px;
  }

  .stats-value {
    font-size: 20px;
  }

  .stats-sub {
    font-size: 10px;
  }
}
</style>

