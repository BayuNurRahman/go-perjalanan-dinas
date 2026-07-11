<template>
  <AppLayout>
    <!-- Header -->
    <div class="flex items-center gap-3 mb-6">
      <router-link to="/trips" class="text-slate-400 hover:text-white transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </router-link>
      <div>
        <h1 class="text-2xl font-bold text-white">Ajukan Perjalanan Dinas</h1>
        <p class="text-slate-400 text-sm">Isi formulir di bawah untuk mengajukan perjalanan dinas</p>
      </div>
    </div>

    <!-- Form Card -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6 max-w-2xl">
      <form @submit.prevent="handleSubmit" class="space-y-5">
        <!-- Nomor Surat -->
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Nomor Surat <span class="text-red-400">*</span></label>
          <input v-model="form.nomor_surat" required placeholder="DIN/2024/001"
            class="input-field" />
        </div>

        <!-- Tujuan -->
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Tujuan <span class="text-red-400">*</span></label>
          <input v-model="form.destination" required placeholder="Jakarta, Indonesia"
            class="input-field" />
        </div>

        <!-- Inisiator -->
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Nama Inisiator <span class="text-red-400">*</span></label>
          <input v-model="form.initiator" required placeholder="Nama lengkap pengaju"
            class="input-field" />
        </div>

        <!-- Tanggal -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Tanggal Mulai <span class="text-red-400">*</span></label>
            <input v-model="form.start_date" type="date" required class="input-field" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Tanggal Selesai <span class="text-red-400">*</span></label>
            <input v-model="form.end_date" type="date" required class="input-field" />
          </div>
        </div>

        <!-- Deskripsi -->
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Deskripsi <span class="text-red-400">*</span></label>
          <textarea v-model="form.description" required rows="3" placeholder="Jelaskan tujuan perjalanan dinas ini..."
            class="input-field resize-none"></textarea>
        </div>

        <!-- Ringkasan -->
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Ringkasan <span class="text-slate-500">(opsional)</span></label>
          <textarea v-model="form.summary" rows="2" placeholder="Ringkasan singkat perjalanan..."
            class="input-field resize-none"></textarea>
        </div>

        <!-- Upload File -->
        <div>
          <label class="block text-sm font-medium text-slate-300 mb-2">Dokumen Pendukung <span class="text-slate-500">(PDF)</span></label>
          <div class="border-2 border-dashed border-slate-700 rounded-xl p-6 text-center hover:border-blue-500 transition cursor-pointer"
            @click="$refs.fileInput.click()" @dragover.prevent @drop.prevent="onDrop">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-slate-500 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
            <p class="text-slate-400 text-sm">Klik atau seret file ke sini</p>
            <p class="text-slate-500 text-xs mt-1">Maks. ukuran: 10MB per file</p>
          </div>
          <input ref="fileInput" type="file" multiple accept=".pdf" class="hidden" @change="onFileChange" />
          <!-- File List -->
          <div v-if="files.length > 0" class="mt-3 space-y-2">
            <div v-for="(f, i) in files" :key="i"
              class="flex items-center justify-between bg-slate-800 rounded-lg px-3 py-2">
              <span class="text-slate-300 text-xs truncate">{{ f.name }}</span>
              <button type="button" @click="removeFile(i)" class="text-red-400 hover:text-red-300 ml-2 flex-shrink-0">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>
        </div>

        <!-- Error -->
        <div v-if="error" class="bg-red-900/40 border border-red-700 text-red-300 rounded-xl px-4 py-3 text-sm">{{ error }}</div>

        <!-- Submit -->
        <div class="flex gap-3 pt-2">
          <router-link to="/trips"
            class="flex-1 text-center bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium px-4 py-3 rounded-xl text-sm transition">
            Batal
          </router-link>
          <button type="submit" :disabled="loading"
            class="flex-1 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold px-4 py-3 rounded-xl text-sm transition flex items-center justify-center gap-2">
            <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            {{ loading ? 'Mengajukan...' : 'Ajukan Perjalanan' }}
          </button>
        </div>
      </form>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { tripsApi } from '@/api/trips'

const router = useRouter()
const loading = ref(false)
const error = ref('')
const files = ref([])
const fileInput = ref(null)

const form = ref({
  destination: '',
  start_date: '',
  end_date: '',
  description: '',
  initiator: '',
  summary: '',
  nomor_surat: '',
})

function onFileChange(e) {
  const selectedFiles = Array.from(e.target.files)
  const hasNonPdf = selectedFiles.some(file => !file.name.toLowerCase().endsWith('.pdf'))
  if (hasNonPdf) {
    error.value = 'format file harus pdf'
    e.target.value = ''
    return
  }
  error.value = ''
  files.value.push(...selectedFiles)
  e.target.value = ''
}

function onDrop(e) {
  const selectedFiles = Array.from(e.dataTransfer.files)
  const hasNonPdf = selectedFiles.some(file => !file.name.toLowerCase().endsWith('.pdf'))
  if (hasNonPdf) {
    error.value = 'format file harus pdf'
    return
  }
  error.value = ''
  files.value.push(...selectedFiles)
}

function removeFile(i) {
  files.value.splice(i, 1)
}

async function handleSubmit() {
  loading.value = true
  error.value = ''
  try {
    const fd = new FormData()
    Object.entries(form.value).forEach(([k, v]) => fd.append(k, v))
    files.value.forEach(f => fd.append('files', f))
    await tripsApi.create(fd)
    router.push('/trips')
  } catch (e) {
    error.value = e.response?.data?.message || 'Gagal mengajukan perjalanan dinas.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
@reference "../../style.css";
.input-field {
  @apply w-full bg-slate-800 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition;
}
</style>
