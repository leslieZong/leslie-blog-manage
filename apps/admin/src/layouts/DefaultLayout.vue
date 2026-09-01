<template>
  <div class="app-layout">
    <Sidebar :collapsed="appStore.sidebarCollapsed" />
    <div class="app-layout__main">
      <Header />
      <div class="app-layout__content">
        <router-view v-slot="{ Component }">
          <keep-alive :include="cachedViews">
            <component :is="Component" />
          </keep-alive>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import Header from '@/components/Header.vue'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()
const router = useRouter()

// 缓存非 hidden 的列表页
const cachedViews = computed(() =>
  router
    .getRoutes()
    .filter((r) => r.name && !r.meta.hidden)
    .map((r) => r.name as string),
)
</script>
