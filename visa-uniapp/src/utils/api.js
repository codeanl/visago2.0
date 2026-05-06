import { getAuthToken, logoutSession, LOGIN_PAGE, setAuthUser } from './auth'

const BASE = (import.meta.env.VITE_API_BASE || 'http://127.0.0.1:8080/api').replace(/\/+$/, '')

function withPath(path) {
  return `${BASE}${path.startsWith('/') ? path : `/${path}`}`
}

function withQuery(path, params = {}) {
  const entries = Object.entries(params).filter(([, value]) => value !== undefined && value !== null && value !== '')
  if (!entries.length) return path
  const search = entries
    .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
    .join('&')
  return `${path}${path.includes('?') ? '&' : '?'}${search}`
}

function request(path, options = {}) {
  const method = options.method || 'GET'
  const auth = options.auth !== false
  const data = options.data || null
  const header = {
    'Content-Type': 'application/json',
    ...(options.header || {}),
  }
  const token = getAuthToken()
  if (auth && token) {
    header.Authorization = `Bearer ${token}`
  }

  return new Promise((resolve, reject) => {
    uni.request({
      url: withPath(path),
      method,
      data,
      header,
      success: (res) => {
        const statusCode = Number(res.statusCode || 0)
        const body = res.data || {}
        if (statusCode >= 200 && statusCode < 300) {
          resolve(body.data)
          return
        }
        if (statusCode === 401) {
          logoutSession()
          const pages = typeof getCurrentPages === 'function' ? getCurrentPages() : []
          const current = pages.length ? `/${pages[pages.length - 1].route}` : ''
          if (current && current !== LOGIN_PAGE) {
            uni.reLaunch({ url: LOGIN_PAGE })
          }
        }
        reject(new Error(body.message || `Request failed (${statusCode})`))
      },
      fail: (err) => {
        reject(new Error((err && err.errMsg) || 'Network request failed'))
      },
    })
  })
}

function rememberUser(user) {
  if (user && typeof user === 'object') {
    setAuthUser(user)
  }
  return user
}

export const api = {
  health: () => request('/health', { auth: false }),

  register: (payload) => request('/auth/register', { method: 'POST', data: payload, auth: false }),
  login: (payload) => request('/auth/login', { method: 'POST', data: payload, auth: false }),
  resetPassword: (payload) => request('/auth/password/reset', { method: 'POST', data: payload, auth: false }),
  me: () => request('/auth/me').then(rememberUser),
  deleteMyAccount: (payload) => request('/auth/me/delete', { method: 'POST', data: payload }),
  updateProfile: (payload) => request('/auth/me/profile', { method: 'PUT', data: payload }).then(rememberUser),
  uploadAvatar: (filePath, fileObject) =>
    new Promise((resolve, reject) => {
      const token = getAuthToken()
      const uploadOptions = {
        url: withPath('/auth/me/avatar'),
        name: 'file',
        header: token ? { Authorization: `Bearer ${token}` } : {},
        success: (res) => {
          const statusCode = Number(res.statusCode || 0)
          let body = {}
          try {
            body = JSON.parse(res.data || '{}')
          } catch (error) {
            body = {}
          }
          if (statusCode >= 200 && statusCode < 300) {
            if (body.data && body.data.profile) {
              rememberUser(body.data.profile)
            }
            resolve(body.data)
            return
          }
          reject(new Error(body.message || `Upload failed (${statusCode})`))
        },
        fail: (err) => {
          reject(new Error((err && err.errMsg) || 'Upload failed'))
        },
      }
      if (filePath) {
        uploadOptions.filePath = filePath
      } else if (fileObject) {
        uploadOptions.file = fileObject
      } else {
        reject(new Error('No upload file found'))
        return
      }
      uni.uploadFile(uploadOptions)
    }),
  uploadImage: (filePath, folder = 'community') =>
    new Promise((resolve, reject) => {
      if (!filePath) {
        reject(new Error('No upload file found'))
        return
      }
      uni.uploadFile({
        url: withPath(withQuery('/uploads/image', { folder })),
        name: 'file',
        filePath,
        success: (res) => {
          const statusCode = Number(res.statusCode || 0)
          let body = {}
          try {
            body = JSON.parse(res.data || '{}')
          } catch (error) {
            body = {}
          }
          if (statusCode >= 200 && statusCode < 300) {
            resolve(body.data)
            return
          }
          reject(new Error(body.message || `Upload failed (${statusCode})`))
        },
        fail: (err) => {
          reject(new Error((err && err.errMsg) || 'Upload failed'))
        },
      })
    }),
  analyzePhotoCheck: (filePath, formData = {}) =>
    new Promise((resolve, reject) => {
      const token = getAuthToken()
      if (!filePath) {
        reject(new Error('No upload file found'))
        return
      }
      uni.uploadFile({
        url: withPath('/tools/photo-check/analyze'),
        name: 'file',
        filePath,
        formData,
        header: token ? { Authorization: `Bearer ${token}` } : {},
        success: (res) => {
          const statusCode = Number(res.statusCode || 0)
          let body = {}
          try {
            body = JSON.parse(res.data || '{}')
          } catch (error) {
            body = {}
          }
          if (statusCode >= 200 && statusCode < 300) {
            resolve(body.data)
            return
          }
          if (statusCode === 401) {
            logoutSession()
            const pages = typeof getCurrentPages === 'function' ? getCurrentPages() : []
            const current = pages.length ? `/${pages[pages.length - 1].route}` : ''
            if (current && current !== LOGIN_PAGE) {
              uni.reLaunch({ url: LOGIN_PAGE })
            }
          }
          reject(new Error(body.message || `Upload failed (${statusCode})`))
        },
        fail: (err) => {
          reject(new Error((err && err.errMsg) || 'Upload failed'))
        },
      })
    }),
  subscribeMembership: (planKey) =>
    request('/auth/me/membership/subscribe', { method: 'POST', data: { planKey } }).then(rememberUser),

  listHotDestinations: (params = {}) => request(withQuery('/visa/hot-destinations', params), { auth: false }),
  listCountries: (q, region) => request(withQuery('/visa/countries', { q, region }), { auth: false }),
  listFreeCountries: (params = {}) => request(withQuery('/visa/free-countries', params), { auth: false }),
  listVisasByCountry: (countryId, q) => request(withQuery(`/visa/countries/${countryId}/visas`, { q }), { auth: false }),
  getVisaDetail: (visaId) => request(`/visa/visas/${visaId}`, { auth: false }),
  getPhotoCheckQuota: () => request('/tools/photo-check/quota'),
  translateText: (payload) => request('/tools/translate/text', { method: 'POST', data: payload, auth: false }),
  listExchangeCountries: (params = {}) => request(withQuery('/tools/exchange/countries', params), { auth: false }),
  getExchangeQuote: (params = {}) => request(withQuery('/tools/exchange/quote', params), { auth: false }),
  getExchangeTrend: (params = {}) => request(withQuery('/tools/exchange/trend', params), { auth: false }),
  listEmbassies: (params = {}) => request(withQuery('/tools/embassies', params), { auth: false }),
  listCommunityPosts: (params = {}) => request(withQuery('/community/posts', params)),
  listMyCommunityPosts: (params = {}) => request(withQuery('/community/me/posts', params)),
  listCommunityFavorites: () => request('/community/favorites'),
  favoriteCommunityPost: (postId) => request('/community/favorites', { method: 'POST', data: { postId } }),
  unfavoriteCommunityPost: (postId) => request(`/community/favorites/${postId}`, { method: 'DELETE' }),
  getCommunityPost: (postId) => request(`/community/posts/${postId}`),
  createCommunityPost: (payload) => request('/community/posts', { method: 'POST', data: payload }),
  deleteCommunityPost: (postId) => request(`/community/posts/${postId}`, { method: 'DELETE' }),
  listCommunityComments: (postId) => request(`/community/posts/${postId}/comments`),
  createCommunityComment: (postId, payload) => request(`/community/posts/${postId}/comments`, { method: 'POST', data: payload }),
  deleteCommunityComment: (commentId) => request(`/community/comments/${commentId}`, { method: 'DELETE' }),
  reportCommunityComment: (commentId, payload) => request(`/community/comments/${commentId}/report`, { method: 'POST', data: payload }),
  likeCommunityPost: (postId) => request(`/community/posts/${postId}/like`, { method: 'POST' }),
  unlikeCommunityPost: (postId) => request(`/community/posts/${postId}/like`, { method: 'DELETE' }),
  reportCommunityPost: (postId, payload) => request(`/community/posts/${postId}/report`, { method: 'POST', data: payload }),
  blockCommunityUser: (userId) => request(`/community/users/${userId}/block`, { method: 'POST' }),
  unblockCommunityUser: (userId) => request(`/community/users/${userId}/block`, { method: 'DELETE' }),

  createPlan: (payload) => request('/plans', { method: 'POST', data: payload }),
  listPlans: () => request('/plans'),
  getPlanDetail: (planId) => request(`/plans/${planId}`),
  updatePlanTask: (planId, taskId, payload) => request(`/plans/${planId}/tasks/${taskId}`, { method: 'PATCH', data: payload }),
  updatePlanResult: (planId, payload) => request(`/plans/${planId}/result`, { method: 'PATCH', data: payload }),
  deletePlan: (planId) => request(`/plans/${planId}`, { method: 'DELETE' }),
}
