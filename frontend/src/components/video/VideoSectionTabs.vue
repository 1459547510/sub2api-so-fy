<template>
  <nav class="border-b border-gray-200 dark:border-dark-700" :aria-label="t('video.tabs.label')">
    <div class="flex min-w-0 gap-6 overflow-x-auto" role="tablist">
      <RouterLink
        v-for="tab in tabs"
        :key="tab.name"
        :to="tab.to"
        role="tab"
        :aria-selected="isTabActive(tab)"
        class="inline-flex h-11 flex-shrink-0 items-center gap-2 border-b-2 px-1 text-sm font-medium transition-colors"
        :class="isTabActive(tab)
          ? 'border-primary-600 text-primary-700 dark:border-primary-400 dark:text-primary-300'
          : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-800 dark:text-dark-300 dark:hover:border-dark-500 dark:hover:text-white'"
      >
        <Icon :name="tab.icon" size="sm" />
        <span>{{ t(tab.label) }}</span>
      </RouterLink>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { RouterLink, useRoute } from 'vue-router'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()
const route = useRoute()
const tabs = [
  { name: 'VideoGeneration', to: '/video-generation', label: 'video.tabs.workbench', icon: 'play' as const },
  { name: 'VideoApiDocs', to: '/video-generation/api-docs', label: 'video.tabs.apiDocs', icon: 'book' as const },
  { name: 'VideoPricing', to: '/video-generation/pricing', label: 'video.tabs.pricing', icon: 'dollar' as const },
]

function isTabActive(tab: (typeof tabs)[number]) {
  if (tab.name === 'VideoApiDocs') {
    return route.name === 'VideoApiDocs' || route.name === 'VideoApiDocsV1'
  }
  return route.name === tab.name
}
</script>
