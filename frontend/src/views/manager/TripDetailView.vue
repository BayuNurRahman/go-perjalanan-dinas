<template>
  <AppLayout>
    <div class="flex items-center gap-3 mb-6">
      <router-link to="/manager/applications" class="text-slate-400 hover:text-white transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </router-link>
      <h1 class="text-2xl font-bold text-white">Detail Pengajuan</h1>
    </div>

    <div v-if="loading" class="space-y-4">
      <div class="h-40 bg-slate-900 border border-slate-800 rounded-2xl animate-pulse" />
      <div class="h-20 bg-slate-900 border border-slate-800 rounded-2xl animate-pulse" />
    </div>

    <div v-else-if="trip" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Main Info Panel (left 2/3) -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Header Summary Card -->
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6 relative">
          <div class="flex items-start justify-between mb-4">
            <div>
              <span class="text-slate-500 text-xs font-mono block mb-1">{{ trip.nomor_surat }}</span>
              <h2 class="text-2xl font-bold text-white">{{ trip.destination }}</h2>
              <p class="text-slate-400 text-sm mt-1">Diajukan oleh: <span class="text-slate-300 font-medium">{{ trip.user?.name || 'Karyawan' }}</span></p>
            </div>
            <StatusBadge :status="trip.status" />
          </div>

          <div class="grid grid-cols-2 gap-4 border-t border-slate-800/60 pt-4">
            <InfoItem label="Tanggal Mulai" :value="formatDate(trip.start_date)" />
            <InfoItem label="Tanggal Selesai" :value="formatDate(trip.end_date)" />
            <InfoItem label="Inisiator" :value="trip.initiator" />
            <InfoItem label="Tanggal Pengajuan" :value="formatDate(trip.created_at)" />
          </div>

          <div class="border-t border-slate-800/60 pt-4 mt-4">
            <p class="text-slate-500 text-xs font-medium mb-1">Deskripsi</p>
            <p class="text-slate-300 text-sm">{{ trip.description }}</p>
          </div>
          
          <div v-if="trip.summary" class="border-t border-slate-800/60 pt-4 mt-4">
            <p class="text-slate-500 text-xs font-medium mb-1">Ringkasan Perjalanan</p>
            <p class="text-slate-300 text-sm">{{ trip.summary }}</p>
          </div>

          <div v-if="trip.notes" class="border-t border-slate-800/60 pt-4 mt-4">
            <p class="text-slate-500 text-xs font-medium mb-1">Catatan Tambahan</p>
            <p class="text-slate-300 text-sm">{{ trip.notes }}</p>
          </div>
        </div>

        <!-- Attachments -->
        <div v-if="attachments.length" class="bg-slate-900 border border-slate-800 rounded-2xl p-6">
          <h3 class="text-white font-semibold mb-3">Lampiran Dokumen</h3>
          <div class="space-y-2">
            <div v-for="fname in attachments" :key="fname"
              class="flex items-center justify-between bg-slate-800 border border-slate-700/40 rounded-xl px-4 py-3">
              <span class="text-slate-300 text-sm truncate max-w-xs md:max-w-md">{{ fname }}</span>
              <div class="flex items-center gap-3 ml-3 flex-shrink-0">
                <button type="button" @click="previewFile(fname)"
                  class="text-teal-400 hover:text-teal-300 text-xs font-medium transition">
                  Pratinjau
                </button>
                <span class="text-slate-600 text-xs">|</span>
                <button type="button" @click="downloadFile(fname)"
                  class="text-blue-400 hover:text-blue-300 text-xs font-medium transition">
                  Unduh
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Claims Section -->
        <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
          <div class="px-6 py-4 border-b border-slate-800">
            <h3 class="text-white font-semibold">Klaim Reimbursement Karyawan</h3>
          </div>

          <div v-if="claims.length === 0" class="p-8 text-center text-slate-500 text-sm">
            Belum ada klaim yang diajukan untuk perjalanan ini
          </div>

          <div v-else class="divide-y divide-slate-800">
            <div v-for="claim in claims" :key="claim.id" class="px-6 py-4">
              <div class="flex items-start justify-between">
                <div>
                  <p class="text-white font-medium">{{ claim.title }}</p>
                  <p class="text-slate-400 text-sm">{{ claim.description }}</p>
                  <p class="text-slate-500 text-xs mt-1">{{ formatDate(claim.transaction_date) }}</p>
                </div>
                <div class="text-right flex-shrink-0 ml-4">
                  <p class="text-white font-semibold">{{ formatCurrency(claim.amount) }}</p>
                  <StatusBadge :status="claim.status" />
                </div>
              </div>
              
              <!-- Claim Attachments -->
              <div v-if="claim.attachment_paths || claim.attachment_path" class="mt-3 flex flex-wrap gap-2">
                <div v-for="fname in getClaimAttachments(claim)" :key="fname"
                  class="inline-flex items-center gap-2 bg-slate-800 border border-slate-700/60 rounded-lg px-2.5 py-1.5 text-xs text-slate-300">
                  <span class="truncate max-w-[150px]">{{ fname }}</span>
                  <div class="flex items-center gap-1.5 ml-1">
                    <button type="button" @click="previewClaimFile(claim.id, fname)"
                      class="text-teal-400 hover:text-teal-300 transition font-medium">
                      Pratinjau
                    </button>
                    <span class="text-slate-600">|</span>
                    <button type="button" @click="downloadClaimFile(claim.id, fname)"
                      class="text-blue-400 hover:text-blue-300 transition font-medium">
                      Unduh
                    </button>
                  </div>
                </div>
              </div>

              <div v-if="claim.rejected_reason" class="mt-2 text-red-400 text-xs">
                Alasan penolakan: {{ claim.rejected_reason }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Action Panel (right 1/3) -->
      <div class="space-y-4">
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6">
          <h3 class="text-white font-semibold mb-4">Update Status</h3>

          <!-- Status buttons -->
          <div class="space-y-2 mb-4">
            <button v-for="opt in statusOptions" :key="opt.value"
              @click="selectedStatus = opt.value"
              :class="[
                selectedStatus === opt.value ? opt.activeClass : 'bg-slate-800 text-slate-300 hover:bg-slate-700',
                'w-full text-left px-4 py-3 rounded-xl text-sm font-medium transition flex items-center gap-3'
              ]">
              <span :class="opt.dotClass" class="w-2 h-2 rounded-full flex-shrink-0"></span>
              {{ opt.label }}
            </button>
          </div>

          <!-- Reason / Catatan -->
          <div v-if="selectedStatus === 'REJECTED' || selectedStatus === 'REVISION_REQUESTED'" class="mb-4">
            <label class="block text-sm font-medium text-slate-300 mb-2">
              {{ selectedStatus === 'REJECTED' ? 'Alasan Penolakan *' : 'Catatan Revisi *' }}
            </label>
            <textarea v-model="reason" rows="3"
              :placeholder="selectedStatus === 'REJECTED' ? 'Tuliskan alasan penolakan...' : 'Tuliskan catatan revisi...'"
              class="w-full bg-slate-800 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500 transition resize-none">
            </textarea>
          </div>

          <!-- Submit -->
          <button @click="handleUpdateStatus"
            :disabled="!selectedStatus || ((selectedStatus === 'REJECTED' || selectedStatus === 'REVISION_REQUESTED') && !reason.trim()) || processing"
            class="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-40 text-white font-semibold py-3 rounded-xl text-sm transition flex items-center justify-center gap-2">
            <svg v-if="processing" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            {{ processing ? 'Memproses...' : 'Simpan Status' }}
          </button>

          <!-- Success message -->
          <div v-if="successMsg" class="mt-3 bg-green-900/40 border border-green-700 text-green-300 rounded-xl px-4 py-2.5 text-sm">
            {{ successMsg }}
          </div>
          <div v-if="errorMsg" class="mt-3 bg-red-900/40 border border-red-700 text-red-300 rounded-xl px-4 py-2.5 text-sm">
            {{ errorMsg }}
          </div>
        </div>

        <!-- Current Status Info -->
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-5">
          <p class="text-slate-400 text-xs font-medium mb-2">Status Saat Ini</p>
          <StatusBadge :status="trip.status" />
        </div>

        <!-- Delete Trip Card for Manager -->
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6">
          <h3 class="text-white font-semibold mb-2">Manajemen Perjalanan</h3>
          <p class="text-slate-400 text-xs mb-4">Sebagai manager, Anda dapat menghapus perjalanan dinas ini agar karyawan dapat mengajukan ulang.</p>
          <button @click="handleDeleteTrip" :disabled="deleteProcessing"
            class="w-full bg-red-600 hover:bg-red-700 disabled:opacity-40 text-white font-semibold py-3 rounded-xl text-sm transition">
            {{ deleteProcessing ? 'Menghapus...' : 'Hapus Perjalanan Dinas' }}
          </button>
        </div>
      </div>
    </div>

    <div v-else class="text-slate-500 text-center py-16">Data tidak ditemukan.</div>

    <!-- Preview Modal -->
    <div v-if="previewUrl" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
      <div class="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-4xl max-h-[90vh] flex flex-col overflow-hidden">
        <!-- Modal Header -->
        <div class="flex items-center justify-between px-6 py-4 border-b border-slate-800">
          <h3 class="text-white font-semibold truncate">{{ previewFilename }}</h3>
          <button @click="closePreview" class="text-slate-400 hover:text-white transition focus:outline-none">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Modal Body (Preview Content) -->
        <div class="flex-1 overflow-auto p-6 bg-slate-950/40 flex items-center justify-center min-h-[300px]">
          <!-- If Image -->
          <img v-if="isImage(previewFilename)" :src="previewUrl" class="max-w-full max-h-[70vh] object-contain rounded-lg border border-slate-800" />
          
          <!-- If PDF -->
          <iframe v-else-if="isPdf(previewFilename)" :src="previewUrl" class="w-full h-[65vh] rounded-lg border border-slate-800 bg-white" />

          <!-- Unsupported type -->
          <div v-else class="text-center text-slate-400 py-10">
            <p>Format file ini tidak mendukung pratinjau langsung.</p>
            <button @click="downloadFile(previewFilename)" class="mt-4 bg-blue-600 hover:bg-blue-700 text-white font-medium px-4 py-2 rounded-xl text-sm transition">
              Unduh File
            </button>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import InfoItem from '@/components/ui/InfoItem.vue'
import { tripsApi } from '@/api/trips'
import { claimsApi } from '@/api/claims'
import { useToast } from '@/composables/useToast'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const tripId = route.params.id

const loading = ref(true)
const processing = ref(false)
const deleteProcessing = ref(false)
const trip = ref(null)
const claims = ref([])
const selectedStatus = ref('')
const reason = ref('')
const successMsg = ref('')
const errorMsg = ref('')

const previewUrl = ref(null)
const previewFilename = ref('')

function formatCurrency(n) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(n)
}

const statusOptions = [
  { value: 'APPROVED',  label: 'Setujui Pengajuan',   activeClass: 'bg-green-600/20 text-green-300 border border-green-600/40', dotClass: 'bg-green-400' },
  { value: 'REJECTED',  label: 'Tolak Pengajuan',     activeClass: 'bg-red-600/20 text-red-300 border border-red-600/40',       dotClass: 'bg-red-400' },
  { value: 'REVISION_REQUESTED', label: 'Minta Revisi', activeClass: 'bg-orange-600/20 text-orange-300 border border-orange-600/40', dotClass: 'bg-orange-400' },
]

const attachments = computed(() => {
  if (!trip.value?.attachment_paths) return []
  return trip.value.attachment_paths.split(',').map(p => {
    const trimmed = p.trim()
    const normalized = trimmed.replace(/\\/g, '/')
    return normalized.split('/').pop()
  }).filter(Boolean)
})

function getClaimAttachments(claim) {
  const paths = claim.attachment_paths || claim.attachment_path || ''
  if (!paths) return []
  return paths.split(',').map(p => {
    const trimmed = p.trim()
    const normalized = trimmed.replace(/\\/g, '/')
    return normalized.split('/').pop()
  }).filter(Boolean)
}

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

function isImage(filename) {
  const ext = filename.toLowerCase().split('.').pop()
  return ['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext)
}

function isPdf(filename) {
  const ext = filename.toLowerCase().split('.').pop()
  return ext === 'pdf'
}

async function previewFile(filename) {
  try {
    const res = await tripsApi.downloadAttachment(tripId, filename)
    const mimeType = res.data.type
    const blob = new Blob([res.data], { type: mimeType })
    previewUrl.value = URL.createObjectURL(blob)
    previewFilename.value = filename
  } catch (e) {
    console.error('Failed to preview file', e)
  }
}

async function previewClaimFile(claimId, filename) {
  try {
    const res = await claimsApi.downloadAttachment(claimId, filename)
    const mimeType = res.data.type
    const blob = new Blob([res.data], { type: mimeType })
    previewUrl.value = URL.createObjectURL(blob)
    previewFilename.value = filename
  } catch (e) {
    console.error('Failed to preview claim file', e)
  }
}

async function downloadFile(filename) {
  try {
    const res = await tripsApi.downloadAttachment(tripId, filename)
    const url = URL.createObjectURL(res.data)
    const a = document.createElement('a'); a.href = url; a.download = filename; a.click()
    URL.revokeObjectURL(url)
  } catch (e) { console.error(e) }
}

async function downloadClaimFile(claimId, filename) {
  try {
    const res = await claimsApi.downloadAttachment(claimId, filename)
    const url = URL.createObjectURL(res.data)
    const a = document.createElement('a'); a.href = url; a.download = filename; a.click()
    URL.revokeObjectURL(url)
  } catch (e) { console.error(e) }
}

function closePreview() {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = null
  }
  previewFilename.value = ''
}

async function handleDeleteTrip() {
  if (!confirm('Apakah Anda yakin ingin menghapus pengajuan perjalanan dinas ini? Tindakan ini tidak dapat dibatalkan.')) return
  deleteProcessing.value = true
  try {
    await tripsApi.remove(tripId)
    toast.showToast('Perjalanan dinas berhasil dihapus.', 'success')
    router.push('/manager/applications')
  } catch (e) {
    toast.showToast(e.response?.data?.message || 'Gagal menghapus perjalanan dinas.', 'error')
  } finally {
    deleteProcessing.value = false
  }
}

async function handleUpdateStatus() {
  if (!selectedStatus.value) return
  processing.value = true
  successMsg.value = ''
  errorMsg.value = ''
  try {
    const sendReason = (selectedStatus.value === 'REJECTED' || selectedStatus.value === 'REVISION_REQUESTED') ? reason.value : ''
    const res = await tripsApi.updateStatus(tripId, { status: selectedStatus.value, reason: sendReason })
    trip.value = res.data.data || trip.value
    successMsg.value = 'Status berhasil diperbarui!'
    selectedStatus.value = ''
    reason.value = ''
  } catch (e) {
    errorMsg.value = e.response?.data?.message || 'Gagal memperbarui status.'
  } finally {
    processing.value = false
  }
}

onMounted(async () => {
  try {
    const [tripRes, claimsRes] = await Promise.all([
      tripsApi.getById(tripId),
      claimsApi.getByTripId(tripId)
    ])
    trip.value = tripRes.data.data || null
    claims.value = claimsRes.data.data || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
})
</script>
