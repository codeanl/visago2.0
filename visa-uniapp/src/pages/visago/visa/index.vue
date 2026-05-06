<template>
  <view class="visa-page">
    <VisagoTopBar page-name="签证" />

    <view class="visa-content visago-page-width">
      <view class="filter-panel">
        <view class="search-box">
          <text class="material-symbols-outlined search-icon" @tap="triggerSearch">search</text>
          <input
            v-model.trim="keyword"
            class="search-input"
            placeholder="搜索目的地或签证类型"
            confirm-type="search"
            @confirm="triggerSearch"
          />
          <text v-if="keyword" class="material-symbols-outlined clear-icon" @tap="clearKeyword">close</text>
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
            v-for="country in filteredCountries"
            :key="country.id"
            class="country-card visago-card"
            @tap="openCountryVisas(country)"
          >
            <image class="country-image" :src="country.image" mode="aspectFill" />
            <view class="country-main">
              <text class="country-name">{{ country.name }}</text>
              <view class="tag-row">
                <text v-for="tag in country.tags" :key="tag" class="tag">{{ tag }}</text>
              </view>
              <text class="country-note">{{ country.note }}</text>
            </view>
            <text class="material-symbols-outlined arrow">chevron_right</text>
          </view>

          <view v-if="loading" class="empty-card visago-card">
            <text class="empty-title">正在加载国家数据</text>
            <text class="empty-sub">请稍等，正在从数据库读取签证国家列表</text>
          </view>

          <view v-else-if="!filteredCountries.length" class="empty-card visago-card">
            <text class="empty-title">没有匹配结果</text>
            <text class="empty-sub">请尝试更换关键词或切换地区筛选</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <VisagoBottomNav active-tab="visa" />
  </view>
</template>

<script>
import VisagoTopBar from '../../../components/VisagoTopBar.vue'
import VisagoBottomNav from '../../../components/VisagoBottomNav.vue'
import { api } from '../../../utils/api'
import { applyTheme, getStoredTheme } from '../../../utils/theme'

const ALL_REGION = '全部地区'

function countryMatches(country, keyword) {
  const word = String(keyword || '').trim().toLowerCase()
  if (!word) return true
  const pool = [
    country.name,
    country.region,
    country.code,
    country.note,
    ...(country.tags || []),
    ...(country.keywords || []),
  ]
  return pool.some((item) => String(item || '').toLowerCase().includes(word))
}

export default {
  components: {
    VisagoTopBar,
    VisagoBottomNav,
  },
  data() {
    return {
      keyword: '',
      selectedRegion: ALL_REGION,
      visaRegions: [ALL_REGION],
      countries: [],
      loading: false,
      regionLoaded: false,
    }
  },
  watch: {
    keyword(value) {
      if (!String(value || '').trim()) {
        this.loadCountries({ syncRegions: !this.regionLoaded })
      }
    },
  },
  computed: {
    filteredCountries() {
      const word = this.keyword.trim().toLowerCase()
      return this.countries.filter((country) => {
        const passRegion = this.selectedRegion === ALL_REGION || country.region === this.selectedRegion
        if (!passRegion) {
          return false
        }
        return countryMatches(country, word)
      })
    },
  },
  onLoad(query) {
    const keyword = String(query.keyword || '').trim()
    if (keyword) {
      this.keyword = keyword
    }
  },
  onShow() {
    applyTheme(getStoredTheme())
    this.loadCountries({ syncRegions: !this.regionLoaded })
  },
  methods: {
    triggerSearch() {
      this.loadCountries({ syncRegions: !this.regionLoaded })
    },
    clearKeyword() {
      if (!this.keyword) {
        return
      }
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
          const allCountries = region ? (await api.listCountries('', '')) || [] : countries
          const regions = Array.from(new Set(allCountries.map((item) => item.region).filter(Boolean)))
          this.visaRegions = [ALL_REGION, ...regions]
          this.regionLoaded = true
        }
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '加载失败',
          icon: 'none',
        })
      } finally {
        this.loading = false
      }
    },
    selectRegion(region) {
      if (this.selectedRegion === region) {
        return
      }
      this.selectedRegion = region
      this.loadCountries()
    },
    openCountryVisas(country) {
      uni.navigateTo({
        url: `/pages/visago/visa-country/index?countryId=${country.id}`,
        fail: () => {
          uni.showToast({
            title: '打开失败',
            icon: 'none',
          })
        },
      })
    },
  },
}
</script>

<style scoped>
.visa-page {
  min-height: 100vh;
  background: var(--visago-bg);
}

.visa-content {
  height: 100vh;
  box-sizing: border-box;
  padding: 88px 16px calc(84px + var(--visago-safe-bottom));
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 0;
}

.filter-panel {
  position: sticky;
  top: 0;
  z-index: 5;
  background: var(--visago-bg);
  padding-top: 8px;
  padding-bottom: 8px;
}

.search-box {
  height: 48px;
  border-radius: 12px;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 12px;
}

.search-icon {
  color: var(--visago-text-soft);
  font-size: 20px;
}

.search-input {
  flex: 1;
  height: 48px;
  font-size: 14px;
}

.clear-icon {
  color: var(--visago-text-soft);
  font-size: 18px;
  padding: 4px;
}

.region-scroll {
  margin-top: 12px;
  width: 100%;
  white-space: nowrap;
}

.region-row {
  display: inline-flex;
  gap: 8px;
  padding-right: 8px;
}

.region-chip {
  padding: 8px 14px;
  border-radius: 9999px;
  background: var(--visago-surface-soft);
  color: var(--visago-text-muted);
  font-size: 13px;
  line-height: 1;
}

.region-chip--active {
  background: var(--visago-primary);
  color: #fff;
}

.country-scroll {
  flex: 1;
  min-height: 0;
}

.country-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding-top: 4px;
  padding-bottom: calc(12px + var(--visago-safe-bottom));
}

.country-card {
  padding: 10px;
  display: flex;
  gap: 10px;
  align-items: center;
}

.country-image {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  flex-shrink: 0;
}

.country-main {
  flex: 1;
  min-width: 0;
}

.country-name {
  display: block;
  font-size: 16px;
  line-height: 1.2;
  font-weight: 600;
}

.tag-row {
  margin-top: 4px;
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.tag {
  font-size: 10px;
  color: var(--visago-primary);
  border-radius: 9999px;
  background: #eaf2ff;
  padding: 2px 7px;
}

.country-note {
  margin-top: 4px;
  display: block;
  color: var(--visago-text-muted);
  font-size: 11px;
  line-height: 1.35;
}

.arrow {
  color: #bac2d7;
  font-size: 20px;
}

.empty-card {
  padding: 18px 14px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.empty-title {
  font-size: 15px;
  font-weight: 600;
}

.empty-sub {
  font-size: 12px;
  color: var(--visago-text-muted);
}
</style>
