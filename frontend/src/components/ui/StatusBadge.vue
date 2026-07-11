<template>
  <span :class="[bgClass, textClass]" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold">
    {{ label }}
  </span>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({ status: { type: String, default: '' } })

const statusMap = {
  PENDING:   { bg: 'bg-yellow-500/15', text: 'text-yellow-400', label: 'Menunggu' },
  APPROVED:  { bg: 'bg-green-500/15',  text: 'text-green-400',  label: 'Disetujui' },
  REJECTED:  { bg: 'bg-red-500/15',    text: 'text-red-400',    label: 'Ditolak' },
  ON_DUTY:   { bg: 'bg-blue-500/15',   text: 'text-blue-400',   label: 'Dalam Perjalanan' },
  COMPLETED: { bg: 'bg-purple-500/15', text: 'text-purple-400', label: 'Selesai' },
  REVISION_REQUESTED: { bg: 'bg-orange-500/15', text: 'text-orange-400', label: 'Revisi' },
  // Claim statuses
  DISBURSED: { bg: 'bg-teal-500/15',   text: 'text-teal-400',   label: 'Dicairkan' },
}

const current = computed(() => statusMap[props.status] || { bg: 'bg-slate-500/15', text: 'text-slate-400', label: props.status })
const bgClass = computed(() => current.value.bg)
const textClass = computed(() => current.value.text)
const label = computed(() => current.value.label)
</script>
