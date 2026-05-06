<template>
  <view class="auth-page">
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

        <view class="agree-row" @tap="toggleAgree">
          <view class="checkbox" :class="{ 'checkbox--checked': agreed }">
            <text v-if="agreed" class="material-symbols-outlined">check</text>
          </view>
          <view class="agree-copy">
            <text class="agree-text">我已阅读并同意</text>
            <text class="agree-link" @tap.stop="openAgreement">《用户协议》</text>
            <text class="agree-text">与</text>
            <text class="agree-link" @tap.stop="openPrivacyPolicy">《隐私政策》</text>
          </view>
        </view>

        <view class="submit-btn" :class="{ 'submit-btn--disabled': submitting }" @tap="submitRegister">
          {{ submitting ? text.submitting : text.register }}
        </view>
      </view>

      <view class="bottom-link">
        <text>{{ text.hasAccount }}</text>
        <text class="link-action" @tap="goLogin">{{ text.goLogin }}</text>
      </view>
    </view>
  </view>
</template>

<script>
import { applyTheme, getStoredTheme } from '../../../../utils/theme'
import { LOGIN_PAGE } from '../../../../utils/auth'
import { api } from '../../../../utils/api'

const REGISTER_TEXT = Object.freeze({
  title: '\u521b\u5efa\u8d26\u53f7',
  subtitle: '\u52a0\u5165 Visago\uff0c\u5f00\u59cb\u4f60\u7684\u7b7e\u8bc1\u6d41\u7a0b',
  nameLabel: '\u59d3\u540d',
  namePlaceholder: '\u8bf7\u8f93\u5165\u59d3\u540d',
  phoneLabel: '\u624b\u673a\u53f7',
  phonePlaceholder: '\u8bf7\u8f93\u5165\u624b\u673a\u53f7',
  passwordLabel: '\u5bc6\u7801',
  passwordPlaceholder: '\u8bf7\u81f3\u5c11\u8f93\u5165 8 \u4f4d\u5bc6\u7801',
  submitting: '\u6ce8\u518c\u4e2d...',
  register: '\u6ce8\u518c',
  hasAccount: '\u5df2\u6709\u8d26\u53f7\uff1f',
  goLogin: '\u53bb\u767b\u5f55',
  needCompleteInfo: '\u8bf7\u5b8c\u6574\u586b\u5199\u4fe1\u606f',
  passwordTooShort: '\u5bc6\u7801\u81f3\u5c11 8 \u4f4d',
  needAgreement: '\u8bf7\u5148\u52fe\u9009\u534f\u8bae',
  registerSuccess: '\u6ce8\u518c\u6210\u529f\uff0c\u8bf7\u767b\u5f55',
  registerFailed: '\u6ce8\u518c\u5931\u8d25',
})

export default {
  data() {
    return {
      form: {
        name: '',
        phone: '',
        password: '',
      },
      showPassword: false,
      agreed: false,
      submitting: false,
      text: REGISTER_TEXT,
    }
  },
  onShow() {
    applyTheme(getStoredTheme())
  },
  methods: {
    toggleAgree() {
      this.agreed = !this.agreed
    },
    async submitRegister() {
      if (this.submitting) return
      if (!this.form.name || !this.form.phone || !this.form.password) {
        uni.showToast({
          title: this.text.needCompleteInfo,
          icon: 'none',
        })
        return
      }
      if (this.form.password.length < 8) {
        uni.showToast({
          title: this.text.passwordTooShort,
          icon: 'none',
        })
        return
      }
      if (!this.agreed) {
        uni.showToast({
          title: this.text.needAgreement,
          icon: 'none',
        })
        return
      }

      this.submitting = true
      try {
        await api.register({
          name: this.form.name,
          phone: this.form.phone,
          password: this.form.password,
        })
        uni.showToast({
          title: this.text.registerSuccess,
          icon: 'none',
        })
        setTimeout(() => {
          uni.reLaunch({ url: LOGIN_PAGE })
        }, 260)
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || this.text.registerFailed,
          icon: 'none',
        })
      } finally {
        this.submitting = false
      }
    },
    goLogin() {
      uni.reLaunch({ url: LOGIN_PAGE })
    },
    openAgreement() {
      uni.navigateTo({ url: '/pages/visago/profile/settings/about/agreement/index' })
    },
    openPrivacyPolicy() {
      uni.navigateTo({ url: '/pages/visago/profile/settings/about/privacy-policy/index' })
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
  padding: 44px 22px calc(30px + var(--visago-safe-bottom));
  display: flex;
  flex-direction: column;
}

.hero {
  margin-top: 14px;
}

.hero-title {
  display: block;
  font-size: 30px;
  line-height: 1.2;
  font-weight: 700;
  color: var(--visago-text);
}

.hero-sub {
  margin-top: 6px;
  display: block;
  font-size: 14px;
  color: var(--visago-text-muted);
}

.form-area {
  margin-top: 24px;
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

.agree-row {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-top: 2px;
}

.agree-copy {
  display: flex;
  align-items: center;
  gap: 2px;
  flex-wrap: wrap;
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
  background: var(--visago-primary);
  border-color: var(--visago-primary);
}

.checkbox .material-symbols-outlined {
  font-size: 14px;
}

.agree-text {
  font-size: 13px;
  line-height: 1.45;
  color: var(--visago-text-muted);
}

.agree-link {
  font-size: 13px;
  line-height: 1.45;
  color: var(--visago-primary);
  font-weight: 600;
}

.submit-btn {
  margin-top: 8px;
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
