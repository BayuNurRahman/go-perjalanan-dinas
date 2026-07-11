<template>
  <AppLayout>
    <div class="flex items-center gap-3 mb-6">
      <router-link :to="`/finance/trips/${tripId}`" class="text-slate-400 hover:text-white transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </router-link>
      <div>
        <h1 class="text-2xl font-bold text-white">Pencairan Dana</h1>
        <p class="text-slate-400 text-sm">Perjalanan #{{ tripId }}</p>
      </div>
    </div>

    <!-- Warning Info -->
    <div class="bg-yellow-500/10 border border-yellow-500/30 rounded-2xl px-5 py-4 mb-5 flex gap-3">
      <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-yellow-400 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <div>
        <p class="text-yellow-300 font-medium text-sm">Perhatian</p>
        <p class="text-yellow-400/80 text-xs mt-0.5">
          Pastikan semua klaim telah direview dan trip sudah berstatus COMPLETED sebelum melakukan pencairan.
        </p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <!-- Form -->
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6">
        <h2 class="text-white font-semibold mb-5">Form Pencairan Dana</h2>

        <form @submit.prevent="handleDisburse" class="space-y-5">
          <!-- Jumlah -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">
              Jumlah Pencairan (Rp) <span class="text-red-400">*</span>
            </label>
            <input v-model="form.amount" type="number" required min="1"
              placeholder="Masukkan nominal pencairan"
              class="w-full bg-slate-800 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500 transition" />
            <p v-if="form.amount" class="text-slate-400 text-xs mt-1">
              {{ formatCurrency(form.amount) }}
            </p>
          </div>

          <!-- Reference ID -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">
              Nomor Referensi <span class="text-red-400">*</span>
            </label>
            <input v-model="form.reference_id" required
              placeholder="cth. TRF-2024-001234"
              class="w-full bg-slate-800 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500 transition" />
            <p class="text-slate-500 text-xs mt-1">Nomor transfer / referensi transaksi bank</p>
          </div>

          <!-- Notes -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Catatan</label>
            <textarea v-model="form.notes" rows="3" placeholder="Catatan pencairan dana (opsional)..."
              class="w-full bg-slate-800 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500 transition resize-none">
            </textarea>
          </div>

          <!-- Error / Success -->
          <div v-if="errorMsg" class="bg-red-900/40 border border-red-700 text-red-300 rounded-xl px-4 py-3 text-sm">{{ errorMsg }}</div>
          <div v-if="successMsg" class="bg-green-900/40 border border-green-700 text-green-300 rounded-xl px-4 py-3 text-sm">{{ successMsg }}</div>

          <!-- Submit -->
          <button type="submit" :disabled="loading"
            class="w-full bg-teal-600 hover:bg-teal-700 disabled:opacity-40 text-white font-semibold py-3 rounded-xl text-sm transition flex items-center justify-center gap-2">
            <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
            {{ loading ? 'Memproses...' : 'Proses Pencairan Dana' }}
          </button>
        </form>
      </div>

      <!-- Claim Summary -->
      <div class="space-y-4">
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6">
          <h3 class="text-white font-semibold mb-4">Ringkasan Klaim Disetujui</h3>

          <div v-if="claimsLoading" class="space-y-2">
            <div v-for="i in 3" :key="i" class="h-12 bg-slate-800 rounded-xl animate-pulse" />
          </div>

          <div v-else-if="approvedClaims.length === 0" class="text-slate-500 text-sm text-center py-4">
            Belum ada klaim yang disetujui
          </div>

          <div v-else class="space-y-2">
            <div v-for="claim in approvedClaims" :key="claim.id"
              class="flex items-center justify-between bg-slate-800 rounded-xl px-4 py-3">
              <p class="text-slate-300 text-sm truncate">{{ claim.title }}</p>
              <p class="text-white text-sm font-semibold ml-3 flex-shrink-0">{{ formatCurrency(claim.amount) }}</p>
            </div>

            <div class="flex items-center justify-between bg-blue-600/10 border border-blue-600/20 rounded-xl px-4 py-3 mt-3">
              <p class="text-blue-300 font-semibold text-sm">Total</p>
              <p class="text-blue-300 font-bold text-lg">{{ formatCurrency(totalApproved) }}</p>
            </div>
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
import { tripsApi } from '@/api/trips'
import { claimsApi } from '@/api/claims'

const route = useRoute()
const router = useRouter()
const tripId = route.params.id

const loading = ref(false)
const claimsLoading = ref(true)
const claims = ref([])
const errorMsg = ref('')
const successMsg = ref('')

const form = ref({ amount: '', reference_id: '', notes: '' })

const approvedClaims = computed(() => claims.value.filter(c => c.status === 'APPROVED'))
const totalApproved = computed(() => approvedClaims.value.reduce((s, c) => s + (Number(c.amount) || 0), 0))

function formatCurrency(n) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(Number(n) || 0)
}

async function handleDisburse() {
  loading.value = true
  errorMsg.value = ''
  successMsg.value = ''
  try {
    await tripsApi.disburse(tripId, {
      amount: Number(form.value.amount),
      reference_id: form.value.reference_id,
      notes: form.value.notes,
    })
    successMsg.value = 'Pencairan dana berhasil diproses!'
    form.value = { amount: '', reference_id: '', notes: '' }
    setTimeout(() => router.push(`/finance/trips/${tripId}`), 1500)
  } catch (e) {
    errorMsg.value = e.response?.data?.message || 'Gagal memproses pencairan.'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const res = await claimsApi.getByTripId(tripId)
    claims.value = res.data.data || []
  } catch (e) { console.error(e) }
  finally { claimsLoading.value = false }
})
</script>
