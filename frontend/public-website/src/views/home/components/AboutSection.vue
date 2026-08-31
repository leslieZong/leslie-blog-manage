<!-- AboutSection.vue -->
/** * About / 关于 个人介绍 */
<template>
  <div class="about-section">
    <div class="about-title">
      <div class="left">
        <h2 class="title">{{ $t('home.about.title') }}</h2>
        <div>{{ $t('home.about.description') }}</div>
      </div>
    </div>

    <div class="about-content">
      <div class="about-bio">
        <div class="name">{{ $t('home.about.name') }}</div>
        <div class="role">{{ $t('home.about.role') }}</div>
        <p class="bio-text">{{ $t('home.about.bio1') }}</p>
        <p class="bio-text">{{ $t('home.about.bio2') }}</p>
      </div>

      <div class="about-currently">
        <div class="section-label">{{ $t('home.about.currentlyTitle') }}</div>
        <div class="currently-list">
          <div
            class="currently-item"
            v-for="item in currentlyList"
            :key="item.label"
          >
            <div class="currently-label">{{ item.label }}</div>
            <div class="currently-desc">{{ item.desc }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="about-connect">
      <div class="section-label">{{ $t('home.about.connectTitle') }}</div>
      <div class="connect-list">
        <div
          class="connect-item"
          v-for="item in connectList"
          :key="item.name"
          @click="handleClick(item.link)"
        >
          {{ item.name }} <span class="arrow">↗</span>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { isSafeUrl } from '@/utils/sanitize'

const { t } = useI18n()

const currentlyList = computed(() => [
  { label: t('home.about.building'), desc: t('home.about.buildingDesc') },
  { label: t('home.about.exploring'), desc: t('home.about.exploringDesc') },
  { label: t('home.about.writing'), desc: t('home.about.writingDesc') },
])

const connectList = computed(() => [
  { name: t('home.about.github'), link: 'https://github.com/leslieZong' },
  { name: t('home.about.linkedin'), link: 'https://www.linkedin.com/' },
  { name: t('home.about.email'), link: 'mailto:lesliezhaozp@gmail.com' },
])

const handleClick = (link: string) => {
  if (!isSafeUrl(link)) return
  window.open(link, '_blank', 'noopener,noreferrer')
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
  }

  .about-content {
    margin-top: 30px;
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 40px;
    .about-bio {
      .name {
        font-size: 32px;
        font-weight: bold;
        letter-spacing: 0.02em;
      }
      .role {
        margin-top: 8px;
        font-size: 16px;
        color: var(--el-text-color-secondary);
      }
      .bio-text {
        margin-top: 20px;
        color: var(--el-text-color-regular);
        line-height: 1.8;
      }
    }
    .about-currently {
      .section-label {
        font-weight: bold;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        color: var(--el-text-color-secondary);
        margin-bottom: 20px;
      }
      .currently-list {
        display: flex;
        flex-direction: column;
      }
      .currently-item {
        display: flex;
        align-items: baseline;
        gap: 16px;
        padding: 16px 0;
        border-bottom: 1px solid var(--el-border-color);
        &:last-child {
          border-bottom: none;
        }
        .currently-label {
          width: 80px;
          flex-shrink: 0;
          font-weight: bold;
          color: var(--el-text-color-primary);
        }
        .currently-desc {
          color: var(--el-text-color-regular);
        }
      }
    }
  }

  .about-connect {
    margin-top: 40px;
    padding-top: 30px;
    border-top: 1px solid var(--el-border-color);
    display: flex;
    align-items: center;
    justify-content: space-between;
    .section-label {
      font-weight: bold;
      letter-spacing: 0.1em;
      text-transform: uppercase;
      color: var(--el-text-color-secondary);
    }
    .connect-list {
      display: flex;
      align-items: center;
      gap: 40px;
    }
    .connect-item {
      display: flex;
      align-items: center;
      gap: 4px;
      cursor: pointer;
      font-weight: bold;
      transition: color 0.5s ease-in-out;
      &:hover {
        color: var(--el-color-primary);
        .arrow {
          transform: translate(2px, -2px);
        }
      }
      .arrow {
        font-size: 12px;
        transition: transform 0.5s ease-in-out;
      }
    }
  }
}
</style>
