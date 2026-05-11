<template>
  <view class="page-screen home-screen">
    <VuPageHeader home logo-src="/static/header-logo.png" />
    <view class="page-top-spacer" :style="topSpacerStyle" />

    <view class="page-width page-body home-body">
      <view class="hero-block">
        <text class="hero-title">想去哪里探索世界？</text>
        <view class="app-card search-card">
          <u-search
            v-model="keyword"
            placeholder="搜索国家、地区或签证类型"
            shape="round"
            :show-action="true"
            action-text="搜索"
            bg-color="#ffffff"
            search-icon-color="#8a94a6"
            @search="searchCountries"
            @custom="searchCountries"
          />
        </view>
      </view>

      <view class="type-row">
        <view
          v-for="item in visaTypeMenu"
          :key="item.key"
          class="type-item"
          @tap="onVisaTypeTap(item)"
        >
          <view class="type-icon-box" :class="{ 'type-icon-box--active': activeVisaType === item.key }">
            <u-icon :name="item.icon" size="26" :color="activeVisaType === item.key ? '#1677ff' : item.color" />
          </view>
          <text class="type-label" :class="{ 'type-label--active': activeVisaType === item.key }">{{ item.label }}</text>
        </view>
      </view>

      <view class="quick-row">
        <view class="app-card quick-card" hover-class="quick-card--hover" @tap="openCountries">
          <view class="quick-icon">
            <u-icon name="map" size="28" color="#1677ff" />
          </view>
          <text class="quick-title">查国家签证</text>
          <text class="quick-desc">按目的地查看签证类型</text>
        </view>
        <view class="app-card quick-card" hover-class="quick-card--hover" @tap="openGoals">
          <view class="quick-icon quick-icon--warm">
            <u-icon name="bookmark" size="28" color="#ff8f1f" />
          </view>
          <text class="quick-title">我的目标</text>
          <text class="quick-desc">同步你的办理进度</text>
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
            hover-class="hot-card--hover"
            @tap="openHotDetail(item)"
          >
            <view class="hot-image-wrap">
              <image v-if="item.image" class="hot-image" :src="item.image" mode="aspectFill" />
              <view v-else class="hot-image hot-image--fallback">
                <text>{{ item.flag || '签' }}</text>
              </view>
              <view v-if="item.hot" class="hot-tag-wrap">
                <u-tag text="热门" type="warning" size="mini" />
              </view>
            </view>
            <view class="hot-body">
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
            <u-icon name="empty-search" size="28" color="#9aa4b5" />
            <text class="hot-empty-text">暂无热门目的地</text>
          </view>
        </view>
      </scroll-view>

      <view class="section-head section-head--countries">
        <text class="app-section-title">先选国家</text>
        <text class="app-link" @tap="openCountries">查看全部</text>
      </view>

      <view class="country-row">
        <view
          v-for="country in countryPreview"
          :key="country.id"
          class="app-card country-card"
          hover-class="country-card--hover"
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

    <VuTabbar active="home" />
  </view>
</template>

<script>
import VuPageHeader from '../../components/uview/VuPageHeader.vue'
import VuTabbar from '../../components/uview/VuTabbar.vue'
import { api } from '../../utils/api'
import { getTopSpacerStyle } from '../../utils/layout'

export default {
  components: {
    VuPageHeader,
    VuTabbar,
  },
  data() {
    return {
      keyword: '',
      activeVisaType: 'tourism',
      loadingHot: false,
      hotDestinations: [],
      countryPreview: [],
      topSpacerStyle: 'height:108px;',
      visaTypeMenu: [
        { key: 'tourism', label: '旅游签', icon: 'map', color: '#2f89ff' },
        { key: 'business', label: '商务签', icon: 'bag', color: '#7c3aed' },
        { key: 'study', label: '学生签', icon: 'star', color: '#f59e0b' },
        { key: 'work', label: '工作签', icon: 'file-text', color: '#18b566' },
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
  background: #f7f8fb;
}

.page-top-spacer {
  height: 108px;
}

.home-body {
  padding-bottom: 116px;
}

.hero-block {
  padding-top: 4px;
}

.hero-title {
  display: block;
  font-size: 28px;
  line-height: 1.24;
  font-weight: 900;
  color: #1f2329;
}

.search-card {
  margin-top: 14px;
  padding: 6px 8px;
  border-radius: 22px;
}

.type-row {
  display: flex;
  margin-left: -6px;
  margin-right: -6px;
  margin-top: 16px;
}

.type-item {
  width: 25%;
  padding-left: 6px;
  padding-right: 6px;
}

.type-icon-box {
  height: 74px;
  border-radius: 16px;
  background: #ffffff;
  border: 1px solid #e8edf5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.type-icon-box--active {
  border-color: rgba(22, 119, 255, 0.32);
  background: #eef6ff;
}

.type-label {
  display: block;
  margin-top: 8px;
  font-size: 13px;
  line-height: 1;
  color: #667085;
  text-align: center;
  font-weight: 600;
}

.type-label--active {
  color: #1677ff;
  font-weight: 800;
}

.quick-row {
  display: flex;
  margin-left: -6px;
  margin-right: -6px;
  margin-top: 18px;
}

.quick-card {
  width: 50%;
  min-height: 128px;
  margin-left: 6px;
  margin-right: 6px;
  padding: 16px;
}

.quick-card--hover {
  background: #f8fbff;
}

.quick-icon {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  background: #eef6ff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.quick-icon--warm {
  background: #fff4e7;
}

.quick-title {
  display: block;
  margin-top: 14px;
  font-size: 18px;
  line-height: 1.25;
  font-weight: 800;
  color: #1f2329;
}

.quick-desc {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.35;
  color: #7b8495;
}

.section-head {
  margin-top: 18px;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.section-head--countries {
  margin-top: 20px;
}

.hot-scroll {
  width: 100%;
  white-space: nowrap;
}

.hot-row {
  display: flex;
  padding-right: 8px;
}

.hot-card,
.hot-empty {
  width: 272px;
  flex: 0 0 272px;
  margin-right: 14px;
  overflow: hidden;
}

.hot-card--hover {
  opacity: 0.92;
}

.hot-image-wrap {
  position: relative;
  height: 130px;
}

.hot-image {
  width: 100%;
  height: 130px;
}

.hot-image--fallback {
  background: #edf2f8;
  color: #1677ff;
  font-size: 42px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
}

.hot-tag-wrap {
  position: absolute;
  right: 10px;
  top: 8px;
}

.hot-body {
  padding: 12px;
}

.hot-title-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
}

.hot-name {
  max-width: 190px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  font-size: 18px;
  line-height: 1.25;
  font-weight: 800;
  color: #1f2329;
}

.hot-flag {
  font-size: 18px;
  line-height: 1;
}

.hot-meta {
  margin-top: 8px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.hot-meta-left,
.hot-meta-right {
  display: flex;
  align-items: center;
}

.hot-meta-text {
  margin-left: 4px;
  max-width: 148px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  font-size: 13px;
  color: #7b8495;
}

.hot-price {
  color: #1677ff;
  font-size: 16px;
  font-weight: 900;
}

.hot-price-suffix {
  margin-left: 3px;
  color: #8a94a6;
  font-size: 11px;
}

.hot-empty {
  min-height: 96px;
  padding: 18px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.hot-empty-text {
  display: block;
  margin-top: 10px;
  font-size: 14px;
  color: #7b8495;
}

.country-row {
  display: flex;
  flex-wrap: wrap;
  margin-left: -6px;
  margin-right: -6px;
}

.country-card {
  width: calc(50% - 12px);
  margin-left: 6px;
  margin-right: 6px;
  margin-bottom: 12px;
  padding: 14px;
}

.country-card--hover {
  opacity: 0.92;
}

.country-top {
  display: flex;
  align-items: center;
}

.country-region {
  display: block;
  margin-top: 8px;
  font-size: 13px;
  line-height: 1.35;
  color: #7b8495;
}

.country-flag {
  margin-right: 8px;
  font-size: 22px;
  line-height: 1;
}

.country-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  font-size: 16px;
  font-weight: 800;
  color: #1f2329;
}
</style>
