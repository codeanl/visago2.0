<script>
import { applyTheme, getStoredTheme } from './utils/theme'
import { ensureAuthByRoute } from './utils/auth'

export default {
  onLaunch() {
    applyTheme(getStoredTheme())
    this.enforceAuth()
    console.log('App Launch')
  },
  onShow() {
    applyTheme(getStoredTheme())
    this.enforceAuth()
    console.log('App Show')
  },
  onHide() {
    console.log('App Hide')
  },
  methods: {
    enforceAuth() {
      if (typeof getCurrentPages !== 'function') return
      const pages = getCurrentPages()
      const current = pages && pages.length ? pages[pages.length - 1] : null
      if (!current || !current.route) return
      ensureAuthByRoute(current.route)
    },
  },
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');
@import url('https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap');

:root {
  --visago-primary: #0f65d8;
  --visago-primary-strong: #0857bf;
  --visago-secondary: #4c4aca;
  --visago-tertiary: #0f9735;
  --visago-danger: #dc2626;

  --visago-bg: #f3f2f8;
  --visago-surface: #ffffff;
  --visago-surface-soft: #ecebf2;
  --visago-surface-softest: #f7f6fb;
  --visago-line: #dde1ea;

  --visago-text: #1a1b1f;
  --visago-text-muted: #667085;
  --visago-text-soft: #8b93a9;

  --visago-radius-sm: 10px;
  --visago-radius-md: 14px;
  --visago-radius-lg: 18px;
  --visago-radius-pill: 9999px;

  --visago-shadow-card: 0 6px 18px rgba(21, 31, 52, 0.06);
  --visago-shadow-fab: 0 10px 24px rgba(15, 101, 216, 0.3);

  --visago-space-1: 4px;
  --visago-space-2: 8px;
  --visago-space-3: 12px;
  --visago-space-4: 16px;
  --visago-space-5: 24px;
  --visago-space-6: 32px;

  --visago-safe-bottom: env(safe-area-inset-bottom, 0px);
}

html.theme-dark {
  --visago-bg: #080a0f;
  --visago-surface: #11141b;
  --visago-surface-soft: #171b24;
  --visago-surface-softest: #202634;
  --visago-line: #252b38;
  --visago-text: #eef1f8;
  --visago-text-muted: #b1b9ce;
  --visago-text-soft: #8a94ad;
  --visago-shadow-card: 0 8px 20px rgba(0, 0, 0, 0.28);
}

page {
  background: var(--visago-bg);
  color: var(--visago-text);
  font-family: Inter, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
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

.visago-page-width {
  width: 100%;
  max-width: 430px;
  margin: 0 auto;
}

.visago-card {
  background: var(--visago-surface);
  border-radius: var(--visago-radius-md);
  box-shadow: var(--visago-shadow-card);
  border: 1px solid var(--visago-line);
}

.visago-chip {
  padding: 10px 20px;
  border-radius: var(--visago-radius-pill);
  background: var(--visago-surface-soft);
  color: var(--visago-text-muted);
  font-size: 18px;
  line-height: 1;
}

.visago-chip--active {
  background: var(--visago-primary);
  color: #fff;
}

.visago-divider {
  height: 1px;
  background: var(--visago-line);
}

.visago-ellipsis-2 {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
</style>

