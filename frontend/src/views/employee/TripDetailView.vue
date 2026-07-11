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
        <div class="flex items-center gap-3">
          <h1 class="text-2xl font-bold text-white">Detail Perjalanan Dinas</h1>
          <StatusBadge v-if="trip" :status="trip.status" />
          
          <!-- Mulai Dinas Button (when APPROVED) -->
          <button v-if="trip && trip.status === 'APPROVED'" type="button" @click="handleStartDuty" :disabled="statusProcessing"
            class="px-3 py-1.5 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold rounded-xl text-xs transition flex items-center gap-1.5 shadow-lg">
            <svg v-if="statusProcessing" class="animate-spin w-3 h-3" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            Mulai Dinas (On Duty)
          </button>

          <!-- Selesaikan Dinas Button (when ON_DUTY) -->
          <button v-if="trip && trip.status === 'ON_DUTY'" type="button" @click="openCompleteModal" :disabled="statusProcessing"
            class="px-3 py-1.5 bg-purple-600 hover:bg-purple-700 disabled:opacity-50 text-white font-semibold rounded-xl text-xs transition flex items-center gap-1.5 shadow-lg">
            Selesaikan Dinas (Complete)
          </button>
        </div>
        <p class="text-slate-400 text-sm mt-1">
          {{ trip && trip.status === 'PENDING' ? 'Perbarui data pengajuan perjalanan dinas Anda di bawah' : 'Informasi detail pengajuan perjalanan dinas Anda' }}
        </p>
      </div>
    </div>

    <!-- Loading Skeleton -->
    <div v-if="loading" class="space-y-4 max-w-2xl">
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6 animate-pulse h-96" />
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6 animate-pulse h-48" />
    </div>

    <div v-else-if="trip" class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start max-w-6xl">
      <!-- Left Column: Form Card -->
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6 lg:col-span-7">
        <form @submit.prevent="handleUpdateTrip" class="space-y-5">
          <!-- Nomor Surat -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Nomor Surat <span class="text-red-400">*</span></label>
            <input v-model="form.nomor_surat" required placeholder="DIN/2024/001"
              :disabled="trip.status !== 'PENDING'" class="input-field" />
          </div>

          <!-- Tujuan -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Tujuan <span class="text-red-400">*</span></label>
            <input v-model="form.destination" required placeholder="Jakarta, Indonesia"
              :disabled="trip.status !== 'PENDING'" class="input-field" />
          </div>

          <!-- Inisiator -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Nama Inisiator <span class="text-red-400">*</span></label>
            <input v-model="form.initiator" required placeholder="Nama lengkap pengaju"
              :disabled="trip.status !== 'PENDING'" class="input-field" />
          </div>

          <!-- Tanggal -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-2">Tanggal Mulai <span class="text-red-400">*</span></label>
              <input v-model="form.start_date" type="date" required
                :disabled="trip.status !== 'PENDING'" class="input-field" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-2">Tanggal Selesai <span class="text-red-400">*</span></label>
              <input v-model="form.end_date" type="date" required
                :disabled="trip.status !== 'PENDING'" class="input-field" />
            </div>
          </div>

          <!-- Deskripsi -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Deskripsi <span class="text-red-400">*</span></label>
            <textarea v-model="form.description" required rows="3" placeholder="Jelaskan tujuan perjalanan dinas ini..."
              :disabled="trip.status !== 'PENDING'" class="input-field resize-none"></textarea>
          </div>

          <!-- Ringkasan -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Ringkasan <span class="text-slate-500">(opsional)</span></label>
            <textarea v-model="form.summary" rows="2" placeholder="Ringkasan singkat perjalanan..."
              :disabled="trip.status !== 'PENDING'" class="input-field resize-none"></textarea>
          </div>

          <!-- Notes/Catatan dari Manager/Finance -->
          <div v-if="trip.notes">
            <label class="block text-sm font-medium text-slate-300 mb-2">Catatan Persetujuan/Penolakan</label>
            <div class="bg-slate-950/60 border border-slate-800 rounded-xl p-4 text-slate-300 text-sm">
              {{ trip.notes }}
            </div>
          </div>

          <!-- Dokumen Surat Perjalanan Dinas (PDF) -->
          <div class="border-t border-slate-800 pt-5 space-y-4">
            <div class="flex items-center justify-between">
              <label class="block text-sm font-medium text-slate-300">Surat Perjalanan Dinas <span class="text-slate-500">(PDF saja)</span></label>
              
              <!-- Preview Button -->
              <button type="button" @click="handlePreviewPdf"
                class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white text-xs font-semibold rounded-lg transition border border-slate-700">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                Preview PDF
              </button>
            </div>

            <!-- Jika sudah ada file yang diunggah (tampilkan dokumen pengajuan awal) -->
            <div v-if="attachments.length > 0" class="bg-slate-950/40 border border-slate-800 rounded-xl p-4 flex items-center justify-between">
              <div class="flex items-center gap-3 truncate">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-red-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
                <div class="truncate">
                  <p class="text-white text-sm font-medium truncate">{{ attachments[0] }}</p>
                  <p class="text-slate-500 text-xs">Dokumen sudah diunggah</p>
                </div>
              </div>
              <button type="button" @click="downloadFile(attachments[0])"
                class="text-blue-400 hover:text-blue-300 text-xs font-semibold px-2 py-1 transition flex-shrink-0">
                Unduh
              </button>
            </div>

            <!-- Dropzone / Input file PDF (selalu tampil) -->
            <div v-if="canUpdate">
              <div class="border-2 border-dashed border-slate-700 rounded-xl p-5 text-center hover:border-blue-500 transition cursor-pointer"
                @click="$refs.pdfInput.click()" @dragover.prevent @drop.prevent="onDropPdf">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-slate-500 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                </svg>
                <p class="text-slate-400 text-xs">Pilih file PDF baru untuk diunggah</p>
                <p class="text-slate-500 text-[10px] mt-1">Maks. ukuran: 5MB</p>
              </div>
              <input ref="pdfInput" type="file" accept=".pdf" class="hidden" @change="onPdfFileChange" />
              
              <!-- Selected File Preview -->
              <div v-if="selectedPdfFile" class="mt-3 flex items-center justify-between bg-blue-950/30 border border-blue-900/50 rounded-xl px-4 py-3">
                <div class="flex items-center gap-2 truncate">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  <span class="text-slate-300 text-xs truncate">{{ selectedPdfFile.name }}</span>
                </div>
                <button type="button" @click="removeSelectedPdf" class="text-red-400 hover:text-red-300 ml-2 flex-shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>
            <div v-else class="bg-slate-900/40 border border-slate-800/80 rounded-xl p-3 text-center text-xs text-slate-500">
              Dokumen tidak dapat diubah karena status perjalanan dinas sudah diproses.
            </div>
          </div>

          <!-- Action Buttons -->
          <div v-if="canUpdate" class="flex flex-col sm:flex-row gap-3 pt-4 border-t border-slate-800">
            <button type="button" @click="resetForm"
              class="flex-1 bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium px-4 py-3 rounded-xl text-sm transition">
              Batal
            </button>
            <button type="submit" :disabled="editProcessing"
              class="flex-1 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold px-4 py-3 rounded-xl text-sm transition flex items-center justify-center gap-2">
              <svg v-if="editProcessing" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
              </svg>
              Simpan Perubahan
            </button>
            <button v-if="trip.status === 'PENDING'" type="button" @click="handleDeleteTrip" :disabled="deleteProcessing"
              class="flex-1 bg-red-600 hover:bg-red-700 disabled:opacity-50 text-white font-semibold px-4 py-3 rounded-xl text-sm transition flex items-center justify-center gap-2">
              <svg v-if="deleteProcessing" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
              </svg>
              Hapus Perjalanan
            </button>
          </div>
        </form>
      </div>

      <!-- Right Column: Bukti Perjalanan Dinas -->
      <div v-if="trip.status === 'COMPLETED' && attachments.length > 1" class="bg-slate-900 border border-slate-800 rounded-2xl p-6 lg:col-span-5 space-y-4">
        <h3 class="text-lg font-bold text-white">Bukti Perjalanan Dinas</h3>
        <p class="text-slate-400 text-sm">
          Berkas bukti perjalanan dinas yang telah diunggah oleh karyawan.
        </p>

        <!-- Document Card -->
        <div class="bg-slate-950/60 border border-slate-800 rounded-xl p-4 flex flex-col gap-4">
          <div class="flex items-center gap-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 text-purple-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
            </svg>
            <div class="truncate flex-1">
              <p class="text-white text-sm font-semibold truncate">{{ attachments[1] }}</p>
              <p class="text-slate-500 text-xs">Format: PDF</p>
            </div>
          </div>

          <div class="flex gap-2">
            <button type="button" @click="handlePreviewProofPdf"
              class="flex-1 inline-flex items-center justify-center gap-1.5 px-3 py-2 bg-purple-600/20 hover:bg-purple-600/30 text-purple-300 font-semibold rounded-lg text-xs transition border border-purple-500/30">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              Pratinjau PDF
            </button>
            <button type="button" @click="openCompleteModal"
              class="flex-1 inline-flex items-center justify-center gap-1.5 px-3 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 font-semibold rounded-lg text-xs transition border border-slate-700">
              Edit
            </button>
          </div>
        </div>
      </div>

      <!-- PDF Preview Modal -->
      <div v-if="showPdfModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm animate-fade-in">
        <div class="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-4xl max-h-[90vh] flex flex-col shadow-2xl overflow-hidden">
          <!-- Modal Header -->
          <div class="flex items-center justify-between px-6 py-4 border-b border-slate-800">
            <h3 class="text-white font-semibold flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              {{ pdfModalTitle }}
            </h3>
            <button @click="closePdfModal" class="text-slate-400 hover:text-white transition">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Modal Body -->
          <div class="flex-1 bg-slate-950 p-4 flex items-center justify-center min-h-[60vh]">
            <div v-if="loadingPdf" class="flex flex-col items-center gap-3">
              <svg class="animate-spin w-8 h-8 text-blue-500" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
              </svg>
              <p class="text-slate-400 text-sm">Mengunduh file untuk pratinjau...</p>
            </div>
            <iframe v-else-if="pdfUrl" :src="pdfUrl" class="w-full h-[70vh] rounded-lg border border-slate-800" />
            <div v-else class="text-slate-500 text-sm">Gagal memuat pratinjau PDF.</div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="text-slate-500 text-center py-16">Data tidak ditemukan.</div>

    <!-- Complete Trip Modal -->
    <div v-if="showCompleteModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm animate-fade-in">
      <div class="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-md flex flex-col shadow-2xl overflow-hidden">
        <!-- Modal Header -->
        <div class="flex items-center justify-between px-6 py-4 border-b border-slate-800">
          <h3 class="text-white font-semibold flex items-center gap-2">
            Selesaikan Perjalanan Dinas
          </h3>
          <button @click="showCompleteModal = false" class="text-slate-400 hover:text-white transition">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6 space-y-4">
          <p class="text-slate-300 text-sm">
            Silakan unggah bukti perjalanan dinas Anda (misal: laporan kegiatan). Berkas harus dalam format PDF.
          </p>

          <!-- File Upload Dropzone -->
          <div class="border-2 border-dashed border-slate-700 rounded-xl p-5 text-center hover:border-blue-500 transition cursor-pointer"
            @click="$refs.completeFileInput.click()" @dragover.prevent @drop.prevent="onDropCompleteFile">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-8 h-8 text-slate-500 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
            </svg>
            <p class="text-slate-400 text-xs font-medium">Klik atau seret file PDF ke sini</p>
            <p class="text-slate-500 text-[10px] mt-1">Maks. ukuran: 5MB</p>
          </div>
          <input ref="completeFileInput" type="file" accept=".pdf" class="hidden" @change="onCompleteFileChange" />

          <!-- Selected File Preview -->
          <div v-if="completeFile" class="flex items-center justify-between bg-blue-950/30 border border-blue-900/50 rounded-xl px-4 py-3">
            <div class="flex items-center gap-2 truncate">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-blue-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <span class="text-slate-300 text-xs truncate">{{ completeFile.name }}</span>
            </div>
            <button type="button" @click="completeFile = null" class="text-red-400 hover:text-red-300 ml-2 flex-shrink-0">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Error message if any -->
          <div v-if="completeError" class="bg-red-900/40 border border-red-700 text-red-300 rounded-xl px-4 py-2.5 text-xs">
            {{ completeError }}
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="px-6 py-4 bg-slate-950 border-t border-slate-800 flex justify-end gap-3">
          <button type="button" @click="showCompleteModal = false" :disabled="statusProcessing"
            class="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl text-xs font-semibold transition">
            Batal
          </button>
          <button type="button" @click="handleCompleteTrip" :disabled="statusProcessing || !completeFile"
            class="px-4 py-2 bg-purple-600 hover:bg-purple-700 disabled:opacity-50 text-white rounded-xl text-xs font-semibold transition flex items-center gap-1.5">
            <svg v-if="statusProcessing" class="animate-spin w-3 h-3" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            Simpan
          </button>
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
import { tripsApi } from '@/api/trips'
import { useToast } from '@/composables/useToast'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const tripId = route.params.id

const loading = ref(true)
const editProcessing = ref(false)
const deleteProcessing = ref(false)
const trip = ref(null)

const selectedPdfFile = ref(null)
const pdfInput = ref(null)

const showPdfModal = ref(false)
const loadingPdf = ref(false)
const pdfUrl = ref('')
const pdfModalTitle = ref('Preview Surat Perjalanan Dinas')

const showCompleteModal = ref(false)
const completeFile = ref(null)
const completeFileInput = ref(null)
const completeError = ref('')
const statusProcessing = ref(false)

function openCompleteModal() {
  completeFile.value = null
  completeError.value = ''
  showCompleteModal.value = true
}

function onCompleteFileChange(e) {
  const file = e.target.files[0]
  if (file) {
    if (file.type !== 'application/pdf') {
      completeError.value = 'format file harus pdf'
      return
    }
    completeError.value = ''
    completeFile.value = file
  }
  e.target.value = ''
}

function onDropCompleteFile(e) {
  const file = e.dataTransfer.files[0]
  if (file) {
    if (file.type !== 'application/pdf') {
      completeError.value = 'format file harus pdf'
      return
    }
    completeError.value = ''
    completeFile.value = file
  }
}

async function handleStartDuty() {
  statusProcessing.value = true
  try {
    const res = await tripsApi.updateStatus(tripId, { status: 'ON_DUTY' })
    trip.value = res.data.data
    resetForm()
    toast.showToast('Status berhasil diubah menjadi Sedang Bertugas!', 'success')
  } catch (e) {
    toast.showToast(e.response?.data?.message || 'Gagal memulai dinas.', 'error')
  } finally {
    statusProcessing.value = false
  }
}

async function handleCompleteTrip() {
  if (!completeFile.value) {
    completeError.value = 'bukti perjalanan dinas harus diunggah'
    return
  }
  statusProcessing.value = true
  completeError.value = ''
  try {
    const fd = new FormData()
    fd.append('status', 'COMPLETED')
    fd.append('files', completeFile.value)
    
    const res = await tripsApi.updateStatus(tripId, fd)
    trip.value = res.data.data
    resetForm()
    showCompleteModal.value = false
    completeFile.value = null
    toast.showToast('Perjalanan dinas telah diselesaikan!', 'success')
  } catch (e) {
    completeError.value = e.response?.data?.message || 'Gagal menyelesaikan perjalanan dinas.'
  } finally {
    statusProcessing.value = false
  }
}

const form = ref({
  nomor_surat: '',
  destination: '',
  initiator: '',
  start_date: '',
  end_date: '',
  description: '',
  summary: ''
})

const attachments = computed(() => {
  if (!trip.value?.attachment_paths) return []
  return trip.value.attachment_paths.split(',').map(p => {
    const normalized = p.trim().replace(/\\/g, '/')
    return normalized.split('/').pop()
  }).filter(Boolean)
})

const canUpdate = computed(() => {
  if (!trip.value) return false
  return trip.value.status === 'PENDING' || !trip.value.attachment_paths
})

function onPdfFileChange(e) {
  const file = e.target.files[0]
  if (file) {
    if (file.type !== 'application/pdf') {
      toast.showToast('format file harus pdf', 'error')
      return
    }
    selectedPdfFile.value = file
  }
  e.target.value = ''
}

function onDropPdf(e) {
  const file = e.dataTransfer.files[0]
  if (file) {
    if (file.type !== 'application/pdf') {
      toast.showToast('format file harus pdf', 'error')
      return
    }
    selectedPdfFile.value = file
  }
}

function removeSelectedPdf() {
  selectedPdfFile.value = null
}

async function handlePreviewPdf() {
  if (attachments.value.length === 0) {
    toast.showToast('anda belum menyertakan surat perjalanan dinas anda', 'error')
    return
  }

  pdfModalTitle.value = 'Preview Surat Perjalanan Dinas'
  showPdfModal.value = true
  loadingPdf.value = true
  pdfUrl.value = ''

  try {
    const filename = attachments.value[0]
    const res = await tripsApi.downloadAttachment(tripId, filename)
    pdfUrl.value = URL.createObjectURL(res.data)
    loadingPdf.value = false
  } catch (e) {
    console.error(e)
    toast.showToast('Gagal memuat preview dokumen.', 'error')
    showPdfModal.value = false
    loadingPdf.value = false
  }
}

async function handlePreviewProofPdf() {
  if (attachments.value.length < 2) {
    toast.showToast('Bukti perjalanan dinas tidak ditemukan.', 'error')
    return
  }

  pdfModalTitle.value = 'Preview Bukti Perjalanan Dinas'
  showPdfModal.value = true
  loadingPdf.value = true
  pdfUrl.value = ''

  try {
    const filename = attachments.value[1]
    const res = await tripsApi.downloadAttachment(tripId, filename)
    pdfUrl.value = URL.createObjectURL(res.data)
    loadingPdf.value = false
  } catch (e) {
    console.error(e)
    toast.showToast('Gagal memuat preview dokumen.', 'error')
    showPdfModal.value = false
    loadingPdf.value = false
  }
}

function closePdfModal() {
  showPdfModal.value = false
  if (pdfUrl.value && pdfUrl.value.startsWith('blob:')) {
    URL.revokeObjectURL(pdfUrl.value)
  }
  pdfUrl.value = ''
}

function resetForm() {
  if (!trip.value) return
  form.value = {
    nomor_surat: trip.value.nomor_surat || '',
    destination: trip.value.destination || '',
    initiator: trip.value.initiator || '',
    start_date: trip.value.start_date ? new Date(trip.value.start_date).toISOString().split('T')[0] : '',
    end_date: trip.value.end_date ? new Date(trip.value.end_date).toISOString().split('T')[0] : '',
    description: trip.value.description || '',
    summary: trip.value.summary || ''
  }
  selectedPdfFile.value = null
}

async function downloadFile(filename) {
  try {
    const res = await tripsApi.downloadAttachment(tripId, filename)
    const url = URL.createObjectURL(res.data)
    const a = document.createElement('a')
    a.href = url; a.download = filename; a.click()
    URL.revokeObjectURL(url)
  } catch (e) { console.error(e) }
}

async function handleUpdateTrip() {
  editProcessing.value = true
  try {
    const fd = new FormData()
    Object.entries(form.value).forEach(([k, v]) => {
      fd.append(k, v)
    })
    
    if (selectedPdfFile.value) {
      fd.append('files', selectedPdfFile.value)
    }

    const res = await tripsApi.update(tripId, fd)
    trip.value = res.data.data
    selectedPdfFile.value = null
    resetForm()
    toast.showToast('Perjalanan dinas berhasil diperbarui!', 'success')
  } catch (e) {
    toast.showToast(e.response?.data?.message || 'Gagal memperbarui perjalanan dinas.', 'error')
  } finally {
    editProcessing.value = false
  }
}

async function handleDeleteTrip() {
  if (!confirm('Apakah Anda yakin ingin menghapus pengajuan perjalanan dinas ini?')) return
  deleteProcessing.value = true
  try {
    await tripsApi.remove(tripId)
    toast.showToast('Perjalanan dinas berhasil dihapus.', 'success')
    router.push('/trips')
  } catch (e) {
    toast.showToast(e.response?.data?.message || 'Gagal menghapus perjalanan dinas.', 'error')
  } finally {
    deleteProcessing.value = false
  }
}

onMounted(async () => {
  try {
    const tripRes = await tripsApi.getById(tripId)
    trip.value = tripRes.data.data || null
    resetForm()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
@reference "../../style.css";
.input-field {
  @apply w-full bg-slate-800 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition;
}
.input-field:disabled {
  @apply bg-slate-900/60 border-slate-800/80 text-slate-400 cursor-not-allowed;
}
</style>
