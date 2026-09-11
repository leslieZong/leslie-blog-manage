<template>
  <el-config-provider :locale="zhCn">
    <router-view />
  </el-config-provider>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import zhCn from 'element-plus/es/locale/lang/zh-cn.mjs'
import { useThemeStore } from '@/stores/theme'
import {  getUserInfo as getUserInfoApi } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
const themeStore = useThemeStore()
themeStore.initTheme()
const authStore = useAuthStore()
onMounted(async () => {
  if (!authStore.token) return
      const res = await getUserInfoApi()
    console.log("🚀 ~ onSubmit ~ res:", res)
    //  authStore.setUserInfo(userInfo)
})
</script>

<style lang="scss">
html,
body,
#app {
  width: 100%;
  height: 100%;
}
</style>
