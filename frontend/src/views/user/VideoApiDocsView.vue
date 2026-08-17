<template>
  <AppLayout>
    <div class="mx-auto max-w-[1320px] space-y-5">
      <header class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-sm font-medium text-primary-600 dark:text-primary-400">{{ t(copy.eyebrow) }}</p>
          <h1 class="mt-1 text-2xl font-semibold tracking-tight text-gray-900 dark:text-white">{{ t(copy.title) }}</h1>
          <p class="mt-1 max-w-3xl text-sm text-gray-500 dark:text-dark-300">{{ t(copy.description) }}</p>
        </div>
        <div class="flex flex-col items-start gap-2 sm:items-end">
          <div class="inline-flex overflow-hidden rounded-md border border-gray-200 text-xs font-medium dark:border-dark-700" :aria-label="t('video.apiDocs.version.label')">
            <RouterLink
              v-for="option in versionOptions"
              :key="option.name"
              :to="{ name: option.name }"
              class="px-3 py-1.5 transition-colors"
              :class="route.name === option.name
                ? 'bg-primary-600 text-white dark:bg-primary-500'
                : 'bg-white text-gray-600 hover:bg-gray-50 dark:bg-dark-900 dark:text-dark-300 dark:hover:bg-dark-800'"
            >
              {{ t(option.label) }}
            </RouterLink>
          </div>
          <div class="inline-flex h-8 flex-shrink-0 items-center gap-2 rounded-md border border-gray-200 px-2.5 font-mono text-xs text-gray-600 dark:border-dark-700 dark:text-dark-300">
            <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
            {{ t(copy.badge) }}
          </div>
        </div>
      </header>

      <VideoSectionTabs />
      <VideoApiDocsV2Content v-if="isV2" />
      <VideoApiDocsV1Content v-else />
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, useRoute } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import VideoApiDocsV1Content from '@/components/video/VideoApiDocsV1Content.vue'
import VideoApiDocsV2Content from '@/components/video/VideoApiDocsV2Content.vue'
import VideoSectionTabs from '@/components/video/VideoSectionTabs.vue'

const { t } = useI18n()
const route = useRoute()
const isV2 = computed(() => route.name !== 'VideoApiDocsV1')

const versionOptions = [
  { name: 'VideoApiDocs', label: 'video.apiDocs.version.v2' },
  { name: 'VideoApiDocsV1', label: 'video.apiDocs.version.v1' },
]

const copy = computed(() => isV2.value
  ? {
      eyebrow: 'video.apiDocs.v2.eyebrow',
      title: 'video.apiDocs.v2.title',
      description: 'video.apiDocs.v2.description',
      badge: 'video.apiDocs.v2.badge',
    }
  : {
      eyebrow: 'video.apiDocs.eyebrow',
      title: 'video.apiDocs.title',
      description: 'video.apiDocs.description',
      badge: 'video.apiDocs.badge',
    })
</script>
