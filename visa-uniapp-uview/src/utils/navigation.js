export const MAIN_PAGE_PADDING_BOTTOM = '112px'

export function pageBottomStyle(withTabbar = false) {
  return withTabbar
    ? `padding-bottom: calc(${MAIN_PAGE_PADDING_BOTTOM} + var(--uview-safe-bottom));`
    : 'padding-bottom: 0;'
}
