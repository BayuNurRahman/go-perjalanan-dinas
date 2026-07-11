<template>
  <AppLayout>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-white">Pengajuan Masuk</h1>
        <p class="text-slate-400 text-sm mt-1">Tinjau dan setujui perjalanan dinas karyawan</p>
      </div>
    </div>

    <!-- Filter Tabs -->
    <div class="flex gap-2 mb-5">
      <button v-for="tab in tabs" :key="tab.value"
        @click="activeTab = tab.value"
        :class="activeTab === tab.value
          ? 'bg-blue-600 text-white'
          : 'bg-slate-800 text-slate-400 hover:text-white hover:bg-slate-700'"
        class="px-4 py-2 rounded-xl text-sm font-medium transition">
        {{ tab.label }}
        <span v-if="tab.value === 'PENDING'" class="ml-1.5 bg-yellow-500/20 text-yellow-400 text-xs px-1.5 py-0.5 rounded-full">
          {{ pendingCount }}
        </span>
      </button>
    </div>

    <!-- Table -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <div v-if="loading" class="p-6 space-y-3">
        <div v-for="i in 5" :key="i" class="h-16 bg-slate-800 rounded-xl animate-pulse" />
      </div>

      <div v-else-if="filteredTrips.length === 0" class="p-16 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 text-slate-700 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-slate-500">Tidak ada pengajuan dengan status ini</p>
      </div>

      <table v-else class="w-full text-sm">
        <thead class="bg-slate-800/60 border-b border-slate-800">
          <tr>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Karyawan</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Tujuan</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Tanggal</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Status</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="trip in filteredTrips" :key="trip.id" class="hover:bg-slate-800/40 transition">
            <td class="px-6 py-4">
              <p class="text-white font-medium">{{ trip.user?.name || 'N/A' }}</p>
              <p class="text-slate-400 text-xs font-mono">{{ trip.nomor_surat }}</p>
            </td>
            <td class="px-6 py-4 text-slate-300">{{ trip.destination }}</td>
            <td class="px-6 py-4 text-slate-400 text-xs">
              {{ formatDate(trip.start_date) }} –<br>{{ formatDate(trip.end_date) }}
            </td>
            <td class="px-6 py-4"><StatusBadge :status="trip.status" /></td>
            <td class="px-6 py-4">
              <div class="flex gap-2">
                <router-link :to="`/manager/trips/${trip.id}`"
                  class="text-blue-400 hover:text-blue-300 text-xs font-medium transition">
                  Detail
                </router-link>
                <!-- Quick approve/reject only for PENDING -->
                <template v-if="trip.status === 'PENDING'">
                  <button @click="quickAction(trip.id, 'APPROVED')"
                    :disabled="processingId === trip.id"
                    class="text-green-400 hover:text-green-300 text-xs font-medium transition disabled:opacity-40">
                    Setujui
                  </button>
                  <button @click="openReject(trip)"
                    :disabled="processingId === trip.id"
                    class="text-red-400 hover:text-red-300 text-xs font-medium transition disabled:opacity-40">
                    Tolak
                  </button>
                </template>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Reject Modal -->
    <div v-if="rejectModal.open"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center px-4">
      <div class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-md">
        <h3 class="text-white font-semibold mb-1">Tolak Pengajuan</h3>
        <p class="text-slate-400 text-sm mb-4">
          Perjalanan ke <strong class="text-white">{{ rejectModal.trip?.destination }}</strong>
        </p>
        <label class="block text-sm font-medium text-slate-300 mb-2">Alasan Penolakan</label>
        <textarea v-model="rejectModal.reason" rows="3" placeholder="Berikan alasan penolakan yang jelas..."
          class="w-full bg-slate-800 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-red-500 transition resize-none mb-4">
        </textarea>
        <div class="flex gap-3">
          <button @click="rejectModal.open = false"
            class="flex-1 bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium py-2.5 rounded-xl text-sm transition">
            Batal
          </button>
          <button @click="confirmReject" :disabled="!rejectModal.reason.trim() || processingId"
            class="flex-1 bg-red-600 hover:bg-red-700 disabled:opacity-40 text-white font-semibold py-2.5 rounded-xl text-sm transition">
            Tolak Pengajuan
          </button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import { tripsApi } from '@/api/trips'

const loading = ref(true)
const trips = ref([])
const processingId = ref(null)
const activeTab = ref('PENDING')

const tabs = [
  { label: 'Menunggu', value: 'PENDING' },
  { label: 'Disetujui', value: 'APPROVED' },
  { label: 'Ditolak', value: 'REJECTED' },
  { label: 'Semua', value: 'ALL' },
]

const rejectModal = ref({ open: false, trip: null, reason: '' })

const pendingCount = computed(() => trips.value.filter(t => t.status === 'PENDING').length)
const filteredTrips = computed(() => {
  if (activeTab.value === 'ALL') return trips.value
  return trips.value.filter(t => t.status === activeTab.value)
})

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

async function fetchApplications() {
  loading.value = true
  try {
    const res = await tripsApi.getIncomingApplications()
    trips.value = res.data.data || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

async function quickAction(id, status, reason = '') {
  processingId.value = id
  try {
    await tripsApi.updateStatus(id, { status, reason })
    await fetchApplications()
  } catch (e) { console.error(e) }
  finally { processingId.value = null }
}

function openReject(trip) {
  rejectModal.value = { open: true, trip, reason: '' }
}

async function confirmReject() {
  const { trip, reason } = rejectModal.value
  rejectModal.value.open = false
  await quickAction(trip.id, 'REJECTED', reason)
}

onMounted(fetchApplications)
</script>
