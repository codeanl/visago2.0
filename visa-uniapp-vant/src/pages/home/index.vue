<template>
  <view class="v-screen home-screen">
    <VNavBar home title="VisaNow" logo-src="/static/header-logo.png" />
    <view class="page-top-spacer" :style="topSpacerStyle" />

    <view class="v-page-width home-body">
      <view class="hero-block">
        <text class="hero-title">想去哪里探索世界？</text>

        <view class="v-card search-card">
          <view class="search-input-wrap">
            <input
              v-model.trim="keyword"
              class="search-input"
              placeholder="搜索国家、地区或签证类型"
              confirm-type="search"
              @confirm="searchCountries"
            />
          </view>
          <view class="v-button-primary search-btn" @tap="searchCountries">搜索</view>
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
            <VIcon class="type-icon" :name="item.icon" />
          </view>
          <text class="type-label">{{ item.label }}</text>
        </view>
      </view>

      <view class="quick-grid">
        <view class="v-card quick-card" @tap="openCountries">
          <VIcon class="quick-icon" name="flag-o" />
          <text class="quick-title">查国家签证</text>
        </view>
        <view class="v-card quick-card" @tap="openGoals">
          <VIcon class="quick-icon" name="todo-list-o" />
          <text class="quick-title">我的目标</text>
        </view>
      </view>

      <view class="section-head">
        <text class="v-section-title">热门目的地</text>
      </view>

      <scroll-view class="hot-scroll" scroll-x>
        <view class="hot-row">
          <view
            v-for="item in hotDestinations"
            :key="item.visaId"
            class="v-card hot-card"
            @tap="openHotDetail(item)"
          >
            <view class="hot-image-wrap">
              <image v-if="item.image" class="hot-image" :src="item.image" mode="aspectFill" />
              <view v-else class="hot-image hot-image--fallback">{{ item.flag || '签' }}</view>
              <VTag v-if="item.hot" text="热门" type="warning" class="hot-badge" />
            </view>

            <view class="hot-card-body">
              <view class="hot-title-row">
                <text class="hot-name">{{ item.name }}</text>
                <text class="hot-flag">{{ item.flag }}</text>
              </view>

              <view class="hot-meta">
                <view class="hot-meta-left">
                  <VIcon class="meta-icon" name="passed" />
                  <text>{{ item.type }} | {{ item.time }}</text>
                </view>
                <view class="hot-meta-right">
                  <text class="hot-price">{{ item.price || '待更新' }}</text>
                  <text class="hot-price-suffix">起</text>
                </view>
              </view>
            </view>
          </view>

          <view v-if="loadingHot" class="v-card hot-empty">正在加载热门目的地...</view>
          <view v-else-if="!hotDestinations.length" class="v-card hot-empty">暂无热门目的地</view>
        </view>
      </scroll-view>

      <view class="section-head">
        <text class="v-section-title">先选国家</text>
        <text class="link-text" @tap="openCountries">查看全部</text>
      </view>

      <view class="country-grid">
        <view
          v-for="country in countryPreview"
          :key="country.id"
          class="v-card country-card"
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

    <VTabbar active="home" />
  </view>
</template>

<script>
import VNavBar from '../../components/vant/VNavBar.vue'
import VTabbar from '../../components/vant/VTabbar.vue'
import VIcon from '../../components/vant/VIcon.vue'
import VTag from '../../components/vant/VTag.vue'
import { api } from '../../utils/api'
import { getTopSpacerStyle } from '../../utils/layout'

export default {
  components: { VNavBar, VTabbar, VIcon, VTag },
  data() {
    return {
      keyword: '',
      activeVisaType: 'tourism',
      loadingHot: false,
      hotDestinations: [],
      countryPreview: [],
      topSpacerStyle: 'height: 108px;',
      visaTypeMenu: [
        { key: 'tourism', label: '旅游签', icon: 'guide-o' },
        { key: 'business', label: '商务签', icon: 'bag-o' },
        { key: 'study', label: '学生签', icon: 'award-o' },
        { key: 'work', label: '工作签', icon: 'notes-o' },
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
  background: var(--vant-bg);
}

.page-top-spacer {
  height: 108px;
  flex-shrink: 0;
}

.home-body {
  padding: 0 16px calc(116px + var(--vant-safe-bottom));
}

.hero-block {
  display: block;
}

.hero-title {
  font-size: 28px;
  line-height: 1.24;
  font-weight: 800;
  color: var(--vant-text);
}

.search-card {
  display: flex;
  align-items: center;
  padding: 10px;
  border-radius: 20px;
  margin-top: 14px;
}

.search-input-wrap {
  flex: 1;
  min-width: 0;
}

.search-input {
  width: 100%;
  height: 44px;
  padding: 0 14px;
  border-radius: 14px;
  background: var(--vant-surface);
  color: var(--vant-text);
  font-size: 14px;
}

.search-btn {
  width: 88px;
  margin-left: 10px;
  flex-shrink: 0;
}

.type-grid,
.quick-grid,
.country-grid {
  display: flex;
  flex-wrap: wrap;
  margin: 0 -6px;
}

.type-grid {
  margin-top: 8px;
  margin-bottom: 14px;
}

.type-item {
  width: 25%;
  padding: 0 6px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.type-icon-wrap {
  width: 100%;
  height: 74px;
  border-radius: 16px;
  background: #fff;
  border: 1px solid var(--vant-border);
  display: flex;
  align-items: center;
  justify-content: center;
}

.type-item--active .type-icon-wrap {
  border-color: rgba(22, 119, 255, 0.28);
  background: rgba(22, 119, 255, 0.04);
}

.type-icon {
  font-size: 24px;
  color: var(--vant-primary);
}

.type-label {
  margin-top: 8px;
  font-size: 13px;
  line-height: 1;
  font-weight: 600;
  color: var(--vant-text-2);
}

.type-item--active .type-label {
  color: var(--vant-primary);
}

.quick-card {
  width: calc(50% - 12px);
  margin: 0 6px 12px;
  min-height: 132px;
  padding: 18px;
  box-sizing: border-box;
}

.quick-icon {
  font-size: 28px;
  color: var(--vant-primary);
}

.quick-title {
  display: block;
  margin-top: 16px;
  font-size: 18px;
  font-weight: 800;
  color: var(--vant-text);
}

.section-head {
  margin: 12px 0;
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

.hot-card {
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
  background: var(--vant-surface-soft);
  font-size: 42px;
}

.hot-badge {
  position: absolute;
  right: 10px;
  top: 8px;
}

.hot-card-body {
  padding: 12px 14px;
}

.hot-title-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
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
  font-size: 13px;
  color: var(--vant-text-3);
}

.hot-meta-left,
.hot-meta-right {
  display: flex;
  align-items: center;
}

.meta-icon {
  margin-right: 4px;
  font-size: 14px;
}

.hot-price {
  color: var(--vant-primary);
  font-size: 16px;
  font-weight: 800;
  margin-right: 4px;
}

.hot-price-suffix {
  font-size: 11px;
  color: var(--vant-text-3);
}

.hot-empty {
  width: 272px;
  min-height: 86px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  font-size: 14px;
  color: var(--vant-text-3);
  text-align: center;
  margin-right: 14px;
}

.country-card {
  width: calc(50% - 12px);
  margin: 0 6px 12px;
  padding: 16px;
  box-sizing: border-box;
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
  color: var(--vant-text-3);
}
</style>
