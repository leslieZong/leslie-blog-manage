import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', {
  state: () => ({
    isDark: localStorage.getItem('isDark') === 'true',
  }),
  actions: {
    toggleDark() {
      localStorage.setItem('isDark', this.isDark ? 'false' : 'true')
      this.isDark = !this.isDark
      if (this.isDark) {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }
    },
    initTheme() {
      if (this.isDark) document.documentElement.classList.add('dark')
    },
  },
})
