<template>
  <div class="post-editor" v-loading="loading">
    <div class="post-editor__header">
      <el-button :icon="ArrowLeft" @click="router.back()">返回</el-button>
      <div class="post-editor__title">{{ isEdit ? '编辑文章' : '新建文章' }}</div>
      <div class="post-editor__actions">
        <el-button @click="onSave(0)">存为草稿</el-button>
        <el-button type="primary" @click="onSave(1)">保存并发布</el-button>
      </div>
    </div>

    <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
      <el-form-item prop="title" label="标题">
        <el-input v-model="form.title" placeholder="请输入文章标题" maxlength="120" show-word-limit />
      </el-form-item>

      <el-row :gutter="16">
        <el-col :xs="24" :sm="12">
          <el-form-item prop="slug" label="Slug（URL 别名）">
            <el-input v-model="form.slug" placeholder="留空则自动生成" />
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item prop="categoryId" label="分类">
            <el-select v-model="form.categoryId" placeholder="选择分类" clearable style="width: 100%">
              <el-option
                v-for="c in categories"
                :key="c.id"
                :label="c.name"
                :value="c.id"
              />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="16">
        <el-col :xs="24" :sm="12">
          <el-form-item label="标签">
            <el-input
              v-model="tagInput"
              placeholder="输入后回车添加标签"
              @keyup.enter="onAddTag"
            />
            <div class="post-editor__tags">
              <el-tag
                v-for="(t, i) in form.tags"
                :key="t"
                closable
                @close="form.tags!.splice(i, 1)"
              >{{ t }}</el-tag>
            </div>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item label="封面图">
            <el-input v-model="form.cover" placeholder="封面图 URL">
              <template #append>
                <el-button @click="coverPickerVisible = true">选择</el-button>
              </template>
            </el-input>
            <el-image
              v-if="form.cover"
              :src="form.cover"
              fit="cover"
              class="post-editor__cover"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item prop="summary" label="摘要">
        <el-input
          v-model="form.summary"
          type="textarea"
          :rows="2"
          maxlength="300"
          show-word-limit
          placeholder="文章摘要（用于列表/SEO）"
        />
      </el-form-item>

      <el-form-item prop="content" label="正文">
        <div class="post-editor__content">
          <Toolbar :editor="editorRef" :default-config="toolbarConfig" style="border-bottom: 1px solid var(--el-border-color-light)" />
          <Editor
            v-model="form.content"
            :default-config="editorConfig"
            style="height: 480px; overflow-y: hidden"
            @onCreated="handleCreated"
          />
        </div>
      </el-form-item>

      <el-form-item label="置顶">
        <el-switch v-model="form.isTop" />
      </el-form-item>
    </el-form>

    <MediaPicker v-model:visible="coverPickerVisible" @select="onPickCover" />
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, shallowRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import type { IDomEditor } from '@wangeditor/editor'
import { getPost, createPost, updatePost } from '@/api/posts'
import { getCategories } from '@/api/categories'
import type { Category, PostForm } from '@/types'
import MediaPicker from '@/components/MediaPicker.vue'

const route = useRoute()
const router = useRouter()
const id = computed(() => Number(route.params.id) || 0)
const isEdit = computed(() => !!id.value)
const loading = ref(false)

const formRef = ref<FormInstance>()
const categories = ref<Category[]>([])
const tagInput = ref('')
const coverPickerVisible = ref(false)

const form = reactive<PostForm>({
  title: '',
  slug: '',
  summary: '',
  content: '',
  cover: '',
  categoryId: undefined,
  tags: [],
  status: 0,
  isTop: false,
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入正文', trigger: 'blur' }],
}

// wangEditor
const editorRef = shallowRef<IDomEditor>()
const toolbarConfig = {}
const editorConfig = {
  placeholder: '开始撰写文章正文…',
  MENU_CONF: {
    uploadImage: {
      // 接入后端上传接口：成功返回 { errno: 0, data: { url } }
      async customUpload(file: File, insertFn: (url: string) => void) {
        const { uploadMedia } = await import('@/api/media')
        const res = await uploadMedia(file)
        insertFn(res.url)
      },
    },
  },
}

function handleCreated(editor: IDomEditor) {
  editorRef.value = editor
}
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor) editor.destroy()
})

function onAddTag() {
  const t = tagInput.value.trim()
  if (!t) return
  if (!form.tags!.includes(t)) form.tags!.push(t)
  tagInput.value = ''
}

function onPickCover(url: string) {
  form.cover = url
  coverPickerVisible.value = false
}

async function loadPost() {
  if (!isEdit.value) return
  loading.value = true
  try {
    const post = await getPost(id.value)
    Object.assign(form, {
      title: post.title,
      slug: post.slug ?? '',
      summary: post.summary ?? '',
      content: post.content,
      cover: post.cover ?? '',
      categoryId: post.categoryId,
      tags: post.tags ?? [],
      status: post.status,
      isTop: post.isTop,
    })
  } finally {
    loading.value = false
  }
}

async function onSave(targetStatus: number) {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  const payload: PostForm = { ...form, status: targetStatus }
  try {
    if (isEdit.value) {
      await updatePost(id.value, payload)
      ElMessage.success(targetStatus === 1 ? '已更新并发布' : '已保存为草稿')
    } else {
      await createPost(payload)
      ElMessage.success(targetStatus === 1 ? '已创建并发布' : '已保存为草稿')
    }
    router.push('/posts')
  } catch {
    /* 拦截器已提示 */
  }
}

onMounted(async () => {
  categories.value = await getCategories()
  loadPost()
})
</script>

<style scoped lang="scss">
.post-editor {
  &__header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
  }
  &__title {
    font-size: 16px;
    font-weight: 600;
    color: var(--el-text-color-primary);
  }
  &__tags {
    margin-top: 8px;
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }
  &__cover {
    width: 160px;
    height: 90px;
    border-radius: 6px;
    margin-top: 8px;
    border: 1px solid var(--el-border-color-light);
  }
  &__content {
    border: 1px solid var(--el-border-color);
    border-radius: 4px;
    overflow: hidden;
    width: 100%;
  }
}
</style>
