<template>
  <view class="mine-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="我的帖子" />

    <view class="mine-filter-shell">
      <view class="visago-page-width mine-filter-inner">
        <scroll-view class="filter-scroll" scroll-x>
          <view class="filter-row">
            <view
              v-for="item in filters"
              :key="item.key"
              class="visago-chip"
              :class="{ 'visago-chip--active': activeStatus === item.key }"
              @tap="selectStatus(item.key)"
            >
              {{ item.label }}
            </view>
          </view>
        </scroll-view>
      </view>
    </view>

    <scroll-view scroll-y class="mine-scroll">
      <view class="mine-content visago-page-width">
        <view v-if="loading" class="empty-card visago-card">
          <text class="material-symbols-outlined">hourglass_top</text>
          <text>正在加载我的帖子</text>
        </view>

        <view v-else-if="posts.length" class="post-list">
          <view v-for="post in posts" :key="post.id" class="visago-card post-card" @tap="openDetail(post)">
            <view class="post-top">
              <text class="post-category">{{ post.category }}</text>
              <text class="post-status" :class="`post-status--${post.status}`">{{ statusLabel(post.status) }}</text>
            </view>
            <text class="post-title visago-ellipsis-2">{{ post.title }}</text>
            <text class="post-content visago-ellipsis-2">{{ post.content }}</text>
            <text v-if="post.reviewNote" class="post-note">原因：{{ post.reviewNote }}</text>
            <view class="post-bottom">
              <text class="post-time">{{ formatDateTime(post.createdAt) }}</text>
              <view class="post-metrics">
                <text>{{ post.likeCount }} 赞</text>
                <text>{{ post.commentCount }} 评论</text>
              </view>
            </view>
          </view>
        </view>

        <view v-else class="empty-card visago-card">
          <text class="material-symbols-outlined">forum</text>
          <text>还没有匹配的帖子</text>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { api } from '../../../../utils/api'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'

const FILTERS = [
  { key: '', label: '全部' },
  { key: 'review', label: '待审核' },
  { key: 'published', label: '已发布' },
  { key: 'hidden', label: '已隐藏' },
  { key: 'rejected', label: '已驳回' },
]

export default {
  components: { VisagoTopBar },
  data() {
    return {
      filters: FILTERS,
      activeStatus: '',
      loading: false,
      posts: [],
    }
  },
  onShow() {
    applyTheme(getStoredTheme())
    this.loadPosts()
  },
  methods: {
    async loadPosts() {
      this.loading = true
      try {
        this.posts = await api.listMyCommunityPosts({ status: this.activeStatus })
      } catch (error) {
        this.posts = []
        uni.showToast({
          title: (error && error.message) || '加载失败',
          icon: 'none',
        })
      } finally {
        this.loading = false
      }
    },
    selectStatus(status) {
      if (this.activeStatus === status) return
      this.activeStatus = status
      this.loadPosts()
    },
    statusLabel(status) {
      switch (status) {
        case 'published':
          return '已发布'
        case 'hidden':
          return '已隐藏'
        case 'rejected':
          return '已驳回'
        default:
          return '待审核'
      }
    },
    formatDateTime(value) {
      const date = new Date(value)
      if (Number.isNaN(date.getTime())) return String(value || '')
      const month = `${date.getMonth() + 1}`.padStart(2, '0')
      const day = `${date.getDate()}`.padStart(2, '0')
      const hours = `${date.getHours()}`.padStart(2, '0')
      const minutes = `${date.getMinutes()}`.padStart(2, '0')
      return `${month}-${day} ${hours}:${minutes}`
    },
    openDetail(post) {
      uni.navigateTo({
        url: `/pages/visago/community/detail/index?postId=${post.id}`,
      })
    },
  },
}
</script>

<style scoped>
.mine-page {
  min-height: 100vh;
  background: var(--visago-bg);
}

.mine-filter-shell {
  position: fixed;
  top: 74px;
  left: 0;
  right: 0;
  z-index: 50;
  background: color-mix(in srgb, var(--visago-surface) 94%, transparent);
  border-bottom: 1px solid var(--visago-line);
  backdrop-filter: blur(12px);
}

.mine-filter-inner {
  padding: 12px 16px;
  box-sizing: border-box;
}

.filter-scroll {
  width: 100%;
  white-space: nowrap;
}

.filter-row {
  display: inline-flex;
  gap: 10px;
  padding-right: 8px;
}

.mine-scroll {
  position: fixed;
  top: 132px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.mine-content {
  padding: 16px;
  box-sizing: border-box;
}

.post-list {
  display: grid;
  gap: 12px;
}

.post-card {
  padding: 14px;
}

.post-top,
.post-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.post-category {
  padding: 4px 8px;
  border-radius: 999px;
  background: rgba(15, 101, 216, 0.12);
  color: var(--visago-primary);
  font-size: 11px;
  font-weight: 700;
}

.post-status {
  font-size: 12px;
  font-weight: 700;
}

.post-status--published {
  color: #16a34a;
}

.post-status--hidden {
  color: #d97706;
}

.post-status--rejected {
  color: #ef4444;
}

.post-status--review {
  color: var(--visago-primary);
}

.post-title,
.post-content,
.post-note {
  display: block;
}

.post-title {
  margin-top: 12px;
  font-size: 16px;
  font-weight: 700;
  color: var(--visago-text);
}

.post-content {
  margin-top: 8px;
  color: var(--visago-text-muted);
  font-size: 13px;
  line-height: 1.55;
}

.post-note {
  margin-top: 10px;
  padding: 10px 12px;
  border-radius: 12px;
  background: rgba(239, 68, 68, 0.08);
  color: #b91c1c;
  font-size: 12px;
  line-height: 1.5;
}

.post-bottom {
  margin-top: 12px;
}

.post-time,
.post-metrics {
  color: var(--visago-text-soft);
  font-size: 12px;
}

.post-metrics {
  display: flex;
  gap: 12px;
}

.empty-card {
  padding: 28px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: var(--visago-text-muted);
}
</style>
