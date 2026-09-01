<!-- AboutSection.vue -->
/** * About / 关于 个人介绍 */
<template>
  <div class="about-section">
    <div class="about-title">
      <div class="left">
        <h2 class="title">{{ $t('home.about.title') }}</h2>
        <div>{{ $t('home.about.description') }}</div>
      </div>
      <div class="right" @click="handleClickAll">
        <div class="explore-btn">{{ $t('home.about.viewAll') }}</div>
        <el-icon><ArrowRightBold /></el-icon>
      </div>
    </div>
    <div class="about-content">
      <div v-for="(item, index) in aboutList" :key="item.title" class="about-item">
        <div class="sort">{{ '0' + (index + 1) }}</div>
        <div class="content">
          <h3 class="title">{{ item.title }}</h3>
          <div class="desc">{{ item.desc }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const aboutList = computed(() => {
  return [
    {
      title: t('home.about.whoAmI'),
      desc: t('home.about.IamDesc'),
    },
    {
      title: t('home.about.whatICareAbout'),
      desc: t('home.about.whatICareAboutDesc'),
    },
    {
      title: t('home.about.currentlyTitle'),
      desc: t('home.about.currentlyDesc'),
    },
    {
      title: t('home.about.principalTitle'),
      desc: t('home.about.principalDesc'),
    },
  ]
})
const handleClickAll = () => {
  console.log('点击了查看所有')
}
</script>
<style scoped lang="scss">
.about-section {
  padding-top: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--el-border-color);
  .about-title {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .left {
      .title {
        margin-bottom: 12px;
      }
    }
    .right {
      display: flex;
      align-items: center;
      gap: 6px;
      cursor: pointer;
      &:hover {
        color: var(--el-color-primary);
      }
    }
  }
  .about-content {
    margin-top: 30px;
    display: flex;
    flex-direction: column;
    gap: 12px;
    .about-item {
      display: flex;
      align-items: flex-start;
      gap: 24px;
      padding-bottom: 12px;
      border-bottom: 1px solid var(--el-border-color);
      &:last-child {
        border-bottom: none;
      }
      .content {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 12px;
        .title {
          font-weight: bold;
          color: var(--el-text-color-primary);
        }
        .desc {
          color: var(--el-text-color-regular);
        }
      }
    }
  }
}
</style>
