<template>
  <div class="dashboard" v-loading="loading">
    <el-row :gutter="16" class="dashboard__stat-row">
      <el-col v-for="card in statCards" :key="card.label" :xs="12" :sm="12" :md="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-card__inner">
            <div class="stat-card__icon" :style="{ background: card.color }">
              <el-icon :size="24">
                <component :is="card.icon" />
              </el-icon>
            </div>
            <div class="stat-card__body">
              <div class="stat-card__value">{{ card.value }}</div>
              <div class="stat-card__label">{{ card.label }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="dashboard__panel-row">
      <el-col :xs="24" :md="12">
        <el-card shadow="never" class="panel">
          <template #header>
            <div class="panel__header">
              <span>最近文章</span>
              <el-button text type="primary" @click="router.push('/posts')">全部</el-button>
            </div>
          </template>
          <el-table :data="recentPosts" size="small" :show-header="false">
            <el-table-column prop="title" label="标题" show-overflow-tooltip />
            <el-table-column prop="createdAt" label="时间" width="170">
              <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card shadow="never" class="panel">
          <template #header>
            <div class="panel__header">
              <span>最近评论</span>
              <el-button text type="primary" @click="router.push('/comments')">全部</el-button>
            </div>
          </template>
          <el-table :data="recentComments" size="small" :show-header="false">
            <el-table-column prop="title" label="内容" show-overflow-tooltip />
            <el-table-column prop="createdAt" label="时间" width="170">
              <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getStats, getRecentPosts, getRecentComments } from '@/api/dashboard'
import type { DashboardStats, RecentItem } from '@/types'
import dayjs from 'dayjs'
import { markRaw } from 'vue'
import { Document, Files, ChatDotRound, View, Folder, Cpu, Picture, Promotion } from '@element-plus/icons-vue'

const router = useRouter()
const loading = ref(false)
const stats = ref<DashboardStats>({
  postCount: 0,
  publishedCount: 0,
  categoryCount: 0,
  commentCount: 0,
  projectCount: 0,
  techStackCount: 0,
  mediaCount: 0,
  viewCount: 0,
})
const recentPosts = ref<RecentItem[]>([])
const recentComments = ref<RecentItem[]>([])

const statCards = ref<{ label: string; value: number; icon: unknown; color: string }[]>([])

function buildCards() {
  statCards.value = [
    { label: '文章总数', value: stats.value.postCount, icon: markRaw(Document), color: '#3b82f6' },
    {
      label: '已发布',
      value: stats.value.publishedCount,
      icon: markRaw(Promotion),
      color: '#10b981',
    },
    { label: '总浏览', value: stats.value.viewCount, icon: markRaw(View), color: '#f59e0b' },
    {
      label: '评论数',
      value: stats.value.commentCount,
      icon: markRaw(ChatDotRound),
      color: '#ef4444',
    },
    {
      label: '分类',
      value: stats.value.categoryCount,
      icon: markRaw(Files),
      color: '#8b5cf6',
    },
    { label: '项目', value: stats.value.projectCount, icon: markRaw(Folder), color: '#06b6d4' },
    {
      label: '技术栈',
      value: stats.value.techStackCount,
      icon: markRaw(Cpu),
      color: '#ec4899',
    },
    { label: '媒体文件', value: stats.value.mediaCount, icon: markRaw(Picture), color: '#14b8a6' },
  ]
}

function formatTime(t: string) {
  return dayjs(t).format('YYYY-MM-DD HH:mm')
}

async function fetchData() {
  loading.value = true
  try {
    const [s, p, c] = await Promise.all([
      getStats(),
      getRecentPosts(6),
      getRecentComments(6),
    ])
    stats.value = s
    recentPosts.value = p
    recentComments.value = c
    buildCards()
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)
</script>

<style scoped lang="scss">
.dashboard {
  &__stat-row {
    margin-bottom: 4px;
  }
  &__panel-row {
    margin-top: 12px;
  }
}

.stat-card {
  margin-bottom: 16px;
  &__inner {
    display: flex;
    align-items: center;
    gap: 16px;
  }
  &__icon {
    width: 48px;
    height: 48px;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    flex-shrink: 0;
  }
  &__value {
    font-size: 24px;
    font-weight: 700;
    color: var(--el-text-color-primary);
    line-height: 1.2;
  }
  &__label {
    font-size: 13px;
    color: var(--el-text-color-secondary);
    margin-top: 2px;
  }
}

.panel {
  &__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>
