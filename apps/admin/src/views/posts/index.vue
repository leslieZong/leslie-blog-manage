<template>
  <div class="page-container">
    <div class="page-toolbar">
      <div class="page-toolbar__left">
        <el-input
          v-model="query.keyword"
          placeholder="搜索标题 / 标签"
          clearable
          style="width: 240px"
          :prefix-icon="Search"
          @clear="loadData"
          @keyup.enter="onSearch"
        />
        <el-select
          v-model="query.categoryId"
          placeholder="分类"
          clearable
          style="width: 160px"
          @change="onSearch"
        >
          <el-option
            v-for="c in categories"
            :key="c.id"
            :label="c.name"
            :value="c.id"
          />
        </el-select>
        <el-select
          v-model="query.status"
          placeholder="状态"
          clearable
          style="width: 120px"
          @change="onSearch"
        >
          <el-option label="草稿" :value="0" />
          <el-option label="已发布" :value="1" />
          <el-option label="已下线" :value="2" />
        </el-select>
        <el-button type="primary" :icon="Search" @click="onSearch">搜索</el-button>
        <el-button @click="onReset">重置</el-button>
      </div>
      <div>
        <el-button type="primary" :icon="Plus" @click="router.push('/posts/editor')">
          新建文章
        </el-button>
      </div>
    </div>

    <el-table v-loading="loading" :data="list" border row-key="id">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="title" label="标题" min-width="220" show-overflow-tooltip>
        <template #default="{ row }">
          <el-link type="primary" @click="router.push(`/posts/editor/${row.id}`)">
            {{ row.title }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column prop="categoryName" label="分类" width="120">
        <template #default="{ row }">{{ row.categoryName || '-' }}</template>
      </el-table-column>
      <el-table-column label="标签" min-width="160">
        <template #default="{ row }">
          <el-tag
            v-for="t in row.tags"
            :key="t"
            size="small"
            style="margin: 2px"
          >{{ t }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)" effect="light">
            {{ statusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="viewCount" label="浏览" width="90" />
      <el-table-column prop="commentCount" label="评论" width="90" />
      <el-table-column label="置顶" width="80">
        <template #default="{ row }">
          <el-switch
            :model-value="!!row.isTop"
            @change="(v: string | number | boolean) => onToggleTop(row, Boolean(v))"
          />
        </template>
      </el-table-column>
      <el-table-column prop="updatedAt" label="更新时间" width="170">
        <template #default="{ row }">{{ formatTime(row.updatedAt) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="240" fixed="right">
        <template #default="{ row }">
          <el-button
            link
            type="primary"
            :icon="Edit"
            @click="router.push(`/posts/editor/${row.id}`)"
          >编辑</el-button>
          <el-button
            v-if="row.status !== 1"
            link
            type="success"
            :icon="Promotion"
            @click="onPublish(row)"
          >发布</el-button>
          <el-button
            v-else
            link
            type="warning"
            :icon="SemiSelect"
            @click="onUnpublish(row)"
          >下线</el-button>
          <el-button link type="danger" :icon="Delete" @click="onDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrap">
      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        background
        @size-change="loadData"
        @current-change="loadData"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Plus,
  Edit,
  Delete,
  Promotion,
  SemiSelect,
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import {
  getPosts,
  deletePost,
  publishPost,
  unpublishPost,
  toggleTop,
  type PostQuery,
} from '@/api/posts'
import { getCategories } from '@/api/categories'
import type { Post, Category } from '@/types'

const router = useRouter()
const loading = ref(false)
const list = ref<Post[]>([])
const total = ref(0)
const categories = ref<Category[]>([])

const query = reactive<PostQuery>({
  page: 1,
  pageSize: 10,
  keyword: '',
  status: undefined,
  categoryId: undefined,
})

function statusText(s: number) {
  return ['草稿', '已发布', '已下线'][s] ?? '未知'
}
function statusType(s: number): 'info' | 'success' | 'warning' {
  return (['info', 'success', 'warning'][s] ?? 'info') as 'info' | 'success' | 'warning'
}
function formatTime(t: string) {
  return dayjs(t).format('YYYY-MM-DD HH:mm')
}

async function loadData() {
  loading.value = true
  try {
    const res = await getPosts(query)
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
function onReset() {
  query.keyword = ''
  query.status = undefined
  query.categoryId = undefined
  onSearch()
}

async function onPublish(row: Post) {
  await ElMessageBox.confirm(`确认发布《${row.title}》？`, '发布确认')
  await publishPost(row.id)
  ElMessage.success('已发布')
  loadData()
}

async function onUnpublish(row: Post) {
  await ElMessageBox.confirm(`确认下线《${row.title}》？`, '下线确认')
  await unpublishPost(row.id)
  ElMessage.success('已下线')
  loadData()
}

async function onToggleTop(row: Post, isTop: boolean) {
  try {
    await toggleTop(row.id, isTop)
    row.isTop = isTop
    ElMessage.success(isTop ? '已置顶' : '已取消置顶')
  } catch {
    /* 失败已由拦截器提示 */
  }
}

async function onDelete(row: Post) {
  await ElMessageBox.confirm(`确认删除《${row.title}》？此操作不可恢复`, '删除确认', {
    type: 'warning',
  })
  await deletePost(row.id)
  ElMessage.success('已删除')
  if (list.value.length === 1 && query.page! > 1) query.page = query.page! - 1
  loadData()
}

onMounted(async () => {
  categories.value = await getCategories()
  loadData()
})
</script>
