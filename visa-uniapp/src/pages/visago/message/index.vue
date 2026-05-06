<template>
  <view class="message-page" :class="{ 'message-page--dark': themeMode === 'dark' }">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="消息" />

    <view class="message-content visago-page-width">
      <view class="search-box">
        <text class="material-symbols-outlined search-icon">search</text>
        <input v-model.trim="keyword" class="search-input" placeholder="Search messages" />
      </view>

      <view class="message-list">
        <view v-for="item in filteredMessages" :key="item.id" class="message-swipe">
          <view class="message-delete" @tap.stop="removeMessage(item.id)">删除</view>

          <view
            class="message-card"
            :class="{ 'message-card--dragging': item.dragging }"
            :style="{ transform: `translateX(${item.offsetX}px)` }"
            @touchstart="onTouchStart($event, item.id)"
            @touchmove="onTouchMove($event, item.id)"
            @touchend="onTouchEnd(item.id)"
            @touchcancel="onTouchEnd(item.id)"
          >
            <view class="message-row" @tap="toggleOpen(item.id)">
              <view v-if="item.unread" class="unread-dot" />
              <view class="avatar-wrap" :class="{ 'avatar-wrap--icon': !item.avatar }">
                <image v-if="item.avatar" class="avatar" :src="item.avatar" mode="aspectFill" />
                <text v-else class="material-symbols-outlined avatar-icon">{{ item.icon || 'description' }}</text>
              </view>

              <view class="message-main">
                <view class="title-row">
                  <text class="msg-title">{{ item.title }}</text>
                  <text class="msg-time" :class="{ 'msg-time--unread': item.unread }">{{ item.time }}</text>
                </view>
                <text class="msg-preview">{{ item.preview }}</text>
              </view>

              <text class="material-symbols-outlined expand-icon" :class="{ 'expand-icon--open': isOpen(item.id) }">expand_more</text>
            </view>

            <view v-if="isOpen(item.id)" class="message-detail">
              <text>{{ item.content }}</text>
            </view>
          </view>
        </view>

        <view v-if="!filteredMessages.length" class="empty-card visago-card">
          <text>没有找到匹配消息</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../components/VisagoTopBar.vue'
import { applyTheme, getStoredTheme, THEME_CHANGE_EVENT } from '../../../utils/theme'

const DELETE_WIDTH = 74

const MOCK_MESSAGES = [
  {
    id: 'support',
    title: 'Visago Support Team',
    time: '10:42 AM',
    unread: true,
    avatar:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuAiMOwJa_C0HY3ocYHP8AGtiZhAP-fX9g0LncZeVFEkt0bBEHGJuaBJpdPy-vspSGChX2d4JuA1gBPyiV-od_KT8Xg2LT-r3gThq1q7DrtoLMb9IZ7wvvyqe965i7CTUECkZsIP24CioO-Z3mnlf5EWbdZvuIwgeLKrlYzDP1ndL_oyshMCQXTnMpByp8r2JMJ68C48x84u-9zK3ARCUhoDHfyvNXn06YXGau_wXp6nLOAl2X3r3kNzWLRt1EWykOfm1R95_-F_Bow',
    preview: 'Your Japan E-Visa application has been approved...',
    content:
      'Your Japan E-Visa application has been approved. Please check your document vault to download the official PDF.',
  },
  {
    id: 'france',
    title: 'France Visa Advisor',
    time: 'Yesterday',
    unread: false,
    avatar:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuB_1HpE-2dH8z5Jio1EsBxcMI7x2HmUCulb7ckAbr-pvnU25LntrQJJjimr3dGXoatqRuhxGS_Zi-9I8I8ZWUeOe93l5jXszwWhIZCIt-xvAUYWwsXffW6tg8M-40zpr75-50KOkPbknabjZf59QLha4Dreo-bENbiaCZMhDKjaz5YUJl6hWXe5tGE1QZRg14tCfXZtp0aqCbDjZObLsC8HcEtMxqBBWumFp8SegMb9Cm6oi__P_NCiqW5YFjWJh0dmeidwdmn3Fck',
    preview: 'We received your financial statements. The processing...',
    content: 'We received your financial statements. Current processing time is about 10-15 business days.',
  },
  {
    id: 'doc',
    title: 'Document Verification',
    time: 'Tue',
    unread: false,
    icon: 'description',
    preview: 'Passport scan verified successfully.',
    content: 'Passport scan verified successfully. No further action is required for this step.',
  },
  {
    id: 'general',
    title: 'James - General Support',
    time: 'Oct 12',
    unread: false,
    avatar:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuBz-089yYNlvbPJ0G31qjBfpeN6Ln3D8OQYA87pnfLKWAvvMbCUUEjvt8ohY6eTyZR3_Q1zzFLJIOJuITPLXWVr8NMmRGp3QZ0_shgraP_ZRFWgBcFNZGL040a5JXuKxaXaODdxkEaCRKQuA4okqg4-bU_E14dCVVRDgDV6CXdpV5K6zXUBN0HLPeSUCiN2i9uplm-EhegHZtTvqSuUhyuQatbbc2XpXXkcFCEOe7wdY6Vp-SVf2lzl7fIeC9wX4_5KoZOPlu-JXm4',
    preview: 'Happy to help! Let me know if you need any further...',
    content: 'Happy to help! Let me know if you need any further assistance planning your itinerary.',
  },
]

export default {
  components: { VisagoTopBar },
  data() {
    return {
      themeMode: 'light',
      keyword: '',
      messages: MOCK_MESSAGES.map((item) => ({ ...item, offsetX: 0, dragging: false })),
      openIds: ['support'],
      touch: {
        id: '',
        startX: 0,
        startOffset: 0,
        moved: false,
      },
      tapLockUntil: 0,
    }
  },
  computed: {
    filteredMessages() {
      const word = this.keyword.trim().toLowerCase()
      if (!word) return this.messages
      return this.messages.filter((item) =>
        [item.title, item.preview, item.content].some((v) => String(v).toLowerCase().includes(word))
      )
    },
  },
  onLoad() {
    if (typeof uni !== 'undefined' && uni.$on) {
      uni.$on(THEME_CHANGE_EVENT, this.onThemeChange)
    }
    this.syncTheme()
  },
  onShow() {
    this.syncTheme()
  },
  onUnload() {
    if (typeof uni !== 'undefined' && uni.$off) {
      uni.$off(THEME_CHANGE_EVENT, this.onThemeChange)
    }
  },
  methods: {
    syncTheme() {
      this.themeMode = applyTheme(getStoredTheme())
    },
    onThemeChange(theme) {
      this.themeMode = theme === 'dark' ? 'dark' : 'light'
      applyTheme(this.themeMode)
    },
    isOpen(id) {
      return this.openIds.includes(id)
    },
    toggleOpen(id) {
      if (Date.now() < this.tapLockUntil) return
      const target = this.messages.find((item) => item.id === id)
      if (target && target.offsetX !== 0) {
        this.setOffset(id, 0)
        return
      }
      if (this.isOpen(id)) {
        this.openIds = this.openIds.filter((x) => x !== id)
      } else {
        this.openIds = [...this.openIds, id]
      }
    },
    setOffset(id, offset) {
      this.messages = this.messages.map((item) =>
        item.id === id ? { ...item, offsetX: offset } : item
      )
    },
    closeOthers(exceptId) {
      this.messages = this.messages.map((item) =>
        item.id !== exceptId && item.offsetX !== 0 ? { ...item, offsetX: 0 } : item
      )
    },
    onTouchStart(event, id) {
      const msg = this.messages.find((item) => item.id === id)
      this.closeOthers(id)
      this.touch.id = id
      this.touch.startX = event.touches[0].clientX
      this.touch.startOffset = msg ? msg.offsetX : 0
      this.touch.moved = false
      this.messages = this.messages.map((item) =>
        item.id === id ? { ...item, dragging: true } : item
      )
    },
    onTouchMove(event, id) {
      if (this.touch.id !== id) return
      const currentX = event.touches[0].clientX
      const deltaX = currentX - this.touch.startX
      if (Math.abs(deltaX) > 6) {
        this.touch.moved = true
      }
      const next = Math.max(-DELETE_WIDTH, Math.min(0, this.touch.startOffset + deltaX))
      this.setOffset(id, next)
    },
    onTouchEnd(id) {
      if (this.touch.id !== id) return
      const msg = this.messages.find((item) => item.id === id)
      const currentOffset = msg ? msg.offsetX : 0
      const finalOffset = currentOffset < -DELETE_WIDTH / 2 ? -DELETE_WIDTH : 0
      this.messages = this.messages.map((item) =>
        item.id === id ? { ...item, offsetX: finalOffset, dragging: false } : { ...item, dragging: false }
      )
      if (this.touch.moved) {
        this.tapLockUntil = Date.now() + 180
      }
      this.touch.id = ''
      this.touch.moved = false
    },
    removeMessage(id) {
      this.messages = this.messages.filter((item) => item.id !== id)
      this.openIds = this.openIds.filter((x) => x !== id)
    },
  },
}
</script>

<style scoped>
.message-page {
  min-height: 100vh;
  background: var(--visago-bg);
}

.message-content {
  box-sizing: border-box;
  padding: 92px 0 calc(18px + var(--visago-safe-bottom));
}

.search-box {
  margin: 12px 16px;
  height: 44px;
  border-radius: 12px;
  background: var(--visago-surface-soft);
  display: flex;
  align-items: center;
  padding: 0 12px;
  gap: 8px;
}

.search-icon {
  color: var(--visago-text-soft);
  font-size: 20px;
}

.search-input {
  flex: 1;
  font-size: 14px;
  height: 44px;
}

.message-list {
  background: var(--visago-surface);
  border-top: 1px solid var(--visago-line);
  border-bottom: 1px solid var(--visago-line);
}

.message-swipe {
  position: relative;
  overflow: hidden;
  border-bottom: 1px solid var(--visago-line);
}

.message-swipe:last-child {
  border-bottom: none;
}

.message-delete {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 74px;
  background: #e74d46;
  color: #fff;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.message-card {
  position: relative;
  background: var(--visago-surface);
  transition: transform 0.22s ease;
}

.message-card--dragging {
  transition: none;
}

.message-row {
  position: relative;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
}

.unread-dot {
  position: absolute;
  left: 6px;
  top: 26px;
  width: 8px;
  height: 8px;
  border-radius: 9999px;
  background: var(--visago-primary);
}

.avatar-wrap {
  width: 46px;
  height: 46px;
  border-radius: 9999px;
  overflow: hidden;
  background: var(--visago-surface-soft);
  border: 1px solid var(--visago-line);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.avatar-wrap--icon {
  background: rgba(76, 74, 202, 0.16);
  color: var(--visago-secondary);
}

.avatar {
  width: 100%;
  height: 100%;
}

.avatar-icon {
  font-size: 22px;
}

.message-main {
  flex: 1;
  min-width: 0;
}

.title-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 8px;
}

.msg-title {
  font-size: 16px;
  line-height: 1.3;
  font-weight: 600;
  color: var(--visago-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.msg-time {
  flex-shrink: 0;
  font-size: 11px;
  color: var(--visago-text-soft);
}

.msg-time--unread {
  color: var(--visago-primary);
}

.msg-preview {
  margin-top: 2px;
  display: block;
  font-size: 13px;
  color: var(--visago-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.expand-icon {
  color: var(--visago-text-soft);
  font-size: 18px;
  transition: transform 0.2s ease;
}

.expand-icon--open {
  transform: rotate(180deg);
}

.message-detail {
  padding: 2px 40px 12px 72px;
}

.message-detail text {
  display: block;
  font-size: 14px;
  line-height: 1.45;
  color: var(--visago-text-muted);
}

.empty-card {
  margin: 12px 16px;
  padding: 18px;
  text-align: center;
  color: var(--visago-text-soft);
  font-size: 14px;
}
</style>
