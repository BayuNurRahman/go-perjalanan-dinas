<template>
  <AppLayout>
    <div class="flex items-center gap-3 mb-6">
      <router-link :to="`/finance/trips/${tripId}`" class="text-slate-400 hover:text-white transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </router-link>
      <div>
        <h1 class="text-2xl font-bold text-white">Klaim Reimbursement</h1>
        <p class="text-slate-400 text-sm">Perjalanan #{{ tripId }}</p>
      </div>
    </div>

    <!-- Summary Bar -->
    <div class="grid grid-cols-3 gap-4 mb-5">
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 text-center">
        <p class="text-2xl font-bold text-white">{{ claims.length }}</p>
        <p class="text-slate-400 text-xs mt-0.5">Total Klaim</p>
      </div>
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 text-center">
        <p class="text-2xl font-bold text-green-400">{{ approvedCount }}</p>
        <p class="text-slate-400 text-xs mt-0.5">Disetujui</p>
      </div>
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 text-center">
        <p class="text-lg font-bold text-white">{{ formatCurrency(totalApproved) }}</p>
        <p class="text-slate-400 text-xs mt-0.5">Nilai Disetujui</p>
      </div>
    </div>

    <!-- Claims List -->
    <div class="space-y-3">
      <div v-if="loading" v-for="i in 4" :key="i"
        class="bg-slate-900 border border-slate-800 rounded-2xl p-5 animate-pulse h-28" />

      <div v-else-if="claims.length === 0" class="bg-slate-900 border border-slate-800 rounded-2xl p-16 text-center">
        <p class="text-slate-500">Belum ada klaim yang diajukan</p>
      </div>

      <div v-else v-for="claim in claims" :key="claim.id"
        class="bg-slate-900 border border-slate-800 rounded-2xl p-5">
        <div class="flex items-start justify-between mb-3">
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <p class="text-white font-semibold">{{ claim.title }}</p>
              <StatusBadge :status="claim.status" />
            </div>
            <p class="text-slate-400 text-sm">{{ claim.description }}</p>
            <p class="text-slate-500 text-xs mt-1">
              Tanggal transaksi: {{ formatDate(claim.transaction_date) }}
            </p>
            <div v-if="claim.rejected_reason" class="mt-2 text-red-400 text-xs">
              Alasan penolakan: {{ claim.rejected_reason }}
            </div>
          </div>
          <div class="text-right ml-4 flex-shrink-0">
            <p class="text-xl font-bold text-white">{{ formatCurrency(claim.amount) }}</p>
            <p v-if="claim.reviewed_at" class="text-slate-500 text-xs mt-1">
              Direview: {{ formatDate(claim.reviewed_at) }}
            </p>
          </div>
        </div>

        <!-- Actions row -->
        <div class="flex items-center justify-between pt-3 border-t border-slate-800">
          <!-- Download attachment -->
          <button v-if="hasAttachment(claim)"
            @click="downloadClaimFile(claim)"
            class="flex items-center gap-1.5 text-blue-400 hover:text-blue-300 text-xs font-medium transition">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            Unduh Bukti
          </button>
          <span v-else class="text-slate-600 text-xs">Tidak ada lampiran</span>

          <!-- Approve / Reject buttons (only for PENDING) -->
          <div v-if="claim.status === 'PENDING'" class="flex gap-2">
            <button @click="openReviewModal(claim, 'APPROVED')"
              :disabled="processingId === claim.id"
              class="bg-green-600/20 hover:bg-green-600/30 text-green-300 text-xs font-medium px-3 py-1.5 rounded-lg transition disabled:opacity-40">
              ✓ Setujui
            </button>
            <button @click="openReviewModal(claim, 'REJECTED')"
              :disabled="processingId === claim.id"
              class="bg-red-600/20 hover:bg-red-600/30 text-red-300 text-xs font-medium px-3 py-1.5 rounded-lg transition disabled:opacity-40">
              ✕ Tolak
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Review Modal -->
    <div v-if="modal.open" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center px-4">
      <div class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-md">
        <h3 class="text-white font-semibold mb-1">
          {{ modal.status === 'APPROVED' ? '✓ Setujui Klaim' : '✕ Tolak Klaim' }}
        </h3>
        <p class="text-slate-400 text-sm mb-4">
          <strong class="text-white">{{ modal.claim?.title }}</strong> —
          {{ formatCurrency(modal.claim?.amount) }}
        </p>

        <div v-if="modal.status === 'REJECTED'">
          <label class="block text-sm font-medium text-slate-300 mb-2">Alasan Penolakan *</label>
          <textarea v-model="modal.reason" rows="3" placeholder="Tuliskan alasan penolakan..."
            class="w-full bg-slate-800 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-red-500 transition resize-none mb-4">
          </textarea>
        </div>

        <div class="flex gap-3">
          <button @click="modal.open = false"
            class="flex-1 bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium py-2.5 rounded-xl text-sm transition">
            Batal
          </button>
          <button @click="confirmReview"
            :disabled="modal.status === 'REJECTED' && !modal.reason.trim()"
            :class="modal.status === 'APPROVED' ? 'bg-green-600 hover:bg-green-700' : 'bg-red-600 hover:bg-red-700'"
            class="flex-1 disabled:opacity-40 text-white font-semibold py-2.5 rounded-xl text-sm transition">
            {{ modal.status === 'APPROVED' ? 'Setujui' : 'Tolak' }}
          </button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import { claimsApi } from '@/api/claims'

const route = useRoute()
const tripId = route.params.id

const loading = ref(true)
const claims = ref([])
const processingId = ref(null)

const modal = ref({ open: false, claim: null, status: '', reason: '' })

const approvedCount = computed(() => claims.value.filter(c => c.status === 'APPROVED').length)
const totalApproved = computed(() => claims.value.filter(c => c.status === 'APPROVED').reduce((s, c) => s + (Number(c.amount) || 0), 0))

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}
function formatCurrency(n) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(n)
}

function hasAttachment(claim) {
  return claim.attachment_paths && claim.attachment_paths.trim() !== ''
}

async function downloadClaimFile(claim) {
  const fname = claim.attachment_paths?.split(',')[0]?.trim()?.split('/').pop()
  if (!fname) return
  try {
    const res = await claimsApi.downloadAttachment(claim.id, fname)
    const url = URL.createObjectURL(res.data)
    const a = document.createElement('a'); a.href = url; a.download = fname; a.click()
    URL.revokeObjectURL(url)
  } catch (e) { console.error(e) }
}

function openReviewModal(claim, status) {
  modal.value = { open: true, claim, status, reason: '' }
}

async function confirmReview() {
  const { claim, status, reason } = modal.value
  modal.value.open = false
  processingId.value = claim.id
  try {
    const payload = { status }
    if (reason) payload.rejected_reason = reason
    const res = await claimsApi.review(claim.id, payload)
    const idx = claims.value.findIndex(c => c.id === claim.id)
    if (idx !== -1) claims.value[idx] = res.data.data || { ...claim, status }
  } catch (e) { console.error(e) }
  finally { processingId.value = null }
}

async function fetchClaims() {
  loading.value = true
  try {
    const res = await claimsApi.getByTripId(tripId)
    claims.value = res.data.data || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

onMounted(fetchClaims)
</script>
