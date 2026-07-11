<template>
  <AppLayout>
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-white">Dashboard</h1>
      <p class="text-slate-400 text-sm mt-1">Selamat datang kembali, {{ authStore.userName }}</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div v-for="i in 4" :key="i" class="bg-slate-900 border border-slate-800 rounded-2xl p-5 animate-pulse h-28" />
    </div>

    <!-- Stats Cards -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <StatCard label="Total Perjalanan" :value="stats.total_trips ?? 0" color="blue"
        icon="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
      <StatCard label="Menunggu" :value="stats.pending ?? 0" color="yellow"
        icon="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      <StatCard label="Disetujui" :value="stats.approved ?? 0" color="green"
        icon="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
      <StatCard label="Selesai" :value="stats.completed ?? 0" color="purple"
        icon="M5 13l4 4L19 7" />
    </div>

    <!-- Recent Trips -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <div class="flex items-center justify-between px-6 py-4 border-b border-slate-800">
        <h2 class="text-white font-semibold">Perjalanan Terbaru</h2>
        <router-link to="/trips" class="text-blue-400 hover:text-blue-300 text-sm transition">Lihat semua →</router-link>
      </div>

      <div v-if="loading" class="p-6 space-y-3">
        <div v-for="i in 3" :key="i" class="h-14 bg-slate-800 rounded-xl animate-pulse" />
      </div>

      <div v-else-if="recentTrips.length === 0" class="p-12 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 text-slate-700 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
        <p class="text-slate-500">Belum ada perjalanan dinas</p>
        <router-link to="/trips/create" class="inline-block mt-3 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium px-4 py-2 rounded-xl transition">
          Ajukan Pertama
        </router-link>
      </div>

      <div v-else class="divide-y divide-slate-800">
        <div v-for="trip in recentTrips" :key="trip.id"
          class="flex items-center justify-between px-6 py-4 hover:bg-slate-800/50 transition cursor-pointer"
          @click="$router.push(`/trips/${trip.id}`)">
          <div class="flex-1 min-w-0">
            <p class="text-white font-medium truncate">{{ trip.destination }}</p>
            <p class="text-slate-400 text-sm">{{ trip.nomor_surat }} · {{ formatDate(trip.start_date) }}</p>
          </div>
          <StatusBadge :status="trip.status" />
        </div>
      </div>
    </div>

    <!-- Quick Action -->
    <div class="mt-4 flex justify-end">
      <router-link to="/trips/create"
        class="inline-flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white font-semibold px-5 py-2.5 rounded-xl text-sm transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Ajukan Perjalanan Baru
      </router-link>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import StatCard from '@/components/ui/StatCard.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import { useAuthStore } from '@/stores/auth'
import { tripsApi } from '@/api/trips'

const authStore = useAuthStore()
const loading = ref(true)
const stats = ref({})
const recentTrips = ref([])

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

onMounted(async () => {
  try {
    const [dashRes, tripsRes] = await Promise.all([
      tripsApi.getEmployeeDashboard(),
      tripsApi.getMyTrips({ page: 1, limit: 5 }),
    ])
    const rawStats = dashRes.data.data?.summary || {}
    stats.value = {
      total_trips: rawStats.total ?? 0,
      pending: rawStats.pending ?? 0,
      approved: rawStats.approved ?? 0,
      completed: rawStats.completed ?? 0,
    }
    recentTrips.value = tripsRes.data.data?.items || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>
