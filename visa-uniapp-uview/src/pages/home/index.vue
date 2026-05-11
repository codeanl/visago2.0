<template>
  <view class="app-screen home-screen">
    <UPageHeader home title="VisaNow" logo-src="/static/header-logo.png" />
    <view class="page-top-spacer" :style="topSpacerStyle" />

    <view class="app-page-width home-body">
      <view class="hero-block">
        <text class="hero-title">想去哪里探索世界？</text>

        <view class="app-card search-shell">
          <u-search
            v-model="keyword"
            :show-action="true"
            action-text="搜索"
            placeholder="搜索国家、地区或签证类型"
            bg-color="#ffffff"
            shape="round"
            search-icon-color="#86909c"
            @search="searchCountries"
            @custom="searchCountries"
          />
        </view>
      </view>

      <view class="type-grid">
        <view
          v-for="item in visaTypeMenu"
          :key="item.key"
          class="type-item"
          :class="{ 'type-item--active': activeVisaType === item.key }"
          @tap="onVisaTypeTap(item)"
        >
          <view class="type-icon-wrap">
            <u-icon :name="item.icon" size="26" :color="activeVisaType === item.key ? '#1677ff' : '#4e5969'" />
          </view>
          <text class="type-label">{{ item.label }}</text>
        </view>
      </view>

      <view class="quick-grid">
        <view class="app-card quick-card" @tap="openCountries">
          <u-icon name="map" size="28" color="#1677ff" />
          <text class="quick-title">查国家签证</text>
        </view>
        <view class="app-card quick-card" @tap="openGoals">
          <u-icon name="bookmark" size="28" color="#1677ff" />
          <text class="quick-title">我的目标</text>
        </view>
      </view>

      <view class="section-head">
        <text class="app-section-title">热门目的地</text>
      </view>

      <scroll-view class="hot-scroll" scroll-x>
        <view class="hot-row">
          <view
            v-for="item in hotDestinations"
            :key="item.visaId"
            class="app-card hot-card"
            @tap="openHotDetail(item)"
          >
            <view class="hot-image-wrap">
              <image v-if="item.image" class="hot-image" :src="item.image" mode="aspectFill" />
              <view v-else class="hot-image hot-image--fallback">{{ item.flag || '签' }}</view>
              <view class="hot-badge">
                <u-tag v-if="item.hot" text="热门" type="warning" size="mini" />
              </view>
            </view>

            <view class="hot-card-body">
              <view class="hot-title-row">
                <text class="hot-name">{{ item.name }}</text>
                <text class="hot-flag">{{ item.flag }}</text>
              </view>
              <view class="hot-meta">
                <view class="hot-meta-left">
                  <u-icon name="checkmark-circle" size="14" color="#1677ff" />
                  <text class="hot-meta-text">{{ item.type }} | {{ item.time }}</text>
                </view>
                <view class="hot-meta-right">
                  <text class="hot-price">{{ item.price || '待更新' }}</text>
                  <text class="hot-price-suffix">起</text>
                </view>
              </view>
            </view>
          </view>

          <view v-if="loadingHot" class="app-card hot-empty">
            <u-loading-icon mode="circle" size="22" color="#1677ff" />
            <text class="hot-empty-text">正在加载热门目的地...</text>
          </view>

          <view v-else-if="!hotDestinations.length" class="app-card hot-empty">
            <text class="hot-empty-text">暂无热门目的地</text>
          </view>
        </view>
      </scroll-view>

      <view class="section-head">
        <text class="app-section-title">先选国家</text>
        <text class="app-link" @tap="openCountries">查看全部</text>
      </view>

      <view class="country-grid">
        <view
          v-for="country in countryPreview"
          :key="country.id"
          class="app-card country-card"
          @tap="openCountry(country)"
        >
          <view class="country-top">
            <text class="country-flag">{{ country.flag || '签' }}</text>
            <text class="country-name">{{ country.name }}</text>
          </view>
          <text class="country-region">{{ country.region || '全球签证' }}</text>
        </view>
      </view>
    </view>

    <UMainTabbar active="home" />
  </view>
</template>

<script>
import UPageHeader from '../../components/shared/UPageHeader.vue'
import UMainTabbar from '../../components/shared/UMainTabbar.vue'
import { api } from '../../utils/api'
import { getTopSpacerStyle } from '../../utils/layout'

export default {
  components: { UPageHeader, UMainTabbar },
  data() {
    return {
      keyword: '',
      activeVisaType: 'tourism',
      loadingHot: false,
      hotDestinations: [],
      countryPreview: [],
      topSpacerStyle: 'height: 108px;',
      visaTypeMenu: [
        { key: 'tourism', label: '旅游签', icon: 'map' },
        { key: 'business', label: '商务签', icon: 'bag' },
        { key: 'study', label: '学生签', icon: 'star' },
        { key: 'work', label: '工作签', icon: 'file-text' },
      ],
    }
  },
  onLoad() {
    this.syncLayout()
  },
  onShow() {
    this.syncLayout()
    this.loadHomeData()
  },
  methods: {
    syncLayout() {
      this.topSpacerStyle = getTopSpacerStyle('home')
    },
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
        uni.showToast({ title: (error && error.message) || '加载失败', icon: 'none' })
      } finally {
        this.loadingHot = false
      }
    },
    openCountries() {
      uni.navigateTo({ url: '/pages/countries/index' })
    },
    openGoals() {
      uni.reLaunch({ url: '/pages/goals/index' })
    },
    searchCountries(keywordOverride = '') {
      const q = encodeURIComponent(keywordOverride || this.keyword || '')
      uni.navigateTo({ url: `/pages/countries/index?q=${q}` })
    },
    onVisaTypeTap(item) {
      this.activeVisaType = item.key
      this.openCountries()
    },
    openCountry(country) {
      uni.navigateTo({
        url: `/pages/visas/index?countryId=${country.id}&countryName=${encodeURIComponent(country.name || '')}`,
      })
    },
    openHotDetail(item) {
      if (!item || !item.countryId) return
      uni.navigateTo({
        url: `/pages/visas/index?countryId=${item.countryId}&countryName=${encodeURIComponent(item.name || '')}`,
      })
    },
  },
}
</script>

<style scoped>
.home-screen {
  background: var(--uview-bg);
}

.page-top-spacer {
  height: 108px;
}

.home-body {
  padding: 0 16px calc(116px + var(--uview-safe-bottom));
}

.hero-block {
  padding-top: 4px;
}

.hero-title {
  font-size: 28px;
  line-height: 1.24;
  font-weight: 800;
  color: var(--uview-text);
}

.search-shell {
  margin-top: 14px;
  padding: 6px 8px;
  border-radius: 22px;
}

.type-grid,
.quick-grid,
.country-grid {
  display: flex;
  flex-wrap: wrap;
  margin: 0 -6px;
}

.type-grid {
  margin-top: 16px;
  margin-bottom: 16px;
}

.type-item {
  width: 25%;
  padding: 0 6px;
}

.type-icon-wrap {
  height: 74px;
  border-radius: 16px;
  background: #ffffff;
  border: 1px solid var(--uview-border);
  display: flex;
  align-items: center;
  justify-content: center;
}

.type-item--active .type-icon-wrap {
  border-color: rgba(22, 119, 255, 0.28);
  background: rgba(22, 119, 255, 0.04);
}

.type-label {
  display: block;
  margin-top: 8px;
  text-align: center;
  font-size: 13px;
  font-weight: 600;
  color: var(--uview-text-2);
}

.type-item--active .type-label {
  color: var(--uview-brand);
}

.quick-card {
  width: calc(50% - 12px);
  margin: 0 6px 12px;
  min-height: 126px;
  padding: 18px 16px;
}

.quick-title {
  display: block;
  margin-top: 16px;
  font-size: 18px;
  font-weight: 800;
  color: var(--uview-text);
}

.section-head {
  margin: 14px 0 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.hot-scroll {
  width: 100%;
  white-space: nowrap;
}

.hot-row {
  display: inline-flex;
  padding-right: 8px;
}

.hot-card,
.hot-empty {
  width: 272px;
  overflow: hidden;
  margin-right: 14px;
}

.hot-image-wrap {
  position: relative;
  height: 130px;
}

.hot-image {
  width: 100%;
  height: 100%;
}

.hot-image--fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--uview-surface-soft);
  font-size: 42px;
}

.hot-badge {
  position: absolute;
  right: 10px;
  top: 8px;
}

.hot-card-body {
  padding: 12px;
}

.hot-title-row {
  display: flex;
  justify-content: space-between;
}

.hot-name {
  font-size: 18px;
  font-weight: 700;
}

.hot-flag {
  font-size: 18px;
}

.hot-meta {
  margin-top: 6px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.hot-meta-left,
.hot-meta-right {
  display: flex;
  align-items: center;
}

.hot-meta-left :deep(.u-icon) {
  margin-right: 4px;
}

.hot-meta-text {
  font-size: 13px;
  color: var(--uview-text-3);
}

.hot-price {
  color: var(--uview-brand);
  font-size: 16px;
  font-weight: 800;
  margin-right: 4px;
}

.hot-price-suffix {
  font-size: 11px;
  color: var(--uview-text-3);
}

.hot-empty {
  min-height: 86px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  padding: 16px;
}

.hot-empty-text {
  margin-top: 10px;
  font-size: 14px;
  color: var(--uview-text-3);
  text-align: center;
}

.country-card {
  width: calc(50% - 12px);
  margin: 0 6px 12px;
  padding: 14px;
}

.country-top {
  display: flex;
  align-items: center;
}

.country-flag {
  font-size: 22px;
  margin-right: 8px;
}

.country-name {
  font-size: 16px;
  font-weight: 800;
}

.country-region {
  display: block;
  margin-top: 8px;
  font-size: 13px;
  color: var(--uview-text-3);
}
</style>
