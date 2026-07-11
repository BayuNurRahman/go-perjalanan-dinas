<template>
  <AppLayout>
    <div class="flex items-center gap-3 mb-6">
      <router-link to="/finance/trips" class="text-slate-400 hover:text-white transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </router-link>
      <h1 class="text-2xl font-bold text-white">Review Finansial</h1>
    </div>

    <div v-if="loading" class="space-y-4">
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6 animate-pulse h-48" />
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6 animate-pulse h-40" />
    </div>

    <div v-else-if="trip" class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Trip Info -->
      <div class="lg:col-span-2 space-y-4">
        <!-- Trip Detail Card -->
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6">
          <div class="flex items-start justify-between mb-4">
            <div>
              <p class="text-slate-400 text-xs font-mono mb-1">{{ trip.nomor_surat }}</p>
              <h2 class="text-xl font-bold text-white">{{ trip.destination }}</h2>
              <p class="text-slate-400 text-sm mt-1">
                Karyawan: <span class="text-white">{{ trip.user?.name }}</span> ·
                Departemen: <span class="text-white">{{ trip.user?.department?.name || '-' }}</span>
              </p>
            </div>
            <StatusBadge :status="trip.status" />
          </div>

          <div class="grid grid-cols-2 gap-4 py-4 border-t border-slate-800">
            <InfoItem label="Tanggal Mulai" :value="formatDate(trip.start_date)" />
            <InfoItem label="Tanggal Selesai" :value="formatDate(trip.end_date)" />
            <InfoItem label="Inisiator" :value="trip.initiator" />
            <InfoItem label="Durasi" :value="duration + ' hari'" />
          </div>

          <div class="mt-4 pt-4 border-t border-slate-800">
            <p class="text-slate-400 text-xs font-medium mb-1">Deskripsi</p>
            <p class="text-slate-300 text-sm">{{ trip.description }}</p>
          </div>
        </div>

        <!-- Claims Summary -->
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-white font-semibold">Ringkasan Klaim</h3>
            <router-link :to="`/finance/trips/${tripId}/claims`"
              class="text-blue-400 hover:text-blue-300 text-sm transition">
              Lihat detail →
            </router-link>
          </div>
          <div class="grid grid-cols-3 gap-4">
            <div class="bg-slate-800 rounded-xl p-4 text-center">
              <p class="text-2xl font-bold text-white">{{ claimSummary.total }}</p>
              <p class="text-slate-400 text-xs mt-1">Total Klaim</p>
            </div>
            <div class="bg-slate-800 rounded-xl p-4 text-center">
              <p class="text-2xl font-bold text-green-400">{{ claimSummary.approved }}</p>
              <p class="text-slate-400 text-xs mt-1">Disetujui</p>
            </div>
            <div class="bg-slate-800 rounded-xl p-4 text-center">
              <p class="text-lg font-bold text-white">{{ formatCurrency(claimSummary.totalAmount) }}</p>
              <p class="text-slate-400 text-xs mt-1">Total Nilai</p>
            </div>
          </div>
        </div>

        <!-- Attachments -->
        <div v-if="attachments.length" class="bg-slate-900 border border-slate-800 rounded-2xl p-6">
          <h3 class="text-white font-semibold mb-3">Lampiran Dokumen</h3>
          <div class="space-y-2">
            <div v-for="fname in attachments" :key="fname"
              class="flex items-center justify-between bg-slate-800 rounded-xl px-4 py-3">
              <span class="text-slate-300 text-sm truncate">{{ fname }}</span>
              <button @click="downloadFile(fname)"
                class="text-blue-400 hover:text-blue-300 text-xs font-medium ml-3 flex-shrink-0 transition">
                Unduh
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Action Panel -->
      <div class="space-y-4">
        <!-- Review Financial Form -->
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6">
          <h3 class="text-white font-semibold mb-4">Review Keuangan</h3>

          <div class="space-y-2 mb-4">
            <button v-for="opt in reviewOptions" :key="opt.value"
              @click="reviewForm.status = opt.value"
              :class="[
                reviewForm.status === opt.value ? opt.activeClass : 'bg-slate-800 text-slate-300 hover:bg-slate-700',
                'w-full text-left px-4 py-3 rounded-xl text-sm font-medium transition flex items-center gap-3'
              ]">
              <span :class="opt.dotClass" class="w-2 h-2 rounded-full flex-shrink-0"></span>
              {{ opt.label }}
            </button>
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium text-slate-300 mb-2">Catatan Keuangan</label>
            <textarea v-model="reviewForm.notes" rows="3" placeholder="Catatan atau keterangan finansial..."
              class="w-full bg-slate-800 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500 transition resize-none">
            </textarea>
          </div>

          <button @click="handleReviewFinancial"
            :disabled="!reviewForm.status || processing"
            class="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-40 text-white font-semibold py-3 rounded-xl text-sm transition flex items-center justify-center gap-2">
            <svg v-if="processing" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            {{ processing ? 'Memproses...' : 'Simpan Review' }}
          </button>

          <div v-if="successMsg" class="mt-3 bg-green-900/40 border border-green-700 text-green-300 rounded-xl px-4 py-2.5 text-sm">{{ successMsg }}</div>
          <div v-if="errorMsg" class="mt-3 bg-red-900/40 border border-red-700 text-red-300 rounded-xl px-4 py-2.5 text-sm">{{ errorMsg }}</div>
        </div>

        <!-- Disburse Button -->
        <router-link :to="`/finance/trips/${tripId}/disburse`"
          class="block bg-teal-600/10 border border-teal-600/30 hover:border-teal-500/50 rounded-2xl p-5 text-center transition group">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-teal-400 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
          <p class="text-teal-300 font-semibold text-sm">Lakukan Pencairan Dana</p>
          <p class="text-teal-400/60 text-xs mt-1">Proses pembayaran perjalanan dinas</p>
        </router-link>

        <!-- Klaim Link -->
        <router-link :to="`/finance/trips/${tripId}/claims`"
          class="block bg-slate-900 border border-slate-800 hover:border-slate-600 rounded-2xl p-5 text-center transition">
          <p class="text-white font-medium text-sm">Lihat & Review Klaim</p>
          <p class="text-slate-400 text-xs mt-1">{{ claimSummary.total }} klaim diajukan</p>
        </router-link>
      </div>
    </div>

    <div v-else class="text-slate-500 text-center py-16">Data tidak ditemukan.</div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import InfoItem from '@/components/ui/InfoItem.vue'
import { tripsApi } from '@/api/trips'
import { claimsApi } from '@/api/claims'

const route = useRoute()
const tripId = route.params.id

const loading = ref(true)
const processing = ref(false)
const trip = ref(null)
const claims = ref([])
const successMsg = ref('')
const errorMsg = ref('')

const reviewForm = ref({ status: '', notes: '' })

const reviewOptions = [
  { value: 'APPROVED', label: 'Setujui Finansial', activeClass: 'bg-green-600/20 text-green-300 border border-green-600/40', dotClass: 'bg-green-400' },
  { value: 'REJECTED', label: 'Tolak Finansial',   activeClass: 'bg-red-600/20 text-red-300 border border-red-600/40',   dotClass: 'bg-red-400' },
]

const attachments = computed(() => {
  if (!trip.value?.attachment_paths) return []
  return trip.value.attachment_paths.split(',').map(p => p.trim().split('/').pop()).filter(Boolean)
})

const duration = computed(() => {
  if (!trip.value) return 0
  const d = Math.ceil((new Date(trip.value.end_date) - new Date(trip.value.start_date)) / (1000 * 60 * 60 * 24))
  return Math.max(1, d + 1)
})

const claimSummary = computed(() => ({
  total: claims.value.length,
  approved: claims.value.filter(c => c.status === 'APPROVED').length,
  totalAmount: claims.value.reduce((s, c) => s + (Number(c.amount) || 0), 0),
}))

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

function formatCurrency(n) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(n)
}

async function downloadFile(filename) {
  try {
    const res = await tripsApi.downloadAttachment(tripId, filename)
    const url = URL.createObjectURL(res.data)
    const a = document.createElement('a'); a.href = url; a.download = filename; a.click()
    URL.revokeObjectURL(url)
  } catch (e) { console.error(e) }
}

async function handleReviewFinancial() {
  processing.value = true
  successMsg.value = ''
  errorMsg.value = ''
  try {
    const res = await tripsApi.reviewFinancial(tripId, reviewForm.value)
    trip.value = res.data.data || trip.value
    successMsg.value = 'Review keuangan berhasil disimpan!'
    reviewForm.value = { status: '', notes: '' }
  } catch (e) {
    errorMsg.value = e.response?.data?.message || 'Gagal menyimpan review.'
  } finally {
    processing.value = false
  }
}

onMounted(async () => {
  try {
    const [tripsRes, claimsRes] = await Promise.all([
      tripsApi.getAll({ page: 1, limit: 100 }),
      claimsApi.getByTripId(tripId),
    ])
    const all = tripsRes.data.data?.items || []
    trip.value = all.find(t => String(t.id) === String(tripId)) || null
    claims.value = claimsRes.data.data || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
})
</script>
