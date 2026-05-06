<template>
  <view class="auth-page">
    <view class="auth-body visago-page-width">
      <view class="hero">
        <view class="hero-icon-wrap">
          <text class="material-symbols-outlined hero-icon">travel_explore</text>
        </view>
        <text class="hero-title">{{ text.welcomeTitle }}</text>
        <text class="hero-sub">{{ text.welcomeSubtitle }}</text>
      </view>

      <view class="form-area">
        <view class="field-block">
          <text class="field-label">{{ text.phoneLabel }}</text>
          <input v-model.trim="form.phone" class="field-input" :placeholder="text.phonePlaceholder" />
        </view>

        <view class="field-block">
          <view class="field-head">
            <text class="field-label">{{ text.passwordLabel }}</text>
            <text class="mini-link" @tap="onForgot">{{ text.forgotPassword }}</text>
          </view>
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

        <view class="submit-btn" :class="{ 'submit-btn--disabled': submitting }" @tap="submitLogin">
          {{ submitting ? text.loggingIn : text.login }}
        </view>

        <view class="legal-row">
          <text class="legal-text">登录即代表你已阅读并同意</text>
          <text class="legal-link" @tap="openAgreement">《用户协议》</text>
          <text class="legal-text">和</text>
          <text class="legal-link" @tap="openPrivacyPolicy">《隐私政策》</text>
        </view>
      </view>

      <view class="bottom-link">
        <text>{{ text.noAccount }}</text>
        <text class="link-action" @tap="goRegister">{{ text.goRegister }}</text>
      </view>
    </view>
  </view>
</template>

<script>
import { applyTheme, getStoredTheme } from '../../../../utils/theme'
import { HOME_PAGE, REGISTER_PAGE, RESET_PASSWORD_PAGE, isLoggedIn, saveLoginSession } from '../../../../utils/auth'
import { api } from '../../../../utils/api'

const LOGIN_TEXT = Object.freeze({
  welcomeTitle: '\u6b22\u8fce\u56de\u6765',
  welcomeSubtitle: '\u767b\u5f55\u540e\u7ee7\u7eed\u4f60\u7684\u7b7e\u8bc1\u6d41\u7a0b',
  phoneLabel: '\u624b\u673a\u53f7',
  phonePlaceholder: '\u8bf7\u8f93\u5165\u624b\u673a\u53f7',
  passwordLabel: '\u5bc6\u7801',
  forgotPassword: '\u5fd8\u8bb0\u5bc6\u7801\uff1f',
  passwordPlaceholder: '\u8bf7\u8f93\u5165\u5bc6\u7801',
  loggingIn: '\u767b\u5f55\u4e2d...',
  login: '\u767b\u5f55',
  noAccount: '\u8fd8\u6ca1\u6709\u8d26\u53f7\uff1f',
  goRegister: '\u53bb\u6ce8\u518c',
  needPhoneAndPassword: '\u8bf7\u8f93\u5165\u624b\u673a\u53f7\u548c\u5bc6\u7801',
  loginResponseInvalid: '\u767b\u5f55\u54cd\u5e94\u5f02\u5e38',
  loginSuccess: '\u767b\u5f55\u6210\u529f',
  loginFailed: '\u767b\u5f55\u5931\u8d25',
})

export default {
  data() {
    return {
      form: {
        phone: '',
        password: '',
      },
      showPassword: false,
      submitting: false,
      text: LOGIN_TEXT,
    }
  },
  onShow() {
    applyTheme(getStoredTheme())
    if (isLoggedIn()) {
      uni.reLaunch({ url: HOME_PAGE })
    }
  },
  methods: {
    resolveLoginError(error) {
      const message = String((error && error.message) || '').trim()
      if (!message) return this.text.loginFailed
      if (message.includes('密码错误')) {
        this.form.password = ''
        return '密码错误，请重新输入'
      }
      if (message.includes('手机号不存在')) {
        return '手机号不存在，请先注册'
      }
      return message
    },
    async submitLogin() {
      if (this.submitting) return
      if (!this.form.phone || !this.form.password) {
        uni.showToast({
          title: this.text.needPhoneAndPassword,
          icon: 'none',
        })
        return
      }

      this.submitting = true
      try {
        const data = await api.login({
          phone: this.form.phone,
          password: this.form.password,
        })
        if (!data || !data.token) {
          throw new Error(this.text.loginResponseInvalid)
        }
        saveLoginSession(data.token, data.user || null)
        uni.showToast({
          title: this.text.loginSuccess,
          icon: 'none',
        })
        setTimeout(() => {
          uni.reLaunch({ url: HOME_PAGE })
        }, 180)
      } catch (error) {
        uni.showToast({
          title: this.resolveLoginError(error),
          icon: 'none',
        })
      } finally {
        this.submitting = false
      }
    },
    goRegister() {
      uni.navigateTo({ url: REGISTER_PAGE })
    },
    onForgot() {
      uni.navigateTo({ url: RESET_PASSWORD_PAGE })
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
  margin-top: 18px;
  text-align: center;
}

.hero-icon-wrap {
  width: 68px;
  height: 68px;
  border-radius: 16px;
  background: var(--visago-primary);
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--visago-shadow-fab);
}

.hero-icon {
  font-size: 40px;
  color: #fff;
  font-variation-settings: 'FILL' 1;
}

.hero-title {
  margin-top: 18px;
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
  margin-top: 30px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.field-block {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.field-label {
  font-size: 13px;
  color: var(--visago-text-muted);
}

.mini-link {
  font-size: 13px;
  color: var(--visago-primary);
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

.submit-btn {
  margin-top: 10px;
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

.legal-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  flex-wrap: wrap;
  margin-top: 4px;
}

.legal-text,
.legal-link {
  font-size: 12px;
  line-height: 1.5;
}

.legal-text {
  color: var(--visago-text-soft);
}

.legal-link {
  color: var(--visago-primary);
  font-weight: 600;
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
