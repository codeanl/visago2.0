<template>
  <view class="profile-page" :class="themeMode === 'dark' ? 'profile-page--dark' : 'profile-page--light'">
    <view class="profile-topbar">
      <view class="profile-topbar__inner visago-page-width">
        <view class="profile-topbar__left" />
        <view class="profile-topbar__title-wrap">
          <text class="profile-topbar__title">我的</text>
        </view>
        <view class="profile-topbar__actions">
          <view class="hero-action-btn" @tap="toggleThemeMode">
            <text class="material-symbols-outlined">{{ themeMode === 'dark' ? 'light_mode' : 'dark_mode' }}</text>
          </view>
        </view>
      </view>
    </view>

    <scroll-view scroll-y class="profile-scroll">
      <view class="profile-hero">
        <view class="profile-content visago-page-width">
          <view class="identity-row" @tap="openAccountProfile">
            <image class="avatar" :src="profileForm.avatar || defaultAvatar" mode="aspectFill" />
            <view class="identity-copy">
              <view class="name-row">
                <text class="name">{{ displayName }}</text>
              </view>
              <!--
              <view class="member-row">
                <text class="member-badge" v-if="profileForm.membership.hasMembership">{{ profileForm.membership.planName }}</text>
              </view>
              -->
            </view>
          </view>

          <!--
          <view class="vip-banner" @tap="openMembershipPage">
            <view class="vip-copy">
              <text class="vip-title">{{ profileForm.membership.hasMembership ? `${profileForm.membership.planName} VIP` : '未开通 VIP' }}</text>
              <text class="vip-sub">
                {{
                  profileForm.membership.hasMembership
                    ? `有效期至 ${formatDate(profileForm.membership.expiresAt)}，20+ 权益已解锁`
                    : '开通会员，解锁 20+ 权益'
                }}
              </text>
              <view class="vip-note-row">
                <text class="vip-note">{{ profileForm.membership.hasMembership ? '查看会员权益与到期时间' : '新客专享，联系客服可领取体验权益' }}</text>
              </view>
            </view>
            <view class="vip-cta">{{ profileForm.membership.hasMembership ? '查看权益' : '开通会员' }}</view>
          </view>
          -->
        </view>
      </view>

      <view class="profile-content profile-content--body visago-page-width">
        <view class="action-grid">
          <view
            v-for="item in profileActionCards"
            :key="item.key"
            class="action-card"
            :class="{ 'action-card--active': item.active }"
            @tap="onActionTap(item)"
          >
            <view class="action-icon-wrap" :style="{ background: item.tint, color: item.color }">
              <text class="material-symbols-outlined">{{ item.icon }}</text>
            </view>
            <text class="action-title">{{ item.title }}</text>
          </view>
        </view>

        <view class="settings-card">
          <view v-for="(item, idx) in profileSettings" :key="item.key" class="setting-row" @tap="onSettingTap(item)">
            <view class="setting-left">
              <text class="material-symbols-outlined setting-icon">{{ item.icon }}</text>
              <text class="setting-label">{{ item.title }}</text>
            </view>
            <text class="material-symbols-outlined setting-arrow">chevron_right</text>
            <view v-if="idx < profileSettings.length - 1" class="row-divider" />
          </view>
        </view>
      </view>
    </scroll-view>

    <VisagoBottomNav active-tab="profile" />
  </view>
</template>

<script>
import VisagoBottomNav from '../../../components/VisagoBottomNav.vue'
import { profileActions, profileSettings } from '../../../utils/visago-data'
import { api } from '../../../utils/api'
import { applyTheme, getStoredTheme, toggleTheme } from '../../../utils/theme'

const defaultAvatar = 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?auto=format&fit=crop&w=300&q=80'
const MY_POSTS_ACTION = {
  key: 'myPosts',
  title: '我的帖子',
  icon: 'forum',
  tint: '#dbeafe',
  color: '#2563eb',
}

function normalizeMembership(profile = {}) {
  const source = profile && typeof profile.membership === 'object' && profile.membership ? profile.membership : {}
  const planKey = String(source.planKey || profile.membershipPlanKey || '').trim()
  const planName = String(source.planName || profile.membershipPlanName || '').trim()
  const startedAt = String(source.startedAt || profile.membershipStartedAt || '').trim()
  const expiresAt = String(source.expiresAt || profile.membershipExpiresAt || '').trim()
  let status = String(source.status || profile.membershipStatus || '').trim()
  const hasMembership =
    typeof source.hasMembership === 'boolean'
      ? source.hasMembership
      : Boolean(planKey || planName || startedAt || expiresAt || status)
  if (!status && expiresAt) {
    const expireTime = new Date(expiresAt).getTime()
    if (!Number.isNaN(expireTime)) {
      status = expireTime > Date.now() ? 'active' : 'expired'
    }
  }
  return {
    hasMembership,
    planKey,
    planName,
    startedAt,
    expiresAt,
    status,
  }
}

function createProfileForm(profile = {}) {
  const membership = normalizeMembership(profile)
  return {
    id: profile.id || 0,
    avatar: profile.avatar || '',
    nickname: profile.nickname || '',
    name: profile.name || '',
    bio: profile.bio || '',
    membership,
  }
}

export default {
  components: {
    VisagoBottomNav,
  },
  data() {
    return {
      defaultAvatar,
      profileActionCards: [
        MY_POSTS_ACTION,
        ...profileActions
          .filter((item) => ['favorite', 'history'].includes(item.key))
          .map((item) => ({ ...item })),
      ],
      profileSettings: profileSettings.filter((item) => item.key !== 'theme'),
      themeMode: 'light',
      profileForm: createProfileForm(),
    }
  },
  computed: {
    displayName() {
      return this.profileForm.nickname || this.profileForm.name || '流浪启条'
    },
  },
  onShow() {
    this.themeMode = applyTheme(getStoredTheme())
    this.loadDashboard()
  },
  methods: {
    formatDate(value) {
      if (!value) return '-'
      const date = new Date(value)
      if (Number.isNaN(date.getTime())) return String(value)
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
    },
    async loadDashboard() {
      try {
        const profile = await api.me()
        this.profileForm = createProfileForm(profile)
      } catch (error) {
      }
    },
    openAccountProfile() {
      uni.navigateTo({
        url: '/pages/visago/profile/account/index',
      })
    },
    onSettingTap(item) {
      const routeMap = {
        notify: '/pages/visago/profile/settings/notify/index',
        privacy: '/pages/visago/profile/settings/privacy/index',
        help: '/pages/visago/profile/settings/help/index',
        about: '/pages/visago/profile/settings/about/index',
      }
      const url = routeMap[item.key]
      if (!url) return
      uni.navigateTo({ url })
    },
    onActionTap(item) {
      const routeMap = {
        myPosts: '/pages/visago/community/mine/index',
        favorite: '/pages/visago/profile/favorites/index',
        history: '/pages/visago/profile/history/index',
      }
      const url = routeMap[item.key]
      if (!url) return
      uni.navigateTo({ url })
    },
    toggleThemeMode() {
      this.themeMode = toggleTheme(this.themeMode)
    },
  },
}
</script>

<style scoped>
.profile-page {
  --profile-page-bg: #f4f6fb;
  --profile-body-bg: #f4f6fb;
  --profile-hero-bg:
    radial-gradient(circle at 16% 0%, rgba(199, 210, 226, 0.72) 0%, rgba(199, 210, 226, 0) 32%),
    linear-gradient(180deg, #eef3fa 0%, #e2e9f4 38%, #f4f6fb 100%);
  --profile-action-bg: rgba(255, 255, 255, 0.76);
  --profile-action-line: rgba(201, 211, 225, 0.82);
  --profile-title: #1f2533;
  --profile-sub: #6c7384;
  --profile-muted: #8b92a6;
  --profile-member-bg: rgba(255, 255, 255, 0.84);
  --profile-member-text: #666d7d;
  --profile-card-bg: #ffffff;
  --profile-card-line: rgba(217, 223, 233, 0.92);
  --profile-card-shadow: 0 10px 24px rgba(26, 36, 61, 0.08);
  --profile-topbar-bg: rgba(244, 246, 251, 0.82);
  min-height: 100vh;
  background: var(--profile-page-bg);
}

.profile-page--dark {
  --profile-page-bg: #050507;
  --profile-body-bg: #050507;
  --profile-hero-bg:
    radial-gradient(circle at 20% 0%, rgba(86, 96, 123, 0.22) 0%, rgba(86, 96, 123, 0) 30%),
    linear-gradient(180deg, #2f333b 0%, #191b20 38%, #0b0b0f 100%);
  --profile-action-bg: rgba(255, 255, 255, 0.08);
  --profile-action-line: rgba(255, 255, 255, 0.06);
  --profile-title: #f4f5f7;
  --profile-sub: #c5c9d3;
  --profile-muted: #7c8192;
  --profile-member-bg: rgba(255, 255, 255, 0.08);
  --profile-member-text: rgba(255, 255, 255, 0.5);
  --profile-card-bg: #1a1a20;
  --profile-card-line: rgba(255, 255, 255, 0.04);
  --profile-card-shadow: none;
  --profile-topbar-bg: rgba(5, 5, 7, 0.82);
}

.profile-page--light {
  --profile-page-bg: #f4f6fb;
}

.profile-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.profile-topbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 120;
  background: var(--profile-topbar-bg);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--profile-card-line);
}

.profile-topbar__inner {
  position: relative;
  height: 74px;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  padding: 0 14px 12px;
  box-sizing: border-box;
}

.profile-topbar__left {
  width: 40px;
  height: 40px;
  flex-shrink: 0;
}

.profile-topbar__title-wrap {
  position: absolute;
  left: 50%;
  bottom: 12px;
  transform: translateX(-50%);
  min-height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.profile-topbar__title {
  font-size: 34rpx;
  font-weight: 700;
  color: var(--profile-title);
}

.profile-topbar__actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.profile-content {
  box-sizing: border-box;
  padding: 0 14px;
}

.profile-content--body {
  padding-top: 18px;
  padding-bottom: calc(116px + var(--visago-safe-bottom));
  background: var(--profile-body-bg);
}

.profile-hero {
  position: relative;
  overflow: hidden;
  padding-top: 12px;
  padding-bottom: 18px;
  background: var(--profile-hero-bg);
}

.hero-action-btn {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  background: var(--profile-action-bg);
  color: var(--profile-title);
  border: 1px solid var(--profile-action-line);
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(10px);
}

.hero-action-btn .material-symbols-outlined {
  font-size: 16px;
}

.identity-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 18px;
}

.avatar {
  width: 58px;
  height: 58px;
  border-radius: 999px;
  border: 2px solid rgba(255, 255, 255, 0.2);
  flex-shrink: 0;
}

.identity-copy {
  min-width: 0;
  flex: 1;
}

.name-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.name {
  font-size: 18px;
  font-weight: 800;
  color: var(--profile-title);
}

.member-row {
  margin-top: 6px;
}

.member-badge {
  display: inline-flex;
  align-items: center;
  height: 20px;
  padding: 0 7px;
  border-radius: 999px;
  background: var(--profile-member-bg);
  color: var(--profile-member-text);
  font-size: 10px;
  font-weight: 700;
}

.vip-banner {
  margin-top: 16px;
  border-radius: 16px;
  padding: 16px;
  background:
    linear-gradient(115deg, rgba(131, 104, 28, 0.2) 0%, rgba(131, 104, 28, 0) 58%),
    linear-gradient(180deg, #2b2412 0%, #201b10 100%);
  border: 1px solid rgba(145, 117, 47, 0.42);
  box-shadow: inset 0 1px 0 rgba(255, 215, 130, 0.08);
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.profile-page--light .vip-banner {
  background:
    linear-gradient(115deg, rgba(211, 185, 110, 0.22) 0%, rgba(211, 185, 110, 0.04) 58%),
    linear-gradient(180deg, #f7efcf 0%, #efe1b5 100%);
  border-color: rgba(175, 145, 74, 0.45);
  box-shadow: 0 10px 24px rgba(156, 129, 63, 0.12);
}

.vip-copy {
  min-width: 0;
  flex: 1;
}

.vip-title,
.vip-sub,
.vip-note,
.action-title {
  display: block;
}

.vip-title {
  font-size: 15px;
  font-weight: 800;
  color: #e2c377;
}

.profile-page--light .vip-title {
  color: #7d6220;
}

.vip-sub {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.5;
  color: rgba(215, 185, 109, 0.72);
}

.profile-page--light .vip-sub {
  color: rgba(117, 92, 35, 0.75);
}

.vip-note-row {
  margin-top: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.vip-note {
  font-size: 11px;
  color: #f2a15f;
}

.profile-page--light .vip-note {
  color: #c57c39;
}

.vip-cta {
  height: 38px;
  min-width: 82px;
  padding: 0 12px;
  border-radius: 10px;
  background: linear-gradient(180deg, #9b8143 0%, #7c6735 100%);
  color: #f7edcf;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.16);
}

.profile-page--light .vip-cta {
  background: linear-gradient(180deg, #baa062 0%, #9f8447 100%);
  color: #fff8e2;
}

.action-grid {
  margin-top: 6px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.action-card,
.settings-card {
  border-radius: 18px;
  background: var(--profile-card-bg);
  border: 1px solid var(--profile-card-line);
  box-shadow: var(--profile-card-shadow);
}

.action-card {
  min-height: 94px;
  padding: 12px 8px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  text-align: center;
}

.action-card--active {
  box-shadow: inset 0 0 0 1px rgba(126, 136, 214, 0.22);
}

.action-icon-wrap {
  width: 40px;
  height: 40px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.action-icon-wrap .material-symbols-outlined {
  font-size: 20px;
}

.action-title {
  font-size: 12px;
  font-weight: 700;
  color: var(--profile-title);
}

.setting-arrow {
  color: var(--profile-muted);
  font-size: 18px;
}

.settings-card {
  margin-top: 18px;
  overflow: hidden;
}

.setting-row {
  position: relative;
  min-height: 54px;
  padding: 0 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.setting-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.setting-icon {
  color: var(--profile-muted);
  font-size: 20px;
}

.setting-label {
  font-size: 14px;
  color: var(--profile-title);
}

.row-divider {
  position: absolute;
  left: 54px;
  right: 0;
  bottom: 0;
  height: 1px;
  background: var(--profile-card-line);
}
</style>
