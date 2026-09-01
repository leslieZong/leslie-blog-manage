<template>
  <div class="page-container">
    <div class="page-toolbar">
      <div class="page-toolbar__left">
        <el-input
          v-model="query.keyword"
          placeholder="搜索文件名"
          clearable
          style="width: 240px"
          :prefix-icon="Search"
          @keyup.enter="onSearch"
        />
        <el-select v-model="query.type" placeholder="类型" clearable style="width: 140px" @change="onSearch">
          <el-option label="图片" value="image" />
          <el-option label="视频" value="video" />
          <el-option label="文档" value="file" />
        </el-select>
        <el-button type="primary" :icon="Search" @click="onSearch">搜索</el-button>
      </div>
      <div>
        <el-upload
          :show-file-list="false"
          multiple
          :http-request="onUpload"
          accept="image/*,video/*,.pdf,.zip,.doc,.docx"
        >
          <el-button type="primary" :icon="Upload" :loading="uploading">上传文件</el-button>
        </el-upload>
      </div>
    </div>

    <div v-loading="loading" class="media-grid">
      <div v-for="m in list" :key="m.id" class="media-card">
        <div class="media-card__preview">
          <el-image
            v-if="m.type === 'image'"
            :src="m.url"
            fit="cover"
            :preview-src-list="[m.url]"
            preview-teleported
            class="media-card__img"
          />
          <div v-else class="media-card__file">
            <el-icon :size="36"><Document /></el-icon>
            <span>{{ m.mimeType || m.type }}</span>
          </div>
        </div>
        <div class="media-card__name" :title="m.name">{{ m.name }}</div>
        <div class="media-card__meta">{{ formatSize(m.size) }} · {{ formatTime(m.createdAt) }}</div>
        <div class="media-card__actions">
          <el-button link :icon="CopyDocument" @click="onCopy(m.url)">复制链接</el-button>
          <el-button link type="danger" :icon="Delete" @click="onDelete(m)">删除</el-button>
        </div>
      </div>
      <el-empty v-if="!loading && !list.length" description="暂无媒体文件" />
    </div>

    <div class="pagination-wrap">
      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.pageSize"
        :total="total"
        :page-sizes="[18, 36, 72]"
        layout="total, sizes, prev, pager, next"
        background
        @size-change="loadData"
        @current-change="loadData"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Upload,
  Delete,
  Document,
  CopyDocument,
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { getMediaList, uploadMedia, deleteMedia } from '@/api/media'
import type { MediaItem, PageQuery } from '@/types'

interface MediaQuery extends PageQuery {
  type?: string
}

const loading = ref(false)
const uploading = ref(false)
const list = ref<MediaItem[]>([])
const total = ref(0)
const query = reactive<MediaQuery>({ page: 1, pageSize: 18, keyword: '' })

function formatSize(bytes: number) {
  if (!bytes) return '-'
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / 1024 / 1024).toFixed(1)} MB`
}
function formatTime(t: string) {
  return dayjs(t).format('YYYY-MM-DD')
}
async function loadData() {
  loading.value = true
  try {
    const res = await getMediaList(query)
    list.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}
function onSearch() {
  query.page = 1
  loadData()
}
async function onUpload(opt: { file: File }) {
  uploading.value = true
  try {
    await uploadMedia(opt.file)
    ElMessage.success('上传成功')
    loadData()
  } finally {
    uploading.value = false
  }
}
async function onCopy(url: string) {
  await navigator.clipboard.writeText(url)
  ElMessage.success('链接已复制')
}
async function onDelete(m: MediaItem) {
  await ElMessageBox.confirm(`确认删除「${m.name}」？`, '删除确认', { type: 'warning' })
  await deleteMedia(m.id)
  ElMessage.success('已删除')
  loadData()
}
onMounted(loadData)
</script>

<style scoped lang="scss">
.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 16px;
  min-height: 200px;
}
.media-card {
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
  overflow: hidden;
  background: var(--el-bg-color);
  &__preview {
    width: 100%;
    height: 130px;
    background: var(--el-fill-color);
    display: flex;
    align-items: center;
    justify-content: center;
  }
  &__img {
    width: 100%;
    height: 100%;
    display: block;
  }
  &__file {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    color: var(--el-text-color-secondary);
    font-size: 12px;
  }
  &__name {
    padding: 6px 8px 0;
    font-size: 13px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  &__meta {
    padding: 0 8px 4px;
    font-size: 12px;
    color: var(--el-text-color-secondary);
  }
  &__actions {
    display: flex;
    justify-content: space-between;
    border-top: 1px solid var(--el-border-color-light);
  }
}
</style>
