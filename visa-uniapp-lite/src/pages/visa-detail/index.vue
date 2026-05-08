<template>
  <view class="screen">
    <lite-header :title="detail.name || '签证详情'" :show-back="true" />

    <view class="page-wrap detail-stack">
      <view v-if="loading" class="card empty-state">
        <text class="empty-state__emoji">签</text>
        <text class="empty-state__title">正在加载详情</text>
      </view>

      <template v-else>
        <view class="card detail-hero">
          <view class="detail-hero__top">
            <view>
              <text class="detail-name">{{ detail.name }}</text>
              <text class="detail-type">{{ detail.visaType || '待补充' }}</text>
            </view>
            <text v-if="savedGoal" class="detail-saved">已加入目标</text>
          </view>
          <text v-if="detail.description" class="detail-desc">{{ detail.description }}</text>
          <view class="meta-grid">
            <view class="meta-item">
              <text class="meta-label">办理时长</text>
              <text class="meta-value">{{ detail.processingTime || '待更新' }}</text>
            </view>
            <view class="meta-item">
              <text class="meta-label">费用参考</text>
              <text class="meta-value">{{ detail.fee || '待更新' }}</text>
            </view>
            <view class="meta-item">
              <text class="meta-label">有效期</text>
              <text class="meta-value">{{ detail.validity || '待更新' }}</text>
            </view>
            <view class="meta-item">
              <text class="meta-label">入境次数</text>
              <text class="meta-value">{{ detail.entries || '待更新' }}</text>
            </view>
          </view>
          <view class="action-row">
            <view v-if="savedGoal" class="button-secondary action-btn" @tap="removeFromGoals">移出目标</view>
            <view v-else class="button-primary action-btn" @tap="saveToGoals">加入我的目标</view>
            <view class="button-secondary action-btn" @tap="openGoals">查看目标页</view>
          </view>
        </view>

        <view class="card prose-card">
          <text class="section-title">签证说明</text>
          <text class="prose-body">{{ detail.longIntro || detail.description || '暂无说明' }}</text>
        </view>

        <view class="card prose-card">
          <text class="section-title">办理步骤</text>
          <view v-if="detail.steps && detail.steps.length" class="step-list">
            <view v-for="(step, index) in detail.steps" :key="step.id || step.stepKey || index" class="step-card">
              <view class="step-top">
                <text class="step-index">{{ index + 1 }}</text>
                <text class="step-title">{{ step.title }}</text>
              </view>
              <view v-if="step.materials && step.materials.length" class="pill-row">
                <text v-for="material in step.materials" :key="material" class="pill pill--muted">{{ material }}</text>
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
import LiteHeader from '../../components/LiteHeader.vue'
import { api } from '../../utils/api'
import { getGoalByVisaId, removeGoal, upsertGoal } from '../../utils/goals'

export default {
  components: {
    LiteHeader,
  },
  data() {
    return {
      visaId: '',
      countryName: '',
      loading: false,
      detail: {
        steps: [],
      },
      savedGoal: null,
    }
  },
  onLoad(query) {
    this.visaId = String(query.visaId || '')
    this.countryName = decodeURIComponent(String(query.countryName || ''))
  },
  onShow() {
    this.savedGoal = getGoalByVisaId(this.visaId)
    this.loadDetail()
  },
  methods: {
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
.detail-stack {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.detail-hero,
.prose-card {
  padding: 18px;
}

.detail-hero__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
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
  color: var(--lite-primary);
}

.detail-saved {
  padding: 8px 10px;
  border-radius: 999px;
  background: rgba(22, 141, 100, 0.12);
  color: var(--lite-success);
  font-size: 12px;
  font-weight: 800;
}

.detail-desc,
.prose-body,
.step-text {
  display: block;
  font-size: 13px;
  line-height: 1.8;
  color: var(--lite-text-muted);
}

.detail-desc {
  margin-top: 12px;
}

.meta-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin-top: 14px;
}

.meta-item {
  padding: 12px;
  border-radius: 16px;
  background: rgba(15, 95, 123, 0.05);
}

.meta-label {
  display: block;
  font-size: 11px;
  color: var(--lite-text-muted);
}

.meta-value {
  display: block;
  margin-top: 6px;
  font-size: 14px;
  font-weight: 800;
}

.action-row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin-top: 16px;
}

.action-btn {
  width: 100%;
}

.prose-body {
  margin-top: 12px;
}

.step-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 12px;
}

.step-card {
  padding: 14px;
  border-radius: 18px;
  background: rgba(15, 95, 123, 0.05);
}

.step-top {
  display: flex;
  align-items: center;
  gap: 10px;
}

.step-index {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  background: var(--lite-primary);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 800;
}

.step-title {
  font-size: 15px;
  font-weight: 800;
}
</style>
