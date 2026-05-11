<script>
export default {
  onLaunch() {
    console.log('VisaNow Vant Launch')
  },
}
</script>

<style>
:root {
  --vant-safe-top: env(safe-area-inset-top, 0px);
  --vant-safe-bottom: env(safe-area-inset-bottom, 0px);

  --vant-primary: #1677ff;
  --vant-success: #07c160;
  --vant-warning: #ff976a;
  --vant-danger: #ee0a24;

  --vant-bg: #f7f8fa;
  --vant-surface: #ffffff;
  --vant-surface-soft: #f2f3f5;
  --vant-border: #ebedf0;
  --vant-text: #1f2329;
  --vant-text-2: #4e5969;
  --vant-text-3: #86909c;

  --vant-radius-lg: 20px;
  --vant-radius-md: 16px;
  --vant-radius-sm: 12px;
  --vant-shadow: 0 8px 24px rgba(31, 35, 41, 0.05);
}

page {
  background: var(--vant-bg);
  color: var(--vant-text);
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Helvetica Neue', sans-serif;
}

view,
text,
input,
textarea,
scroll-view,
image {
  box-sizing: border-box;
}

.v-page-width {
  width: 100%;
  max-width: 430px;
  margin: 0 auto;
}

.v-screen {
  min-height: 100vh;
  background: var(--vant-bg);
}

.v-card {
  background: var(--vant-surface);
  border: 1px solid var(--vant-border);
  border-radius: var(--vant-radius-lg);
  box-shadow: var(--vant-shadow);
}

.v-section-title {
  font-size: 22px;
  line-height: 1.2;
  font-weight: 800;
  color: var(--vant-text);
}

.v-text-muted {
  color: var(--vant-text-3);
}

.v-pill {
  display: inline-block;
  min-height: 28px;
  padding: 0 12px;
  border-radius: 999px;
  background: var(--vant-surface-soft);
  color: var(--vant-text-2);
  font-size: 12px;
  line-height: 28px;
  font-weight: 600;
  text-align: center;
}

.v-button-primary,
.v-button-secondary {
  min-height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 18px;
  font-size: 15px;
  font-weight: 700;
  box-sizing: border-box;
}

.v-button-primary {
  background: var(--vant-primary);
  color: #fff;
}

.v-button-secondary {
  background: rgba(22, 119, 255, 0.08);
  color: var(--vant-primary);
  border: 1px solid rgba(22, 119, 255, 0.14);
}
</style>
