<template>
  <div class="page-container">
    <div class="page-toolbar">
      <div class="page-toolbar__left">
        <el-input
          v-model="query.keyword"
          placeholder="搜索项目"
          clearable
          style="width: 240px"
          :prefix-icon="Search"
          @keyup.enter="onSearch"
        />
        <el-select v-model="query.status" placeholder="状态" clearable style="width: 140px" @change="onSearch">
          <el-option label="进行中" :value="0" />
          <el-option label="已完成" :value="1" />
          <el-option label="已归档" :value="2" />
        </el-select>
        <el-button type="primary" :icon="Search" @click="onSearch">搜索</el-button>
      </div>
      <div>
        <el-button type="primary" :icon="Plus" @click="onAdd">新建项目</el-button>
      </div>
    </div>

    <el-table v-loading="loading" :data="list" border row-key="id">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="封面" width="90">
        <template #default="{ row }">
          <el-image
            v-if="row.cover"
            :src="row.cover"
            fit="cover"
            style="width: 50px; height: 50px; border-radius: 4px"
          />
          <el-icon v-else><Picture /></el-icon>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="160">
        <template #default="{ row }">
          <el-link type="primary" @click="onEdit(row)">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column label="技术栈" min-width="180">
        <template #default="{ row }">
          <el-tag v-for="t in row.techStack" :key="t" size="small" style="margin: 2px">
            {{ t }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="链接" width="160">
        <template #default="{ row }">
          <el-link v-if="row.demoUrl" :href="row.demoUrl" target="_blank" type="primary">
            Demo
          </el-link>
          <el-link v-if="row.repoUrl" :href="row.repoUrl" target="_blank" type="info" style="margin-left: 8px">
            Repo
          </el-link>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="sortOrder" label="排序" width="80" />
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" :icon="Edit" @click="onEdit(row)">编辑</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" @closed="resetForm">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item prop="name" label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="封面">
          <el-input v-model="form.cover" placeholder="封面图 URL">
            <template #append>
              <el-button @click="pickerVisible = true">选择</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="技术栈">
          <el-input v-model="tagInput" placeholder="回车添加" @keyup.enter="onAddTech" />
          <div style="margin-top: 8px">
            <el-tag
              v-for="(t, i) in form.techStack"
              :key="t"
              closable
              style="margin: 2px"
              @close="form.techStack!.splice(i, 1)"
            >{{ t }}</el-tag>
          </div>
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="Demo URL">
              <el-input v-model="form.demoUrl" placeholder="https://" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Repo URL">
              <el-input v-model="form.repoUrl" placeholder="https://" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="状态">
              <el-select v-model="form.status" style="width: 100%">
                <el-option label="进行中" :value="0" />
                <el-option label="已完成" :value="1" />
                <el-option label="已归档" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序">
              <el-input-number v-model="form.sortOrder" :min="0" :max="9999" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="onSave">保存</el-button>
      </template>
    </el-dialog>

    <MediaPicker v-model:visible="pickerVisible" @select="onPickCover" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Search, Plus, Edit, Delete, Picture } from '@element-plus/icons-vue'
import {
  getProjects,
  createProject,
  updateProject,
  deleteProject,
} from '@/api/projects'
import type { Project, ProjectForm, PageQuery } from '@/types'
import MediaPicker from '@/components/MediaPicker.vue'

interface ProjectQuery extends PageQuery {
  status?: number
}

const loading = ref(false)
const list = ref<Project[]>([])
const total = ref(0)
const query = reactive<ProjectQuery>({ page: 1, pageSize: 10, keyword: '' })

const dialogVisible = ref(false)
const saving = ref(false)
const pickerVisible = ref(false)
const tagInput = ref('')
const formRef = ref<FormInstance>()
const form = reactive<ProjectForm>({
  name: '',
  description: '',
  cover: '',
  demoUrl: '',
  repoUrl: '',
  techStack: [],
  status: 0,
  sortOrder: 0,
})
const rules: FormRules = {
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
}

const isEdit = computed(() => !!form.id)
const dialogTitle = computed(() => (isEdit.value ? '编辑项目' : '新建项目'))

function statusText(s: number) {
  return ['进行中', '已完成', '已归档'][s] ?? '未知'
}
function statusType(s: number): 'success' | 'warning' | 'info' {
  return (['warning', 'success', 'info'][s] ?? 'info') as 'success' | 'warning' | 'info'
}

async function loadData() {
  loading.value = true
  try {
    const res = await getProjects(query)
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
function resetForm() {
  Object.assign(form, {
    id: undefined,
    name: '',
    description: '',
    cover: '',
    demoUrl: '',
    repoUrl: '',
    techStack: [],
    status: 0,
    sortOrder: 0,
  })
  tagInput.value = ''
  formRef.value?.clearValidate()
}
function onAdd() {
  resetForm()
  dialogVisible.value = true
}
function onEdit(row: Project) {
  resetForm()
  Object.assign(form, {
    id: row.id,
    name: row.name,
    description: row.description ?? '',
    cover: row.cover ?? '',
    demoUrl: row.demoUrl ?? '',
    repoUrl: row.repoUrl ?? '',
    techStack: row.techStack ? [...row.techStack] : [],
    status: row.status,
    sortOrder: row.sortOrder ?? 0,
  })
  dialogVisible.value = true
}
function onAddTech() {
  const t = tagInput.value.trim()
  if (!t) return
  if (!form.techStack!.includes(t)) form.techStack!.push(t)
  tagInput.value = ''
}
function onPickCover(url: string) {
  form.cover = url
  pickerVisible.value = false
}
async function onSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (isEdit.value) {
      await updateProject(form.id!, { ...form })
      ElMessage.success('已更新')
    } else {
      await createProject({ ...form })
      ElMessage.success('已创建')
    }
    dialogVisible.value = false
    loadData()
  } finally {
    saving.value = false
  }
}
async function onDelete(row: Project) {
  await ElMessageBox.confirm(`确认删除项目「${row.name}」？`, '删除确认', { type: 'warning' })
  await deleteProject(row.id)
  ElMessage.success('已删除')
  loadData()
}

onMounted(loadData)
</script>
