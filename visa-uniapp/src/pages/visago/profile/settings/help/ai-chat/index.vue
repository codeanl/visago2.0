<template>
  <view class="chat-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="AI客服" />
    <view class="chat-wrap">
      <scroll-view scroll-y class="chat-scroll" :scroll-into-view="scrollIntoView">
        <view class="chat-content visago-page-width">
          <view v-for="item in messages" :key="item.id" :id="`msg-${item.id}`" class="bubble-row" :class="{ 'bubble-row--me': item.role === 'user' }">
            <view class="bubble" :class="{ 'bubble--me': item.role === 'user' }">
              <text>{{ item.text }}</text>
            </view>
          </view>
        </view>
      </scroll-view>

      <view class="composer visago-page-width">
        <input
          v-model.trim="draft"
          class="composer-input"
          placeholder="请输入问题，例如：日本旅游签需要哪些材料？"
          confirm-type="send"
          @confirm="sendMessage"
        />
        <view class="send-btn" @tap="sendMessage">发送</view>
      </view>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../../../components/VisagoTopBar.vue'
import { applyTheme, getStoredTheme } from '../../../../../../utils/theme'

let seed = 1

function nextId() {
  seed += 1
  return seed
}

export default {
  components: { VisagoTopBar },
  data() {
    return {
      draft: '',
      scrollIntoView: '',
      messages: [
        {
          id: 1,
          role: 'assistant',
          text: '你好，我是 Visago AI 客服。你可以问我签证材料、办理流程、时间安排或常见拒签原因。',
        },
      ],
    }
  },
  onShow() {
    applyTheme(getStoredTheme())
    this.scrollToBottom()
  },
  methods: {
    sendMessage() {
      if (!this.draft) return
      const text = this.draft
      this.draft = ''
      const userMsg = { id: nextId(), role: 'user', text }
      this.messages.push(userMsg)
      this.scrollToBottom()

      const reply = { id: nextId(), role: 'assistant', text: this.buildReply(text) }
      setTimeout(() => {
        this.messages.push(reply)
        this.scrollToBottom()
      }, 300)
    },
    buildReply(text) {
      const t = text.toLowerCase()
      if (t.includes('日本') || t.includes('japan')) {
        return '日本旅游签通常需要护照、照片、行程单、在职或在读证明、资产证明。建议先在“我的资料”里按清单逐项补齐。'
      }
      if (t.includes('拒签') || t.includes('不过')) {
        return '常见拒签原因是材料不完整、行程逻辑不清、资金证明不足。你可以把当前材料情况发我，我帮你按风险点逐条排查。'
      }
      if (t.includes('多久') || t.includes('时间')) {
        return '多数旅游签从递签到结果通常在 5-15 个工作日，旺季会更久。建议至少提前 4-6 周启动准备。'
      }
      return '已收到你的问题。建议先确认目标国家、出行时间和签证类型，我可以按这三项给你一份可执行的准备清单。'
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const last = this.messages[this.messages.length - 1]
        if (!last) return
        this.scrollIntoView = `msg-${last.id}`
      })
    },
  },
}
</script>

<style scoped>
.chat-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.chat-wrap {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 0;
  left: 0;
  display: flex;
  flex-direction: column;
}

.chat-scroll {
  flex: 1;
  min-height: 0;
}

.chat-content {
  padding: 14px 16px 16px;
  box-sizing: border-box;
}

.bubble-row {
  display: flex;
  margin-top: 10px;
}

.bubble-row--me {
  justify-content: flex-end;
}

.bubble {
  max-width: 78%;
  padding: 10px 12px;
  border-radius: 14px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  color: var(--visago-text);
  font-size: 14px;
  line-height: 1.55;
}

.bubble--me {
  background: var(--visago-primary);
  border-color: var(--visago-primary);
  color: #fff;
}

.composer {
  padding: 10px 16px calc(10px + var(--visago-safe-bottom));
  box-sizing: border-box;
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--visago-bg);
}

.composer-input {
  flex: 1;
  height: 42px;
  border-radius: 12px;
  padding: 0 12px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  font-size: 14px;
}

.send-btn {
  width: 62px;
  height: 42px;
  border-radius: 12px;
  background: var(--visago-primary);
  color: #fff;
  font-size: 14px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
