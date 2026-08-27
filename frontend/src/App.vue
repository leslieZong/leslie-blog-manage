<template>
  <div class="app">
    <el-config-provider :locale="elementLocale">
      <Header v-if="hasHeader" class="header"></Header>
      <div class="content">
        <router-view v-slot="{ Component }">
          <keep-alive>
            <component :is="Component" />
          </keep-alive>
        </router-view>
        <Footer v-if="hasHeader" class="footer"></Footer>
      </div>
    </el-config-provider>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Header from '@/components/header/index.vue'
import Footer from '@/components/footer/index.vue'
import { useRoute } from 'vue-router'
import { useLangStore } from '@/stores/language'
import zhCn from 'element-plus/es/locale/lang/zh-cn.mjs'
import enUs from 'element-plus/es/locale/lang/en.mjs'

const langStore = useLangStore()
const elementLocale = computed(() => (langStore.lang === 'zh' ? zhCn : enUs))

const route = useRoute()
const hasHeader = computed(() => route.meta.hasHeader)
</script>

<style scoped lang="scss">
.app {
  width: 100%;
  height: 100%;
  background-color: var(--el-bg-color-page);
  color: var(--el-text-color-primary);
  .content {
    padding: 140px 20px 20px;
    flex: 1;
  }
}
</style>
