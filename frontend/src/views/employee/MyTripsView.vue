<template>
  <AppLayout>
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-white">Perjalanan Dinas Saya</h1>
        <p class="text-slate-400 text-sm mt-1">Daftar pengajuan perjalanan dinas Anda</p>
      </div>
      <router-link to="/trips/create"
        class="inline-flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white font-semibold px-4 py-2.5 rounded-xl text-sm transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Ajukan Perjalanan
      </router-link>
    </div>

    <!-- Search -->
    <div class="mb-4">
      <input v-model="search" @input="onSearch" type="text" placeholder="Cari tujuan atau nomor surat..."
        class="w-full md:w-80 bg-slate-900 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-blue-500 transition" />
    </div>

    <!-- Table -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <div v-if="loading" class="p-6 space-y-3">
        <div v-for="i in 5" :key="i" class="h-14 bg-slate-800 rounded-xl animate-pulse" />
      </div>

      <div v-else-if="trips.length === 0" class="p-16 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 text-slate-700 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
        <p class="text-slate-500">Belum ada perjalanan dinas</p>
      </div>

      <table v-else class="w-full text-sm">
        <thead class="bg-slate-800/60 border-b border-slate-800">
          <tr>
            <th class="text-left text-slate-400 font-medium px-6 py-3">No. Surat</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Tujuan</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Tanggal</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Status</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="trip in trips" :key="trip.id" class="hover:bg-slate-800/40 transition">
            <td class="px-6 py-4 text-slate-300 font-mono text-xs">{{ trip.nomor_surat }}</td>
            <td class="px-6 py-4">
              <p class="text-white font-medium">{{ trip.destination }}</p>
              <p class="text-slate-400 text-xs">{{ trip.initiator }}</p>
            </td>
            <td class="px-6 py-4 text-slate-300">
              {{ formatDate(trip.start_date) }} – {{ formatDate(trip.end_date) }}
            </td>
            <td class="px-6 py-4"><StatusBadge :status="trip.status" /></td>
            <td class="px-6 py-4">
              <router-link :to="`/trips/${trip.id}`"
                class="text-blue-400 hover:text-blue-300 text-xs font-medium transition">
                Detail →
              </router-link>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <div v-if="pagination.total_pages > 1" class="flex items-center justify-between px-6 py-4 border-t border-slate-800">
        <p class="text-slate-400 text-sm">
          Halaman {{ pagination.current_page }} dari {{ pagination.total_pages }}
          ({{ pagination.total_items }} data)
        </p>
        <div class="flex gap-2">
          <button @click="changePage(pagination.current_page - 1)" :disabled="pagination.current_page <= 1"
            class="px-3 py-1.5 text-sm bg-slate-800 text-slate-300 rounded-lg disabled:opacity-40 hover:bg-slate-700 transition">← Prev</button>
          <button @click="changePage(pagination.current_page + 1)" :disabled="pagination.current_page >= pagination.total_pages"
            class="px-3 py-1.5 text-sm bg-slate-800 text-slate-300 rounded-lg disabled:opacity-40 hover:bg-slate-700 transition">Next →</button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import { tripsApi } from '@/api/trips'

const trips = ref([])
const pagination = ref({})
const loading = ref(true)
const search = ref('')
let searchTimer = null

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

async function fetchTrips(page = 1) {
  loading.value = true
  try {
    const res = await tripsApi.getMyTrips({ page, limit: 10, search: search.value })
    trips.value = res.data.data?.items || []
    pagination.value = res.data.data?.pagination || {}
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function onSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => fetchTrips(1), 400)
}

function changePage(page) {
  if (page < 1 || page > pagination.value.total_pages) return
  fetchTrips(page)
}

onMounted(() => fetchTrips())
</script>
