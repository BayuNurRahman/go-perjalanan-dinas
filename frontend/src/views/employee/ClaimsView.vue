<template>
  <AppLayout>
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-white">Klaim Reimbursement Saya</h1>
        <p class="text-slate-400 text-sm mt-1">Daftar seluruh pengajuan klaim reimbursement Anda</p>
      </div>
      <router-link to="/claims/create"
        class="inline-flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white font-semibold px-4 py-2.5 rounded-xl text-sm transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Ajukan Klaim
      </router-link>
    </div>

    <!-- Table Card -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <div v-if="loading" class="p-6 space-y-3">
        <div v-for="i in 5" :key="i" class="h-14 bg-slate-800 rounded-xl animate-pulse" />
      </div>

      <div v-else-if="claims.length === 0" class="p-16 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 text-slate-700 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
        <p class="text-slate-500">Belum ada klaim reimbursement yang diajukan</p>
      </div>

      <table v-else class="w-full text-sm">
        <thead class="bg-slate-800/60 border-b border-slate-800">
          <tr>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Klaim</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Perjalanan Dinas</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Tanggal Transaksi</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Nominal</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Status</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Bukti</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="claim in claims" :key="claim.id" class="hover:bg-slate-800/40 transition">
            <td class="px-6 py-4">
              <p class="text-white font-medium">{{ claim.title }}</p>
              <p class="text-slate-400 text-xs mt-0.5 max-w-xs truncate">{{ claim.description }}</p>
            </td>
            <td class="px-6 py-4 text-slate-300">
              <router-link :to="`/trips/${claim.trip_id}`" class="text-blue-400 hover:underline text-xs block font-mono">
                {{ claim.tripNomorSurat }}
              </router-link>
              <span class="text-slate-400 text-xs">{{ claim.tripDestination }}</span>
            </td>
            <td class="px-6 py-4 text-slate-300">
              {{ formatDate(claim.transaction_date) }}
            </td>
            <td class="px-6 py-4 text-white font-semibold">
              {{ formatCurrency(claim.amount) }}
            </td>
            <td class="px-6 py-4">
              <StatusBadge :status="claim.status" />
              <p v-if="claim.rejected_reason" class="text-red-400 text-xs mt-1 max-w-xs leading-tight">
                Ket: {{ claim.rejected_reason }}
              </p>
            </td>
            <td class="px-6 py-4">
              <div v-if="claim.attachment_paths || claim.attachment_path" class="flex flex-col gap-1">
                <button v-for="fname in getAttachments(claim)" :key="fname"
                  @click="downloadFile(claim.id, fname)"
                  class="text-blue-400 hover:text-blue-300 text-xs font-medium text-left transition truncate max-w-[120px]">
                  📄 {{ fname }}
                </button>
              </div>
              <span v-else class="text-slate-500 text-xs">-</span>
            </td>
            <td class="px-6 py-4">
              <div v-if="claim.status === 'PENDING'" class="flex items-center gap-3">
                <router-link :to="`/claims/${claim.id}/edit`" class="text-blue-400 hover:text-blue-300 transition" title="Edit Klaim">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </router-link>
                <button @click="handleDelete(claim.id)" class="text-red-400 hover:text-red-300 transition" title="Hapus Klaim">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
              <span v-else class="text-slate-500 text-xs">-</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import { tripsApi } from '@/api/trips'
import { claimsApi } from '@/api/claims'
import { useToast } from '@/composables/useToast'

const toast = useToast()

const claims = ref([])
const loading = ref(true)

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

function formatCurrency(n) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(n)
}

function getAttachments(claim) {
  const paths = claim.attachment_paths || claim.attachment_path || ''
  if (!paths) return []
  return paths.split(',').map(p => {
    const trimmed = p.trim()
    const normalized = trimmed.replace(/\\/g, '/')
    return normalized.split('/').pop()
  }).filter(Boolean)
}

async function downloadFile(claimId, filename) {
  try {
    const res = await claimsApi.downloadAttachment(claimId, filename)
    const url = URL.createObjectURL(res.data)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)
  } catch (e) {
    console.error(e)
  }
}

async function handleDelete(claimId) {
  if (!confirm('Apakah Anda yakin ingin menghapus klaim reimbursement ini? Berkas bukti di komputer/lokal juga akan ikut terhapus.')) {
    return
  }
  try {
    await claimsApi.delete(claimId)
    toast.success('Klaim reimbursement berhasil dihapus!')
    await reloadClaims()
  } catch (e) {
    console.error(e)
    toast.error(e.response?.data?.message || 'Gagal menghapus klaim.')
  }
}

async function reloadClaims() {
  loading.value = true
  try {
    const res = await tripsApi.getMyTrips({ page: 1, limit: 100 })
    const trips = res.data.data?.items || []

    const claimsPromises = trips.map(t => claimsApi.getByTripId(t.id))
    const claimsResponses = await Promise.all(claimsPromises)

    const allClaims = []
    claimsResponses.forEach((resp, idx) => {
      const tripClaims = resp.data.data || []
      const t = trips[idx]
      tripClaims.forEach(c => {
        allClaims.push({
          ...c,
          tripDestination: t.destination,
          tripNomorSurat: t.nomor_surat,
          tripStatus: t.status
        })
      })
    })

    allClaims.sort((a, b) => new Date(b.created_at || b.transaction_date) - new Date(a.created_at || a.transaction_date))
    claims.value = allClaims
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await reloadClaims()
})
</script>
