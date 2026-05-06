<template>
  <view class="photo-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="照片指南" />

    <scroll-view scroll-y class="photo-scroll">
      <view class="photo-content visago-page-width">
        <view class="search-panel">
          <view class="search-box">
            <text class="material-symbols-outlined search-icon">search</text>
            <input v-model.trim="keyword" class="search-input" placeholder="搜索国家或签证类型..." />
          </view>

          <picker :range="countryLabels" :value="countryIndex" @change="onCountryChange">
            <view class="country-card">
              <view class="country-left">
                <view class="flag-wrap">{{ currentStandard.flag }}</view>
                <view>
                  <text class="country-name">{{ currentStandard.country }}</text>
                  <text class="country-visa">{{ currentStandard.visa }}</text>
                </view>
              </view>
              <text class="material-symbols-outlined picker-arrow">expand_more</text>
            </view>
          </picker>
        </view>

        <view class="sample-card">
          <image class="sample-image" :src="sampleImage" mode="aspectFill" />
          <view class="sample-mask">
            <view class="sample-badge">
              <text class="material-symbols-outlined sample-badge-icon">check_circle</text>
              <text>理想样本</text>
            </view>
            <text class="sample-title">面部居中、背景干净、五官清晰</text>
            <text class="sample-desc">AI 会对照 {{ currentStandard.country }} {{ currentStandard.visa }} 标准即时检查。</text>
          </view>
        </view>

        <view class="standard-grid">
          <view v-for="item in standardCards" :key="item.key" class="standard-card">
            <view class="standard-icon">
              <text class="material-symbols-outlined">{{ item.icon }}</text>
            </view>
            <text class="standard-title">{{ item.title }}</text>
            <text class="standard-desc">{{ item.value }}</text>
          </view>
        </view>

        <view v-if="uploadedImage" class="result-card">
          <view class="result-head">
            <view>
              <text class="result-title">AI 检测结果</text>
              <text class="result-sub">已按当前国家标准完成初步校验</text>
            </view>
            <view class="score-ring" :style="scoreStyle">
              <view class="score-inner">{{ photoScore }}%</view>
            </view>
          </view>
          <image class="uploaded-image" :src="uploadedImage" mode="aspectFill" />
          <view class="result-tags">
            <view v-for="item in resultTags" :key="item.label" class="result-tag" :class="`result-tag--${item.type}`">
              <text class="material-symbols-outlined">{{ item.icon }}</text>
              <text>{{ item.label }}</text>
            </view>
          </view>
        </view>

        <view class="check-section">
          <text class="section-label">合规清单</text>
          <view class="check-list">
            <view v-for="item in checklist" :key="item.key" class="check-row" @tap="toggleCheck(item.key)">
              <view class="check-left">
                <text class="material-symbols-outlined check-icon" :class="{ 'check-icon--ok': item.done }">{{ item.icon }}</text>
                <view>
                  <text class="check-title">{{ item.title }}</text>
                  <text class="check-desc">{{ item.desc }}</text>
                </view>
              </view>
              <text class="material-symbols-outlined check-state" :class="{ 'check-state--ok': item.done }">
                {{ item.done ? 'check_circle' : 'info' }}
              </text>
            </view>
          </view>
        </view>

        <view class="tips-card">
          <view class="tips-head">
            <text class="material-symbols-outlined">auto_fix_high</text>
            <text>拍摄优化建议</text>
          </view>
          <text class="tips-text">建议在白天靠近窗户拍摄，脸部不要有明显阴影；照片不加滤镜，头发不要遮挡眉毛和耳朵轮廓。</text>
        </view>
        <view class="quota-card">
          <view class="quota-head">
            <text class="material-symbols-outlined">shield</text>
            <text>检测额度</text>
          </view>
          <text class="quota-title">{{ quotaTitle }}</text>
          <text class="quota-desc">{{ quotaDesc }}</text>
          <text class="quota-note">{{ quotaNote }}</text>
        </view>
      </view>
    </scroll-view>

    <view class="photo-footer">
      <view class="footer-btn" @tap="choosePhoto">
        <text class="material-symbols-outlined">photo_camera</text>
        <text>{{ uploadedImage ? '重新上传照片' : '拍摄或上传照片' }}</text>
      </view>
      <text class="footer-note">Visago 仅在本地演示检测，不会上传你的照片。</text>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { api } from '../../../../utils/api'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'

const PHOTO_STANDARDS = [
  {
    country: '美国',
    flag: '🇺🇸',
    visa: '旅游签证 B1/B2',
    size: '51 x 51 毫米',
    background: '纯白或乳白',
    face: '头部占 50%-69%',
    file: 'JPG，≤ 240KB',
  },
  {
    country: '日本',
    flag: '🇯🇵',
    visa: '单次旅游签',
    size: '45 x 35 毫米',
    background: '纯白背景',
    face: '头顶到下巴 32-36mm',
    file: '清晰彩色照片',
  },
  {
    country: '英国',
    flag: '🇬🇧',
    visa: 'Standard Visitor',
    size: '45 x 35 毫米',
    background: '浅灰或乳白',
    face: '面部无遮挡',
    file: '600 x 750 像素以上',
  },
  {
    country: '澳大利亚',
    flag: '🇦🇺',
    visa: '访客签证',
    size: '45 x 35 毫米',
    background: '浅色纯背景',
    face: '自然表情，双眼睁开',
    file: '高分辨率彩色',
  },
]

export default {
  components: {
    VisagoTopBar,
  },
  data() {
    return {
      keyword: '',
      countryIndex: 0,
      uploadedImage: '',
      photoScore: 92,
      sampleImage: 'https://images.unsplash.com/photo-1607746882042-944635dfe10e?auto=format&fit=crop&w=900&q=80',
      checks: {
        face: true,
        glasses: false,
        size: false,
        resolution: true,
      },
      checkDetails: {
        face: '眼睛平视镜头，脸部完整露出',
        glasses: '避免反光、遮挡眉眼或改变面部轮廓',
        size: '',
        resolution: '背景均匀，面部没有明显阴影',
      },
      analyzing: false,
      quotaLoading: false,
      quotaError: '',
      quota: {
        dailyLimit: 5,
        usedCount: 0,
        remaining: 5,
        date: '',
        configured: false,
        provider: 'Aliyun',
      },
    }
  },
  computed: {
    countryLabels() {
      return this.visibleStandards.map((item) => `${item.flag} ${item.country} · ${item.visa}`)
    },
    visibleStandards() {
      const word = this.keyword.trim()
      const list = word
        ? PHOTO_STANDARDS.filter((item) => `${item.country}${item.visa}`.includes(word))
        : PHOTO_STANDARDS
      return list.length ? list : PHOTO_STANDARDS
    },
    currentStandard() {
      return PHOTO_STANDARDS[this.countryIndex] || PHOTO_STANDARDS[0]
    },
    standardCards() {
      return [
        { key: 'size', icon: 'straighten', title: '尺寸', value: this.currentStandard.size },
        { key: 'background', icon: 'wallpaper', title: '背景', value: this.currentStandard.background },
        { key: 'face', icon: 'face', title: '面部', value: this.currentStandard.face },
        { key: 'file', icon: 'account_box', title: '文件', value: this.currentStandard.file },
      ]
    },
    checklist() {
      return [
        { key: 'face', icon: 'face', title: '头部居中且比例合适', desc: '眼睛平视镜头，脸部完整露出', done: this.checks.face },
        { key: 'glasses', icon: 'eyeglasses', title: '不佩戴眼镜和夸张饰品', desc: '避免反光、遮挡眉眼或改变面部轮廓', done: this.checks.glasses },
        { key: 'size', icon: 'photo_size_select_large', title: '尺寸与文件大小符合要求', desc: `${this.currentStandard.size}，${this.currentStandard.file}`, done: this.checks.size },
        { key: 'resolution', icon: 'high_quality', title: '照片清晰，无滤镜和噪点', desc: '背景均匀，面部没有明显阴影', done: this.checks.resolution },
      ]
    },
    resultTags() {
      return [
        { label: '脸部清晰', icon: 'verified', type: 'ok' },
        { label: this.checks.glasses ? '无遮挡' : '需确认眼镜', icon: this.checks.glasses ? 'check_circle' : 'warning', type: this.checks.glasses ? 'ok' : 'warn' },
        { label: this.checks.size ? '尺寸合格' : '建议裁剪', icon: this.checks.size ? 'check_circle' : 'crop', type: this.checks.size ? 'ok' : 'todo' },
      ]
    },
    scoreStyle() {
      return {
        background: `conic-gradient(var(--visago-primary) ${this.photoScore}%, rgba(15, 101, 216, 0.16) ${this.photoScore}% 100%)`,
      }
    },
    quotaTitle() {
      if (this.quotaLoading) return '正在获取今日额度...'
      if (this.quotaError) return '暂时无法获取今日额度'
      return `当前账号今日已用 ${this.quota.usedCount}/${this.quota.dailyLimit} 次`
    },
    quotaDesc() {
      if (this.quotaLoading) return '请稍候，正在同步你的检测次数'
      if (this.quotaError) return this.quotaError
      return `剩余 ${this.quota.remaining} 次，按用户维度每天最多 ${this.quota.dailyLimit} 次`
    },
    quotaNote() {
      if (this.quotaLoading) return '额度信息会在每次进入页面时自动刷新'
      return this.quota.configured
        ? `${this.quota.provider} 服务已配置，后续接入正式检测后将按当前账号执行每日限额`
        : `${this.quota.provider} 服务尚未配置，当前页面仍为本地演示模式`
    },
  },
  onShow() {
    applyTheme(getStoredTheme())
    this.loadQuota()
  },
  methods: {
    async loadQuota() {
      this.quotaLoading = true
      this.quotaError = ''
      try {
        const quota = await api.getPhotoCheckQuota()
        this.quota = {
          dailyLimit: Number(quota.dailyLimit || 5),
          usedCount: Number(quota.usedCount || 0),
          remaining: Number(quota.remaining || 0),
          date: quota.date || '',
          configured: !!quota.configured,
          provider: quota.provider || 'Aliyun',
        }
      } catch (error) {
        this.quotaError = (error && error.message) || '额度信息获取失败'
      } finally {
        this.quotaLoading = false
      }
    },
    onCountryChange(e) {
      const selected = this.visibleStandards[Number(e.detail.value || 0)]
      const index = PHOTO_STANDARDS.findIndex((item) => item.country === selected.country && item.visa === selected.visa)
      this.countryIndex = Math.max(0, index)
    },
    toggleCheck(key) {
      this.checks[key] = !this.checks[key]
      this.updateScore()
    },
    updateScore() {
      const doneCount = Object.values(this.checks).filter(Boolean).length
      this.photoScore = Math.min(98, 64 + doneCount * 8 + (this.uploadedImage ? 2 : 0))
    },
    choosePhoto() {
      uni.chooseImage({
        count: 1,
        sizeType: ['compressed'],
        sourceType: ['album', 'camera'],
        success: (res) => {
          this.uploadedImage = res.tempFilePaths[0]
          this.checks = {
            face: true,
            glasses: true,
            size: false,
            resolution: true,
          }
          this.updateScore()
          uni.showToast({ title: '检测完成', icon: 'none' })
        },
        fail: () => {
          uni.showToast({ title: '未选择照片', icon: 'none' })
        },
      })
    },
  },
}
</script>

<style scoped>
.photo-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.photo-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: calc(92px + var(--visago-safe-bottom));
  left: 0;
  height: auto;
}

.photo-content {
  box-sizing: border-box;
  padding: 18px 16px 24px;
}

.search-panel {
  display: grid;
  gap: 12px;
}

.search-box {
  height: 44px;
  border-radius: 14px;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 14px;
  box-sizing: border-box;
}

.search-icon,
.picker-arrow {
  color: var(--visago-text-soft);
  font-size: 22px;
}

.search-input {
  flex: 1;
  height: 100%;
  color: var(--visago-text);
  font-size: 14px;
}

.country-card,
.standard-card,
.result-card,
.check-list,
.tips-card {
  border: 1px solid var(--visago-line);
  background: color-mix(in srgb, var(--visago-surface) 88%, transparent);
  box-shadow: var(--visago-shadow-card);
}

.country-card {
  min-height: 66px;
  border-radius: 16px;
  padding: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
}

.country-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.flag-wrap {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.country-name,
.country-visa,
.standard-title,
.standard-desc,
.result-title,
.result-sub,
.check-title,
.check-desc,
.tips-text {
  display: block;
}

.country-name {
  font-size: 15px;
  font-weight: 800;
}

.country-visa {
  margin-top: 4px;
  font-size: 12px;
  color: var(--visago-text-muted);
}

.sample-card {
  position: relative;
  margin-top: 18px;
  height: 214px;
  border-radius: 20px;
  overflow: hidden;
  border: 1px solid var(--visago-line);
  box-shadow: 0 16px 32px rgba(15, 23, 42, 0.14);
}

.sample-image {
  width: 100%;
  height: 100%;
}

.sample-mask {
  position: absolute;
  inset: auto 0 0;
  padding: 54px 16px 16px;
  background: linear-gradient(180deg, transparent 0%, rgba(0, 0, 0, 0.82) 72%);
  color: #fff;
}

.sample-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  font-weight: 800;
  letter-spacing: 0.08em;
}

.sample-badge-icon {
  color: #53e16f;
  font-size: 17px;
  font-variation-settings: 'FILL' 1;
}

.sample-title {
  display: block;
  margin-top: 8px;
  font-size: 17px;
  font-weight: 800;
}

.sample-desc {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.78);
}

.standard-grid {
  margin-top: 18px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.standard-card {
  min-height: 124px;
  border-radius: 16px;
  padding: 14px;
  box-sizing: border-box;
}

.standard-icon {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: rgba(15, 101, 216, 0.14);
  color: var(--visago-primary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.standard-title {
  margin-top: 12px;
  font-size: 14px;
  font-weight: 800;
}

.standard-desc {
  margin-top: 5px;
  font-size: 12px;
  line-height: 1.45;
  color: var(--visago-text-muted);
}

.result-card {
  margin-top: 18px;
  border-radius: 18px;
  padding: 14px;
}

.result-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.result-title {
  font-size: 17px;
  font-weight: 800;
}

.result-sub {
  margin-top: 4px;
  font-size: 12px;
  color: var(--visago-text-muted);
}

.score-ring {
  width: 58px;
  height: 58px;
  border-radius: 50%;
  padding: 5px;
  box-sizing: border-box;
  flex-shrink: 0;
}

.score-inner {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: var(--visago-surface);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 900;
  color: var(--visago-primary);
}

.uploaded-image {
  width: 100%;
  height: 170px;
  margin-top: 14px;
  border-radius: 14px;
}

.result-tags {
  margin-top: 12px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.result-tag {
  height: 30px;
  border-radius: 999px;
  padding: 0 10px;
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 700;
}

.result-tag .material-symbols-outlined {
  font-size: 16px;
}

.result-tag--ok {
  background: rgba(22, 163, 74, 0.14);
  color: #16a34a;
}

.result-tag--warn {
  background: rgba(245, 158, 11, 0.14);
  color: #d97706;
}

.result-tag--todo {
  background: rgba(15, 101, 216, 0.14);
  color: var(--visago-primary);
}

.check-section {
  margin-top: 20px;
}

.section-label {
  display: block;
  margin: 0 0 9px 4px;
  font-size: 12px;
  font-weight: 800;
  color: var(--visago-text-soft);
  letter-spacing: 0.12em;
}

.check-list {
  border-radius: 16px;
  overflow: hidden;
}

.check-row {
  min-height: 66px;
  padding: 12px 14px;
  border-bottom: 1px solid var(--visago-line);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  box-sizing: border-box;
}

.check-row:last-child {
  border-bottom: none;
}

.check-left {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.check-icon {
  color: var(--visago-text-soft);
  font-size: 22px;
  flex-shrink: 0;
}

.check-icon--ok,
.check-state--ok {
  color: var(--visago-primary);
}

.check-title {
  font-size: 14px;
  font-weight: 700;
}

.check-desc {
  margin-top: 3px;
  font-size: 11px;
  line-height: 1.35;
  color: var(--visago-text-muted);
}

.check-state {
  color: var(--visago-text-soft);
  font-size: 21px;
  flex-shrink: 0;
}

.tips-card {
  margin-top: 18px;
  border-radius: 16px;
  padding: 14px;
}

.tips-head {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--visago-primary);
  font-size: 15px;
  font-weight: 800;
}

.tips-text {
  margin-top: 8px;
  font-size: 12px;
  line-height: 1.55;
  color: var(--visago-text-muted);
}

.quota-card {
  margin-top: 16px;
  border: 1px solid var(--visago-line);
  border-radius: 20px;
  padding: 16px;
  background: color-mix(in srgb, var(--visago-surface) 88%, transparent);
  box-shadow: var(--visago-shadow-card);
  box-sizing: border-box;
}

.quota-head {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--visago-primary);
  font-size: 15px;
  font-weight: 800;
}

.quota-title {
  display: block;
  margin-top: 12px;
  font-size: 16px;
  font-weight: 800;
  color: var(--visago-text);
}

.quota-desc {
  display: block;
  margin-top: 8px;
  font-size: 13px;
  line-height: 1.6;
  color: var(--visago-text);
}

.quota-note {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.6;
  color: var(--visago-text-soft);
}

.photo-footer {
  position: fixed;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 50;
  padding: 12px 16px calc(10px + var(--visago-safe-bottom));
  background: color-mix(in srgb, var(--visago-surface) 94%, transparent);
  border-top: 1px solid var(--visago-line);
  backdrop-filter: blur(16px);
  box-sizing: border-box;
}

.footer-btn {
  height: 52px;
  max-width: 398px;
  margin: 0 auto;
  border-radius: 16px;
  color: #fff;
  background: linear-gradient(135deg, var(--visago-primary) 0%, #2f89ff 100%);
  box-shadow: 0 10px 22px rgba(15, 101, 216, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 800;
}

.footer-note {
  display: block;
  margin-top: 7px;
  text-align: center;
  font-size: 11px;
  color: var(--visago-text-soft);
}

:global(html.theme-dark) .country-card,
:global(html.theme-dark) .standard-card,
:global(html.theme-dark) .result-card,
:global(html.theme-dark) .check-list,
:global(html.theme-dark) .tips-card,
:global(html.theme-dark) .quota-card {
  background: rgba(23, 27, 36, 0.86);
}

:global(html.theme-dark) .sample-card {
  box-shadow: 0 18px 36px rgba(0, 0, 0, 0.32);
}
</style>
