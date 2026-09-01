import { defineStore } from 'pinia'
import { ElConfigProvider } from 'element-plus'

import zhCn from 'element-plus/es/locale/lang/zh-cn.mjs'
import enUs from 'element-plus/es/locale/lang/en.mjs'

export const useLangStore = defineStore('lang', {
  state() {
    return {
      lang: localStorage.getItem('lang') || 'zh'
    }
  },
  actions: {
    setLang(lang: 'zh' | 'en') {
      this.lang = lang
      localStorage.setItem('lang', lang)

      // 更新 Element Plus 组件库语言
      if(lang === 'zh'){
        ElConfigProvider.locale = zhCn
      }else{
        ElConfigProvider.locale = enUs
      }
    }
  }
})
