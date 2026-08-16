<template>
  <AppLayout>
    <div class="mx-auto max-w-[1320px] space-y-5">
      <header class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-sm font-medium text-primary-600 dark:text-primary-400">{{ t('video.pricing.eyebrow') }}</p>
          <h1 class="mt-1 text-2xl font-semibold tracking-tight text-gray-900 dark:text-white">{{ t('video.pricing.title') }}</h1>
          <p class="mt-1 max-w-3xl text-sm text-gray-500 dark:text-dark-300">{{ t('video.pricing.description') }}</p>
        </div>
        <button
          type="button"
          class="btn btn-secondary inline-flex h-9 items-center gap-2 self-start px-3 sm:self-auto"
          :disabled="loading"
          :title="t('video.pricing.refresh')"
          @click="loadPricing"
        >
          <Icon name="refresh" size="sm" :class="loading ? 'animate-spin' : ''" />
          <span>{{ t('video.pricing.refresh') }}</span>
        </button>
      </header>

      <VideoSectionTabs />

      <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('video.pricing.source') }}</p>

      <div v-if="loading && sections.length === 0" class="rounded-lg border border-dashed border-gray-300 px-5 py-8 text-center text-sm text-gray-500 dark:border-dark-600 dark:text-dark-400">
        {{ t('video.pricing.loading') }}
      </div>
      <div v-else-if="unavailable" class="rounded-lg border border-red-200 bg-red-50 px-5 py-8 text-center text-sm text-red-600 dark:border-red-500/30 dark:bg-red-500/10 dark:text-red-300">
        {{ t('video.pricing.unavailable') }}
      </div>
      <div v-else-if="sections.length === 0" class="rounded-lg border border-dashed border-gray-300 px-5 py-8 text-center text-sm text-gray-500 dark:border-dark-600 dark:text-dark-400">
        {{ t('video.pricing.empty') }}
      </div>
      <div v-else class="space-y-8">
        <section v-for="section in sections" :key="section.id" class="space-y-4">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ section.name }}</h2>
          <div class="grid gap-4 md:grid-cols-2">
            <article
              v-for="card in section.cards"
              :key="card.key"
              class="rounded-lg border border-gray-200 p-4 dark:border-dark-700"
              :data-testid="`price-card-${card.key}`"
            >
              <div class="mb-3 flex items-center justify-between gap-2">
                <span
                  class="rounded-full px-2 py-0.5 text-xs font-medium"
                  :class="card.kind === 'image'
                    ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-300'
                    : 'bg-sky-50 text-sky-700 dark:bg-sky-500/10 dark:text-sky-300'"
                >
                  {{ card.kind === 'image' ? t('video.pricing.image') : t('video.pricing.video') }}
                </span>
                <span class="text-xs text-gray-400 dark:text-dark-400">
                  {{ card.unit === 'per_image' ? t('video.pricing.perImage') : t('video.pricing.perSecond') }}
                </span>
              </div>
              <h3 class="font-mono text-sm font-semibold text-gray-900 dark:text-white">{{ card.title }}</h3>
              <div v-if="card.tiers.length === 1 && card.tiers[0].label == null" class="mt-3 text-2xl font-extrabold text-gray-900 dark:text-white">
                ${{ formatUsd(card.tiers[0].value) }}
              </div>
              <div v-else class="mt-3 grid grid-cols-3 gap-3">
                <div v-for="tier in card.tiers" :key="tier.label ?? 'price'">
                  <div class="text-xs text-gray-500 dark:text-dark-400">{{ tier.label }}</div>
                  <div class="text-lg font-extrabold text-gray-900 dark:text-white">
                    ${{ formatUsd(tier.value) }}<span v-if="card.unit === 'per_second'" class="text-xs font-medium text-gray-500">/s</span>
                  </div>
                </div>
              </div>
            </article>
          </div>
        </section>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import VideoSectionTabs from '@/components/video/VideoSectionTabs.vue'
import userGroupsAPI from '@/api/groups'
import userChannelsAPI, { type UserAvailableChannel } from '@/api/channels'
import modelPlazaAPI, { type ModelPlazaGroup } from '@/api/modelPlaza'
import type { Group } from '@/types'
import { buildMediaPricingSections } from '@/utils/mediaPricing'

const { t } = useI18n()
const groups = ref<Group[]>([])
const plazaGroups = ref<ModelPlazaGroup[]>([])
const channels = ref<UserAvailableChannel[]>([])
const loading = ref(false)
const groupsFailed = ref(false)

const sections = computed(() => buildMediaPricingSections(groups.value, plazaGroups.value, channels.value))
const unavailable = computed(() => groupsFailed.value && sections.value.length === 0)

function formatUsd(value: number): string {
  return String(Math.round(value * 1_000_000) / 1_000_000)
}

async function loadPricing() {
  loading.value = true
  groupsFailed.value = false
  const [groupsResult, plazaResult, channelsResult] = await Promise.allSettled([
    userGroupsAPI.getAvailable(),
    modelPlazaAPI.getModelPlaza(),
    userChannelsAPI.getAvailable(),
  ])
  if (groupsResult.status === 'fulfilled') {
    groups.value = groupsResult.value
  } else {
    groupsFailed.value = true
  }
  if (plazaResult.status === 'fulfilled') {
    plazaGroups.value = plazaResult.value.groups
  }
  if (channelsResult.status === 'fulfilled') {
    channels.value = channelsResult.value
  }
  loading.value = false
}

onMounted(() => {
  void loadPricing()
})
</script>
