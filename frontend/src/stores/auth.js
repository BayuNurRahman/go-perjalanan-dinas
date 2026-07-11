import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || null)
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)
  const role = computed(() => user.value?.role || null)
  const userName = computed(() => user.value?.name || '')
  const departmentCode = computed(() => user.value?.department?.code || '')
  const departmentName = computed(() => user.value?.department?.name || '')

  const isFinanceStaff = computed(() => {
    if (role.value === 'ADMIN_FIN') return true
    if (role.value === 'MANAGER') {
      const code = departmentCode.value.toUpperCase()
      const name = departmentName.value.toLowerCase()
      return code === 'FIN' || name.includes('finance') || name.includes('keuangan')
    }
    return false
  })

  function setAuth(data) {
    token.value = data.token
    user.value = { name: data.name, role: data.role, department: data.department }
    localStorage.setItem('token', data.token)
    localStorage.setItem('user', JSON.stringify(user.value))
  }

  function clearAuth() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  // Redirect path berdasarkan role setelah login
  function getDashboardPath() {
    switch (role.value) {
      case 'MANAGER':    return '/manager/dashboard'
      case 'ADMIN_FIN':  return '/finance/trips'
      case 'ADMIN_HR':   return '/admin/users'
      case 'ADMIN_IT':   return '/admin/departments'
      case 'SUPER_ADMIN': return '/admin/users'
      default:           return '/dashboard'  // EMPLOYEE
    }
  }

  return { token, user, isLoggedIn, role, userName, departmentCode, departmentName, isFinanceStaff, setAuth, clearAuth, getDashboardPath }
})
