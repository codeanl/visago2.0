<template>
  <view class="page-screen detail-screen">
    <VuPageHeader :title="detail.name || '签证详情'" :show-back="true" />
    <view class="page-top-spacer" :style="topSpacerStyle" />

    <view class="page-width page-body detail-body">
      <view v-if="loading" class="app-card app-empty">
        <u-loading-icon mode="circle" size="24" color="#1677ff" />
        <text class="app-empty-title">正在加载详情</text>
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

          <view class="meta-row">
            <view class="meta-item">
              <u-icon name="clock" size="16" color="#1677ff" />
              <text class="meta-label">办理时长</text>
              <text class="meta-value">{{ detail.processingTime || '待更新' }}</text>
            </view>
            <view class="meta-item">
              <u-icon name="rmb" size="16" color="#1677ff" />
              <text class="meta-label">费用参考</text>
              <text class="meta-value">{{ detail.fee || '待更新' }}</text>
            </view>
            <view class="meta-item">
              <u-icon name="bookmark" size="16" color="#1677ff" />
              <text class="meta-label">有效期</text>
              <text class="meta-value">{{ detail.validity || '待更新' }}</text>
            </view>
            <view class="meta-item">
              <u-icon name="calendar" size="16" color="#1677ff" />
              <text class="meta-label">入境次数</text>
              <text class="meta-value">{{ detail.entries || '待更新' }}</text>
            </view>
          </view>

          <view class="button-row">
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
          <text class="app-section-title section-title-sm">签证说明</text>
          <text class="prose-body">{{ detail.longIntro || detail.description || '暂无说明' }}</text>
        </view>

        <view class="app-card prose-card">
          <text class="app-section-title section-title-sm">办理步骤</text>
          <view v-if="detail.steps && detail.steps.length" class="step-list">
            <view v-for="(step, index) in detail.steps" :key="step.id || step.stepKey || index" class="step-card">
              <view class="step-top">
                <view class="step-index">{{ index + 1 }}</view>
                <text class="step-title">{{ step.title }}</text>
              </view>
              <view v-if="step.materials && step.materials.length" class="materials-row">
                <view v-for="material in step.materials" :key="material" class="material-tag-cell">
                  <u-tag
                    :text="material"
                    type="info"
                    size="mini"
                    plain
                    plain-fill
                  />
                </view>
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
import VuPageHeader from '../../components/uview/VuPageHeader.vue'
import { api } from '../../utils/api'
import { getGoalByVisaId, removeGoal, upsertGoal } from '../../utils/goals'
import { getTopSpacerStyle } from '../../utils/layout'

export default {
  components: {
    VuPageHeader,
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
      topSpacerStyle: 'height:96px;',
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
  background: #f7f8fb;
}

.page-top-spacer {
  height: 96px;
}

.detail-body {
  padding-bottom: 32px;
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
  margin-right: 10px;
}

.detail-name {
  display: block;
  color: #1f2329;
  font-size: 24px;
  line-height: 1.25;
  font-weight: 900;
}

.detail-type {
  display: block;
  margin-top: 8px;
  color: #1677ff;
  font-size: 13px;
  line-height: 1.3;
}

.detail-desc,
.prose-body,
.step-text {
  display: block;
  color: #6b7280;
  font-size: 13px;
  line-height: 1.75;
}

.detail-desc {
  margin-top: 12px;
}

.meta-row {
  margin-top: 14px;
  display: flex;
  flex-wrap: wrap;
  margin-left: -5px;
  margin-right: -5px;
}

.meta-item {
  width: calc(50% - 10px);
  margin-left: 5px;
  margin-right: 5px;
  margin-bottom: 10px;
  padding: 12px;
  border-radius: 16px;
  background: #eef6ff;
}

.meta-label {
  display: block;
  margin-top: 8px;
  color: #7b8495;
  font-size: 11px;
}

.meta-value {
  display: block;
  margin-top: 6px;
  color: #1f2329;
  font-size: 14px;
  line-height: 1.35;
  font-weight: 900;
}

.button-row {
  margin-top: 8px;
}

.prose-card {
  margin-top: 14px;
}

.section-title-sm {
  font-size: 18px;
}

.prose-body {
  margin-top: 12px;
}

.step-list {
  margin-top: 12px;
}

.step-card {
  margin-bottom: 12px;
  padding: 14px;
  border-radius: 18px;
  background: #f4f7fb;
}

.step-top {
  display: flex;
  align-items: center;
}

.step-index {
  width: 28px;
  height: 28px;
  margin-right: 10px;
  border-radius: 14px;
  background: #1677ff;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 900;
  flex-shrink: 0;
}

.step-title {
  color: #1f2329;
  font-size: 15px;
  font-weight: 900;
}

.materials-row {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
}

.material-tag-cell {
  margin-right: 6px;
  margin-bottom: 6px;
}

.step-text {
  margin-top: 8px;
}
</style>
