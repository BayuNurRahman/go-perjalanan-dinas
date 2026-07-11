import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useLoadingStore } from '@/stores/loading'
import { useToast } from '@/composables/useToast'
import router from '@/router'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor: tambahkan JWT token ke setiap request
api.interceptors.request.use(
  (config) => {
    // Start global loading
    try {
      const loadingStore = useLoadingStore()
      loadingStore.startLoading()
    } catch (e) {}

    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    try {
      const loadingStore = useLoadingStore()
      loadingStore.stopLoading()
    } catch (e) {}
    return Promise.reject(error)
  }
)

// Response interceptor: handle 401 (token expired) → auto logout & redirect login
api.interceptors.response.use(
  (response) => {
    // Stop global loading
    try {
      const loadingStore = useLoadingStore()
      loadingStore.stopLoading()
    } catch (e) {}
    return response
  },
  (error) => {
    // Stop global loading
    try {
      const loadingStore = useLoadingStore()
      loadingStore.stopLoading()
    } catch (e) {}

    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      const toast = useToast()
      authStore.clearAuth()
      toast.error('Sesi Anda telah berakhir. Silakan login kembali.')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export default api

