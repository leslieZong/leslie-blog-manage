<template>
  <aside class="sidebar" :class="{ 'is-collapsed': collapsed }">
    <div class="sidebar__logo">
      <img src="/favicon.svg" alt="logo" class="sidebar__logo-img" />
      <span v-show="!collapsed" class="sidebar__logo-text">Leslie Admin</span>
    </div>
    <el-scrollbar class="sidebar__scroll">
      <el-menu
        :default-active="activeMenu"
        :collapse="collapsed"
        :collapse-transition="false"
        router
        background-color="transparent"
        text-color="#c9cdd4"
        active-text-color="#fff"
      >
        <template v-for="item in menuItems" :key="item.path">
          <el-menu-item :index="item.path">
            <el-icon>
              <component :is="item.icon" />
            </el-icon>
            <template #title>{{ item.title }}</template>
          </el-menu-item>
        </template>
      </el-menu>
    </el-scrollbar>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import * as ElIcons from '@element-plus/icons-vue'
import { routerMenu } from './menu'

defineProps<{ collapsed: boolean }>()

const route = useRoute()
const activeMenu = computed(() => {
  // 编辑页高亮父级 posts
  if (route.meta.parent) return `/${route.meta.parent}`
  return route.path
})

const menuItems = computed(() =>
  routerMenu.map((m) => ({
    ...m,
    icon: (ElIcons as Record<string, unknown>)[m.icon],
  })),
)
</script>

<style scoped lang="scss">
.sidebar {
  width: var(--admin-sidebar-width);
  height: 100%;
  background-color: #1a1f29;
  display: flex;
  flex-direction: column;
  transition: width 0.28s;
  overflow: hidden;
  flex-shrink: 0;

  &.is-collapsed {
    width: var(--admin-sidebar-collapsed-width);
  }

  &__logo {
    height: var(--admin-header-height);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    color: #fff;
    font-weight: 700;
    font-size: 16px;
    border-bottom: 1px solid #2c3038;
    flex-shrink: 0;
  }

  &__logo-img {
    width: 28px;
    height: 28px;
    border-radius: 6px;
  }

  &__scroll {
    flex: 1;
    min-height: 0;
  }

  :deep(.el-menu) {
    border-right: none;
  }

  :deep(.el-menu-item) {
    &.is-active {
      background-color: var(--el-color-primary) !important;
    }
    &:hover {
      background-color: rgba(255, 255, 255, 0.06) !important;
    }
  }
}
</style>
