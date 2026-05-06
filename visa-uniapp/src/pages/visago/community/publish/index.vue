<template>
  <view class="publish-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="发布帖子" />

    <scroll-view scroll-y class="publish-scroll">
      <view class="publish-content visago-page-width">
        <view class="visago-card form-card">
          <view class="field">
            <view class="field-head">
              <text class="field-label">分类</text>
            </view>
            <picker :range="categories" :value="categoryIndex" @change="onCategoryChange">
              <view class="field-picker">
                <text>{{ categories[categoryIndex] }}</text>
                <text class="material-symbols-outlined">expand_more</text>
              </view>
            </picker>
          </view>

          <view class="field">
            <view class="field-head">
              <text class="field-label">标题</text>
              <text class="field-tip">{{ title.length }}/80</text>
            </view>
            <input v-model="title" class="field-input" maxlength="80" placeholder="写一个清晰的标题" />
          </view>

          <view class="field">
            <view class="field-head">
              <text class="field-label">正文</text>
              <text class="field-tip">{{ content.length }}/2000</text>
            </view>
            <textarea v-model="content" class="field-textarea" maxlength="2000" placeholder="分享你的签证经验、避坑提示或材料建议..." />
          </view>

          <view class="field">
            <view class="field-head">
              <text class="field-label">配图</text>
              <text class="field-tip">可选，最多 9 张</text>
            </view>
            <view class="image-grid">
              <view
                v-for="(image, index) in imagePreviews"
                :key="`${image}-${index}`"
                class="image-cell"
              >
                <image class="image-preview" :src="image" mode="aspectFill" />
                <view class="image-remove" @tap.stop="removeImage(index)">
                  <text class="material-symbols-outlined">close</text>
                </view>
              </view>

              <view v-if="imagePreviews.length < 9" class="image-cell image-cell--add" @tap="chooseImages">
                <text class="material-symbols-outlined">add_photo_alternate</text>
                <text>添加图片</text>
              </view>
            </view>
          </view>
        </view>

        <view class="visago-card guideline-card">
          <text class="guideline-title">发布前提示</text>
          <text class="guideline-text">标题最多 80 字，正文最多 2000 字。图片不是必填项，不上传图片也可以直接发帖。含敏感词的内容会被拦截或进入审核状态。</text>
        </view>
      </view>
    </scroll-view>

    <view class="submit-bar">
      <view class="submit-btn" :class="{ 'submit-btn--disabled': submitting }" @tap="submitPost">
        {{ submitting ? '提交中...' : '发布帖子' }}
      </view>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { api } from '../../../../utils/api'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'

const COMMUNITY_CATEGORIES = ['推荐', '攻略', '问答', '签证经验', '材料模板']

export default {
  components: { VisagoTopBar },
  data() {
    return {
      categories: COMMUNITY_CATEGORIES,
      categoryIndex: 0,
      title: '',
      content: '',
      imageFilePaths: [],
      imagePreviews: [],
      submitting: false,
    }
  },
  onShow() {
    applyTheme(getStoredTheme())
  },
  methods: {
    onCategoryChange(event) {
      this.categoryIndex = Number(event.detail.value || 0)
    },
    chooseImages() {
      const remain = Math.max(0, 9 - this.imagePreviews.length)
      if (!remain) return
      uni.chooseImage({
        count: remain,
        sizeType: ['compressed'],
        sourceType: ['album', 'camera'],
        success: (res) => {
          const next = Array.isArray(res.tempFilePaths) ? res.tempFilePaths : []
          this.imageFilePaths = [...this.imageFilePaths, ...next].slice(0, 9)
          this.imagePreviews = [...this.imageFilePaths]
        },
      })
    },
    removeImage(index) {
      this.imageFilePaths.splice(index, 1)
      this.imagePreviews.splice(index, 1)
    },
    async submitPost() {
      if (this.submitting) return
      const title = this.title.trim()
      const content = this.content.trim()
      if (!title || !content) {
        uni.showToast({ title: '标题和正文必填', icon: 'none' })
        return
      }
      if (title.length > 80) {
        uni.showToast({ title: '标题最多 80 字', icon: 'none' })
        return
      }
      if (content.length > 2000) {
        uni.showToast({ title: '正文最多 2000 字', icon: 'none' })
        return
      }

      this.submitting = true
      try {
        const images = []
        for (const filePath of this.imageFilePaths) {
          const uploaded = await api.uploadImage(filePath, 'community')
          if (uploaded && uploaded.url) {
            images.push(uploaded.url)
          }
        }
        const post = await api.createCommunityPost({
          category: this.categories[this.categoryIndex],
          title,
          content,
          image: images[0] || '',
          images,
        })
        const review = post && post.status === 'review'
        uni.showToast({
          title: review ? '已提交审核' : '发布成功',
          icon: 'none',
        })
        setTimeout(() => {
          uni.navigateBack({ delta: 1 })
        }, 500)
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '发布失败',
          icon: 'none',
        })
      } finally {
        this.submitting = false
      }
    },
  },
}
</script>

<style scoped>
.publish-page {
  min-height: 100vh;
  background: var(--visago-bg);
}

.publish-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: calc(84px + var(--visago-safe-bottom));
  left: 0;
  height: auto;
}

.publish-content {
  padding: 18px 16px 24px;
  box-sizing: border-box;
}

.form-card,
.guideline-card {
  padding: 14px;
}

.guideline-card {
  margin-top: 14px;
}

.field + .field {
  margin-top: 16px;
}

.field-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.field-label {
  display: block;
  font-size: 13px;
  font-weight: 700;
  color: var(--visago-text);
}

.field-tip {
  font-size: 11px;
  color: var(--visago-text-soft);
}

.field-picker,
.field-input,
.field-textarea {
  width: 100%;
  box-sizing: border-box;
  border-radius: 14px;
  background: var(--visago-surface-soft);
}

.field-picker {
  height: 46px;
  padding: 0 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.field-input {
  height: 46px;
  padding: 0 14px;
  font-size: 15px;
}

.field-textarea {
  min-height: 180px;
  padding: 14px;
  font-size: 15px;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.image-cell {
  position: relative;
  aspect-ratio: 1 / 1;
  border-radius: 16px;
  overflow: hidden;
  background: var(--visago-surface-soft);
}

.image-cell--add {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  color: var(--visago-text-soft);
  font-size: 12px;
}

.image-preview {
  width: 100%;
  height: 100%;
}

.image-remove {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 26px;
  height: 26px;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.6);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.guideline-title,
.guideline-text {
  display: block;
}

.guideline-title {
  font-size: 14px;
  font-weight: 700;
}

.guideline-text {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.6;
  color: var(--visago-text-muted);
}

.submit-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 12px 16px calc(12px + var(--visago-safe-bottom));
  background: color-mix(in srgb, var(--visago-surface) 94%, transparent);
  border-top: 1px solid var(--visago-line);
  backdrop-filter: blur(14px);
}

.submit-btn {
  height: 50px;
  border-radius: 16px;
  background: var(--visago-primary);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
}

.submit-btn--disabled {
  opacity: 0.6;
}
</style>
