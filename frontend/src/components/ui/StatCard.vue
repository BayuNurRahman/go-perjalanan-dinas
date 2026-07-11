<template>
  <div class="bg-slate-900 border border-slate-800 rounded-2xl p-5">
    <div class="flex items-start justify-between">
      <div>
        <p class="text-slate-400 text-sm">{{ label }}</p>
        <p class="text-3xl font-bold text-white mt-1">{{ value }}</p>
      </div>
      <div :class="iconBg" class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0">
        <svg xmlns="http://www.w3.org/2000/svg" :class="iconColor" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="icon" />
        </svg>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({
  label: String,
  value: [Number, String],
  icon: String,
  color: { type: String, default: 'blue' },
})
const colorMap = {
  blue:   { bg: 'bg-blue-600/20',   text: 'text-blue-400' },
  green:  { bg: 'bg-green-600/20',  text: 'text-green-400' },
  yellow: { bg: 'bg-yellow-600/20', text: 'text-yellow-400' },
  purple: { bg: 'bg-purple-600/20', text: 'text-purple-400' },
  red:    { bg: 'bg-red-600/20',    text: 'text-red-400' },
}
const iconBg = computed(() => colorMap[props.color]?.bg || colorMap.blue.bg)
const iconColor = computed(() => colorMap[props.color]?.text || colorMap.blue.text)
</script>
