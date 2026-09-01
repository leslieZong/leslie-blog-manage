<template>
  <div class="page-container">
    <div class="page-toolbar">
      <div class="page-toolbar__left">
        <el-input
          v-model="keyword"
          placeholder="搜索技术"
          clearable
          style="width: 240px"
          :prefix-icon="Search"
          @keyup.enter="loadData"
        />
        <el-button type="primary" :icon="Search" @click="loadData">搜索</el-button>
      </div>
      <div>
        <el-button type="primary" :icon="Plus" @click="onAdd">新建</el-button>
      </div>
    </div>

    <el-table v-loading="loading" :data="list" border row-key="id">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="图标" width="80">
        <template #default="{ row }">
          <el-image
            v-if="row.icon"
            :src="row.icon"
            fit="contain"
            style="width: 32px; height: 32px"
          />
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="140">
        <template #default="{ row }">
          <el-link type="primary" @click="onEdit(row)">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="category" label="分类" width="140">
        <template #default="{ row }">{{ row.category || '-' }}</template>
      </el-table-column>
      <el-table-column label="熟练度" width="220">
        <template #default="{ row }">
          <el-progress
            :percentage="row.level ?? 0"
            :stroke-width="10"
            :color="levelColor(row.level)"
          />
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="180" show-overflow-tooltip />
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="520px" @closed="resetForm">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item prop="name" label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.icon" placeholder="图标 URL">
            <template #append>
              <el-button @click="pickerVisible = true">选择</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="分类">
          <el-input v-model="form.category" placeholder="如：前端 / 后端 / 工具" />
        </el-form-item>
        <el-form-item label="熟练度">
          <el-slider v-model="form.level" :max="100" show-input />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sortOrder" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="onSave">保存</el-button>
      </template>
    </el-dialog>

    <MediaPicker v-model:visible="pickerVisible" @select="onPickIcon" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Search, Plus, Edit, Delete } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import {
  getTechStack,
  createTechStack,
  updateTechStack,
  deleteTechStack,
} from '@/api/tech-stack'
import type { TechStack, TechStackForm } from '@/types'
import MediaPicker from '@/components/MediaPicker.vue'

const loading = ref(false)
const list = ref<TechStack[]>([])
const keyword = ref('')

const dialogVisible = ref(false)
const saving = ref(false)
const pickerVisible = ref(false)
const formRef = ref<FormInstance>()
const form = reactive<TechStackForm>({
  name: '',
  icon: '',
  category: '',
  level: 60,
  description: '',
  sortOrder: 0,
})
const rules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
}

const isEdit = computed(() => !!form.id)
const dialogTitle = computed(() => (isEdit.value ? '编辑技术' : '新建技术'))

function levelColor(level?: number) {
  if ((level ?? 0) >= 80) return '#10b981'
  if ((level ?? 0) >= 50) return '#3b82f6'
  return '#f59e0b'
}
function formatTime(t: string) {
  return dayjs(t).format('YYYY-MM-DD HH:mm')
}
async function loadData() {
  loading.value = true
  try {
    list.value = await getTechStack({ keyword: keyword.value })
  } finally {
    loading.value = false
  }
}
function resetForm() {
  Object.assign(form, {
    id: undefined,
    name: '',
    icon: '',
    category: '',
    level: 60,
    description: '',
    sortOrder: 0,
  })
  formRef.value?.clearValidate()
}
function onAdd() {
  resetForm()
  dialogVisible.value = true
}
function onEdit(row: TechStack) {
  resetForm()
  Object.assign(form, {
    id: row.id,
    name: row.name,
    icon: row.icon ?? '',
    category: row.category ?? '',
    level: row.level ?? 60,
    description: row.description ?? '',
    sortOrder: row.sortOrder ?? 0,
  })
  dialogVisible.value = true
}
function onPickIcon(url: string) {
  form.icon = url
  pickerVisible.value = false
}
async function onSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (isEdit.value) {
      await updateTechStack(form.id!, { ...form })
      ElMessage.success('已更新')
    } else {
      await createTechStack({ ...form })
      ElMessage.success('已创建')
    }
    dialogVisible.value = false
    loadData()
  } finally {
    saving.value = false
  }
}
async function onDelete(row: TechStack) {
  await ElMessageBox.confirm(`确认删除「${row.name}」？`, '删除确认', { type: 'warning' })
  await deleteTechStack(row.id)
  ElMessage.success('已删除')
  loadData()
}
onMounted(loadData)
</script>
