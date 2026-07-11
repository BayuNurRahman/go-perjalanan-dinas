import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const theme = ref(localStorage.getItem('app-theme') || 'dark')
  const customFrom = ref(localStorage.getItem('app-theme-custom-from') || '#1e3a8a')
  const customTo = ref(localStorage.getItem('app-theme-custom-to') || '#0f172a')

  function applyTheme(newTheme = theme.value) {
    theme.value = newTheme
    localStorage.setItem('app-theme', newTheme)

    const root = document.documentElement
    
    // Clear custom theme style attributes
    root.removeAttribute('data-theme')
    root.style.removeProperty('--custom-gradient-from')
    root.style.removeProperty('--custom-gradient-to')
    root.style.removeProperty('--bg-main-gradient')

    if (newTheme === 'system') {
      const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      root.setAttribute('data-theme', isDark ? 'dark' : 'light')
    } else if (newTheme === 'custom') {
      root.setAttribute('data-theme', 'custom')
      root.style.setProperty('--custom-gradient-from', customFrom.value)
      root.style.setProperty('--custom-gradient-to', customTo.value)
      root.style.setProperty('--bg-main-gradient', `linear-gradient(to bottom, ${customFrom.value}, ${customTo.value})`)
      localStorage.setItem('app-theme-custom-from', customFrom.value)
      localStorage.setItem('app-theme-custom-to', customTo.value)
    } else {
      root.setAttribute('data-theme', newTheme)
    }
  }

  function setCustomColors(from, to) {
    customFrom.value = from
    customTo.value = to
    if (theme.value === 'custom') {
      applyTheme('custom')
    }
  }

  // Setup system listener
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  mediaQuery.addEventListener('change', () => {
    if (theme.value === 'system') {
      applyTheme('system')
    }
  })

  return { theme, customFrom, customTo, applyTheme, setCustomColors }
})
