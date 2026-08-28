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
      <div class="category-item" v-for="(item, index) in 6" :key="index">
        <div class="category-sort">{{ '0' + (index + 1) }}</div>
        <div class="category-content">
          <div class="category-info">
            <div class="category-name">Frontend</div>
            <div class="category-count">
              10 <el-icon class="arrow-icon"><ArrowRightBold /></el-icon>
            </div>
          </div>
          <div class="category-description">Modern web interfaces, Vue, TypeScript & CSS</div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
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
    gap: 24px;
    margin-top: 30px;
    .category-item {
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
