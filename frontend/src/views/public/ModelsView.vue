<template>
  <div class="models-page">
    <div class="models-bg" aria-hidden="true">
      <span class="models-orb models-orb-a"></span>
      <span class="models-orb models-orb-b"></span>
    </div>

    <header class="models-header">
      <router-link to="/home" class="models-brand">
        <span class="models-logo-wrap">
          <img :src="siteLogo || '/logo.png'" alt="Logo" class="models-logo" />
        </span>
        <span>{{ siteName }}</span>
      </router-link>
      <router-link to="/home" class="models-home-link">返回首页</router-link>
    </header>

    <main class="models-main">
      <section class="models-hero">
        <p>Model Catalog</p>
        <h1>模型列表</h1>
        <span>展示当前前端配置内置的模型集合，便于用户在接入前快速确认模型命名和平台覆盖。</span>
      </section>

      <section class="models-toolbar" aria-label="Model filters">
        <label class="models-search">
          <Icon name="search" size="sm" />
          <input v-model="searchText" type="search" placeholder="搜索模型名称" />
        </label>
        <div class="models-tabs">
          <button
            v-for="platform in platforms"
            :key="platform.key"
            type="button"
            :class="{ active: activePlatform === platform.key }"
            @click="activePlatform = platform.key"
          >
            {{ platform.label }}
            <span>{{ platform.models.length }}</span>
          </button>
        </div>
      </section>

      <section class="models-summary" aria-label="Model summary">
        <article>
          <strong>{{ totalCount }}</strong>
          <span>全部模型</span>
        </article>
        <article>
          <strong>{{ platforms.length - 1 }}</strong>
          <span>平台分组</span>
        </article>
        <article>
          <strong>{{ filteredModels.length }}</strong>
          <span>当前结果</span>
        </article>
      </section>

      <section class="models-grid" aria-live="polite">
        <article v-for="model in filteredModels" :key="`${model.platform}-${model.name}`" class="model-card">
          <div>
            <span class="model-provider">{{ model.platformLabel }}</span>
            <h2>{{ model.name }}</h2>
          </div>
          <span class="model-chip">{{ model.family }}</span>
        </article>
      </section>

      <p v-if="filteredModels.length === 0" class="models-empty">没有匹配的模型。</p>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useAppStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import { getModelsByPlatform } from '@/composables/useModelWhitelist'

const appStore = useAppStore()
const searchText = ref('')
const activePlatform = ref('all')

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')

const platformDefinitions = [
  { key: 'openai', label: 'OpenAI' },
  { key: 'claude', label: 'Claude' },
  { key: 'gemini', label: 'Gemini' },
  { key: 'antigravity', label: 'Antigravity' },
  { key: 'qwen', label: 'Qwen' },
  { key: 'deepseek', label: 'DeepSeek' },
  { key: 'xai', label: 'Grok' },
  { key: 'doubao', label: 'Doubao' },
  { key: 'zhipu', label: 'GLM' },
  { key: 'mistral', label: 'Mistral' },
  { key: 'meta', label: 'Llama' },
  { key: 'moonshot', label: 'Kimi' },
  { key: 'perplexity', label: 'Perplexity' }
]

const platformGroups = platformDefinitions.map((platform) => ({
  ...platform,
  models: getModelsByPlatform(platform.key)
}))

const platforms = [
  {
    key: 'all',
    label: '全部',
    models: platformGroups.flatMap((platform) => platform.models)
  },
  ...platformGroups
]

const totalCount = computed(() => platforms[0].models.length)

const filteredModels = computed(() => {
  const keyword = searchText.value.trim().toLowerCase()
  const selectedPlatforms = activePlatform.value === 'all'
    ? platformGroups
    : platformGroups.filter((platform) => platform.key === activePlatform.value)

  return selectedPlatforms
    .flatMap((platform) => platform.models.map((name) => ({
      name,
      platform: platform.key,
      platformLabel: platform.label,
      family: resolveModelFamily(name)
    })))
    .filter((model) => !keyword || model.name.toLowerCase().includes(keyword))
})

function resolveModelFamily(model: string): string {
  if (model.includes('image') || model.includes('vision') || model.includes('cogview')) return 'multimodal'
  if (model.includes('coder') || model.includes('code') || model.includes('codex')) return 'coding'
  if (model.includes('reasoner') || model.includes('thinking') || model.includes('r1') || model.includes('qwq')) return 'reasoning'
  if (model.includes('flash') || model.includes('lite') || model.includes('mini') || model.includes('haiku')) return 'fast'
  return 'chat'
}

onMounted(() => {
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.models-page {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
  background:
    radial-gradient(circle at 20% 0%, rgba(20, 184, 166, 0.18), transparent 34rem),
    linear-gradient(135deg, #030712, #0f172a 48%, #111827);
  color: #f8fafc;
}

.models-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background-image:
    linear-gradient(rgba(148, 163, 184, 0.12) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148, 163, 184, 0.1) 1px, transparent 1px);
  background-size: 76px 76px;
  mask-image: linear-gradient(to bottom, black, transparent 70%);
}

.models-orb {
  position: absolute;
  width: 30rem;
  height: 30rem;
  border-radius: 999px;
  filter: blur(80px);
}

.models-orb-a {
  right: -8rem;
  top: -8rem;
  background: rgba(244, 63, 94, 0.48);
}

.models-orb-b {
  left: -12rem;
  bottom: 10rem;
  background: rgba(59, 130, 246, 0.32);
}

.models-header,
.models-main {
  position: relative;
  z-index: 1;
  max-width: 1180px;
  margin: 0 auto;
  padding-right: 1.5rem;
  padding-left: 1.5rem;
}

.models-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 1.25rem;
  padding-bottom: 1.25rem;
}

.models-brand,
.models-home-link {
  display: inline-flex;
  align-items: center;
  color: #fff;
  text-decoration: none;
}

.models-brand {
  gap: 0.75rem;
  font-weight: 800;
}

.models-logo-wrap {
  display: grid;
  width: 2.5rem;
  height: 2.5rem;
  place-items: center;
  overflow: hidden;
  border-radius: 0.9rem;
  background: rgba(255, 255, 255, 0.1);
}

.models-logo {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.models-home-link {
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 999px;
  padding: 0.65rem 1rem;
  background: rgba(255, 255, 255, 0.07);
  font-size: 0.875rem;
  font-weight: 700;
}

.models-main {
  padding-bottom: 5rem;
}

.models-hero {
  max-width: 54rem;
  padding: 4rem 0 2rem;
}

.models-hero p {
  margin: 0;
  color: #67e8f9;
  font-size: 0.82rem;
  font-weight: 900;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.models-hero h1 {
  margin: 0.8rem 0;
  font-size: clamp(3rem, 7vw, 6rem);
  font-weight: 900;
  letter-spacing: -0.08em;
  line-height: 0.9;
}

.models-hero span {
  display: block;
  max-width: 42rem;
  color: rgba(226, 232, 240, 0.72);
  font-size: 1.05rem;
  line-height: 1.8;
}

.models-toolbar,
.models-summary,
.model-card {
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(15, 23, 42, 0.62);
  box-shadow: 0 24px 70px rgba(0, 0, 0, 0.16);
  backdrop-filter: blur(20px);
}

.models-toolbar {
  display: grid;
  gap: 1rem;
  border-radius: 1.5rem;
  padding: 1rem;
}

.models-search {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 999px;
  padding: 0.75rem 1rem;
  background: rgba(255, 255, 255, 0.07);
  color: rgba(226, 232, 240, 0.76);
}

.models-search input {
  width: 100%;
  border: 0;
  outline: 0;
  background: transparent;
  color: #fff;
}

.models-search input::placeholder {
  color: rgba(226, 232, 240, 0.42);
}

.models-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 0.65rem;
}

.models-tabs button {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 999px;
  padding: 0.55rem 0.85rem;
  background: rgba(255, 255, 255, 0.06);
  color: rgba(226, 232, 240, 0.78);
  font-size: 0.82rem;
  font-weight: 800;
}

.models-tabs button.active {
  border-color: rgba(45, 212, 191, 0.42);
  background: rgba(20, 184, 166, 0.2);
  color: #ccfbf1;
}

.models-tabs span {
  color: rgba(226, 232, 240, 0.5);
}

.models-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1px;
  overflow: hidden;
  margin: 1rem 0;
  border-radius: 1.5rem;
}

.models-summary article {
  padding: 1.25rem;
  background: rgba(255, 255, 255, 0.04);
}

.models-summary strong,
.models-summary span {
  display: block;
}

.models-summary strong {
  font-size: 1.8rem;
}

.models-summary span {
  color: rgba(226, 232, 240, 0.6);
  font-size: 0.82rem;
}

.models-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
}

.model-card {
  display: flex;
  min-height: 9.5rem;
  flex-direction: column;
  justify-content: space-between;
  border-radius: 1.4rem;
  padding: 1.15rem;
  transition: transform 0.2s ease, border-color 0.2s ease;
}

.model-card:hover {
  transform: translateY(-3px);
  border-color: rgba(45, 212, 191, 0.38);
}

.model-provider {
  color: #67e8f9;
  font-size: 0.75rem;
  font-weight: 900;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.model-card h2 {
  overflow-wrap: anywhere;
  margin: 0.65rem 0 0;
  font-size: 1.05rem;
  line-height: 1.4;
}

.model-chip {
  align-self: flex-start;
  border-radius: 999px;
  padding: 0.35rem 0.65rem;
  background: rgba(255, 255, 255, 0.08);
  color: rgba(226, 232, 240, 0.76);
  font-size: 0.75rem;
  font-weight: 800;
}

.models-empty {
  padding: 4rem 0;
  color: rgba(226, 232, 240, 0.66);
  text-align: center;
}

@media (max-width: 980px) {
  .models-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .models-header,
  .models-main {
    padding-right: 1rem;
    padding-left: 1rem;
  }

  .models-summary,
  .models-grid {
    grid-template-columns: 1fr;
  }
}
</style>
