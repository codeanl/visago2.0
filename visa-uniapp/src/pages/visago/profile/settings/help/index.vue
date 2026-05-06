<template>
  <view class="help-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="帮助中心" />

    <scroll-view scroll-y class="help-scroll">
      <view class="help-content visago-page-width">
        <view class="faq-card">
          <view v-for="item in faqs" :key="item.q" class="faq-row" @tap="toggleFaq(item.q)">
            <view class="faq-head">
              <text class="faq-q">{{ item.q }}</text>
              <text class="material-symbols-outlined faq-arrow">{{ openFaq === item.q ? 'expand_less' : 'expand_more' }}</text>
            </view>
            <text v-if="openFaq === item.q" class="faq-a">{{ item.a }}</text>
          </view>
        </view>

        <text class="help-note">如仍未解决，可在“关于我们”中查看联系邮箱。</text>
      </view>
    </scroll-view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../../components/VisagoTopBar.vue'
import { applyTheme, getStoredTheme } from '../../../../../utils/theme'

export default {
  components: { VisagoTopBar },
  data() {
    return {
      openFaq: '签证材料多久会审核完成？',
      faqs: [
        { q: '签证材料多久会审核完成？', a: '普通材料审核通常需要 1 到 2 个工作日。遇到高峰期或材料缺失时，系统会在计划页提示补充。' },
        { q: '如何修改已经提交的资料？', a: '进入计划页，找到对应任务后重新上传或编辑材料即可。' },
        { q: '为什么 AI 预测结果和实际出签不同？', a: 'AI 预测只做辅助评估，最终结果仍以使领馆审核为准。' },
        { q: '可以删除我的个人数据吗？', a: '可以，在“隐私与安全”里可申请下载数据或注销账户。' },
      ],
    }
  },
  onShow() {
    applyTheme(getStoredTheme())
  },
  methods: {
    toggleFaq(q) {
      this.openFaq = this.openFaq === q ? '' : q
    },
  },
}
</script>

<style scoped>
.help-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.help-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.help-content {
  padding: 18px 16px 34px;
  box-sizing: border-box;
}

.faq-card {
  overflow: hidden;
  border-radius: 18px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
}

.faq-row {
  padding: 14px;
  border-bottom: 1px solid var(--visago-line);
}

.faq-row:last-child {
  border-bottom: none;
}

.faq-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.faq-q,
.faq-a,
.help-note {
  display: block;
}

.faq-q {
  font-size: 15px;
  line-height: 1.45;
  font-weight: 800;
  color: var(--visago-text);
}

.faq-arrow {
  color: var(--visago-text-soft);
}

.faq-a {
  margin-top: 9px;
  color: var(--visago-text-muted);
  font-size: 13px;
  line-height: 1.6;
}

.help-note {
  margin-top: 12px;
  padding: 0 4px;
  font-size: 12px;
  line-height: 1.5;
  color: var(--visago-text-soft);
}
</style>
