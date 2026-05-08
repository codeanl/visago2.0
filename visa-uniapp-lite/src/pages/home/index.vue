<template>
  <view class="screen screen--home">
    <VisagoTopBar :home="true" :show-notice="false" logo-src="/static/header-logo.png" />

    <view class="page-wrap home-stack">
      <view class="hero-section">
        <text class="hero-title">想去哪里探索世界？</text>
        <view class="search-strip">
          <input
            v-model.trim="keyword"
            class="search-input"
            placeholder="搜索国家、地区或签证类型"
            confirm-type="search"
            @confirm="searchCountries"
          />
          <view class="search-btn" @tap="searchCountries">搜索</view>
        </view>
      </view>

      <view class="visa-type-wrap">
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

      <view class="action-grid">
        <view class="card action-card action-card--primary" @tap="openCountries">
          <text class="action-emoji">🌍</text>
          <text class="action-title">查国家签证</text>
        </view>
        <view class="card action-card" @tap="openGoals">
          <text class="action-emoji">🎯</text>
          <text class="action-title">我的目标</text>
        </view>
      </view>

      <view class="section-header">
        <text class="section-title">热门目的地</text>
        <text class="section-link"></text>
      </view>

      <scroll-view class="destination-scroll" scroll-x>
        <view class="destination-row">
          <view v-for="item in hotDestinations" :key="item.visaId" class="destination-card visago-card" @tap="openHotDetail(item)">
            <view class="destination-image-wrap">
              <image v-if="item.image" class="destination-image" :src="item.image" mode="aspectFill" />
              <view v-else class="destination-image destination-image--fallback">{{ item.flag || '签' }}</view>
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
                  <text class="destination-price">{{ item.price || '待更新' }}</text>
                  <text class="price-suffix">起</text>
                </view>
              </view>
            </view>
          </view>

          <view v-if="loadingHot" class="empty-tip visago-card">
            正在加载热门目的地...
          </view>

          <view v-else-if="!hotDestinations.length" class="empty-tip visago-card">
            暂无热门目的地
          </view>
        </view>
      </scroll-view>

      <view class="section-block">
        <view class="section-head">
          <text class="section-title">先选国家</text>
          <text class="link-text" @tap="openCountries">查看全部</text>
        </view>
        <view class="country-grid">
          <view v-for="country in countryPreview" :key="country.id" class="card country-card" @tap="openCountry(country)">
            <view class="country-top">
              <text class="country-flag">{{ country.flag || '签' }}</text>
              <text class="country-name">{{ country.name }}</text>
            </view>
            <text class="country-region">{{ country.region || '全球签证' }}</text>
          </view>
        </view>
      </view>
    </view>

    <bottom-dock active="home" />
  </view>
</template>

<script>
import VisagoTopBar from '../../components/VisagoTopBar.vue'
import BottomDock from '../../components/BottomDock.vue'
import { api } from '../../utils/api'

export default {
  components: {
    VisagoTopBar,
    BottomDock,
  },
  data() {
    return {
      keyword: '',
      activeVisaType: 'tourism',
      loadingHot: false,
      hotDestinations: [],
      countryPreview: [],
      visaTypeMenu: [
        { key: 'tourism', label: '旅游签', icon: 'flight', color: '#2f89ff' },
        { key: 'business', label: '商务签', icon: 'business_center', color: '#9333ea' },
        { key: 'study', label: '学生签', icon: 'school', color: '#f59e0b' },
        { key: 'work', label: '工作签', icon: 'description', color: '#22c55e' },
      ],
    }
  },
  onShow() {
    this.loadHomeData()
  },
  methods: {
    async loadHomeData() {
      this.loadingHot = true
      try {
        const [hot, countries] = await Promise.all([
          api.listHotDestinations({ limit: 6 }),
          api.listCountries('', ''),
        ])
        this.hotDestinations = Array.isArray(hot) ? hot : []
        this.countryPreview = Array.isArray(countries) ? countries.slice(0, 6) : []
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '加载失败',
          icon: 'none',
        })
      } finally {
        this.loadingHot = false
      }
    },
    openCountries() {
      uni.navigateTo({ url: '/pages/countries/index' })
    },
    searchCountries(keywordOverride = '') {
      const q = encodeURIComponent(keywordOverride || this.keyword || '')
      uni.navigateTo({ url: `/pages/countries/index?q=${q}` })
    },
    onVisaTypeTap(item) {
      this.activeVisaType = item.key
      this.openCountries()
    },
    openGoals() {
      uni.reLaunch({ url: '/pages/goals/index' })
    },
    openCountry(country) {
      uni.navigateTo({
        url: `/pages/visas/index?countryId=${country.id}&countryName=${encodeURIComponent(country.name || '')}`,
      })
    },
    openHotDetail(item) {
      if (!item || !item.countryId) {
        return
      }
      uni.navigateTo({
        url: `/pages/visas/index?countryId=${item.countryId}&countryName=${encodeURIComponent(item.name || '')}`,
      })
    },
  },
}
</script>

<style scoped>
.screen--home {
  background: var(--visago-bg);
}

.home-stack {
  padding-top: 96px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.hero-section {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.hero-title {
  font-size: 26px;
  line-height: 1.25;
  font-weight: 800;
  color: var(--lite-text);
}

.visa-type-wrap {
  margin-top: -2px;
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
  background: #ffffff;
  border: 1px solid #e7edf8;
}

.visa-type-icon-box--active {
  background: #f8fbff;
  border-color: rgba(15, 101, 216, 0.28);
  transform: translateY(-1px);
}

.visa-type-icon {
  font-size: 26px;
}

.visa-type-label {
  font-size: 13px;
  line-height: 1;
  font-weight: 500;
  color: #65708a;
}

.visa-type-label--active {
  color: #1d4fa8;
  font-weight: 600;
}

.search-strip {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border-radius: 18px;
  background: #ffffff;
  border: 1px solid var(--lite-line);
  box-shadow: var(--lite-shadow);
}

.search-input {
  flex: 1;
  height: 42px;
  padding: 0 14px;
  border-radius: 12px;
  background: #fff;
  color: #173140;
  font-size: 14px;
}

.search-btn {
  min-width: 78px;
  height: 42px;
  border-radius: 12px;
  background: var(--lite-accent);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 800;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.action-card {
  min-height: 138px;
  padding: 18px;
}

.action-card--primary {
  background: #ffffff;
}

.action-emoji {
  font-size: 28px;
}

.action-title {
  display: block;
  margin-top: 16px;
  font-size: 18px;
  font-weight: 800;
}

.section-block {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.section-header {
  margin-top: 10px;
  margin-bottom: 2px;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.section-link {
  font-size: 14px;
  color: var(--visago-primary);
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

.destination-image--fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #eef3f7;
  font-size: 42px;
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

.country-region {
  display: block;
  margin-top: 8px;
  font-size: 13px;
  color: var(--lite-text-muted);
}

.country-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.country-card {
  padding: 16px;
}

.country-top {
  display: flex;
  align-items: center;
  gap: 8px;
}

.country-flag {
  font-size: 22px;
  font-weight: 800;
  color: var(--lite-primary);
}

.country-name {
  font-size: 16px;
  font-weight: 800;
}
</style>
