<template>
  <view class="auth-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="找回密码" />

    <view class="auth-body visago-page-width">
      <view class="hero">
        <text class="hero-title">{{ text.title }}</text>
        <text class="hero-sub">{{ text.subtitle }}</text>
      </view>

      <view class="form-area">
        <view class="field-block">
          <text class="field-label">{{ text.nameLabel }}</text>
          <input v-model.trim="form.name" class="field-input" :placeholder="text.namePlaceholder" />
        </view>

        <view class="field-block">
          <text class="field-label">{{ text.phoneLabel }}</text>
          <input v-model.trim="form.phone" class="field-input" :placeholder="text.phonePlaceholder" />
        </view>

        <view class="field-block">
          <text class="field-label">{{ text.passwordLabel }}</text>
          <view class="password-wrap">
            <input
              v-model="form.password"
              :password="!showPassword"
              class="field-input field-input--password"
              :placeholder="text.passwordPlaceholder"
            />
            <text class="material-symbols-outlined eye-btn" @tap="showPassword = !showPassword">
              {{ showPassword ? 'visibility_off' : 'visibility' }}
            </text>
          </view>
        </view>

        <view class="field-block">
          <text class="field-label">{{ text.confirmPasswordLabel }}</text>
          <input
            v-model="form.confirmPassword"
            :password="!showPassword"
            class="field-input"
            :placeholder="text.confirmPasswordPlaceholder"
          />
        </view>

        <text class="helper-text">{{ text.helper }}</text>

        <view class="submit-btn" :class="{ 'submit-btn--disabled': submitting }" @tap="submitReset">
          {{ submitting ? text.submitting : text.submit }}
        </view>
      </view>

      <view class="bottom-link">
        <text>{{ text.backToLoginLabel }}</text>
        <text class="link-action" @tap="goLogin">{{ text.backToLogin }}</text>
      </view>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'
import { LOGIN_PAGE } from '../../../../utils/auth'
import { api } from '../../../../utils/api'

const RESET_TEXT = Object.freeze({
  title: '重置登录密码',
  subtitle: '通过手机号和昵称核验后，重新设置新的登录密码',
  nameLabel: '昵称',
  namePlaceholder: '请输入当前昵称',
  phoneLabel: '手机号',
  phonePlaceholder: '请输入注册手机号',
  passwordLabel: '新密码',
  passwordPlaceholder: '请至少输入 8 位新密码',
  confirmPasswordLabel: '确认新密码',
  confirmPasswordPlaceholder: '请再次输入新密码',
  helper: '当前版本使用手机号 + 昵称进行核验，重置成功后请使用新密码登录。',
  submit: '重置密码',
  submitting: '提交中...',
  backToLoginLabel: '已经想起密码？',
  backToLogin: '返回登录',
  needCompleteInfo: '请完整填写信息',
  passwordTooShort: '密码至少 8 位',
  passwordNotSame: '两次输入的密码不一致',
  resetSuccess: '密码已重置，请重新登录',
  resetFailed: '重置失败',
})

export default {
  components: { VisagoTopBar },
  data() {
    return {
      form: {
        name: '',
        phone: '',
        password: '',
        confirmPassword: '',
      },
      showPassword: false,
      submitting: false,
      text: RESET_TEXT,
    }
  },
  onShow() {
    applyTheme(getStoredTheme())
  },
  methods: {
    async submitReset() {
      if (this.submitting) return
      if (!this.form.name || !this.form.phone || !this.form.password || !this.form.confirmPassword) {
        uni.showToast({ title: this.text.needCompleteInfo, icon: 'none' })
        return
      }
      if (this.form.password.length < 8) {
        uni.showToast({ title: this.text.passwordTooShort, icon: 'none' })
        return
      }
      if (this.form.password !== this.form.confirmPassword) {
        uni.showToast({ title: this.text.passwordNotSame, icon: 'none' })
        return
      }

      this.submitting = true
      try {
        await api.resetPassword({
          name: this.form.name,
          phone: this.form.phone,
          password: this.form.password,
        })
        uni.showToast({ title: this.text.resetSuccess, icon: 'none' })
        setTimeout(() => {
          uni.reLaunch({ url: LOGIN_PAGE })
        }, 220)
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || this.text.resetFailed,
          icon: 'none',
        })
      } finally {
        this.submitting = false
      }
    },
    goLogin() {
      uni.reLaunch({ url: LOGIN_PAGE })
    },
  },
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  background: var(--visago-bg);
}

.auth-body {
  box-sizing: border-box;
  min-height: 100vh;
  padding: 118px 22px calc(30px + var(--visago-safe-bottom));
  display: flex;
  flex-direction: column;
}

.hero-title,
.hero-sub,
.helper-text {
  display: block;
}

.hero-title {
  font-size: 30px;
  line-height: 1.2;
  font-weight: 700;
  color: var(--visago-text);
}

.hero-sub {
  margin-top: 6px;
  font-size: 14px;
  line-height: 1.5;
  color: var(--visago-text-muted);
}

.form-area {
  margin-top: 26px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.field-block {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-size: 13px;
  color: var(--visago-text-muted);
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

.password-wrap {
  position: relative;
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

.helper-text {
  font-size: 12px;
  line-height: 1.55;
  color: var(--visago-text-soft);
}

.submit-btn {
  margin-top: 6px;
  height: 46px;
  border-radius: 12px;
  background: var(--visago-primary);
  color: #fff;
  font-size: 16px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.submit-btn--disabled {
  opacity: 0.7;
}

.bottom-link {
  margin-top: auto;
  text-align: center;
  font-size: 14px;
  color: var(--visago-text-muted);
}

.link-action {
  margin-left: 6px;
  color: var(--visago-primary);
  font-weight: 600;
}
</style>
