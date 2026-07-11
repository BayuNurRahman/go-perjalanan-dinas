<template>
  <AppLayout>
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-white">Dashboard Manager</h1>
      <p class="text-slate-400 text-sm mt-1">Ringkasan perjalanan dinas tim Anda</p>
    </div>

    <!-- Stats -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div v-for="i in 4" :key="i" class="bg-slate-900 border border-slate-800 rounded-2xl p-5 animate-pulse h-28" />
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <StatCard label="Total Pengajuan" :value="stats.total_trips ?? 0" color="blue"
        icon="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
      <StatCard label="Menunggu Review" :value="stats.pending ?? 0" color="yellow"
        icon="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      <StatCard label="Disetujui" :value="stats.approved ?? 0" color="green"
        icon="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
      <StatCard label="Dalam Perjalanan" :value="stats.on_duty ?? 0" color="purple"
        icon="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </div>

    <!-- Quick Links -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
      <router-link to="/manager/applications"
        class="bg-slate-900 border border-slate-800 hover:border-blue-600/50 rounded-2xl p-5 flex items-center gap-4 transition group">
        <div class="w-10 h-10 bg-blue-600/20 rounded-xl flex items-center justify-center group-hover:bg-blue-600/30 transition">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2m-6 9l2 2 4-4" />
          </svg>
        </div>
        <div>
          <p class="text-white font-medium text-sm">Pengajuan Masuk</p>
          <p class="text-slate-400 text-xs">Review & approve trip</p>
        </div>
      </router-link>

      <router-link to="/manager/team"
        class="bg-slate-900 border border-slate-800 hover:border-purple-600/50 rounded-2xl p-5 flex items-center gap-4 transition group">
        <div class="w-10 h-10 bg-purple-600/20 rounded-xl flex items-center justify-center group-hover:bg-purple-600/30 transition">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </div>
        <div>
          <p class="text-white font-medium text-sm">Distribusi Tim</p>
          <p class="text-slate-400 text-xs">Lihat status anggota</p>
        </div>
      </router-link>

      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-5 flex items-center gap-4">
        <div class="w-10 h-10 bg-green-600/20 rounded-xl flex items-center justify-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
        </div>
        <div>
          <p class="text-white font-medium text-sm">Approval Rate</p>
          <p class="text-green-400 text-lg font-bold">{{ approvalRate }}%</p>
        </div>
      </div>
    </div>

    <!-- Pending Applications Table -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <div class="flex items-center justify-between px-6 py-4 border-b border-slate-800">
        <h2 class="text-white font-semibold">Pengajuan Menunggu Persetujuan</h2>
        <router-link to="/manager/applications" class="text-blue-400 hover:text-blue-300 text-sm transition">
          Lihat semua →
        </router-link>
      </div>

      <div v-if="loading" class="p-6 space-y-3">
        <div v-for="i in 3" :key="i" class="h-14 bg-slate-800 rounded-xl animate-pulse" />
      </div>

      <div v-else-if="pendingTrips.length === 0" class="p-10 text-center">
        <p class="text-slate-500 text-sm">Tidak ada pengajuan yang menunggu persetujuan</p>
      </div>

      <div v-else class="divide-y divide-slate-800">
        <div v-for="trip in pendingTrips" :key="trip.id"
          class="flex items-center justify-between px-6 py-4 hover:bg-slate-800/40 transition cursor-pointer"
          @click="$router.push(`/manager/trips/${trip.id}`)">
          <div class="flex-1 min-w-0">
            <p class="text-white font-medium truncate">{{ trip.destination }}</p>
            <p class="text-slate-400 text-xs">{{ trip.user?.name || 'Unknown' }} · {{ trip.nomor_surat }}</p>
          </div>
          <div class="flex items-center gap-3 ml-4">
            <span class="text-slate-400 text-xs">{{ formatDate(trip.start_date) }}</span>
            <StatusBadge :status="trip.status" />
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import StatCard from '@/components/ui/StatCard.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import { tripsApi } from '@/api/trips'

const loading = ref(true)
const stats = ref({})
const pendingTrips = ref([])

const approvalRate = computed(() => {
  const total = stats.value.total_trips || 0
  const approved = (stats.value.approved || 0) + (stats.value.completed || 0)
  if (!total) return 0
  return Math.round((approved / total) * 100)
})

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

onMounted(async () => {
  try {
    const [dashRes, appRes] = await Promise.all([
      tripsApi.getManagerDashboard(),
      tripsApi.getIncomingApplications(),
    ])
    const rawStats = dashRes.data.data?.summary || {}
    stats.value = {
      total_trips: rawStats.total ?? 0,
      pending: rawStats.pending ?? 0,
      approved: rawStats.approved ?? 0,
      on_duty: rawStats.on_duty ?? 0,
      completed: rawStats.completed ?? 0,
    }
    const all = appRes.data.data || []
    pendingTrips.value = all.filter(t => t.status === 'PENDING').slice(0, 5)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>
