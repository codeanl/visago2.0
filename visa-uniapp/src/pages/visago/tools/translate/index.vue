<template>
  <view class="translate-page" :class="{ 'translate-page--dark': themeMode === 'dark' }">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="翻译助手" />

    <view class="translate-content visago-page-width">
      <view class="lang-switch visago-card">
        <view class="lang-btn" @tap="toggleFromLang">
          <text class="lang-text">{{ fromLang.label }}</text>
          <text class="material-symbols-outlined lang-drop">expand_more</text>
        </view>
        <view class="swap-btn" @tap="swapLangs">
          <text class="material-symbols-outlined">swap_horiz</text>
        </view>
        <view class="lang-btn" @tap="toggleToLang">
          <text class="lang-text">{{ toLang.label }}</text>
          <text class="material-symbols-outlined lang-drop">expand_more</text>
        </view>
      </view>

      <view class="input-wrap visago-card">
        <textarea v-model.trim="inputText" class="translate-input" placeholder="输入要翻译的文字..." />
        <view class="input-actions">
          <view class="round-action" @tap="clearInput">
            <text class="material-symbols-outlined">close</text>
          </view>
          <view class="round-action round-action--primary" :class="{ 'round-action--disabled': !canTranslate }" @tap="translateNow">
            <text class="material-symbols-outlined">{{ translating ? 'hourglass_top' : 'translate' }}</text>
          </view>
        </view>
      </view>

      <view class="mode-card visago-card">
        <view class="mode-copy">
          <text class="mode-title">签证材料翻译</text>
          <text class="mode-desc">适合翻译申请说明、行程描述、在职证明片段、地址和常用材料字段。</text>
        </view>
        <view class="mode-pill">MVP</view>
      </view>

      <view class="result-wrap visago-card">
        <text class="result-text">{{ resultText }}</text>
        <view class="result-actions">
          <view class="icon-action" @tap="copyResult">
            <text class="material-symbols-outlined">content_copy</text>
          </view>
        </view>
      </view>

      <view class="action-grid">
        <view class="action-card visago-card" @tap="comingSoon('拍照翻译')">
          <view class="action-icon action-icon--primary">
            <text class="material-symbols-outlined">photo_camera</text>
          </view>
          <view>
            <text class="action-title">拍照翻译</text>
            <text class="action-sub">下一阶段接 OCR + 翻译</text>
          </view>
        </view>
        <view class="action-card visago-card" @tap="comingSoon('文档翻译')">
          <view class="action-icon action-icon--tertiary">
            <text class="material-symbols-outlined">description</text>
          </view>
          <view>
            <text class="action-title">文档翻译</text>
            <text class="action-sub">后续支持 PDF / DOC</text>
          </view>
        </view>
      </view>

      <view class="history-head">
        <text class="history-title">最近记录</text>
        <text class="history-more" @tap="clearHistory">清空</text>
      </view>

      <view class="history-list">
        <view v-for="item in histories" :key="item.id" class="history-item visago-card" @tap="applyHistory(item)">
          <view class="history-main">
            <text class="history-origin">{{ item.origin }}</text>
            <text class="history-result">{{ item.result }}</text>
            <text class="history-time">{{ item.time }}</text>
          </view>
          <text class="material-symbols-outlined history-star">history</text>
        </view>

        <view v-if="!histories.length" class="history-empty visago-card">
          <text>还没有翻译记录</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { api } from '../../../../utils/api'
import { applyTheme, getStoredTheme, THEME_CHANGE_EVENT } from '../../../../utils/theme'

const LANGS = [
  { key: 'zh', label: '中文' },
  { key: 'en', label: 'English' },
]

const HISTORY_KEY = 'visago_translate_history_v1'

function nowLabel() {
  const date = new Date()
  const month = `${date.getMonth() + 1}`.padStart(2, '0')
  const day = `${date.getDate()}`.padStart(2, '0')
  const hours = `${date.getHours()}`.padStart(2, '0')
  const minutes = `${date.getMinutes()}`.padStart(2, '0')
  return `${month}-${day} ${hours}:${minutes}`
}

export default {
  components: { VisagoTopBar },
  data() {
    return {
      themeMode: 'light',
      fromLangIndex: 0,
      toLangIndex: 1,
      inputText: '',
      translatedText: '',
      translating: false,
      histories: [],
    }
  },
  computed: {
    fromLang() {
      return LANGS[this.fromLangIndex]
    },
    toLang() {
      return LANGS[this.toLangIndex]
    },
    canTranslate() {
      return !this.translating && !!this.inputText.trim()
    },
    resultText() {
      if (this.translating) {
        return '翻译中...'
      }
      if (this.translatedText) {
        return this.translatedText
      }
      return '翻译结果会显示在这里...'
    },
  },
  onLoad() {
    if (typeof uni !== 'undefined' && uni.$on) {
      uni.$on(THEME_CHANGE_EVENT, this.onThemeChange)
    }
    this.syncTheme()
    this.restoreHistory()
  },
  onShow() {
    this.syncTheme()
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
    toggleFromLang() {
      this.fromLangIndex = this.fromLangIndex === 0 ? 1 : 0
      if (this.fromLangIndex === this.toLangIndex) {
        this.toLangIndex = this.toLangIndex === 0 ? 1 : 0
      }
      this.translatedText = ''
    },
    toggleToLang() {
      this.toLangIndex = this.toLangIndex === 0 ? 1 : 0
      if (this.fromLangIndex === this.toLangIndex) {
        this.fromLangIndex = this.fromLangIndex === 0 ? 1 : 0
      }
      this.translatedText = ''
    },
    swapLangs() {
      const temp = this.fromLangIndex
      this.fromLangIndex = this.toLangIndex
      this.toLangIndex = temp
      this.translatedText = ''
    },
    clearInput() {
      this.inputText = ''
      this.translatedText = ''
    },
    async translateNow() {
      if (!this.canTranslate) return
      this.translating = true
      try {
        const result = await api.translateText({
          text: this.inputText,
          sourceLang: this.fromLang.key,
          targetLang: this.toLang.key,
          mode: 'visa',
        })
        this.translatedText = result.translatedText || ''
        this.saveHistory({
          id: `history-${Date.now()}`,
          origin: this.inputText,
          result: this.translatedText,
          sourceLang: this.fromLang.key,
          targetLang: this.toLang.key,
          time: nowLabel(),
        })
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '翻译失败',
          icon: 'none',
        })
      } finally {
        this.translating = false
      }
    },
    copyResult() {
      const text = this.translatedText.trim()
      if (!text) {
        uni.showToast({
          title: '暂无翻译结果',
          icon: 'none',
        })
        return
      }
      uni.setClipboardData({
        data: text,
        success: () => {
          uni.showToast({
            title: '已复制',
            icon: 'none',
          })
        },
      })
    },
    saveHistory(entry) {
      const list = [entry, ...this.histories.filter((item) => item.origin !== entry.origin || item.result !== entry.result)].slice(0, 10)
      this.histories = list
      uni.setStorageSync(HISTORY_KEY, list)
    },
    restoreHistory() {
      try {
        const stored = uni.getStorageSync(HISTORY_KEY)
        this.histories = Array.isArray(stored) ? stored : []
      } catch (error) {
        this.histories = []
      }
    },
    clearHistory() {
      this.histories = []
      uni.removeStorageSync(HISTORY_KEY)
      uni.showToast({
        title: '已清空',
        icon: 'none',
      })
    },
    applyHistory(item) {
      const sourceIndex = LANGS.findIndex((lang) => lang.key === item.sourceLang)
      const targetIndex = LANGS.findIndex((lang) => lang.key === item.targetLang)
      if (sourceIndex >= 0) this.fromLangIndex = sourceIndex
      if (targetIndex >= 0) this.toLangIndex = targetIndex
      this.inputText = item.origin || ''
      this.translatedText = item.result || ''
    },
    comingSoon(title) {
      uni.showToast({
        title: `${title}开发中`,
        icon: 'none',
      })
    },
  },
}
</script>

<style scoped>
.translate-page {
  min-height: 100vh;
  background: var(--visago-bg);
}

.translate-content {
  box-sizing: border-box;
  padding: 92px 16px 22px;
}

.lang-switch {
  border-radius: 16px;
  background: var(--visago-surface-soft);
  border-color: transparent;
  box-shadow: none;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.lang-btn {
  display: flex;
  align-items: center;
  gap: 2px;
  color: var(--visago-primary);
}

.lang-text {
  font-size: 22px;
  font-weight: 600;
}

.lang-drop {
  font-size: 18px;
}

.swap-btn {
  width: 40px;
  height: 40px;
  border-radius: 9999px;
  background: var(--visago-surface);
  color: var(--visago-primary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.input-wrap {
  margin-top: 14px;
  padding: 12px;
  border-radius: 16px;
  background: var(--visago-surface-soft);
  border-color: transparent;
  box-shadow: none;
}

.translate-input {
  width: 100%;
  height: 160px;
  color: var(--visago-text);
  font-size: 16px;
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.round-action {
  width: 40px;
  height: 40px;
  border-radius: 9999px;
  background: rgba(10, 132, 255, 0.14);
  color: var(--visago-primary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.round-action--primary {
  background: var(--visago-primary);
  color: #fff;
}

.round-action--disabled {
  opacity: 0.45;
}

.mode-card {
  margin-top: 12px;
  border-radius: 16px;
  padding: 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.mode-copy {
  min-width: 0;
}

.mode-title,
.mode-desc {
  display: block;
}

.mode-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--visago-text);
}

.mode-desc {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.5;
  color: var(--visago-text-muted);
}

.mode-pill {
  flex-shrink: 0;
  padding: 5px 10px;
  border-radius: 999px;
  background: rgba(10, 132, 255, 0.12);
  color: var(--visago-primary);
  font-size: 11px;
  font-weight: 700;
}

.result-wrap {
  margin-top: 12px;
  min-height: 122px;
  border-radius: 16px;
  padding: 12px;
  position: relative;
}

.result-text {
  color: var(--visago-text-muted);
  font-size: 16px;
  line-height: 1.45;
  white-space: pre-wrap;
}

.result-actions {
  position: absolute;
  top: 12px;
  right: 12px;
  display: flex;
  gap: 8px;
}

.icon-action {
  color: var(--visago-text-soft);
}

.action-grid {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.action-card {
  border-radius: 16px;
  padding: 12px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.action-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-icon--primary {
  background: rgba(10, 132, 255, 0.2);
  color: var(--visago-primary);
}

.action-icon--tertiary {
  background: rgba(0, 151, 83, 0.2);
  color: #0f9735;
}

.action-title {
  display: block;
  font-size: 14px;
  font-weight: 600;
}

.action-sub {
  display: block;
  margin-top: 2px;
  font-size: 11px;
  color: var(--visago-text-soft);
}

.history-head {
  margin-top: 18px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.history-title {
  font-size: 22px;
  font-weight: 600;
}

.history-more {
  color: var(--visago-primary);
  font-size: 13px;
}

.history-list {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.history-item,
.history-empty {
  border-radius: 14px;
  padding: 12px;
}

.history-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
}

.history-empty {
  text-align: center;
  color: var(--visago-text-soft);
}

.history-main {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.history-origin {
  color: var(--visago-text);
  font-size: 14px;
  line-height: 1.4;
}

.history-result {
  color: var(--visago-primary);
  font-size: 14px;
  line-height: 1.4;
}

.history-time {
  color: var(--visago-text-soft);
  font-size: 11px;
}

.history-star {
  color: var(--visago-text-soft);
  font-size: 20px;
}

.translate-page--dark,
:global(html.theme-dark) .translate-page {
  background: #000;
}

.translate-page--dark .lang-switch,
:global(html.theme-dark) .lang-switch {
  background: #1c1c1e;
}

.translate-page--dark .swap-btn,
:global(html.theme-dark) .swap-btn {
  background: #2c2c2e;
}
</style>
