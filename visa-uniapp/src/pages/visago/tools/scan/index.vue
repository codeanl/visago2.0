<template>
  <view class="scan-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="文件扫描" />

    <view class="mode-bar visago-page-width">
      <view class="mode-pill">
        <view class="mode-item" :class="{ 'mode-item--active': mode === 'library' }" @tap="mode = 'library'">文件库</view>
        <view class="mode-item" :class="{ 'mode-item--active': mode === 'camera' }" @tap="mode = 'camera'">拍摄扫描</view>
      </view>
    </view>

    <scroll-view v-if="mode === 'library'" scroll-y class="library-scroll">
      <view class="scan-content visago-page-width">
        <view class="title-block">
          <text class="page-title">文件扫描件</text>
          <text class="page-sub">护照、签证、保险、流水统一收纳，申请时一键复用。</text>
        </view>

        <view class="search-box">
          <text class="material-symbols-outlined search-icon">search</text>
          <input v-model.trim="keyword" class="search-input" placeholder="搜索文件名称..." />
        </view>

        <scroll-view scroll-x class="chip-scroll">
          <view class="chip-row">
            <view v-for="item in categories" :key="item.key" class="chip" :class="{ 'chip--active': activeCategory === item.key }" @tap="activeCategory = item.key">
              {{ item.label }}
            </view>
          </view>
        </scroll-view>

        <view v-if="featuredDoc" class="featured-card" @tap="previewDoc(featuredDoc)">
          <image class="featured-image" :src="featuredDoc.image" mode="aspectFill" />
          <view class="featured-main">
            <view class="tag-row">
              <text class="doc-tag">{{ featuredDoc.typeLabel }}</text>
              <text class="material-symbols-outlined more-icon">more_horiz</text>
            </view>
            <text class="featured-title">{{ featuredDoc.title }}</text>
            <text class="featured-date">{{ featuredDoc.date }}</text>
            <view class="action-row">
              <view class="mini-action">
                <text class="material-symbols-outlined">edit</text>
                <text>编辑</text>
              </view>
              <view class="mini-action mini-action--ghost">
                <text class="material-symbols-outlined">ios_share</text>
                <text>导出</text>
              </view>
            </view>
          </view>
        </view>

        <view class="doc-grid">
          <view v-for="doc in secondaryDocs" :key="doc.id" class="doc-card" @tap="previewDoc(doc)">
            <view class="doc-cover">
              <image class="doc-image" :src="doc.image" mode="aspectFill" />
              <text v-if="doc.status" class="status-badge" :class="`status-badge--${doc.statusType}`">{{ doc.status }}</text>
            </view>
            <text class="doc-title">{{ doc.title }}</text>
            <text class="doc-date">{{ doc.date }}</text>
          </view>

          <view class="add-card" @tap="startScan">
            <text class="material-symbols-outlined">add_circle</text>
            <text>扫描新文档</text>
          </view>
        </view>

        <view v-if="!filteredDocs.length" class="empty-card">
          <text class="material-symbols-outlined">folder_off</text>
          <text>没有找到相关文件</text>
        </view>
      </view>
    </scroll-view>

    <view v-else class="camera-panel">
      <image class="camera-bg" src="https://images.unsplash.com/photo-1586953208448-b95a79798f07?auto=format&fit=crop&w=1100&q=80" mode="aspectFill" />
      <view class="camera-shade" />

      <view class="camera-controls">
        <view class="circle-btn">
          <text class="material-symbols-outlined">flash_on</text>
        </view>
        <view class="focus-chip">正在对焦...</view>
        <view class="circle-btn" @tap="mode = 'library'">
          <text class="material-symbols-outlined">close</text>
        </view>
      </view>

      <view class="scan-frame">
        <view class="corner corner--tl" />
        <view class="corner corner--tr" />
        <view class="corner corner--bl" />
        <view class="corner corner--br" />
        <view class="scan-line" />
      </view>

      <view class="capture-area">
        <text class="capture-tip">将证件置于框内，我们将自动捕捉并优化文档边缘</text>
        <view class="capture-row">
          <view class="gallery-thumb" @tap="mode = 'library'">
            <image src="https://images.unsplash.com/photo-1554224155-6726b3ff858f?auto=format&fit=crop&w=300&q=80" mode="aspectFill" />
          </view>
          <view class="shutter" @tap="captureMock">
            <view class="shutter-inner" />
          </view>
          <view class="auto-mode">
            <view class="auto-icon">
              <text class="material-symbols-outlined">auto_awesome</text>
            </view>
            <text>自动模式</text>
          </view>
        </view>
      </view>
    </view>

    <view v-if="mode === 'library'" class="scan-fab" @tap="startScan">
      <text class="material-symbols-outlined">add</text>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'

const BASE_DOCS = [
  {
    id: 'passport',
    title: '个人护照 - 陈先生',
    type: 'passport',
    typeLabel: 'Passport',
    date: '更新于 2026年04月12日',
    status: '已验证',
    statusType: 'ok',
    image: 'https://images.unsplash.com/photo-1624138784614-87fd1b6528f8?auto=format&fit=crop&w=800&q=80',
  },
  {
    id: 'visa-fr',
    title: '申根签证 - 法国',
    type: 'visa',
    typeLabel: 'Visa',
    date: '2026年03月05日',
    status: '已验证',
    statusType: 'ok',
    image: 'https://images.unsplash.com/photo-1554224154-26032ffc0d07?auto=format&fit=crop&w=800&q=80',
  },
  {
    id: 'insurance',
    title: '安联全球旅行险',
    type: 'insurance',
    typeLabel: 'Insurance',
    date: '2026年02月20日',
    status: '即将过期',
    statusType: 'warn',
    image: 'https://images.unsplash.com/photo-1450101499163-c8848c66ca85?auto=format&fit=crop&w=800&q=80',
  },
  {
    id: 'bank',
    title: '银行存款证明',
    type: 'bank',
    typeLabel: 'Bank',
    date: '2026年04月01日',
    status: '',
    statusType: '',
    image: 'https://images.unsplash.com/photo-1554224155-6726b3ff858f?auto=format&fit=crop&w=800&q=80',
  },
]

export default {
  components: {
    VisagoTopBar,
  },
  data() {
    return {
      mode: 'library',
      keyword: '',
      activeCategory: 'all',
      docs: BASE_DOCS,
      categories: [
        { key: 'all', label: '全部' },
        { key: 'passport', label: '护照' },
        { key: 'visa', label: '签证' },
        { key: 'insurance', label: '保险' },
        { key: 'bank', label: '资金证明' },
      ],
    }
  },
  computed: {
    filteredDocs() {
      const word = this.keyword.trim().toLowerCase()
      return this.docs.filter((doc) => {
        const matchType = this.activeCategory === 'all' || doc.type === this.activeCategory
        const matchWord = !word || `${doc.title}${doc.typeLabel}`.toLowerCase().includes(word)
        return matchType && matchWord
      })
    },
    featuredDoc() {
      return this.filteredDocs[0]
    },
    secondaryDocs() {
      return this.filteredDocs.slice(1)
    },
  },
  onShow() {
    applyTheme(getStoredTheme())
  },
  methods: {
    startScan() {
      this.mode = 'camera'
    },
    captureMock() {
      const nextDoc = {
        id: `scan-${Date.now()}`,
        title: '新扫描文件',
        type: 'passport',
        typeLabel: 'Scan',
        date: '刚刚更新',
        status: '待命名',
        statusType: 'todo',
        image: 'https://images.unsplash.com/photo-1586953208448-b95a79798f07?auto=format&fit=crop&w=800&q=80',
      }
      this.docs = [nextDoc, ...this.docs]
      this.mode = 'library'
      this.activeCategory = 'all'
      uni.showToast({ title: '已生成扫描件', icon: 'none' })
    },
    previewDoc(doc) {
      uni.showToast({ title: `${doc.title} 已打开`, icon: 'none' })
    },
  },
}
</script>

<style scoped>
.scan-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.mode-bar {
  position: fixed;
  top: 74px;
  left: 0;
  right: 0;
  z-index: 30;
  padding: 12px 16px;
  box-sizing: border-box;
  background: var(--visago-bg);
}

.mode-pill {
  height: 42px;
  padding: 4px;
  border-radius: 14px;
  background: var(--visago-surface-soft);
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 4px;
}

.mode-item {
  border-radius: 11px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--visago-text-muted);
  font-size: 14px;
  font-weight: 800;
}

.mode-item--active {
  color: var(--visago-primary);
  background: var(--visago-surface);
  box-shadow: var(--visago-shadow-card);
}

.library-scroll {
  position: fixed;
  top: 140px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.scan-content {
  padding: 8px 16px calc(92px + var(--visago-safe-bottom));
  box-sizing: border-box;
}

.title-block {
  margin-bottom: 18px;
}

.page-title,
.page-sub,
.featured-title,
.featured-date,
.doc-title,
.doc-date {
  display: block;
}

.page-title {
  font-size: 27px;
  line-height: 1.18;
  font-weight: 900;
}

.page-sub {
  margin-top: 6px;
  font-size: 13px;
  line-height: 1.45;
  color: var(--visago-text-muted);
}

.search-box {
  height: 44px;
  border-radius: 14px;
  background: var(--visago-surface-soft);
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
  margin: 14px -16px 0;
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
  font-weight: 700;
}

.chip--active {
  color: #fff;
  background: var(--visago-primary);
  border-color: var(--visago-primary);
  box-shadow: 0 8px 18px rgba(15, 101, 216, 0.18);
}

.featured-card {
  margin-top: 16px;
  border-radius: 18px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
  padding: 14px;
  display: flex;
  gap: 14px;
  overflow: hidden;
}

.featured-image {
  width: 84px;
  height: 118px;
  border-radius: 12px;
  background: var(--visago-surface-soft);
  flex-shrink: 0;
}

.featured-main {
  flex: 1;
  min-width: 0;
}

.tag-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.doc-tag {
  padding: 4px 8px;
  border-radius: 999px;
  background: rgba(15, 101, 216, 0.12);
  color: var(--visago-primary);
  font-size: 10px;
  font-weight: 900;
  text-transform: uppercase;
}

.more-icon {
  color: var(--visago-text-soft);
}

.featured-title {
  margin-top: 12px;
  font-size: 16px;
  font-weight: 900;
}

.featured-date {
  margin-top: 5px;
  font-size: 12px;
  color: var(--visago-text-muted);
}

.action-row {
  display: flex;
  gap: 8px;
  margin-top: 16px;
}

.mini-action {
  height: 31px;
  padding: 0 10px;
  border-radius: 9px;
  background: rgba(15, 101, 216, 0.13);
  color: var(--visago-primary);
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 800;
}

.mini-action--ghost {
  color: var(--visago-text-muted);
  background: var(--visago-surface-soft);
}

.mini-action .material-symbols-outlined {
  font-size: 16px;
}

.doc-grid {
  margin-top: 16px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 13px;
}

.doc-card,
.add-card,
.empty-card {
  border-radius: 18px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
}

.doc-card {
  padding: 12px;
}

.doc-cover {
  position: relative;
  width: 100%;
  aspect-ratio: 3 / 4;
  border-radius: 13px;
  overflow: hidden;
  background: var(--visago-surface-soft);
}

.doc-image {
  width: 100%;
  height: 100%;
}

.status-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 3px 7px;
  border-radius: 999px;
  font-size: 10px;
  font-weight: 800;
  color: #fff;
}

.status-badge--ok {
  background: #16a34a;
}

.status-badge--warn {
  background: #f59e0b;
}

.status-badge--todo {
  background: var(--visago-primary);
}

.doc-title {
  margin-top: 10px;
  font-size: 14px;
  font-weight: 800;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.doc-date {
  margin-top: 4px;
  font-size: 11px;
  color: var(--visago-text-muted);
}

.add-card {
  min-height: 178px;
  border-style: dashed;
  color: var(--visago-text-muted);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 800;
}

.add-card .material-symbols-outlined {
  color: var(--visago-primary);
  font-size: 34px;
}

.empty-card {
  margin-top: 16px;
  padding: 28px;
  color: var(--visago-text-muted);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.camera-panel {
  position: fixed;
  top: 140px;
  right: 0;
  bottom: 0;
  left: 0;
  overflow: hidden;
  background: #000;
}

.camera-bg,
.camera-shade {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

.camera-bg {
  opacity: 0.78;
}

.camera-shade {
  background: radial-gradient(circle at center, transparent 0%, rgba(0, 0, 0, 0.14) 40%, rgba(0, 0, 0, 0.7) 100%);
}

.camera-controls {
  position: absolute;
  top: 20px;
  left: 20px;
  right: 20px;
  z-index: 3;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.circle-btn,
.focus-chip,
.auto-icon {
  background: rgba(0, 0, 0, 0.42);
  color: #fff;
  backdrop-filter: blur(14px);
  border: 1px solid rgba(255, 255, 255, 0.12);
}

.circle-btn {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.focus-chip {
  height: 32px;
  padding: 0 15px;
  border-radius: 999px;
  display: flex;
  align-items: center;
  font-size: 13px;
  font-weight: 700;
}

.scan-frame {
  position: absolute;
  z-index: 2;
  top: 16%;
  left: 50%;
  width: 78%;
  max-width: 330px;
  aspect-ratio: 1 / 1.414;
  transform: translateX(-50%);
  border: 2px solid rgba(0, 122, 255, 0.28);
  border-radius: 20px;
  overflow: hidden;
}

.corner {
  position: absolute;
  width: 42px;
  height: 42px;
  border-color: #0a84ff;
}

.corner--tl {
  top: 0;
  left: 0;
  border-top: 5px solid;
  border-left: 5px solid;
  border-top-left-radius: 18px;
}

.corner--tr {
  top: 0;
  right: 0;
  border-top: 5px solid;
  border-right: 5px solid;
  border-top-right-radius: 18px;
}

.corner--bl {
  bottom: 0;
  left: 0;
  border-bottom: 5px solid;
  border-left: 5px solid;
  border-bottom-left-radius: 18px;
}

.corner--br {
  right: 0;
  bottom: 0;
  border-right: 5px solid;
  border-bottom: 5px solid;
  border-bottom-right-radius: 18px;
}

.scan-line {
  position: absolute;
  left: 0;
  right: 0;
  top: 28%;
  height: 72px;
  background: linear-gradient(180deg, transparent, rgba(0, 122, 255, 0.18), transparent);
}

.capture-area {
  position: absolute;
  z-index: 3;
  left: 0;
  right: 0;
  bottom: calc(26px + var(--visago-safe-bottom));
  display: flex;
  flex-direction: column;
  align-items: center;
}

.capture-tip {
  width: 280px;
  text-align: center;
  color: rgba(255, 255, 255, 0.75);
  font-size: 13px;
  line-height: 1.45;
}

.capture-row {
  width: 100%;
  margin-top: 26px;
  padding: 0 34px;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.gallery-thumb {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  overflow: hidden;
  border: 2px solid rgba(255, 255, 255, 0.26);
}

.gallery-thumb image {
  width: 100%;
  height: 100%;
}

.shutter {
  width: 78px;
  height: 78px;
  border-radius: 50%;
  border: 4px solid rgba(255, 255, 255, 0.34);
  display: flex;
  align-items: center;
  justify-content: center;
}

.shutter-inner {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: #fff;
  border: 2px solid #111;
}

.auto-mode {
  width: 62px;
  color: #fff;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 800;
}

.auto-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.scan-fab {
  position: fixed;
  right: 24px;
  bottom: calc(28px + var(--visago-safe-bottom));
  z-index: 40;
  width: 58px;
  height: 58px;
  border-radius: 50%;
  background: var(--visago-primary);
  color: #fff;
  box-shadow: var(--visago-shadow-fab);
  display: flex;
  align-items: center;
  justify-content: center;
}

:global(html.theme-dark) .mode-item--active,
:global(html.theme-dark) .featured-card,
:global(html.theme-dark) .doc-card,
:global(html.theme-dark) .add-card,
:global(html.theme-dark) .empty-card {
  background: #171b24;
}
</style>
