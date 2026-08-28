<script setup lang="ts">
/**
 * BrandIcon - 品牌图标组件(GitHub / RSS / Email)
 *
 * 用法:
 *   <BrandIcon name="github" :size="24" color="#181717" />
 *   <BrandIcon name="rss"   :size="32" />          <!-- 颜色默认继承 currentColor -->
 *   <BrandIcon name="email" :size="20" label="邮箱" />  <!-- label 用于无障碍 -->
 *
 * 说明:
 *   - 图标基于 24x24 viewBox 的矢量 path,任意缩放不失真
 *   - color 默认 currentColor,可随文字颜色自动变化
 *   - 也可直接引入 /icons 下的 .svg 文件(支持按需加载)
 */
import { computed } from 'vue'

const props = defineProps({
  // 图标名: github | rss | email
  name: { type: String, required: true },
  // 尺寸(px),默认 24
  size: { type: [Number, String], default: 24 },
  // 无障碍标签(可选)
  label: { type: String, default: '' },
})

// 图标 path 配置表,新增图标只需在此追加
const icons = {
  github:
    'M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12',
  rss: 'M6.18 15.64a2.18 2.18 0 0 1 2.18 2.18C8.36 19 7.38 20 6.18 20 5 20 4 19 4 17.82a2.18 2.18 0 0 1 2.18-2.18M4 4.44A15.56 15.56 0 0 1 19.56 20h-2.83A12.73 12.73 0 0 0 4 7.27V4.44m0 5.66a9.9 9.9 0 0 1 9.9 9.9h-2.83A7.07 7.07 0 0 0 4 12.93V10.1Z',
  email:
    'M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z',
}

const path = computed(() => icons[props.name as keyof typeof icons] || '')
</script>

<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="currentColor"
    :role="label ? 'img' : 'presentation'"
    :aria-label="label || undefined"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path :d="path" />
  </svg>
</template>
