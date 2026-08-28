import { fileURLToPath, URL } from 'node:url'

import { defineConfig, type PluginOption } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// 生产环境 CSP：禁止内联脚本与外部源、禁 object/iframe 嵌入；
// style-src 放开 'unsafe-inline' 以兼容 element-plus 运行时注入的 <style>
const CSP = [
  "default-src 'self'",
  "script-src 'self'",
  "style-src 'self' 'unsafe-inline'",
  "img-src 'self' data: https:",
  "font-src 'self' data:",
  "connect-src 'self' https:",
  "base-uri 'self'",
  "form-action 'self'",
  "object-src 'none'",
  "frame-ancestors 'none'",
].join('; ')

// 仅在 build 注入 CSP meta，避免破坏 dev 的 HMR / vue-devtools（依赖 eval / 内联脚本）
function cspMetaPlugin(): PluginOption {
  return {
    name: 'inject-csp-meta',
    apply: 'build',
    transformIndexHtml(html) {
      const tag = `<meta http-equiv="Content-Security-Policy" content="${CSP}">`
      return html.replace('<meta charset="UTF-8">', `<meta charset="UTF-8">\n    ${tag}`)
    },
  }
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), vueDevTools(), cspMetaPlugin()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
