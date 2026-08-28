<!-- FeaturedPosts.vue -->
/** * Featured / 精选文章 */
<template>
  <div class="featured-posts">
    <div class="featured-posts-title">
      <div class="left">
        <h2 class="title">{{ $t('home.featuredArticles.title') }}</h2>
        <div>{{ $t('home.featuredArticles.description1') }}</div>
        <div>{{ $t('home.featuredArticles.description2') }}</div>
      </div>
      <div class="right" @click="handleClickAll">
        <div class="explore-btn">{{ $t('home.featuredArticles.viewAll') }}</div>
        <el-icon><ArrowRightBold /></el-icon>
      </div>
    </div>
    <div class="featured-posts-content">
      <div class="posts-content-left">
        <img class="cover" src="@/assets/img/home/cover.png" alt="技术插画型封面" />
        <div class="type">Frontend</div>
        <div class="title">Vue3 响应式原理深度解析</div>
        <div class="summary">
          深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解
          Vue3 的响应式系统, 深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统, 深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统, 深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统, 深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统, 深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统, 深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统, 深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统, 深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统, 深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统, 深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统, 深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统, 深入理解 Vue3
          的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
          的响应式系统
        </div>
        <div class="meta">
          <div class="date">发布于： 2026-08-20</div>
          <div class="read-time">预估阅读时间： 8 min</div>
        </div>
      </div>
      <div class="posts-content-right">
        <div class="category">
          <img class="cover" src="@/assets/img/home/cover.png" alt="技术插画型封面" />
          <div class="type">Frontend</div>
          <div class="title">Vue3 响应式原理深度解析</div>
          <div class="summary">
            深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
            的响应式系统,深入理解 Vue3 的响应式系统
          </div>
          <div class="meta">
            <div class="date">发布于： 2026-08-20</div>
            <div class="read-time">预估阅读时间： 8 min</div>
          </div>
        </div>
        <div class="category">
          <img class="cover" src="@/assets/img/home/cover.png" alt="技术插画型封面" />
          <div class="type">Frontend</div>
          <div class="title">Vue3 响应式原理深度解析</div>
          <div class="summary">
            深入理解 Vue3 的响应式系统,深入理解 Vue3 的响应式系统,深入理解 Vue3
            的响应式系统,深入理解 Vue3 的响应式系统
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
.featured-posts {
  padding-top: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--el-border-color);
  .featured-posts-title {
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
  .featured-posts-content {
    margin-top: 30px;
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 24px;
    .posts-content-left {
      background-color: var(--el-fill-color);
      padding: 24px;
      border-radius: 12px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      display: flex;
      flex-direction: column;
      gap: 20px;
      cursor: pointer;
      .cover {
        width: 100%;
        object-fit: cover;
      }
      .type {
        font-size: 20px;
        font-weight: bold;
        color: var(--el-color-primary);
      }
      .title {
        font-size: 24px;
        font-weight: bold;
        color: var(--el-text-color-primary);
        margin-bottom: 12px;
      }
      .summary {
        color: var(--el-text-color-secondary);
        display: -webkit-box;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 6; /* 最多显示 6 行 */
        overflow: hidden;
        text-overflow: ellipsis;
      }
      .meta {
        display: flex;
        align-items: center;
        justify-content: space-between;
        color: var(--el-text-color-secondary);
      }
    }
    .posts-content-right {
      display: flex;
      flex-direction: column;
      gap: 24px;
      .category {
        flex: 1;
        background-color: var(--el-fill-color);
        padding: 24px;
        border-radius: 12px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        display: flex;
        flex-direction: column;
        gap: 10px;
        cursor: pointer;
        .cover {
          width: 100%;
        }
        .type {
          font-size: 18px;
          font-weight: bold;
          color: var(--el-color-primary);
        }
        .title {
          font-size: 20px;
          font-weight: bold;
          color: var(--el-text-color-primary);
        }
        .summary {
          color: var(--el-text-color-secondary);
          display: -webkit-box;
          -webkit-box-orient: vertical;
          -webkit-line-clamp: 2; /* 最多显示 2 行 */
          overflow: hidden;
          text-overflow: ellipsis;
        }
        .meta {
          display: flex;
          align-items: center;
          justify-content: space-between;
          color: var(--el-text-color-secondary);
        }
      }
    }
  }
}
</style>
