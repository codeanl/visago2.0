<template>
  <view class="delete-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="注销账户" />

    <scroll-view scroll-y class="delete-scroll">
      <view class="delete-content visago-page-width">
        <view class="warning-card">
          <view class="warning-icon">
            <text class="material-symbols-outlined">warning</text>
          </view>
          <view class="warning-copy">
            <text class="warning-title">注销后不可恢复</text>
            <text class="warning-desc">账号、计划、会员信息以及社区内容都会被一并删除，请谨慎操作。</text>
          </view>
        </view>

        <view class="section-card">
          <text class="section-title">注销后会发生什么</text>
          <view class="bullet-list">
            <view v-for="item in consequences" :key="item" class="bullet-row">
              <view class="bullet-dot" />
              <text class="bullet-text">{{ item }}</text>
            </view>
          </view>
        </view>

        <view class="section-card">
          <text class="section-title">确认身份</text>

          <view class="field-block">
            <text class="field-label">当前手机号</text>
            <text class="readonly-value">{{ profile.phone || '未读取到手机号' }}</text>
          </view>

          <view class="field-block">
            <text class="field-label">登录密码</text>
            <view class="password-wrap">
              <input
                v-model="password"
                :password="!showPassword"
                class="field-input field-input--password"
                placeholder="请输入当前登录密码"
              />
              <text class="material-symbols-outlined eye-btn" @tap="showPassword = !showPassword">
                {{ showPassword ? 'visibility_off' : 'visibility' }}
              </text>
            </view>
          </view>

          <view class="agree-row" @tap="confirmed = !confirmed">
            <view class="checkbox" :class="{ 'checkbox--checked': confirmed }">
              <text v-if="confirmed" class="material-symbols-outlined">check</text>
            </view>
            <text class="agree-text">我已知晓：社区帖子、评论、收藏、举报记录会随账号一起删除，且无法恢复。</text>
          </view>
        </view>

        <view class="action-row">
          <view class="cancel-btn" @tap="goBack">再想想</view>
          <view class="delete-btn" :class="{ 'delete-btn--disabled': submitting }" @tap="submitDelete">
            {{ submitting ? '注销中...' : '确认注销' }}
          </view>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../../../components/VisagoTopBar.vue'
import { api } from '../../../../../../utils/api'
import { LOGIN_PAGE, logoutSession } from '../../../../../../utils/auth'
import { applyTheme, getStoredTheme } from '../../../../../../utils/theme'

export default {
  components: { VisagoTopBar },
  data() {
    return {
      profile: {},
      password: '',
      showPassword: false,
      confirmed: false,
      submitting: false,
      consequences: [
        '账号资料、头像、偏好设置会被删除。',
        '签证计划、申请历史、会员信息与相关记录会被删除。',
        '你发布的社区帖子、评论、收藏、举报和屏蔽记录会随账号一起删除。',
        '注销完成后无法恢复，请在操作前自行备份需要的信息。',
      ],
    }
  },
  onShow() {
    applyTheme(getStoredTheme())
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
    goBack() {
      uni.navigateBack({
        fail: () => {
          uni.reLaunch({ url: '/pages/visago/profile/settings/privacy/index' })
        },
      })
    },
    async submitDelete() {
      if (this.submitting) return
      if (!this.password.trim()) {
        uni.showToast({ title: '请输入当前密码', icon: 'none' })
        return
      }
      if (!this.confirmed) {
        uni.showToast({ title: '请先勾选删除确认', icon: 'none' })
        return
      }

      this.submitting = true
      try {
        await api.deleteMyAccount({
          password: this.password.trim(),
          confirm: true,
        })
        logoutSession()
        uni.showToast({
          title: '账户已注销',
          icon: 'none',
        })
        setTimeout(() => {
          uni.reLaunch({ url: LOGIN_PAGE })
        }, 220)
      } catch (error) {
        if (error && error.message && String(error.message).includes('当前密码错误')) {
          this.password = ''
        }
        uni.showToast({
          title: (error && error.message) || '注销失败',
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
.delete-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.delete-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.delete-content {
  padding: 18px 16px 34px;
  box-sizing: border-box;
}

.warning-card,
.section-card {
  border-radius: 18px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
}

.warning-card {
  padding: 16px;
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.warning-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: rgba(239, 68, 68, 0.12);
  color: #ef4444;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.warning-title,
.warning-desc,
.section-title,
.bullet-text,
.field-label,
.readonly-value,
.agree-text {
  display: block;
}

.warning-title {
  font-size: 16px;
  font-weight: 800;
  color: var(--visago-text);
}

.warning-desc {
  margin-top: 6px;
  font-size: 13px;
  line-height: 1.6;
  color: var(--visago-text-muted);
}

.section-card {
  margin-top: 16px;
  padding: 16px;
}

.section-title {
  font-size: 15px;
  font-weight: 800;
  color: var(--visago-text);
}

.bullet-list {
  margin-top: 12px;
  display: grid;
  gap: 10px;
}

.bullet-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.bullet-dot {
  width: 8px;
  height: 8px;
  margin-top: 7px;
  border-radius: 999px;
  background: #ef4444;
  flex-shrink: 0;
}

.bullet-text {
  font-size: 13px;
  line-height: 1.65;
  color: var(--visago-text-muted);
}

.field-block {
  margin-top: 14px;
  display: grid;
  gap: 6px;
}

.field-label {
  font-size: 13px;
  color: var(--visago-text-muted);
}

.readonly-value {
  font-size: 15px;
  font-weight: 700;
  color: var(--visago-text);
}

.password-wrap {
  position: relative;
}

.field-input {
  height: 50px;
  border-radius: 12px;
  background: var(--visago-surface-soft);
  padding: 0 14px;
  box-sizing: border-box;
  font-size: 15px;
  color: var(--visago-text);
}

.field-input--password {
  padding-right: 46px;
}

.eye-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--visago-text-soft);
  font-size: 21px;
}

.agree-row {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-top: 16px;
}

.checkbox {
  width: 19px;
  height: 19px;
  border-radius: 6px;
  border: 1px solid var(--visago-line);
  background: var(--visago-surface);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.checkbox--checked {
  background: #ef4444;
  border-color: #ef4444;
}

.checkbox .material-symbols-outlined {
  font-size: 14px;
}

.agree-text {
  font-size: 13px;
  line-height: 1.55;
  color: var(--visago-text-muted);
}

.action-row {
  margin-top: 24px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.cancel-btn,
.delete-btn {
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 15px;
  font-weight: 700;
}

.cancel-btn {
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  color: var(--visago-text-muted);
}

.delete-btn {
  background: #ef4444;
  color: #fff;
}

.delete-btn--disabled {
  opacity: 0.72;
}
</style>
