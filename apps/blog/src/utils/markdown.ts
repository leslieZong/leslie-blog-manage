import MarkdownIt, { type Token } from 'markdown-it'
import Prism from 'prismjs'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-python'
import 'prismjs/components/prism-css'
import 'prismjs/components/prism-yaml'
import { sanitizeHtml } from './sanitize'

const md = new MarkdownIt({
  // 禁止 markdown 内嵌裸 HTML（深度防御：输出仍会再过一遍 DOMPurify）
  html: false,
  linkify: true,
  breaks: false,
  highlight(code, lang) {
    const grammar = lang && Prism.languages[lang]
    if (grammar) {
      try {
        // 仅返回高亮后的 HTML 片段，markdown-it 会包裹 <pre><code>
        return Prism.highlight(code, grammar, lang)
      } catch {
        // 降级到默认转义
      }
    }
    return ''
  },
})

// 给所有 <a> 追加 target=_blank rel=noopener noreferrer，防反向标签页钓鱼
const defaultLinkOpen =
  md.renderer.rules.link_open ||
  ((tokens, idx, options, _env, self) => self.renderToken(tokens, idx, options))

function setAttr(token: Token, name: string, value: string) {
  const i = token.attrIndex(name)
  if (i < 0) token.attrPush([name, value])
  else if (token.attrs && token.attrs[i]) token.attrs[i][1] = value
}

md.renderer.rules.link_open = (tokens, idx, options, env, self) => {
  const token = tokens[idx]
  if (token) {
    setAttr(token, 'target', '_blank')
    setAttr(token, 'rel', 'noopener noreferrer')
  }
  return defaultLinkOpen(tokens, idx, options, env, self)
}

// 渲染 markdown 为已净化的 HTML 字符串，供 v-html 使用
export function renderMarkdown(src: string): string {
  const html = md.render(src ?? '')
  // ADD_ATTR: 允许我们注入的 target/rel 通过 DOMPurify
  return sanitizeHtml(html, {
    ADD_ATTR: ['target', 'rel'],
  })
}

export default md
