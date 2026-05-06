<template>
  <view class="fav-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="我的收藏" />

    <view class="search-shell">
      <view class="search-wrap visago-page-width">
        <view class="search-box">
          <text class="material-symbols-outlined">search</text>
          <input
            v-model.trim="keyword"
            class="search-input"
            placeholder="搜索收藏的社区帖子..."
          />
        </view>
      </view>
    </view>

    <scroll-view scroll-y class="fav-scroll">
      <view class="fav-content visago-page-width">
        <view class="search-placeholder" />

        <view v-if="loading" class="empty-card visago-card">
          <text class="material-symbols-outlined empty-icon">hourglass_top</text>
          <text class="empty-title">正在加载收藏</text>
        </view>

        <view v-else-if="filteredList.length" class="fav-list">
          <view v-for="item in filteredList" :key="item.id" class="fav-card" @tap="openDetail(item)">
            <image v-if="coverImage(item)" class="fav-image" :src="coverImage(item)" mode="aspectFill" />

            <view class="fav-body">
              <view class="fav-top">
                <text class="fav-type fav-type--community">社区帖子</text>
                <text class="material-symbols-outlined fav-icon">bookmark</text>
              </view>

              <text class="fav-title">{{ item.title }}</text>
              <text class="fav-desc">{{ item.content }}</text>
            </view>
          </view>
        </view>

        <view v-else class="empty-card visago-card">
          <text class="material-symbols-outlined empty-icon">search_off</text>
          <text class="empty-title">未找到匹配收藏</text>
          <text class="empty-sub">试试更短关键词，或先去社区详情里收藏帖子。</text>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { api } from '../../../../utils/api'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'

export default {
  components: { VisagoTopBar },
  data() {
    return {
      keyword: '',
      loading: false,
      list: [],
    }
  },
  computed: {
    filteredList() {
      const word = this.keyword.trim()
      if (!word) return this.list
      return this.list.filter((item) => `${item.title || ''}${item.content || ''}`.includes(word))
    },
  },
  onShow() {
    applyTheme(getStoredTheme())
    this.loadFavorites()
  },
  methods: {
    async loadFavorites() {
      this.loading = true
      try {
        const items = await api.listCommunityFavorites()
        this.list = Array.isArray(items) ? items : []
      } catch (error) {
        this.list = []
        uni.showToast({
          title: (error && error.message) || '加载失败',
          icon: 'none',
        })
      } finally {
        this.loading = false
      }
    },
    coverImage(item) {
      const images = Array.isArray(item && item.images) ? item.images.filter(Boolean) : []
      if (images.length) return images[0]
      return item && item.image ? item.image : ''
    },
    openDetail(item) {
      uni.navigateTo({
        url: `/pages/visago/community/detail/index?postId=${item.id}`,
      })
    },
  },
}
</script>

<style scoped>
.fav-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.search-shell {
  position: fixed;
  top: 74px;
  left: 0;
  right: 0;
  z-index: 30;
  background: color-mix(in srgb, var(--visago-bg) 92%, #fff 8%);
  border-bottom: 1px solid var(--visago-line);
}

.search-wrap {
  padding: 10px 16px 12px;
  box-sizing: border-box;
}

.search-box {
  height: 46px;
  border-radius: 14px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
  color: var(--visago-text-soft);
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 13px;
}

.search-input {
  flex: 1;
  height: 100%;
  color: var(--visago-text);
  font-size: 14px;
}

.fav-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.fav-content {
  padding: 0 16px calc(34px + var(--visago-safe-bottom));
  box-sizing: border-box;
}

.search-placeholder {
  height: 70px;
}

.fav-list {
  display: grid;
  gap: 12px;
}

.fav-card {
  padding: 12px;
  border-radius: 18px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
  display: flex;
  gap: 12px;
}

.fav-image {
  width: 92px;
  height: 92px;
  border-radius: 14px;
  flex-shrink: 0;
}

.fav-body {
  min-width: 0;
  flex: 1;
}

.fav-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.fav-type {
  padding: 4px 8px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 900;
}

.fav-type--community {
  color: #16a34a;
  background: color-mix(in srgb, #16a34a 14%, var(--visago-surface) 86%);
}

.fav-title,
.fav-desc {
  display: -webkit-box;
  overflow: hidden;
  text-overflow: ellipsis;
  -webkit-box-orient: vertical;
}

.fav-title {
  margin-top: 8px;
  -webkit-line-clamp: 2;
  font-size: 15px;
  line-height: 1.38;
  font-weight: 900;
  color: var(--visago-text);
}

.fav-desc {
  margin-top: 8px;
  -webkit-line-clamp: 2;
  color: var(--visago-text-muted);
  font-size: 12px;
  line-height: 1.45;
}

.fav-icon {
  color: var(--visago-primary);
  font-size: 18px;
  font-variation-settings: 'FILL' 1;
}

.empty-card {
  margin-top: 2px;
  padding: 26px 14px;
  text-align: center;
}

.empty-icon {
  font-size: 30px;
  color: var(--visago-text-soft);
}

.empty-title,
.empty-sub {
  display: block;
}

.empty-title {
  margin-top: 8px;
  font-size: 15px;
  font-weight: 800;
}

.empty-sub {
  margin-top: 5px;
  color: var(--visago-text-muted);
  font-size: 12px;
}
</style>
