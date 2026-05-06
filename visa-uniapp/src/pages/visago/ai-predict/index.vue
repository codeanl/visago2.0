<template>
  <view class="ai-page">
    <view v-if="screen === 'form'" class="form-screen">
      <VisagoTopBar :show-back="true" :show-notice="false" page-name="AI预测" />

      <view class="stepper">
        <view v-for="(step, index) in formSteps" :key="step.key" class="stepper-item" :class="{ 'stepper-item--active': index === formStepIndex, 'stepper-item--done': index < formStepIndex }">
          <view class="stepper-bar" />
          <text class="stepper-label">{{ step.code }} {{ step.label }}</text>
        </view>
      </view>

      <scroll-view scroll-y class="form-scroll">
        <view class="form-content">
          <view class="title-block">
            <text class="page-title">{{ currentFormStep.title }}</text>
            <text class="page-subtitle">{{ currentFormStep.sub }}</text>
          </view>

          <view class="ai-tip">
            <view class="ai-tip-icon">
              <text class="material-symbols-outlined">auto_awesome</text>
            </view>
            <view class="ai-tip-copy">
              <text class="ai-tip-title">AI 智能校验</text>
              <text class="ai-tip-text">我们会根据全球签证要求校验资料细节，并预测你的通过率和材料风险。</text>
            </view>
          </view>

          <view v-if="currentFormStep.key === 'personal'" class="field-list">
            <view class="field-block">
              <text class="field-label">评估模式</text>
              <view class="mode-tabs">
                <view class="mode-tab" :class="{ 'mode-tab--active': form.mode === 'country' }" @tap="setMode('country')">指定国家</view>
                <view class="mode-tab" :class="{ 'mode-tab--active': form.mode === 'discover' }" @tap="setMode('discover')">推荐目的地</view>
              </view>
            </view>

            <view v-if="form.mode === 'country'" class="field-block">
              <text class="field-label">目标国家</text>
              <picker :range="countryOptions" :value="form.countryIndex" @change="onCountryChange">
                <view class="input-like">
                  <text>{{ countryOptions[form.countryIndex] }}</text>
                  <text class="material-symbols-outlined input-icon">expand_more</text>
                </view>
              </picker>
            </view>

            <view class="field-block">
              <text class="field-label">签证类型</text>
              <picker :range="visaTypeOptions" :value="form.visaTypeIndex" @change="onVisaTypeChange">
                <view class="input-like">
                  <text>{{ visaTypeOptions[form.visaTypeIndex] }}</text>
                  <text class="material-symbols-outlined input-icon">expand_more</text>
                </view>
              </picker>
            </view>

            <view class="field-block">
              <text class="field-label">职业 / 身份</text>
              <input v-model.trim="form.occupation" class="input-like" placeholder="例如：产品经理 / 学生 / 自由职业" />
            </view>
          </view>

          <view v-if="currentFormStep.key === 'financial'" class="field-list">
            <view class="field-block">
              <text class="field-label">月收入区间</text>
              <picker :range="incomeOptions" :value="form.incomeIndex" @change="onIncomeChange">
                <view class="input-like">
                  <text>{{ incomeOptions[form.incomeIndex] }}</text>
                  <text class="material-symbols-outlined input-icon">expand_more</text>
                </view>
              </picker>
            </view>

            <view class="field-block">
              <text class="field-label">近 6 月流水</text>
              <picker :range="fundOptions" :value="form.fundIndex" @change="onFundChange">
                <view class="input-like">
                  <text>{{ fundOptions[form.fundIndex] }}</text>
                  <text class="material-symbols-outlined input-icon">expand_more</text>
                </view>
              </picker>
            </view>

            <view class="field-block">
              <text class="field-label">在职 / 在读稳定性</text>
              <picker :range="tenureOptions" :value="form.tenureIndex" @change="onTenureChange">
                <view class="input-like">
                  <text>{{ tenureOptions[form.tenureIndex] }}</text>
                  <text class="material-symbols-outlined input-icon">expand_more</text>
                </view>
              </picker>
            </view>
          </view>

          <view v-if="currentFormStep.key === 'travel'" class="field-list">
            <view class="field-block">
              <text class="field-label">出行月份</text>
              <input v-model.trim="form.tripMonth" class="input-like" placeholder="例如：2026-09" />
            </view>

            <view class="field-block">
              <text class="field-label">预计停留</text>
              <picker :range="stayOptions" :value="form.stayIndex" @change="onStayChange">
                <view class="input-like">
                  <text>{{ stayOptions[form.stayIndex] }}</text>
                  <text class="material-symbols-outlined input-icon">expand_more</text>
                </view>
              </picker>
            </view>

            <view class="field-block">
              <text class="field-label">发达国家出行记录</text>
              <picker :range="travelHistoryOptions" :value="form.travelHistoryIndex" @change="onTravelHistoryChange">
                <view class="input-like">
                  <text>{{ travelHistoryOptions[form.travelHistoryIndex] }}</text>
                  <text class="material-symbols-outlined input-icon">expand_more</text>
                </view>
              </picker>
            </view>

            <view class="security-card">
              <text class="material-symbols-outlined security-icon">privacy_tip</text>
              <view>
                <text class="security-title">数据保护</text>
                <text class="security-text">你的信息仅用于签证通行力预测，不会分享给第三方。</text>
              </view>
            </view>
          </view>

          <view v-if="currentFormStep.key === 'docs'" class="field-list">
            <view class="switch-card">
              <view>
                <text class="switch-title">是否有拒签记录</text>
                <text class="switch-desc">拒签记录会影响整体风险判断</text>
              </view>
              <switch :checked="form.hasRefusal" color="#0f65d8" @change="onRefusalChange" />
            </view>

            <view class="doc-list">
              <view v-for="item in materialOptions" :key="item.key" class="doc-item" :class="{ 'doc-item--active': form.materials[item.key] }" @tap="toggleMaterial(item.key)">
                <view class="doc-icon-box">
                  <text class="material-symbols-outlined">{{ item.icon }}</text>
                </view>
                <view class="doc-copy">
                  <text class="doc-title">{{ item.label }}</text>
                  <text class="doc-desc">{{ item.desc }}</text>
                </view>
                <text class="material-symbols-outlined doc-check">{{ form.materials[item.key] ? 'check_circle' : 'radio_button_unchecked' }}</text>
              </view>
            </view>
          </view>
        </view>
      </scroll-view>

      <view class="footer">
        <view v-if="formStepIndex > 0" class="footer-secondary" @tap="prevStep">上一步</view>
        <view class="footer-primary" :class="{ 'footer-primary--wide': formStepIndex === 0 }" @tap="nextStep">
          <text>{{ formStepIndex === formSteps.length - 1 ? '生成报告' : '下一步' }}</text>
          <text class="material-symbols-outlined footer-icon">arrow_forward</text>
        </view>
      </view>
    </view>

    <view v-else class="report-screen">
      <VisagoTopBar :show-back="true" :custom-back="true" :show-notice="false" page-name="AI智能评估报告" @back="backToForm" />

      <scroll-view scroll-y class="report-scroll">
        <view class="score-panel">
          <view class="score-gauge" :style="scoreRingStyle">
            <view class="score-inner">
              <text class="score-number">{{ report.score }}%</text>
              <text class="score-label">Prediction Score</text>
            </view>
          </view>
          <text class="score-title">{{ reportHeadline }}</text>
          <text class="score-desc">{{ reportSummary }}</text>
        </view>

        <view class="report-section">
          <text class="section-title">核心维度诊断</text>
          <view v-for="item in report.breakdown" :key="item.key" class="dimension-card">
            <view class="dimension-head">
              <text class="dimension-title">{{ item.title }}</text>
              <text class="dimension-tag" :class="`dimension-tag--${item.level}`">{{ item.label }}</text>
            </view>
            <view class="dimension-track">
              <view class="dimension-fill" :class="`dimension-fill--${item.level}`" :style="{ width: item.score + '%' }" />
            </view>
            <text class="dimension-note">{{ item.note }}</text>
          </view>
        </view>

        <view v-if="report.mode === 'discover'" class="report-section">
          <text class="section-title">推荐目的地</text>
          <view v-for="item in report.recommendations" :key="item.country" class="recommend-card">
            <view>
              <text class="recommend-country">{{ item.country }}</text>
              <text class="recommend-reason">{{ item.reason }}</text>
            </view>
            <text class="recommend-score">{{ item.match }}%</text>
          </view>
        </view>

        <view class="report-section">
          <text class="section-title">材料风险预警</text>
          <view class="risk-box">
            <view v-for="risk in report.risks" :key="risk.title" class="risk-item">
              <text class="material-symbols-outlined risk-icon">{{ risk.icon }}</text>
              <view>
                <text class="risk-title">{{ risk.title }}</text>
                <text class="risk-desc">{{ risk.desc }}</text>
              </view>
            </view>
          </view>
        </view>

        <view class="report-section">
          <text class="section-title">AI 提升方案</text>
          <view v-for="todo in report.todos" :key="todo.title" class="improve-card">
            <view class="improve-icon">
              <text class="material-symbols-outlined">{{ todo.icon }}</text>
            </view>
            <view class="improve-copy">
              <text class="improve-title">{{ todo.title }}</text>
              <text class="improve-desc">{{ todo.desc }}</text>
            </view>
            <text class="material-symbols-outlined improve-arrow">chevron_right</text>
          </view>
        </view>

        <view class="report-actions">
          <view class="start-btn" @tap="saveToPlan">
            <text>确定并开始申请</text>
            <text class="material-symbols-outlined">rocket_launch</text>
          </view>
          <text class="powered">Powered by Visago AI</text>
          <text class="generated">报告生成于 {{ report.createdAt }}</text>
        </view>
      </scroll-view>
    </view>
  </view>
</template>

<script>
import { applyTheme, getStoredTheme } from '../../../utils/theme'
import VisagoTopBar from '../../../components/VisagoTopBar.vue'

const FORM_KEY = 'visago_ai_predict_form_test_style'
const REPORT_KEY = 'visago_ai_predict_report_test_style'

const COUNTRY_POOL = [
  { country: '日本', difficulty: 8, tags: ['旅游签', '商务签'] },
  { country: '新加坡', difficulty: 6, tags: ['旅游签', '商务签'] },
  { country: '泰国', difficulty: 5, tags: ['旅游签', '商务签'] },
  { country: '法国', difficulty: 13, tags: ['旅游签', '商务签', '学生签'] },
  { country: '英国', difficulty: 15, tags: ['旅游签', '商务签', '学生签', '工作签'] },
  { country: '澳大利亚', difficulty: 12, tags: ['旅游签', '商务签', '学生签', '工作签'] },
  { country: '加拿大', difficulty: 16, tags: ['旅游签', '学生签', '工作签'] },
]

function defaultMaterials() {
  return {
    passport: false,
    photo: false,
    bank: false,
    itinerary: false,
    employment: false,
    assets: false,
  }
}

function clamp(value, min, max) {
  return Math.max(min, Math.min(max, value))
}

function nowText() {
  const date = new Date()
  const y = date.getFullYear()
  const m = `${date.getMonth() + 1}`.padStart(2, '0')
  const d = `${date.getDate()}`.padStart(2, '0')
  const h = `${date.getHours()}`.padStart(2, '0')
  const min = `${date.getMinutes()}`.padStart(2, '0')
  return `${y}-${m}-${d} ${h}:${min}`
}

export default {
  components: {
    VisagoTopBar,
  },
  data() {
    return {
      screen: 'form',
      formStepIndex: 0,
      formSteps: [
        { key: 'personal', code: '01', label: 'PERSONAL', title: '个人详细信息', sub: '请选择申请目标，并提供与申请一致的身份信息。' },
        { key: 'financial', code: '02', label: 'FINANCIAL', title: '资金与稳定性', sub: '资金能力和身份稳定性会影响签证官判断。' },
        { key: 'travel', code: '03', label: 'TRAVEL', title: '出行计划', sub: '请提供明确的出行时间、停留周期和过往旅行记录。' },
        { key: 'docs', code: '04', label: 'DOCS', title: '材料准备情况', sub: '勾选已经准备好的材料，AI 将识别风险缺口。' },
      ],
      countryOptions: ['日本', '英国', '澳大利亚', '法国', '德国', '新加坡', '加拿大'],
      visaTypeOptions: ['旅游签', '商务签', '学生签', '工作签'],
      stayOptions: ['7天以内', '8-15天', '16-30天', '30天以上'],
      tenureOptions: ['6个月以内', '6个月-2年', '2年以上'],
      incomeOptions: ['8千以下', '8千-2万', '2万-5万', '5万以上'],
      fundOptions: ['5万以下', '5万-15万', '15万-30万', '30万以上'],
      travelHistoryOptions: ['无', '1-2次', '3次以上'],
      materialOptions: [
        { key: 'passport', label: '护照扫描件', desc: '有效期至少6个月', icon: 'badge' },
        { key: 'photo', label: '证件照片', desc: '符合目标国家规格', icon: 'photo_camera' },
        { key: 'bank', label: '银行流水', desc: '近6个月稳定记录', icon: 'account_balance' },
        { key: 'itinerary', label: '行程单', desc: '机票、酒店、每日安排', icon: 'route' },
        { key: 'employment', label: '在职/在读证明', desc: '增强回国约束力', icon: 'work' },
        { key: 'assets', label: '资产证明', desc: '房产、车产、存款等', icon: 'payments' },
      ],
      form: {
        mode: 'country',
        countryIndex: 0,
        visaTypeIndex: 0,
        occupation: '',
        tenureIndex: 1,
        incomeIndex: 1,
        fundIndex: 1,
        tripMonth: '',
        stayIndex: 1,
        travelHistoryIndex: 0,
        hasRefusal: false,
        materials: defaultMaterials(),
      },
      report: null,
    }
  },
  computed: {
    currentFormStep() {
      return this.formSteps[this.formStepIndex]
    },
    reportHeadline() {
      if (!this.report) return ''
      if (this.report.score >= 88) return '极高出签概率'
      if (this.report.score >= 76) return '较高出签概率'
      if (this.report.score >= 62) return '中等出签概率'
      return '存在明显补强空间'
    },
    reportSummary() {
      if (!this.report) return ''
      if (this.report.mode === 'discover') return '根据当前资料质量，以下国家更适合作为优先申请目标。'
      return `基于${this.report.country}${this.report.visaType}要求，你的资料质量处于${this.report.levelText}水平。`
    },
    scoreRingStyle() {
      const score = this.report ? this.report.score : 0
      return {
        background: `conic-gradient(#0f73ed 0% ${score}%, #e8edf5 ${score}% 100%)`,
      }
    },
  },
  onShow() {
    applyTheme(getStoredTheme())
    this.restore()
  },
  methods: {
    goBack() {
      uni.navigateBack({ fail: () => uni.reLaunch({ url: '/pages/visago/home/index' }) })
    },
    setMode(mode) {
      this.form.mode = mode
      this.persistForm()
    },
    onCountryChange(e) {
      this.form.countryIndex = Number(e.detail.value || 0)
      this.persistForm()
    },
    onVisaTypeChange(e) {
      this.form.visaTypeIndex = Number(e.detail.value || 0)
      this.persistForm()
    },
    onTenureChange(e) {
      this.form.tenureIndex = Number(e.detail.value || 0)
      this.persistForm()
    },
    onIncomeChange(e) {
      this.form.incomeIndex = Number(e.detail.value || 0)
      this.persistForm()
    },
    onFundChange(e) {
      this.form.fundIndex = Number(e.detail.value || 0)
      this.persistForm()
    },
    onStayChange(e) {
      this.form.stayIndex = Number(e.detail.value || 0)
      this.persistForm()
    },
    onTravelHistoryChange(e) {
      this.form.travelHistoryIndex = Number(e.detail.value || 0)
      this.persistForm()
    },
    onRefusalChange(e) {
      this.form.hasRefusal = !!e.detail.value
      this.persistForm()
    },
    toggleMaterial(key) {
      this.form.materials[key] = !this.form.materials[key]
      this.persistForm()
    },
    prevStep() {
      this.formStepIndex = Math.max(0, this.formStepIndex - 1)
      this.persistForm()
    },
    nextStep() {
      const error = this.validateStep()
      if (error) {
        uni.showToast({ title: error, icon: 'none' })
        return
      }
      if (this.formStepIndex === this.formSteps.length - 1) {
        this.report = this.buildReport()
        this.screen = 'report'
        uni.setStorageSync(REPORT_KEY, this.report)
        return
      }
      this.formStepIndex += 1
      this.persistForm()
    },
    validateStep() {
      const key = this.currentFormStep.key
      if (key === 'personal' && (!this.form.occupation || this.form.occupation.length < 2)) return '请填写职业/身份'
      if (key === 'travel' && !/^\d{4}-(0[1-9]|1[0-2])$/.test(this.form.tripMonth)) return '出行月份格式应为 YYYY-MM'
      if (key === 'docs') {
        if (!this.form.materials.passport) return '请先准备护照扫描件'
        if (!this.form.materials.photo) return '请先准备证件照片'
        if (Object.values(this.form.materials).filter(Boolean).length < 3) return '请至少勾选3项材料'
      }
      return ''
    },
    buildScores() {
      const count = Object.values(this.form.materials).filter(Boolean).length
      let finance = 62 + [-10, -3, 8, 14][this.form.fundIndex] + [-6, 0, 6, 10][this.form.incomeIndex]
      let compliance = 48 + count * 8
      let travel = 58 + [0, 8, 14][this.form.travelHistoryIndex]
      let career = 60 + [0, 7, 13][this.form.tenureIndex]
      if (this.form.hasRefusal) travel -= 14
      if (!this.form.materials.bank) finance -= 12
      if (!this.form.materials.itinerary) travel -= 10
      if (!this.form.materials.employment) career -= 8
      return {
        finance: clamp(finance, 20, 98),
        compliance: clamp(compliance, 20, 98),
        travel: clamp(travel, 20, 98),
        career: clamp(career, 20, 98),
      }
    },
    difficulty() {
      const country = this.countryOptions[this.form.countryIndex]
      const found = COUNTRY_POOL.find((item) => item.country === country)
      return found ? found.difficulty : 10
    },
    buildReport() {
      const scores = this.buildScores()
      let score = Math.round(scores.finance * 0.3 + scores.compliance * 0.3 + scores.travel * 0.18 + scores.career * 0.22)
      const visaType = this.visaTypeOptions[this.form.visaTypeIndex]
      if (visaType === '学生签') score -= 4
      if (visaType === '工作签') score -= 6
      if (visaType === '商务签') score += 2
      score -= this.form.mode === 'country' ? this.difficulty() : 4
      score = clamp(score, 28, 96)
      const levelText = score >= 85 ? '领先' : score >= 72 ? '较优' : score >= 60 ? '中等' : '待补强'
      return {
        mode: this.form.mode,
        country: this.countryOptions[this.form.countryIndex],
        visaType,
        score,
        levelText,
        createdAt: this.nowText(),
        breakdown: [
          this.dimension('经济实力', scores.finance),
          this.dimension('资料合规性', scores.compliance),
          this.dimension('过往旅行史', scores.travel),
          this.dimension('职业稳定性', scores.career),
        ],
        risks: this.buildRisks(scores),
        todos: this.buildTodos(),
        recommendations: this.form.mode === 'discover' ? this.buildRecommendations(score, visaType) : [],
      }
    },
    dimension(title, score) {
      const level = score >= 85 ? 'strong' : score >= 70 ? 'good' : score >= 55 ? 'mid' : 'weak'
      const labelMap = { strong: '极强', good: '优秀', mid: '中等', weak: '需补强' }
      const noteMap = {
        strong: '该项明显优于常规申请水平。',
        good: '该项整体稳定，递交前保持一致即可。',
        mid: '该项存在轻微不足，建议补充说明材料。',
        weak: '该项风险较高，建议优先补强。',
      }
      return { key: title, title, score, level, label: labelMap[level], note: noteMap[level] }
    },
    buildRisks(scores) {
      const risks = []
      if (scores.finance < 72) risks.push({ title: '银行流水', desc: '建议补充近6个月工资对账单或存款证明，增强收入稳定性。', icon: 'account_balance_wallet' })
      if (scores.travel < 70) risks.push({ title: '行程单', desc: '每日行程可进一步细化，特别是住宿与景点之间的衔接。', icon: 'route' })
      if (this.form.hasRefusal) risks.push({ title: '拒签说明', desc: '需要提供解释信，说明上次拒签原因以及本次资料变化。', icon: 'warning' })
      if (!risks.length) risks.push({ title: '资料一致性', desc: '当前风险较低，递交前重点核对姓名、日期、证件号一致。', icon: 'verified' })
      return risks.slice(0, 3)
    },
    buildTodos() {
      const todos = []
      if (!this.form.materials.bank) todos.push({ title: '补充资料', desc: '提供一份详细的银行流水与资金说明，证明本次行程支付能力。', icon: 'description' })
      if (!this.form.materials.itinerary) todos.push({ title: '完善行程', desc: '补齐机票、酒店和每日安排，提升出行真实性。', icon: 'map' })
      if (!this.form.materials.employment) todos.push({ title: '专家审核', desc: '建议在递交前人工预审，重点检查身份稳定性材料。', icon: 'support_agent' })
      if (!todos.length) todos.push({ title: '开始申请', desc: '材料质量良好，可进入预约递交阶段。', icon: 'rocket_launch' })
      return todos.slice(0, 3)
    },
    buildRecommendations(score, visaType) {
      return COUNTRY_POOL.map((item) => {
        let match = score - item.difficulty + (item.tags.includes(visaType) ? 8 : 0)
        match = clamp(match, 38, 95)
        let reason = '资料匹配度较高，可作为优先申请目标。'
        if (match < 65) reason = '建议先补充资金和行程材料后再申请。'
        return { country: item.country, match, reason }
      }).sort((a, b) => b.match - a.match).slice(0, 4)
    },
    backToForm() {
      this.screen = 'form'
    },
    saveToPlan() {
      const country = this.report && this.report.mode === 'discover' && this.report.recommendations[0]
        ? this.report.recommendations[0].country
        : this.countryOptions[this.form.countryIndex]
      const visaType = this.visaTypeOptions[this.form.visaTypeIndex]
      uni.navigateTo({ url: `/pages/visago/plan/index?countryName=${encodeURIComponent(country)}&visaTitle=${encodeURIComponent(visaType)}` })
    },
    nowText() {
      const d = new Date()
      const y = d.getFullYear()
      const m = `${d.getMonth() + 1}`.padStart(2, '0')
      const day = `${d.getDate()}`.padStart(2, '0')
      const h = `${d.getHours()}`.padStart(2, '0')
      const min = `${d.getMinutes()}`.padStart(2, '0')
      return `${y}-${m}-${day} ${h}:${min}`
    },
    persistForm() {
      uni.setStorageSync(FORM_KEY, { form: this.form, formStepIndex: this.formStepIndex })
    },
    restore() {
      try {
        const saved = uni.getStorageSync(FORM_KEY)
        if (saved && saved.form) {
          this.form = { ...this.form, ...saved.form, materials: { ...defaultMaterials(), ...(saved.form.materials || {}) } }
          this.formStepIndex = Number(saved.formStepIndex || 0)
        }
        const report = uni.getStorageSync(REPORT_KEY)
        if (report && typeof report === 'object') this.report = report
      } catch (error) {
      }
    },
  },
}
</script>

<style scoped>
.ai-page {
  --ai-topbar-height: 74px;
  --ai-stepper-height: 76px;
  --ai-footer-height: calc(82px + var(--visago-safe-bottom));
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.stepper {
  position: fixed;
  top: var(--ai-topbar-height);
  left: 0;
  right: 0;
  z-index: 30;
  height: var(--ai-stepper-height);
  padding: 18px 16px 10px;
  background: var(--visago-bg);
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
  box-sizing: border-box;
  border-bottom: 1px solid var(--visago-line);
}

.stepper-bar {
  height: 4px;
  border-radius: 999px;
  background: var(--visago-line);
}

.stepper-item--active .stepper-bar,
.stepper-item--done .stepper-bar {
  background: #0f73ed;
}

.stepper-label {
  display: block;
  margin-top: 6px;
  font-size: 11px;
  font-weight: 700;
  color: var(--visago-text-soft);
  white-space: nowrap;
}

.stepper-item--active .stepper-label,
.stepper-item--done .stepper-label {
  color: #0f65d8;
}

.form-scroll {
  position: fixed;
  top: calc(var(--ai-topbar-height) + var(--ai-stepper-height));
  right: 0;
  bottom: var(--ai-footer-height);
  left: 0;
  box-sizing: border-box;
  height: auto;
}

.form-content,
.report-scroll {
  padding: 0 16px;
  box-sizing: border-box;
}

.form-content {
  padding-bottom: 22px;
}

.title-block {
  padding-top: 18px;
}

.page-title {
  display: block;
  font-size: 31px;
  line-height: 1.18;
  font-weight: 800;
  color: var(--visago-text);
}

.page-subtitle {
  display: block;
  margin-top: 8px;
  font-size: 16px;
  line-height: 1.5;
  color: var(--visago-text-muted);
}

.ai-tip,
.security-card {
  margin-top: 28px;
  border-radius: 14px;
  border: 1px solid var(--visago-line);
  background: var(--visago-surface);
  padding: 16px;
  display: flex;
  gap: 14px;
  align-items: flex-start;
}

.ai-tip-icon {
  width: 58px;
  height: 58px;
  border-radius: 12px;
  background: #e6f0ff;
  color: #0f73ed;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.ai-tip-icon .material-symbols-outlined {
  font-size: 30px;
  font-variation-settings: 'FILL' 1;
}

.ai-tip-title,
.security-title {
  display: block;
  font-size: 16px;
  font-weight: 700;
}

.ai-tip-text,
.security-text {
  display: block;
  margin-top: 4px;
  font-size: 13px;
  line-height: 1.5;
  color: var(--visago-text-muted);
}

.field-list {
  margin-top: 28px;
  display: grid;
  gap: 20px;
}

.field-label {
  display: block;
  margin: 0 0 8px 8px;
  font-size: 14px;
  color: var(--visago-text);
}

.input-like {
  width: 100%;
  height: 54px;
  border-radius: 13px;
  background: var(--visago-surface-soft);
  border: none;
  padding: 0 16px;
  box-sizing: border-box;
  font-size: 17px;
  color: var(--visago-text);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.input-icon {
  color: var(--visago-text-soft);
}

.mode-tabs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.mode-tab {
  height: 46px;
  border-radius: 12px;
  background: var(--visago-surface-soft);
  color: var(--visago-text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
}

.mode-tab--active {
  background: #0f73ed;
  color: #fff;
}

.switch-card,
.doc-item {
  border-radius: 14px;
  border: 1px solid var(--visago-line);
  background: var(--visago-surface);
  padding: 14px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.switch-card {
  justify-content: space-between;
}

.switch-title,
.doc-title {
  display: block;
  font-size: 15px;
  font-weight: 700;
}

.switch-desc,
.doc-desc {
  display: block;
  margin-top: 3px;
  font-size: 12px;
  color: var(--visago-text-muted);
}

.doc-list {
  display: grid;
  gap: 10px;
}

.doc-item--active {
  border-color: #bdd9ff;
  background: color-mix(in srgb, var(--visago-primary) 8%, var(--visago-surface));
}

.doc-icon-box {
  width: 42px;
  height: 42px;
  border-radius: 11px;
  background: #eaf3ff;
  color: #0f73ed;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.doc-copy {
  flex: 1;
}

.doc-check {
  color: #0f73ed;
}

.footer {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  min-height: var(--ai-footer-height);
  padding: 14px 16px calc(12px + var(--visago-safe-bottom));
  background: color-mix(in srgb, var(--visago-surface) 94%, transparent);
  border-top: 1px solid var(--visago-line);
  backdrop-filter: blur(14px);
  display: grid;
  grid-template-columns: 104px 1fr;
  gap: 10px;
  box-sizing: border-box;
  z-index: 50;
}

.footer-secondary,
.footer-primary {
  height: 54px;
  border-radius: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 700;
}

.footer-secondary {
  background: var(--visago-surface-soft);
  color: var(--visago-text-muted);
}

.footer-primary {
  color: #fff;
  background: #0f73ed;
  box-shadow: 0 8px 18px rgba(15, 115, 237, 0.22);
}

.footer-primary--wide {
  grid-column: 1 / -1;
}

.footer-icon {
  font-size: 22px;
}

.report-scroll {
  position: fixed;
  top: var(--ai-topbar-height);
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
  padding-top: 18px;
  padding-bottom: calc(24px + var(--visago-safe-bottom));
}

.score-panel {
  border-radius: 14px;
  border: 1px solid var(--visago-line);
  background: linear-gradient(145deg, var(--visago-surface) 0%, color-mix(in srgb, var(--visago-primary) 10%, var(--visago-surface)) 100%);
  padding: 26px 18px;
  text-align: center;
}

.score-gauge {
  width: 168px;
  height: 168px;
  border-radius: 50%;
  padding: 14px;
  margin: 0 auto;
  box-sizing: border-box;
}

.score-inner {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: color-mix(in srgb, var(--visago-primary) 10%, var(--visago-surface));
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.score-number {
  font-size: 34px;
  font-weight: 800;
  color: #0d5fcb;
}

.score-label {
  margin-top: 2px;
  font-size: 10px;
  color: #7a8496;
  text-transform: uppercase;
}

.score-title {
  display: block;
  margin-top: 22px;
  font-size: 24px;
  font-weight: 800;
  color: var(--visago-text);
}

.score-desc {
  display: block;
  margin-top: 8px;
  font-size: 14px;
  line-height: 1.5;
  color: var(--visago-text-muted);
}

.report-section {
  margin-top: 24px;
}

.section-title {
  display: block;
  margin-bottom: 12px;
  font-size: 21px;
  font-weight: 800;
  color: var(--visago-text);
}

.dimension-card,
.recommend-card,
.improve-card {
  margin-top: 12px;
  border-radius: 12px;
  border: 1px solid var(--visago-line);
  background: var(--visago-surface);
  padding: 14px;
}

.dimension-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.dimension-title,
.risk-title,
.improve-title,
.recommend-country {
  font-size: 15px;
  font-weight: 800;
}

.dimension-tag,
.recommend-score {
  border-radius: 999px;
  padding: 4px 9px;
  background: #e7f8ed;
  color: #087b34;
  font-size: 11px;
  font-weight: 700;
}

.dimension-tag--mid {
  background: #eef0ff;
  color: #4f46e5;
}

.dimension-tag--weak {
  background: #fff1f1;
  color: #d82727;
}

.dimension-track {
  margin-top: 10px;
  height: 7px;
  border-radius: 999px;
  background: #edf0f6;
  overflow: hidden;
}

.dimension-fill {
  height: 100%;
  border-radius: inherit;
  background: #087b34;
}

.dimension-fill--mid {
  background: #4f46e5;
}

.dimension-fill--weak {
  background: #d82727;
}

.dimension-note,
.risk-desc,
.improve-desc,
.recommend-reason {
  display: block;
  margin-top: 8px;
  font-size: 13px;
  line-height: 1.45;
  color: var(--visago-text-muted);
}

.risk-box {
  border-radius: 12px;
  border: 1px solid #ffc9c9;
  background: #fff5f5;
  overflow: hidden;
}

.risk-item {
  padding: 14px;
  display: flex;
  gap: 12px;
  border-bottom: 1px solid #ffdada;
}

.risk-item:last-child {
  border-bottom: none;
}

.risk-icon {
  color: #d82727;
  font-size: 22px;
}

.improve-card,
.recommend-card {
  display: flex;
  align-items: center;
  gap: 14px;
}

.improve-icon {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: #eaf3ff;
  color: #0f73ed;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.improve-copy,
.recommend-card > view {
  flex: 1;
}

.improve-arrow {
  color: #7a8496;
}

.report-actions {
  margin-top: 28px;
  text-align: center;
}

.start-btn {
  height: 54px;
  border-radius: 10px;
  color: #fff;
  background: #0b64c7;
  box-shadow: 0 8px 18px rgba(11, 100, 199, 0.24);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 800;
}

.powered {
  display: block;
  margin-top: 28px;
  font-size: 12px;
  color: #7a8496;
}

.generated {
  display: block;
  margin-top: 4px;
  font-size: 11px;
  color: #b0b7c5;
}
</style>
