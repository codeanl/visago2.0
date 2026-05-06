<template>
  <view class="doc-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="用户协议" />

    <scroll-view scroll-y class="doc-scroll">
      <view class="doc-content visago-page-width">
        <text class="doc-meta">生效日期：2026-05-03</text>

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
    title: '服务说明',
    text: 'Visago 提供签证信息查询、材料建议、流程跟进与提醒等工具能力，页面内容仅作辅助参考。',
  },
  {
    index: '02',
    title: '用户责任',
    text: '你需要保证提交资料真实、完整、有效，因资料问题导致的审核失败或延误由用户自行承担。',
  },
  {
    index: '03',
    title: '使用规范',
    text: '不得利用本服务从事违法活动，也不得攻击、抓取或批量传播平台内容。',
  },
  {
    index: '04',
    title: '免责说明',
    text: '签证政策与使领馆要求可能变化，最终结果以官方通知为准。',
  },
  {
    index: '05',
    title: '协议更新',
    text: '协议调整后会在本页面更新展示，继续使用即视为接受更新内容。',
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
