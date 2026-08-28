<!-- LatestPosts.vue -->
/** * Latest / 最新文章 */
<template>
  <div class="latest-posts">
    <div class="latest-posts-title">
      <div class="left">
        <h2 class="title">{{ $t('home.latestPosts.title') }}</h2>
        <div>{{ $t('home.latestPosts.description') }}</div>
      </div>
      <div class="right" @click="handleClickAll">
        <div class="explore-btn">{{ $t('home.latestPosts.viewAll') }}</div>
        <el-icon><ArrowRightBold /></el-icon>
      </div>
    </div>
    <div class="latest-posts-content">
      <div class="post-item" v-for="item in 5" :key="item">
        <img class="cover" src="@/assets/img/home/cover.png" alt="技术插画型封面" />
        <div class="post-info">
          <div class="top">
            <div class="type">Frontend</div>
            <div class="title">Vue3 响应式原理深度解析</div>
          </div>
          <div class="summary">
            深入理解 Vue3 深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
            的响应式系统,深入理解 Vue3 的响应式系统, 深入理解 Vue3 的响应式系统,深入理解 Vue3
            的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统, 深入理解 Vue3
            的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
            的响应式系统, 深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
            的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
            的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统, 深入理解 Vue3
            的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
            的响应式系统, 深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
            的响应式系统,深入理解 Vue3 的响应式系统, 深入理解 V,
          </div>
          <div class="meta">
            <div class="date">发布于： 2026-08-20</div>
            <div class="read-time">预估阅读时间： 8 min</div>
          </div>
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
.latest-posts {
  padding-top: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--el-border-color);
  .latest-posts-title {
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
  .latest-posts-content {
    margin-top: 30px;
    display: flex;
    flex-direction: column;
    gap: 20px;
    .post-item {
      padding: 20px;
      display: flex;
      align-items: center;
      gap: 20px;
      cursor: pointer;
      transition: transform 0.5s ease-in-out;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
      border-radius: 10px;
      background-color: var(--el-fill-color);
      &:hover {
        transform: scale(1.03);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
      }
      .cover {
        width: 100px;
        object-fit: cover;
      }
      .post-info {
        flex: 1;
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        gap: 6px;
        overflow: hidden;
        .top {
          display: flex;
          align-items: center;
          gap: 20px;
          .type {
            font-weight: bold;
            color: var(--el-color-primary);
          }
          .title {
            font-weight: bold;
            color: var(--el-text-color-primary);
          }
        }
        .summary {
          font-size: 16px;
          color: var(--el-text-color-secondary);
          text-overflow: ellipsis;
          white-space: nowrap;
          overflow: hidden;
        }
        .meta {
          display: flex;
          align-items: center;
          gap: 20px;
          color: var(--el-text-color-secondary);
          font-size: 14px;
        }
      }
    }
  }
}
</style>
