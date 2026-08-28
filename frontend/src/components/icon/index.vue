<script setup lang="ts">
/**
 * BrandIcon - 品牌/通用图标组件
 * 支持 GitHub / RSS / Email / 向下箭头系列(arrowDown, arrowDownLine, arrowDownScroll)
 *
 * 用法:
 *   <BrandIcon name="github" :size="24" color="#181717" />
 *   <BrandIcon name="arrowDown" :size="32" />              <!-- 粗圆角下箭头 -->
 *   <BrandIcon name="arrowDownLine" :size="24" />          <!-- 细线条下箭头 -->
 *   <BrandIcon name="arrowDownScroll" :size="28" />        <!-- 圆点+箭头(滚动指示器) -->
 *   <BrandIcon name="email" :size="20" label="邮箱" />     <!-- label 用于无障碍 -->
 *
 * 说明:
 *   - 所有图标基于 24x24 viewBox,任意缩放不失真
 *   - color 默认 currentColor,可随文字颜色自动变化
 *   - fill 型图标:配置 fill 字段;线型图标:配置 paths + strokeWidth
 *   - 也可直接引入 /icons 下的 .svg 文件(支持按需加载)
 */
import { computed } from 'vue'

const props = defineProps({
  // 图标名: github | rss | email | arrowDown | arrowDownLine | arrowDownScroll
  name: { type: String, required: true },
  // 尺寸(px),默认 24
  size: { type: [Number, String], default: 24 },
  // 颜色,默认继承文字颜色
  color: { type: String, default: 'currentColor' },
  // 线型图标的描边宽度(仅对 stroke 图标生效)
  strokeWidth: { type: [Number, String], default: 2 },
  // 无障碍标签(可选)
  label: { type: String, default: '' },
})

// 图标配置表,新增图标只需在此追加
//  - fill: 纯填充图标的 path d
//  - paths: 线型图标的 path 数组(描边、圆头圆角)
//  - strokeWidth: 可选,自定义描边宽度
const icons = {
  github: {
    fill: 'M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12',
  },
  rss: {
    fill: 'M6.18 15.64a2.18 2.18 0 0 1 2.18 2.18C8.36 19 7.38 20 6.18 20 5 20 4 19 4 17.82a2.18 2.18 0 0 1 2.18-2.18M4 4.44A15.56 15.56 0 0 1 19.56 20h-2.83A12.73 12.73 0 0 0 4 7.27V4.44m0 5.66a9.9 9.9 0 0 1 9.9 9.9h-2.83A7.07 7.07 0 0 0 4 12.93V10.1Z',
  },
  email: {
    fill: 'M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z',
  },
  // 粗圆角下箭头(饱满、圆润,默认推荐)
  arrowDown: {
    paths: ['M12 4v10', 'M6 11l6 6 6-6'],
    strokeWidth: 3.4,
  },
  // 细线条圆角下箭头(简洁)
  arrowDownLine: {
    paths: ['M12 4v10', 'M6 11l6 6 6-6'],
  },
  // 顶部圆点 + 箭头(向下滚动指示器,最有造型)
  arrowDownScroll: {
    paths: [
      'M12 3.6a1.8 1.8 0 1 1 0 3.6a1.8 1.8 0 1 1 0-3.6z',
      'M12 9v7',
      'M6.5 12.5l5.5 5.5 5.5-5.5',
    ],
  },
}

const current = computed(() => icons[props.name as keyof typeof icons] || {})
const isFill = computed(() => 'fill' in current.value)
const isStroke = computed(() => 'strokeWidth' in current.value)
const effectiveStrokeWidth = computed(() =>
  isStroke.value ? (current.value as { strokeWidth: number }).strokeWidth : props.strokeWidth,
)
</script>

<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    :fill="isFill ? color : 'none'"
    :stroke="isFill ? 'none' : color"
    :stroke-width="isFill ? undefined : effectiveStrokeWidth"
    stroke-linecap="round"
    stroke-linejoin="round"
    :role="label ? 'img' : 'presentation'"
    :aria-label="label || undefined"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path v-if="isFill" :d="(current as { fill: string }).fill" />
    <path v-for="(p, i) in (current as { paths: string[] })?.paths" v-else :key="i" :d="p" />
  </svg>
</template>
