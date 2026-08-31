<!-- GithubActivity.vue -->
/** * * Github Activity / Github 活动 */
<template>
  <!-- GitHub 活跃度 -->
  <div class="github-activity">
    <div class="github-activity-title">
      <div class="left">
        <h2 class="title">{{ $t('home.githubActivity.title') }}</h2>
        <div>{{ $t('home.githubActivity.description') }}</div>
      </div>
      <div class="right" @click="handleClickActivity">
        <div class="explore-btn">{{ $t('home.githubActivity.viewAll') }}</div>
      </div>
    </div>
    <div class="github-activity-info">
      {{ total }} contributions · 12 repositories · 8 stars
    </div>
    <div class="github-activity-detail">
      <el-date-picker v-model="value" type="year" placeholder="Pick one or more years" />
      <CalendarHeatmap
        :data="data"
        range="year"
        :levels="5"
        :colors="colors"
        :cell-size="11"
        :cell-gap="4"
        :cell-radius="3"
        :locale="locale"
        :theme="themeStore.isDark ? 'dark' : 'light'"
        :today="{ color: 'var(--el-color-primary)', style: 'ring', size: 2 }"
        :tooltip-enabled="true"
        :tooltip-formatter="tooltipFormatter"
        :show-legend="true"
        :show-months="true"
        :show-weekdays="true"
      />
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed } from 'vue'
import { CalendarHeatmap } from 'vue-calendar-activity'
import type { HeatmapLocale, HeatmapDay } from 'vue-calendar-activity'
import 'vue-calendar-activity/style.css'
import dayjs from 'dayjs'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()

const handleClickActivity = () => {
  window.open('https://github.com/leslieZong', '_blank')
}

// 引用 CSS 变量（随 html.dark 自动切换明暗），levels=5 对应 index 0~4
const colors = [
  'var(--github-0)',
  'var(--github-1)',
  'var(--github-2)',
  'var(--github-3)',
  'var(--github-4)',
]

// 单字母星期标签（Mon→Sun），与 11px 格子宽度匹配
const locale = {
  weekdaysShort: ['S', 'M', 'T', 'W', 'T', 'F', 'S'],
} as HeatmapLocale

const tooltipFormatter = (day: HeatmapDay) => {
  const date = dayjs(day.date).format('MMM D, YYYY')
  return day.value === 0 ? `${date} No contributions` : `${date} ${day.value} contributions`
}

// 演示数据：覆盖滚动的一年（range="year" 渲染过去 365 天），稳定伪随机分布
const buildDemoData = () => {
  const items: { date: string; value: number }[] = []
  let seed = 20260831
  const rand = () => {
    seed = (seed * 9301 + 49297) % 233280
    return seed / 233280
  }
  const today = dayjs()
  for (let i = 364; i >= 0; i--) {
    const r = rand()
    if (r < 0.45) continue
    items.push({
      date: today.subtract(i, 'day').format('YYYY-MM-DD'),
      value: Math.min(4, Math.floor(rand() * 4) + 1),
    })
  }
  return items
}

const currentYear = new Date().getFullYear().toString()
const value = ref(currentYear)
const data = ref(buildDemoData())
const total = computed(() => data.value.reduce((sum, d) => sum + d.value, 0))
</script>
<style scoped lang="scss">
.github-activity {
  padding-top: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--el-border-color);
  .github-activity-title {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .left {
      .title {
        margin-bottom: 12px;
      }
    }

    .right {
      display: flex;
      align-items: center;
      gap: 6px;
      cursor: pointer;
      &:hover {
        color: var(--el-color-primary);
      }
    }
  }
  .github-activity-info {
    margin-top: 30px;
    margin-bottom: 30px;
    color: var(--el-text-color-regular);
  }
  .github-activity-detail {
    display: flex;
    flex-direction: column;
    gap: 20px;
    :deep(.calendar-heatmap) {
      --ch-font-color: var(--el-text-color-secondary);
      --ch-bg: var(--el-bg-color-overlay);
      --ch-tooltip-bg: var(--el-bg-color-overlay);
      --ch-tooltip-color: var(--el-text-color-primary);
    }
  }
}
</style>
