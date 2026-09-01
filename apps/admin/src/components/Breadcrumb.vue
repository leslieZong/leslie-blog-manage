<template>
  <el-breadcrumb class="breadcrumb" :separator-icon="ArrowRight">
    <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
    <el-breadcrumb-item v-if="parentTitle" :to="parentPath">
      {{ parentTitle }}
    </el-breadcrumb-item>
    <el-breadcrumb-item>{{ currentTitle }}</el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowRight } from '@element-plus/icons-vue'
import { routerMenu } from './menu'

const route = useRoute()

const parentTitle = computed(() => {
  const parent = route.meta.parent as string | undefined
  if (!parent) return ''
  return routerMenu.find((m) => m.path === `/${parent}`)?.title ?? ''
})

const parentPath = computed(() => {
  const parent = route.meta.parent as string | undefined
  return parent ? `/${parent}` : ''
})

const currentTitle = computed(() => (route.meta.title as string) || '')
</script>

<style scoped lang="scss">
.breadcrumb {
  font-size: 13px;
  line-height: 1;
}
</style>
