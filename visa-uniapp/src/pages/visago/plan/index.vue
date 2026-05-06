<template>
  <view class="plan-page" :class="{ 'plan-page--dark': themeMode === 'dark' }">
    <VisagoTopBar page-name="计划" />

    <view class="plan-content visago-page-width">
      <swiper
        v-if="planCards.length"
        class="hero-swiper"
        :current="activePlanIndex"
        circular
        :indicator-dots="false"
        @change="onPlanChange"
      >
        <swiper-item v-for="plan in planCards" :key="plan.id">
          <view class="hero-card visago-card">
            <view class="hero-top">
              <view class="visa-info">
                <view class="flag-box">
                  <text class="flag-text">{{ plan.countryFlag || '🌍' }}</text>
                </view>
                <view>
                  <view class="plan-top">
                    <text class="visa-title">{{ plan.countryName }}</text>
                    <view class="plan-delete-inline" @tap.stop="confirmDeletePlan(plan.id)">
                      <text class="material-symbols-outlined">delete</text>
                    </view>
                  </view>
                  <text class="visa-sub">{{ plan.visaTitle }} · 申请编号 {{ formatPlanCode(plan.id) }}</text>
                </view>
              </view>

              <view class="ring" :style="ringStyle(plan.progress)">
                <view class="ring-inner">{{ plan.progress }}%</view>
              </view>
            </view>

            <view
              class="steps-track"
              :style="{ gridTemplateColumns: `repeat(${Math.max(displaySteps.length, 1)}, minmax(0, 1fr))` }"
            >
              <view class="steps-line" />
              <view class="steps-line-progress" :style="{ width: progressLineWidth(displaySteps) }" />

              <view v-for="step in displaySteps" :key="step.stepKey" class="step-item" @tap="onMajorStepTap(step.stepKey)">
                <view class="step-dot" :class="`step-dot--${step.status}`">
                  <text v-if="step.status === 'done'" class="material-symbols-outlined done-icon">check</text>
                  <view v-else-if="step.status === 'active'" class="step-dot-inner" />
                </view>
                <text class="step-label" :class="{ 'step-label--active': step.status === 'active' }">{{ step.title }}</text>
              </view>
            </view>
          </view>
        </swiper-item>
      </swiper>

      <view v-if="planCards.length" class="hero-dots">
        <view
          v-for="(item, idx) in planCards"
          :key="item.id"
          class="hero-dot"
          :class="{ 'hero-dot--active': idx === activePlanIndex }"
        />
      </view>

      <view v-if="loadingPlans" class="empty-plan visago-card">
        <text class="material-symbols-outlined">hourglass_top</text>
        <text>正在加载签证计划</text>
      </view>

      <view v-else-if="!planCards.length" class="empty-plan visago-card">
        <text class="material-symbols-outlined">assignment</text>
        <text>暂无签证计划</text>
      </view>

      <template v-if="planCards.length">
        <view class="task-section-head">
          <view>
            <text class="task-section-title">进度任务</text>
            <text class="task-section-sub">点击任务查看攻略，右侧勾选可直接更新状态</text>
          </view>
          <text class="task-section-count">{{ currentStepTasks.length }} 项</text>
        </view>

        <view class="task-list">
          <view
            v-for="task in currentStepTasks"
            :key="task.id"
            class="task-card visago-card"
            :class="{ 'task-card--missing': task.status === 'missing' }"
          >
            <view class="task-left" @tap="openTaskDrawer(task.id)">
              <view class="task-icon-wrap" :class="`task-icon-wrap--${task.status}`">
                <text class="material-symbols-outlined">{{ task.icon }}</text>
              </view>
              <view>
                <text class="task-title">{{ task.title }}</text>
                <text class="task-status" :class="`task-status--${task.status}`">{{ task.statusText }}</text>
              </view>
            </view>

            <view class="task-checkbox" :class="{ 'task-checkbox--checked': task.status === 'done' }" @tap.stop="toggleTaskChecked(task)">
              <text v-if="task.status === 'done'" class="material-symbols-outlined">check</text>
            </view>
          </view>

          <view v-if="!currentStepTasks.length" class="empty-tasks visago-card">
            当前阶段暂无任务
          </view>
        </view>

        <view v-if="showResultSection" class="result-section">
          <text class="tips-title">申请结果</text>
          <view class="result-card visago-card">
            <view class="result-head">
              <view>
                <text class="result-title">{{ currentResultMeta.title }}</text>
                <text class="result-desc">{{ currentResultMeta.desc }}</text>
              </view>
              <text class="result-badge" :class="`result-badge--${currentResultMeta.status}`">{{ currentResultMeta.label }}</text>
            </view>
            <view v-if="canChooseResult" class="result-actions">
              <view class="ghost-btn" @tap="choosePlanResult">{{ currentPlan.resultStatus === 'pending' ? '选择结果' : '修改结果' }}</view>
            </view>
          </view>
        </view>

        <view class="tips-section">
          <text class="tips-title">办理提示</text>
          <view class="tips-card visago-card">
            <view v-for="tip in currentPlan.tips || []" :key="tip" class="tip-item">
              <text class="material-symbols-outlined tip-icon">check_circle</text>
              <text class="tip-text">{{ tip }}</text>
            </view>
          </view>
        </view>
      </template>
    </view>

    <view v-if="drawerVisible" class="drawer-mask" @tap="closeTaskDrawer">
      <view class="task-drawer" @tap.stop>
        <view class="drawer-handle" />
        <view v-if="selectedTask" class="drawer-body">
          <view class="drawer-fixed-head">
            <view class="drawer-head">
              <view class="drawer-icon" :class="`drawer-icon--${selectedTask.status}`">
                <text class="material-symbols-outlined">{{ selectedTask.icon }}</text>
              </view>
              <view class="drawer-title-wrap">
                <text class="drawer-title">{{ selectedTask.title }}</text>
                <text class="drawer-sub">{{ selectedTaskGuide.summary }}</text>
              </view>
              <view class="drawer-close" @tap="closeTaskDrawer">
                <text class="material-symbols-outlined">close</text>
              </view>
            </view>

            <view class="drawer-status-row">
              <view class="drawer-check" :class="{ 'drawer-check--checked': selectedTask.status === 'done' }" @tap="toggleTaskChecked(selectedTask)">
                <text v-if="selectedTask.status === 'done'" class="material-symbols-outlined">check</text>
              </view>
              <view>
                <text class="drawer-status-title">{{ selectedTask.status === 'done' ? '已完成' : '未完成' }}</text>
                <text class="drawer-status-sub">点击勾选框即可切换任务状态</text>
              </view>
            </view>

            <view class="drawer-tabs">
              <view
                v-for="tab in drawerTabs"
                :key="tab.key"
                class="drawer-tab"
                :class="{ 'drawer-tab--active': activeDrawerTab === tab.key }"
                @tap="activeDrawerTab = tab.key"
              >
                {{ tab.label }}
              </view>
            </view>
          </view>

          <scroll-view scroll-y class="drawer-scroll">
            <view v-if="activeDrawerTab === 'strategy'" class="drawer-section">
              <text class="drawer-section-title">办理攻略</text>
              <view v-for="item in selectedTaskGuide.steps" :key="item" class="drawer-guide-row">
                <text class="material-symbols-outlined">check_circle</text>
                <text>{{ item }}</text>
              </view>
              <view class="drawer-note">
                <text class="material-symbols-outlined">tips_and_updates</text>
                <text>{{ selectedTaskGuide.note }}</text>
              </view>
            </view>

            <view v-if="activeDrawerTab === 'guide'" class="drawer-section">
              <text class="drawer-section-title">签证指南</text>
              <view v-for="(step, idx) in currentGuideSteps" :key="`${step.title}-${idx}`" class="guide-detail-card">
                <view class="guide-step-head">
                  <view class="guide-index">{{ idx + 1 }}</view>
                  <text class="guide-step-title">{{ step.title }}</text>
                </view>
                <text class="guide-step-desc">{{ step.desc }}</text>
                <image v-if="step.image" class="guide-step-image" :src="step.image" mode="aspectFill" />
                <text v-if="step.url" class="guide-url">{{ step.url }}</text>
                <view v-if="step.cta || step.url" class="guide-cta" @tap="openGuideLink(step)">
                  <text class="guide-cta-text">{{ step.cta || '打开官方链接' }}</text>
                  <text class="material-symbols-outlined guide-cta-icon">open_in_new</text>
                </view>
              </view>
              <view v-if="!currentGuideSteps.length" class="drawer-note">
                <text class="material-symbols-outlined">tips_and_updates</text>
                <text>当前步骤还没有配置签证指南。</text>
              </view>
            </view>

            <view v-if="activeDrawerTab === 'materials'" class="drawer-section">
              <text class="drawer-section-title">材料要求</text>
              <view v-for="item in selectedTaskGuide.materials" :key="item" class="material-row">
                <text class="material-symbols-outlined">description</text>
                <text>{{ item }}</text>
              </view>
              <view v-if="!selectedTaskGuide.materials.length" class="drawer-note">
                <text class="material-symbols-outlined">tips_and_updates</text>
                <text>当前步骤还没有配置材料要求。</text>
              </view>
            </view>
          </scroll-view>
        </view>
      </view>
    </view>

    <VisagoBottomNav active-tab="plan" />
  </view>
</template>

<script>
import VisagoTopBar from '../../../components/VisagoTopBar.vue'
import VisagoBottomNav from '../../../components/VisagoBottomNav.vue'
import { api } from '../../../utils/api'
import { applyTheme, getStoredTheme } from '../../../utils/theme'

const FALLBACK_STEPS = [
  { stepKey: 'apply', title: '申请', status: 'done' },
  { stepKey: 'docs', title: '材料', status: 'active' },
  { stepKey: 'book', title: '预约', status: 'todo' },
  { stepKey: 'result', title: '结果', status: 'todo' },
]

export default {
  components: {
    VisagoTopBar,
    VisagoBottomNav,
  },
  data() {
    return {
      themeMode: 'light',
      loadingPlans: false,
      updatingTask: false,
      updatingResult: false,
      queryPlanId: 0,
      historyPlanOnly: false,
      activePlanIndex: 0,
      activeMajorStepKey: 'docs',
      drawerVisible: false,
      selectedTaskId: 0,
      activeDrawerTab: 'strategy',
      drawerTabs: [
        { key: 'strategy', label: '办理攻略' },
        { key: 'guide', label: '签证指南' },
        { key: 'materials', label: '材料要求' },
      ],
      planCards: [],
      planDetailMap: {},
    }
  },
  computed: {
    currentPlanSummary() {
      return this.planCards[this.activePlanIndex] || null
    },
    currentPlan() {
      if (!this.currentPlanSummary) {
        return { id: 0, steps: [], tips: [] }
      }
      return this.planDetailMap[this.currentPlanSummary.id] || { ...this.currentPlanSummary, steps: [], tips: [] }
    },
    displaySteps() {
      const steps = this.currentPlan.steps || []
      return steps.length ? steps : FALLBACK_STEPS
    },
    currentStep() {
      return this.displaySteps.find((step) => step.stepKey === this.activeMajorStepKey) || this.displaySteps[0] || { tasks: [], guides: [], materials: [], strategies: [] }
    },
    currentStepTasks() {
      return this.currentStep.tasks || []
    },
    currentGuideSteps() {
      return this.currentStep.guides || []
    },
    selectedTask() {
      return this.currentStepTasks.find((task) => Number(task.id) === Number(this.selectedTaskId))
    },
    selectedTaskGuide() {
      const strategies = this.currentStep.strategies || []
      const materials = this.currentStep.materials || []
      const firstGuide = (this.currentStep.guides || [])[0]
      return {
        summary: firstGuide && firstGuide.desc ? firstGuide.desc : `当前步骤：${this.currentStep.title || '未命名步骤'}`,
        steps: strategies,
        materials,
        note: firstGuide && firstGuide.cta ? firstGuide.cta : '更新任务状态后，后台管理端与数据库中的计划进度会同步刷新。',
      }
    },
    showResultSection() {
      return this.planCards.length && Number(this.currentPlan.progress || 0) >= 100
    },
    canChooseResult() {
      return this.showResultSection && (this.currentPlan.resultStatus || 'pending') === 'pending'
    },
    currentResultMeta() {
      const status = this.currentPlan.resultStatus || 'pending'
      if (status === 'approved') {
        return {
          status: 'approved',
          label: '已出签',
          title: '申请结果已记录',
          desc: this.currentPlan.resultAt ? `结果时间：${this.formatDateTime(this.currentPlan.resultAt)}` : '这次申请已记录为出签。',
        }
      }
      if (status === 'rejected') {
        return {
          status: 'rejected',
          label: '被拒签',
          title: '申请结果已记录',
          desc: this.currentPlan.resultNote || '这次申请已记录为拒签，可在历史中随时回看当时步骤。',
        }
      }
      if (status === 'withdrawn') {
        return {
          status: 'withdrawn',
          label: '已撤签',
          title: '申请结果已记录',
          desc: this.currentPlan.resultNote || '这次申请已记录为撤签。',
        }
      }
      return {
        status: 'pending',
        label: '待选择',
        title: '任务已完成',
        desc: '当前计划进度已到 100%，可以记录这次 DIY 申请的最终结果。',
      }
    },
  },
  onLoad(query) {
    this.queryPlanId = Number(query.planId || 0)
  },
  onShow() {
    this.themeMode = applyTheme(getStoredTheme())
    this.loadPlans()
  },
  methods: {
    formatPlanCode(planId) {
      return `#${String(planId).padStart(4, '0')}`
    },
    formatDateTime(value) {
      if (!value) return ''
      const date = new Date(value)
      if (Number.isNaN(date.getTime())) return String(value)
      const y = date.getFullYear()
      const m = String(date.getMonth() + 1).padStart(2, '0')
      const d = String(date.getDate()).padStart(2, '0')
      return `${y}-${m}-${d}`
    },
    ringStyle(percent) {
      const safe = Math.max(0, Math.min(100, Number(percent) || 0))
      const deg = Math.round((360 * safe) / 100)
      return {
        background: `conic-gradient(var(--visago-primary) 0deg ${deg}deg, var(--visago-surface-soft) ${deg}deg 360deg)`,
      }
    },
    progressLineWidth(steps) {
      if (!steps.length || steps.length <= 1) return '0%'
      const activeIndex = Math.max(0, steps.findIndex((step) => step.status === 'active'))
      const ratio = activeIndex / (steps.length - 1)
      return `calc((100% - 36px) * ${ratio})`
    },
    async loadPlans() {
      this.loadingPlans = true
      try {
        const plans = (await api.listPlans()) || []
        const activePlans = plans.filter((item) => (item.resultStatus || 'pending') === 'pending')
        const targetId = this.queryPlanId || this.currentPlanSummary?.id

        if (targetId) {
          const targetPlan = plans.find((item) => Number(item.id) === Number(targetId))
          if (targetPlan && (targetPlan.resultStatus || 'pending') !== 'pending') {
            this.historyPlanOnly = true
            this.planCards = [targetPlan]
            this.activePlanIndex = 0
            await this.loadPlanDetail(targetPlan.id)
            this.queryPlanId = 0
            return
          }
        }

        this.historyPlanOnly = false
        this.planCards = activePlans
        if (!activePlans.length) {
          this.planDetailMap = {}
          this.activePlanIndex = 0
          this.activeMajorStepKey = 'docs'
          return
        }

        const nextIndex = targetId ? activePlans.findIndex((item) => Number(item.id) === Number(targetId)) : 0
        this.activePlanIndex = nextIndex >= 0 ? nextIndex : 0
        this.queryPlanId = 0
        await this.loadPlanDetail(this.planCards[this.activePlanIndex].id)
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '加载失败',
          icon: 'none',
        })
      } finally {
        this.loadingPlans = false
      }
    },
    async loadPlanDetail(planId) {
      const detail = await api.getPlanDetail(planId)
      this.planDetailMap = {
        ...this.planDetailMap,
        [planId]: detail,
      }
      const index = this.planCards.findIndex((item) => Number(item.id) === Number(planId))
      if (index >= 0) {
        this.planCards.splice(index, 1, {
          ...this.planCards[index],
          progress: detail.progress,
          activeStepKey: detail.activeStepKey,
          countryFlag: detail.countryFlag,
        })
      }
      this.activeMajorStepKey = detail.activeStepKey || detail.steps?.[0]?.stepKey || 'docs'
      if (!this.currentStepTasks.some((task) => Number(task.id) === Number(this.selectedTaskId))) {
        this.selectedTaskId = 0
      }
    },
    async onPlanChange(event) {
      this.activePlanIndex = event.detail.current || 0
      this.closeTaskDrawer()
      const summary = this.currentPlanSummary
      if (!summary) return
      if (!this.planDetailMap[summary.id]) {
        await this.loadPlanDetail(summary.id)
      } else {
        this.activeMajorStepKey = this.planDetailMap[summary.id].activeStepKey || 'docs'
      }
    },
    async confirmDeletePlan(planId) {
      const target = this.planCards.find((plan) => Number(plan.id) === Number(planId))
      if (!target) return
      uni.showModal({
        title: '确认删除该计划？',
        content: `你正在删除「${target.visaTitle}」。删除后该计划会从当前列表移除，请再次确认。`,
        cancelText: '再想想',
        confirmText: '确认删除',
        confirmColor: '#ef4444',
        success: async (res) => {
          if (!res.confirm) return
          await this.deletePlan(planId)
        },
      })
    },
    async deletePlan(planId) {
      try {
        await api.deletePlan(planId)
        this.planCards = this.planCards.filter((plan) => Number(plan.id) !== Number(planId))
        delete this.planDetailMap[planId]
        this.closeTaskDrawer()
        if (this.planCards.length) {
          this.activePlanIndex = Math.min(this.activePlanIndex, this.planCards.length - 1)
          await this.loadPlanDetail(this.planCards[this.activePlanIndex].id)
        } else {
          this.activePlanIndex = 0
          this.activeMajorStepKey = 'docs'
        }
        uni.showToast({ title: '计划已删除', icon: 'none' })
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '删除失败',
          icon: 'none',
        })
      }
    },
    onMajorStepTap(stepKey) {
      this.activeMajorStepKey = stepKey
      this.closeTaskDrawer()
    },
    openTaskDrawer(taskId) {
      this.selectedTaskId = taskId
      this.activeDrawerTab = 'strategy'
      this.drawerVisible = true
    },
    closeTaskDrawer() {
      this.drawerVisible = false
      this.selectedTaskId = 0
    },
    async toggleTaskChecked(task) {
      if (!task || !this.currentPlan.id || this.updatingTask) return
      this.updatingTask = true
      try {
        const nextDone = task.status !== 'done'
        const detail = await api.updatePlanTask(this.currentPlan.id, task.id, {
          status: nextDone ? 'done' : 'todo',
          statusText: nextDone ? '已完成' : '待处理',
        })
        this.planDetailMap = {
          ...this.planDetailMap,
          [detail.id]: detail,
        }
        const index = this.planCards.findIndex((item) => Number(item.id) === Number(detail.id))
        if (index >= 0) {
          this.planCards.splice(index, 1, {
            ...this.planCards[index],
            progress: detail.progress,
            activeStepKey: detail.activeStepKey,
          })
        }
        this.activeMajorStepKey = detail.activeStepKey || this.activeMajorStepKey
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '更新失败',
          icon: 'none',
        })
      } finally {
        this.updatingTask = false
      }
    },
    choosePlanResult() {
      if (!this.currentPlan.id || this.updatingResult || Number(this.currentPlan.progress || 0) < 100) {
        return
      }
      const itemList = ['已出签', '被拒签', '已撤签', '暂不记录']
      uni.showActionSheet({
        itemList,
        success: async (res) => {
          const mapping = ['approved', 'rejected', 'withdrawn', 'pending']
          const resultStatus = mapping[res.tapIndex]
          if (!resultStatus) return
          await this.updateCurrentPlanResult(resultStatus)
        },
      })
    },
    async updateCurrentPlanResult(resultStatus) {
      if (!this.currentPlan.id) return
      this.updatingResult = true
      try {
        const detail = await api.updatePlanResult(this.currentPlan.id, {
          resultStatus,
          resultNote: '',
        })
        this.planDetailMap = {
          ...this.planDetailMap,
          [detail.id]: detail,
        }
        const index = this.planCards.findIndex((item) => Number(item.id) === Number(detail.id))
        if (index >= 0) {
          this.planCards.splice(index, 1, {
            ...this.planCards[index],
            progress: detail.progress,
            activeStepKey: detail.activeStepKey,
            status: detail.status,
            resultStatus: detail.resultStatus,
            resultNote: detail.resultNote,
            resultAt: detail.resultAt,
          })
        }
        if (resultStatus !== 'pending' && !this.historyPlanOnly) {
          this.planCards = this.planCards.filter((item) => Number(item.id) !== Number(detail.id))
          delete this.planDetailMap[detail.id]
          this.closeTaskDrawer()
          if (this.planCards.length) {
            this.activePlanIndex = Math.min(this.activePlanIndex, this.planCards.length - 1)
            await this.loadPlanDetail(this.planCards[this.activePlanIndex].id)
          } else {
            this.activePlanIndex = 0
            this.activeMajorStepKey = 'docs'
          }
        }
        uni.showToast({
          title:
            resultStatus === 'approved'
              ? '已移至申请历史'
              : resultStatus === 'rejected'
                ? '已移至申请历史'
                : resultStatus === 'withdrawn'
                  ? '已移至申请历史'
                  : '已恢复为待结果',
          icon: 'none',
        })
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '结果更新失败',
          icon: 'none',
        })
      } finally {
        this.updatingResult = false
      }
    },
    openGuideLink(step) {
      const url = String((step && step.url) || '').trim()
      if (!url) return
      if (typeof window !== 'undefined' && typeof window.open === 'function') {
        window.open(url, '_blank')
        return
      }
      uni.setClipboardData({
        data: url,
        success: () => {
          uni.showToast({
            title: '链接已复制',
            icon: 'none',
          })
        },
      })
    },
  },
}
</script>

<style scoped>
.plan-page {
  --plan-success: #16a34a;
  --plan-warning: #0f65d8;
  --plan-danger: #ef4444;
  min-height: 100vh;
  background: var(--visago-bg);
}

.plan-content {
  box-sizing: border-box;
  padding: 96px 16px calc(118px + var(--visago-safe-bottom));
}

.hero-swiper {
  height: 182px;
}

.hero-card {
  position: relative;
  padding: 16px;
  margin: 0 1px;
  box-sizing: border-box;
}

.hero-top {
  display: flex;
  justify-content: space-between;
  gap: 10px;
}

.visa-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.flag-box {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  background: #f5f7fc;
  border: 1px solid var(--visago-line);
  display: flex;
  align-items: center;
  justify-content: center;
}

.flag-text {
  font-size: 28px;
}

.plan-top {
  display: flex;
  align-items: center;
  gap: 8px;
}

.visa-title {
  display: block;
  font-size: 18px;
  font-weight: 700;
  color: var(--visago-text);
}

.visa-sub {
  font-size: 12px;
  color: var(--visago-text-muted);
}

.ring {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  padding: 5px;
  box-sizing: border-box;
}

.ring-inner {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: var(--visago-surface);
  color: var(--visago-text);
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.plan-delete-inline {
  width: 24px;
  height: 24px;
  padding: 0;
  border-radius: 999px;
  background: rgba(239, 68, 68, 0.1);
  color: var(--plan-danger);
  border: 1px solid rgba(239, 68, 68, 0.16);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.plan-delete-inline .material-symbols-outlined {
  font-size: 14px;
}

.steps-track {
  position: relative;
  margin-top: 18px;
  display: grid;
  gap: 0;
}

.steps-line,
.steps-line-progress {
  position: absolute;
  left: 18px;
  right: 18px;
  top: 9px;
  height: 2px;
  border-radius: 999px;
}

.steps-line {
  background: var(--visago-line);
}

.steps-line-progress {
  right: auto;
  background: var(--visago-primary);
  transition: width 0.25s ease;
}

.step-item {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.step-dot {
  width: 18px;
  height: 18px;
  border-radius: 999px;
  border: 2px solid var(--visago-line);
  background: var(--visago-surface);
  display: flex;
  align-items: center;
  justify-content: center;
}

.step-dot--done {
  border-color: var(--visago-primary);
  background: var(--visago-primary);
}

.step-dot--active {
  border-color: var(--visago-primary);
  background: var(--visago-surface);
}

.done-icon {
  font-size: 12px;
  color: #fff;
}

.step-dot-inner {
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: var(--visago-primary);
}

.step-label {
  font-size: 12px;
  line-height: 1.2;
  color: var(--visago-text-muted);
  font-weight: 600;
}

.step-label--active {
  color: var(--visago-primary);
}

.hero-dots {
  margin: -14px 0 0;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 6px;
}

.hero-dot {
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: var(--visago-line);
}

.hero-dot--active {
  width: 16px;
  background: var(--visago-primary);
}

.task-section-head {
  margin: -2px 0 8px;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 12px;
}

.empty-plan {
  margin: 8px 0 14px;
  padding: 22px;
  color: var(--visago-text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 700;
}

.task-section-title,
.task-section-sub {
  display: block;
}

.task-section-title {
  font-size: 17px;
  font-weight: 800;
  color: var(--visago-text);
}

.task-section-sub {
  margin-top: 3px;
  font-size: 12px;
  color: var(--visago-text-muted);
}

.task-section-count {
  padding: 5px 9px;
  border-radius: 999px;
  background: rgba(15, 101, 216, 0.12);
  color: var(--visago-primary);
  font-size: 12px;
  font-weight: 800;
  flex-shrink: 0;
}

.task-list {
  display: grid;
  gap: 10px;
}

.task-card {
  padding: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.task-card--missing {
  border-style: dashed;
  border-color: color-mix(in srgb, #ff4d4f 30%, var(--visago-line));
}

.task-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.task-icon-wrap {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--visago-surface-soft);
  color: var(--visago-primary);
}

.task-icon-wrap .material-symbols-outlined {
  font-size: 20px;
}

.task-icon-wrap--done {
  background: rgba(22, 163, 74, 0.14);
  color: var(--plan-success);
}

.task-icon-wrap--review {
  background: rgba(15, 101, 216, 0.14);
  color: var(--plan-warning);
}

.task-icon-wrap--missing,
.task-icon-wrap--todo {
  background: rgba(239, 68, 68, 0.13);
  color: var(--plan-danger);
}

.task-title {
  display: block;
  font-size: 15px;
  font-weight: 700;
  color: var(--visago-text);
}

.task-status {
  display: block;
  margin-top: 4px;
  font-size: 13px;
  font-weight: 600;
}

.task-status--done {
  color: var(--plan-success);
}

.task-status--review {
  color: var(--plan-warning);
}

.task-status--missing,
.task-status--todo {
  color: var(--plan-danger);
}

.task-checkbox {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  border: 2px solid var(--visago-line);
  background: var(--visago-surface);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.task-checkbox--checked {
  border-color: var(--plan-success);
  background: var(--plan-success);
}

.task-checkbox .material-symbols-outlined {
  font-size: 16px;
  color: #fff;
}

.empty-tasks {
  text-align: center;
  padding: 14px;
  font-size: 13px;
  color: var(--visago-text-muted);
}

.tips-section {
  margin-top: 14px;
  display: grid;
  gap: 10px;
}

.result-section {
  margin-top: 14px;
  display: grid;
  gap: 10px;
}

.result-card {
  padding: 14px;
  display: grid;
  gap: 12px;
}

.result-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.result-title,
.result-desc {
  display: block;
}

.result-title {
  font-size: 15px;
  font-weight: 800;
  color: var(--visago-text);
}

.result-desc {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.5;
  color: var(--visago-text-muted);
}

.result-badge {
  padding: 5px 10px;
  border-radius: 999px;
  background: rgba(15, 101, 216, 0.12);
  color: var(--visago-primary);
  font-size: 12px;
  font-weight: 800;
  flex-shrink: 0;
}

.result-badge--approved {
  background: rgba(22, 163, 74, 0.12);
  color: #16a34a;
}

.result-badge--rejected {
  background: rgba(239, 68, 68, 0.12);
  color: #ef4444;
}

.result-badge--withdrawn {
  background: rgba(100, 116, 139, 0.14);
  color: #475569;
}

.result-actions {
  display: flex;
  justify-content: flex-end;
}

.tips-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--visago-text);
}

.tips-card {
  padding: 12px;
  display: grid;
  gap: 10px;
}

.tip-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.tip-icon {
  color: var(--visago-primary);
  font-size: 18px;
  margin-top: 1px;
}

.tip-text {
  flex: 1;
  color: var(--visago-text-soft);
  font-size: 13px;
  line-height: 1.5;
}

.drawer-mask {
  position: fixed;
  top: 0;
  right: 0;
  bottom: calc(82px + var(--visago-safe-bottom));
  left: 0;
  z-index: 120;
  background: rgba(0, 0, 0, 0.42);
  display: flex;
  align-items: flex-end;
}

.task-drawer {
  width: 100%;
  height: calc(100vh - 106px - 82px - var(--visago-safe-bottom));
  border-radius: 24px 24px 0 0;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: 0 -18px 38px rgba(15, 23, 42, 0.2);
  padding: 10px 16px 18px;
  box-sizing: border-box;
}

.drawer-handle {
  width: 42px;
  height: 4px;
  border-radius: 999px;
  background: var(--visago-line);
  margin: 0 auto 16px;
}

.drawer-body {
  height: calc(100% - 20px);
  max-width: 430px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.drawer-fixed-head {
  flex-shrink: 0;
}

.drawer-head {
  display: flex;
  align-items: center;
  gap: 12px;
}

.drawer-icon {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  background: var(--visago-surface-soft);
  color: var(--visago-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.drawer-icon--done {
  background: rgba(22, 163, 74, 0.14);
  color: var(--plan-success);
}

.drawer-icon--review {
  background: rgba(15, 101, 216, 0.14);
  color: var(--plan-warning);
}

.drawer-icon--missing,
.drawer-icon--todo {
  background: rgba(239, 68, 68, 0.13);
  color: var(--plan-danger);
}

.drawer-title-wrap {
  flex: 1;
  min-width: 0;
}

.drawer-title,
.drawer-sub,
.drawer-section-title {
  display: block;
}

.drawer-title {
  font-size: 18px;
  font-weight: 800;
  color: var(--visago-text);
}

.drawer-sub {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.45;
  color: var(--visago-text-muted);
}

.drawer-close {
  width: 34px;
  height: 34px;
  border-radius: 999px;
  background: var(--visago-surface-soft);
  color: var(--visago-text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.drawer-close .material-symbols-outlined {
  font-size: 19px;
}

.drawer-status-row {
  margin-top: 16px;
  padding: 12px;
  border-radius: 16px;
  background: color-mix(in srgb, var(--visago-primary) 8%, var(--visago-surface));
  border: 1px solid rgba(15, 101, 216, 0.14);
  display: flex;
  align-items: center;
  gap: 10px;
}

.drawer-check {
  width: 30px;
  height: 30px;
  border-radius: 999px;
  border: 2px solid var(--visago-line);
  background: var(--visago-surface);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.drawer-check--checked {
  border-color: var(--plan-success);
  background: var(--plan-success);
}

.drawer-check .material-symbols-outlined {
  color: #fff;
  font-size: 17px;
}

.drawer-status-title,
.drawer-status-sub {
  display: block;
}

.drawer-status-title {
  font-size: 14px;
  font-weight: 800;
  color: var(--visago-text);
}

.drawer-status-sub {
  margin-top: 2px;
  font-size: 11px;
  color: var(--visago-text-muted);
}

.drawer-tabs {
  margin-top: 14px;
  padding: 4px;
  border-radius: 14px;
  background: var(--visago-surface-soft);
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 4px;
}

.drawer-tab {
  height: 36px;
  border-radius: 11px;
  color: var(--visago-text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 800;
}

.drawer-tab--active {
  color: var(--visago-primary);
  background: var(--visago-surface);
  box-shadow: var(--visago-shadow-card);
}

.drawer-scroll {
  flex: 1;
  min-height: 0;
  margin-top: 12px;
}

.drawer-section {
  padding-bottom: 16px;
}

.drawer-section-title {
  margin-bottom: 10px;
  font-size: 14px;
  font-weight: 800;
  color: var(--visago-text);
}

.drawer-guide-row {
  margin-top: 9px;
  border-radius: 14px;
  padding: 11px 12px;
  background: var(--visago-surface-soft);
  color: var(--visago-text-muted);
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 13px;
  line-height: 1.45;
}

.drawer-guide-row .material-symbols-outlined {
  color: var(--visago-primary);
  font-size: 18px;
  margin-top: 1px;
  flex-shrink: 0;
}

.drawer-note,
.material-row,
.guide-detail-card {
  border-radius: 16px;
  background: var(--visago-surface-soft);
}

.drawer-note {
  margin-top: 12px;
  padding: 12px;
  color: var(--visago-text-muted);
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 12px;
  line-height: 1.5;
}

.drawer-note .material-symbols-outlined {
  color: var(--visago-primary);
  font-size: 19px;
  flex-shrink: 0;
}

.guide-detail-card {
  padding: 14px;
  margin-top: 10px;
}

.guide-step-head {
  display: flex;
  align-items: center;
  gap: 10px;
}

.guide-index {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--visago-primary) 14%, transparent);
  color: var(--visago-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
}

.guide-step-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--visago-text);
}

.guide-step-desc {
  display: block;
  margin-top: 10px;
  color: var(--visago-text-soft);
  line-height: 1.7;
  font-size: 15px;
}

.guide-url {
  display: block;
  margin-top: 10px;
  font-size: 11px;
  line-height: 1.5;
  color: var(--visago-primary);
  word-break: break-all;
}

.guide-step-image {
  margin-top: 12px;
  width: 100%;
  height: 140px;
  border-radius: 10px;
}

.guide-cta {
  margin-top: 14px;
  height: 42px;
  border-radius: 10px;
  background: var(--visago-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.guide-cta-text,
.guide-cta-icon {
  color: #fff;
  font-size: 15px;
  font-weight: 600;
}

.material-row {
  margin-top: 9px;
  padding: 12px;
  color: var(--visago-text-muted);
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 13px;
  line-height: 1.45;
}

.material-row .material-symbols-outlined {
  color: var(--visago-primary);
  font-size: 18px;
  flex-shrink: 0;
}

.plan-page--dark .task-card {
  background: #0f1420;
  border-color: #1e293b;
}

.plan-page--dark .task-title {
  color: #f4f7ff;
}

.plan-page--dark .task-status--done {
  color: #20d39a;
}

.plan-page--dark .task-status--review {
  color: #56a8ff;
}

.plan-page--dark .task-drawer {
  background: #171b24;
  border-color: #30394c;
}

@media (max-width: 360px) {
  .plan-content {
    padding-left: 12px;
    padding-right: 12px;
  }

  .visa-title {
    font-size: 16px;
  }

  .ring {
    width: 62px;
    height: 62px;
  }

  .guide-step-title {
    font-size: 18px;
  }
}
</style>
