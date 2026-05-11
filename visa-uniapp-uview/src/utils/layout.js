const DEFAULT_HEADER_HEIGHT = 88
const DEFAULT_HOME_TOP = 108
const DEFAULT_SUB_TOP = 96

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

    if (!statusBarHeight || !menuTop) {
      return {
        headerHeight: DEFAULT_HEADER_HEIGHT,
        homeTop: DEFAULT_HOME_TOP,
        subTop: DEFAULT_SUB_TOP,
        safeTop: 0,
        safeBottom: toNumber(safeAreaInsets.bottom),
      }
    }

    const gap = Math.max(menuTop - statusBarHeight, 6)
    const navHeight = menuHeight + gap * 2
    const headerHeight = Math.ceil(statusBarHeight + navHeight)

    return {
      headerHeight,
      homeTop: headerHeight + 18,
      subTop: headerHeight + 12,
      safeTop: statusBarHeight,
      safeBottom: toNumber(safeAreaInsets.bottom),
    }
  } catch (error) {
    return {
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
