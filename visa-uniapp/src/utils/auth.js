export const AUTH_TOKEN_KEY = 'visago_auth_token'
export const AUTH_USER_KEY = 'visago_auth_user'

export const LOGIN_PAGE = '/pages/visago/auth/login/index'
export const REGISTER_PAGE = '/pages/visago/auth/register/index'
export const RESET_PASSWORD_PAGE = '/pages/visago/auth/forgot-password/index'
export const AGREEMENT_PAGE = '/pages/visago/profile/settings/about/agreement/index'
export const PRIVACY_POLICY_PAGE = '/pages/visago/profile/settings/about/privacy-policy/index'
export const HOME_PAGE = '/pages/visago/home/index'

export function getAuthToken() {
  try {
    return uni.getStorageSync(AUTH_TOKEN_KEY) || ''
  } catch (error) {
    return ''
  }
}

export function setAuthToken(token) {
  const next = String(token || '').trim()
  try {
    if (next) {
      uni.setStorageSync(AUTH_TOKEN_KEY, next)
    } else {
      uni.removeStorageSync(AUTH_TOKEN_KEY)
    }
  } catch (error) {
    return ''
  }
  return next
}

export function getAuthUser() {
  try {
    const user = uni.getStorageSync(AUTH_USER_KEY)
    return user && typeof user === 'object' ? user : null
  } catch (error) {
    return null
  }
}

export function setAuthUser(user) {
  try {
    if (user && typeof user === 'object') {
      uni.setStorageSync(AUTH_USER_KEY, user)
      return user
    }
    uni.removeStorageSync(AUTH_USER_KEY)
  } catch (error) {
    return null
  }
  return null
}

export function isLoggedIn() {
  return !!getAuthToken()
}

export function setLoggedIn(loggedIn) {
  if (!loggedIn) {
    logoutSession()
    return false
  }
  return isLoggedIn()
}

export function saveLoginSession(token, user) {
  const savedToken = setAuthToken(token)
  setAuthUser(user || null)
  return !!savedToken
}

export function logoutSession() {
  setAuthToken('')
  setAuthUser(null)
}

export function ensureAuthByRoute(route) {
  const path = route ? (route.startsWith('/') ? route : `/${route}`) : ''
  const publicPages = [LOGIN_PAGE, REGISTER_PAGE, RESET_PASSWORD_PAGE, AGREEMENT_PAGE, PRIVACY_POLICY_PAGE]
  const logged = isLoggedIn()

  if (!logged && !publicPages.includes(path)) {
    uni.reLaunch({ url: LOGIN_PAGE })
    return false
  }

  if (logged && path === LOGIN_PAGE) {
    uni.reLaunch({ url: HOME_PAGE })
    return false
  }

  return true
}
