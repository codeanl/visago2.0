<template>
  <view class="membership-page" :class="{ 'membership-page--dark': themeMode === 'dark' }">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="开通会员" />

    <scroll-view scroll-y class="membership-scroll">
      <view class="membership-content visago-page-width">
        <view class="hero-card">
          <view class="hero-copy">
            <text class="hero-badge">{{ heroBadge }}</text>
            <text class="hero-title">把签证办理升级成一条顺畅流程</text>
            <text class="hero-desc">{{ heroDesc }}</text>
          </view>

          <view class="hero-metrics">
            <view class="metric-item">
              <text class="metric-value">24h</text>
              <text class="metric-label">资料预审反馈</text>
            </view>
            <view class="metric-item">
              <text class="metric-value">1v1</text>
              <text class="metric-label">顾问答疑支持</text>
            </view>
            <view class="metric-item">
              <text class="metric-value">+31%</text>
              <text class="metric-label">资料完善效率</text>
            </view>
          </view>
        </view>

        <view class="benefit-card visago-card">
          <view class="section-head">
            <text class="section-title">会员权益</text>
            <text class="section-sub">围绕签前评估、材料准备、流程跟进三阶段增强</text>
          </view>

          <view class="benefit-list">
            <view v-for="item in benefits" :key="item.title" class="benefit-item">
              <view class="benefit-icon-wrap" :style="{ background: item.tint, color: item.color }">
                <text class="material-symbols-outlined benefit-icon">{{ item.icon }}</text>
              </view>
              <view class="benefit-copy">
                <text class="benefit-title">{{ item.title }}</text>
                <text class="benefit-desc">{{ item.desc }}</text>
              </view>
            </view>
          </view>
        </view>

        <view class="plans-card visago-card">
          <view class="section-head">
            <text class="section-title">选择方案</text>
            <text class="section-sub">{{ planSubText }}</text>
          </view>

          <view class="plan-list">
            <view
              v-for="plan in plans"
              :key="plan.key"
              class="plan-item"
              :class="{ 'plan-item--active': selectedPlan === plan.key, 'plan-item--featured': plan.featured }"
              @tap="selectedPlan = plan.key"
            >
              <view class="plan-main">
                <view class="plan-title-row">
                  <text class="plan-title">{{ plan.title }}</text>
                  <text v-if="plan.tag" class="plan-tag">{{ plan.tag }}</text>
                </view>
                <text class="plan-desc">{{ plan.desc }}</text>
              </view>

              <view class="plan-price-wrap">
                <view class="plan-price-row">
                  <text class="plan-price-symbol">¥</text>
                  <text class="plan-price">{{ plan.price }}</text>
                </view>
                <text class="plan-period">{{ plan.period }}</text>
              </view>
            </view>
          </view>
        </view>

        <view class="compare-card visago-card">
          <view class="section-head">
            <text class="section-title">适合谁开通</text>
          </view>

          <view class="scenario-list">
            <view v-for="scenario in scenarios" :key="scenario.title" class="scenario-item">
              <text class="material-symbols-outlined scenario-icon">{{ scenario.icon }}</text>
              <view class="scenario-copy">
                <text class="scenario-title">{{ scenario.title }}</text>
                <text class="scenario-desc">{{ scenario.desc }}</text>
              </view>
            </view>
          </view>
        </view>

        <view class="faq-card visago-card">
          <view class="section-head">
            <text class="section-title">说明</text>
          </view>

          <view v-for="item in faqList" :key="item.q" class="faq-item">
            <text class="faq-q">{{ item.q }}</text>
            <text class="faq-a">{{ item.a }}</text>
          </view>
        </view>
      </view>
    </scroll-view>

    <view class="cta-bar">
      <view class="cta-copy">
        <text class="cta-label">{{ currentPlan.title }}</text>
        <view class="cta-price-row">
          <text class="cta-price-symbol">¥</text>
          <text class="cta-price">{{ currentPlan.price }}</text>
          <text class="cta-period">{{ currentPlan.period }}</text>
        </view>
      </view>
      <view class="cta-btn" :class="{ 'cta-btn--disabled': submitting }" @tap="subscribeNow">
        {{ submitting ? '提交中...' : ctaText }}
      </view>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { api } from '../../../../utils/api'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'

const plans = [
  {
    key: 'month',
    title: '月度会员',
    desc: '适合近期准备递签，优先拿到预审和顾问响应。',
    price: '8.8',
    period: '/月',
    tag: '',
    featured: false,
  },
  {
    key: 'season',
    title: '季度会员',
    desc: '覆盖完整准备周期，适合首次办理或多轮补件。',
    price: '29.8',
    period: '/季',
    tag: '推荐',
    featured: true,
  },
  {
    key: 'year',
    title: '年度会员',
    desc: '适合多次出行、多成员家庭与长期签证管理。',
    price: '88.8',
    period: '/年',
    tag: '省心',
    featured: false,
  },
]

function formatDate(value) {
  if (!value) return '--'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return String(value)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

export default {
  components: {
    VisagoTopBar,
  },
  data() {
    return {
      themeMode: 'light',
      selectedPlan: 'season',
      submitting: false,
      profile: {},
      benefits: [
        {
          icon: 'workspace_premium',
          title: '优先顾问审核',
          desc: '复杂材料可插队预审，减少因细节导致的返工。',
          tint: 'rgba(59, 130, 246, 0.12)',
          color: '#2563eb',
        },
        {
          icon: 'auto_awesome',
          title: 'AI 风险诊断增强',
          desc: '针对目的地与签证类型给出更细的风险拆解和建议。',
          tint: 'rgba(14, 165, 233, 0.12)',
          color: '#0284c7',
        },
        {
          icon: 'folder_managed',
          title: '资料模板与清单',
          desc: '解锁更多模板、样例与材料收集顺序建议。',
          tint: 'rgba(16, 185, 129, 0.12)',
          color: '#0f9f6e',
        },
        {
          icon: 'support_agent',
          title: '办理过程跟进',
          desc: '从预约到递签的关键节点都能收到更及时提醒。',
          tint: 'rgba(245, 158, 11, 0.14)',
          color: '#d97706',
        },
      ],
      scenarios: [
        {
          icon: 'travel',
          title: '首次办理签证',
          desc: '不熟悉材料顺序、容易漏项的人，会员能显著减少试错成本。',
        },
        {
          icon: 'groups',
          title: '多人同行出行',
          desc: '情侣、家庭、同学团体办理时，统一整理资料更省心。',
        },
        {
          icon: 'event_available',
          title: '时间紧张用户',
          desc: '需要快速补齐材料、尽快拿到建议和提醒的用户更适合开通。',
        },
      ],
      faqList: [
        {
          q: '开通后是否立即生效？',
          a: '支付成功后立即生效，当前账号下的签证计划会自动获得会员加速权益。',
        },
        {
          q: '历史计划是否能继续使用会员能力？',
          a: '可以，已创建的计划仍可使用 AI 评估增强、资料模板和提醒能力。',
        },
        {
          q: '如果短期内只办理一次，选哪档更合适？',
          a: '通常建议先开季度会员，能覆盖准备到递签的完整窗口，性价比更高。',
        },
      ],
      plans,
    }
  },
  computed: {
    membershipInfo() {
      const source = this.profile && typeof this.profile.membership === 'object' && this.profile.membership ? this.profile.membership : {}
      const planKey = String(source.planKey || this.profile.membershipPlanKey || '').trim()
      const planName = String(source.planName || this.profile.membershipPlanName || '').trim()
      const startedAt = String(source.startedAt || this.profile.membershipStartedAt || '').trim()
      const expiresAt = String(source.expiresAt || this.profile.membershipExpiresAt || '').trim()
      let status = String(source.status || this.profile.membershipStatus || '').trim()
      const hasMembership =
        typeof source.hasMembership === 'boolean'
          ? source.hasMembership
          : Boolean(planKey || planName || startedAt || expiresAt || status)
      if (!status && expiresAt) {
        const expireTime = new Date(expiresAt).getTime()
        if (!Number.isNaN(expireTime)) {
          status = expireTime > Date.now() ? 'active' : 'expired'
        }
      }
      return {
        hasMembership,
        planKey,
        planName,
        startedAt,
        expiresAt,
        status,
      }
    },
    currentPlan() {
      return this.plans.find((item) => item.key === this.selectedPlan) || this.plans[0]
    },
    heroBadge() {
      if (this.membershipInfo.status === 'active') {
        return `${this.membershipInfo.planName || 'Visago Pro'} · 已开通`
      }
      if (this.membershipInfo.status === 'expired') {
        return 'Visago Pro · 已过期'
      }
      return 'Visago Pro'
    },
    heroDesc() {
      if (this.membershipInfo.status === 'active') {
        return `当前会员有效期至 ${formatDate(this.membershipInfo.expiresAt)}，后台已同步会员状态与截止时间。`
      }
      return '适合有明确出行计划、希望提升材料通过率并获得更快顾问响应的用户。'
    },
    planSubText() {
      if (this.membershipInfo.status === 'active') {
        return '可直接续费或升级，新的时长会继续同步到数据库和后台管理。'
      }
      return '按当前办理节奏选择，后续可补差价升级。'
    },
    ctaText() {
      return this.membershipInfo.status === 'active' ? '立即续费' : '立即开通'
    },
  },
  onShow() {
    this.themeMode = applyTheme(getStoredTheme())
    this.loadProfile()
  },
  methods: {
    async loadProfile() {
      try {
        this.profile = (await api.me()) || {}
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '加载失败',
          icon: 'none',
        })
      }
    },
    async subscribeNow() {
      if (this.submitting) return
      this.submitting = true
      try {
        this.profile = (await api.subscribeMembership(this.selectedPlan)) || {}
        uni.showToast({
          title: '会员已同步开通',
          icon: 'none',
        })
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '开通失败',
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
.membership-page {
  min-height: 100vh;
  background: var(--visago-bg);
}

.membership-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 88px;
  left: 0;
  height: auto;
}

.membership-content {
  padding: 18px 16px calc(28px + var(--visago-safe-bottom));
  box-sizing: border-box;
}

.hero-card {
  padding: 20px 18px;
  border-radius: 26px;
  background:
    radial-gradient(circle at 12% 18%, rgba(255, 255, 255, 0.34) 0, rgba(255, 255, 255, 0) 26%),
    radial-gradient(circle at 86% 18%, rgba(255, 219, 125, 0.32) 0, rgba(255, 219, 125, 0) 28%),
    linear-gradient(135deg, #2f75ff 0%, #5ea1ff 52%, #7db6ff 100%);
  box-shadow: 0 16px 30px rgba(39, 102, 220, 0.18);
}

.hero-badge,
.hero-title,
.hero-desc,
.section-title,
.section-sub,
.benefit-title,
.benefit-desc,
.plan-title,
.plan-desc,
.plan-period,
.scenario-title,
.scenario-desc,
.faq-q,
.faq-a,
.cta-label,
.cta-period {
  display: block;
}

.hero-badge {
  width: fit-content;
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.18);
  color: #eff6ff;
  font-size: 12px;
  font-weight: 700;
}

.hero-title {
  margin-top: 14px;
  font-size: 22px;
  line-height: 1.28;
  font-weight: 800;
  color: #fff;
}

.hero-desc {
  margin-top: 10px;
  font-size: 13px;
  line-height: 1.6;
  color: rgba(239, 246, 255, 0.9);
}

.hero-metrics {
  margin-top: 18px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.metric-item {
  min-height: 82px;
  padding: 12px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.16);
  backdrop-filter: blur(10px);
}

.metric-value {
  display: block;
  font-size: 24px;
  line-height: 1;
  font-weight: 800;
  color: #fff;
}

.metric-label {
  display: block;
  margin-top: 10px;
  font-size: 12px;
  line-height: 1.4;
  color: rgba(239, 246, 255, 0.88);
}

.benefit-card,
.plans-card,
.compare-card,
.faq-card {
  margin-top: 14px;
  padding: 18px 16px;
}

.section-head {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.section-title {
  font-size: 16px;
  font-weight: 800;
  color: var(--visago-text);
}

.section-sub {
  font-size: 12px;
  line-height: 1.5;
  color: var(--visago-text-soft);
}

.benefit-list {
  margin-top: 14px;
  display: grid;
  gap: 12px;
}

.benefit-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.benefit-icon-wrap {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.benefit-icon {
  font-size: 22px;
}

.benefit-copy {
  min-width: 0;
  flex: 1;
}

.benefit-title {
  font-size: 14px;
  font-weight: 700;
  color: var(--visago-text);
}

.benefit-desc {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.6;
  color: var(--visago-text-soft);
}

.plan-list {
  margin-top: 14px;
  display: grid;
  gap: 12px;
}

.plan-item {
  padding: 16px 14px;
  border-radius: 18px;
  border: 1px solid var(--visago-line);
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.plan-item--featured {
  background: linear-gradient(135deg, rgba(239, 246, 255, 0.94) 0%, rgba(230, 241, 255, 0.98) 100%);
}

.plan-item--active {
  border-color: rgba(37, 99, 235, 0.32);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.plan-main {
  min-width: 0;
  flex: 1;
}

.plan-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.plan-title {
  font-size: 15px;
  font-weight: 800;
  color: var(--visago-text);
}

.plan-tag {
  padding: 2px 7px;
  border-radius: 999px;
  background: rgba(37, 99, 235, 0.12);
  color: #2563eb;
  font-size: 11px;
  font-weight: 700;
}

.plan-desc {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.55;
  color: var(--visago-text-soft);
}

.plan-price-wrap {
  text-align: right;
  flex-shrink: 0;
}

.plan-price-row,
.cta-price-row {
  display: flex;
  align-items: flex-end;
  justify-content: flex-end;
  gap: 2px;
}

.plan-price-symbol,
.cta-price-symbol {
  font-size: 13px;
  font-weight: 700;
  color: #2563eb;
}

.plan-price {
  font-size: 28px;
  line-height: 1;
  font-weight: 800;
  color: #2563eb;
}

.plan-period {
  margin-top: 4px;
}

.scenario-list {
  margin-top: 14px;
  display: grid;
  gap: 12px;
}

.scenario-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.scenario-icon {
  width: 40px;
  height: 40px;
  border-radius: 14px;
  background: rgba(59, 130, 246, 0.12);
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.scenario-copy {
  min-width: 0;
  flex: 1;
}

.scenario-title {
  font-size: 14px;
  font-weight: 700;
  color: var(--visago-text);
}

.scenario-desc {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.6;
  color: var(--visago-text-soft);
}

.faq-item + .faq-item {
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px solid var(--visago-line);
}

.faq-q {
  font-size: 14px;
  font-weight: 700;
  color: var(--visago-text);
}

.faq-a {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.65;
  color: var(--visago-text-soft);
}

.cta-bar {
  position: fixed;
  left: 12px;
  right: 12px;
  bottom: calc(12px + var(--visago-safe-bottom));
  z-index: 30;
  padding: 12px 12px 12px 14px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.94);
  border: 1px solid rgba(217, 226, 239, 0.9);
  box-shadow: 0 16px 34px rgba(15, 23, 42, 0.1);
  backdrop-filter: blur(16px);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.cta-copy {
  min-width: 0;
}

.cta-label {
  font-size: 13px;
  color: var(--visago-text-soft);
}

.cta-price {
  font-size: 26px;
  line-height: 1;
  font-weight: 800;
  color: var(--visago-text);
}

.cta-period {
  margin-left: 4px;
  margin-bottom: 2px;
  font-size: 12px;
  color: var(--visago-text-soft);
}

.cta-btn {
  height: 46px;
  min-width: 118px;
  padding: 0 18px;
  border-radius: 14px;
  background: linear-gradient(135deg, #2f75ff 0%, #5a9bff 100%);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 15px;
  font-weight: 700;
  flex-shrink: 0;
}

.cta-btn--disabled {
  opacity: 0.72;
}

.membership-page--dark .hero-card,
:global(html.theme-dark) .hero-card {
  box-shadow: 0 18px 34px rgba(0, 0, 0, 0.32);
}

.membership-page--dark .metric-item,
:global(html.theme-dark) .metric-item {
  background: rgba(8, 15, 28, 0.22);
  border-color: rgba(255, 255, 255, 0.1);
}

.membership-page--dark .plan-item,
:global(html.theme-dark) .plan-item {
  background: rgba(18, 24, 35, 0.96);
}

.membership-page--dark .plan-item--featured,
:global(html.theme-dark) .plan-item--featured {
  background: linear-gradient(135deg, rgba(17, 26, 43, 0.96) 0%, rgba(21, 35, 58, 0.94) 100%);
}

.membership-page--dark .plan-item--active,
:global(html.theme-dark) .plan-item--active {
  box-shadow: 0 0 0 3px rgba(96, 165, 250, 0.12);
}

.membership-page--dark .scenario-icon,
:global(html.theme-dark) .scenario-icon {
  background: rgba(59, 130, 246, 0.18);
}

.membership-page--dark .cta-bar,
:global(html.theme-dark) .cta-bar {
  background: rgba(12, 18, 28, 0.94);
  border-color: rgba(48, 58, 74, 0.92);
  box-shadow: 0 18px 34px rgba(0, 0, 0, 0.28);
}
</style>
