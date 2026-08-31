<template>
  <main>
    <!-- 个人介绍 -->
    <HeroSection />

    <!-- 精选文章 -->
    <FeaturedPosts />

    <!-- 最新文章 -->
    <LatestPosts />

    <!-- 分类 -->
    <CategorySection />

    <!-- 项目 -->
    <ProjectSection />

    <!-- 技术栈 -->
    <TechStack />

    <!-- GitHub 活跃度 -->
    <GithubActivity />

    <!-- 关于我 -->
    <!-- <AboutSection /> -->
  </main>
</template>

<script setup lang="ts">
import HeroSection from './components/HeroSection.vue'
import FeaturedPosts from './components/FeaturedPosts.vue'
import LatestPosts from './components/LatestPosts.vue'
import CategorySection from './components/CategorySection.vue'
import ProjectSection from './components/ProjectSection.vue'
import TechStack from './components/TechStack.vue'
import GithubActivity from './components/GithubActivity.vue'
// import AboutSection from './components/AboutSection.vue'
</script>

<style scoped lang="scss"></style>
