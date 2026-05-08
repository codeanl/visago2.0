<script>
export default {
  onLaunch() {
    console.log('Visago Lite Launch')
  },
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap');

:root {
  --visago-primary: #0f65d8;
  --visago-primary-strong: #0857bf;
  --visago-danger: #dc2626;
  --visago-bg: #f3f2f8;
  --visago-surface: #ffffff;
  --visago-surface-soft: #ecebf2;
  --visago-line: #dde1ea;
  --visago-text: #1a1b1f;
  --visago-text-muted: #667085;
  --visago-text-soft: #8b93a9;
  --visago-radius-md: 14px;
  --visago-radius-pill: 9999px;
  --visago-shadow-card: 0 6px 18px rgba(21, 31, 52, 0.06);
  --visago-shadow-fab: 0 10px 24px rgba(15, 101, 216, 0.3);
  --visago-safe-bottom: env(safe-area-inset-bottom, 0px);

  --lite-bg: #ffffff;
  --lite-surface: #ffffff;
  --lite-surface-strong: #f8fafc;
  --lite-primary: #0f5f7b;
  --lite-primary-strong: #0b4760;
  --lite-accent: #ff8a3d;
  --lite-text: #1f2a33;
  --lite-text-muted: #6a7682;
  --lite-line: #e7edf3;
  --lite-success: #168d64;
  --lite-warning: #b86814;
  --lite-danger: #c53e3e;
  --lite-shadow: 0 10px 28px rgba(37, 44, 52, 0.05);
  --lite-radius-lg: 24px;
  --lite-radius-md: 18px;
  --lite-radius-sm: 12px;
  --lite-safe-bottom: env(safe-area-inset-bottom, 0px);
}

page {
  background: var(--lite-bg);
  color: var(--lite-text);
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

.material-symbols-outlined {
  font-family: 'Material Symbols Outlined';
  font-style: normal;
  font-weight: 400;
  font-size: 24px;
  line-height: 1;
  letter-spacing: normal;
  text-transform: none;
  display: inline-block;
  white-space: nowrap;
  word-wrap: normal;
  direction: ltr;
  font-feature-settings: 'liga';
}

.page-wrap,
.visago-page-width {
  width: 100%;
  max-width: 430px;
  margin: 0 auto;
}

.screen {
  min-height: 100vh;
  padding: 0 16px calc(104px + var(--lite-safe-bottom));
}

.card {
  background: #ffffff;
  border: 1px solid var(--lite-line);
  border-radius: var(--lite-radius-lg);
  box-shadow: var(--lite-shadow);
}

.visago-card {
  background: var(--visago-surface);
  border-radius: var(--visago-radius-md);
  box-shadow: var(--visago-shadow-card);
  border: 1px solid var(--visago-line);
}

.section-title {
  font-size: 22px;
  line-height: 1.2;
  font-weight: 800;
  color: var(--lite-text);
}

.section-subtitle {
  margin-top: 6px;
  font-size: 13px;
  line-height: 1.6;
  color: var(--lite-text-muted);
}

.button-primary,
.button-secondary {
  height: 52px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 18px;
  font-size: 15px;
  font-weight: 700;
}

.button-primary {
  color: #fff;
  background: linear-gradient(135deg, var(--lite-primary) 0%, var(--lite-primary-strong) 100%);
}

.button-secondary {
  color: var(--lite-primary);
  background: rgba(15, 95, 123, 0.08);
  border: 1px solid rgba(15, 95, 123, 0.16);
}

.pill-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.pill {
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(15, 95, 123, 0.08);
  color: var(--lite-primary);
  font-size: 12px;
  line-height: 1;
}

.pill--muted {
  background: rgba(31, 42, 51, 0.06);
  color: var(--lite-text-muted);
}

.muted {
  color: var(--lite-text-muted);
}

.link-text {
  color: var(--lite-primary);
  font-weight: 700;
}

.empty-state {
  padding: 28px 20px;
  text-align: center;
}

.empty-state__emoji {
  font-size: 34px;
}

.empty-state__title {
  display: block;
  margin-top: 14px;
  font-size: 18px;
  font-weight: 800;
}

.empty-state__desc {
  display: block;
  margin-top: 8px;
  font-size: 13px;
  line-height: 1.7;
  color: var(--lite-text-muted);
}
</style>
