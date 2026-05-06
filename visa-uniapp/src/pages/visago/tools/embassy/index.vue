<template>
  <view class="embassy-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="驻华使领馆" />

    <scroll-view scroll-y class="embassy-scroll">
      <view class="embassy-content visago-page-width">
        <view class="search-box">
          <text class="material-symbols-outlined search-icon">search</text>
          <input
            v-model.trim="keyword"
            class="search-input"
            placeholder="搜索国家、城市或使领馆..."
            confirm-type="search"
            @confirm="loadEmbassies"
          />
        </view>

        <scroll-view scroll-x class="chip-scroll">
          <view class="chip-row">
            <view
              v-for="item in regionTabs"
              :key="item.key"
              class="chip"
              :class="{ 'chip--active': activeRegion === item.key }"
              @tap="activeRegion = item.key"
            >
              {{ item.label }}
            </view>
          </view>
        </scroll-view>

        <view class="urgent-card">
          <view>
            <text class="urgent-title">模块用途说明</text>
            <text class="urgent-sub">这里展示的是目标国家在中国的使领馆或签证办理联系信息，主要用于签前查询、电话咨询和地址定位。</text>
          </view>
          <view class="urgent-btn" @tap="showModuleIntro">
            <text class="material-symbols-outlined">info</text>
            <text>查看说明</text>
          </view>
        </view>

        <view v-if="loadingEmbassies" class="empty-card">
          <text class="material-symbols-outlined">hourglass_top</text>
          <text>正在加载使领馆数据</text>
        </view>

        <view v-else-if="filteredEmbassies.length" class="embassy-grid">
          <view v-for="item in filteredEmbassies" :key="item.id" class="embassy-card">
            <view class="cover-wrap">
              <image class="cover-image" :src="item.image" mode="aspectFill" />
              <view class="country-badge">
                <text>{{ item.flag || countryFlag(item.countryCode) }}</text>
                <text>{{ item.country }}</text>
              </view>
            </view>

            <view class="card-body">
              <view class="card-head">
                <text class="embassy-name">{{ item.name }}</text>
                <text v-if="item.distance" class="distance">{{ item.distance }}</text>
              </view>
              <view class="meta-row">
                <text class="material-symbols-outlined">location_on</text>
                <text>{{ item.address || item.city || '-' }}</text>
              </view>
              <view class="meta-row">
                <text class="material-symbols-outlined">schedule</text>
                <text>{{ item.hours || '请以官方最新通知为准' }}</text>
              </view>
              <view v-if="item.services && item.services.length" class="service-row">
                <text v-for="tag in item.services" :key="tag" class="service-tag">{{ tag }}</text>
              </view>
              <view class="button-row">
                <view class="map-btn" @tap="openMap(item)">
                  <text class="material-symbols-outlined">map</text>
                  <text>地图查看</text>
                </view>
                <view class="call-btn" @tap="showCall(item.phone)">
                  <text class="material-symbols-outlined">call</text>
                </view>
              </view>
            </view>
          </view>
        </view>

        <view v-else class="empty-card">
          <text class="material-symbols-outlined">travel_explore</text>
          <text>没有找到对应使领馆</text>
          <view class="submit-link" @tap="submitMissing">提交缺失信息</view>
        </view>

        <view v-if="!loadingEmbassies && filteredEmbassies.length" class="submit-box">
          <text>没有找到你要找的使领馆？</text>
          <view @tap="submitMissing">提交缺失信息</view>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { api } from '../../../../utils/api'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'

const REGION_LABELS = {
  all: '全部',
  asia: '亚洲',
  europe: '欧洲',
  'north-america': '北美',
  'south-america': '南美',
  oceania: '大洋洲',
  africa: '非洲',
  'middle-east': '中东',
  express: '签证加急',
}

export default {
  components: {
    VisagoTopBar,
  },
  data() {
    return {
      keyword: '',
      activeRegion: 'all',
      loadingEmbassies: false,
      loadTimer: null,
      requestSeq: 0,
      embassies: [],
    }
  },
  watch: {
    keyword() {
      this.scheduleLoadEmbassies()
    },
    activeRegion() {
      this.scheduleLoadEmbassies()
    },
  },
  computed: {
    regionTabs() {
      const regions = Array.from(new Set(this.embassies.map((item) => String(item.region || '').trim()).filter(Boolean)))
      return [
        { key: 'all', label: REGION_LABELS.all },
        ...regions.map((region) => ({
          key: region,
          label: REGION_LABELS[region] || region,
        })),
        { key: 'express', label: REGION_LABELS.express },
      ]
    },
    filteredEmbassies() {
      const word = this.keyword.trim().toLowerCase()
      return this.embassies.filter((item) => {
        const services = Array.isArray(item.services) ? item.services : []
        const matchRegion =
          this.activeRegion === 'all' ||
          item.region === this.activeRegion ||
          (this.activeRegion === 'express' && services.some((tag) => String(tag || '').includes('加急')))
        const haystack = `${item.country || ''}${item.name || ''}${item.city || ''}${item.address || ''}${services.join('')}`.toLowerCase()
        return matchRegion && (!word || haystack.includes(word))
      })
    },
  },
  onShow() {
    applyTheme(getStoredTheme())
    this.loadEmbassies()
  },
  onUnload() {
    this.clearLoadTimer()
  },
  methods: {
    clearLoadTimer() {
      if (this.loadTimer) {
        clearTimeout(this.loadTimer)
        this.loadTimer = null
      }
    },
    scheduleLoadEmbassies() {
      this.clearLoadTimer()
      const delay = this.keyword ? 280 : 0
      this.loadTimer = setTimeout(() => {
        this.loadEmbassies()
      }, delay)
    },
    async loadEmbassies() {
      this.clearLoadTimer()
      const seq = ++this.requestSeq
      this.loadingEmbassies = true
      try {
        const params = {
          enabled: 1,
          q: this.keyword,
        }
        if (this.activeRegion && this.activeRegion !== 'all' && this.activeRegion !== 'express') {
          params.region = this.activeRegion
        }
        const items = await api.listEmbassies(params)
        if (seq !== this.requestSeq) {
          return
        }
        this.embassies = Array.isArray(items) ? items : []
      } catch (error) {
        if (seq !== this.requestSeq) {
          return
        }
        this.embassies = []
        uni.showToast({
          title: (error && error.message) || '加载失败',
          icon: 'none',
        })
      } finally {
        if (seq === this.requestSeq) {
          this.loadingEmbassies = false
        }
      }
    },
    countryFlag(code) {
      const value = String(code || '').trim().toUpperCase()
      if (!/^[A-Z]{2}$/.test(value)) {
        return ''
      }
      return String.fromCodePoint(...value.split('').map((char) => 127397 + char.charCodeAt(0)))
    },
    openMap(item) {
      const latitude = Number(item.latitude || 0)
      const longitude = Number(item.longitude || 0)
      if (latitude && longitude) {
        uni.openLocation({
          latitude,
          longitude,
          name: item.name,
          address: item.address,
          fail: () => {
            this.copyAddress(item)
          },
        })
        return
      }
      this.copyAddress(item)
    },
    copyAddress(item) {
      const text = String(item.address || item.name || '').trim()
      if (!text) {
        uni.showToast({
          title: '暂无地址信息',
          icon: 'none',
        })
        return
      }
      uni.setClipboardData({
        data: text,
        success: () => {
          uni.showToast({
            title: '地址已复制',
            icon: 'none',
          })
        },
      })
    },
    showCall(phone) {
      const text = String(phone || '').trim()
      if (!text) {
        uni.showToast({
          title: '暂无联系电话',
          icon: 'none',
        })
        return
      }
      uni.showModal({
        title: '联系电话',
        content: text,
        confirmText: '复制',
        success: (res) => {
          if (res.confirm) {
            uni.setClipboardData({ data: text })
          }
        },
      })
    },
    submitMissing() {
      uni.showToast({
        title: '已收到反馈入口请求',
        icon: 'none',
      })
    },
    showModuleIntro() {
      uni.showModal({
        title: '驻华使领馆模块',
        content: '当前展示的是目标国家驻华使领馆/签证中心信息，用于办理签证前查看地址、电话、办公时间和基础服务信息，不是境外中国领保模块。',
        showCancel: false,
        confirmText: '知道了',
      })
    },
  },
}
</script>

<style scoped>
.embassy-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.embassy-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.embassy-content {
  padding: 18px 16px 32px;
  box-sizing: border-box;
}

.urgent-title,
.urgent-sub,
.embassy-name,
.submit-box text {
  display: block;
}

.search-box {
  height: 46px;
  border-radius: 15px;
  background: var(--visago-surface-soft);
  border: 1px solid var(--visago-line);
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 13px;
}

.search-icon {
  color: var(--visago-text-soft);
  font-size: 21px;
}

.search-input {
  flex: 1;
  height: 100%;
  font-size: 14px;
  color: var(--visago-text);
}

.chip-scroll {
  margin: 13px -16px 0;
  white-space: nowrap;
}

.chip-row {
  display: flex;
  gap: 8px;
  padding: 0 16px 2px;
}

.chip {
  padding: 9px 15px;
  border-radius: 999px;
  background: var(--visago-surface);
  color: var(--visago-text-muted);
  border: 1px solid var(--visago-line);
  font-size: 13px;
  font-weight: 800;
}

.chip--active {
  background: var(--visago-primary);
  border-color: var(--visago-primary);
  color: #fff;
  box-shadow: 0 8px 18px rgba(15, 101, 216, 0.2);
}

.urgent-card {
  margin-top: 16px;
  border-radius: 18px;
  padding: 14px;
  background: color-mix(in srgb, var(--visago-primary) 9%, var(--visago-surface));
  border: 1px solid rgba(15, 101, 216, 0.2);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.urgent-title {
  font-size: 15px;
  font-weight: 900;
}

.urgent-sub {
  margin-top: 4px;
  font-size: 11px;
  line-height: 1.4;
  color: var(--visago-text-muted);
}

.urgent-btn {
  height: 38px;
  padding: 0 12px;
  border-radius: 999px;
  background: var(--visago-primary);
  color: #fff;
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 13px;
  font-weight: 900;
  flex-shrink: 0;
}

.urgent-btn .material-symbols-outlined {
  font-size: 18px;
}

.embassy-grid {
  margin-top: 16px;
  display: grid;
  gap: 14px;
}

.embassy-card {
  border-radius: 20px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
  overflow: hidden;
}

.cover-wrap {
  position: relative;
  height: 132px;
  overflow: hidden;
  background: var(--visago-surface-soft);
}

.cover-image {
  width: 100%;
  height: 100%;
}

.country-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  height: 28px;
  padding: 0 9px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.58);
  color: #fff;
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.14);
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  font-weight: 800;
}

.card-body {
  padding: 14px;
}

.card-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.embassy-name {
  font-size: 17px;
  line-height: 1.3;
  font-weight: 900;
}

.distance {
  padding: 4px 8px;
  border-radius: 999px;
  background: rgba(15, 101, 216, 0.13);
  color: var(--visago-primary);
  font-size: 11px;
  font-weight: 900;
  flex-shrink: 0;
}

.meta-row {
  margin-top: 10px;
  color: var(--visago-text-muted);
  display: flex;
  align-items: flex-start;
  gap: 7px;
  font-size: 12px;
  line-height: 1.45;
}

.meta-row .material-symbols-outlined {
  color: var(--visago-primary);
  font-size: 18px;
}

.service-row {
  margin-top: 12px;
  display: flex;
  gap: 7px;
  flex-wrap: wrap;
}

.service-tag {
  padding: 5px 8px;
  border-radius: 999px;
  background: var(--visago-surface-soft);
  color: var(--visago-text-muted);
  font-size: 11px;
  font-weight: 800;
}

.button-row {
  display: grid;
  grid-template-columns: 1fr 46px;
  gap: 9px;
  margin-top: 14px;
}

.map-btn,
.call-btn {
  height: 44px;
  border-radius: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  font-size: 13px;
  font-weight: 900;
}

.map-btn {
  background: var(--visago-primary);
  color: #fff;
  box-shadow: 0 8px 18px rgba(15, 101, 216, 0.2);
}

.call-btn {
  background: var(--visago-surface-soft);
  color: var(--visago-primary);
}

.empty-card,
.submit-box {
  margin-top: 22px;
  color: var(--visago-text-muted);
  text-align: center;
}

.empty-card {
  border-radius: 18px;
  padding: 28px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.submit-box {
  padding-bottom: 14px;
  font-size: 13px;
}

.submit-link,
.submit-box view {
  margin-top: 8px;
  color: var(--visago-primary);
  font-weight: 900;
}

:global(html.theme-dark) .embassy-card,
:global(html.theme-dark) .empty-card {
  background: #171b24;
}
</style>
