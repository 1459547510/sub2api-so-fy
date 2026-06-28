<template>
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="homeContent"></div>
  </div>

  <div v-else class="ag-home" :class="{ 'ag-home-dark': isDark }">
    <header class="ag-header">
      <router-link to="/home" class="ag-brand" :aria-label="copy.homeAria">
        <img v-if="siteLogo" :src="siteLogo" :alt="siteName" class="ag-logo" />
        <span v-else class="ag-mark" aria-hidden="true"></span>
        <span>{{ siteName }} {{ copy.brandSuffix }}</span>
      </router-link>

      <nav class="ag-nav" :aria-label="copy.mainNavAria">
        <a href="#products">{{ copy.navProducts }}</a>
        <a href="#use-cases">{{ copy.navUseCases }}</a>
        <router-link to="/models">{{ copy.navModels }}</router-link>
        <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer">{{ copy.navResources }}</a>
      </nav>

      <div class="ag-actions">
        <div class="ag-language">
          <LocaleSwitcher />
        </div>
        <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="ag-console">
          {{ isAuthenticated ? t('home.dashboard') : t('home.getStarted') }}
          <Icon name="arrowRight" size="sm" :stroke-width="2" />
        </router-link>
      </div>
    </header>

    <main>
      <section class="ag-hero">
        <div class="ag-brand-fade">
          <img v-if="siteLogo" :src="siteLogo" :alt="siteName" class="ag-logo small" />
          <span v-else class="ag-mark small" aria-hidden="true"></span>
          <span>{{ siteName }} {{ copy.brandSuffix }}</span>
        </div>

        <h1>{{ copy.heroTitle }}</h1>
        <p v-if="siteSubtitle" class="ag-site-subtitle">{{ siteSubtitle }}</p>

        <div class="ag-hero-orbit" aria-hidden="true">
          <div v-for="tool in toolIcons" :key="tool.className" class="ag-tool" :class="tool.className">
            <Icon :name="tool.icon" size="lg" />
            <span>{{ tool.label }}</span>
          </div>
          <div class="ag-liftoff-core">
            <img v-if="siteLogo" :src="siteLogo" :alt="siteName" class="ag-logo core" />
            <span v-else class="ag-mark core" aria-hidden="true"></span>
          </div>
        </div>

        <div class="ag-hero-actions">
          <router-link :to="isAuthenticated ? dashboardPath : '/login'" class="ag-primary">
            <Icon name="terminal" size="sm" />
            {{ isAuthenticated ? t('home.dashboard') : t('home.getStarted') }}
          </router-link>
          <router-link to="/models" class="ag-secondary">{{ copy.exploreModels }}</router-link>
        </div>
      </section>

      <section class="ag-orbit ag-liftoff" :aria-label="copy.previewAria">
        <div class="ag-preview-card ag-float-window">
          <div class="ag-window-bar">
            <span></span><span></span><span></span>
            <strong>{{ copy.previewTitle }}</strong>
          </div>
          <div class="ag-workspace">
            <aside>
              <span>{{ copy.previewProjects }}</span>
              <strong>{{ siteName }}</strong>
              <small v-if="siteSubtitle">{{ siteSubtitle }}</small>
              <em>{{ copy.previewRouting }}</em>
              <em>{{ copy.previewMapping }}</em>
              <em>{{ copy.previewBilling }}</em>
              <em v-if="apiBaseUrl">{{ copy.previewApiBase }} {{ apiBaseUrl }}</em>
              <em v-if="contactInfo">{{ copy.previewContact }} {{ contactInfo }}</em>
            </aside>
            <section>
              <p class="ag-code-line"><b>{{ copy.previewAgentLabel }}</b>{{ copy.previewAgentText }}</p>
              <p class="ag-code-line"><b>{{ copy.previewModelLabel }}</b>{{ copy.previewModelText }}</p>
              <p class="ag-code-line"><b>{{ copy.previewUsageLabel }}</b>{{ copy.previewUsageText }}</p>
            </section>
          </div>
        </div>
      </section>

      <section id="use-cases" class="ag-product-grid">
        <article v-for="product in products" :key="product.title" class="ag-product-card">
          <div class="ag-product-icon">
            <Icon :name="product.icon" size="lg" />
          </div>
          <span>{{ product.kicker }}</span>
          <h2>{{ product.title }}</h2>
          <p>{{ product.description }}</p>
        </article>
      </section>
    </main>

    <footer class="ag-footer">
      <span>&copy; {{ currentYear }} {{ siteName }}</span>
      <span v-if="siteSubtitle" class="ag-footer-meta">{{ siteSubtitle }}</span>
      <span v-if="apiBaseUrl" class="ag-footer-meta">API: {{ apiBaseUrl }}</span>
      <span v-if="contactInfo" class="ag-footer-meta">{{ contactInfo }}</span>
      <LocaleSwitcher placement="top" />
      <button type="button" @click="toggleTheme">{{ isDark ? copy.lightTheme : copy.darkTheme }}</button>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'

const { t, locale } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() =>
  sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', {
    allowRelative: true,
    allowDataUrl: true
  })
)
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || '')
const apiBaseUrl = computed(() => appStore.cachedPublicSettings?.api_base_url || appStore.apiBaseUrl || '')
const contactInfo = computed(() => appStore.cachedPublicSettings?.contact_info || appStore.contactInfo || '')
const docUrl = computed(() =>
  sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '', {
    allowRelative: true
  })
)
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

const isDark = ref(document.documentElement.classList.contains('dark'))
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const currentYear = computed(() => new Date().getFullYear())
const isZh = computed(() => locale.value.startsWith('zh'))

type IconName = InstanceType<typeof Icon>['$props']['name']

const copy = computed(() => {
  if (isZh.value) {
    return {
      homeAria: '首页',
      mainNavAria: '主导航',
      brandSuffix: 'AI API 网关',
      navProducts: '解决方案',
      navUseCases: '核心功能',
      navModels: '支持模型',
      navResources: '文档',
      heroTitle: '一个密钥，畅用多个 AI 模型',
      exploreModels: '查看支持模型',
      previewAria: '平台预览',
      previewTitle: 'AI 服务中转平台',
      previewProjects: '快速开始',
      previewRouting: '创建 API Key',
      previewMapping: '配置模型映射',
      previewBilling: '设置计费配额',
      previewApiBase: 'API 地址',
      previewContact: '联系',
      previewAgentLabel: '接入',
      previewAgentText: ' 一个 API 密钥调用所有已接入模型',
      previewModelLabel: '映射',
      previewModelText: ' 按模型别名匹配已接入服务',
      previewUsageLabel: '计费',
      previewUsageText: ' 按实际使用量计费，支持配额上限',
      lightTheme: '浅色',
      darkTheme: '深色',
      toolLabels: ['订阅转 API', '会话保持', '按量计费', '一键接入', '稳定可靠', '模型映射', '配额控制', '用量明细']
    }
  }

  return {
    homeAria: 'Home',
    mainNavAria: 'Main navigation',
    brandSuffix: 'AI API Gateway',
    navProducts: 'Solutions',
    navUseCases: 'Features',
    navModels: 'Models',
    navResources: 'Docs',
    heroTitle: 'One Key, All AI Models',
    exploreModels: 'View supported models',
    previewAria: 'Platform preview',
    previewTitle: 'AI service gateway',
    previewProjects: 'Quick start',
    previewRouting: 'Create API Key',
    previewMapping: 'Configure model mapping',
    previewBilling: 'Set billing quotas',
    previewApiBase: 'API base',
    previewContact: 'Contact',
    previewAgentLabel: 'access',
    previewAgentText: ' call all connected models with one API key',
    previewModelLabel: 'mapping',
    previewModelText: ' map model aliases to connected services',
    previewUsageLabel: 'billing',
    previewUsageText: ' usage-based billing with configurable quota limits',
    lightTheme: 'Light',
    darkTheme: 'Dark',
    toolLabels: ['subscription to API', 'sticky session', 'pay as you go', 'one-click access', 'reliable', 'model mapping', 'quota control', 'usage details']
  }
})

const toolIcons = computed<Array<{ icon: IconName; label: string; className: string }>>(() => [
  { icon: 'grid', label: copy.value.toolLabels[0], className: 'tool-a' },
  { icon: 'terminal', label: copy.value.toolLabels[1], className: 'tool-b' },
  { icon: 'search', label: copy.value.toolLabels[2], className: 'tool-c' },
  { icon: 'terminal', label: copy.value.toolLabels[3], className: 'tool-d' },
  { icon: 'sparkles', label: copy.value.toolLabels[4], className: 'tool-e' },
  { icon: 'sync', label: copy.value.toolLabels[5], className: 'tool-f' },
  { icon: 'document', label: copy.value.toolLabels[6], className: 'tool-g' },
  { icon: 'cpu', label: copy.value.toolLabels[7], className: 'tool-h' }
])

const products = computed<Array<{ icon: IconName; kicker: string; title: string; description: string }>>(() => isZh.value ? [
  {
    icon: 'grid',
    kicker: '一键接入',
    title: '获取一个 API 密钥，调用所有已接入模型',
    description: '使用统一 API 密钥接入平台已配置模型，减少接入和维护成本。'
  },
  {
    icon: 'terminal',
    kicker: '用多少付多少',
    title: '按实际使用量计费，团队用量一目了然',
    description: '按 Token 级别统计消耗并计算成本，支持配额上限、余额扣减和用量明细查询。'
  }
] : [
  {
    icon: 'grid',
    kicker: 'One-Click Access',
    title: 'Get one API key to call all connected models',
    description: 'Use one unified API key for configured models and reduce integration and maintenance overhead.'
  },
  {
    icon: 'terminal',
    kicker: 'Pay What You Use',
    title: 'Usage-based billing with clear team visibility',
    description: 'Track token-level consumption and cost with quota limits, balance deduction and detailed usage records.'
  }
])

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

onMounted(() => {
  initTheme()
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.ag-home {
  min-height: 100vh;
  overflow-x: hidden;
  background: #faf9f6;
  color: #111114;
  font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

.ag-home::before {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  content: '';
  opacity: 0.42;
  background:
    radial-gradient(circle at 18% 12%, rgba(66, 133, 244, 0.12), transparent 18rem),
    radial-gradient(circle at 82% 20%, rgba(234, 67, 53, 0.1), transparent 18rem),
    radial-gradient(circle at 50% 70%, rgba(52, 168, 83, 0.1), transparent 24rem),
    repeating-radial-gradient(circle at 50% 0%, rgba(20, 20, 24, 0.045) 0 1px, transparent 1px 34px);
}

.ag-header,
.ag-hero,
.ag-orbit,
.ag-product-grid,
.ag-footer {
  position: relative;
  z-index: 1;
}

.ag-header {
  z-index: 20;
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  gap: 1.25rem;
  max-width: 1160px;
  margin: 0 auto;
  padding: 0.9rem 1.5rem;
}

.ag-brand,
.ag-nav,
.ag-actions,
.ag-language,
.ag-console,
.ag-hero-actions,
.ag-primary,
.ag-secondary,
.ag-footer,
.ag-brand-fade,
.ag-product-icon,
.ag-tool,
.ag-liftoff-core {
  display: flex;
  align-items: center;
}

.ag-brand {
  gap: 0.55rem;
  color: #1f1f24;
  font-size: 1.05rem;
  font-weight: 650;
  text-decoration: none;
}

.ag-mark {
  display: inline-grid;
  width: 1.35rem;
  height: 1.35rem;
  place-items: center;
  background: conic-gradient(from 210deg, #4285f4, #34a853, #fbbc04, #ea4335, #8b5cf6, #4285f4);
  color: transparent;
  font-size: 0.1rem;
  clip-path: polygon(50% 0, 100% 100%, 72% 100%, 50% 50%, 28% 100%, 0 100%);
}

.ag-mark.small {
  width: 1.5rem;
  height: 1.5rem;
  opacity: 0.45;
}

.ag-logo {
  width: 1.8rem;
  height: 1.8rem;
  flex: 0 0 auto;
  border-radius: 0.55rem;
  object-fit: cover;
}

.ag-logo.small {
  width: 2rem;
  height: 2rem;
  opacity: 0.72;
}

.ag-logo.core {
  width: 2.6rem;
  height: 2.6rem;
  border-radius: 0.85rem;
}

.ag-nav {
  justify-content: center;
  gap: 2rem;
  color: #2f3035;
  font-size: 0.98rem;
}

.ag-nav a,
.ag-console,
.ag-primary,
.ag-secondary,
.ag-footer a {
  color: inherit;
  text-decoration: none;
}

.ag-actions {
  position: relative;
  z-index: 21;
  justify-self: end;
  justify-content: flex-end;
  gap: 0.75rem;
}

.ag-language {
  position: relative;
  z-index: 22;
  justify-content: center;
  min-height: 2.75rem;
  border: 1px solid rgba(17, 17, 20, 0.08);
  border-radius: 999px;
  padding: 0 0.35rem;
  background: rgba(255, 255, 255, 0.72);
  box-shadow: 0 0.7rem 1.8rem rgba(31, 31, 35, 0.06);
  backdrop-filter: blur(14px);
}

.ag-console {
  justify-self: end;
  justify-content: center;
  gap: 0.45rem;
  border-radius: 999px;
  padding: 0.75rem 1.1rem;
  background: #0d0e12;
  color: #fff;
  font-weight: 700;
  box-shadow: 0 0.7rem 1.8rem rgba(0, 0, 0, 0.12);
}

.ag-hero {
  position: relative;
  display: flex;
  min-height: 38rem;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  max-width: 1080px;
  margin: 0 auto;
  overflow: visible;
  padding: 6.5rem 1.5rem 4.25rem;
  text-align: center;
}

.ag-hero-orbit {
  position: absolute;
  left: 50%;
  top: 54%;
  z-index: 0;
  width: min(1040px, calc(100vw - 2rem));
  height: 31rem;
  pointer-events: none;
  transform: translate(-50%, -50%);
}

.ag-hero-orbit::before,
.ag-hero-orbit::after {
  position: absolute;
  left: 50%;
  top: 50%;
  width: 44rem;
  height: 20rem;
  border: 1px solid rgba(17, 17, 20, 0.055);
  border-radius: 50%;
  content: '';
  transform: translate(-50%, -50%) rotate(-8deg);
  animation: ag-ellipse-drift 16s linear infinite;
}

.ag-hero-orbit::after {
  width: 34rem;
  height: 15rem;
  animation-duration: 12s;
  animation-direction: reverse;
}

.ag-liftoff-core {
  position: absolute;
  left: 50%;
  top: 50%;
  z-index: 1;
  justify-content: center;
  width: 5rem;
  height: 5rem;
  border: 1px solid rgba(17, 17, 20, 0.06);
  border-radius: 1.45rem;
  background: rgba(255, 255, 255, 0.54);
  box-shadow: 0 1.5rem 3.5rem rgba(31, 31, 35, 0.1);
  transform: translate(-50%, -50%);
  backdrop-filter: blur(18px);
}

.ag-mark.core {
  width: 2.25rem;
  height: 2.25rem;
}

.ag-brand-fade,
.ag-hero h1,
.ag-hero-actions {
  position: relative;
  z-index: 2;
}

.ag-brand-fade {
  justify-content: center;
  gap: 0.7rem;
  margin-bottom: 1.7rem;
  color: rgba(17, 17, 20, 0.16);
  font-size: 1.35rem;
  font-weight: 700;
}

.ag-hero h1 {
  max-width: 1040px;
  margin: 0;
  color: #111114;
  font-size: clamp(4rem, 8.4vw, 7.7rem);
  font-weight: 620;
  letter-spacing: -0.085em;
  line-height: 0.92;
}

.ag-site-subtitle {
  position: relative;
  z-index: 2;
  max-width: 48rem;
  margin: 1.35rem 0 0;
  color: rgba(17, 17, 20, 0.56);
  font-size: clamp(1rem, 2vw, 1.35rem);
  line-height: 1.7;
}

.ag-hero-actions {
  flex-wrap: wrap;
  justify-content: center;
  gap: 1rem;
  margin-top: 6.2rem;
}

.ag-primary,
.ag-secondary {
  justify-content: center;
  gap: 0.55rem;
  border: 0;
  border-radius: 999px;
  padding: 0.95rem 1.35rem;
  font-size: 1rem;
  font-weight: 750;
}

.ag-primary {
  background: #d2d2d2;
  color: #fff;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.45);
}

.ag-secondary {
  background: rgba(255, 255, 255, 0.58);
  color: rgba(17, 17, 20, 0.24);
}

.ag-orbit {
  max-width: 1180px;
  min-height: 38rem;
  margin: 0 auto;
  padding: 1rem 1.5rem 5.5rem;
}

.ag-preview-card {
  position: relative;
  z-index: 2;
  width: min(920px, 100%);
  margin: 0 auto;
  overflow: hidden;
  border: 1px solid rgba(17, 17, 20, 0.08);
  border-radius: 2rem;
  background: rgba(255, 255, 255, 0.72);
  box-shadow: 0 2rem 5rem rgba(31, 31, 35, 0.1);
  backdrop-filter: blur(18px);
}

.ag-window-bar {
  display: flex;
  align-items: center;
  gap: 0.45rem;
  border-bottom: 1px solid rgba(17, 17, 20, 0.08);
  padding: 1rem 1.1rem;
  color: rgba(17, 17, 20, 0.48);
}

.ag-window-bar span {
  width: 0.68rem;
  height: 0.68rem;
  border-radius: 50%;
  background: #ea4335;
}

.ag-window-bar span:nth-child(2) { background: #fbbc04; }
.ag-window-bar span:nth-child(3) { background: #34a853; }
.ag-window-bar strong { margin-left: auto; font-size: 0.85rem; }

.ag-workspace {
  display: grid;
  grid-template-columns: 15rem 1fr;
  min-height: 22rem;
}

.ag-workspace aside {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  border-right: 1px solid rgba(17, 17, 20, 0.08);
  padding: 1.2rem;
  color: rgba(17, 17, 20, 0.54);
}

.ag-workspace aside strong {
  color: #111114;
  font-size: 1.15rem;
}

.ag-workspace aside small {
  color: rgba(17, 17, 20, 0.5);
  line-height: 1.55;
}

.ag-workspace aside em,
.ag-code-line {
  border-radius: 1rem;
  background: rgba(17, 17, 20, 0.045);
  font-style: normal;
}

.ag-workspace aside em {
  padding: 0.72rem 0.85rem;
}

.ag-workspace section {
  display: grid;
  align-content: center;
  gap: 1rem;
  padding: 2rem;
}

.ag-code-line {
  margin: 0;
  padding: 1rem 1.1rem;
  color: rgba(17, 17, 20, 0.68);
}

.ag-code-line b {
  display: inline-flex;
  margin-right: 0.75rem;
  color: #1a73e8;
}

.ag-tool {
  position: absolute;
  left: calc(50% - 2.8rem);
  top: calc(50% - 2.8rem);
  z-index: 1;
  justify-content: center;
  flex-direction: column;
  gap: 0.35rem;
  width: 5.6rem;
  height: 5.6rem;
  border: 1px solid rgba(17, 17, 20, 0.08);
  border-radius: 1.4rem;
  background: rgba(255, 255, 255, 0.68);
  color: rgba(17, 17, 20, 0.74);
  box-shadow: 0 1rem 2.4rem rgba(31, 31, 35, 0.08);
  backdrop-filter: blur(14px);
}

.ag-tool span {
  font-size: 0.68rem;
  font-weight: 750;
}

.tool-a { color: #4285f4; --lift-x: -25rem; --lift-y: -3.2rem; --rot: -12deg; }
.tool-b { color: #34a853; --lift-x: 25rem; --lift-y: -4.2rem; --rot: 10deg; }
.tool-c { color: #ea4335; --lift-x: -19rem; --lift-y: 10.5rem; --rot: 14deg; }
.tool-d { color: #1a73e8; --lift-x: 19rem; --lift-y: 10rem; --rot: -9deg; }
.tool-e { color: #fbbc04; --lift-x: 0rem; --lift-y: -16rem; --rot: 6deg; }
.tool-f { color: #8b5cf6; --lift-x: 8.5rem; --lift-y: 15.5rem; --rot: 12deg; }
.tool-g { color: #34a853; --lift-x: -9.5rem; --lift-y: 15.8rem; --rot: -14deg; }
.tool-h { color: #ea4335; --lift-x: 10.5rem; --lift-y: -13.8rem; --rot: -6deg; }

.ag-product-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1.2rem;
  max-width: 1120px;
  margin: 0 auto;
  padding: 0 1.5rem 6rem;
}

.ag-product-card {
  min-height: 22rem;
  border: 1px solid rgba(17, 17, 20, 0.08);
  border-radius: 2rem;
  padding: 2rem;
  background: rgba(255, 255, 255, 0.7);
  box-shadow: 0 1.4rem 3.4rem rgba(31, 31, 35, 0.08);
}

.ag-product-icon {
  justify-content: center;
  width: 3rem;
  height: 3rem;
  border-radius: 1rem;
  background: #f1f3f4;
  color: #1a73e8;
}

.ag-product-card span {
  display: block;
  margin-top: 2rem;
  color: rgba(17, 17, 20, 0.48);
  font-size: 0.85rem;
  font-weight: 800;
}

.ag-product-card h2 {
  margin: 0.7rem 0;
  font-size: clamp(1.8rem, 3vw, 3rem);
  font-weight: 620;
  letter-spacing: -0.065em;
  line-height: 1;
}

.ag-product-card p {
  color: rgba(17, 17, 20, 0.62);
  font-size: 1.02rem;
  line-height: 1.75;
}

.ag-footer {
  flex-wrap: wrap;
  justify-content: center;
  gap: 1rem;
  padding: 2rem;
  color: rgba(17, 17, 20, 0.56);
}

.ag-footer-meta {
  max-width: min(42rem, 100%);
  overflow-wrap: anywhere;
}

.ag-footer button {
  border: 1px solid rgba(17, 17, 20, 0.08);
  border-radius: 999px;
  padding: 0.45rem 0.8rem;
  background: rgba(255, 255, 255, 0.62);
}

.ag-home.ag-home-dark {
  background: #070a12;
  color: #f8fafc;
}

.ag-home.ag-home-dark::before {
  opacity: 0.76;
  background:
    radial-gradient(circle at 18% 12%, rgba(66, 133, 244, 0.28), transparent 19rem),
    radial-gradient(circle at 82% 20%, rgba(139, 92, 246, 0.22), transparent 18rem),
    radial-gradient(circle at 50% 70%, rgba(52, 168, 83, 0.16), transparent 24rem),
    repeating-radial-gradient(circle at 50% 0%, rgba(255, 255, 255, 0.055) 0 1px, transparent 1px 34px);
}

.ag-home-dark .ag-brand,
.ag-home-dark .ag-nav,
.ag-home-dark .ag-footer {
  color: rgba(248, 250, 252, 0.82);
}

.ag-home-dark .ag-language,
.ag-home-dark .ag-preview-card,
.ag-home-dark .ag-product-card,
.ag-home-dark .ag-tool,
.ag-home-dark .ag-liftoff-core {
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(15, 23, 42, 0.72);
  box-shadow: 0 1.6rem 4rem rgba(0, 0, 0, 0.34);
}

.ag-home-dark .ag-console,
.ag-home-dark .ag-primary {
  background: #f8fafc;
  color: #0f172a;
}

.ag-home-dark .ag-secondary,
.ag-home-dark .ag-footer button {
  border-color: rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.08);
  color: rgba(248, 250, 252, 0.8);
}

.ag-home-dark .ag-brand-fade {
  color: rgba(248, 250, 252, 0.28);
}

.ag-home-dark .ag-hero h1,
.ag-home-dark .ag-workspace aside strong {
  color: #f8fafc;
}

.ag-home-dark .ag-site-subtitle,
.ag-home-dark .ag-workspace aside small,
.ag-home-dark .ag-footer-meta {
  color: rgba(248, 250, 252, 0.62);
}

.ag-home-dark .ag-hero-orbit::before,
.ag-home-dark .ag-hero-orbit::after {
  border-color: rgba(248, 250, 252, 0.1);
}

.ag-home-dark .ag-window-bar,
.ag-home-dark .ag-workspace aside {
  border-color: rgba(255, 255, 255, 0.1);
  color: rgba(248, 250, 252, 0.58);
}

.ag-home-dark .ag-workspace aside em,
.ag-home-dark .ag-code-line {
  background: rgba(255, 255, 255, 0.06);
}

.ag-home-dark .ag-code-line {
  color: rgba(248, 250, 252, 0.72);
}

.ag-home-dark .ag-code-line b {
  color: #8ab4f8;
}

.ag-home-dark .ag-product-icon {
  background: rgba(26, 115, 232, 0.16);
  color: #8ab4f8;
}

.ag-home-dark .ag-product-card span {
  color: rgba(248, 250, 252, 0.52);
}

.ag-home-dark .ag-product-card p {
  color: rgba(248, 250, 252, 0.68);
}

.ag-home-dark .ag-orbit::before,
.ag-home-dark .ag-orbit::after {
  opacity: 0.3;
}



/* Antigravity-style liftoff motion */
.ag-brand-fade,
.ag-hero h1,
.ag-hero-actions,
.ag-preview-card,
.ag-product-card {
  animation: ag-rise 0.9s cubic-bezier(0.2, 0.8, 0.2, 1) both;
}

.ag-brand-fade { animation-delay: 0.05s; }
.ag-hero h1 { animation-delay: 0.16s; }
.ag-hero-actions { animation-delay: 0.32s; }
.ag-preview-card { animation-delay: 0.5s; }
.ag-product-card:nth-child(1) { animation-delay: 0.74s; }
.ag-product-card:nth-child(2) { animation-delay: 0.82s; }
.ag-product-card:nth-child(3) { animation-delay: 0.9s; }
.ag-product-card:nth-child(4) { animation-delay: 0.98s; }

.ag-liftoff-core {
  animation:
    ag-core-pop 0.7s cubic-bezier(0.2, 0.9, 0.15, 1) 0.06s both,
    ag-core-pulse 2.8s ease-in-out 0.8s infinite;
}

.ag-liftoff-core::before,
.ag-liftoff-core::after {
  position: absolute;
  inset: -1.2rem;
  border: 1px solid rgba(66, 133, 244, 0.18);
  border-radius: 2.1rem;
  content: '';
  opacity: 0;
  animation: ag-core-ring 2.6s ease-out 0.5s infinite;
}

.ag-liftoff-core::after {
  animation-delay: 1.05s;
}

.ag-orbit::before,
.ag-orbit::after {
  position: absolute;
  left: 50%;
  top: 34%;
  width: 44rem;
  height: 44rem;
  border: 1px solid rgba(17, 17, 20, 0.055);
  border-radius: 999px;
  content: '';
  transform: translate(-50%, -50%);
  animation: ag-orbit-spin 18s linear infinite;
}

.ag-orbit::after {
  width: 30rem;
  height: 30rem;
  animation-duration: 13s;
  animation-direction: reverse;
}

.ag-tool {
  opacity: 0;
  transform: translate(0, 0) scale(0.2);
  animation:
    ag-tool-liftoff 1.25s cubic-bezier(0.16, 1, 0.3, 1) forwards,
    ag-tool-breathe 5.2s ease-in-out infinite;
}

.tool-a { animation-delay: 0.14s, 1.3s; --float-y: -0.7rem; }
.tool-b { animation-delay: 0.2s, 1.38s; --float-y: -0.95rem; }
.tool-c { animation-delay: 0.26s, 1.46s; --float-y: -0.8rem; }
.tool-d { animation-delay: 0.32s, 1.54s; --float-y: -1.05rem; }
.tool-e { animation-delay: 0.38s, 1.62s; --float-y: -0.75rem; }
.tool-f { animation-delay: 0.44s, 1.7s; --float-y: -0.9rem; }
.tool-g { animation-delay: 0.5s, 1.78s; --float-y: -0.85rem; }
.tool-h { animation-delay: 0.56s, 1.86s; --float-y: -1rem; }

.ag-float-window {
  animation:
    ag-rise 0.9s cubic-bezier(0.2, 0.8, 0.2, 1) 0.5s both,
    ag-window-hover 5.8s ease-in-out 1.45s infinite;
}

.ag-primary,
.ag-secondary,
.ag-console {
  transition: transform 0.18s ease, box-shadow 0.18s ease, background 0.18s ease;
}

.ag-primary:hover,
.ag-secondary:hover,
.ag-console:hover {
  transform: translateY(-2px);
}

.ag-primary:hover {
  background: #c4c4c4;
}

.ag-console:hover {
  box-shadow: 0 1rem 2.3rem rgba(0, 0, 0, 0.18);
}

@keyframes ag-rise {
  from {
    opacity: 0;
    transform: translateY(2.2rem) scale(0.98);
    filter: blur(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
    filter: blur(0);
  }
}

@keyframes ag-tool-liftoff {
  0% {
    opacity: 0;
    transform: translate(calc(var(--lift-x) * -0.06), calc(var(--lift-y) * -0.06)) scale(0.2) rotate(0deg);
    filter: blur(8px);
  }
  55% {
    opacity: 1;
    transform: translate(calc(var(--lift-x) * 1.08), calc(var(--lift-y) * 1.08 - 0.75rem)) scale(1.05) rotate(calc(var(--rot) * 1.2));
    filter: blur(0);
  }
  78% {
    transform: translate(calc(var(--lift-x) * 0.96), calc(var(--lift-y) * 0.96 + 0.2rem)) scale(0.98) rotate(calc(var(--rot) * 0.82));
  }
  100% {
    opacity: 1;
    transform: translate(var(--lift-x), var(--lift-y)) scale(1) rotate(var(--rot));
    filter: blur(0);
  }
}

@keyframes ag-tool-breathe {
  0%, 100% {
    translate: 0 0;
  }
  50% {
    translate: 0 var(--float-y);
  }
}

@keyframes ag-core-pop {
  from {
    opacity: 0;
    transform: translate(-50%, -50%) scale(0.32) rotate(-10deg);
    filter: blur(10px);
  }
  to {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1) rotate(0deg);
    filter: blur(0);
  }
}

@keyframes ag-core-pulse {
  0%, 100% {
    box-shadow: 0 1.5rem 3.5rem rgba(31, 31, 35, 0.1);
  }
  50% {
    box-shadow: 0 1.8rem 4.4rem rgba(66, 133, 244, 0.18);
  }
}

@keyframes ag-core-ring {
  0% {
    opacity: 0.55;
    transform: scale(0.6);
  }
  100% {
    opacity: 0;
    transform: scale(1.75);
  }
}

@keyframes ag-window-hover {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-0.75rem); }
}

@keyframes ag-orbit-spin {
  from { transform: translate(-50%, -50%) rotate(0deg); }
  to { transform: translate(-50%, -50%) rotate(360deg); }
}

@keyframes ag-ellipse-drift {
  from { transform: translate(-50%, -50%) rotate(-8deg); }
  to { transform: translate(-50%, -50%) rotate(352deg); }
}

@media (prefers-reduced-motion: reduce) {
  .ag-brand-fade,
  .ag-hero h1,
  .ag-hero-actions,
  .ag-preview-card,
  .ag-product-card,
  .ag-tool,
  .ag-liftoff-core,
  .ag-liftoff-core::before,
  .ag-liftoff-core::after,
  .ag-float-window,
  .ag-hero-orbit::before,
  .ag-hero-orbit::after,
  .ag-orbit::before,
  .ag-orbit::after {
    animation: none;
    opacity: 1;
  }

  .ag-tool {
    transform: translate(var(--lift-x), var(--lift-y)) scale(1) rotate(var(--rot));
  }
}

@media (max-width: 960px) {
  .ag-header {
    grid-template-columns: 1fr auto;
  }

  .ag-nav {
    display: none;
  }

  .ag-hero {
    min-height: 31rem;
    padding-top: 4rem;
  }

  .ag-hero-actions {
    margin-top: 3rem;
  }

  .ag-tool {
    display: none;
  }

  .ag-liftoff-core,
  .ag-hero-orbit {
    display: none;
  }

  .ag-workspace {
    grid-template-columns: 1fr;
  }

  .ag-workspace aside {
    border-right: 0;
    border-bottom: 1px solid rgba(17, 17, 20, 0.08);
  }
}

@media (max-width: 720px) {
  .ag-header {
    grid-template-columns: 1fr;
    padding-right: 1rem;
    padding-left: 1rem;
  }

  .ag-brand span:last-child {
    max-width: 10rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .ag-actions {
    width: 100%;
    justify-self: stretch;
    justify-content: space-between;
    gap: 0.5rem;
  }

  .ag-language :deep(.absolute.right-0) {
    right: auto;
    left: 0;
  }

  .ag-console {
    padding: 0.65rem 0.85rem;
  }

  .ag-hero h1 {
    font-size: clamp(3rem, 14vw, 5rem);
  }

  .ag-product-grid {
    grid-template-columns: 1fr;
  }
}
</style>
