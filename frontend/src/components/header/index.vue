<template>
  <div class="header">
    <div class="logo" @click="handleClickLogo">
      <img class="logo-img" src="@/assets/img/header/logo-light.png" alt="Leslie Blog Logo" />
      Leslie Blog
    </div>
    <div class="nav">
      <!-- 导航栏 -->
      <el-menu
        :popper-effect="effect"
        class="nav-menu"
        mode="horizontal"
        :default-active="activeIndex"
        router
      >
        <el-menu-item :index="item.index" v-for="item in menuList" :key="item.index">{{
          item.label
        }}</el-menu-item>
      </el-menu>
      <!-- 主题切换 -->
      <el-switch
        :model-value="bgSwitch"
        @change="handleChange"
        style="--el-switch-on-color: #f8fafc; --el-switch-off-color: #121826"
      >
        <template #active-action>
          <el-icon class="moon-icon"><Moon /></el-icon>
        </template>
        <template #inactive-action>
          <el-icon class="sunny-icon"><Sunny /></el-icon>
        </template>
      </el-switch>

      <!-- 语言切换 -->
      <el-dropdown placement="bottom" :effect="effect" @command="handleLangChange">
        <img
          class="language-icon"
          src="@/assets/img/header/switch-language-icon.png"
          alt="Switch Language Icon"
        />
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item
              v-for="item in dropdownMenuList"
              :key="item.command"
              :disabled="langStore.lang === item.command"
              :command="item.command"
              >{{ item.label }}</el-dropdown-item
            >
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <!-- GitHub 仓库 -->
      <GitHubIcon @click="handleClickGitHub" class="github-icon" />
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Sunny, Moon } from '@element-plus/icons-vue'
import { useThemeStore } from '@/stores/theme'
import { useLangStore } from '@/stores/language'
import { useI18n } from 'vue-i18n'
import GitHubIcon from '@/components/github-icon/index.vue'

// 初始化语言
const { locale, t } = useI18n()
const langStore = useLangStore()
const handleLangChange = (lang: 'zh' | 'en') => {
  langStore.setLang(lang)
  locale.value = lang
}
const dropdownMenuList = ref([
  {
    label: 'English',
    command: 'en',
  },
  {
    label: '简体中文',
    command: 'zh',
  },
])
// 初始化主题
const themeStore = useThemeStore()
const bgSwitch = computed(() => themeStore.isDark)
const effect = computed(() => (themeStore.isDark ? 'dark' : 'light'))
themeStore.initTheme()
const handleChange = () => {
  themeStore.toggleDark()
}
// 点击logo返回首页
const router = useRouter()
const handleClickLogo = () => {
  router.push({ path: '/' })
}

// 点击GitHub仓库跳转
const handleClickGitHub = () => {
  window.open('https://github.com/leslieZong/leslie-blog-manage', '_blank')
}
// 首页 文章 分类 项目 关于
const activeIndex = ref('/')
const menuList = computed(() => [
  {
    label: t('header.menu.home'),
    index: '/',
  },
  {
    label: t('header.menu.article'),
    index: '/article',
  },
  {
    label: t('header.menu.category'),
    index: '/category',
  },
  {
    label: t('header.menu.project'),
    index: '/project',
  },
  {
    label: t('header.menu.about'),
    index: '/about',
  },
])
watch(
  () => router.currentRoute.value.path,
  (newPath) => {
    activeIndex.value = newPath || '/'
  },
  { immediate: true },
)
</script>
<style scoped lang="scss">
.header {
  width: 100%;
  height: 60px;
  background-color: var(--el-bg-color-page);
  color: var(--el-text-color-primary);
  border-bottom: 1px solid var(--el-border-color);
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  .logo {
    display: flex;
    align-items: center;
    gap: 20px;
    cursor: pointer;
    font-size: 20px;
    font-weight: bold;
    .logo-img {
      height: 60px;
    }
  }
  .nav {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 10px;
    .nav-menu {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: flex-end;
    }
    .moon-icon,
    .sunny-icon {
      background-color: var(--el-bg-color-page);
      color: var(--el-text-color-primary);
      border-radius: 50%;
    }
    .language-icon {
      height: 40px;
      cursor: pointer;
    }
    .github-icon {
      height: 40px;
      cursor: pointer;
    }
  }
}
</style>
