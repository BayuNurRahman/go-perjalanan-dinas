<template>
  <div>
    <p class="text-slate-400 text-xs mb-0.5">{{ label }}</p>
    <p class="text-white text-sm font-medium">{{ value || '-' }}</p>
  </div>
</template>
<script setup>
defineProps({ label: String, value: [String, Number] })
</script>
