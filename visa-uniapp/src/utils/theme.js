const THEME_STORAGE_KEY = 'visago_theme_mode'
export const THEME_CHANGE_EVENT = 'visago_theme_change'

export function getStoredTheme() {
  try {
    const value = uni.getStorageSync(THEME_STORAGE_KEY)
    return value === 'dark' ? 'dark' : 'light'
  } catch (error) {
    return 'light'
  }
}

export function setStoredTheme(theme) {
  const nextTheme = theme === 'dark' ? 'dark' : 'light'
  uni.setStorageSync(THEME_STORAGE_KEY, nextTheme)
  return nextTheme
}

export function applyTheme(theme) {
  const normalizedTheme = theme === 'dark' ? 'dark' : 'light'
  if (typeof document !== 'undefined' && document.documentElement) {
    document.documentElement.classList.toggle('theme-dark', normalizedTheme === 'dark')
  }
  return normalizedTheme
}

export function updateTheme(theme) {
  const normalizedTheme = setStoredTheme(theme)
  const appliedTheme = applyTheme(normalizedTheme)
  if (typeof uni !== 'undefined' && uni.$emit) {
    uni.$emit(THEME_CHANGE_EVENT, appliedTheme)
  }
  return appliedTheme
}

export function toggleTheme(currentTheme) {
  const nextTheme = currentTheme === 'dark' ? 'light' : 'dark'
  return updateTheme(nextTheme)
}

