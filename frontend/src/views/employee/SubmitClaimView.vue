<template>
  <AppLayout>
    <!-- Header -->
    <div class="max-w-2xl mx-auto py-8 px-4">
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-white mb-2">{{ isEdit ? 'Edit Klaim Reimbursement' : 'Ajukan Klaim Reimbursement' }}</h1>
        <p class="text-slate-400 text-sm">{{ isEdit ? 'Ubah formulir di bawah untuk memperbarui klaim' : 'Isi formulir di bawah untuk mengajukan reimbursement klaim' }}</p>
      </div>
    </div>

    <!-- Form Card -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6 max-w-xl">
      <form @submit.prevent="handleSubmit" class="space-y-5">
        <!-- Trip Selector (Only visible if tripId is NOT in route parameters) -->
        <div v-if="!tripId">
          <label class="block text-sm font-medium text-slate-300 mb-2">Perjalanan Dinas <span class="text-red-400">*</span></label>
          <select v-model="form.trip_id" required class="input-field select-field">
            <option value="" disabled>Pilih perjalanan dinas yang sedang berlangsung / selesai...</option>
            <option v-for="t in eligibleTrips" :key="t.id" :value="t.id">
              {{ t.nomor_surat }} - {{ t.destination }} ({{ t.status }})
            </option>
          </select>
          <p v-if="eligibleTrips.length === 0 && !loadingTrips" class="text-amber-400 text-xs mt-1">
            Belum ada perjalanan dinas dalam status ON_DUTY atau COMPLETED untuk diajukan klaim.
          </p>
        </div>

        <!-- Judul Klaim -->
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Judul Klaim <span class="text-red-400">*</span></label>
          <input v-model="form.title" required placeholder="cth. Biaya transportasi"
            class="input-field" />
        </div>

        <!-- Deskripsi -->
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Keterangan</label>
          <textarea v-model="form.description" rows="2" placeholder="Detail pengeluaran..."
            class="input-field resize-none"></textarea>
        </div>

        <!-- Nominal -->
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Nominal (Rp) <span class="text-red-400">*</span></label>
          <input :value="displayAmount" @input="handleAmountInput" type="text" required placeholder="500.000"
            class="input-field" />
        </div>

        <!-- Tanggal Transaksi (Pilihan dari tanggal dinas) -->
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Tanggal Transaksi <span class="text-red-400">*</span></label>
          <select v-model="form.transaction_date" required class="input-field select-field" :disabled="!form.trip_id || loadingDates">
            <option value="" disabled>{{ loadingDates ? 'Memuat tanggal...' : 'Pilih tanggal dari perjalanan dinas...' }}</option>
            <option v-for="d in tripDates" :key="d" :value="d" :disabled="isDateDisabled(d)">
              {{ formatDate(d) }} {{ isDateDisabled(d) ? '(Sudah Dicairkan)' : '' }}
            </option>
          </select>
          <p v-if="!form.trip_id" class="text-slate-500 text-xs mt-1">Silakan pilih perjalanan dinas terlebih dahulu.</p>
        </div>

        <!-- Bukti Pembayaran -->
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Bukti Pembayaran (PDF, PNG, JPG, JPEG)</label>
          <div class="border-2 border-dashed border-slate-700 rounded-xl p-5 text-center hover:border-blue-500 transition cursor-pointer"
            @click="$refs.fileInput.click()">
            <p class="text-slate-400 text-sm">Klik untuk pilih file bukti</p>
          </div>
          <input ref="fileInput" type="file" multiple accept=".pdf,.png,.jpg,.jpeg" class="hidden" @change="onFileChange" />
          <div v-if="files.length > 0" class="mt-3 space-y-2">
            <div v-for="(f, i) in files" :key="i"
              class="flex items-center justify-between bg-slate-800 rounded-lg px-3 py-2">
              <span class="text-slate-300 text-xs truncate">{{ f.name }}</span>
              <button type="button" @click="files.splice(i, 1)" class="text-red-400 hover:text-red-300 ml-2">✕</button>
            </div>
          </div>
        </div>

        <!-- Error -->
        <div v-if="error" class="bg-red-900/40 border border-red-700 text-red-300 rounded-xl px-4 py-3 text-sm">{{ error }}</div>

        <!-- Submit -->
        <div class="flex gap-3 pt-2">
          <router-link :to="tripId ? `/trips/${tripId}` : '/claims'"
            class="flex-1 text-center bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium px-4 py-3 rounded-xl text-sm transition">
            Batal
          </router-link>
          <button type="submit" :disabled="loading"
            class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 rounded-xl disabled:opacity-50 transition">
            {{ loading ? (isEdit ? 'Menyimpan...' : 'Mengajukan...') : (isEdit ? 'Simpan Perubahan' : 'Ajukan Klaim') }}
          </button>
        </div>
      </form>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { claimsApi } from '@/api/claims'
import { tripsApi } from '@/api/trips'
import { useToast } from '@/composables/useToast'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const isEdit = computed(() => route.path.includes('/edit'))
const claimId = route.params.id
const tripId = computed(() => !isEdit.value && route.path.includes('/trips/') ? route.params.id : null)

const loading = ref(false)
const loadingTrips = ref(false)
const loadingDates = ref(false)
const error = ref('')
const files = ref([])
const fileInput = ref(null)
const eligibleTrips = ref([])

const tripDates = ref([])
const existingClaims = ref([])

const form = ref({
  trip_id: tripId.value || '',
  title: '',
  description: '',
  amount: '',
  transaction_date: ''
})

const displayAmount = ref('')

function handleAmountInput(e) {
  let val = e.target.value.replace(/\D/g, '')
  if (val) {
    displayAmount.value = new Intl.NumberFormat('id-ID').format(val)
    form.value.amount = parseFloat(val)
  } else {
    displayAmount.value = ''
    form.value.amount = ''
  }
}

watch(() => form.value.amount, (newVal) => {
  if (newVal) {
    displayAmount.value = new Intl.NumberFormat('id-ID').format(newVal)
  } else {
    displayAmount.value = ''
  }
})

function onFileChange(e) {
  const selectedFiles = Array.from(e.target.files)
  const allowedExtensions = ['.pdf', '.png', '.jpg', '.jpeg']
  const hasInvalid = selectedFiles.some(file => {
    const name = file.name.toLowerCase()
    return !allowedExtensions.some(ext => name.endsWith(ext))
  })
  if (hasInvalid) {
    error.value = 'format file harus pdf, png, jpg, atau jpeg'
    e.target.value = ''
    return
  }
  error.value = ''
  files.value.push(...selectedFiles)
  e.target.value = ''
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

function generateTripDates(trip) {
  if (!trip?.start_date || !trip?.end_date) {
    tripDates.value = []
    return
  }

  const start = new Date(trip.start_date)
  const end = new Date(trip.end_date)
  const dates = []

  let current = new Date(start)
  while (current <= end) {
    const yyyy = current.getFullYear()
    const mm = String(current.getMonth() + 1).padStart(2, '0')
    const dd = String(current.getDate()).padStart(2, '0')
    dates.push(`${yyyy}-${mm}-${dd}`)
    current.setDate(current.getDate() + 1)
  }

  tripDates.value = dates
}

function isDateDisabled(dateStr) {
  return existingClaims.value.some(c => {
    if (isEdit.value && c.id === Number(claimId)) return false
    if (!c.transaction_date) return false
    const claimDate = new Date(c.transaction_date).toISOString().split('T')[0]
    return claimDate === dateStr && c.status === 'APPROVED'
  })
}

async function loadTripDetailsAndClaims(id) {
  loadingDates.value = true
  try {
    const [tripRes, claimsRes] = await Promise.all([
      tripsApi.getById(id),
      claimsApi.getByTripId(id)
    ])
    const trip = tripRes.data.data
    existingClaims.value = claimsRes.data.data || []
    generateTripDates(trip)
  } catch (e) {
    console.error('Failed to load trip dates/claims', e)
  } finally {
    loadingDates.value = false
  }
}

watch(() => form.value.trip_id, async (newVal) => {
  if (!newVal) {
    tripDates.value = []
    existingClaims.value = []
    form.value.transaction_date = ''
    return
  }
  form.value.transaction_date = ''
  await loadTripDetailsAndClaims(newVal)
})

async function handleSubmit() {
  loading.value = true
  error.value = ''
  try {
    const fd = new FormData()
    fd.append('trip_id', form.value.trip_id)
    fd.append('title', form.value.title)
    fd.append('description', form.value.description)
    fd.append('amount', form.value.amount)
    fd.append('transaction_date', form.value.transaction_date)
    files.value.forEach(f => fd.append('files', f))
    
    if (isEdit.value) {
      await claimsApi.update(claimId, fd)
      toast.success('Klaim reimbursement berhasil diperbarui!')
    } else {
      await claimsApi.submit(fd)
      toast.success('Klaim reimbursement berhasil diajukan!')
    }
    
    if (tripId.value) {
      router.push(`/trips/${tripId.value}`)
    } else {
      router.push('/claims')
    }
  } catch (e) {
    error.value = e.response?.data?.message || 'Gagal menyimpan klaim.'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  if (isEdit.value) {
    loading.value = true
    try {
      const res = await claimsApi.getById(claimId)
      const claimData = res.data.data
      form.value.trip_id = claimData.trip_id
      form.value.title = claimData.title
      form.value.description = claimData.description
      form.value.amount = claimData.amount
      if (claimData.transaction_date) {
        form.value.transaction_date = claimData.transaction_date.split('T')[0]
      }
      await loadTripDetailsAndClaims(claimData.trip_id)
    } catch (e) {
      console.error(e)
      error.value = 'Gagal memuat data klaim.'
    } finally {
      loading.value = false
    }
  } else if (tripId.value) {
    await loadTripDetailsAndClaims(tripId.value)
  } else {
    loadingTrips.value = true
    try {
      const res = await tripsApi.getMyTrips({ page: 1, limit: 100 })
      const all = res.data.data?.items || []
      // Filter for COMPLETED or ON_DUTY trips
      eligibleTrips.value = all.filter(t => t.status === 'COMPLETED' || t.status === 'ON_DUTY')
    } catch (e) {
      console.error(e)
    } finally {
      loadingTrips.value = false
    }
  }
})
</script>

<style scoped>
@reference "../../style.css";
.input-field {
  @apply w-full bg-slate-800 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition;
}
.select-field {
  @apply appearance-none pr-10 bg-slate-800 text-slate-300;
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%2394a3b8' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
  background-position: right 0.75rem center;
  background-repeat: no-repeat;
  background-size: 1.5em 1.5em;
}
</style>
