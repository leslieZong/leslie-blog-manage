<!-- ProjectSection.vue -->
/** * Project / 项目 */
<template>
  <div class="project-section">
    <div class="project-title">
      <div class="left">
        <h2 class="title">{{ $t('home.project.title') }}</h2>
        <div>{{ $t('home.project.description') }}</div>
      </div>
      <div class="right" @click="handleClickAll">
        <div class="explore-btn">{{ $t('home.project.viewAll') }}</div>
        <el-icon><ArrowRightBold /></el-icon>
      </div>
    </div>
    <div class="project-list">
      <div class="project-item" v-for="(item, index) in projectList" :key="item.id">
        <div class="project-sort">{{ '0' + (index + 1) }}</div>
        <div class="project-info">
          <div class="project-info-top">
            <div class="project-name">{{ item.name }}</div>
            <div class="project-status-date">
              <div class="project-status">{{ item.status }}</div>
              <div>·</div>
              <div class="project-date">{{ item.date }}</div>
            </div>
          </div>
          <div class="project-desc">{{ item.desc }}</div>
          <div class="project-tech-stack">{{ item.techStack }}</div>
          <div class="project-info-bottom">
            <div class="project-live">{{ item.live }} ↗</div>
            <div class="project-github">{{ item.github }} ↗</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
// 点击查看所有
const handleClickAll = () => {
  console.log('点击了查看所有')
}
const projectList = ref([
  {
    id: 1,
    sort: '01',
    name: 'LeslieBlog',
    desc: 'Personal technical blog and knowledge base.',
    techStack: 'Vue 3 · TypeScript · Go · MySQL',
    date: '2026',
    status: 'Active',
    live: 'Live Demo',
    github: 'GitHub',
  },
])
</script>
<style scoped lang="scss">
.project-section {
  padding-top: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--el-border-color);
  .project-title {
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
  .project-list {
    display: flex;
    flex-direction: column;
    gap: 24px;
    margin-top: 30px;
    .project-item {
      display: flex;
      align-items: center;
      gap: 12px;
      overflow: hidden;
      padding: 12px;
      background-color: linear-gradient(to right, #111, #666, #e5e5e5);
      border-bottom: 1px solid var(--el-border-color);
      cursor: pointer;
      transition: all 0.5s ease-in-out;
      &:last-child {
        border-bottom: none;
      }
      &:hover {
        transform: translateX(4px);
        background-color: var(--el-fill-color);
        border-radius: 8px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
      }
      .project-sort {
        width: 40px;
        font-size: 16px;
        font-weight: bold;
      }
      .project-info {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 12px;
        .project-info-top {
          display: flex;
          justify-content: space-between;
          align-items: center;
          .project-name {
            font-size: 18px;
            font-weight: bold;
          }
          .project-status-date {
            display: flex;
            align-items: center;
            gap: 12px;
          }
          .project-date {
            font-size: 14px;
          }
        }
        .project-desc {
          color: var(--el-text-color-secondary);
        }
        .project-info-bottom {
          display: flex;
          align-items: center;
          gap: 30px;
        }
      }
    }
  }
}
</style>
