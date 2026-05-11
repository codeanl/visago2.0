<template>
  <view class="app-screen detail-screen">
    <UPageHeader :title="detail.name || '签证详情'" :show-back="true" />
    <view class="page-top-spacer" :style="topSpacerStyle" />

    <view class="app-page-width detail-body">
      <view v-if="loading" class="app-card empty-card">
        <u-loading-icon mode="circle" size="24" color="#1677ff" />
        <text class="empty-title">正在加载详情</text>
      </view>

      <template v-else>
        <view class="app-card detail-card">
          <view class="detail-head">
            <view class="detail-main">
              <text class="detail-name">{{ detail.name }}</text>
              <text class="detail-type">{{ detail.visaType || '待补充' }}</text>
            </view>
            <u-tag v-if="savedGoal" text="已加入目标" type="success" size="mini" />
          </view>

          <text class="detail-desc">{{ detail.description || detail.longIntro || '暂无说明' }}</text>

          <view class="detail-grid">
            <view class="detail-item">
              <text class="detail-label">办理时长</text>
              <text class="detail-value">{{ detail.processingTime || '待更新' }}</text>
            </view>
            <view class="detail-item">
              <text class="detail-label">费用参考</text>
              <text class="detail-value">{{ detail.fee || '待更新' }}</text>
            </view>
            <view class="detail-item">
              <text class="detail-label">有效期</text>
              <text class="detail-value">{{ detail.validity || '待更新' }}</text>
            </view>
            <view class="detail-item">
              <text class="detail-label">入境次数</text>
              <text class="detail-value">{{ detail.entries || '待更新' }}</text>
            </view>
          </view>

          <view class="action-row">
            <u-button
              v-if="savedGoal"
              text="移出目标"
              type="error"
              shape="circle"
              customStyle="margin-bottom: 10px;"
              @click="removeFromGoals"
            />
            <u-button
              v-else
              text="加入我的目标"
              type="primary"
              shape="circle"
              customStyle="margin-bottom: 10px;"
              @click="saveToGoals"
            />
            <u-button text="查看目标页" type="info" shape="circle" @click="openGoals" />
          </view>
        </view>

        <view class="app-card prose-card">
          <text class="app-section-title section-title--sm">签证说明</text>
          <text class="prose-body">{{ detail.longIntro || detail.description || '暂无说明' }}</text>
        </view>

        <view class="app-card prose-card">
          <text class="app-section-title section-title--sm">办理步骤</text>
          <view v-if="detail.steps && detail.steps.length" class="step-list">
            <view v-for="(step, index) in detail.steps" :key="step.id || step.stepKey || index" class="step-card">
              <view class="step-top">
                <view class="step-index">{{ index + 1 }}</view>
                <text class="step-title">{{ step.title }}</text>
              </view>
              <text v-if="step.strategies && step.strategies.length" class="step-text">{{ step.strategies.join('；') }}</text>
            </view>
          </view>
          <text v-else class="prose-body">暂无步骤</text>
        </view>
      </template>
    </view>
  </view>
</template>

<script>
import UPageHeader from '../../components/shared/UPageHeader.vue'
import { api } from '../../utils/api'
import { getGoalByVisaId, removeGoal, upsertGoal } from '../../utils/goals'
import { getTopSpacerStyle } from '../../utils/layout'

export default {
  components: { UPageHeader },
  data() {
    return {
      visaId: '',
      countryName: '',
      loading: false,
      detail: {
        steps: [],
      },
      savedGoal: null,
      topSpacerStyle: 'height: 96px;',
    }
  },
  onLoad(query) {
    this.syncLayout()
    this.visaId = String(query.visaId || '')
    this.countryName = decodeURIComponent(String(query.countryName || ''))
  },
  onShow() {
    this.syncLayout()
    this.savedGoal = getGoalByVisaId(this.visaId)
    this.loadDetail()
  },
  methods: {
    syncLayout() {
      this.topSpacerStyle = getTopSpacerStyle('sub')
    },
    async loadDetail() {
      if (!this.visaId) {
        uni.showToast({ title: '缺少签证信息', icon: 'none' })
        return
      }
      this.loading = true
      try {
        const detail = await api.getVisaDetail(this.visaId)
        this.detail = detail || { steps: [] }
      } catch (error) {
        uni.showToast({ title: (error && error.message) || '加载失败', icon: 'none' })
      } finally {
        this.loading = false
      }
    },
    saveToGoals() {
      upsertGoal({
        visaId: this.detail.id,
        countryId: this.detail.countryId,
        countryName: this.countryName || this.detail.countryName || '',
        countryFlag: this.detail.countryFlag || ((this.countryName || this.detail.countryName || '').slice(0, 1) || '签'),
        visaName: this.detail.name,
        visaType: this.detail.visaType,
        processingTime: this.detail.processingTime,
        fee: this.detail.fee,
        validity: this.detail.validity,
        entries: this.detail.entries,
        description: this.detail.description,
        visaSteps: this.detail.steps || [],
      })
      this.savedGoal = getGoalByVisaId(this.visaId)
      uni.showToast({ title: '已加入目标', icon: 'none' })
    },
    removeFromGoals() {
      removeGoal(this.visaId)
      this.savedGoal = null
      uni.showToast({ title: '已移出目标', icon: 'none' })
    },
    openGoals() {
      uni.reLaunch({ url: '/pages/goals/index' })
    },
  },
}
</script>

<style scoped>
.detail-screen {
  background: var(--uview-bg);
}

.page-top-spacer {
  height: 96px;
}

.detail-body {
  padding: 0 16px 32px;
}

.detail-card,
.prose-card {
  padding: 18px;
}

.detail-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
}

.detail-main {
  flex: 1;
  min-width: 0;
}

.detail-name {
  display: block;
  font-size: 24px;
  line-height: 1.25;
  font-weight: 900;
}

.detail-type {
  display: block;
  margin-top: 8px;
  font-size: 13px;
  color: var(--uview-brand);
}

.detail-desc,
.prose-body,
.step-text {
  display: block;
  font-size: 13px;
  line-height: 1.75;
  color: var(--uview-text-3);
}

.detail-desc {
  margin-top: 12px;
}

.detail-grid {
  display: flex;
  flex-wrap: wrap;
  margin: 14px -5px 0;
}

.detail-item {
  width: calc(50% - 10px);
  margin: 0 5px 10px;
  padding: 12px;
  border-radius: 16px;
  background: rgba(22, 119, 255, 0.05);
}

.detail-label {
  display: block;
  font-size: 11px;
  color: var(--uview-text-3);
}

.detail-value {
  display: block;
  margin-top: 6px;
  font-size: 14px;
  font-weight: 800;
}

.action-row {
  margin-top: 16px;
}

.prose-card {
  margin-top: 14px;
}

.section-title--sm {
  font-size: 18px;
}

.prose-body {
  margin-top: 12px;
}

.step-list {
  display: block;
  margin-top: 12px;
}

.step-card {
  padding: 14px;
  border-radius: 18px;
  background: rgba(22, 119, 255, 0.05);
  margin-bottom: 12px;
}

.step-top {
  display: flex;
  align-items: center;
}

.step-index {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  background: var(--uview-brand);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 800;
  margin-right: 10px;
}

.step-title {
  font-size: 15px;
  font-weight: 800;
}

.empty-card {
  text-align: center;
  padding: 28px 20px;
}

.empty-title {
  display: block;
  margin-top: 14px;
  font-size: 18px;
  font-weight: 800;
}
</style>
