<template>
  <view class="detail-page" :class="{ 'detail-page--dark': themeMode === 'dark' }">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="帖子详情" />

    <scroll-view scroll-y class="detail-scroll">
      <view class="detail-content visago-page-width">
        <view v-if="loading" class="empty-block">
          <text class="material-symbols-outlined">hourglass_top</text>
          <text>正在加载帖子详情</text>
        </view>

        <template v-else-if="post">
          <view class="hero-wrap">
            <swiper
              v-if="displayImages.length"
              class="hero-swiper"
              :indicator-dots="displayImages.length > 1"
              :autoplay="displayImages.length > 1"
              :interval="3200"
              :duration="320"
              circular
              @change="onSwiperChange"
            >
              <swiper-item v-for="(image, index) in displayImages" :key="`${image}-${index}`">
                <image class="hero-image" :src="image" mode="aspectFill" @tap="previewImages(index)" />
              </swiper-item>
            </swiper>

            <view v-if="displayImages.length > 1" class="hero-counter">{{ currentImageIndex + 1 }}/{{ displayImages.length }}</view>

            <view class="post-panel">
              <view class="post-head">
                <text class="post-category">{{ post.category }}</text>
                <text class="post-time">{{ formatTime(post.createdAt) }}</text>
              </view>

              <view v-if="post.status !== 'published'" class="status-banner" :class="`status-banner--${post.status}`">
                <text class="status-banner__title">{{ statusLabel(post.status) }}</text>
                <text v-if="post.reviewNote" class="status-banner__desc">原因：{{ post.reviewNote }}</text>
              </view>

              <text class="post-title">{{ post.title }}</text>

              <view class="author-row">
                <image v-if="post.author.avatar" class="avatar" :src="post.author.avatar" mode="aspectFill" />
                <view v-else class="avatar avatar--fallback">{{ firstLetter(post.author.nickname) }}</view>
                <view class="author-copy">
                  <text class="author">{{ post.author.nickname }}</text>
                  <text class="author-sub">发布于 {{ formatDateTime(post.createdAt) }}</text>
                </view>
              </view>

              <view class="section-divider" />
              <text class="post-content">{{ post.content }}</text>

              <view class="section-divider section-divider--compact" />
              <view class="post-footer">
                <view class="metric">
                  <text class="material-symbols-outlined metric-icon">favorite</text>
                  <text>{{ post.likeCount }}</text>
                </view>
                <view class="metric">
                  <text class="material-symbols-outlined metric-icon">chat_bubble</text>
                  <text>{{ post.commentCount }}</text>
                </view>
              </view>
            </view>
          </view>

          <view class="comment-head">
            <text class="comment-title">评论区</text>
            <text class="comment-count">{{ visibleCommentCount }} 条</text>
          </view>

          <view class="comment-support" @tap="copyCommunityContact">
            <text class="material-symbols-outlined comment-support__icon">support_agent</text>
            <text class="comment-support__text">举报或申诉请联系 jiaanliao23@gmail.com</text>
            <text class="comment-support__action">复制</text>
          </view>

          <view v-if="comments.length" class="comment-list">
            <view v-for="comment in comments" :key="comment.id" class="comment-card">
              <view class="comment-main">
                <view class="comment-author-row">
                  <image v-if="comment.author.avatar" class="avatar" :src="comment.author.avatar" mode="aspectFill" />
                  <view v-else class="avatar avatar--fallback">{{ firstLetter(comment.author.nickname) }}</view>
                  <view class="author-copy">
                    <text class="author">{{ comment.author.nickname }}</text>
                    <text class="author-sub">{{ formatDateTime(comment.createdAt) }}</text>
                  </view>
                </view>

                <view v-if="comment.status !== 'published'" class="comment-status" :class="`comment-status--${comment.status}`">
                  <text>{{ commentStatusLabel(comment.status) }}</text>
                  <text v-if="comment.reviewNote" class="comment-status__note">：{{ comment.reviewNote }}</text>
                </view>

                <text v-if="comment.content" class="comment-content">{{ comment.content }}</text>
                <image v-if="comment.image" class="comment-image" :src="comment.image" mode="aspectFill" @tap="previewSingle(comment.image)" />

                <view class="comment-actions">
                  <text class="comment-action" @tap="startReply(comment)">回复</text>
                  <text v-if="canDelete(comment)" class="comment-action comment-action--danger" @tap="removeComment(comment)">删除</text>
                  <text v-else-if="canReport(comment)" class="comment-action" @tap="reportComment(comment)">举报</text>
                  <text
                    v-if="comment.replies && comment.replies.length"
                    class="comment-action comment-action--toggle"
                    @tap="toggleReplies(comment.id)"
                  >
                    {{ isRepliesExpanded(comment.id) ? '收起回复' : `展开 ${comment.replies.length} 条回复` }}
                  </text>
                </view>
              </view>

              <view v-if="comment.replies && comment.replies.length && isRepliesExpanded(comment.id)" class="reply-list">
                <view v-for="reply in comment.replies" :key="reply.id" class="reply-item">
                  <view class="comment-author-row">
                    <image v-if="reply.author.avatar" class="avatar avatar--small" :src="reply.author.avatar" mode="aspectFill" />
                    <view v-else class="avatar avatar--fallback avatar--small">{{ firstLetter(reply.author.nickname) }}</view>
                    <view class="author-copy">
                      <text class="author">
                        {{ reply.author.nickname }}
                        <text v-if="reply.replyToNickname" class="reply-target"> 回复 {{ reply.replyToNickname }}</text>
                      </text>
                      <text class="author-sub">{{ formatDateTime(reply.createdAt) }}</text>
                    </view>
                  </view>

                  <view v-if="reply.status !== 'published'" class="comment-status" :class="`comment-status--${reply.status}`">
                    <text>{{ commentStatusLabel(reply.status) }}</text>
                    <text v-if="reply.reviewNote" class="comment-status__note">：{{ reply.reviewNote }}</text>
                  </view>

                  <text v-if="reply.content" class="comment-content">{{ reply.content }}</text>
                  <image v-if="reply.image" class="comment-image comment-image--reply" :src="reply.image" mode="aspectFill" @tap="previewSingle(reply.image)" />

                  <view class="comment-actions">
                    <text class="comment-action" @tap="startReply(reply)">回复</text>
                    <text v-if="canDelete(reply)" class="comment-action comment-action--danger" @tap="removeComment(reply)">删除</text>
                    <text v-else-if="canReport(reply)" class="comment-action" @tap="reportComment(reply)">举报</text>
                  </view>
                </view>
              </view>
            </view>
          </view>

          <view v-else class="empty-block">
            <text class="material-symbols-outlined">chat</text>
            <text>还没有评论，来抢个沙发吧</text>
          </view>
        </template>
      </view>
    </scroll-view>

    <view v-if="!composerVisible" class="compact-bar">
      <view class="compact-input" @tap="openComposer">
        <text class="material-symbols-outlined compact-input__icon">edit</text>
        <text class="compact-input__placeholder">{{ replyTarget ? `回复 ${replyTarget.author.nickname}` : '说点什么...' }}</text>
      </view>

      <view class="compact-metric" @tap="toggleLike">
        <text class="material-symbols-outlined compact-metric__icon" :class="{ 'compact-metric__icon--active': post && post.liked }">favorite</text>
        <text>{{ post ? post.likeCount : 0 }}</text>
      </view>

      <view class="compact-metric" @tap="toggleFavorite">
        <text class="material-symbols-outlined compact-metric__icon" :class="{ 'compact-metric__icon--favorite': favorite }">bookmark</text>
        <text>{{ favorite ? '已藏' : '收藏' }}</text>
      </view>

      <view class="compact-metric">
        <text class="material-symbols-outlined compact-metric__icon">chat_bubble</text>
        <text>{{ post ? post.commentCount : 0 }}</text>
      </view>
    </view>

    <view v-if="composerVisible" class="composer-mask" @tap="closeComposer">
      <view class="composer-panel" @tap.stop>
        <textarea v-model.trim="commentText" class="composer-textarea" maxlength="50" placeholder="说点什么..." />

        <view v-if="commentImage" class="composer-image-wrap">
          <image class="composer-image" :src="commentImage" mode="aspectFill" />
          <view class="composer-image-remove" @tap="clearCommentImage">
            <text class="material-symbols-outlined">close</text>
          </view>
        </view>

        <view class="composer-toolbar">
          <view class="composer-tools">
            <view class="composer-tool" @tap="chooseCommentImage">
              <text class="material-symbols-outlined">image</text>
            </view>
          </view>

          <view class="composer-right">
            <text class="composer-limit">{{ commentLength }}/50</text>
            <view class="composer-submit" :class="{ 'composer-submit--disabled': submittingComment }" @tap="submitComment">发送</view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { api } from '../../../../utils/api'
import { getAuthUser, getAuthToken, isLoggedIn, LOGIN_PAGE } from '../../../../utils/auth'
import { applyTheme, getStoredTheme, THEME_CHANGE_EVENT } from '../../../../utils/theme'

const API_BASE = (import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8080/api').replace(/\/+$/, '')

function reportCommentFallback(commentId, payload = {}) {
  const token = getAuthToken()
  return new Promise((resolve, reject) => {
    uni.request({
      url: `${API_BASE}/community/comments/${commentId}/report`,
      method: 'POST',
      data: payload,
      header: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      success: (res) => {
        const statusCode = Number(res.statusCode || 0)
        const body = res.data || {}
        if (statusCode >= 200 && statusCode < 300) {
          resolve(body.data)
          return
        }
        reject(new Error(body.message || `Request failed (${statusCode})`))
      },
      fail: (err) => {
        reject(new Error((err && err.errMsg) || 'Network request failed'))
      },
    })
  })
}

export default {
  components: { VisagoTopBar },
  data() {
    return {
      postId: 0,
      themeMode: 'light',
      loading: false,
      submittingComment: false,
      post: null,
      comments: [],
      commentText: '',
      commentImage: '',
      replyTarget: null,
      currentImageIndex: 0,
      favorite: false,
      replyExpandedMap: {},
      composerVisible: false,
    }
  },
  computed: {
    displayImages() {
      const images = this.post && Array.isArray(this.post.images) ? this.post.images.filter(Boolean) : []
      if (images.length) return images
      const image = this.post && this.post.image ? [this.post.image] : []
      return image
    },
    commentLength() {
      return this.commentText.length
    },
    visibleCommentCount() {
      return this.comments.reduce((total, item) => total + 1 + (Array.isArray(item.replies) ? item.replies.length : 0), 0)
    },
  },
  onLoad(query) {
    this.postId = Number(query.postId || 0)
    if (typeof uni !== 'undefined' && uni.$on) {
      uni.$on(THEME_CHANGE_EVENT, this.onThemeChange)
    }
  },
  onShow() {
    this.syncTheme()
    this.loadDetail()
  },
  onUnload() {
    if (typeof uni !== 'undefined' && uni.$off) {
      uni.$off(THEME_CHANGE_EVENT, this.onThemeChange)
    }
  },
  methods: {
    syncTheme() {
      this.themeMode = applyTheme(getStoredTheme())
    },
    onThemeChange(theme) {
      this.themeMode = theme === 'dark' ? 'dark' : 'light'
      applyTheme(this.themeMode)
    },
    async loadDetail() {
      if (!this.postId) return
      this.loading = true
      try {
        const [post, comments] = await Promise.all([
          api.getCommunityPost(this.postId),
          api.listCommunityComments(this.postId),
        ])
        this.post = post || null
        this.comments = Array.isArray(comments) ? comments : []
        this.replyExpandedMap = {}
        this.favorite = !!(this.post && this.post.favorited)
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '加载失败',
          icon: 'none',
        })
      } finally {
        this.loading = false
      }
    },
    firstLetter(name) {
      return String(name || 'V').slice(0, 1)
    },
    formatTime(value) {
      const date = new Date(value)
      if (Number.isNaN(date.getTime())) return ''
      const month = `${date.getMonth() + 1}`.padStart(2, '0')
      const day = `${date.getDate()}`.padStart(2, '0')
      return `${month}-${day}`
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
    statusLabel(status) {
      switch (status) {
        case 'hidden':
          return '当前帖子已被隐藏'
        case 'rejected':
          return '当前帖子已被驳回'
        case 'review':
          return '当前帖子正在审核中'
        default:
          return '已发布'
      }
    },
    commentStatusLabel(status) {
      switch (status) {
        case 'hidden':
          return '评论已隐藏'
        case 'rejected':
          return '评论已驳回'
        case 'review':
          return '评论审核中'
        default:
          return '已发布'
      }
    },
    previewImages(current) {
      const urls = this.displayImages
      if (!urls.length) return
      uni.previewImage({
        current,
        urls,
      })
    },
    previewSingle(url) {
      if (!url) return
      uni.previewImage({
        current: 0,
        urls: [url],
      })
    },
    onSwiperChange(event) {
      this.currentImageIndex = Number(event.detail.current || 0)
    },
    startReply(comment) {
      if (comment.rootId) {
        this.replyExpandedMap = {
          ...this.replyExpandedMap,
          [comment.rootId]: true,
        }
      }
      this.replyTarget = comment
      this.commentText = ''
      this.commentImage = ''
      this.composerVisible = true
    },
    cancelReply() {
      this.replyTarget = null
      this.commentText = ''
      this.commentImage = ''
    },
    canDelete(comment) {
      const user = getAuthUser()
      return Number(user && user.id) === Number(comment.userId)
    },
    canReport(comment) {
      const user = getAuthUser()
      return isLoggedIn() && Number(user && user.id) > 0 && Number(user && user.id) !== Number(comment.userId)
    },
    isRepliesExpanded(commentId) {
      return !!this.replyExpandedMap[commentId]
    },
    toggleReplies(commentId) {
      this.replyExpandedMap = {
        ...this.replyExpandedMap,
        [commentId]: !this.replyExpandedMap[commentId],
      }
    },
    openComposer() {
      this.composerVisible = true
    },
    closeComposer() {
      this.composerVisible = false
    },
    chooseCommentImage() {
      uni.chooseImage({
        count: 1,
        sizeType: ['compressed'],
        sourceType: ['album', 'camera'],
        success: (res) => {
          this.commentImage = (res.tempFilePaths && res.tempFilePaths[0]) || ''
        },
      })
    },
    clearCommentImage() {
      this.commentImage = ''
    },
    async submitComment() {
      if (this.submittingComment || (!this.commentText.trim() && !this.commentImage) || !this.postId) return
      if (this.commentText.length > 50) {
        uni.showToast({ title: '评论最多 50 字', icon: 'none' })
        return
      }
      this.submittingComment = true
      try {
        let image = ''
        if (this.commentImage) {
          const uploaded = await api.uploadImage(this.commentImage, 'community')
          image = uploaded && uploaded.url ? uploaded.url : ''
        }
        const payload = {
          content: this.commentText,
          image,
          parentId: this.replyTarget ? this.replyTarget.id : 0,
        }
        const comment = await api.createCommunityComment(this.postId, payload)
        const review = comment && comment.status === 'review'
        uni.showToast({
          title: review ? '评论已提交审核' : '评论成功',
          icon: 'none',
        })
        if (comment && comment.status === 'published') {
          this.insertComment(comment, this.replyTarget)
          if (this.post) {
            this.post.commentCount = Number(this.post.commentCount || 0) + 1
          }
        } else {
          await this.loadDetail()
        }
        this.commentText = ''
        this.commentImage = ''
        this.replyTarget = null
        this.composerVisible = false
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '评论失败',
          icon: 'none',
        })
      } finally {
        this.submittingComment = false
      }
    },
    insertComment(comment, replyTarget) {
      const next = {
        ...comment,
        replies: Array.isArray(comment.replies) ? comment.replies : [],
      }
      if (!replyTarget || !replyTarget.id) {
        this.comments.push(next)
        return
      }
      const rootId = Number(replyTarget.rootId || replyTarget.id)
      const root = this.comments.find((item) => Number(item.id) === rootId)
      if (!root) {
        this.comments.push(next)
        return
      }
      if (!Array.isArray(root.replies)) {
        root.replies = []
      }
      root.replies.push(next)
      root.replyCount = Number(root.replyCount || 0) + 1
      this.replyExpandedMap = {
        ...this.replyExpandedMap,
        [rootId]: true,
      }
    },
    async removeComment(comment) {
      try {
        await api.deleteCommunityComment(comment.id)
        uni.showToast({ title: '已删除', icon: 'none' })
        await this.loadDetail()
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '删除失败',
          icon: 'none',
        })
      }
    },
    async reportComment(comment) {
      if (!isLoggedIn()) {
        uni.showToast({
          title: '请先登录后再举报',
          icon: 'none',
        })
        setTimeout(() => {
          uni.reLaunch({ url: LOGIN_PAGE })
        }, 180)
        return
      }
      try {
        if (typeof api.reportCommunityComment === 'function') {
          await api.reportCommunityComment(comment.id, { reason: '不当内容', detail: '' })
        } else {
          await reportCommentFallback(comment.id, { reason: '不当内容', detail: '' })
        }
        uni.showToast({ title: '已举报', icon: 'none' })
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '举报失败',
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
    async toggleLike() {
      if (!this.post) return
      try {
        const next = this.post.liked
          ? await api.unlikeCommunityPost(this.post.id)
          : await api.likeCommunityPost(this.post.id)
        this.post = next
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '操作失败',
          icon: 'none',
        })
      }
    },
    async toggleFavorite() {
      if (!this.post) return
      try {
        const next = this.favorite
          ? await api.unfavoriteCommunityPost(this.post.id)
          : await api.favoriteCommunityPost(this.post.id)
        this.post = next
        this.favorite = !!(next && next.favorited)
        uni.showToast({
          title: this.favorite ? '已收藏' : '已取消收藏',
          icon: 'none',
        })
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '操作失败',
          icon: 'none',
        })
      }
    },
  },
}
</script>

<style scoped>
.detail-page {
  --detail-surface: #fcfdff;
  --detail-line: rgba(148, 163, 184, 0.18);
  min-height: 100vh;
  background: var(--detail-surface);
}

.detail-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: calc(72px + var(--visago-safe-bottom));
  left: 0;
  height: auto;
}

.detail-content {
  padding: 18px 16px 24px;
  box-sizing: border-box;
}

.hero-wrap {
  overflow: hidden;
  background: transparent;
}

.hero-swiper {
  width: 100%;
  height: 280px;
}

.hero-image {
  width: 100%;
  height: 100%;
  display: block;
}

.hero-counter {
  position: absolute;
  top: 18px;
  right: 18px;
  z-index: 2;
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.54);
  color: #fff;
  font-size: 11px;
}

.post-panel {
  padding: 18px 4px 10px;
}

.section-divider {
  margin-top: 14px;
  height: 1px;
  background: var(--detail-line);
}

.section-divider--compact {
  margin-top: 16px;
}

.post-head,
.author-row,
.comment-author-row,
.post-footer,
.comment-actions,
.comment-head {
  display: flex;
  align-items: center;
}

.post-head,
.post-footer,
.comment-head {
  justify-content: space-between;
}

.post-category {
  padding: 4px 8px;
  border-radius: 999px;
  background: rgba(15, 101, 216, 0.12);
  color: var(--visago-primary);
  font-size: 11px;
  font-weight: 700;
}

.post-time,
.comment-time,
.comment-count {
  color: var(--visago-text-soft);
  font-size: 11px;
}

.post-title {
  display: block;
  margin-top: 14px;
  font-size: 20px;
  font-weight: 800;
  color: var(--visago-text);
  line-height: 1.45;
}

.author-row,
.comment-author-row {
  gap: 8px;
  margin-top: 14px;
}

.avatar {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  flex-shrink: 0;
}

.avatar--small {
  width: 20px;
  height: 20px;
}

.avatar--fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--visago-surface-soft);
  color: var(--visago-primary);
  font-size: 12px;
  font-weight: 700;
}

.author-copy {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.author {
  color: var(--visago-text);
  font-size: 13px;
  font-weight: 600;
}

.author-sub {
  color: var(--visago-text-soft);
  font-size: 11px;
}

.reply-target {
  color: var(--visago-text-soft);
  font-weight: 400;
}

.post-content,
.comment-content {
  display: block;
  margin-top: 14px;
  color: var(--visago-text);
  font-size: 14px;
  line-height: 1.8;
  white-space: pre-wrap;
}

.comment-image {
  margin-top: 10px;
  width: 120px;
  height: 120px;
  border-radius: 14px;
  display: block;
}

.comment-image--reply {
  width: 96px;
  height: 96px;
}

.comment-status {
  margin-top: 10px;
  font-size: 12px;
  line-height: 1.5;
}

.comment-status--review {
  color: var(--visago-primary);
}

.comment-status--hidden {
  color: #b45309;
}

.comment-status--rejected {
  color: #b91c1c;
}

.comment-status__note {
  opacity: 0.92;
}

.status-banner {
  margin-top: 12px;
  padding: 12px;
  border-radius: 14px;
  display: grid;
  gap: 4px;
}

.status-banner--review {
  background: rgba(15, 101, 216, 0.08);
  color: var(--visago-primary);
}

.status-banner--hidden {
  background: rgba(217, 119, 6, 0.08);
  color: #b45309;
}

.status-banner--rejected {
  background: rgba(239, 68, 68, 0.08);
  color: #b91c1c;
}

.status-banner__title {
  font-size: 13px;
  font-weight: 700;
}

.status-banner__desc {
  font-size: 12px;
  line-height: 1.55;
}

.post-footer {
  margin-top: 16px;
}

.metric {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--visago-text-soft);
  font-size: 12px;
}

.metric-icon {
  font-size: 16px;
}

.comment-head {
  margin-top: 18px;
  margin-bottom: 10px;
  padding-top: 10px;
  border-top: 1px solid var(--detail-line);
}

.comment-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--visago-text);
}

.comment-list {
  display: grid;
  gap: 12px;
}

.comment-support {
  margin: 0 0 12px;
  min-height: 38px;
  padding: 0 10px;
  border-radius: 12px;
  background: color-mix(in srgb, var(--visago-surface) 88%, #ffffff 12%);
  border: 1px solid var(--detail-line);
  display: flex;
  align-items: center;
  gap: 8px;
}

.comment-support__icon,
.comment-support__text,
.comment-support__action {
  display: block;
}

.comment-support__icon {
  font-size: 16px;
  color: var(--visago-primary);
}

.comment-support__text {
  min-width: 0;
  flex: 1;
  font-size: 12px;
  color: var(--visago-text-soft);
}

.comment-support__action {
  flex-shrink: 0;
  font-size: 12px;
  font-weight: 700;
  color: var(--visago-primary);
}

.comment-card {
  padding: 12px 0 0;
  border-bottom: 1px solid var(--detail-line);
}

.comment-main {
  display: grid;
}

.comment-actions {
  gap: 12px;
  margin-top: 10px;
}

.comment-action {
  color: var(--visago-primary);
  font-size: 12px;
}

.comment-action--danger {
  color: #ef4444;
}

.comment-action--toggle {
  margin-left: auto;
}

.reply-list {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--detail-line);
  display: grid;
  gap: 10px;
}

.reply-item {
  padding: 10px 0 0 12px;
  border-left: 2px solid color-mix(in srgb, var(--visago-primary) 18%, #ffffff 82%);
}

.empty-block {
  margin-top: 10px;
  padding: 24px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: var(--visago-text-muted);
}

.compact-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 10px 16px calc(10px + var(--visago-safe-bottom));
  background: color-mix(in srgb, #ffffff 94%, transparent);
  border-top: 1px solid var(--visago-line);
  display: flex;
  align-items: center;
  gap: 14px;
  box-sizing: border-box;
}

.compact-input {
  flex: 1;
  min-width: 0;
  height: 42px;
  padding: 0 14px;
  border-radius: 999px;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  gap: 8px;
}

.compact-input__icon {
  font-size: 18px;
  color: var(--visago-text-soft);
}

.compact-input__placeholder {
  color: var(--visago-text-soft);
  font-size: 14px;
}

.compact-metric {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--visago-text);
  font-size: 13px;
}

.compact-metric__icon {
  font-size: 22px;
  color: var(--visago-text);
}

.compact-metric__icon--active {
  color: #ef4444;
  font-variation-settings: 'FILL' 1;
}

.compact-metric__icon--favorite {
  color: var(--visago-primary);
  font-variation-settings: 'FILL' 1;
}

.composer-mask {
  position: fixed;
  inset: 0;
  z-index: 10020;
  background: rgba(15, 23, 42, 0.24);
  display: flex;
  align-items: flex-end;
}

.composer-panel {
  width: 100%;
  padding: 14px 14px calc(14px + var(--visago-safe-bottom));
  border-radius: 24px 24px 0 0;
  background: #ffffff;
  box-sizing: border-box;
}

.composer-textarea {
  width: 100%;
  min-height: 120px;
  padding: 14px;
  border-radius: 18px;
  background: var(--visago-surface-soft);
  box-sizing: border-box;
  font-size: 14px;
}

.composer-image-wrap {
  position: relative;
  width: 96px;
  height: 96px;
  margin-top: 12px;
}

.composer-image {
  width: 100%;
  height: 100%;
  border-radius: 14px;
}

.composer-image-remove {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 24px;
  height: 24px;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.6);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.composer-toolbar {
  margin-top: 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.composer-tools {
  display: flex;
  align-items: center;
  gap: 14px;
}

.composer-tool {
  color: var(--visago-text-soft);
}

.composer-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.composer-limit {
  font-size: 12px;
  color: var(--visago-text-soft);
}

.composer-submit {
  min-width: 60px;
  height: 34px;
  padding: 0 16px;
  border-radius: 999px;
  background: #fcb9c7;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
}

.composer-submit--disabled {
  opacity: 0.5;
}

.detail-page--dark,
:global(html.theme-dark) .detail-page {
  --detail-surface: var(--visago-bg);
  --detail-line: rgba(148, 163, 184, 0.14);
}

.detail-page--dark .compact-bar,
:global(html.theme-dark) .compact-bar {
  background: color-mix(in srgb, var(--visago-surface) 94%, transparent);
}

.detail-page--dark .composer-panel,
:global(html.theme-dark) .composer-panel {
  background: var(--visago-surface);
}

.detail-page--dark .reply-item,
:global(html.theme-dark) .reply-item {
  border-left-color: rgba(148, 163, 184, 0.22);
}
</style>
