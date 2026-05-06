<template>
  <view class="sub-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="申请历史" />
    <scroll-view scroll-y class="sub-scroll">
      <view class="sub-content visago-page-width">
        <view class="filter-row">
          <view
            v-for="item in filters"
            :key="item"
            class="filter-chip"
            :class="{ 'filter-chip--active': activeFilter === item }"
            @tap="activeFilter = item"
          >
            {{ item }}
          </view>
        </view>

        <view class="timeline">
          <view v-for="item in filteredList" :key="item.id" class="history-card">
            <view class="timeline-dot" :class="`timeline-dot--${item.status}`" />
            <view class="history-main">
              <view class="history-head">
                <text class="history-title">{{ item.title }}</text>
                <text class="history-status" :class="`history-status--${item.status}`">{{ item.statusText }}</text>
              </view>
              <text class="history-date">{{ item.date }}</text>
              <text class="history-desc">{{ item.desc }}</text>
              <view class="history-actions">
                <view class="ghost-btn" @tap="openPlanDetail(item)">查看详情</view>
                <view class="primary-btn" @tap="reapply(item)">再次申请</view>
              </view>
            </view>
          </view>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { api } from '../../../../utils/api'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'

export default {
  components: { VisagoTopBar },
  data() {
    return {
      activeFilter: '全部',
      filters: ['全部', '已出签', '被拒签', '已撤签'],
      list: [],
    }
  },
  computed: {
    filteredList() {
      if (this.activeFilter === '全部') return this.list
      return this.list.filter((item) => item.statusText === this.activeFilter)
    },
  },
  onShow() {
    applyTheme(getStoredTheme())
    this.loadHistory()
  },
  methods: {
    formatDate(value) {
      if (!value) return ''
      const date = new Date(value)
      if (Number.isNaN(date.getTime())) return String(value)
      const y = date.getFullYear()
      const m = String(date.getMonth() + 1).padStart(2, '0')
      const d = String(date.getDate()).padStart(2, '0')
      return `${y}-${m}-${d}`
    },
    normalizeHistoryItem(item) {
      if (item.resultStatus === 'approved') {
        return {
          ...item,
          title: `${item.countryName} ${item.visaTitle}`,
          status: 'approved',
          statusText: '已出签',
          date: this.formatDate(item.resultAt || item.createdAt),
          desc: '本次申请结果已记录为出签，可点击查看当时的办理步骤。',
        }
      }
      if (item.resultStatus === 'rejected') {
        return {
          ...item,
          title: `${item.countryName} ${item.visaTitle}`,
          status: 'rejected',
          statusText: '被拒签',
          date: this.formatDate(item.resultAt || item.createdAt),
          desc: item.resultNote || '本次申请结果已记录为拒签，可回看当时步骤并准备再次申请。',
        }
      }
      if (item.resultStatus === 'withdrawn') {
        return {
          ...item,
          title: `${item.countryName} ${item.visaTitle}`,
          status: 'withdrawn',
          statusText: '已撤签',
          date: this.formatDate(item.resultAt || item.createdAt),
          desc: item.resultNote || '本次申请已记录为撤签。',
        }
      }
      if (Number(item.progress || 0) >= 100) {
        return {
          ...item,
          title: `${item.countryName} ${item.visaTitle}`,
          status: 'completed',
          statusText: '待结果',
          date: this.formatDate(item.createdAt),
          desc: '申请流程已完成，等待你确认这次 DIY 申请的最终结果。',
        }
      }
      return {
        ...item,
        title: `${item.countryName} ${item.visaTitle}`,
        status: 'progress',
        statusText: '进行中',
        date: this.formatDate(item.createdAt),
        desc: `当前申请进度 ${item.progress || 0}% ，可点击查看当时生成的步骤快照。`,
      }
    },
    async loadHistory() {
      try {
        const plans = (await api.listPlans()) || []
        this.list = plans
          .filter((item) => (item.resultStatus || 'pending') !== 'pending')
          .map((item) => this.normalizeHistoryItem(item))
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '历史记录加载失败',
          icon: 'none',
        })
      }
    },
    openPlanDetail(item) {
      uni.navigateTo({
        url: `/pages/visago/profile/history/detail/index?planId=${item.id}`,
      })
    },
    reapply(item) {
      uni.navigateTo({
        url: `/pages/visago/visa-country/index?countryId=${item.countryId}`,
      })
    },
  },
}
</script>

<style scoped>
.sub-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.sub-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.sub-content {
  padding: 18px 16px 34px;
  box-sizing: border-box;
}

.filter-row {
  margin: 0 0 12px;
  display: flex;
  gap: 8px;
}

.filter-chip {
  padding: 8px 13px;
  border-radius: 999px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  color: var(--visago-text-muted);
  font-size: 13px;
  font-weight: 800;
}

.filter-chip--active {
  background: var(--visago-primary);
  border-color: var(--visago-primary);
  color: #fff;
}

.timeline {
  display: grid;
  gap: 12px;
}

.history-card {
  position: relative;
  padding-left: 18px;
}

.timeline-dot {
  position: absolute;
  left: 0;
  top: 22px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--visago-primary);
}

.timeline-dot--approved {
  background: #16a34a;
}

.timeline-dot--rejected {
  background: #ef4444;
}

.timeline-dot--withdrawn {
  background: #64748b;
}

.history-main {
  padding: 14px;
  border-radius: 18px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
}

.history-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.history-title,
.history-date,
.history-desc {
  display: block;
}

.history-title {
  font-size: 16px;
  font-weight: 900;
}

.history-status {
  padding: 4px 8px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 900;
  color: var(--visago-primary);
  background: rgba(15, 101, 216, 0.12);
  flex-shrink: 0;
}

.history-status--approved {
  color: #16a34a;
  background: rgba(22, 163, 74, 0.12);
}

.history-status--rejected {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.12);
}

.history-status--withdrawn {
  color: #475569;
  background: rgba(100, 116, 139, 0.14);
}

.history-date {
  margin-top: 5px;
  font-size: 12px;
  color: var(--visago-text-soft);
}

.history-desc {
  margin-top: 8px;
  font-size: 13px;
  line-height: 1.5;
  color: var(--visago-text-muted);
}

.history-actions {
  margin-top: 12px;
  display: flex;
  gap: 8px;
}

.ghost-btn,
.primary-btn {
  height: 34px;
  padding: 0 12px;
  border-radius: 999px;
  display: flex;
  align-items: center;
  font-size: 12px;
  font-weight: 900;
}

.ghost-btn {
  background: var(--visago-surface-soft);
  color: var(--visago-text-muted);
}

.primary-btn {
  background: var(--visago-primary);
  color: #fff;
}
</style>
