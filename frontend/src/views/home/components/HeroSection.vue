<template>
  <div class="hero-container">
    <!-- 背景点缀-->
    <div class="deco">
      <span
        v-for="d in decoList"
        :key="d.id"
        :style="{
          top: d.top,
          left: d.left,
          fontSize: d.size + 'px',
          animationDelay: d.delay + 's',
        }"
        >{{ d.symbol }}</span
      >
    </div>

    <div class="hero-content-container">
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
            <Icon :name="item.icon" size="24" class="social-icon" :alt="item.name" />
            {{ item.name }}
          </div>
        </div>
      </div>

      <!-- 右侧：个人技术生态图 -->
      <div class="hero-visual">
        <svg class="connections" viewBox="0 0 280 280" preserveAspectRatio="xMidYMid meet">
          <line class="line line-top" x1="140" y1="140" x2="140" y2="30" />
          <line class="line line-right" x1="140" y1="140" x2="250" y2="140" />
          <line class="line line-bottom" x1="140" y1="140" x2="140" y2="250" />
          <line class="line line-left" x1="140" y1="140" x2="30" y2="140" />
        </svg>

        <div class="hero-logo">
          <img src="@/assets/img/header/logo-light.png" />
        </div>

        <div class="tech-node vue">Vue</div>
        <div class="tech-node typescript">TypeScript</div>
        <div class="tech-node node">Go</div>
        <div class="tech-node ai">AI</div>
      </div>
    </div>
    <div class="link-up">
      <Icon name="arrowDown" size="24" class="link-icon" :alt="$t('home.hero.exploreBlog')" />
      {{ $t('home.hero.exploreBlog') }}
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icon/index.vue'
const stackList = ref(['Vue', 'TypeScript', 'Vite', 'Go', 'AI'])

// 背景点缀：固定位置，仅 twinkle 呼吸
interface DecoItem {
  id: number
  symbol: string
  top: string
  left: string
  size: number
  delay: number
}
const decoList: DecoItem[] = [
  { id: 0, symbol: '○', top: '10%', left: '4%', size: 14, delay: 0 },
  { id: 1, symbol: '✦', top: '8%', left: '38%', size: 18, delay: 0.8 },
  { id: 2, symbol: '/', top: '20%', left: '88%', size: 20, delay: 1.6 },
  { id: 3, symbol: '·', top: '85%', left: '12%', size: 22, delay: 2.2 },
  { id: 4, symbol: '○', top: '88%', left: '42%', size: 18, delay: 1.2 },
  { id: 5, symbol: '✦', top: '55%', left: '3%', size: 14, delay: 0.4 },
  { id: 6, symbol: '/', top: '30%', left: '20%', size: 16, delay: 1 },
  { id: 7, symbol: '·', top: '70%', left: '30%', size: 20, delay: 2.6 },
  { id: 8, symbol: '○', top: '15%', left: '70%', size: 16, delay: 0.6 },
  { id: 9, symbol: '✦', top: '75%', left: '75%', size: 18, delay: 1.8 },
  { id: 10, symbol: '/', top: '45%', left: '50%', size: 14, delay: 2 },
  { id: 11, symbol: '·', top: '60%', left: '65%', size: 22, delay: 0.2 },
]
const stack = computed(() => stackList.value.join('  ·  '))
const t = useI18n().t
const socialList = computed(() => [
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
  position: relative;
  isolation: isolate;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--el-border-color);
  .deco {
    position: absolute;
    inset: 0;
    z-index: -1;
    pointer-events: none;
    span {
      position: absolute;
      color: var(--el-text-color-secondary);
      font-size: 14px;
      opacity: 0.18;
      animation: twinkle 5s ease-in-out infinite;
    }
  }
  .hero-content-container {
    display: flex;
    justify-content: space-between;
    align-items: center;
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
          .social-icon {
            height: 20px;
          }
        }
      }
    }
    .hero-visual {
      position: relative;
      width: 360px;
      height: 360px;
      .connections {
        position: absolute;
        inset: 0;
        width: 100%;
        height: 100%;
        pointer-events: none;
        .line {
          stroke: var(--el-border-color);
          stroke-width: 1;
          stroke-dasharray: 3 6;
          animation: flow 4s linear infinite;
        }
      }
      .hero-logo {
        position: absolute;
        top: 50%;
        left: 50%;
        width: 90px;
        height: 90px;
        transform: translate(-50%, -50%);
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        background: var(--el-bg-color);
        box-shadow: 0 0 0 1px var(--el-border-color);
        animation: breathe 2s ease-in-out infinite;
        img {
          max-width: 70px;
          max-height: 70px;
          object-fit: contain;
        }
      }
      .tech-node {
        position: absolute;
        padding: 3px 10px;
        border-radius: 12px;
        border: 1px solid var(--el-border-color);
        background: var(--el-bg-color);
        font-size: 12px;
        color: var(--el-text-color-regular);
        white-space: nowrap;
        user-select: none;
        &.typescript {
          top: 8px;
          left: 50%;
          animation: floatY 2.5s ease-in-out infinite;
        }
        &.ai {
          top: 50%;
          right: 8px;
          animation: floatX 3s ease-in-out infinite;
        }
        &.node {
          bottom: 8px;
          left: 50%;
          animation: floatY 3.5s ease-in-out infinite;
        }
        &.vue {
          top: 50%;
          left: 8px;
          animation: floatX 4s ease-in-out infinite;
        }
      }
    }
  }
  .link-up {
    margin-top: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    color: var(--el-text-color-primary);
    .link-icon {
      animation: breathe2 3s ease-in-out infinite;
    }
  }
}

@keyframes flow {
  to {
    stroke-dashoffset: -9;
  }
}
@keyframes twinkle {
  0%,
  100% {
    opacity: 0.12;
    transform: scale(1);
  }
  50% {
    opacity: 0.3;
    transform: scale(1.15);
  }
}
@keyframes breathe {
  0%,
  100% {
    transform: translate(-50%, -50%) scale(1);
  }
  50% {
    transform: translate(-50%, -50%) scale(1.05);
  }
}
@keyframes breathe2 {
  0%,
  100% {
    transform: translate(-50%, -50%) scale(1);
  }
  50% {
    transform: translate(-50%, -50%) scale(1.5);
  }
}
@keyframes floatY {
  0%,
  100% {
    transform: translate(-50%, 0);
  }
  50% {
    transform: translate(-50%, -5px);
  }
}
@keyframes floatX {
  0%,
  100% {
    transform: translate(0, -50%);
  }
  50% {
    transform: translate(5px, -50%);
  }
}
</style>
