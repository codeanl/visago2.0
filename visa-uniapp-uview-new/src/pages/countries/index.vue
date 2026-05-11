<template>
  <view class="page-screen countries-screen">
    <VuPageHeader title="签证国家" :show-back="true" />
    <view class="page-top-spacer" :style="topSpacerStyle" />

    <view class="page-width page-body countries-body">
      <view class="app-card filter-card">
        <u-search
          v-model="keyword"
          placeholder="搜索目的地或签证类型"
          shape="round"
          :show-action="true"
          action-text="搜索"
          bg-color="#ffffff"
          search-icon-color="#8a94a6"
          @search="triggerSearch"
          @custom="triggerSearch"
        />

        <scroll-view class="region-scroll" scroll-x>
          <view class="region-row">
            <view
              v-for="region in visaRegions"
              :key="region"
              class="region-chip"
              :class="{ 'region-chip--active': selectedRegion === region }"
              @tap="selectRegion(region)"
            >
              <text>{{ region }}</text>
            </view>
          </view>
        </scroll-view>
      </view>

      <view class="country-list">
        <view
          v-for="country in countries"
          :key="country.id"
          class="app-card country-card"
          hover-class="country-card--hover"
          @tap="openCountry(country)"
        >
          <image v-if="country.image" class="country-image" :src="country.image" mode="aspectFill" />
          <view v-else class="country-image country-image--fallback">
            <text>{{ country.flag || '签' }}</text>
          </view>

          <view class="country-main">
            <text class="country-name">{{ country.name }}</text>
            <view v-if="country.tags && country.tags.length" class="tag-row">
              <view v-for="tag in country.tags" :key="tag" class="tag-cell">
                <u-tag
                  :text="tag"
                  type="primary"
                  size="mini"
                  plain
                  plain-fill
                />
              </view>
            </view>
            <text class="country-note">{{ country.note || '查看该国家的签证类型与办理信息' }}</text>
          </view>

          <u-icon name="arrow-right" size="18" color="#b0b8c5" />
        </view>

        <view v-if="loading" class="app-card app-empty empty-card">
          <u-loading-icon mode="circle" size="22" color="#1677ff" />
          <text class="app-empty-title">正在加载国家数据</text>
          <text class="app-empty-desc">请稍等，正在读取签证国家列表</text>
        </view>

        <view v-else-if="!countries.length" class="app-card app-empty empty-card">
          <u-icon name="empty-search" size="30" color="#9aa4b5" />
          <text class="app-empty-title">没有匹配结果</text>
          <text class="app-empty-desc">请尝试更换关键词或切换地区筛选</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import VuPageHeader from '../../components/uview/VuPageHeader.vue'
import { api } from '../../utils/api'
import { getTopSpacerStyle } from '../../utils/layout'

const ALL_REGION = '全部地区'

export default {
  components: {
    VuPageHeader,
  },
  data() {
    return {
      keyword: '',
      selectedRegion: ALL_REGION,
      visaRegions: [ALL_REGION],
      countries: [],
      loading: false,
      regionLoaded: false,
      topSpacerStyle: 'height:96px;',
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
    if (keyword) {
      this.keyword = decodeURIComponent(keyword)
    }
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
.countries-screen {
  background: #f7f8fb;
}

.page-top-spacer {
  height: 96px;
}

.countries-body {
  padding-bottom: 26px;
}

.filter-card {
  padding: 8px;
}

.region-scroll {
  margin-top: 12px;
  white-space: nowrap;
}

.region-row {
  display: flex;
  padding-right: 8px;
}

.region-chip {
  min-width: 68px;
  margin-right: 8px;
  padding: 8px 14px;
  border-radius: 16px;
  background: #f0f3f8;
  color: #5f6b7a;
  font-size: 13px;
  line-height: 1;
  text-align: center;
  font-weight: 700;
}

.region-chip--active {
  color: #ffffff;
  background: #1677ff;
}

.country-list {
  margin-top: 12px;
  padding-bottom: 12px;
}

.country-card {
  min-height: 92px;
  margin-bottom: 10px;
  padding: 12px;
  display: flex;
  align-items: center;
}

.country-card--hover {
  background: #f8fbff;
}

.country-image {
  width: 64px;
  height: 64px;
  border-radius: 14px;
  flex-shrink: 0;
  background: #eef3f8;
}

.country-image--fallback {
  color: #1677ff;
  font-size: 26px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
}

.country-main {
  flex: 1;
  min-width: 0;
  margin-left: 10px;
  margin-right: 10px;
}

.country-name {
  display: block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  color: #1f2329;
  font-size: 17px;
  line-height: 1.25;
  font-weight: 800;
}

.tag-row {
  margin-top: 6px;
  display: flex;
  flex-wrap: wrap;
}

.tag-cell {
  margin-right: 6px;
  margin-bottom: 4px;
}

.country-note {
  display: block;
  margin-top: 5px;
  color: #7b8495;
  font-size: 12px;
  line-height: 1.45;
}

.empty-card {
  margin-top: 14px;
}
</style>
