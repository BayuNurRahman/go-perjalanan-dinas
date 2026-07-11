<template>
  <AppLayout>
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-white">Distribusi Tim</h1>
      <p class="text-slate-400 text-sm mt-1">Ringkasan perjalanan dinas per anggota tim</p>
    </div>

    <!-- Summary Stats -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
      <StatCard label="Total Anggota" :value="distribution.total_members ?? 0" color="blue"
        icon="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
      <StatCard label="Sedang Tugas" :value="distribution.currently_on_duty ?? 0" color="purple"
        icon="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      <StatCard label="Total Trip Bulan Ini" :value="distribution.total_trips_this_month ?? 0" color="green"
        icon="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
      <StatCard label="Pending Approval" :value="distribution.pending_approval ?? 0" color="yellow"
        icon="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
    </div>

    <!-- Per Member Table -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <div class="px-6 py-4 border-b border-slate-800">
        <h2 class="text-white font-semibold">Detail Per Anggota</h2>
      </div>

      <div v-if="loading" class="p-6 space-y-3">
        <div v-for="i in 5" :key="i" class="h-14 bg-slate-800 rounded-xl animate-pulse" />
      </div>

      <div v-else-if="!members.length" class="p-12 text-center">
        <p class="text-slate-500 text-sm">Tidak ada data anggota tim</p>
      </div>

      <table v-else class="w-full text-sm">
        <thead class="bg-slate-800/60 border-b border-slate-800">
          <tr>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Anggota</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Total Trip</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Pending</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Disetujui</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Selesai</th>
            <th class="text-left text-slate-400 font-medium px-6 py-3">Status</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="member in members" :key="member.user_id" class="hover:bg-slate-800/40 transition">
            <td class="px-6 py-4">
              <p class="text-white font-medium">{{ member.name }}</p>
              <p class="text-slate-400 text-xs">{{ member.email }}</p>
            </td>
            <td class="px-6 py-4 text-slate-300 font-semibold">{{ member.total_trips ?? 0 }}</td>
            <td class="px-6 py-4">
              <span class="text-yellow-400 font-medium">{{ member.pending ?? 0 }}</span>
            </td>
            <td class="px-6 py-4">
              <span class="text-green-400 font-medium">{{ member.approved ?? 0 }}</span>
            </td>
            <td class="px-6 py-4">
              <span class="text-purple-400 font-medium">{{ member.completed ?? 0 }}</span>
            </td>
            <td class="px-6 py-4">
              <span v-if="member.on_duty > 0"
                class="inline-flex items-center gap-1 text-blue-400 text-xs font-medium">
                <span class="w-1.5 h-1.5 bg-blue-400 rounded-full animate-pulse"></span>
                Sedang Bertugas
              </span>
              <span v-else class="text-slate-500 text-xs">Di Kantor</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import StatCard from '@/components/ui/StatCard.vue'
import { tripsApi } from '@/api/trips'

const loading = ref(true)
const distribution = ref({})
const members = computed(() => distribution.value.members || [])

onMounted(async () => {
  try {
    const res = await tripsApi.getTeamDistribution()
    distribution.value = res.data.data || {}
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>
