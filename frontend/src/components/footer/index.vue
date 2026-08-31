<!-- Footer.vue -->
/** * Footer / 页脚 */
<template>
  <footer class="footer">
    <!-- 主区域 -->
    <div class="footer-main">
      <div class="footer-col">
        <div class="col-label">{{ $t('footer.exploreTitle') }}</div>
        <div class="col-links">
          <div
            class="link-item"
            v-for="item in navList"
            :key="item.path"
            @click="handleNav(item.path)"
          >
            {{ item.label }}
          </div>
        </div>
      </div>
      <div class="footer-col">
        <div class="col-label">{{ $t('footer.connectTitle') }}</div>
        <div class="col-links">
          <div
            class="link-item"
            v-for="item in socialList"
            :key="item.name"
            @click="handleClick(item.link)"
          >
            {{ item.name }} <span class="arrow">↗</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部栏 -->
    <div class="footer-bottom">
      <div class="copyright">{{ $t('footer.copyright') }}</div>
      <div class="built-with" @click="handleClick('https://vuejs.org')">
        {{ $t('footer.builtWith') }} ↗
      </div>
      <div class="back-to-top" @click="handleBackToTop">{{ $t('footer.backToTop') }} ↑</div>
    </div>
  </footer>
</template>
<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { isSafeUrl } from '@/utils/sanitize'

const { t } = useI18n()
const router = useRouter()

const navList = computed(() => [
  { label: t('footer.nav.home'), path: '/' },
  { label: t('footer.nav.blog'), path: '/article' },
  { label: t('footer.nav.projects'), path: '/project' },
  { label: t('footer.nav.about'), path: '/about' },
])

const socialList = computed(() => [
  { name: t('footer.social.github'), link: 'https://github.com/leslieZong' },
  { name: t('footer.social.linkedin'), link: 'https://www.linkedin.com/' },
  { name: t('footer.social.rss'), link: 'https://leslie-blog-manage.vercel.app/rss' },
  { name: t('footer.social.email'), link: 'mailto:lesliezhaozp@gmail.com' },
])

const handleNav = (path: string) => {
  router.push(path)
}

const handleClick = (link: string) => {
  if (!isSafeUrl(link)) return
  window.open(link, '_blank', 'noopener,noreferrer')
}

const handleBackToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>
<style scoped lang="scss">
.footer {
  margin-top: 40px;

  .footer-main {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    .footer-col {
      display: flex;
      align-items: flex-start;
      gap: 24px;
      .col-label {
        font-weight: bold;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        color: var(--el-text-color-secondary);
        margin-bottom: 16px;
      }
      .col-links {
        display: flex;
        align-items: flex-start;
        gap: 30px;
      }
      .link-item {
        cursor: pointer;
        color: var(--el-text-color-regular);
        transition: all 0.5s ease-in-out;
        &:hover {
          color: var(--el-color-primary);
          transform: translateX(4px);
        }
        .arrow {
          font-size: 12px;
        }
      }
    }
  }

  .footer-bottom {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 14px;
    color: var(--el-text-color-secondary);
    .built-with,
    .back-to-top {
      cursor: pointer;
      transition: color 0.5s ease-in-out;
      &:hover {
        color: var(--el-color-primary);
      }
    }
  }
}
</style>
