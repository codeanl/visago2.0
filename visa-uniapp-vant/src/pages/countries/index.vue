<template>
  <view class="v-screen country-screen">
    <VNavBar title="签证国家" :show-back="true" />
    <view class="page-top-spacer" :style="topSpacerStyle" />

    <view class="v-page-width country-body">
      <view class="filter-card v-card">
        <view class="search-box">
          <VIcon class="search-icon" name="search" />
          <view class="search-input-wrap">
            <input
              v-model.trim="keyword"
              class="search-input"
              placeholder="搜索目的地或签证类型"
              confirm-type="search"
              @confirm="triggerSearch"
            />
          </view>
          <view v-if="keyword" class="clear-btn" @tap="clearKeyword">
            <VIcon class="clear-icon" name="close" />
          </view>
        </view>

        <scroll-view class="region-scroll" scroll-x>
          <view class="region-row">
            <view
              v-for="region in visaRegions"
              :key="region"
              class="region-chip"
              :class="{ 'region-chip--active': selectedRegion === region }"
              @tap="selectRegion(region)"
            >
              {{ region }}
            </view>
          </view>
        </scroll-view>
      </view>

      <scroll-view class="country-scroll" scroll-y>
        <view class="country-list">
          <view
            v-for="country in countries"
            :key="country.id"
            class="v-card country-card"
            @tap="openCountry(country)"
          >
            <image v-if="country.image" class="country-image" :src="country.image" mode="aspectFill" />
            <view v-else class="country-image country-image--fallback">{{ country.flag || '签' }}</view>

            <view class="country-main">
              <text class="country-name">{{ country.name }}</text>
              <view class="tag-row">
                <VTag v-for="tag in country.tags || []" :key="tag" :text="tag" type="default" />
              </view>
              <text class="country-note">{{ country.note || '查看该国家的签证类型与办理信息' }}</text>
            </view>

            <VIcon class="arrow" name="arrow" />
          </view>

          <view v-if="loading" class="v-card empty-card">
            <text class="empty-title">正在加载国家数据</text>
            <text class="empty-sub">请稍等，正在读取签证国家列表</text>
          </view>

          <view v-else-if="!countries.length" class="v-card empty-card">
            <text class="empty-title">没有匹配结果</text>
            <text class="empty-sub">请尝试更换关键词或切换地区筛选</text>
          </view>
        </view>
      </scroll-view>
    </view>
  </view>
</template>

<script>
import VIcon from '../../components/vant/VIcon.vue'
import VNavBar from '../../components/vant/VNavBar.vue'
import VTag from '../../components/vant/VTag.vue'
import { api } from '../../utils/api'
import { getTopSpacerStyle } from '../../utils/layout'

const ALL_REGION = '全部地区'

export default {
  components: { VIcon, VNavBar, VTag },
  data() {
    return {
      keyword: '',
      selectedRegion: ALL_REGION,
      visaRegions: [ALL_REGION],
      countries: [],
      loading: false,
      regionLoaded: false,
      topSpacerStyle: 'height: 96px;',
    }
  },
  watch: {
    keyword(value) {
      if (!String(value || '').trim()) {
        this.loadCountries({ syncRegions: !this.regionLoaded })
      }
    },
  },
  onLoad(query) {
    this.syncLayout()
    const keyword = String(query.q || '').trim()
    if (keyword) this.keyword = decodeURIComponent(keyword)
  },
  onShow() {
    this.syncLayout()
    this.loadCountries({ syncRegions: !this.regionLoaded })
  },
  methods: {
    syncLayout() {
      this.topSpacerStyle = getTopSpacerStyle('sub')
    },
    triggerSearch() {
      this.loadCountries({ syncRegions: !this.regionLoaded })
    },
    clearKeyword() {
      if (!this.keyword) return
      this.keyword = ''
      this.loadCountries({ syncRegions: !this.regionLoaded })
    },
    async loadCountries(options = {}) {
      const { syncRegions = false } = options
      this.loading = true
      try {
        const region = this.selectedRegion === ALL_REGION ? '' : this.selectedRegion
        const keyword = String(this.keyword || '').trim()
        const countries = (await api.listCountries(keyword, region)) || []
        this.countries = countries
        if (syncRegions) {
          const allCountries = region || keyword ? (await api.listCountries('', '')) || [] : countries
          const regions = Array.from(new Set(allCountries.map((item) => item.region).filter(Boolean)))
          this.visaRegions = [ALL_REGION, ...regions]
          this.regionLoaded = true
        }
      } catch (error) {
        uni.showToast({ title: (error && error.message) || '加载失败', icon: 'none' })
      } finally {
        this.loading = false
      }
    },
    selectRegion(region) {
      if (this.selectedRegion === region) return
      this.selectedRegion = region
      this.loadCountries()
    },
    openCountry(country) {
      uni.navigateTo({
        url: `/pages/visas/index?countryId=${country.id}&countryName=${encodeURIComponent(country.name || '')}`,
      })
    },
  },
}
</script>

<style scoped>
.country-screen {
  background: var(--vant-bg);
  height: 100vh;
  overflow: hidden;
}

.page-top-spacer {
  height: 96px;
}

.country-body {
  height: calc(100vh - 96px);
  display: flex;
  flex-direction: column;
  padding: 0 16px 20px;
}

.filter-card {
  padding: 12px;
}

.search-box {
  height: 48px;
  border-radius: 14px;
  background: var(--vant-surface-soft);
  display: flex;
  align-items: center;
  padding: 0 12px;
}

.search-icon {
  font-size: 18px;
  color: var(--vant-text-3);
}

.search-input-wrap {
  flex: 1;
  min-width: 0;
}

.search-input {
  width: 100%;
  height: 48px;
  margin: 0 8px;
  font-size: 14px;
}

.clear-btn {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.clear-icon {
  font-size: 16px;
  color: var(--vant-text-3);
}

.region-scroll {
  margin-top: 12px;
  white-space: nowrap;
}

.region-row {
  display: inline-flex;
  padding-right: 8px;
}

.region-chip {
  min-width: 64px;
  padding: 8px 14px;
  border-radius: 999px;
  background: var(--vant-surface-soft);
  color: var(--vant-text-2);
  font-size: 13px;
  line-height: 1;
  font-weight: 700;
  text-align: center;
  margin-right: 8px;
}

.region-chip--active {
  background: var(--vant-primary);
  color: #fff;
}

.country-scroll {
  flex: 1;
  min-height: 0;
  margin-top: 12px;
}

.country-list {
  display: block;
  padding-bottom: 12px;
}

.country-card {
  display: flex;
  align-items: center;
  padding: 12px;
  margin-bottom: 10px;
}

.country-image {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  flex-shrink: 0;
  background: var(--vant-surface-soft);
}

.country-image--fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26px;
}

.country-main {
  flex: 1;
  min-width: 0;
  margin-left: 10px;
}

.country-name {
  display: block;
  font-size: 17px;
  line-height: 1.2;
  font-weight: 700;
  color: var(--vant-text);
}

.tag-row {
  margin-top: 6px;
}

.tag-row :deep(.v-tag) {
  margin-right: 6px;
  margin-bottom: 6px;
}

.country-note {
  display: block;
  margin-top: 4px;
  color: var(--vant-text-3);
  font-size: 12px;
  line-height: 1.45;
}

.arrow {
  margin-left: 10px;
  font-size: 18px;
  color: var(--vant-text-3);
}

.empty-card {
  padding: 22px 18px;
  text-align: center;
}

.empty-title {
  font-size: 16px;
  font-weight: 700;
}

.empty-sub {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: var(--vant-text-3);
}
</style>
