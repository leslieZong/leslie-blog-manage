<template>
  <div class="page-container">
    <div class="page-toolbar">
      <div class="page-toolbar__left">
        <el-input
          v-model="query.keyword"
          placeholder="搜索评论内容 / 作者"
          clearable
          style="width: 260px"
          :prefix-icon="Search"
          @keyup.enter="onSearch"
        />
        <el-select v-model="query.status" placeholder="状态" clearable style="width: 140px" @change="onSearch">
          <el-option label="待审核" :value="0" />
          <el-option label="已通过" :value="1" />
          <el-option label="已拒绝" :value="2" />
        </el-select>
        <el-button type="primary" :icon="Search" @click="onSearch">搜索</el-button>
      </div>
    </div>

    <el-table v-loading="loading" :data="list" border row-key="id">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="作者" min-width="180">
        <template #default="{ row }">
          <div class="comment-author">
            <el-avatar :size="28" :src="row.avatar">{{ row.author?.charAt(0) }}</el-avatar>
            <div>
              <div class="comment-author__name">{{ row.author }}</div>
              <div class="comment-author__email">{{ row.email || '-' }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="content" label="内容" min-width="280" show-overflow-tooltip />
      <el-table-column prop="postTitle" label="所属文章" min-width="160" show-overflow-tooltip />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="ip" label="IP" width="130">
        <template #default="{ row }">{{ row.ip || '-' }}</template>
      </el-table-column>
      <el-table-column prop="createdAt" label="时间" width="170">
        <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="240" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status !== 1"
            link
            type="success"
            :icon="CircleCheck"
            @click="onApprove(row)"
          >通过</el-button>
          <el-button
            v-if="row.status !== 2"
            link
            type="warning"
            :icon="CircleClose"
            @click="onReject(row)"
          >拒绝</el-button>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Delete,
  CircleCheck,
  CircleClose,
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import {
  getComments,
  approveComment,
  rejectComment,
  deleteComment,
} from '@/api/comments'
import type { Comment, PageQuery } from '@/types'

interface CommentQuery extends PageQuery {
  status?: number
}

const loading = ref(false)
const list = ref<Comment[]>([])
const total = ref(0)
const query = reactive<CommentQuery>({ page: 1, pageSize: 10, keyword: '' })

function statusText(s: number) {
  return ['待审核', '已通过', '已拒绝'][s] ?? '未知'
}
function statusType(s: number): 'warning' | 'success' | 'info' {
  return (['warning', 'success', 'info'][s] ?? 'info') as 'warning' | 'success' | 'info'
}
function formatTime(t: string) {
  return dayjs(t).format('YYYY-MM-DD HH:mm')
}
async function loadData() {
  loading.value = true
  try {
    const res = await getComments(query)
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
async function onApprove(row: Comment) {
  await approveComment(row.id)
  ElMessage.success('已通过')
  loadData()
}
async function onReject(row: Comment) {
  await ElMessageBox.confirm('确认拒绝该评论？', '操作确认', { type: 'warning' })
  await rejectComment(row.id)
  ElMessage.success('已拒绝')
  loadData()
}
async function onDelete(row: Comment) {
  await ElMessageBox.confirm('确认删除该评论？', '删除确认', { type: 'warning' })
  await deleteComment(row.id)
  ElMessage.success('已删除')
  loadData()
}
onMounted(loadData)
</script>

<style scoped lang="scss">
.comment-author {
  display: flex;
  align-items: center;
  gap: 8px;
  &__name {
    font-size: 13px;
    color: var(--el-text-color-primary);
  }
  &__email {
    font-size: 12px;
    color: var(--el-text-color-secondary);
  }
}
</style>
