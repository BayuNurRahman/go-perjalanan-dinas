<template>
  <AppLayout>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-white">Kelola User</h1>
        <p class="text-slate-400 text-sm mt-1">Manajemen seluruh akun karyawan</p>
      </div>
      <router-link to="/admin/register"
        class="inline-flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white font-semibold px-4 py-2.5 rounded-xl text-sm transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
        </svg>
        Daftarkan User
      </router-link>
    </div>

    <!-- Search -->
    <div class="mb-4">
      <input v-model="search" @input="onSearch" type="text" placeholder="Cari nama atau email..."
        class="w-full md:w-80 bg-slate-900 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-blue-500 transition" />
    </div>

    <!-- Table -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <div v-if="loading" class="p-6 space-y-3">
        <div v-for="i in 6" :key="i" class="h-14 bg-slate-800 rounded-xl animate-pulse" />
      </div>

      <div v-else-if="filteredUsers.length === 0" class="p-16 text-center">
        <p class="text-slate-500">Tidak ada user ditemukan</p>
      </div>

      <table v-else class="w-full text-sm">
        <thead class="bg-slate-800/60 border-b border-slate-800">
          <tr>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Nama</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Email</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Role</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Departemen</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="user in filteredUsers" :key="user.id" class="hover:bg-slate-800/40 transition">
            <td class="px-6 py-4 text-white font-medium">{{ user.name }}</td>
            <td class="px-6 py-4 text-slate-300">{{ user.email }}</td>
            <td class="px-6 py-4">
              <span class="bg-blue-600/15 text-blue-400 text-xs font-semibold px-2.5 py-1 rounded-full">
                {{ user.role?.name || '-' }}
              </span>
            </td>
            <td class="px-6 py-4 text-slate-400">{{ user.department?.name || '-' }}</td>
            <td class="px-6 py-4">
              <div class="flex gap-3">
                <button @click="openEdit(user)"
                  class="text-blue-400 hover:text-blue-300 text-xs font-medium transition">Edit</button>
                <button @click="openDelete(user)"
                  class="text-red-400 hover:text-red-300 text-xs font-medium transition">Hapus</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Edit Modal -->
    <div v-if="editModal.open" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center px-4">
      <div class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-md">
        <h3 class="text-white font-semibold mb-4">Edit User</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-xs font-medium text-slate-400 mb-1.5">Nama</label>
            <input v-model="editModal.form.name" class="input-field" placeholder="Nama lengkap" />
          </div>
          <div>
            <label class="block text-xs font-medium text-slate-400 mb-1.5">Email</label>
            <input v-model="editModal.form.email" type="email" class="input-field" placeholder="Email" />
          </div>
          <div>
            <label class="block text-xs font-medium text-slate-400 mb-1.5">Role ID</label>
            <select v-model="editModal.form.role_id" class="input-field">
              <option v-for="r in roles" :key="r.id" :value="r.id">{{ r.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-medium text-slate-400 mb-1.5">Departemen</label>
            <select v-model="editModal.form.department_id" class="input-field">
              <option :value="null">— Tidak ada —</option>
              <option v-for="d in departments" :key="d.id" :value="d.id">{{ d.name }}</option>
            </select>
          </div>
          <div v-if="editModal.error" class="bg-red-900/40 border border-red-700 text-red-300 rounded-xl px-4 py-2.5 text-sm">
            {{ editModal.error }}
          </div>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="editModal.open = false"
            class="flex-1 bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium py-2.5 rounded-xl text-sm transition">
            Batal
          </button>
          <button @click="confirmEdit" :disabled="editModal.loading"
            class="flex-1 bg-blue-600 hover:bg-blue-700 disabled:opacity-40 text-white font-semibold py-2.5 rounded-xl text-sm transition">
            {{ editModal.loading ? 'Menyimpan...' : 'Simpan' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirm Modal -->
    <div v-if="deleteModal.open" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center px-4">
      <div class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-sm text-center">
        <div class="w-12 h-12 bg-red-600/20 rounded-2xl flex items-center justify-center mx-auto mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </div>
        <h3 class="text-white font-semibold mb-1">Hapus User?</h3>
        <p class="text-slate-400 text-sm mb-5">
          <strong class="text-white">{{ deleteModal.user?.name }}</strong> akan dihapus secara permanen.
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
import { ref, computed, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import { usersApi } from '@/api/users'
import { rolesApi } from '@/api/roles'
import { departmentsApi } from '@/api/departments'

const loading = ref(true)
const users = ref([])
const roles = ref([])
const departments = ref([])
const search = ref('')
let searchTimer = null

const filteredUsers = computed(() => {
  if (!search.value) return users.value
  const q = search.value.toLowerCase()
  return users.value.filter(u => u.name?.toLowerCase().includes(q) || u.email?.toLowerCase().includes(q))
})

const editModal = ref({ open: false, user: null, form: {}, loading: false, error: '' })
const deleteModal = ref({ open: false, user: null, loading: false })

function onSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {}, 300)
}

function openEdit(user) {
  editModal.value = {
    open: true, user,
    form: { name: user.name, email: user.email, role_id: user.role?.id, department_id: user.department?.id || null },
    loading: false, error: ''
  }
}

async function confirmEdit() {
  editModal.value.loading = true
  editModal.value.error = ''
  try {
    const res = await usersApi.update(editModal.value.user.id, editModal.value.form)
    const idx = users.value.findIndex(u => u.id === editModal.value.user.id)
    if (idx !== -1) users.value[idx] = res.data.data || { ...users.value[idx], ...editModal.value.form }
    editModal.value.open = false
  } catch (e) {
    editModal.value.error = e.response?.data?.message || 'Gagal menyimpan perubahan.'
  } finally {
    editModal.value.loading = false
  }
}

function openDelete(user) {
  deleteModal.value = { open: true, user, loading: false }
}

async function confirmDelete() {
  deleteModal.value.loading = true
  try {
    await usersApi.remove(deleteModal.value.user.id)
    users.value = users.value.filter(u => u.id !== deleteModal.value.user.id)
    deleteModal.value.open = false
  } catch (e) { console.error(e) }
  finally { deleteModal.value.loading = false }
}

onMounted(async () => {
  try {
    const [usersRes, rolesRes, deptsRes] = await Promise.all([
      usersApi.getAll(),
      rolesApi.getAll(),
      departmentsApi.getAll(),
    ])
    users.value = usersRes.data.data || []
    roles.value = rolesRes.data.data || []
    departments.value = deptsRes.data.data || []
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
