<template>
  <div class="page-container">
    <div class="page-toolbar">
      <div class="page-toolbar__left">
        <el-input
          v-model="keyword"
          placeholder="搜索分类"
          clearable
          style="width: 240px"
          :prefix-icon="Search"
          @keyup.enter="loadData"
        />
        <el-button type="primary" :icon="Search" @click="loadData">搜索</el-button>
      </div>
      <div>
        <el-button type="primary" :icon="Plus" @click="onAdd">新建分类</el-button>
      </div>
    </div>

    <el-table v-loading="loading" :data="list" border row-key="id" default-expand-all>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" label="名称" min-width="160">
        <template #default="{ row }">
          <el-link type="primary" @click="onEdit(row)">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="slug" label="Slug" width="160">
        <template #default="{ row }">{{ row.slug || '-' }}</template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column prop="postCount" label="文章数" width="100" />
      <el-table-column prop="sortOrder" label="排序" width="90" />
      <el-table-column prop="createdAt" label="创建时间" width="170">
        <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" :icon="Edit" @click="onEdit(row)">编辑</el-button>
          <el-button link type="danger" :icon="Delete" @click="onDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      @closed="resetForm"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item prop="name" label="名称">
          <el-input v-model="form.name" placeholder="分类名称" />
        </el-form-item>
        <el-form-item prop="slug" label="Slug">
          <el-input v-model="form.slug" placeholder="URL 别名，留空自动生成" />
        </el-form-item>
        <el-form-item label="父分类">
          <el-select v-model="form.parentId" placeholder="顶级分类" clearable style="width: 100%">
            <el-option
              v-for="c in parentOptions"
              :key="c.id"
              :label="c.name"
              :value="c.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sortOrder" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="onSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Search, Plus, Edit, Delete } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import {
  getCategories,
  createCategory,
  updateCategory,
  deleteCategory,
} from '@/api/categories'
import type { Category, CategoryForm } from '@/types'

const loading = ref(false)
const list = ref<Category[]>([])
const keyword = ref('')

const dialogVisible = ref(false)
const saving = ref(false)
const formRef = ref<FormInstance>()
const form = reactive<CategoryForm>({
  name: '',
  slug: '',
  description: '',
  parentId: undefined,
  sortOrder: 0,
})
const rules: FormRules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }],
}

const isEdit = computed(() => !!form.id)
const dialogTitle = computed(() => (isEdit.value ? '编辑分类' : '新建分类'))
// 父分类下拉只列顶级分类（避免环路）
const parentOptions = computed(() => list.value.filter((c) => !c.parentId))

function formatTime(t: string) {
  return dayjs(t).format('YYYY-MM-DD HH:mm')
}

async function loadData() {
  loading.value = true
  try {
    list.value = await getCategories({ keyword: keyword.value })
  } finally {
    loading.value = false
  }
}

function resetForm() {
  Object.assign(form, {
    id: undefined,
    name: '',
    slug: '',
    description: '',
    parentId: undefined,
    sortOrder: 0,
  })
  formRef.value?.clearValidate()
}

function onAdd() {
  resetForm()
  dialogVisible.value = true
}

function onEdit(row: Category) {
  resetForm()
  Object.assign(form, {
    id: row.id,
    name: row.name,
    slug: row.slug ?? '',
    description: row.description ?? '',
    parentId: row.parentId,
    sortOrder: row.sortOrder ?? 0,
  })
  dialogVisible.value = true
}

async function onSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (isEdit.value) {
      await updateCategory(form.id!, { ...form })
      ElMessage.success('已更新')
    } else {
      await createCategory({ ...form })
      ElMessage.success('已创建')
    }
    dialogVisible.value = false
    loadData()
  } finally {
    saving.value = false
  }
}

async function onDelete(row: Category) {
  await ElMessageBox.confirm(`确认删除分类「${row.name}」？`, '删除确认', { type: 'warning' })
  await deleteCategory(row.id)
  ElMessage.success('已删除')
  loadData()
}

onMounted(loadData)
</script>
