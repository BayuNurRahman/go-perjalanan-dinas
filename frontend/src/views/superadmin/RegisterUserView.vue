<template>
  <AppLayout>
    <div class="flex items-center gap-3 mb-6">
      <router-link to="/admin/users" class="text-slate-400 hover:text-white transition">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </router-link>
      <div>
        <h1 class="text-2xl font-bold text-white">Daftarkan User Baru</h1>
        <p class="text-slate-400 text-sm">Buat akun baru untuk karyawan</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Form -->
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6">
        <form @submit.prevent="handleRegister" class="space-y-5">
          <!-- Nama -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Nama Lengkap <span class="text-red-400">*</span></label>
            <input v-model="form.name" required placeholder="Nama lengkap karyawan" class="input-field" />
          </div>

          <!-- Email -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Email <span class="text-red-400">*</span></label>
            <input v-model="form.email" type="email" required placeholder="email@perusahaan.com" class="input-field" />
          </div>

          <!-- Password -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Password <span class="text-red-400">*</span></label>
            <input v-model="form.password" type="password" required placeholder="Minimal 8 karakter" minlength="8"
              class="input-field" />
          </div>

          <!-- Role -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Role <span class="text-red-400">*</span></label>
            <select v-model="form.role_id" required class="input-field">
              <option value="" disabled>Pilih role...</option>
              <option v-for="role in roles" :key="role.id" :value="role.id">{{ role.name }}</option>
            </select>
          </div>

          <!-- Departemen -->
          <div>
            <label class="block text-sm font-medium text-slate-300 mb-2">Departemen</label>
            <select v-model="form.department_id" class="input-field">
              <option :value="null">— Tidak ada —</option>
              <option v-for="dept in departments" :key="dept.id" :value="dept.id">{{ dept.name }}</option>
            </select>
          </div>

          <!-- Error / Success -->
          <div v-if="error" class="bg-red-900/40 border border-red-700 text-red-300 rounded-xl px-4 py-3 text-sm">{{ error }}</div>
          <div v-if="success" class="bg-green-900/40 border border-green-700 text-green-300 rounded-xl px-4 py-3 text-sm">
            ✓ {{ success }}
          </div>

          <!-- Submit -->
          <div class="flex gap-3 pt-2">
            <router-link to="/admin/users"
              class="flex-1 text-center bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium py-3 rounded-xl text-sm transition">
              Batal
            </router-link>
            <button type="submit" :disabled="loading"
              class="flex-1 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white font-semibold py-3 rounded-xl text-sm transition flex items-center justify-center gap-2">
              <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
              </svg>
              {{ loading ? 'Mendaftarkan...' : 'Daftarkan User' }}
            </button>
          </div>
        </form>
      </div>

      <!-- Info Panel -->
      <div class="space-y-4">
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-6">
          <h3 class="text-white font-semibold mb-4">Panduan Role Akses</h3>
          <div class="space-y-3">
            <div v-for="info in roleInfos" :key="info.role"
              class="flex gap-3 p-3 bg-slate-800/60 rounded-xl">
              <span :class="info.color" class="text-xs font-bold px-2 py-1 rounded-lg flex-shrink-0">
                {{ info.role }}
              </span>
              <p class="text-slate-400 text-xs leading-relaxed">{{ info.desc }}</p>
            </div>
          </div>
        </div>

        <div class="bg-blue-500/10 border border-blue-500/30 rounded-2xl p-5">
          <p class="text-blue-300 font-medium text-sm mb-1">💡 Tips</p>
          <p class="text-blue-400/80 text-xs leading-relaxed">
            Setelah user didaftarkan, mereka dapat langsung login menggunakan email dan password yang Anda tetapkan. Pastikan untuk memberitahukan password kepada user yang bersangkutan.
          </p>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import { authApi } from '@/api/auth'
import { rolesApi } from '@/api/roles'
import { departmentsApi } from '@/api/departments'

const loading = ref(false)
const error = ref('')
const success = ref('')
const roles = ref([])
const departments = ref([])

const form = ref({ name: '', email: '', password: '', role_id: '', department_id: null })

const roleInfos = [
  { role: 'EMPLOYEE',    color: 'bg-slate-700 text-slate-300', desc: 'Karyawan biasa. Dapat mengajukan perjalanan dinas dan klaim reimbursement.' },
  { role: 'MANAGER',    color: 'bg-green-700/50 text-green-300', desc: 'Manajer departemen. Dapat menyetujui atau menolak pengajuan karyawan di departemennya.' },
  { role: 'ADMIN_FIN',  color: 'bg-teal-700/50 text-teal-300', desc: 'Admin Keuangan. Dapat mereview klaim finansial dan melakukan pencairan dana.' },
  { role: 'SUPER_ADMIN',color: 'bg-purple-700/50 text-purple-300', desc: 'Super Admin. Akses penuh ke semua fitur termasuk manajemen user, departemen, dan role.' },
]

async function handleRegister() {
  loading.value = true
  error.value = ''
  success.value = ''
  try {
    await authApi.register(form.value)
    success.value = `User ${form.value.name} berhasil didaftarkan!`
    form.value = { name: '', email: '', password: '', role_id: '', department_id: null }
  } catch (e) {
    error.value = e.response?.data?.message || 'Gagal mendaftarkan user.'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const [rolesRes, deptsRes] = await Promise.all([rolesApi.getAll(), departmentsApi.getAll()])
    roles.value = rolesRes.data.data || []
    departments.value = deptsRes.data.data || []
  } catch (e) { console.error(e) }
})
</script>

<style scoped>
@reference "../../style.css";
.input-field {
  @apply w-full bg-slate-800 border border-slate-700 text-white placeholder-slate-500 rounded-xl px-4 py-3 text-sm focus:outline-none focus:border-blue-500 transition;
}
</style>
