<template>
  <view class="doc-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="隐私政策" />

    <scroll-view scroll-y class="doc-scroll">
      <view class="doc-content visago-page-width">
        <text class="doc-meta">更新日期：2026-05-03</text>

        <view class="doc-card">
          <view v-for="(item, idx) in sections" :key="item.title" class="doc-row">
            <view class="doc-index">{{ item.index }}</view>
            <view class="doc-main">
              <text class="doc-title">{{ item.title }}</text>
              <text class="doc-text">{{ item.text }}</text>
            </view>
            <view v-if="idx < sections.length - 1" class="doc-divider" />
          </view>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../../../components/VisagoTopBar.vue'
import { applyTheme, getStoredTheme } from '../../../../../../utils/theme'

const sections = [
  {
    index: '01',
    title: '收集的信息',
    text: '我们会收集账号信息、基础身份信息、上传材料以及必要的设备日志，范围遵循最小必要原则。',
  },
  {
    index: '02',
    title: '使用目的',
    text: '相关信息用于签证流程辅助、消息通知、风险评估与账号安全保护。',
  },
  {
    index: '03',
    title: '存储与保护',
    text: '平台采用加密传输、权限隔离与访问审计机制保护你的个人信息。',
  },
  {
    index: '04',
    title: '第三方共享',
    text: '除法律法规要求或经你明确授权外，我们不会向无关第三方出售或共享个人信息。',
  },
  {
    index: '05',
    title: '你的权利',
    text: '你可以查询、更正、删除或导出自己的数据，也可以申请注销账户。',
  },
]

export default {
  components: { VisagoTopBar },
  data() {
    return {
      sections,
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

.doc-meta,
.doc-title,
.doc-text {
  display: block;
}

.doc-meta {
  margin: 0 0 10px 4px;
  font-size: 12px;
  color: var(--visago-text-soft);
}

.doc-card {
  overflow: hidden;
  border-radius: 18px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
}

.doc-row {
  position: relative;
  padding: 16px 14px;
  display: flex;
  gap: 12px;
}

.doc-index {
  width: 36px;
  height: 36px;
  border-radius: 12px;
  background: rgba(15, 101, 216, 0.1);
  color: var(--visago-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 12px;
  font-weight: 800;
}

.doc-main {
  min-width: 0;
  flex: 1;
}

.doc-title {
  font-size: 15px;
  font-weight: 800;
  color: var(--visago-text);
}

.doc-text {
  margin-top: 7px;
  font-size: 13px;
  line-height: 1.7;
  color: var(--visago-text-muted);
}

.doc-divider {
  position: absolute;
  left: 62px;
  right: 0;
  bottom: 0;
  height: 1px;
  background: var(--visago-line);
}
</style>
