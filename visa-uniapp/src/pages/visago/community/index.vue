<template>
  <view class="community-page">
    <VisagoTopBar page-name="社区" />

    <view class="community-filter-shell">
      <view class="visago-page-width community-filter-inner">
        <scroll-view class="category-scroll" scroll-x>
          <view class="category-row">
            <view
              v-for="item in categories"
              :key="item"
              class="visago-chip"
              :class="{ 'visago-chip--active': selectedCategory === item }"
              @tap="selectCategory(item)"
            >
              {{ item }}
            </view>
          </view>
        </scroll-view>
      </view>
    </view>

    <view class="community-support-shell">
      <view class="visago-page-width community-support-strip" @tap="copyCommunityContact">
        <text class="material-symbols-outlined community-support-icon">support_agent</text>
        <text class="community-support-text">举报或申诉联系 jiaanliao23@gmail.com</text>
        <text class="community-support-action">复制</text>
      </view>
    </view>

    <scroll-view scroll-y class="community-scroll">
      <view class="community-content visago-page-width">
        <view v-if="loading" class="empty-card visago-card">
          <text class="material-symbols-outlined">hourglass_top</text>
          <text>正在加载社区内容</text>
        </view>

        <view v-else-if="posts.length" class="masonry-grid">
          <view v-for="post in posts" :key="post.id" class="post-card visago-card" @tap="openDetail(post)">
            <image v-if="coverImage(post)" class="post-image" :src="coverImage(post)" mode="aspectFill" />
            <view class="post-body">
              <view class="post-head">
                <text class="post-category">{{ post.category }}</text>
                <view class="post-more" @tap.stop="openPostActions(post)">
                  <text class="material-symbols-outlined">more_horiz</text>
                </view>
              </view>
              <text class="post-title visago-ellipsis-2">{{ post.title }}</text>
              <text class="post-content visago-ellipsis-2">{{ post.content }}</text>
              <view class="post-meta">
                <view class="author-wrap">
                  <image v-if="post.author.avatar" class="avatar" :src="post.author.avatar" mode="aspectFill" />
                  <view v-else class="avatar avatar--fallback">{{ firstLetter(post.author.nickname) }}</view>
                  <text class="author">{{ post.author.nickname }}</text>
                </view>
                <view class="like-wrap" @tap.stop="toggleLike(post)">
                  <text class="material-symbols-outlined like-icon" :class="{ 'like-icon--active': post.liked }">favorite</text>
                  <text class="like-count">{{ post.likeCount }}</text>
                </view>
              </view>
            </view>
          </view>
        </view>

        <view v-else class="empty-card visago-card">
          <text class="material-symbols-outlined">forum</text>
          <text>还没有社区内容</text>
        </view>
      </view>
    </scroll-view>

    <view class="floating-tool visago-card">
      <view class="floating-profile" @tap="openMine">
        <view v-if="currentUser.avatar" class="floating-avatar-wrap">
          <image class="floating-avatar" :src="currentUser.avatar" mode="aspectFill" />
        </view>
        <view v-else class="floating-avatar-wrap floating-avatar-wrap--fallback">
          <text class="material-symbols-outlined floating-avatar-icon">person</text>
        </view>
        <text class="floating-profile-label">我</text>
      </view>

      <view class="floating-add" @tap.stop="openPublish">
        <text class="material-symbols-outlined floating-add-icon">add</text>
      </view>
    </view>

    <view v-if="actionSheetVisible" class="action-sheet-mask" @tap="closePostActions">
      <view class="action-sheet visago-card" @tap.stop>
        <view class="action-sheet__head">
          <text class="action-sheet__title">帖子操作</text>
          <view class="action-sheet__close" @tap="closePostActions">
            <text class="material-symbols-outlined">close</text>
          </view>
        </view>
        <view
          v-for="item in actionItems"
          :key="item.key"
          class="action-sheet__item"
          @tap="handleActionItem(item.key)"
        >
          <text class="action-sheet__label">{{ item.label }}</text>
        </view>
      </view>
    </view>

    <VisagoBottomNav active-tab="community" />
  </view>
</template>

<script>
import VisagoTopBar from '../../../components/VisagoTopBar.vue'
import VisagoBottomNav from '../../../components/VisagoBottomNav.vue'
import { api } from '../../../utils/api'
import { getAuthUser } from '../../../utils/auth'
import { applyTheme, getStoredTheme } from '../../../utils/theme'

const COMMUNITY_CATEGORIES = ['全部', '推荐', '攻略', '问答', '签证经验', '材料模板']

export default {
  components: {
    VisagoTopBar,
    VisagoBottomNav,
  },
  data() {
    return {
      loading: false,
      selectedCategory: '全部',
      categories: COMMUNITY_CATEGORIES,
      posts: [],
      actionSheetVisible: false,
      activeActionPost: null,
      ownActionItems: [{ key: 'delete', label: '删除帖子' }],
      normalActionItems: [
        { key: 'report', label: '举报内容' },
        { key: 'block', label: '屏蔽作者' },
      ],
    }
  },
  computed: {
    currentUser() {
      return getAuthUser() || {}
    },
    actionItems() {
      const post = this.activeActionPost
      if (!post) return []
      const currentUser = getAuthUser()
      const ownPost = Number(currentUser && currentUser.id) === Number(post.userId)
      return ownPost ? this.ownActionItems : this.normalActionItems
    },
  },
  onShow() {
    applyTheme(getStoredTheme())
    this.loadPosts()
  },
  methods: {
    async loadPosts() {
      this.loading = true
      try {
        const params = {}
        if (this.selectedCategory !== '全部') {
          params.category = this.selectedCategory
        }
        const items = await api.listCommunityPosts(params)
        this.posts = Array.isArray(items) ? items : []
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
    selectCategory(category) {
      if (this.selectedCategory === category) return
      this.selectedCategory = category
      this.loadPosts()
    },
    coverImage(post) {
      const images = Array.isArray(post && post.images) ? post.images.filter(Boolean) : []
      if (images.length) return images[0]
      return post && post.image ? post.image : ''
    },
    firstLetter(name) {
      return String(name || 'V').slice(0, 1)
    },
    openPublish() {
      uni.navigateTo({ url: '/pages/visago/community/publish/index' })
    },
    openMine() {
      uni.navigateTo({ url: '/pages/visago/community/mine/index' })
    },
    openDetail(post) {
      uni.navigateTo({
        url: `/pages/visago/community/detail/index?postId=${post.id}`,
      })
    },
    async toggleLike(post) {
      try {
        const item = post.liked ? await api.unlikeCommunityPost(post.id) : await api.likeCommunityPost(post.id)
        const index = this.posts.findIndex((entry) => entry.id === post.id)
        if (index >= 0) {
          this.posts.splice(index, 1, item)
        }
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '操作失败',
          icon: 'none',
        })
      }
    },
    openPostActions(post) {
      this.activeActionPost = post
      this.actionSheetVisible = true
    },
    closePostActions() {
      this.actionSheetVisible = false
      this.activeActionPost = null
    },
    async handleActionItem(action) {
      const post = this.activeActionPost
      if (!post) return
      this.closePostActions()
      if (action === 'delete') {
        await this.deletePost(post)
        return
      }
      if (action === 'report') {
        await this.reportPost(post)
        return
      }
      if (action === 'block') {
        await this.blockAuthor(post)
      }
    },
    async reportPost(post) {
      try {
        await api.reportCommunityPost(post.id, { reason: '不当内容', detail: '' })
        uni.showToast({ title: '已举报', icon: 'none' })
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '举报失败',
          icon: 'none',
        })
      }
    },
    async blockAuthor(post) {
      try {
        await api.blockCommunityUser(post.author.id)
        uni.showToast({ title: '已屏蔽该作者', icon: 'none' })
        this.loadPosts()
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '屏蔽失败',
          icon: 'none',
        })
      }
    },
    async deletePost(post) {
      try {
        await api.deleteCommunityPost(post.id)
        this.posts = this.posts.filter((item) => item.id !== post.id)
        uni.showToast({ title: '已删除', icon: 'none' })
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '删除失败',
          icon: 'none',
        })
      }
    },
    copyCommunityContact() {
      uni.setClipboardData({
        data: 'jiaanliao23@gmail.com',
        fail: () => {
          uni.showToast({ title: '复制失败', icon: 'none' })
        },
      })
    },
  },
}
</script>

<style scoped>
.community-page {
  min-height: 100vh;
  background: var(--visago-bg);
}

.community-filter-shell {
  position: fixed;
  top: 74px;
  left: 0;
  right: 0;
  z-index: 50;
  background: color-mix(in srgb, var(--visago-surface) 94%, transparent);
  border-bottom: 1px solid var(--visago-line);
  backdrop-filter: blur(12px);
}

.community-filter-inner {
  padding: 12px 16px;
  box-sizing: border-box;
}

.community-support-shell {
  position: fixed;
  top: 122px;
  left: 0;
  right: 0;
  z-index: 49;
  background: color-mix(in srgb, var(--visago-surface) 96%, transparent);
  border-bottom: 1px solid var(--visago-line);
}

.community-support-strip {
  min-height: 38px;
  padding: 0 16px;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  gap: 8px;
}

.community-support-icon,
.community-support-text,
.community-support-action {
  display: block;
}

.community-support-icon {
  font-size: 16px;
  color: var(--visago-primary);
}

.community-support-text {
  min-width: 0;
  flex: 1;
  font-size: 12px;
  color: var(--visago-text-soft);
}

.community-support-action {
  flex-shrink: 0;
  font-size: 12px;
  font-weight: 700;
  color: var(--visago-primary);
}

.category-scroll {
  width: 100%;
  white-space: nowrap;
}

.category-row {
  display: inline-flex;
  gap: 10px;
  padding-right: 8px;
}

.community-scroll {
  position: fixed;
  top: 160px;
  right: 0;
  bottom: calc(94px + var(--visago-safe-bottom));
  left: 0;
  height: auto;
}

.community-content {
  padding: 16px;
  box-sizing: border-box;
}

.masonry-grid {
  column-count: 2;
  column-gap: 12px;
}

.post-card {
  overflow: hidden;
  margin-bottom: 12px;
  break-inside: avoid;
  display: inline-block;
  width: 100%;
}

.post-image {
  width: 100%;
  height: 220px;
  display: block;
}

.post-body {
  padding: 12px;
}

.post-head {
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

.post-more {
  color: var(--visago-text-soft);
}

.post-title {
  margin-top: 10px;
  font-size: 34rpx;
  line-height: 1.35;
  font-weight: 700;
}

.post-content {
  display: block;
  margin-top: 8px;
  color: var(--visago-text-muted);
  font-size: 13px;
  line-height: 1.55;
}

.post-meta {
  margin-top: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.author-wrap {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}

.avatar {
  width: 20px;
  height: 20px;
  border-radius: 9999px;
  flex-shrink: 0;
}

.avatar--fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--visago-surface-soft);
  color: var(--visago-primary);
  font-size: 11px;
  font-weight: 700;
}

.author {
  font-size: 15px;
  color: var(--visago-text-muted);
}

.like-wrap {
  display: flex;
  align-items: center;
  gap: 2px;
}

.like-icon {
  font-size: 16px;
  color: var(--visago-text-soft);
}

.like-icon--active {
  color: #ef4444;
  font-variation-settings: 'FILL' 1;
}

.like-count {
  font-size: 14px;
  color: var(--visago-text-muted);
}

.empty-card {
  padding: 28px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: var(--visago-text-muted);
}

.floating-tool {
  position: fixed;
  right: 20px;
  bottom: calc(124px + var(--visago-safe-bottom));
  z-index: 110;
  width: 76px;
  padding: 12px 10px 10px;
  border-radius: 999px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.floating-profile {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.floating-avatar-wrap {
  width: 40px;
  height: 40px;
  border-radius: 999px;
  overflow: hidden;
  background: #eef2ff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.floating-avatar-wrap--fallback {
  background: #e8edff;
}

.floating-avatar {
  width: 100%;
  height: 100%;
}

.floating-avatar-icon {
  font-size: 22px;
  color: #a9b4d6;
}

.floating-profile-label {
  font-size: 13px;
  color: var(--visago-text-muted);
  line-height: 1;
}

.floating-add {
  width: 56px;
  height: 56px;
  border-radius: 999px;
  background: linear-gradient(180deg, #56e8a1 0%, #30d39a 100%);
  box-shadow: 0 12px 22px rgba(48, 211, 154, 0.28);
  display: flex;
  align-items: center;
  justify-content: center;
}

.floating-add-icon {
  font-size: 34px;
  color: #fff;
  font-variation-settings: 'wght' 500;
}

.action-sheet-mask {
  position: fixed;
  inset: 0;
  z-index: 10020;
  background: rgba(15, 23, 42, 0.34);
  display: flex;
  align-items: flex-end;
}

.action-sheet {
  width: 100%;
  border-radius: 24px 24px 0 0;
  padding: 16px 16px calc(18px + var(--visago-safe-bottom));
  box-sizing: border-box;
}

.action-sheet__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.action-sheet__title {
  font-size: 16px;
  font-weight: 700;
  color: var(--visago-text);
}

.action-sheet__close {
  width: 34px;
  height: 34px;
  border-radius: 999px;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--visago-text-soft);
}

.action-sheet__item {
  margin-top: 12px;
  height: 50px;
  border-radius: 14px;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-sheet__label {
  font-size: 15px;
  font-weight: 600;
  color: var(--visago-text);
}
</style>
