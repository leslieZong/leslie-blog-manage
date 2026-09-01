<template>
  <el-dialog
    :model-value="visible"
    title="选择媒体"
    width="780px"
    @update:model-value="(v: boolean) => emit('update:visible', v)"
  >
    <div class="media-picker">
      <el-upload
        class="media-picker__upload"
        :show-file-list="false"
        :http-request="onUpload"
        accept="image/*"
      >
        <el-button type="primary" :icon="Upload" :loading="uploading">上传图片</el-button>
      </el-upload>

      <div v-loading="loading" class="media-picker__grid">
        <div
          v-for="m in list"
          :key="m.id"
          class="media-picker__item"
          @click="onSelect(m.url)"
        >
          <el-image :src="m.url" fit="cover" class="media-picker__img" />
          <div class="media-picker__name" :title="m.name">{{ m.name }}</div>
        </div>
        <el-empty v-if="!loading && !list.length" description="暂无媒体" />
      </div>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.pageSize"
          :total="total"
          layout="prev, pager, next"
          background
          @current-change="loadData"
        />
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import { getMediaList, uploadMedia } from '@/api/media'
import type { MediaItem, PageQuery } from '@/types'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{
  'update:visible': [v: boolean]
  select: [url: string]
}>()

const loading = ref(false)
const uploading = ref(false)
const list = ref<MediaItem[]>([])
const total = ref(0)
const query = reactive<PageQuery & { type: string }>({
  page: 1,
  pageSize: 18,
  type: 'image',
})

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

async function onUpload(opt: { file: File }) {
  uploading.value = true
  try {
    await uploadMedia(opt.file)
    ElMessage.success('上传成功')
    query.page = 1
    loadData()
  } finally {
    uploading.value = false
  }
}

function onSelect(url: string) {
  emit('select', url)
  emit('update:visible', false)
}

watch(
  () => props.visible,
  (v) => {
    if (v) loadData()
  },
)
</script>

<style scoped lang="scss">
.media-picker {
  &__upload {
    margin-bottom: 12px;
  }
  &__grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 12px;
    max-height: 420px;
    overflow-y: auto;
  }
  &__item {
    cursor: pointer;
    border: 1px solid var(--el-border-color-light);
    border-radius: 6px;
    overflow: hidden;
    transition: border-color 0.2s;
    &:hover {
      border-color: var(--el-color-primary);
    }
  }
  &__img {
    width: 100%;
    height: 90px;
    display: block;
  }
  &__name {
    padding: 4px 6px;
    font-size: 12px;
    color: var(--el-text-color-secondary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}
</style>
