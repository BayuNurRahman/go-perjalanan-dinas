<template>
  <div class="min-h-screen bg-slate-950 flex">
    <!-- Sidebar -->
    <Sidebar :isOpen="isSidebarOpen" @close="isSidebarOpen = false" />

    <!-- Main Content Wrapper -->
    <div class="flex-1 flex flex-col min-h-screen md:ml-64 transition-all duration-300">
      <!-- Top Navbar -->
      <header class="h-16 bg-slate-900 border-b border-slate-800 flex items-center justify-between px-4 md:px-6 sticky top-0 z-10">
        <div class="flex items-center gap-3">
          <!-- Toggle Hamburger Button (Mobile Only) -->
          <button
            @click="isSidebarOpen = !isSidebarOpen"
            class="md:hidden text-slate-400 hover:text-white p-1 rounded-lg hover:bg-slate-800 transition"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          
          <h2 class="text-white font-semibold text-base md:text-lg">{{ pageTitle }}</h2>
        </div>

        <div class="flex items-center gap-3">
          <div class="flex flex-col items-end text-right hidden sm:flex">
            <span class="text-slate-200 text-sm font-medium">{{ authStore.userName }}</span>
            <span v-if="authStore.departmentName" class="text-slate-500 text-xs">Dept: {{ authStore.departmentName }}</span>
          </div>
          <span class="bg-blue-600/20 text-blue-400 text-xs font-medium px-2.5 py-1 rounded-full">{{ authStore.role }}</span>
          
          <!-- Theme Switcher Button -->
          <button @click="isThemeModalOpen = true" class="text-slate-400 hover:text-white p-1.5 rounded-lg hover:bg-slate-800 transition" title="Pilih Tema">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01" />
            </svg>
          </button>

          <button @click="handleLogout" class="text-slate-400 hover:text-white text-sm transition ml-1">Logout</button>
        </div>
      </header>

      <!-- Page Content -->
      <main class="flex-1 p-4 md:p-6">
        <slot />
      </main>
    </div>

    <!-- Theme Switcher Modal -->
    <div v-if="isThemeModalOpen" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 shadow-2xl relative">
        <h3 class="text-lg font-bold text-white mb-4">Pengaturan Tema</h3>
        
        <div class="space-y-4">
          <!-- Theme Preset Options -->
          <div class="grid grid-cols-2 gap-3">
            <button @click="setTheme('light')" :class="[themeStore.theme === 'light' ? 'border-blue-500 bg-blue-500/10 text-blue-400' : 'border-slate-800 hover:border-slate-700 bg-slate-950 text-slate-300', 'border rounded-xl p-3 text-sm font-medium transition text-center']">
              ☀️ Light Theme
            </button>
            <button @click="setTheme('dark')" :class="[themeStore.theme === 'dark' ? 'border-blue-500 bg-blue-500/10 text-blue-400' : 'border-slate-800 hover:border-slate-700 bg-slate-950 text-slate-300', 'border rounded-xl p-3 text-sm font-medium transition text-center']">
              🌙 Dark Theme
            </button>
            <button @click="setTheme('system')" :class="[themeStore.theme === 'system' ? 'border-blue-500 bg-blue-500/10 text-blue-400' : 'border-slate-800 hover:border-slate-700 bg-slate-950 text-slate-300', 'border rounded-xl p-3 text-sm font-medium transition text-center']">
              💻 Sesuai Sistem
            </button>
            <button @click="setTheme('custom')" :class="[themeStore.theme === 'custom' ? 'border-blue-500 bg-blue-500/10 text-blue-400' : 'border-slate-800 hover:border-slate-700 bg-slate-950 text-slate-300', 'border rounded-xl p-3 text-sm font-medium transition text-center']">
              🎨 Custom Gradient
            </button>
          </div>

          <!-- Custom Colors Settings (Visible only when theme is 'custom') -->
          <div v-if="themeStore.theme === 'custom'" class="bg-slate-950 border border-slate-800 rounded-xl p-4 space-y-4">
            <p class="text-xs font-semibold text-slate-400 uppercase tracking-wider">Kustomisasi Gradasi</p>
            
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs text-slate-400 mb-1">Warna Atas (Top)</label>
                <div class="flex items-center gap-2">
                  <input type="color" v-model="tempCustomFrom" @input="updateCustomColors" class="w-8 h-8 rounded border-0 bg-transparent cursor-pointer" />
                  <input type="text" v-model="tempCustomFrom" @input="updateCustomColors" class="flex-1 bg-slate-900 border border-slate-800 rounded-lg px-2 py-1 text-xs text-white font-mono uppercase" />
                </div>
              </div>
              
              <div>
                <label class="block text-xs text-slate-400 mb-1">Warna Bawah (Bottom)</label>
                <div class="flex items-center gap-2">
                  <input type="color" v-model="tempCustomTo" @input="updateCustomColors" class="w-8 h-8 rounded border-0 bg-transparent cursor-pointer" />
                  <input type="text" v-model="tempCustomTo" @input="updateCustomColors" class="flex-1 bg-slate-900 border border-slate-800 rounded-lg px-2 py-1 text-xs text-white font-mono uppercase" />
                </div>
              </div>
            </div>

            <!-- Color Palette Presets -->
            <div class="space-y-2 pt-2">
              <label class="block text-xs text-slate-400">Pilihan Palet Preset:</label>
              <div class="flex flex-wrap gap-2">
                <button v-for="preset in colorPresets" :key="preset.name" 
                  @click="applyPreset(preset)" 
                  class="flex items-center gap-1.5 px-2 py-1 bg-slate-900 hover:bg-slate-800 border border-slate-800 hover:border-slate-700 rounded-lg text-[10px] text-slate-300 transition">
                  <span class="w-3 h-3 rounded-full" :style="{ background: `linear-gradient(to right, ${preset.from}, ${preset.to})` }"></span>
                  {{ preset.name }}
                </button>
              </div>
            </div>

            <!-- Live Preview -->
            <div class="pt-3">
              <label class="block text-xs text-slate-400 mb-1.5">Preview Gradasi:</label>
              <div class="h-16 rounded-xl border border-slate-800 overflow-hidden"
                :style="{ background: `linear-gradient(to bottom, ${tempCustomFrom}, ${tempCustomTo})` }">
              </div>
            </div>
          </div>
        </div>

        <div class="flex justify-end gap-3 mt-6">
          <button @click="isThemeModalOpen = false" class="bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium px-4 py-2 rounded-xl transition">
            Selesai
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'
import { useThemeStore } from '@/stores/theme'
import Sidebar from './Sidebar.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const themeStore = useThemeStore()

const isSidebarOpen = ref(false)
const isThemeModalOpen = ref(false)

const tempCustomFrom = ref(themeStore.customFrom)
const tempCustomTo = ref(themeStore.customTo)

const pageTitle = computed(() => route.meta?.title || 'Perjalanan Dinas')

const colorPresets = [
  { name: 'Deep Blue', from: '#1e3a8a', to: '#0f172a' },
  { name: 'Midnight Purple', from: '#4c1d95', to: '#0f172a' },
  { name: 'Sunset Rose', from: '#881337', to: '#1e1b4b' },
  { name: 'Forest Moss', from: '#064e3b', to: '#022c22' },
  { name: 'Burgundy', from: '#581c87', to: '#180020' },
  { name: 'Cyberpunk', from: '#701a75', to: '#0f172a' }
]

function setTheme(themeName) {
  themeStore.applyTheme(themeName)
}

function updateCustomColors() {
  themeStore.setCustomColors(tempCustomFrom.value, tempCustomTo.value)
}

function applyPreset(preset) {
  tempCustomFrom.value = preset.from
  tempCustomTo.value = preset.to
  themeStore.setCustomColors(preset.from, preset.to)
}

async function handleLogout() {
  try {
    await authApi.logout()
  } catch (_) {}
  authStore.clearAuth()
  router.push('/login')
}
</script>
