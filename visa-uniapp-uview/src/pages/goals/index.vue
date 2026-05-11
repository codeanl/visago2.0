<template>
  <view class="app-screen goals-screen">
    <UPageHeader title="目标" />
    <view class="page-top-spacer" :style="topSpacerStyle" />

    <view class="app-page-width goals-body">
      <swiper
        v-if="planCards.length"
        class="hero-swiper"
        :current="activePlanIndex"
        circular
        :indicator-dots="false"
        @change="onPlanChange"
      >
        <swiper-item v-for="plan in planCards" :key="plan.id">
          <view class="app-card hero-card">
            <view class="hero-top">
              <view class="hero-country">
                <view class="flag-box">{{ plan.countryFlag || '🌍' }}</view>
                <view class="hero-main">
                  <view class="hero-title-row">
                    <text class="hero-title">{{ plan.countryName }}</text>
                    <view class="delete-btn" @tap.stop="confirmDeletePlan(plan.id)">
                      <u-icon name="trash" size="14" color="#ee0a24" />
                    </view>
                  </view>
                  <text class="hero-sub">{{ plan.visaTitle }} · 申请编号 {{ formatPlanCode(plan.id) }}</text>
                </view>
              </view>

              <view class="progress-badge">
                <text class="progress-text">{{ plan.progress }}%</text>
              </view>
            </view>

            <view class="step-track">
              <view class="step-line" />
              <view class="step-line-progress" :style="progressLineStyle(displaySteps)" />
              <view
                v-for="step in displaySteps"
                :key="step.stepKey"
                class="step-item"
                @tap="onMajorStepTap(step.stepKey)"
              >
                <view class="step-dot" :class="`step-dot--${step.status}`">
                  <u-icon v-if="step.status === 'done'" name="checkmark" size="10" color="#ffffff" />
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

      <view v-if="loadingPlans" class="app-card empty-card">
        <u-loading-icon mode="circle" size="22" color="#1677ff" />
        <text class="empty-title">正在加载签证计划</text>
        <text class="empty-sub">请稍等，正在同步你的目标进度。</text>
      </view>

      <view v-else-if="!planCards.length" class="app-card empty-card">
        <text class="empty-title">暂无签证计划</text>
        <text class="empty-sub">先从首页选择一个国家或签证，加入到你的目标里。</text>
      </view>

      <template v-if="planCards.length">
        <view class="section-head">
          <view>
            <text class="app-section-title section-title--sm">进度任务</text>
            <text class="section-sub">点击任务查看攻略，右侧勾选可直接更新状态</text>
          </view>
          <u-tag :text="`${currentStepTasks.length} 项`" type="primary" size="mini" />
        </view>

        <view class="task-list">
          <view
            v-for="task in currentStepTasks"
            :key="task.id"
            class="app-card task-card"
            :class="{ 'task-card--locked': taskLocked }"
          >
            <view class="task-left" @tap="openTaskDrawer(task.id)">
              <view class="task-icon" :class="`task-icon--${task.status}`">
                <u-icon :name="task.icon" size="18" :color="task.status === 'done' ? '#18b566' : task.status === 'review' ? '#1677ff' : '#4e5969'" />
              </view>
              <view class="task-main">
                <text class="task-title">{{ task.title }}</text>
                <text class="task-status" :class="`task-status--${task.status}`">{{ task.statusText }}</text>
              </view>
            </view>
            <view
              class="task-check"
              :class="{ 'task-check--done': task.status === 'done', 'task-check--locked': taskLocked }"
              @tap.stop="toggleTaskChecked(task)"
            >
              <u-icon v-if="task.status === 'done'" name="checkmark" size="14" color="#ffffff" />
            </view>
          </view>
        </view>

        <view v-if="showResultSection" class="result-section">
          <text class="app-section-title section-title--sm">申请结果</text>
          <view class="app-card result-card">
            <view class="result-head">
              <view>
                <text class="result-title">{{ currentResultMeta.title }}</text>
                <text class="result-desc">{{ currentResultMeta.desc }}</text>
              </view>
              <u-tag :text="currentResultMeta.label" :type="resultTagType(currentResultMeta.status)" size="mini" />
            </view>
            <u-button
              v-if="canChooseResult"
              text="选择结果"
              type="primary"
              shape="circle"
              size="small"
              customStyle="margin-top: 12px; width: 128px; align-self: flex-end;"
              @click="choosePlanResult"
            />
          </view>
        </view>

        <view class="tips-section">
          <text class="app-section-title section-title--sm">办理提示</text>
          <view class="app-card tips-card">
            <view v-for="tip in currentPlan.tips || []" :key="tip" class="tip-row">
              <u-icon name="checkmark-circle" size="16" color="#1677ff" />
              <text class="tip-text">{{ tip }}</text>
            </view>
          </view>
        </view>
      </template>
    </view>

    <view v-if="drawerVisible" class="drawer-mask" @tap="closeTaskDrawer">
      <view class="drawer-sheet" @tap.stop>
        <view class="drawer-handle" />
        <view v-if="selectedTask" class="drawer-body">
          <view class="drawer-head">
            <view class="drawer-icon">
              <u-icon :name="selectedTask.icon" size="18" color="#1677ff" />
            </view>
            <view class="drawer-title-wrap">
              <text class="drawer-title">{{ selectedTask.title }}</text>
              <text class="drawer-sub">{{ selectedTaskGuide.summary }}</text>
            </view>
            <view class="drawer-close" @tap="closeTaskDrawer">
              <u-icon name="close" size="16" color="#86909c" />
            </view>
          </view>

          <view class="drawer-status-row">
            <view
              class="task-check"
              :class="{ 'task-check--done': selectedTask.status === 'done', 'task-check--locked': taskLocked }"
              @tap="toggleTaskChecked(selectedTask)"
            >
              <u-icon v-if="selectedTask.status === 'done'" name="checkmark" size="14" color="#ffffff" />
            </view>
            <view>
              <text class="drawer-status-title">{{ selectedTask.status === 'done' ? '已完成' : '未完成' }}</text>
              <text class="drawer-status-sub">{{ taskLocked ? '结果已记录，当前任务状态已锁定' : '点击勾选框即可切换任务状态' }}</text>
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

          <scroll-view scroll-y class="drawer-scroll">
            <view v-if="activeDrawerTab === 'strategy'" class="drawer-block">
              <text class="drawer-block-title">办理攻略</text>
              <view v-for="item in selectedTaskGuide.steps" :key="item" class="drawer-line">
                <u-icon name="checkmark-circle" size="16" color="#1677ff" />
                <text class="drawer-line-text">{{ item }}</text>
              </view>
              <view class="drawer-note">
                <u-icon name="info-circle" size="16" color="#1677ff" />
                <text class="drawer-line-text">{{ selectedTaskGuide.note }}</text>
              </view>
            </view>

            <view v-if="activeDrawerTab === 'guide'" class="drawer-block">
              <text class="drawer-block-title">签证指南</text>
              <view v-for="(step, idx) in currentGuideSteps" :key="`${step.title}-${idx}`" class="guide-card">
                <view class="guide-top">
                  <view class="guide-index">{{ idx + 1 }}</view>
                  <text class="guide-title">{{ step.title }}</text>
                </view>
                <text class="guide-desc">{{ step.desc }}</text>
                <image v-if="step.image" class="guide-image" :src="step.image" mode="aspectFill" />
                <view v-if="step.cta || step.url" class="guide-link" @tap="openGuideLink(step)">
                  <text class="guide-link-text">{{ step.cta || '打开官方链接' }}</text>
                  <u-icon name="arrow-right" size="14" color="#ffffff" />
                </view>
              </view>
            </view>

            <view v-if="activeDrawerTab === 'materials'" class="drawer-block">
              <text class="drawer-block-title">材料要求</text>
              <view v-for="item in selectedTaskGuide.materials" :key="item" class="drawer-line">
                <u-icon name="file-text" size="16" color="#1677ff" />
                <text class="drawer-line-text">{{ item }}</text>
              </view>
            </view>
          </scroll-view>
        </view>
      </view>
    </view>

    <UMainTabbar v-if="!drawerVisible" active="goals" />
  </view>
</template>

<script>
import UPageHeader from '../../components/shared/UPageHeader.vue'
import UMainTabbar from '../../components/shared/UMainTabbar.vue'
import { api } from '../../utils/api'
import { buildGoalView, loadGoals, removeGoal, syncGoalVisaDetail, toggleGoalTask, updateGoalResult } from '../../utils/goals'
import { getTopSpacerStyle } from '../../utils/layout'

const FALLBACK_STEPS = [
  { stepKey: 'apply', title: '申请', status: 'active' },
  { stepKey: 'docs', title: '材料', status: 'todo' },
  { stepKey: 'book', title: '预约', status: 'todo' },
  { stepKey: 'result', title: '结果', status: 'todo' },
]

export default {
  components: { UPageHeader, UMainTabbar },
  data() {
    return {
      loadingPlans: false,
      updatingTask: false,
      updatingResult: false,
      queryPlanId: 0,
      activePlanIndex: 0,
      activeMajorStepKey: 'apply',
      drawerVisible: false,
      selectedTaskId: '',
      activeDrawerTab: 'strategy',
      drawerTabs: [
        { key: 'strategy', label: '办理攻略' },
        { key: 'guide', label: '签证指南' },
        { key: 'materials', label: '材料要求' },
      ],
      planCards: [],
      planDetailMap: {},
      topSpacerStyle: 'height: 108px;',
    }
  },
  computed: {
    currentPlanSummary() {
      return this.planCards[this.activePlanIndex] || null
    },
    currentPlan() {
      if (!this.currentPlanSummary) {
        return { id: 0, resultStatus: 'pending', steps: [], tips: [] }
      }
      return this.planDetailMap[this.currentPlanSummary.id] || { ...this.currentPlanSummary, steps: [], tips: [] }
    },
    displaySteps() {
      const steps = this.currentPlan.steps || []
      return steps.length ? steps : FALLBACK_STEPS
    },
    currentStep() {
      return (
        this.displaySteps.find((step) => step.stepKey === this.activeMajorStepKey) ||
        this.displaySteps[0] || { tasks: [], guides: [], materials: [], strategies: [] }
      )
    },
    currentStepTasks() {
      return this.currentStep.tasks || []
    },
    currentGuideSteps() {
      return this.currentStep.guides || []
    },
    taskLocked() {
      return (this.currentPlan.resultStatus || 'pending') !== 'pending'
    },
    selectedTask() {
      return this.currentStepTasks.find((task) => String(task.id) === String(this.selectedTaskId))
    },
    selectedTaskGuide() {
      const strategies = this.currentStep.strategies || []
      const materials = this.currentStep.materials || []
      const firstGuide = (this.currentStep.guides || [])[0]
      return {
        summary: firstGuide?.desc || this.currentStep.summary || `当前步骤：${this.currentStep.title || '未命名步骤'}`,
        steps: strategies.length ? strategies : ['先查看当前签证详情，再安排这一阶段的事项。'],
        materials,
        note: this.currentStep.note || firstGuide?.cta || '完成当前任务后，计划进度会同步更新。',
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
          status: 'success',
          label: '已完成',
          title: '申请结果已记录',
          desc: this.currentPlan.resultAt ? `结果时间：${this.formatDateTime(this.currentPlan.resultAt)}` : '这次申请已记录为完成。',
        }
      }
      if (status === 'rejected') {
        return {
          status: 'error',
          label: '已放弃',
          title: '申请结果已记录',
          desc: this.currentPlan.resultNote || '这次申请已记录为放弃，可随时继续调整目标。',
        }
      }
      if (status === 'withdrawn') {
        return {
          status: 'warning',
          label: '已取消',
          title: '申请结果已记录',
          desc: this.currentPlan.resultNote || '这次申请已记录为取消。',
        }
      }
      return {
        status: 'primary',
        label: '待选择',
        title: '任务已完成',
        desc: '当前计划进度已到 100%，可以记录这次目标的最终结果。',
      }
    },
  },
  onLoad(query) {
    this.syncLayout()
    this.queryPlanId = Number(query.planId || 0)
  },
  onShow() {
    this.syncLayout()
    this.loadPlans()
  },
  methods: {
    syncLayout() {
      this.topSpacerStyle = getTopSpacerStyle('home')
    },
    resultTagType(status) {
      if (status === 'success') return 'success'
      if (status === 'error') return 'error'
      if (status === 'warning') return 'warning'
      return 'primary'
    },
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
    progressLineStyle(steps) {
      if (!steps.length || steps.length <= 1) return { transform: 'scaleX(0)' }
      const foundIndex = steps.findIndex((step) => step.status === 'active')
      const activeIndex = foundIndex >= 0 ? foundIndex : steps.every((step) => step.status === 'done') ? steps.length - 1 : 0
      const ratio = Math.max(0, Math.min(1, activeIndex / (steps.length - 1)))
      return { transform: `scaleX(${ratio.toFixed(4)})` }
    },
    async loadPlans() {
      this.loadingPlans = true
      try {
        const currentPlanId = this.queryPlanId || this.currentPlanSummary?.id
        const rawGoals = loadGoals()
        const missingVisaSteps = rawGoals.filter((item) => !Array.isArray(item.visaSteps) || !item.visaSteps.length)

        if (missingVisaSteps.length) {
          await Promise.all(
            missingVisaSteps.map(async (item) => {
              try {
                const detail = await api.getVisaDetail(item.visaId)
                syncGoalVisaDetail(item.visaId, detail)
              } catch (error) {
                return null
              }
              return null
            }),
          )
        }

        const plans = loadGoals().map((item) => buildGoalView(item))
        this.planCards = plans
        this.planDetailMap = Object.fromEntries(plans.map((item) => [item.id, item]))

        if (!plans.length) {
          this.activePlanIndex = 0
          this.activeMajorStepKey = 'apply'
          this.closeTaskDrawer()
          return
        }

        const nextIndex = currentPlanId ? plans.findIndex((item) => String(item.id) === String(currentPlanId)) : 0
        this.activePlanIndex = nextIndex >= 0 ? nextIndex : Math.min(this.activePlanIndex, plans.length - 1)
        this.queryPlanId = 0

        const summary = this.currentPlanSummary
        this.activeMajorStepKey = summary?.activeStepKey || summary?.steps?.[0]?.stepKey || 'apply'

        if (!this.currentStepTasks.some((task) => String(task.id) === String(this.selectedTaskId))) {
          this.closeTaskDrawer()
        }
      } catch (error) {
        uni.showToast({ title: (error && error.message) || '加载失败', icon: 'none' })
      } finally {
        this.loadingPlans = false
      }
    },
    onPlanChange(event) {
      this.activePlanIndex = event.detail.current || 0
      this.closeTaskDrawer()
      const summary = this.currentPlanSummary
      this.activeMajorStepKey = summary?.activeStepKey || summary?.steps?.[0]?.stepKey || 'apply'
    },
    confirmDeletePlan(planId) {
      const target = this.planCards.find((plan) => String(plan.id) === String(planId))
      if (!target) return
      uni.showModal({
        title: '确认删除该计划？',
        content: `你正在删除《${target.visaTitle}》。删除后该计划会从当前列表移除，请再次确认。`,
        cancelText: '再想想',
        confirmText: '确认删除',
        confirmColor: '#ee0a24',
        success: (res) => {
          if (!res.confirm) return
          this.deletePlan(planId)
        },
      })
    },
    deletePlan(planId) {
      try {
        removeGoal(planId)
        this.closeTaskDrawer()
        void this.loadPlans()
        uni.showToast({ title: '计划已删除', icon: 'none' })
      } catch (error) {
        uni.showToast({ title: (error && error.message) || '删除失败', icon: 'none' })
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
      this.selectedTaskId = ''
    },
    toggleTaskChecked(task) {
      if (!task || !this.currentPlan.id || this.updatingTask || this.taskLocked) return
      this.updatingTask = true
      try {
        toggleGoalTask(this.currentPlan.id, task.id)
        void this.loadPlans()
      } catch (error) {
        uni.showToast({ title: (error && error.message) || '更新失败', icon: 'none' })
      } finally {
        this.updatingTask = false
      }
    },
    choosePlanResult() {
      if (!this.currentPlan.id || this.updatingResult || Number(this.currentPlan.progress || 0) < 100) return
      uni.showActionSheet({
        itemList: ['已完成', '已放弃', '已取消', '暂不记录'],
        success: (res) => {
          const mapping = ['approved', 'rejected', 'withdrawn', 'pending']
          const resultStatus = mapping[res.tapIndex]
          if (!resultStatus) return
          this.updateCurrentPlanResult(resultStatus)
        },
      })
    },
    updateCurrentPlanResult(resultStatus) {
      if (!this.currentPlan.id) return
      this.updatingResult = true
      try {
        updateGoalResult(this.currentPlan.id, resultStatus, '')
        void this.loadPlans()
        uni.showToast({ title: '已记录结果', icon: 'none' })
      } catch (error) {
        uni.showToast({ title: (error && error.message) || '结果更新失败', icon: 'none' })
      } finally {
        this.updatingResult = false
      }
    },
    openGuideLink(step) {
      const url = String((step && step.url) || '').trim()
      if (!url) return
      if (url.startsWith('/pages/')) {
        uni.navigateTo({ url })
        return
      }
      if (typeof window !== 'undefined' && typeof window.open === 'function') {
        window.open(url, '_blank')
        return
      }
      uni.setClipboardData({
        data: url,
        success: () => uni.showToast({ title: '链接已复制', icon: 'none' }),
      })
    },
  },
}
</script>

<style scoped>
.goals-screen {
  background: var(--uview-bg);
}

.page-top-spacer {
  height: 108px;
}

.goals-body {
  padding: 0 16px calc(116px + var(--uview-safe-bottom));
}

.goals-body {
  padding-top: 2px;
}

.hero-swiper {
  height: 182px;
}

.hero-card {
  padding: 16px;
}

.section-head,
.tips-section,
.result-section {
  position: relative;
  z-index: 1;
}

.hero-top {
  display: flex;
  justify-content: space-between;
}

.hero-country {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
}

.flag-box {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  background: rgba(22, 119, 255, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  flex-shrink: 0;
}

.hero-main {
  flex: 1;
  min-width: 0;
  margin-left: 12px;
}

.hero-title-row {
  display: flex;
  align-items: center;
}

.hero-title {
  font-size: 18px;
  font-weight: 800;
}

.delete-btn {
  width: 24px;
  height: 24px;
  border-radius: 999px;
  background: rgba(238, 10, 36, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 8px;
}

.hero-sub {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: var(--uview-text-3);
}

.progress-badge {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  padding: 5px;
  border: 2px solid rgba(22, 119, 255, 0.18);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 10px;
  flex-shrink: 0;
}

.progress-text {
  font-size: 12px;
  font-weight: 800;
}

.step-track {
  position: relative;
  margin-top: 18px;
  display: flex;
}

.step-line,
.step-line-progress {
  position: absolute;
  left: 12.5%;
  right: 12.5%;
  top: 9px;
  height: 2px;
  border-radius: 999px;
}

.step-line {
  background: var(--uview-border);
}

.step-line-progress {
  background: var(--uview-brand);
  transform-origin: left center;
  transform: scaleX(0);
  transition: transform 0.25s ease;
}

.step-item {
  flex: 1;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.step-dot {
  width: 18px;
  height: 18px;
  border-radius: 999px;
  border: 2px solid var(--uview-border);
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.step-dot--done {
  border-color: var(--uview-brand);
  background: var(--uview-brand);
}

.step-dot--active {
  border-color: var(--uview-brand);
}

.step-dot-inner {
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: var(--uview-brand);
}

.step-label {
  margin-top: 8px;
  font-size: 12px;
  color: var(--uview-text-3);
  font-weight: 600;
}

.step-label--active {
  color: var(--uview-brand);
}

.hero-dots {
  margin: 10px 0 0;
  display: flex;
  justify-content: center;
}

.hero-dot {
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: var(--uview-border);
  margin: 0 3px;
}

.hero-dot--active {
  width: 16px;
  background: var(--uview-brand);
}

.section-head {
  margin: 16px 0 10px;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.section-title--sm {
  font-size: 18px;
}

.section-sub {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: var(--uview-text-3);
}

.task-list {
  display: block;
}

.task-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  margin-bottom: 10px;
  border-radius: 18px;
}

.task-card--locked {
  opacity: 0.88;
}

.task-left {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
}

.task-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(22, 119, 255, 0.08);
  margin-right: 10px;
  flex-shrink: 0;
}

.task-icon--done {
  background: rgba(24, 181, 102, 0.12);
}

.task-icon--review {
  background: rgba(22, 119, 255, 0.12);
}

.task-icon--todo,
.task-icon--missing {
  background: rgba(238, 10, 36, 0.08);
}

.task-main {
  flex: 1;
  min-width: 0;
}

.task-title {
  display: block;
  font-size: 15px;
  font-weight: 700;
}

.task-status {
  display: block;
  margin-top: 4px;
  font-size: 13px;
  font-weight: 600;
}

.task-status--done { color: var(--uview-success); }
.task-status--review { color: var(--uview-brand); }
.task-status--todo,
.task-status--missing { color: var(--uview-danger); }

.task-check {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  border: 2px solid var(--uview-border);
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.task-check--done {
  border-color: var(--uview-success);
  background: var(--uview-success);
}

.task-check--locked {
  opacity: 0.62;
}

.empty-card {
  padding: 22px 18px;
  text-align: center;
}

.empty-title {
  display: block;
  margin-top: 10px;
  font-size: 16px;
  font-weight: 700;
}

.empty-sub {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  color: var(--uview-text-3);
  line-height: 1.5;
}

.result-section,
.tips-section {
  margin-top: 10px;
}

.result-card,
.tips-card {
  padding: 14px;
  border-radius: 18px;
}

.result-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
}

.result-title {
  display: block;
  font-size: 15px;
  font-weight: 800;
}

.result-desc {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: var(--uview-text-3);
  line-height: 1.5;
}

.tip-row {
  display: flex;
  align-items: flex-start;
  margin-bottom: 10px;
}

.tip-row :deep(.u-icon) {
  margin-right: 8px;
  flex-shrink: 0;
}

.tip-text {
  flex: 1;
  font-size: 13px;
  line-height: 1.5;
  color: var(--uview-text-3);
}

.drawer-mask {
  position: fixed;
  inset: 0;
  z-index: 120;
  background: rgba(18, 24, 38, 0.38);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.drawer-sheet {
  width: 100%;
  max-width: 430px;
  height: 72vh;
  min-height: 440px;
  max-height: 680px;
  border-radius: 22px 22px 0 0;
  background: #fff;
  padding: 10px 16px 18px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.drawer-handle {
  width: 42px;
  height: 4px;
  border-radius: 999px;
  background: var(--uview-border);
  margin: 0 auto 16px;
}

.drawer-body {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.drawer-head {
  display: flex;
  align-items: center;
}

.drawer-icon {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  background: rgba(22, 119, 255, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  flex-shrink: 0;
}

.drawer-title-wrap {
  flex: 1;
  min-width: 0;
}

.drawer-title {
  display: block;
  font-size: 17px;
  font-weight: 800;
}

.drawer-sub {
  display: block;
  margin-top: 4px;
  font-size: 11px;
  line-height: 1.45;
  color: var(--uview-text-3);
}

.drawer-close {
  width: 34px;
  height: 34px;
  border-radius: 999px;
  background: var(--uview-surface-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.drawer-status-row {
  margin-top: 16px;
  padding: 12px;
  border-radius: 16px;
  background: rgba(22, 119, 255, 0.08);
  border: 1px solid rgba(22, 119, 255, 0.14);
  display: flex;
  align-items: center;
}

.drawer-status-row .task-check {
  margin-right: 10px;
}

.drawer-status-title {
  display: block;
  font-size: 14px;
  font-weight: 800;
}

.drawer-status-sub {
  display: block;
  margin-top: 2px;
  font-size: 11px;
  color: var(--uview-text-3);
}

.drawer-tabs {
  margin-top: 14px;
  padding: 4px;
  border-radius: 14px;
  background: var(--uview-surface-soft);
  display: flex;
}

.drawer-tab {
  flex: 1;
  height: 36px;
  border-radius: 11px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--uview-text-3);
  font-size: 13px;
  font-weight: 800;
}

.drawer-tab + .drawer-tab {
  margin-left: 4px;
}

.drawer-tab--active {
  color: var(--uview-brand);
  background: #fff;
  box-shadow: 0 8px 24px rgba(31, 35, 41, 0.05);
}

.drawer-scroll {
  flex: 1;
  min-height: 0;
  margin-top: 12px;
}

.drawer-block {
  padding-bottom: 16px;
}

.drawer-block-title {
  display: block;
  margin-bottom: 10px;
  font-size: 14px;
  font-weight: 800;
}

.drawer-line,
.drawer-note,
.guide-card {
  border-radius: 16px;
  background: var(--uview-surface-soft);
}

.drawer-line {
  margin-top: 9px;
  padding: 11px 12px;
  display: flex;
  align-items: flex-start;
}

.drawer-line :deep(.u-icon),
.drawer-note :deep(.u-icon) {
  margin-right: 8px;
  flex-shrink: 0;
}

.drawer-line-text {
  flex: 1;
  font-size: 13px;
  line-height: 1.45;
  color: var(--uview-text-3);
}

.drawer-note {
  margin-top: 12px;
  padding: 12px;
  display: flex;
  align-items: flex-start;
}

.guide-card {
  margin-top: 10px;
  padding: 14px;
}

.guide-top {
  display: flex;
  align-items: center;
}

.guide-index {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  background: rgba(22, 119, 255, 0.14);
  color: var(--uview-brand);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 800;
  margin-right: 10px;
}

.guide-title {
  font-size: 16px;
  font-weight: 800;
}

.guide-desc {
  display: block;
  margin-top: 10px;
  font-size: 13px;
  line-height: 1.6;
  color: var(--uview-text-3);
}

.guide-image {
  margin-top: 12px;
  width: 100%;
  height: 140px;
  border-radius: 10px;
}

.guide-link {
  margin-top: 14px;
  height: 42px;
  border-radius: 12px;
  background: var(--uview-brand);
  display: flex;
  align-items: center;
  justify-content: center;
}

.guide-link-text {
  color: #fff;
  font-size: 14px;
  font-weight: 700;
  margin-right: 6px;
}
</style>
