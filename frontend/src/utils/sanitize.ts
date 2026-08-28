import DOMPurify from 'dompurify'

// 允许的 URL 协议白名单；其余（含 javascript:、data: 等）一律视为不安全
const SAFE_URL_PROTOCOLS = ['http:', 'https:', 'mailto:', 'tel:']

// 安全的相对/锚点形式：以 / 或 # 开头，或 data:image（仅限图片）
function isSafeRelative(url: string): boolean {
  return /^(\/|#)/.test(url) || /^data:image\//i.test(url)
}

// 判定 URL 是否可安全用于 href/src/window.open
export function isSafeUrl(url: string): boolean {
  if (!url) return false
  const trimmed = url.trim().toLowerCase()
  if (isSafeRelative(trimmed)) return true
  // 禁止协议相对 URL（//evil.com）与控制字符拼接绕过
  if (trimmed.startsWith('//')) return false
  try {
    const { protocol } = new URL(trimmed, window.location.origin)
    return SAFE_URL_PROTOCOLS.includes(protocol)
  } catch {
    return false
  }
}

// 返回安全 URL；不安全时回退到 about:blank，避免拼接出 javascript: 等
export function sanitizeUrl(url: string): string {
  return isSafeUrl(url) ? url : 'about:blank'
}

// DOMPurify 封装：统一出口，便于后续扩展全局配置（如 ALLOWED_TAGS）
// 用 Parameters 提取 sanitize 的第二参数类型，兼容 dompurify 自带类型与 @types/dompurify
type SanitizeConfig = Parameters<typeof DOMPurify.sanitize>[1]

export function sanitizeHtml(dirty: string, config?: SanitizeConfig): string {
  return DOMPurify.sanitize(dirty, config) as string
}
