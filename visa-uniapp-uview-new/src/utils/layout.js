const DEFAULT_HEADER_HEIGHT = 88
const DEFAULT_HOME_TOP = 108
const DEFAULT_SUB_TOP = 96
const DEFAULT_TABBAR_HEIGHT = 92

function toNumber(value, fallback = 0) {
  const n = Number(value)
  return Number.isFinite(n) ? n : fallback
}

export function getLayoutMetrics() {
  try {
    if (typeof uni === 'undefined' || typeof uni.getSystemInfoSync !== 'function') {
      return {
        headerHeight: DEFAULT_HEADER_HEIGHT,
        homeTop: DEFAULT_HOME_TOP,
        subTop: DEFAULT_SUB_TOP,
        safeTop: 0,
        safeBottom: 0,
        tabbarHeight: DEFAULT_TABBAR_HEIGHT,
        platform: 'h5',
      }
    }

    const systemInfo = uni.getSystemInfoSync() || {}
    const safeAreaInsets = systemInfo.safeAreaInsets || {}
    const statusBarHeight = toNumber(systemInfo.statusBarHeight)
    const menuRect = typeof uni.getMenuButtonBoundingClientRect === 'function'
      ? uni.getMenuButtonBoundingClientRect()
      : null
    const menuTop = toNumber(menuRect && menuRect.top)
    const menuHeight = toNumber(menuRect && menuRect.height, 32)
    const platform = String(systemInfo.uniPlatform || systemInfo.platform || '').toLowerCase()
    const safeBottom = toNumber(safeAreaInsets.bottom)

    if (!statusBarHeight || !menuTop) {
      return {
        headerHeight: DEFAULT_HEADER_HEIGHT,
        homeTop: DEFAULT_HOME_TOP,
        subTop: DEFAULT_SUB_TOP,
        safeTop: 0,
        safeBottom,
        tabbarHeight: DEFAULT_TABBAR_HEIGHT + safeBottom,
        platform,
      }
    }

    const navGap = Math.max(menuTop - statusBarHeight, 6)
    const navHeight = menuHeight + navGap * 2
    const headerHeight = Math.ceil(statusBarHeight + navHeight)

    return {
      headerHeight,
      homeTop: headerHeight + 18,
      subTop: headerHeight + 12,
      safeTop: statusBarHeight,
      safeBottom,
      tabbarHeight: DEFAULT_TABBAR_HEIGHT + safeBottom,
      platform,
    }
  } catch (error) {
    return {
      headerHeight: DEFAULT_HEADER_HEIGHT,
      homeTop: DEFAULT_HOME_TOP,
      subTop: DEFAULT_SUB_TOP,
      safeTop: 0,
      safeBottom: 0,
      tabbarHeight: DEFAULT_TABBAR_HEIGHT,
      platform: '',
    }
  }
}

export function getTopSpacerStyle(kind = 'sub') {
  const metrics = getLayoutMetrics()
  const height = kind === 'home' ? metrics.homeTop : metrics.subTop
  return `height:${height}px;`
}

export function getBottomSafeStyle(extra = 0) {
  const metrics = getLayoutMetrics()
  return `padding-bottom:${metrics.safeBottom + extra}px;`
}
