import { createI18n } from 'vue-i18n'
import zh from './locales/zh'
import en from './locales/en'
// element-plus 内置语言包
import zhCn from 'element-plus/es/locale/lang/zh-cn.mjs'
import enUs from 'element-plus/es/locale/lang/en.mjs'

// 读取本地缓存
const localLang = localStorage.getItem('lang') || 'zh'

const i18n = createI18n({
  legacy: false, // ✅ setup语法糖必须关闭legacy
  locale: localLang,
  fallbackLocale: 'zh',
  messages: {
    zh: {
      ...zhCn, // 把element中文合并进来
      ...zh,
    },
    en: {
      ...enUs, // element英文合并进来
      ...en,
    },
  },
})

export default i18n
