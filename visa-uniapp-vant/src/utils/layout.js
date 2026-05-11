const DEFAULT_HEADER_HEIGHT = 88
const DEFAULT_HOME_TOP = 108
const DEFAULT_SUB_TOP = 96

function safeNumber(value, fallback = 0) {
  const n = Number(value)
  return Number.isFinite(n) ? n : fallback
}

export function getLayoutMetrics() {
  try {
    if (typeof uni === 'undefined' || typeof uni.getSystemInfoSync !== 'function') {
      return {
        enabled: false,
        headerHeight: DEFAULT_HEADER_HEIGHT,
        homeTop: DEFAULT_HOME_TOP,
        subTop: DEFAULT_SUB_TOP,
        safeTop: 0,
        safeBottom: 0,
      }
    }

    const systemInfo = uni.getSystemInfoSync() || {}
    const statusBarHeight = safeNumber(systemInfo.statusBarHeight)
    const safeAreaInsets = systemInfo.safeAreaInsets || {}
    const menuRect = typeof uni.getMenuButtonBoundingClientRect === 'function' ? uni.getMenuButtonBoundingClientRect() : null
    const menuTop = safeNumber(menuRect && menuRect.top)
    const menuHeight = safeNumber(menuRect && menuRect.height, 32)

    if (!statusBarHeight || !menuTop) {
      return {
        enabled: false,
        headerHeight: DEFAULT_HEADER_HEIGHT,
        homeTop: DEFAULT_HOME_TOP,
        subTop: DEFAULT_SUB_TOP,
        safeTop: 0,
        safeBottom: safeNumber(safeAreaInsets.bottom),
      }
    }

    const navGap = Math.max(menuTop - statusBarHeight, 6)
    const navHeight = menuHeight + navGap * 2
    const headerHeight = Math.ceil(statusBarHeight + navHeight)

    return {
      enabled: true,
      safeTop: statusBarHeight,
      safeBottom: safeNumber(safeAreaInsets.bottom),
      headerHeight,
      homeTop: headerHeight + 18,
      subTop: headerHeight + 12,
    }
  } catch (error) {
    return {
      enabled: false,
      headerHeight: DEFAULT_HEADER_HEIGHT,
      homeTop: DEFAULT_HOME_TOP,
      subTop: DEFAULT_SUB_TOP,
      safeTop: 0,
      safeBottom: 0,
    }
  }
}

export function getTopSpacerStyle(kind = 'sub') {
  const metrics = getLayoutMetrics()
  const value = kind === 'home' ? metrics.homeTop : metrics.subTop
  return `height:${value}px;`
}
