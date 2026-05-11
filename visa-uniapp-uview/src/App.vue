<script>
export default {
  onLaunch() {
    console.log('VisaNow uView Launch')
  },
}
</script>

<style lang="scss">
:root {
  --uview-safe-top: env(safe-area-inset-top, 0px);
  --uview-safe-bottom: env(safe-area-inset-bottom, 0px);

  --uview-brand: #1677ff;
  --uview-brand-soft: #eff6ff;
  --uview-success: #18b566;
  --uview-warning: #ff9f43;
  --uview-danger: #ee0a24;

  --uview-bg: #ffffff;
  --uview-surface: #ffffff;
  --uview-surface-soft: #f7f8fa;
  --uview-border: #eef2f7;
  --uview-text: #1f2329;
  --uview-text-2: #4e5969;
  --uview-text-3: #86909c;

  --uview-radius-lg: 20px;
  --uview-radius-md: 16px;
  --uview-radius-sm: 12px;
  --uview-shadow: 0 6px 18px rgba(31, 35, 41, 0.04);
}

page {
  background: var(--uview-bg);
  color: var(--uview-text);
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Helvetica Neue', sans-serif;
}

view,
text,
input,
textarea,
scroll-view,
image,
swiper,
swiper-item {
  box-sizing: border-box;
}

.app-page-width {
  width: 100%;
  max-width: 430px;
  margin: 0 auto;
}

.app-screen {
  min-height: 100vh;
  background: var(--uview-bg);
}

.app-card {
  background: var(--uview-surface);
  border: 1px solid var(--uview-border);
  border-radius: var(--uview-radius-lg);
  box-shadow: var(--uview-shadow);
}

.app-section-title {
  font-size: 22px;
  line-height: 1.2;
  font-weight: 800;
  color: var(--uview-text);
}

.app-muted {
  color: var(--uview-text-3);
}

.app-link {
  color: var(--uview-brand);
  font-weight: 700;
}
</style>
