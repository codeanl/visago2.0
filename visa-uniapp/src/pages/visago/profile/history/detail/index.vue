<template>
  <view class="history-detail-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="申请详情" />
    <scroll-view scroll-y class="detail-scroll">
      <view class="detail-content visago-page-width">
        <view v-if="loading" class="empty-card visago-card">
          <text class="material-symbols-outlined">hourglass_top</text>
          <text>正在加载申请详情</text>
        </view>

        <template v-else-if="detail.id">
          <view class="summary-card visago-card">
            <view class="summary-head">
              <view>
                <text class="summary-country">{{ detail.countryName || '-' }}</text>
                <text class="summary-visa">{{ detail.visaTitle || '-' }}</text>
              </view>
              <text class="summary-badge" :class="`summary-badge--${resultMeta.status}`">{{ resultMeta.label }}</text>
            </view>
            <view class="summary-grid">
              <view class="summary-item">
                <text class="summary-label">申请编号</text>
                <text class="summary-value">#{{ String(detail.id).padStart(4, '0') }}</text>
              </view>
              <view class="summary-item">
                <text class="summary-label">创建时间</text>
                <text class="summary-value">{{ formatDate(detail.createdAt) || '-' }}</text>
              </view>
              <view class="summary-item">
                <text class="summary-label">结果时间</text>
                <text class="summary-value">{{ formatDate(detail.resultAt) || '-' }}</text>
              </view>
              <view class="summary-item">
                <text class="summary-label">完成进度</text>
                <text class="summary-value">{{ detail.progress || 0 }}%</text>
              </view>
            </view>
            <text class="summary-desc">{{ resultMeta.desc }}</text>
          </view>

          <view class="steps-section">
            <text class="steps-title">当次申请步骤快照</text>
            <view class="steps-list">
              <view v-for="step in detail.steps || []" :key="step.stepKey" class="step-card visago-card">
                <view class="step-head" @tap="toggleStep(step.stepKey)">
                  <view class="step-title-wrap">
                    <text class="step-title">{{ step.title || '未命名步骤' }}</text>
                    <text class="step-key">{{ step.stepKey || '-' }}</text>
                  </view>
                  <view class="step-head-right">
                    <text class="step-status" :class="`step-status--${step.status}`">{{ formatStepStatus(step.status) }}</text>
                    <text class="material-symbols-outlined step-arrow">{{ isStepExpanded(step.stepKey) ? 'expand_less' : 'expand_more' }}</text>
                  </view>
                </view>

                <view v-if="isStepExpanded(step.stepKey)" class="step-body">
                  <view v-if="step.strategies && step.strategies.length" class="detail-block">
                    <text class="detail-block__title">办理攻略</text>
                    <view v-for="item in step.strategies" :key="item" class="detail-bullet">
                      <text class="material-symbols-outlined detail-bullet__icon">check_circle</text>
                      <text>{{ item }}</text>
                    </view>
                  </view>

                  <view v-if="step.guides && step.guides.length" class="detail-block">
                    <text class="detail-block__title">签证指南</text>
                    <view v-for="(guide, idx) in step.guides" :key="`${guide.title}-${idx}`" class="guide-card">
                      <text class="guide-title">{{ guide.title || '未命名指南' }}</text>
                      <text class="guide-desc">{{ guide.desc || '-' }}</text>
                      <image v-if="guide.image" class="guide-image" :src="guide.image" mode="aspectFill" />
                      <text v-if="guide.url" class="guide-url">{{ guide.url }}</text>
                    </view>
                  </view>

                  <view v-if="step.materials && step.materials.length" class="detail-block">
                    <text class="detail-block__title">材料要求</text>
                    <view v-for="item in step.materials" :key="item" class="detail-bullet detail-bullet--muted">
                      <text class="material-symbols-outlined detail-bullet__icon">description</text>
                      <text>{{ item }}</text>
                    </view>
                  </view>

                  <view v-if="step.tasks && step.tasks.length" class="detail-block">
                    <text class="detail-block__title">任务清单</text>
                    <view v-for="task in step.tasks" :key="task.id || task.taskKey" class="task-row">
                      <view class="task-row__left">
                        <text class="material-symbols-outlined task-row__icon">{{ task.icon || 'task_alt' }}</text>
                        <view>
                          <text class="task-row__title">{{ task.title || '未命名任务' }}</text>
                          <text class="task-row__key">{{ task.taskKey || '-' }}</text>
                        </view>
                      </view>
                      <text class="task-row__status" :class="`task-row__status--${task.status}`">{{ task.statusText || formatTaskStatus(task.status) }}</text>
                    </view>
                  </view>
                </view>
              </view>
            </view>
          </view>
        </template>

        <view v-else class="empty-card visago-card">
          <text class="material-symbols-outlined">assignment</text>
          <text>没有找到这条申请记录</text>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../../components/VisagoTopBar.vue'
import { api } from '../../../../../utils/api'
import { applyTheme, getStoredTheme } from '../../../../../utils/theme'

export default {
  components: { VisagoTopBar },
  data() {
    return {
      loading: false,
      planId: 0,
      detail: {
        id: 0,
        countryName: '',
        visaTitle: '',
        progress: 0,
        resultStatus: 'pending',
        resultNote: '',
        resultAt: '',
        createdAt: '',
        steps: [],
      },
      expandedStepKeys: [],
    }
  },
  computed: {
    resultMeta() {
      const status = this.detail.resultStatus || 'pending'
      if (status === 'approved') {
        return {
          status: 'approved',
          label: '已出签',
          desc: '这条记录对应的是当时那次申请的完整快照，后台模板后续修改不会影响这里。',
        }
      }
      if (status === 'rejected') {
        return {
          status: 'rejected',
          label: '被拒签',
          desc: this.detail.resultNote || '这条记录保留了当时提交时的步骤和任务状态，方便复盘拒签原因。',
        }
      }
      if (status === 'withdrawn') {
        return {
          status: 'withdrawn',
          label: '已撤签',
          desc: this.detail.resultNote || '这次申请已撤签，这里保留的是当次申请时的快照内容。',
        }
      }
      return {
        status: 'pending',
        label: '待结果',
        desc: '当前申请尚未记录最终结果。',
      }
    },
  },
  onLoad(query) {
    this.planId = Number(query.planId || 0)
  },
  onShow() {
    applyTheme(getStoredTheme())
    this.loadDetail()
  },
  methods: {
    async loadDetail() {
      if (!this.planId) return
      this.loading = true
      try {
        const detail = await api.getPlanDetail(this.planId)
        this.detail = detail || this.detail
        this.expandedStepKeys = detail?.steps?.length ? [detail.steps[0].stepKey] : []
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '申请详情加载失败',
          icon: 'none',
        })
      } finally {
        this.loading = false
      }
    },
    formatDate(value) {
      if (!value) return ''
      const date = new Date(value)
      if (Number.isNaN(date.getTime())) return String(value)
      const y = date.getFullYear()
      const m = String(date.getMonth() + 1).padStart(2, '0')
      const d = String(date.getDate()).padStart(2, '0')
      return `${y}-${m}-${d}`
    },
    formatStepStatus(status) {
      if (status === 'done') return '已完成'
      if (status === 'active') return '进行中'
      return '待处理'
    },
    formatTaskStatus(status) {
      if (status === 'done') return '已完成'
      if (status === 'review') return '审核中'
      if (status === 'missing') return '有缺失'
      return '待处理'
    },
    isStepExpanded(stepKey) {
      return this.expandedStepKeys.includes(stepKey)
    },
    toggleStep(stepKey) {
      if (this.isStepExpanded(stepKey)) {
        this.expandedStepKeys = this.expandedStepKeys.filter((item) => item !== stepKey)
      } else {
        this.expandedStepKeys = [...this.expandedStepKeys, stepKey]
      }
    },
  },
}
</script>

<style scoped>
.history-detail-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.detail-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.detail-content {
  padding: 18px 16px 34px;
  box-sizing: border-box;
}

.empty-card {
  padding: 22px;
  color: var(--visago-text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 700;
}

.summary-card {
  padding: 16px;
}

.summary-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.summary-country,
.summary-visa,
.summary-label,
.summary-value,
.summary-desc,
.steps-title,
.step-title,
.step-key,
.detail-block__title,
.guide-title,
.guide-desc,
.guide-url,
.task-row__title,
.task-row__key,
.task-row__status {
  display: block;
}

.summary-country {
  font-size: 18px;
  font-weight: 900;
}

.summary-visa {
  margin-top: 4px;
  font-size: 13px;
  color: var(--visago-text-muted);
}

.summary-badge {
  padding: 5px 10px;
  border-radius: 999px;
  background: rgba(15, 101, 216, 0.12);
  color: var(--visago-primary);
  font-size: 12px;
  font-weight: 800;
}

.summary-badge--approved {
  color: #16a34a;
  background: rgba(22, 163, 74, 0.12);
}

.summary-badge--rejected {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.12);
}

.summary-badge--withdrawn {
  color: #475569;
  background: rgba(100, 116, 139, 0.14);
}

.summary-grid {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.summary-item {
  padding: 10px 12px;
  border-radius: 12px;
  background: var(--visago-surface-soft);
}

.summary-label {
  font-size: 12px;
  color: var(--visago-text-soft);
}

.summary-value {
  margin-top: 4px;
  font-size: 14px;
  font-weight: 700;
}

.summary-desc {
  margin-top: 14px;
  font-size: 13px;
  line-height: 1.6;
  color: var(--visago-text-muted);
}

.steps-section {
  margin-top: 14px;
}

.steps-title {
  font-size: 15px;
  font-weight: 800;
  color: var(--visago-text);
}

.steps-list {
  margin-top: 10px;
  display: grid;
  gap: 12px;
}

.step-card {
  padding: 14px;
}

.step-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.step-title {
  font-size: 15px;
  font-weight: 800;
}

.step-key {
  margin-top: 4px;
  font-size: 12px;
  color: var(--visago-text-soft);
}

.step-head-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.step-status {
  padding: 4px 8px;
  border-radius: 999px;
  background: var(--visago-surface-soft);
  color: var(--visago-text-muted);
  font-size: 11px;
  font-weight: 800;
}

.step-status--done {
  color: #16a34a;
  background: rgba(22, 163, 74, 0.12);
}

.step-status--active {
  color: var(--visago-primary);
  background: rgba(15, 101, 216, 0.12);
}

.step-arrow {
  font-size: 18px;
  color: var(--visago-text-soft);
}

.step-body {
  margin-top: 14px;
  display: grid;
  gap: 14px;
}

.detail-block__title {
  margin-bottom: 8px;
  font-size: 13px;
  font-weight: 800;
}

.detail-bullet {
  margin-top: 8px;
  padding: 10px 12px;
  border-radius: 12px;
  background: rgba(15, 101, 216, 0.08);
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 13px;
  line-height: 1.5;
  color: var(--visago-text-muted);
}

.detail-bullet--muted {
  background: var(--visago-surface-soft);
}

.detail-bullet__icon {
  color: var(--visago-primary);
  font-size: 18px;
  flex-shrink: 0;
}

.guide-card {
  margin-top: 10px;
  padding: 12px;
  border-radius: 14px;
  background: var(--visago-surface-soft);
}

.guide-title {
  font-size: 14px;
  font-weight: 700;
}

.guide-desc {
  margin-top: 8px;
  font-size: 13px;
  line-height: 1.6;
  color: var(--visago-text-muted);
}

.guide-image {
  margin-top: 10px;
  width: 100%;
  height: 150px;
  border-radius: 10px;
}

.guide-url {
  margin-top: 10px;
  font-size: 11px;
  line-height: 1.5;
  color: var(--visago-primary);
  word-break: break-all;
}

.task-row {
  margin-top: 9px;
  padding: 11px 12px;
  border-radius: 12px;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.task-row__left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  flex: 1;
}

.task-row__icon {
  font-size: 20px;
  color: var(--visago-primary);
}

.task-row__title {
  font-size: 14px;
  font-weight: 700;
}

.task-row__key {
  margin-top: 4px;
  font-size: 12px;
  color: var(--visago-text-soft);
}

.task-row__status {
  font-size: 12px;
  font-weight: 800;
  color: var(--visago-text-muted);
  flex-shrink: 0;
}

.task-row__status--done {
  color: #16a34a;
}

.task-row__status--review {
  color: var(--visago-primary);
}

.task-row__status--missing,
.task-row__status--todo {
  color: #ef4444;
}
</style>
