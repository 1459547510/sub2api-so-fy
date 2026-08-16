<template>
  <AppLayout>
    <div class="mx-auto max-w-[1320px] space-y-5">
      <header class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-sm font-medium text-primary-600 dark:text-primary-400">{{ t('video.apiDocs.v2.eyebrow') }}</p>
          <h1 class="mt-1 text-2xl font-semibold tracking-tight text-gray-900 dark:text-white">{{ t('video.apiDocs.v2.title') }}</h1>
          <p class="mt-1 max-w-3xl text-sm text-gray-500 dark:text-dark-300">{{ t('video.apiDocs.v2.description') }}</p>
        </div>
        <div class="inline-flex h-8 flex-shrink-0 items-center gap-2 self-start rounded-md border border-gray-200 px-2.5 font-mono text-xs text-gray-600 dark:border-dark-700 dark:text-dark-300 sm:self-auto">
          <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
          {{ t('video.apiDocs.v2.badge') }}
        </div>
      </header>

      <VideoSectionTabs />

      <div class="grid min-w-0 gap-8 lg:grid-cols-[210px_minmax(0,1fr)] lg:gap-12">
        <aside class="hidden lg:block">
          <nav class="sticky top-6 border-l border-gray-200 pl-4 dark:border-dark-700" :aria-label="t('video.apiDocs.v2.tocLabel')">
            <a
              v-for="item in navigation"
              :key="item.href"
              :href="item.href"
              class="block py-2 text-sm text-gray-500 transition-colors hover:text-primary-700 dark:text-dark-300 dark:hover:text-primary-300"
            >
              {{ t(item.label) }}
            </a>
          </nav>
        </aside>

        <main class="min-w-0 space-y-12 pb-12">
          <section id="quick-start" class="scroll-mt-6">
            <SectionHeading :title="t('video.apiDocs.v2.quick.title')" :description="t('video.apiDocs.v2.quick.description')" />
            <dl class="mt-5 grid gap-px overflow-hidden rounded-lg border border-gray-200 bg-gray-200 sm:grid-cols-3 dark:border-dark-700 dark:bg-dark-700">
              <div class="min-w-0 bg-white p-4 dark:bg-dark-900">
                <dt class="text-xs font-medium uppercase text-gray-400 dark:text-dark-400">{{ t('video.apiDocs.v2.quick.baseUrl') }}</dt>
                <dd class="mt-2 break-all font-mono text-sm text-gray-800 dark:text-gray-100">{{ baseUrl }}</dd>
              </div>
              <div class="min-w-0 bg-white p-4 dark:bg-dark-900">
                <dt class="text-xs font-medium uppercase text-gray-400 dark:text-dark-400">{{ t('video.apiDocs.v2.quick.auth') }}</dt>
                <dd class="mt-2 font-mono text-sm text-gray-800 dark:text-gray-100">Bearer API Key</dd>
              </div>
              <div class="min-w-0 bg-white p-4 dark:bg-dark-900">
                <dt class="text-xs font-medium uppercase text-gray-400 dark:text-dark-400">{{ t('video.apiDocs.v2.quick.mode') }}</dt>
                <dd class="mt-2 font-mono text-sm text-gray-800 dark:text-gray-100">Prefer: respond-async</dd>
              </div>
            </dl>
            <div class="mt-5">
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.videoRequest')" :code="textRequestExample" />
            </div>
            <div class="mt-4">
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.acceptedResponse')" :code="acceptedResponseExample" />
            </div>
          </section>

          <section id="compatibility" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.v2.compatibility.title')" :description="t('video.apiDocs.v2.compatibility.description')" />
            <div class="mt-5 grid gap-4 md:grid-cols-2">
              <div v-for="item in compatibilityItems" :key="item.title" class="rounded-lg border border-gray-200 p-4 dark:border-dark-700">
                <h3 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t(item.title) }}</h3>
                <p class="mt-2 text-sm leading-6 text-gray-600 dark:text-dark-300">{{ t(item.description) }}</p>
              </div>
            </div>
          </section>

          <section id="model-discovery" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.v2.discovery.title')" :description="t('video.apiDocs.v2.discovery.description')" />
            <EndpointTitle method="GET" path="/v1/models" class="mt-5" />
            <p class="mt-3 text-sm leading-6 text-gray-600 dark:text-dark-300">{{ t('video.apiDocs.v2.discovery.details') }}</p>
            <div class="mt-4">
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.modelList')" :code="modelListExample" />
            </div>
          </section>

          <section id="pricing" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <div class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
              <SectionHeading :title="t('video.apiDocs.v2.pricing.title')" :description="t('video.apiDocs.v2.pricing.description')" />
              <button
                type="button"
                class="btn btn-secondary inline-flex h-9 items-center gap-2 self-start px-3 sm:self-auto"
                :disabled="pricingLoading"
                :title="t('video.apiDocs.v2.pricing.refresh')"
                @click="loadPricing"
              >
                <Icon name="refresh" size="sm" :class="pricingLoading ? 'animate-spin' : ''" />
                <span>{{ t('video.apiDocs.v2.pricing.refresh') }}</span>
              </button>
            </div>
            <div class="mt-4 flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-gray-500 dark:text-dark-400">
              <span>{{ t('video.apiDocs.v2.pricing.source') }}</span>
              <span v-if="pricingUpdatedAt">{{ t('video.apiDocs.v2.pricing.updatedAt', { time: formatUpdatedAt(pricingUpdatedAt) }) }}</span>
            </div>
            <div v-if="pricingLoading && pricingRows.length === 0" class="mt-5 rounded-lg border border-dashed border-gray-300 px-5 py-8 text-center text-sm text-gray-500 dark:border-dark-600 dark:text-dark-400">
              {{ t('video.apiDocs.v2.pricing.loading') }}
            </div>
            <div v-else-if="pricingError && pricingRows.length === 0" class="mt-5 rounded-lg border border-red-200 bg-red-50 px-5 py-8 text-center text-sm text-red-600 dark:border-red-500/30 dark:bg-red-500/10 dark:text-red-300">
              {{ t('video.apiDocs.v2.pricing.unavailable') }}
            </div>
            <div v-else-if="pricingRows.length === 0" class="mt-5 rounded-lg border border-dashed border-gray-300 px-5 py-8 text-center text-sm text-gray-500 dark:border-dark-600 dark:text-dark-400">
              {{ t('video.apiDocs.v2.pricing.empty') }}
            </div>
            <div v-else class="mt-5 overflow-x-auto rounded-lg border border-gray-200 dark:border-dark-700">
              <table class="min-w-[720px] w-full text-left text-sm">
                <thead class="bg-gray-50 text-xs uppercase text-gray-500 dark:bg-dark-800 dark:text-dark-300">
                  <tr>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.pricing.group') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.pricing.model') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.pricing.type') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.pricing.price') }}</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100 bg-white dark:divide-dark-700 dark:bg-dark-900">
                  <tr v-for="row in pricingRows" :key="row.key">
                    <td class="px-4 py-3 text-gray-700 dark:text-dark-200">{{ row.group }}</td>
                    <td class="px-4 py-3 font-mono text-xs text-gray-900 dark:text-white">{{ row.model }}</td>
                    <td class="px-4 py-3 text-gray-600 dark:text-dark-300">{{ row.type }}</td>
                    <td class="px-4 py-3 leading-6 text-gray-600 dark:text-dark-300">{{ row.price }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </section>

          <section id="image-api" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.v2.images.title')" :description="t('video.apiDocs.v2.images.description')" />
            <div class="mt-5 grid gap-4 lg:grid-cols-2">
              <div>
                <EndpointTitle method="POST" path="/v1/images/generations" />
                <p class="mt-3 text-sm leading-6 text-gray-600 dark:text-dark-300">{{ t('video.apiDocs.v2.images.generations') }}</p>
              </div>
              <div>
                <EndpointTitle method="POST" path="/v1/images/edits" />
                <p class="mt-3 text-sm leading-6 text-gray-600 dark:text-dark-300">{{ t('video.apiDocs.v2.images.edits') }}</p>
              </div>
            </div>
            <div class="mt-5 overflow-x-auto rounded-lg border border-gray-200 dark:border-dark-700">
              <table class="min-w-[720px] w-full text-left text-sm">
                <thead class="bg-gray-50 text-xs uppercase text-gray-500 dark:bg-dark-800 dark:text-dark-300">
                  <tr>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.fields.field') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.fields.type') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.fields.required') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.fields.description') }}</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100 bg-white dark:divide-dark-700 dark:bg-dark-900">
                  <tr v-for="field in imageFields" :key="field.name">
                    <td class="whitespace-nowrap px-4 py-3 font-mono text-xs text-gray-900 dark:text-white">{{ field.name }}</td>
                    <td class="whitespace-nowrap px-4 py-3 text-gray-500 dark:text-dark-300">{{ field.type }}</td>
                    <td class="whitespace-nowrap px-4 py-3 text-gray-500 dark:text-dark-300">{{ t(field.required ? 'video.apiDocs.v2.fields.yes' : 'video.apiDocs.v2.fields.no') }}</td>
                    <td class="px-4 py-3 leading-6 text-gray-600 dark:text-dark-300">{{ t(field.description) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div class="mt-5 grid min-w-0 gap-4 xl:grid-cols-2">
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.imageGeneration')" :code="imageGenerationExample" />
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.imageEdit')" :code="imageEditExample" />
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.imageMultipart')" :code="imageMultipartExample" />
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.imageResponse')" :code="imageResponseExample" />
            </div>
          </section>

          <section id="upload-media" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.v2.upload.title')" :description="t('video.apiDocs.v2.upload.description')" />
            <EndpointTitle method="POST" path="/v1/videos/uploads" class="mt-5" />
            <p class="mt-3 text-sm leading-6 text-gray-600 dark:text-dark-300">{{ t('video.apiDocs.v2.upload.constraints') }}</p>
            <div class="mt-4">
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.upload')" :code="uploadExample" />
            </div>
            <div class="mt-4">
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.uploadResponse')" :code="uploadResponseExample" />
            </div>
            <div class="mt-4">
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.uploadMedia')" :code="uploadMediaExample" />
            </div>
          </section>

          <section id="create-job" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.v2.video.title')" :description="t('video.apiDocs.v2.video.description')" />
            <EndpointTitle method="POST" path="/v1/videos/generations" class="mt-5" />
            <div class="mt-5 overflow-x-auto rounded-lg border border-gray-200 dark:border-dark-700">
              <table class="min-w-[720px] w-full text-left text-sm">
                <thead class="bg-gray-50 text-xs uppercase text-gray-500 dark:bg-dark-800 dark:text-dark-300">
                  <tr>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.fields.field') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.fields.type') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.fields.required') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.fields.description') }}</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100 bg-white dark:divide-dark-700 dark:bg-dark-900">
                  <tr v-for="field in videoFields" :key="field.name">
                    <td class="whitespace-nowrap px-4 py-3 font-mono text-xs text-gray-900 dark:text-white">{{ field.name }}</td>
                    <td class="whitespace-nowrap px-4 py-3 text-gray-500 dark:text-dark-300">{{ field.type }}</td>
                    <td class="whitespace-nowrap px-4 py-3 text-gray-500 dark:text-dark-300">{{ t(field.required ? 'video.apiDocs.v2.fields.yes' : 'video.apiDocs.v2.fields.no') }}</td>
                    <td class="px-4 py-3 leading-6 text-gray-600 dark:text-dark-300">{{ t(field.description) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div class="mt-5 grid min-w-0 gap-4 xl:grid-cols-2">
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.videoStartFrame')" :code="videoStartFrameExample" />
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.videoFramePair')" :code="videoFramePairExample" />
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.videoReferences')" :code="videoReferencesExample" />
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.videoAudio')" :code="videoAudioExample" />
            </div>
          </section>

          <section id="model-matrix" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.v2.matrix.title')" :description="t('video.apiDocs.v2.matrix.description')" />
            <div class="mt-5 overflow-x-auto rounded-lg border border-gray-200 dark:border-dark-700">
              <table class="min-w-[900px] w-full text-left text-sm">
                <thead class="bg-gray-50 text-xs uppercase text-gray-500 dark:bg-dark-800 dark:text-dark-300">
                  <tr>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.matrix.model') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.matrix.resolution') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.matrix.duration') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.matrix.aspectRatio') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.v2.matrix.references') }}</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100 bg-white dark:divide-dark-700 dark:bg-dark-900">
                  <tr v-for="row in modelMatrixRows" :key="row.model">
                    <td class="whitespace-nowrap px-4 py-3 font-mono text-xs font-semibold text-gray-900 dark:text-white">{{ row.model }}</td>
                    <td class="px-4 py-3 leading-6 text-gray-600 dark:text-dark-300">{{ t(row.resolution) }}</td>
                    <td class="px-4 py-3 leading-6 text-gray-600 dark:text-dark-300">{{ t(row.duration) }}</td>
                    <td class="px-4 py-3 leading-6 text-gray-600 dark:text-dark-300">{{ t(row.aspectRatio) }}</td>
                    <td class="px-4 py-3 leading-6 text-gray-600 dark:text-dark-300">{{ t(row.references) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </section>

          <section id="model-examples" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.v2.examples.modelTitle')" :description="t('video.apiDocs.v2.examples.modelDescription')" />
            <div class="mt-5 grid min-w-0 gap-5">
              <div v-for="example in modelExamples" :key="example.model" class="min-w-0 border-b border-gray-200 pb-5 last:border-b-0 last:pb-0 dark:border-dark-700">
                <h3 class="font-mono text-sm font-semibold text-gray-900 dark:text-white">{{ example.model }}</h3>
                <p class="mt-1 text-sm leading-6 text-gray-600 dark:text-dark-300">{{ t(example.description) }}</p>
                <div class="mt-3">
                  <ApiCodeBlock :label="example.model" :code="example.code" />
                </div>
              </div>
            </div>
          </section>

          <section id="job-operations" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.v2.jobs.title')" :description="t('video.apiDocs.v2.jobs.description')" />
            <div class="mt-5 divide-y divide-gray-200 border-y border-gray-200 dark:divide-dark-700 dark:border-dark-700">
              <div v-for="endpoint in jobEndpoints" :key="`${endpoint.method}-${endpoint.path}`" class="grid gap-2 py-4 md:grid-cols-[90px_minmax(280px,1fr)_minmax(220px,0.8fr)] md:items-center md:gap-4">
                <span class="font-mono text-xs font-semibold" :class="endpoint.method === 'DELETE' ? 'text-red-600 dark:text-red-300' : 'text-emerald-700 dark:text-emerald-300'">{{ endpoint.method }}</span>
                <code class="min-w-0 break-all text-xs text-gray-900 dark:text-white">{{ endpoint.path }}</code>
                <p class="text-sm text-gray-500 dark:text-dark-300">{{ t(endpoint.description) }}</p>
              </div>
            </div>
            <div class="mt-5">
              <ApiCodeBlock :label="t('video.apiDocs.v2.examples.poll')" :code="pollExample" />
            </div>
          </section>

          <section id="status-errors" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.v2.status.title')" :description="t('video.apiDocs.v2.status.description')" />
            <div class="mt-5 grid gap-8 xl:grid-cols-2">
              <div>
                <h3 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('video.apiDocs.v2.status.lifecycle') }}</h3>
                <ol class="mt-4 border-l border-gray-200 pl-5 dark:border-dark-700">
                  <li v-for="status in statuses" :key="status.name" class="relative pb-5 last:pb-0">
                    <span class="absolute -left-[23px] top-1.5 h-2 w-2 rounded-full bg-primary-500 ring-4 ring-white dark:ring-dark-900"></span>
                    <code class="text-xs font-semibold text-gray-900 dark:text-white">{{ status.name }}</code>
                    <p class="mt-1 text-sm leading-6 text-gray-500 dark:text-dark-300">{{ t(status.description) }}</p>
                  </li>
                </ol>
              </div>
              <div>
                <h3 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('video.apiDocs.v2.errors.title') }}</h3>
                <dl class="mt-4 divide-y divide-gray-200 border-y border-gray-200 dark:divide-dark-700 dark:border-dark-700">
                  <div v-for="error in errors" :key="error.code" class="grid grid-cols-[52px_minmax(0,1fr)] gap-3 py-3">
                    <dt class="font-mono text-xs font-semibold text-red-600 dark:text-red-300">{{ error.code }}</dt>
                    <dd class="text-sm leading-6 text-gray-500 dark:text-dark-300">{{ t(error.description) }}</dd>
                  </div>
                </dl>
              </div>
            </div>
          </section>
        </main>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import ApiCodeBlock from '@/components/video/ApiCodeBlock.vue'
import EndpointTitle from '@/components/video/EndpointTitle.vue'
import Icon from '@/components/icons/Icon.vue'
import SectionHeading from '@/components/video/SectionHeading.vue'
import VideoSectionTabs from '@/components/video/VideoSectionTabs.vue'
import { buildGatewayUrl } from '@/api/client'
import userGroupsAPI from '@/api/groups'
import modelPlazaAPI, { type ModelPlazaGroup, type PlazaModel } from '@/api/modelPlaza'
import type { Group } from '@/types'

const { t } = useI18n()

const baseUrl = (() => {
  try {
    const fallback = typeof window === 'undefined' ? 'https://your-sub2-domain.example' : window.location.origin
    return new URL(buildGatewayUrl('/v1/models'), fallback).origin
  } catch {
    return 'https://your-sub2-domain.example'
  }
})()

const navigation = [
  { href: '#quick-start', label: 'video.apiDocs.v2.nav.quick' },
  { href: '#compatibility', label: 'video.apiDocs.v2.nav.compatibility' },
  { href: '#model-discovery', label: 'video.apiDocs.v2.nav.discovery' },
  { href: '#pricing', label: 'video.apiDocs.v2.nav.pricing' },
  { href: '#image-api', label: 'video.apiDocs.v2.nav.images' },
  { href: '#upload-media', label: 'video.apiDocs.v2.nav.upload' },
  { href: '#create-job', label: 'video.apiDocs.v2.nav.video' },
  { href: '#model-matrix', label: 'video.apiDocs.v2.nav.matrix' },
  { href: '#job-operations', label: 'video.apiDocs.v2.nav.jobs' },
  { href: '#status-errors', label: 'video.apiDocs.v2.nav.status' },
]

const compatibilityItems = [
  { title: 'video.apiDocs.v2.compatibility.pathsTitle', description: 'video.apiDocs.v2.compatibility.paths' },
  { title: 'video.apiDocs.v2.compatibility.keysTitle', description: 'video.apiDocs.v2.compatibility.keys' },
  { title: 'video.apiDocs.v2.compatibility.modelsTitle', description: 'video.apiDocs.v2.compatibility.models' },
  { title: 'video.apiDocs.v2.compatibility.contractTitle', description: 'video.apiDocs.v2.compatibility.contract' },
]

const imageFields = [
  { name: 'model', type: 'string', required: false, description: 'video.apiDocs.v2.imageFields.model' },
  { name: 'prompt', type: 'string', required: true, description: 'video.apiDocs.v2.imageFields.prompt' },
  { name: 'n', type: 'integer', required: false, description: 'video.apiDocs.v2.imageFields.n' },
  { name: 'size', type: 'string', required: false, description: 'video.apiDocs.v2.imageFields.size' },
  { name: 'response_format', type: 'string', required: false, description: 'video.apiDocs.v2.imageFields.responseFormat' },
  { name: 'quality', type: 'string', required: false, description: 'video.apiDocs.v2.imageFields.quality' },
  { name: 'images[] / mask', type: 'object / file', required: false, description: 'video.apiDocs.v2.imageFields.edits' },
]

const videoFields = [
  { name: 'model', type: 'string', required: true, description: 'video.apiDocs.v2.videoFields.model' },
  { name: 'prompt', type: 'string', required: true, description: 'video.apiDocs.v2.videoFields.prompt' },
  { name: 'resolution', type: 'string', required: false, description: 'video.apiDocs.v2.videoFields.resolution' },
  { name: 'duration', type: 'integer', required: false, description: 'video.apiDocs.v2.videoFields.duration' },
  { name: 'aspect_ratio', type: 'string', required: false, description: 'video.apiDocs.v2.videoFields.aspectRatio' },
  { name: 'audio', type: 'boolean', required: false, description: 'video.apiDocs.v2.videoFields.audio' },
  { name: 'prompt_enhance', type: 'string', required: false, description: 'video.apiDocs.v2.videoFields.promptEnhance' },
  { name: 'image_url', type: 'string', required: false, description: 'video.apiDocs.v2.videoFields.imageUrl' },
  { name: 'start_frame_url', type: 'string', required: false, description: 'video.apiDocs.v2.videoFields.startFrameUrl' },
  { name: 'end_frame_url', type: 'string', required: false, description: 'video.apiDocs.v2.videoFields.endFrameUrl' },
  { name: 'guidances', type: 'object', required: false, description: 'video.apiDocs.v2.videoFields.guidances' },
]

const modelMatrixRows = [
  {
    model: '<video-model-id>',
    resolution: 'video.apiDocs.v2.matrix.video.resolution',
    duration: 'video.apiDocs.v2.matrix.video.duration',
    aspectRatio: 'video.apiDocs.v2.matrix.video.aspectRatio',
    references: 'video.apiDocs.v2.matrix.video.references',
  },
  {
    model: '<image-model-id>',
    resolution: 'video.apiDocs.v2.matrix.image.resolution',
    duration: 'video.apiDocs.v2.matrix.image.duration',
    aspectRatio: 'video.apiDocs.v2.matrix.image.aspectRatio',
    references: 'video.apiDocs.v2.matrix.image.references',
  },
]

const jobEndpoints = [
  { method: 'GET', path: '/v1/videos/jobs?limit=50&status=running', description: 'video.apiDocs.v2.jobs.list' },
  { method: 'GET', path: '/v1/videos/jobs/{job_id}', description: 'video.apiDocs.v2.jobs.detail' },
  { method: 'DELETE', path: '/v1/videos/jobs/{job_id}', description: 'video.apiDocs.v2.jobs.cancel' },
  { method: 'GET', path: '/v1/videos/jobs/{job_id}/content', description: 'video.apiDocs.v2.jobs.content' },
]

const statuses = [
  { name: 'pending', description: 'video.apiDocs.v2.status.pending' },
  { name: 'running', description: 'video.apiDocs.v2.status.running' },
  { name: 'settling', description: 'video.apiDocs.v2.status.settling' },
  { name: 'completed', description: 'video.apiDocs.v2.status.completed' },
  { name: 'failed / canceled', description: 'video.apiDocs.v2.status.terminal' },
]

const errors = [
  { code: '400', description: 'video.apiDocs.v2.errors.badRequest' },
  { code: '401', description: 'video.apiDocs.v2.errors.unauthorized' },
  { code: '402', description: 'video.apiDocs.v2.errors.payment' },
  { code: '403', description: 'video.apiDocs.v2.errors.forbidden' },
  { code: '404', description: 'video.apiDocs.v2.errors.notFound' },
  { code: '409', description: 'video.apiDocs.v2.errors.conflict' },
  { code: '422', description: 'video.apiDocs.v2.errors.validation' },
  { code: '429', description: 'video.apiDocs.v2.errors.rateLimit' },
  { code: '502', description: 'video.apiDocs.v2.errors.serviceUnavailable' },
]

type PricingRow = {
  key: string
  group: string
  model: string
  type: string
  price: string
}

const groups = ref<Group[]>([])
const plazaGroups = ref<ModelPlazaGroup[]>([])
const userGroupRates = ref<Record<number, number>>({})
const groupsLoaded = ref(false)
const pricingLoading = ref(false)
const pricingError = ref(false)
const pricingUpdatedAt = ref<Date | null>(null)
let pricingRefreshTimer: ReturnType<typeof setInterval> | undefined

const mediaPlatforms = new Set(['leo', 'openai_media', 'video', 'composite'])
const isMediaGroup = (group: Pick<Group, 'platform'>) => mediaPlatforms.has(group.platform)
const isMediaPlazaGroup = (group: Pick<ModelPlazaGroup, 'platform'>) => mediaPlatforms.has(group.platform)

function trimPrice(value: number): string {
  return String(Math.round(value * 1_000_000) / 1_000_000)
}

function effectiveRate(group: Group): number {
  return userGroupRates.value[group.id] ?? group.rate_multiplier ?? 1
}

function imageRate(group: Group): number {
  return group.image_rate_independent ? group.image_rate_multiplier : effectiveRate(group)
}

function videoRate(group: Group): number {
  return group.video_rate_independent ? group.video_rate_multiplier : effectiveRate(group)
}

function formatUpdatedAt(value: Date): string {
  return value.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function modelType(model: PlazaModel): 'image' | 'video' {
  return model.pricing?.billing_mode === 'image' ? 'image' : 'video'
}

function formatPlazaPrice(model: PlazaModel, group: ModelPlazaGroup): string {
  const pricing = model.pricing
  if (!pricing) return t('video.apiDocs.v2.pricing.notConfigured')
  const rate = modelType(model) === 'image'
    ? (group.image_rate_independent ? group.image_rate_multiplier : (group.user_rate_multiplier ?? group.rate_multiplier))
    : (group.user_rate_multiplier ?? group.rate_multiplier)
  if (pricing.billing_mode === 'image' || pricing.billing_mode === 'per_request') {
    const tiers = (pricing.intervals ?? [])
      .filter((interval) => interval.per_request_price != null)
      .map((interval) => `${interval.tier_label || t('video.apiDocs.v2.pricing.perRequest')}: $${trimPrice((interval.per_request_price || 0) * rate)}`)
    const base = pricing.per_request_price == null ? '' : `$${trimPrice(pricing.per_request_price * rate)}`
    return tiers.length > 0 ? tiers.join(' / ') : (base || t('video.apiDocs.v2.pricing.notConfigured'))
  }
  const input = pricing.input_price == null ? '-' : `$${trimPrice(pricing.input_price * 1_000_000 * rate)}`
  const output = pricing.output_price == null ? '-' : `$${trimPrice(pricing.output_price * 1_000_000 * rate)}`
  return `${input} / ${output} ${t('video.apiDocs.v2.pricing.perMillionTokens')}`
}

function buildPricingRows(): PricingRow[] {
  const rows: PricingRow[] = []
  const availableGroupIDs = new Set(groups.value.filter(isMediaGroup).map((group) => group.id))
  const seen = new Set<string>()

  for (const plazaGroup of plazaGroups.value.filter(isMediaPlazaGroup)) {
    if (groupsLoaded.value && !availableGroupIDs.has(plazaGroup.id)) continue
    for (const model of plazaGroup.models) {
      const type = modelType(model)
      const key = `model:${plazaGroup.id}:${model.platform}:${model.name}`
      if (seen.has(key)) continue
      seen.add(key)
      rows.push({
        key,
        group: plazaGroup.name,
        model: model.name,
        type: type === 'image' ? t('video.apiDocs.v2.pricing.image') : t('video.apiDocs.v2.pricing.video'),
        price: formatPlazaPrice(model, plazaGroup),
      })
    }
  }

  for (const group of groups.value.filter(isMediaGroup)) {
    const imagePrices = [
      ['1K', group.image_price_1k],
      ['2K', group.image_price_2k],
      ['4K', group.image_price_4k],
    ].filter((entry): entry is [string, number] => entry[1] != null)
    if (imagePrices.length > 0 && !rows.some((row) => row.key.startsWith(`model:${group.id}:`) && row.type === t('video.apiDocs.v2.pricing.image'))) {
      rows.push({
        key: `image:${group.id}`,
        group: group.name,
        model: t('video.apiDocs.v2.pricing.allImageModels'),
        type: t('video.apiDocs.v2.pricing.image'),
        price: imagePrices.map(([tier, value]) => `${tier}: $${trimPrice(value * imageRate(group))}`).join(' / '),
      })
    }

    const modelPrices = Object.entries(group.video_model_prices ?? {})
    if (modelPrices.length > 0) {
      for (const [model, tiers] of modelPrices) {
        rows.push({
          key: `video:${group.id}:${model}`,
          group: group.name,
          model,
          type: t('video.apiDocs.v2.pricing.video'),
          price: Object.entries(tiers).map(([resolution, value]) => `${resolution}: $${trimPrice(value * videoRate(group))}/s`).join(' / '),
        })
      }
    } else {
      const flatPrices = [
        ['480p', group.video_price_480p],
        ['720p', group.video_price_720p],
        ['1080p', group.video_price_1080p],
      ].filter((entry): entry is [string, number] => entry[1] != null)
      if (flatPrices.length > 0) {
        rows.push({
          key: `video:${group.id}:all`,
          group: group.name,
          model: t('video.apiDocs.v2.pricing.allVideoModels'),
          type: t('video.apiDocs.v2.pricing.video'),
          price: flatPrices.map(([resolution, value]) => `${resolution}: $${trimPrice(value * videoRate(group))}/s`).join(' / '),
        })
      }
    }
  }
  return rows.sort((a, b) => a.group.localeCompare(b.group) || a.type.localeCompare(b.type) || a.model.localeCompare(b.model))
}

const pricingRows = computed(buildPricingRows)

async function loadPricing() {
  pricingLoading.value = true
  pricingError.value = false
  const [groupsResult, ratesResult, plazaResult] = await Promise.allSettled([
    userGroupsAPI.getAvailable(),
    userGroupsAPI.getUserGroupRates(),
    modelPlazaAPI.getModelPlaza(),
  ])
  if (groupsResult.status === 'fulfilled') {
    groups.value = groupsResult.value
    groupsLoaded.value = true
  }
  if (ratesResult.status === 'fulfilled') userGroupRates.value = ratesResult.value
  if (plazaResult.status === 'fulfilled') plazaGroups.value = plazaResult.value.groups
  pricingError.value = groupsResult.status === 'rejected' && plazaResult.status === 'rejected'
  pricingUpdatedAt.value = new Date()
  pricingLoading.value = false
}

onMounted(() => {
  void loadPricing()
  pricingRefreshTimer = setInterval(() => void loadPricing(), 60_000)
})

onBeforeUnmount(() => {
  if (pricingRefreshTimer) clearInterval(pricingRefreshTimer)
})

const authHeaders = `-H "Authorization: Bearer $SUB2_API_KEY" \\
  -H "Content-Type: application/json"`

const textRequestExample = `curl -X POST "${baseUrl}/v1/videos/generations" \\
  ${authHeaders} \\
  -H "Prefer: respond-async" \\
  -d '{
    "model": "<video-model-id>",
    "prompt": "A slow aerial shot over a coastal city at sunrise",
    "resolution": "720p",
    "duration": 8,
    "aspect_ratio": "16:9",
    "audio": false
  }'`

const acceptedResponseExample = `HTTP/1.1 202 Accepted
Preference-Applied: respond-async
Location: /v1/videos/jobs/vidjob_example

{
  "job_id": "vidjob_example",
  "status": "pending",
  "status_url": "/v1/videos/jobs/vidjob_example"
}`

const modelListExample = `curl "${baseUrl}/v1/models" \\
  -H "Authorization: Bearer $SUB2_API_KEY"

{
  "object": "list",
  "data": [
    { "id": "<video-model-id>", "type": "model" },
    { "id": "<image-model-id>", "type": "model" }
  ]
}`

const imageGenerationExample = `curl -X POST "${baseUrl}/v1/images/generations" \\
  ${authHeaders} \\
  -d '{
    "model": "<image-model-id>",
    "prompt": "A clean product photo on a white background",
    "n": 1,
    "size": "1024x1024",
    "response_format": "b64_json"
  }'`

const imageEditExample = `curl -X POST "${baseUrl}/v1/images/edits" \\
  ${authHeaders} \\
  -d '{
    "model": "<image-model-id>",
    "prompt": "Replace the background with a quiet studio",
    "images": [{ "image_url": "https://media.example/input.png" }],
    "mask": { "image_url": "https://media.example/mask.png" }
  }'`

const imageMultipartExample = `curl -X POST "${baseUrl}/v1/images/edits" \\
  -H "Authorization: Bearer $SUB2_API_KEY" \\
  -F "model=<image-model-id>" \\
  -F "prompt=Turn the sketch into a polished illustration" \\
  -F "image[]=@./input.png" \\
  -F "mask=@./mask.png"`

const imageResponseExample = `{
  "created": 1780000000,
  "data": [
    { "b64_json": "<base64-image>" }
  ]
}`

const uploadExample = `curl -X POST "${baseUrl}/v1/videos/uploads" \\
  -H "Authorization: Bearer $SUB2_API_KEY" \\
  -F "image=@./reference.png"`

const uploadResponseExample = `{
  "upload_id": "upload_example",
  "media_url": "https://media.example/uploaded/reference.png",
  "media_type": "image",
  "content_type": "image/png",
  "size": 428516
}`

const uploadMediaExample = [
  `curl -X POST "${baseUrl}/v1/videos/uploads" -H "Authorization: Bearer $SUB2_API_KEY" -F "video=@./reference.mp4"`,
  `curl -X POST "${baseUrl}/v1/videos/uploads" -H "Authorization: Bearer $SUB2_API_KEY" -F "audio=@./reference.mp3"`,
].join('\n\n')

const videoStartFrameExample = `{
  "model": "<video-model-id>",
  "prompt": "Animate the camera through this scene",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "start_frame_url": "https://media.example/start-frame.png"
}`

const videoFramePairExample = `{
  "model": "<video-model-id>",
  "prompt": "Move smoothly from the opening composition to the closing composition",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "start_frame_url": "https://media.example/start-frame.png",
  "end_frame_url": "https://media.example/end-frame.png"
}`

const videoReferencesExample = `{
  "model": "<video-model-id>",
  "prompt": "Keep the subject and product consistent",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "guidances": {
    "image_reference": [{
      "image": { "url": "https://media.example/reference.png", "type": "UPLOADED" },
      "strength": "MID",
      "order": 0
    }],
    "video_reference_base": [{
      "video": { "url": "https://media.example/reference.mp4", "type": "UPLOADED" }
    }]
  }
}`

const videoAudioExample = `{
  "model": "<video-model-id>",
  "prompt": "Match the rhythm and pacing of the reference audio",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "audio": true,
  "guidances": {
    "audio_reference": [{
      "audio": { "url": "https://media.example/reference.mp3", "type": "UPLOADED" }
    }]
  }
}`

const pollExample = `curl "${baseUrl}/v1/videos/jobs/vidjob_example" \\
  -H "Authorization: Bearer $SUB2_API_KEY"

curl -o output.mp4 "${baseUrl}/v1/videos/jobs/vidjob_example/content" \\
  -H "Authorization: Bearer $SUB2_API_KEY"`

const modelExamples = [
  {
    model: '<video-model-id>',
    description: 'video.apiDocs.v2.examples.videoModelDescription',
    code: textRequestExample,
  },
  {
    model: '<image-model-id>',
    description: 'video.apiDocs.v2.examples.imageModelDescription',
    code: imageGenerationExample,
  },
]
</script>
