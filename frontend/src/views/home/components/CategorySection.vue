<!-- CategorySection.vue -->
/** * * Category / 分类 */
<template>
  <div class="category-section">
    <div class="category-title">
      <div class="left">
        <h2 class="title">{{ $t('home.category.title') }}</h2>
        <div>{{ $t('home.category.description') }}</div>
      </div>
      <div class="right" @click="handleClickAll">
        <div class="explore-btn">{{ $t('home.category.viewAll') }}</div>
        <el-icon><ArrowRightBold /></el-icon>
      </div>
    </div>
    <div class="category-list">
      <div class="category-item" v-for="(item, index) in categories" :key="index">
        <div class="category-sort">{{ '0' + (index + 1) }}</div>
        <div class="category-content">
          <div class="category-info">
            <div class="category-name">
              <Icon :name="item.icon" size="20" class="link-icon" />
              {{ item.name }}
            </div>
            <div class="category-count">
              {{ item.count }} <el-icon class="arrow-icon"><ArrowRightBold /></el-icon>
            </div>
          </div>
          <div class="category-description">{{ item.description }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import Icon from '@/components/icon/index.vue'

const categories = ref([
  {
    name: 'Frontend',
    count: 10,
    description: 'Modern web development, Vue, TypeScript & CSS',
    icon: 'frontend',
  },
  {
    name: 'AI',
    count: 5,
    description: 'AI applications, agents and practical experiments',
    icon: 'ai',
  },
  {
    name: 'Engineering',
    count: 8,
    description: 'Architecture, patterns and engineering practices',
    icon: 'engineering',
  },
  {
    name: 'Backend',
    count: 3,
    description: 'Go services, APIs and backend engineering',
    icon: 'backend',
  },
  {
    name: 'Developer Tools ',
    count: 2,
    description: 'Tools and workflows that improve development',
    icon: 'developerTools',
  },
  {
    name: 'Notes',
    count: 2,
    description: 'Thoughts, learning notes and things worth sharing',
    icon: 'notes',
  },
])

// 点击查看所有
const handleClickAll = () => {
  console.log('点击了查看所有')
}
</script>
<style scoped lang="scss">
.category-section {
  padding-top: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--el-border-color);
  .category-title {
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
  .category-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-top: 30px;
    .category-item {
      display: flex;
      align-items: center;
      gap: 12px;
      overflow: hidden;
      padding: 12px;
      border-bottom: 1px solid var(--el-border-color);
      cursor: pointer;
      transition: all 0.5s ease-in-out;
      &:last-child {
        border-bottom: none;
      }
      .category-sort {
        width: 40px;
        font-size: 16px;
      }
      &:hover {
        transform: translateX(4px);
        .category-sort {
          font-weight: bold;
        }
        .category-name {
          color: var(--el-color-primary);
        }

        .arrow-icon {
          opacity: 1;
        }
        .category-description {
          color: var(--el-text-color-primary);
        }
      }
      .arrow-icon {
        transition: all 0.5s ease-in-out;
        opacity: 0;
      }
      .category-content {
        flex: 1;
        display: flex;
        flex-direction: column;
      }
      .category-info {
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 6px;
        .category-name {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 18px;
          font-weight: bold;
        }
        .category-count {
          font-size: 16px;
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
      .category-description {
        font-size: 16px;
        color: var(--el-text-color-secondary);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        transition: all 0.5s ease-in-out;
      }
    }
  }
}
</style>
