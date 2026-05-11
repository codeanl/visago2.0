export const MAIN_TABS = Object.freeze([
  {
    key: 'home',
    label: '首页',
    icon: 'home',
    activeIcon: 'home-fill',
    route: '/pages/home/index',
  },
  {
    key: 'goals',
    label: '目标',
    icon: 'bookmark',
    activeIcon: 'bookmark-fill',
    route: '/pages/goals/index',
  },
])

export function findMainTab(key) {
  return MAIN_TABS.find((item) => item.key === key) || null
}
