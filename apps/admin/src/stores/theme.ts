import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', {
  state: () => ({
    isDark: localStorage.getItem('isDark') === 'true',
  }),
  actions: {
    toggleDark() {
      this.isDark = !this.isDark
      localStorage.setItem('isDark', String(this.isDark))
      this.applyTheme()
    },
    applyTheme() {
      if (this.isDark) {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }
    },
    initTheme() {
      this.applyTheme()
    },
  },
})
