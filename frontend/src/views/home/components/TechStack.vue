/** * Tech / 技术栈 */
<template>
  <div class="tech-stack">
    <div class="tech-stack-title">
      <div class="left">
        <h2 class="title">{{ $t('home.techStack.title') }}</h2>
        <div>{{ $t('home.techStack.description') }}</div>
      </div>
      <div class="right" @click="handleClickAll">
        <div class="explore-btn">{{ $t('home.techStack.viewAll') }}</div>
        <el-icon><ArrowRightBold /></el-icon>
      </div>
    </div>
    <div class="tech-stack-content">
      <div v-for="item in techStackList" :key="item.name" class="tech-node">
        <div class="node-title">{{ item.name }}</div>
        <div v-for="tech in item.list" :key="tech" class="tech-item">
          <div>◉</div>
          {{ tech }}
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
const techStackList = ref([
  {
    name: 'FRONTEND',
    list: ['Vue 3', 'TypeScript', 'Vite', 'JavaScript', 'CSS / SCSS / Less'],
  },
  {
    name: 'BACKEND',
    list: ['Go', 'Node.js', 'MySQL', 'REST API', 'SSE'],
  },
  {
    name: 'AI',
    list: ['LLM', 'AI Agents', 'Ollama', 'OpenAI', 'AI SDK'],
  },
  {
    name: 'TOOLS',
    list: ['Git', 'Docker', 'VS Code', 'pnpm'],
  },
  {
    name: 'ENGINEERING',
    list: ['Agile', 'Scrum', 'Kanban'],
  },
  {
    name: 'Notes',
    list: ['Markdown'],
  },
])
</script>
<style scoped lang="scss">
.tech-stack {
  padding-top: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--el-border-color);
  .tech-stack-title {
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
  .tech-stack-content {
    margin-top: 30px;
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    .tech-node {
      display: flex;
      flex-direction: column;
      justify-content: center;
      gap: 12px;
      padding: 20px;

      &:nth-child(-n + 4) {
        border-bottom: 1px solid var(--el-border-color);
      }
      .node-title {
        font-weight: bold;
        letter-spacing: 0.08em;
        text-transform: uppercase;
      }
      .tech-item {
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 6px;
        color: var(--el-text-color-regular);
        &:hover {
          transform: translateX(4px);
          color: var(--el-color-primary);
        }
      }
    }
  }
}
</style>
