<template>
  <AppLayout>
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-white">Semua Perjalanan Dinas</h1>
      <p class="text-slate-400 text-sm mt-1">Monitoring seluruh data perjalanan dinas karyawan</p>
    </div>

    <!-- Filter Bar -->
    <div class="flex flex-wrap gap-3 mb-5">
      <input v-model="search" @input="onSearch" type="text" placeholder="Cari tujuan, nama, nomor surat..."
        class="bg-slate-900 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-blue-500 transition w-full md:w-72" />

      <select v-model="filterStatus" @change="fetchTrips(1)"
        class="bg-slate-900 border border-slate-700 text-slate-300 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-blue-500 transition">
        <option value="">Semua Status</option>
        <option value="PENDING">Menunggu</option>
        <option value="APPROVED">Disetujui</option>
        <option value="ON_DUTY">Dalam Perjalanan</option>
        <option value="COMPLETED">Selesai</option>
        <option value="REJECTED">Ditolak</option>
      </select>
    </div>

    <!-- Table -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <div v-if="loading" class="p-6 space-y-3">
        <div v-for="i in 8" :key="i" class="h-14 bg-slate-800 rounded-xl animate-pulse" />
      </div>

      <div v-else-if="trips.length === 0" class="p-16 text-center">
        <p class="text-slate-500">Tidak ada data perjalanan dinas</p>
      </div>

      <table v-else class="w-full text-sm">
        <thead class="bg-slate-800/60 border-b border-slate-800">
          <tr>
            <th class="text-left text-slate-400 font-medium px-5 py-3">Karyawan</th>
            <th class="text-left text-slate-400 font-medium px-5 py-3">Tujuan</th>
            <th class="text-left text-slate-400 font-medium px-5 py-3">Periode</th>
            <th class="text-left text-slate-400 font-medium px-5 py-3">Status</th>
            <th class="text-left text-slate-400 font-medium px-5 py-3">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="trip in trips" :key="trip.id" class="hover:bg-slate-800/40 transition">
            <td class="px-5 py-4">
              <p class="text-white font-medium">{{ trip.user?.name || 'N/A' }}</p>
              <p class="text-slate-400 text-xs font-mono">{{ trip.nomor_surat }}</p>
            </td>
            <td class="px-5 py-4 text-slate-300">{{ trip.destination }}</td>
            <td class="px-5 py-4 text-slate-400 text-xs">
              {{ formatDate(trip.start_date) }} –<br />{{ formatDate(trip.end_date) }}
            </td>
            <td class="px-5 py-4"><StatusBadge :status="trip.status" /></td>
            <td class="px-5 py-4">
              <div class="flex gap-3">
                <router-link :to="`/finance/trips/${trip.id}`"
                  class="text-blue-400 hover:text-blue-300 text-xs font-medium transition">
                  Review Finansial
                </router-link>
                <router-link :to="`/finance/trips/${trip.id}/claims`"
                  class="text-teal-400 hover:text-teal-300 text-xs font-medium transition">
                  Klaim
                </router-link>
              </div>
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
            class="px-3 py-1.5 text-sm bg-slate-800 text-slate-300 rounded-lg disabled:opacity-40 hover:bg-slate-700 transition">
            ← Prev
          </button>
          <button @click="changePage(pagination.current_page + 1)" :disabled="pagination.current_page >= pagination.total_pages"
            class="px-3 py-1.5 text-sm bg-slate-800 text-slate-300 rounded-lg disabled:opacity-40 hover:bg-slate-700 transition">
            Next →
          </button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref } from 'vue'
import { onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import { tripsApi } from '@/api/trips'

const trips = ref([])
const pagination = ref({})
const loading = ref(true)
const search = ref('')
const filterStatus = ref('')
let searchTimer = null

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

async function fetchTrips(page = 1) {
  loading.value = true
  try {
    const params = { page, limit: 10, search: search.value }
    if (filterStatus.value) params.status = filterStatus.value
    const res = await tripsApi.getAll(params)
    trips.value = res.data.data?.items || []
    pagination.value = res.data.data?.pagination || {}
  } catch (e) { console.error(e) }
  finally { loading.value = false }
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
