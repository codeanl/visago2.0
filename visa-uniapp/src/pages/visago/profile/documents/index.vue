<template>
  <view class="doc-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="我的资料" />
    <scroll-view scroll-y class="doc-scroll">
      <view class="doc-content visago-page-width">
        <view class="vault-card">
          <view>
            <text class="vault-title">资料保险箱</text>
            <text class="vault-sub">护照、照片、流水和申请材料集中管理</text>
          </view>
          <text class="material-symbols-outlined vault-icon">folder_shared</text>
        </view>

        <view class="doc-grid">
          <view v-for="item in docs" :key="item.title" class="doc-card">
            <view class="doc-icon" :style="{ color: item.color, background: item.tint }">
              <text class="material-symbols-outlined">{{ item.icon }}</text>
            </view>
            <text class="doc-title">{{ item.title }}</text>
            <text class="doc-meta">{{ item.meta }}</text>
            <view class="doc-status">{{ item.status }}</view>
          </view>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'

export default {
  components: { VisagoTopBar },
  data() {
    return {
      docs: [
        { title: '护照扫描件', meta: '2 个文件', status: '已验证', icon: 'badge', color: '#0f65d8', tint: 'rgba(15,101,216,0.12)' },
        { title: '证件照片', meta: '日本/美国规格', status: '合规', icon: 'photo_camera', color: '#16a34a', tint: 'rgba(22,163,74,0.12)' },
        { title: '银行流水', meta: '近 6 个月', status: '待更新', icon: 'account_balance', color: '#f59e0b', tint: 'rgba(245,158,11,0.13)' },
        { title: '在职证明', meta: '英文模板', status: '草稿', icon: 'work', color: '#8b5cf6', tint: 'rgba(139,92,246,0.12)' },
      ],
    }
  },
  onShow() {
    applyTheme(getStoredTheme())
  },
}
</script>

<style scoped>
.doc-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.doc-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.doc-content {
  padding: 18px 16px 34px;
  box-sizing: border-box;
}

.vault-card {
  border-radius: 22px;
  padding: 18px;
  background: linear-gradient(135deg, #111827, var(--visago-primary));
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 16px 34px rgba(15, 101, 216, 0.22);
}

.vault-title,
.vault-sub,
.doc-title,
.doc-meta {
  display: block;
}

.vault-title {
  font-size: 22px;
  font-weight: 900;
}

.vault-sub {
  margin-top: 5px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.78);
}

.vault-icon {
  font-size: 44px;
  opacity: 0.8;
}

.doc-grid {
  margin-top: 16px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.doc-card {
  min-height: 146px;
  padding: 14px;
  border-radius: 18px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
  box-sizing: border-box;
}

.doc-icon {
  width: 42px;
  height: 42px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.doc-title {
  margin-top: 14px;
  font-size: 15px;
  font-weight: 900;
}

.doc-meta {
  margin-top: 5px;
  font-size: 12px;
  color: var(--visago-text-muted);
}

.doc-status {
  width: fit-content;
  margin-top: 12px;
  padding: 4px 8px;
  border-radius: 999px;
  background: var(--visago-surface-soft);
  color: var(--visago-text-muted);
  font-size: 11px;
  font-weight: 800;
}
</style>
