<template>
  <header class="header">
    <div class="header__left">
      <el-icon class="header__collapse" @click="appStore.toggleSidebar()">
        <Fold v-if="!appStore.sidebarCollapsed" />
        <Expand v-else />
      </el-icon>
      <Breadcrumb />
    </div>
    <div class="header__right">
      <el-tooltip content="切换主题" placement="bottom">
        <el-icon class="header__action" @click="themeStore.toggleDark()">
          <Moon v-if="!themeStore.isDark" />
          <Sunny v-else />
        </el-icon>
      </el-tooltip>
      <el-dropdown @command="onCommand">
        <span class="header__user">
          <el-avatar :size="28" :src="authStore.userInfo?.avatar">
            {{ authStore.userInfo?.nickname?.charAt(0) || 'A' }}
          </el-avatar>
          <span class="header__username">{{ authStore.userInfo?.nickname || 'Admin' }}</span>
          <el-icon><ArrowDown /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="settings">个人设置</el-dropdown-item>
            <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { logout as logoutApi } from '@/api/auth'

const appStore = useAppStore()
const authStore = useAuthStore()
const themeStore = useThemeStore()
const router = useRouter()

async function onCommand(cmd: string) {
  if (cmd === 'logout') {
    await ElMessageBox.confirm('确定退出登录吗？', '提示', { type: 'warning' })
    try {
      await logoutApi()
    } catch {
      /* 忽略，仍清理本地 */
    }
    authStore.clearAuth()
    router.push('/login')
  } else if (cmd === 'settings') {
    router.push('/settings')
  }
}
</script>

<style scoped lang="scss">
.header {
  height: var(--admin-header-height);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  background-color: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  flex-shrink: 0;

  &__left {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  &__collapse {
    font-size: 20px;
    cursor: pointer;
    color: var(--el-text-color-regular);
  }

  &__right {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  &__action {
    font-size: 20px;
    cursor: pointer;
    color: var(--el-text-color-regular);
  }

  &__user {
    display: flex;
    align-items: center;
    gap: 6px;
    cursor: pointer;
    color: var(--el-text-color-regular);
  }

  &__username {
    font-size: 14px;
  }
}
</style>
