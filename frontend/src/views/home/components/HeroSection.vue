<template>
  <div class="hero-container">
    <!-- 左侧 -->
    <div class="hero-content">
      <div>{{ $t('home.hero.greeting') }}</div>

      <div class="hero-title">
        <span v-html="$t('home.hero.title')"></span>
        <span v-html="$t('home.hero.subTitle')"></span>
      </div>

      <p class="hero-role">{{ $t('home.hero.role') }}</p>

      <p class="hero-description">
        {{ $t('home.hero.description') }}
      </p>

      <div class="hero-stack">{{ stack }}</div>

      <div class="hero-actions">
        <el-button type="primary">{{ $t('home.hero.exploreArticles') }}</el-button>
        <el-button>{{ $t('home.hero.viewProjects') }}</el-button>
      </div>

      <div class="hero-social">
        <div
          class="hero-social-item"
          v-for="item in socialList"
          :key="item.name"
          @click="handleClick(item.link)"
        >
          <BrandIcon :name="item.icon" :size="18" />
          {{ $t(item.name) }}
        </div>
      </div>
    </div>

    <!-- 右侧 -->
    <div class="hero-visual">
      <div class="hero-logo">
        <img src="@/assets/img/header/logo-light.png" />
      </div>

      <div class="tech-node vue">Vue</div>

      <div class="tech-node typescript">TypeScript</div>

      <div class="tech-node node">GoLang</div>

      <div class="tech-node ai">AI</div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import BrandIcon from '@/components/icon/index.vue'

const stackList = ref(['Vue', 'TypeScript', 'Vite', 'GoLang', 'AI'])
const stack = computed(() => stackList.value.join('  ·  '))
const t = useI18n().t
const socialList = ref([
  {
    name: t('home.hero.github'),
    link: 'https://github.com/leslieZong',
    icon: 'github',
  },
  {
    name: t('home.hero.rss'),
    link: 'https://leslie-blog-manage.vercel.app/rss',
    icon: 'rss',
  },
  {
    name: t('home.hero.email'),
    link: 'mailto:lesliezhaozp@gmail.com',
    icon: 'email',
  },
])

const handleClick = (link: string) => {
  window.open(link, '_blank')
}
</script>

<style scoped lang="scss">
.hero-container {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--el-border-color);
  .hero-content {
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 20px;
    .hero-title {
      display: flex;
      flex-direction: column;
      justify-content: center;
      font-size: 32px;
      font-weight: bold;
      :deep(.web),
      :deep(.ai) {
        background: linear-gradient(90deg, #3b82f6, #8b5cf6, #ec4899);
        background-clip: text;
        color: transparent;
      }
    }
    .hero-social {
      display: flex;
      align-items: center;
      gap: 20px;
      .hero-social-item {
        display: flex;
        align-items: center;
        gap: 5px;
        cursor: pointer;
        &:hover {
          color: var(--el-color-primary);
        }
      }
    }
  }
  .hero-visual {
    .hero-logo {
      width: 100px;
    }
  }
}
</style>
