<template>
  <AppLayout>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-white">Kelola Role</h1>
        <p class="text-slate-400 text-sm mt-1">Manajemen hak akses sistem</p>
      </div>
      <button @click="openCreate"
        class="inline-flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white font-semibold px-4 py-2.5 rounded-xl text-sm transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Tambah Role
      </button>
    </div>

    <!-- Table -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <div v-if="loading" class="p-6 space-y-3">
        <div v-for="i in 4" :key="i" class="h-14 bg-slate-800 rounded-xl animate-pulse" />
      </div>

      <div v-else-if="roles.length === 0" class="p-16 text-center">
        <p class="text-slate-500">Belum ada role</p>
      </div>

      <table v-else class="w-full text-sm">
        <thead class="bg-slate-800/60 border-b border-slate-800">
          <tr>
            <th class="text-left text-slate-400 font-medium px-6 py-3">ID</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Nama Role</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Deskripsi</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="role in roles" :key="role.id" class="hover:bg-slate-800/40 transition">
            <td class="px-6 py-4 text-slate-500 text-xs">#{{ role.id }}</td>
            <td class="px-6 py-4">
              <span class="bg-purple-600/15 text-purple-300 font-semibold text-xs px-2.5 py-1 rounded-full">
                {{ role.name }}
              </span>
            </td>
            <td class="px-6 py-4 text-slate-400 text-sm">{{ role.description || '-' }}</td>
            <td class="px-6 py-4">
              <div class="flex gap-3">
                <button @click="openEdit(role)" class="text-blue-400 hover:text-blue-300 text-xs font-medium transition">Edit</button>
                <button @click="openDelete(role)" class="text-red-400 hover:text-red-300 text-xs font-medium transition">Hapus</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create / Edit Modal -->
    <div v-if="modal.open" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center px-4">
      <div class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-md">
        <h3 class="text-white font-semibold mb-4">{{ modal.isEdit ? 'Edit Role' : 'Tambah Role' }}</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-xs font-medium text-slate-400 mb-1.5">Nama Role *</label>
            <input v-model="modal.form.name" class="input-field" placeholder="cth. MANAGER" />
          </div>
          <div>
            <label class="block text-xs font-medium text-slate-400 mb-1.5">Deskripsi</label>
            <textarea v-model="modal.form.description" rows="2" class="input-field resize-none"
              placeholder="Deskripsi hak akses role ini..."></textarea>
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
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <h3 class="text-white font-semibold mb-1">Hapus Role?</h3>
        <p class="text-slate-400 text-sm mb-5">
          Role <strong class="text-white">{{ deleteModal.role?.name }}</strong> akan dihapus.
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
import { rolesApi } from '@/api/roles'

const loading = ref(true)
const roles = ref([])
const modal = ref({ open: false, isEdit: false, role: null, form: {}, loading: false, error: '' })
const deleteModal = ref({ open: false, role: null, loading: false })

function openCreate() {
  modal.value = { open: true, isEdit: false, role: null, form: { name: '', description: '' }, loading: false, error: '' }
}
function openEdit(role) {
  modal.value = { open: true, isEdit: true, role, form: { name: role.name, description: role.description || '' }, loading: false, error: '' }
}

async function confirmSave() {
  modal.value.loading = true
  modal.value.error = ''
  try {
    if (modal.value.isEdit) {
      const res = await rolesApi.update(modal.value.role.id, modal.value.form)
      const idx = roles.value.findIndex(r => r.id === modal.value.role.id)
      if (idx !== -1) roles.value[idx] = res.data.data || { ...roles.value[idx], ...modal.value.form }
    } else {
      const res = await rolesApi.create(modal.value.form)
      roles.value.push(res.data.data)
    }
    modal.value.open = false
  } catch (e) {
    modal.value.error = e.response?.data?.message || 'Gagal menyimpan.'
  } finally { modal.value.loading = false }
}

function openDelete(role) {
  deleteModal.value = { open: true, role, loading: false }
}
async function confirmDelete() {
  deleteModal.value.loading = true
  try {
    await rolesApi.remove(deleteModal.value.role.id)
    roles.value = roles.value.filter(r => r.id !== deleteModal.value.role.id)
    deleteModal.value.open = false
  } catch (e) { console.error(e) }
  finally { deleteModal.value.loading = false }
}

onMounted(async () => {
  try {
    const res = await rolesApi.getAll()
    roles.value = res.data.data || []
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
