<template>
  <div class="header" :class="{ 'is-scrolled': isScrolled }">
    <div class="header-inner">
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
        <!-- 查询 -->
        <el-icon class="search-icon" @click="handleClickSearch"><Search /></el-icon>
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
    <!-- 查询弹窗 -->
    <SearchDialog v-if="searchDialogVisible" v-model="searchDialogVisible" />
  </div>
</template>
<script setup lang="ts">
import { computed, ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { Sunny, Moon } from '@element-plus/icons-vue'
import { useThemeStore } from '@/stores/theme'
import { useLangStore } from '@/stores/language'
import { useI18n } from 'vue-i18n'
import GitHubIcon from '@/components/github-icon/index.vue'
import SearchDialog from '@/components/search-dialog/index.vue'

// 初始化语言
const { locale, t } = useI18n()
const langStore = useLangStore()
const handleLangChange = (lang: 'zh' | 'en') => {
  langStore.setLang(lang)
  locale.value = lang
}
// 语言切换菜单
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
// 导航栏
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
// 监听路由变化，更新 activeIndex
watch(
  () => router.currentRoute.value.path,
  (newPath) => {
    activeIndex.value = newPath || '/'
  },
  { immediate: true },
)

// 点击查询
const searchDialogVisible = ref(false)
const handleClickSearch = () => {
  searchDialogVisible.value = true
}

const isScrolled = ref(false)
// 监听滚动事件，更新 isScrolled
const handleScroll = () => {
  console.log(window.scrollY)
  isScrolled.value = window.scrollY > 20
}
onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>
<style scoped lang="scss">
.header {
  position: fixed;
  top: 20px;
  left: 50%;
  z-index: 1000;
  width: calc(100% - 40px);
  max-width: 1180px;
  transform: translateX(-50%);
  transition:
    top 0.3s ease,
    width 0.3s ease;
  &.is-scrolled {
    top: 12px;

    .header-inner {
      height: 58px;
      border-radius: 16px;
      box-shadow: 0 10px 35px rgba(0, 0, 0, 0.1);
    }
  }
  .header-inner {
    padding-right: 20px;
    border: 1px solid rgba(0, 0, 0, 0.06);
    height: 68px;
    border-radius: 20px;
    backdrop-filter: blur(18px);
    -webkit-backdrop-filter: blur(18px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.06);
    background-color: var(--el-bg-color);
    color: var(--el-text-color-primary);
    display: flex;
    align-items: center;
    justify-content: space-between;
    .logo {
      display: flex;
      align-items: center;
      cursor: pointer;
      font-size: 20px;
      font-weight: bold;
      .logo-img {
        height: 40px;
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
      :deep(.el-menu--horizontal.el-menu) {
        border-bottom: none;
      }
      :deep(.el-menu--horizontal > .el-menu-item) {
        &.is-active {
          font-weight: bold;
        }
      }
      .search-icon {
        height: 40px;
        cursor: pointer;
        margin-right: 12px;
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
}
</style>
