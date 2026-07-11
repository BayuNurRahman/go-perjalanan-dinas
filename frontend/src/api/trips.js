import api from '@/utils/axios'

export const tripsApi = {
  // Employee
  create: (formData) => api.post('/trips/', formData, { headers: { 'Content-Type': 'multipart/form-data' } }),
  getMyTrips: (params) => api.get('/trips/me', { params }),
  getById: (id) => api.get(`/trips/${id}`),
  update: (id, data) => api.put(`/trips/${id}`, data),
  remove: (id) => api.delete(`/trips/${id}`),
  updateClaim: (id, data) => api.patch(`/trips/${id}/claim`, data),
  getEmployeeDashboard: () => api.get('/trips/dashboard'),
  downloadAttachment: (id, filename) => api.get(`/trips/${id}/files/${filename}`, { responseType: 'blob' }),

  // Manager / Super Admin
  getAll: (params) => api.get('/trips', { params }),
  getManagerDashboard: () => api.get('/trips/manager/dashboard'),
  getIncomingApplications: () => api.get('/trips/manager/applications'),
  getTeamDistribution: () => api.get('/trips/manager/team-distribution'),
  updateStatus: (id, data) => {
    if (data instanceof FormData) {
      return api.patch(`/trips/${id}/status`, data, { headers: { 'Content-Type': 'multipart/form-data' } })
    }
    return api.patch(`/trips/${id}/status`, data)
  },

  // Finance
  reviewFinancial: (id, data) => api.patch(`/trips/${id}/review-financial`, data),
  disburse: (id, data) => api.patch(`/trips/${id}/disburse`, data),
}
