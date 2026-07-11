import api from '@/utils/axios'

export const claimsApi = {
  submit: (formData) => api.post('/claims/', formData, { headers: { 'Content-Type': 'multipart/form-data' } }),
  getByTripId: (tripId) => api.get(`/claims/trip/${tripId}`),
  review: (id, data) => api.patch(`/claims/${id}/review`, data),
  downloadAttachment: (id, filename) => api.get(`/claims/${id}/files/${filename}`, { responseType: 'blob' }),
  getById: (id) => api.get(`/claims/${id}`),
  update: (id, formData) => api.put(`/claims/${id}`, formData, { headers: { 'Content-Type': 'multipart/form-data' } }),
  delete: (id) => api.delete(`/claims/${id}`),
}
