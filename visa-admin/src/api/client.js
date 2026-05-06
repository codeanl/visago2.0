import { clearAdminSession, getAdminToken } from '../auth'

const API_BASE = import.meta.env.VITE_API_BASE || '/api'

function buildHeaders(options = {}) {
  const headers = new Headers(options.headers || {})
  const token = getAdminToken()
  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`)
  }
  if (!(options.body instanceof FormData) && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }
  return headers
}

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: buildHeaders(options),
    ...options,
  })
  const payload = await response.json().catch(() => ({}))
  if (!response.ok) {
    if (response.status === 401 || response.status === 403) {
      clearAdminSession()
    }
    throw new Error(payload.message || 'Request failed')
  }
  return payload.data
}

function withQuery(path, params = {}) {
  const search = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      search.set(key, value)
    }
  })
  const qs = search.toString()
  return qs ? `${path}?${qs}` : path
}

export const api = {
  health: () => request('/health'),
  loginAdmin: (data) => request('/admin/login', { method: 'POST', body: JSON.stringify(data) }),
  getAdminMe: () => request('/admin/me'),
  listAdminAccounts: (q) => request(withQuery('/admin/accounts', { q })),
  createAdminAccount: (data) => request('/admin/accounts', { method: 'POST', body: JSON.stringify(data) }),
  updateAdminAccount: (id, data) => request(`/admin/accounts/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteAdminAccount: (id) => request(`/admin/accounts/${id}`, { method: 'DELETE' }),

  listUsers: (q) => request(withQuery('/users', { q })),
  createUser: (data) => request('/users', { method: 'POST', body: JSON.stringify(data) }),
  updateUser: (id, data) => request(`/users/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteUser: (id) => request(`/users/${id}`, { method: 'DELETE' }),
  uploadUserAvatar: (id, file) => {
    const body = new FormData()
    body.append('file', file)
    return request(`/users/${id}/avatar`, { method: 'POST', body })
  },
  uploadImage: (file, folder = 'guides') => {
    const body = new FormData()
    body.append('file', file)
    return request(withQuery('/uploads/image', { folder }), { method: 'POST', body })
  },
  updateUserMembership: (id, data) => request(`/users/${id}/membership`, { method: 'PUT', body: JSON.stringify(data) }),

  listCountries: (q) => request(withQuery('/visa/countries', { q })),
  createCountry: (data) => request('/visa/countries', { method: 'POST', body: JSON.stringify(data) }),
  updateCountry: (id, data) => request(`/visa/countries/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteCountry: (id) => request(`/visa/countries/${id}`, { method: 'DELETE' }),

  listFreeCountries: (q) => request(withQuery('/visa/free-countries', { q })),
  createFreeCountry: (data) => request('/visa/free-countries', { method: 'POST', body: JSON.stringify(data) }),
  updateFreeCountry: (id, data) => request(`/visa/free-countries/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteFreeCountry: (id) => request(`/visa/free-countries/${id}`, { method: 'DELETE' }),

  listEmbassies: (params = {}) => request(withQuery('/tools/embassies', params)),
  createEmbassy: (data) => request('/tools/embassies', { method: 'POST', body: JSON.stringify(data) }),
  updateEmbassy: (id, data) => request(`/tools/embassies/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteEmbassy: (id) => request(`/tools/embassies/${id}`, { method: 'DELETE' }),

  listCommunityAdminPosts: (params = {}) => request(withQuery('/community/admin/posts', params)),
  updateCommunityAdminPost: (id, data) => request(`/community/admin/posts/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteCommunityAdminPost: (id) => request(`/community/admin/posts/${id}`, { method: 'DELETE' }),
  listCommunityAdminComments: (params = {}) => request(withQuery('/community/admin/comments', params)),
  updateCommunityAdminComment: (id, data) => request(`/community/admin/comments/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteCommunityAdminComment: (id) => request(`/community/admin/comments/${id}`, { method: 'DELETE' }),
  listCommunityAdminReports: (params = {}) => request(withQuery('/community/admin/reports', params)),
  updateCommunityAdminReport: (id, data) => request(`/community/admin/reports/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  listCommunityAdminCommentReports: (params = {}) => request(withQuery('/community/admin/comment-reports', params)),
  updateCommunityAdminCommentReport: (id, data) => request(`/community/admin/comment-reports/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  listCommunityKeywords: () => request('/community/admin/keywords'),
  createCommunityKeyword: (data) => request('/community/admin/keywords', { method: 'POST', body: JSON.stringify(data) }),
  updateCommunityKeyword: (id, data) => request(`/community/admin/keywords/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteCommunityKeyword: (id) => request(`/community/admin/keywords/${id}`, { method: 'DELETE' }),

  listCountryVisas: (params = {}) => request(withQuery('/visa/country-visas', params)),
  listCountryVisasByCountry: (countryId, q) => request(withQuery(`/visa/countries/${countryId}/visas`, { q })),
  createCountryVisaForCountry: (countryId, data) =>
    request(`/visa/countries/${countryId}/visas`, { method: 'POST', body: JSON.stringify(data) }),
  updateCountryVisa: (id, data) => request(`/visa/country-visas/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteCountryVisa: (id) => request(`/visa/country-visas/${id}`, { method: 'DELETE' }),

  getVisaDetail: (id) => request(`/visa/visas/${id}`),
  updateVisaDetail: (id, steps) => request(`/visa/visas/${id}/detail`, { method: 'PUT', body: JSON.stringify({ steps }) }),
}
