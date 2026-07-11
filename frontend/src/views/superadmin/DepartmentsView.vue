<template>
  <AppLayout>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-white">Kelola Departemen</h1>
        <p class="text-slate-400 text-sm mt-1">Tambah, edit, dan hapus departemen</p>
      </div>
      <button @click="openCreate"
        class="inline-flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white font-semibold px-4 py-2.5 rounded-xl text-sm transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Tambah Departemen
      </button>
    </div>

    <!-- Table -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <div v-if="loading" class="p-6 space-y-3">
        <div v-for="i in 4" :key="i" class="h-14 bg-slate-800 rounded-xl animate-pulse" />
      </div>

      <div v-else-if="departments.length === 0" class="p-16 text-center">
        <p class="text-slate-500">Belum ada departemen</p>
      </div>

      <table v-else class="w-full text-sm">
        <thead class="bg-slate-800/60 border-b border-slate-800">
          <tr>
            <th class="text-left text-slate-400 font-medium px-6 py-3">ID</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Nama Departemen</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Kode</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Deskripsi</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="dept in departments" :key="dept.id" class="hover:bg-slate-800/40 transition">
            <td class="px-6 py-4 text-slate-500 text-xs">#{{ dept.id }}</td>
            <td class="px-6 py-4 text-white font-medium">{{ dept.name }}</td>
            <td class="px-6 py-4">
              <span class="font-mono text-xs text-slate-300 bg-slate-800 px-2 py-1 rounded-lg">{{ dept.code || '-' }}</span>
            </td>
            <td class="px-6 py-4 text-slate-400 text-sm">{{ dept.description || '-' }}</td>
            <td class="px-6 py-4">
              <div class="flex gap-3">
                <button @click="openEdit(dept)" class="text-blue-400 hover:text-blue-300 text-xs font-medium transition">Edit</button>
                <button @click="openDelete(dept)" class="text-red-400 hover:text-red-300 text-xs font-medium transition">Hapus</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create / Edit Modal -->
    <div v-if="modal.open" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center px-4">
      <div class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-md">
        <h3 class="text-white font-semibold mb-4">{{ modal.isEdit ? 'Edit Departemen' : 'Tambah Departemen' }}</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-xs font-medium text-slate-400 mb-1.5">Nama Departemen *</label>
            <input v-model="modal.form.name" class="input-field" placeholder="cth. Information Technology" />
          </div>
          <div>
            <label class="block text-xs font-medium text-slate-400 mb-1.5">Kode</label>
            <input v-model="modal.form.code" class="input-field" placeholder="cth. IT" />
          </div>
          <div>
            <label class="block text-xs font-medium text-slate-400 mb-1.5">Deskripsi</label>
            <textarea v-model="modal.form.description" rows="2" class="input-field resize-none"
              placeholder="Deskripsi singkat departemen..."></textarea>
          </div>
          <div v-if="modal.error" class="bg-red-900/40 border border-red-700 text-red-300 rounded-xl px-4 py-2.5 text-sm">
            {{ modal.error }}
          </div>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="modal.open = false"
            class="flex-1 bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium py-2.5 rounded-xl text-sm transition">
            Batal
          </button>
          <button @click="confirmSave" :disabled="!modal.form.name || modal.loading"
            class="flex-1 bg-blue-600 hover:bg-blue-700 disabled:opacity-40 text-white font-semibold py-2.5 rounded-xl text-sm transition">
            {{ modal.loading ? 'Menyimpan...' : (modal.isEdit ? 'Simpan' : 'Tambah') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Modal -->
    <div v-if="deleteModal.open" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center px-4">
      <div class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-sm text-center">
        <div class="w-12 h-12 bg-red-600/20 rounded-2xl flex items-center justify-center mx-auto mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </div>
        <h3 class="text-white font-semibold mb-1">Hapus Departemen?</h3>
        <p class="text-slate-400 text-sm mb-5">
          <strong class="text-white">{{ deleteModal.dept?.name }}</strong> akan dihapus.
        </p>
        <div class="flex gap-3">
          <button @click="deleteModal.open = false"
            class="flex-1 bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium py-2.5 rounded-xl text-sm transition">
            Batal
          </button>
          <button @click="confirmDelete" :disabled="deleteModal.loading"
            class="flex-1 bg-red-600 hover:bg-red-700 disabled:opacity-40 text-white font-semibold py-2.5 rounded-xl text-sm transition">
            {{ deleteModal.loading ? 'Menghapus...' : 'Hapus' }}
          </button>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import { departmentsApi } from '@/api/departments'

const loading = ref(true)
const departments = ref([])

const modal = ref({ open: false, isEdit: false, dept: null, form: {}, loading: false, error: '' })
const deleteModal = ref({ open: false, dept: null, loading: false })

function openCreate() {
  modal.value = { open: true, isEdit: false, dept: null, form: { name: '', code: '', description: '' }, loading: false, error: '' }
}

function openEdit(dept) {
  modal.value = { open: true, isEdit: true, dept, form: { name: dept.name, code: dept.code || '', description: dept.description || '' }, loading: false, error: '' }
}

async function confirmSave() {
  modal.value.loading = true
  modal.value.error = ''
  try {
    if (modal.value.isEdit) {
      const res = await departmentsApi.update(modal.value.dept.id, modal.value.form)
      const idx = departments.value.findIndex(d => d.id === modal.value.dept.id)
      if (idx !== -1) departments.value[idx] = res.data.data || { ...departments.value[idx], ...modal.value.form }
    } else {
      const res = await departmentsApi.create(modal.value.form)
      departments.value.push(res.data.data)
    }
    modal.value.open = false
  } catch (e) {
    modal.value.error = e.response?.data?.message || 'Gagal menyimpan.'
  } finally {
    modal.value.loading = false
  }
}

function openDelete(dept) {
  deleteModal.value = { open: true, dept, loading: false }
}

async function confirmDelete() {
  deleteModal.value.loading = true
  try {
    await departmentsApi.remove(deleteModal.value.dept.id)
    departments.value = departments.value.filter(d => d.id !== deleteModal.value.dept.id)
    deleteModal.value.open = false
  } catch (e) { console.error(e) }
  finally { deleteModal.value.loading = false }
}

onMounted(async () => {
  try {
    const res = await departmentsApi.getAll()
    departments.value = res.data.data || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
})
</script>

<style scoped>
@reference "../../style.css";
.input-field {
  @apply w-full bg-slate-800 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-blue-500 transition;
}
</style>
