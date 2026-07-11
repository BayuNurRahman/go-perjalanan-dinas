import { ref } from 'vue'

const toasts = ref([])
let nextId = 0

export function useToast() {
  function add(message, type = 'info', duration = 3500) {
    const id = ++nextId
    toasts.value.push({ id, message, type })
    setTimeout(() => remove(id), duration)
  }

  function remove(id) {
    const idx = toasts.value.findIndex(t => t.id === id)
    if (idx !== -1) toasts.value.splice(idx, 1)
  }

  return {
    toasts,
    success: (msg, d) => add(msg, 'success', d),
    error:   (msg, d) => add(msg, 'error', d),
    info:    (msg, d) => add(msg, 'info', d),
    warning: (msg, d) => add(msg, 'warning', d),
    remove,
  }
}
