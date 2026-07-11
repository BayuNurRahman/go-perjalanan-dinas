<template>
  <div>
    <!-- Top Progress Loading Bar -->
    <div
      v-if="loadingStore.isLoading"
      class="fixed top-0 left-0 right-0 h-1 bg-gradient-to-r from-blue-500 via-indigo-500 to-purple-500 z-50 animate-pulse"
    ></div>

    <router-view />
    <ToastContainer />
  </div>
</template>

<script setup>
import { useLoadingStore } from '@/stores/loading'
import ToastContainer from '@/components/ui/ToastContainer.vue'

const loadingStore = useLoadingStore()
</script>

<style>
/* Global fade transitions or styling if needed */
</style>
