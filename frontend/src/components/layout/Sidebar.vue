<template>
  <div>
    <!-- Backdrop for Mobile -->
    <div
      v-if="isOpen"
      @click="$emit('close')"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm z-30 md:hidden transition-opacity"
    ></div>

    <!-- Sidebar Container -->
    <aside
      :class="[
        isOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0',
        'fixed left-0 top-0 h-full w-64 bg-slate-900 border-r border-slate-800 flex flex-col z-40 transition-transform duration-300'
      ]"
    >
      <!-- Brand -->
      <div class="h-16 flex items-center justify-between px-6 border-b border-slate-800">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064" />
            </svg>
          </div>
          <span class="text-white font-semibold text-sm">Travel Dinas</span>
        </div>

        <!-- Close Button for Mobile -->
        <button
          @click="$emit('close')"
          class="md:hidden text-slate-400 hover:text-white p-1 rounded-lg hover:bg-slate-800 transition"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Menu Navigation -->
      <nav class="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
        <template v-for="group in visibleMenuGroups" :key="group.label">
          <p class="text-slate-500 text-xs font-semibold uppercase tracking-wider px-3 mb-2 mt-4 first:mt-0">
            {{ group.label }}
          </p>
          <router-link
            v-for="item in group.items"
            :key="item.to"
            :to="item.to"
            @click="$emit('close')"
            class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm transition"
            :class="[$route.path.startsWith(item.to) ? 'bg-blue-600/20 text-blue-400' : 'text-slate-400 hover:bg-slate-800 hover:text-white']"
          >
            <component :is="item.icon" class="w-4 h-4 flex-shrink-0" />
            {{ item.label }}
          </router-link>
        </template>
      </nav>
    </aside>
  </div>
</template>

<script setup>
import { computed, defineComponent, h } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRoute } from 'vue-router'

defineProps({
  isOpen: { type: Boolean, default: false }
})

defineEmits(['close'])

const authStore = useAuthStore()
const $route = useRoute()

// Simple SVG icon factory
const Icon = (path) => defineComponent({ render: () => h('svg', { xmlns: 'http://www.w3.org/2000/svg', fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor' }, [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: path })]) })

const menuConfig = {
  EMPLOYEE: [
    { label: 'Menu', items: [
      { to: '/dashboard', label: 'Dashboard', icon: Icon('M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6') },
      { to: '/trips', label: 'Dinas Saya', icon: Icon('M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2') },
      { to: '/claims', label: 'Klaim Reimbursement', icon: Icon('M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4') },
    ]},
  ],
  MANAGER: [
    { label: 'Dashboard', items: [
      { to: '/dashboard', label: 'Dashboard Saya', icon: Icon('M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6') },
      { to: '/manager/dashboard', label: 'Dashboard Tim', icon: Icon('M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z') },
    ]},
    { label: 'Manajemen', items: [
      { to: '/manager/applications', label: 'Pengajuan Masuk', icon: Icon('M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4') },
      { to: '/manager/team', label: 'Distribusi Tim', icon: Icon('M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z') },
    ]},
  ],
  ADMIN_FIN: [
    { label: 'Keuangan', items: [
      { to: '/finance/trips', label: 'Semua Perjalanan', icon: Icon('M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2') },
    ]},
  ],
  ADMIN_HR: [
    { label: 'Administrasi', items: [
      { to: '/admin/users', label: 'Kelola User', icon: Icon('M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z') },
      { to: '/admin/register', label: 'Daftarkan User', icon: Icon('M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z') },
    ]},
  ],
  ADMIN_IT: [
    { label: 'Administrasi', items: [
      { to: '/admin/departments', label: 'Departemen', icon: Icon('M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4') },
      { to: '/admin/roles', label: 'Roles', icon: Icon('M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z') },
    ]},
  ],
  SUPER_ADMIN: [
    { label: 'Monitoring', items: [
      { to: '/dashboard', label: 'Dashboard', icon: Icon('M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6') },
      { to: '/manager/dashboard', label: 'Dashboard Manager', icon: Icon('M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z') },
    ]},
    { label: 'Administrasi', items: [
      { to: '/admin/users', label: 'Kelola User', icon: Icon('M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z') },
      { to: '/admin/departments', label: 'Departemen', icon: Icon('M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4') },
      { to: '/admin/roles', label: 'Roles', icon: Icon('M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z') },
      { to: '/admin/register', label: 'Daftarkan User', icon: Icon('M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z') },
    ]},
  ],
}

const visibleMenuGroups = computed(() => {
  const groups = [...(menuConfig[authStore.role] || [])]
  if (authStore.role === 'MANAGER' && authStore.isFinanceStaff) {
    groups.push({
      label: 'Keuangan',
      items: [
        { to: '/finance/trips', label: 'Semua Perjalanan', icon: Icon('M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2') },
      ]
    })
  }
  return groups
})
</script>
