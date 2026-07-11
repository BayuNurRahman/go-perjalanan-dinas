import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  // Public
  { path: '/', redirect: '/login' },
  { path: '/login', component: () => import('@/views/auth/LoginView.vue'), meta: { guest: true } },

  // Employee
  { path: '/dashboard',      component: () => import('@/views/employee/DashboardView.vue'), meta: { roles: ['EMPLOYEE', 'MANAGER', 'SUPER_ADMIN'] } },
  { path: '/trips',          component: () => import('@/views/employee/MyTripsView.vue'),   meta: { roles: ['EMPLOYEE'] } },
  { path: '/trips/create',   component: () => import('@/views/employee/CreateTripView.vue'), meta: { roles: ['EMPLOYEE'] } },
  { path: '/trips/:id',      component: () => import('@/views/employee/TripDetailView.vue'), meta: { roles: ['EMPLOYEE'] } },
  { path: '/trips/:id/claim',component: () => import('@/views/employee/SubmitClaimView.vue'), meta: { roles: ['EMPLOYEE'] } },
  { path: '/claims',         component: () => import('@/views/employee/ClaimsView.vue'),      meta: { roles: ['EMPLOYEE'] } },
  { path: '/claims/create',  component: () => import('@/views/employee/SubmitClaimView.vue'), meta: { roles: ['EMPLOYEE'] } },
  { path: '/claims/:id/edit',component: () => import('@/views/employee/SubmitClaimView.vue'), meta: { roles: ['EMPLOYEE'] } },

  // Manager
  { path: '/manager/dashboard',    component: () => import('@/views/manager/DashboardView.vue'),       meta: { roles: ['MANAGER', 'SUPER_ADMIN'] } },
  { path: '/manager/applications', component: () => import('@/views/manager/ApplicationsView.vue'),    meta: { roles: ['MANAGER', 'SUPER_ADMIN'] } },
  { path: '/manager/team',         component: () => import('@/views/manager/TeamDistributionView.vue'),meta: { roles: ['MANAGER', 'SUPER_ADMIN'] } },
  { path: '/manager/trips/:id',    component: () => import('@/views/manager/TripDetailView.vue'),      meta: { roles: ['MANAGER', 'SUPER_ADMIN'] } },

  // Finance
  // Finance
  { path: '/finance/trips',          component: () => import('@/views/finance/AllTripsView.vue'),       meta: { roles: ['ADMIN_FIN', 'MANAGER'] } },
  { path: '/finance/trips/:id',      component: () => import('@/views/finance/TripFinancialView.vue'),  meta: { roles: ['ADMIN_FIN', 'MANAGER'] } },
  { path: '/finance/trips/:id/claims', component: () => import('@/views/finance/ClaimsView.vue'),      meta: { roles: ['ADMIN_FIN', 'MANAGER'] } },
  { path: '/finance/trips/:id/disburse', component: () => import('@/views/finance/DisburseView.vue'),  meta: { roles: ['ADMIN_FIN', 'MANAGER'] } },

  // Super Admin / HR Admin / IT Admin
  { path: '/admin/users',       component: () => import('@/views/superadmin/UsersView.vue'),         meta: { roles: ['SUPER_ADMIN', 'ADMIN_HR'] } },
  { path: '/admin/departments', component: () => import('@/views/superadmin/DepartmentsView.vue'),   meta: { roles: ['SUPER_ADMIN', 'ADMIN_IT'] } },
  { path: '/admin/roles',       component: () => import('@/views/superadmin/RolesView.vue'),         meta: { roles: ['SUPER_ADMIN', 'ADMIN_IT'] } },
  { path: '/admin/register',    component: () => import('@/views/superadmin/RegisterUserView.vue'),  meta: { roles: ['SUPER_ADMIN', 'ADMIN_HR'] } },

  // 404
  { path: '/:pathMatch(.*)*', redirect: '/login' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation Guard
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  // Sudah login tapi coba akses halaman guest (login) → redirect ke dashboard
  if (to.meta.guest && authStore.isLoggedIn) {
    return next(authStore.getDashboardPath())
  }

  // Belum login tapi coba akses halaman yang butuh auth → redirect ke login
  if (to.meta.roles && !authStore.isLoggedIn) {
    return next('/login')
  }

  // Sudah login tapi role tidak diizinkan → redirect ke dashboard role-nya
  if (to.meta.roles && !to.meta.roles.includes(authStore.role)) {
    return next(authStore.getDashboardPath())
  }

  // Guard khusus finance: Hanya ADMIN_FIN atau MANAGER dari departemen keuangan yang boleh masuk
  if (to.path.startsWith('/finance') && !authStore.isFinanceStaff) {
    return next(authStore.getDashboardPath())
  }

  next()
})

export default router
