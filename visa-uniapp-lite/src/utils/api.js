const BASE = (import.meta.env.VITE_API_BASE || 'https://visagoapi.nova2026.top/api').replace(/\/+$/, '')

function withPath(path) {
  return `${BASE}${path.startsWith('/') ? path : `/${path}`}`
}

function withQuery(path, params = {}) {
  const entries = Object.entries(params).filter(([, value]) => value !== undefined && value !== null && value !== '')
  if (!entries.length) return path
  const query = entries.map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`).join('&')
  return `${path}${path.includes('?') ? '&' : '?'}${query}`
}

function request(path) {
  return new Promise((resolve, reject) => {
    uni.request({
      url: withPath(path),
      method: 'GET',
      success: (res) => {
        const statusCode = Number(res.statusCode || 0)
        const body = res.data || {}
        if (statusCode >= 200 && statusCode < 300) {
          resolve(body.data)
          return
        }
        reject(new Error(body.message || `Request failed (${statusCode})`))
      },
      fail: (err) => {
        reject(new Error((err && err.errMsg) || 'Network request failed'))
      },
    })
  })
}

export const api = {
  listHotDestinations: (params = {}) => request(withQuery('/visa/hot-destinations', params)),
  listCountries: (q = '', region = '') => request(withQuery('/visa/countries', { q, region })),
  listVisasByCountry: (countryId, q = '') => request(withQuery(`/visa/countries/${countryId}/visas`, { q })),
  getVisaDetail: (visaId) => request(`/visa/visas/${visaId}`),
}
