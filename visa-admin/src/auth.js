const ADMIN_TOKEN_KEY = 'visago_admin_token'
const ADMIN_USER_KEY = 'visago_admin_user'

export function getAdminToken() {
  try {
    return localStorage.getItem(ADMIN_TOKEN_KEY) || ''
  } catch (error) {
    return ''
  }
}

export function setAdminToken(token) {
  const next = String(token || '').trim()
  try {
    if (next) {
      localStorage.setItem(ADMIN_TOKEN_KEY, next)
    } else {
      localStorage.removeItem(ADMIN_TOKEN_KEY)
    }
  } catch (error) {
    return ''
  }
  return next
}

export function getAdminUser() {
  try {
    const raw = localStorage.getItem(ADMIN_USER_KEY)
    return raw ? JSON.parse(raw) : null
  } catch (error) {
    return null
  }
}

export function setAdminUser(user) {
  try {
    if (user && typeof user === 'object') {
      localStorage.setItem(ADMIN_USER_KEY, JSON.stringify(user))
      return user
    }
    localStorage.removeItem(ADMIN_USER_KEY)
  } catch (error) {
    return null
  }
  return null
}

export function saveAdminSession(token, user) {
  const savedToken = setAdminToken(token)
  setAdminUser(user || null)
  return !!savedToken
}

export function clearAdminSession() {
  setAdminToken('')
  setAdminUser(null)
}
