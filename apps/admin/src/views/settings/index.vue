<template>
  <div class="page-container settings" v-loading="loading">
    <el-form :model="form" label-width="100px" class="settings__form">
      <el-card shadow="never" class="settings__card">
        <template #header>站点信息</template>
        <el-form-item label="站点名称" prop="siteName">
          <el-input v-model="form.siteName" placeholder="Leslie Blog" />
        </el-form-item>
        <el-form-item label="Logo" prop="logo">
          <el-input v-model="form.logo" placeholder="Logo URL">
            <template #append>
              <el-button @click="pickerVisible = true">选择</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="站点描述">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="form.keywords" placeholder="多个用英文逗号分隔" />
        </el-form-item>
        <el-form-item label="作者">
          <el-input v-model="form.author" />
        </el-form-item>
        <el-form-item label="备案号">
          <el-input v-model="form.icp" placeholder="如：粤ICP备xxxx号" />
        </el-form-item>
      </el-card>

      <el-card shadow="never" class="settings__card">
        <template #header>社交链接</template>
        <el-form-item label="GitHub">
          <el-input v-model="form.social!.github" placeholder="https://github.com/xxx" />
        </el-form-item>
        <el-form-item label="Email">
          <el-input v-model="form.social!.email" placeholder="mailto:xxx" />
        </el-form-item>
        <el-form-item label="Twitter">
          <el-input v-model="form.social!.twitter" placeholder="https://twitter.com/xxx" />
        </el-form-item>
      </el-card>

      <div class="settings__footer">
        <el-button type="primary" :loading="saving" @click="onSave">保存设置</el-button>
      </div>
    </el-form>

    <MediaPicker v-model:visible="pickerVisible" @select="onPickLogo" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getSettings, updateSettings } from '@/api/settings'
import type { SiteSettings } from '@/types'
import MediaPicker from '@/components/MediaPicker.vue'

const loading = ref(false)
const saving = ref(false)
const pickerVisible = ref(false)

const form = reactive<SiteSettings>({
  siteName: '',
  logo: '',
  description: '',
  keywords: '',
  author: '',
  icp: '',
  social: { github: '', email: '', twitter: '' },
})

async function loadData() {
  loading.value = true
  try {
    const data = await getSettings()
    Object.assign(form, data)
    form.social = { github: '', email: '', twitter: '', ...(data.social ?? {}) }
  } finally {
    loading.value = false
  }
}
function onPickLogo(url: string) {
  form.logo = url
  pickerVisible.value = false
}
async function onSave() {
  saving.value = true
  try {
    await updateSettings({ ...form })
    ElMessage.success('设置已保存')
  } finally {
    saving.value = false
  }
}
onMounted(loadData)
</script>

<style scoped lang="scss">
.settings {
  &__form {
    max-width: 720px;
  margin: 0 auto;
  }
  &__card {
    margin-bottom: 16px;
  }
  &__footer {
    text-align: center;
    padding: 8px 0;
  }
}
</style>
