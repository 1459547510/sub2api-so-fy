<template>
  <AppLayout>
    <div class="mx-auto max-w-[1320px] space-y-5">
      <header class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-sm font-medium text-primary-600 dark:text-primary-400">{{ t('video.apiDocs.eyebrow') }}</p>
          <h1 class="mt-1 text-2xl font-semibold tracking-tight text-gray-900 dark:text-white">{{ t('video.apiDocs.title') }}</h1>
          <p class="mt-1 max-w-3xl text-sm text-gray-500 dark:text-dark-300">{{ t('video.apiDocs.description') }}</p>
        </div>
        <div class="inline-flex h-8 flex-shrink-0 items-center gap-2 self-start rounded-md border border-gray-200 px-2.5 font-mono text-xs text-gray-600 dark:border-dark-700 dark:text-dark-300 sm:self-auto">
          <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
          {{ t('video.apiDocs.badge') }}
        </div>
      </header>

      <VideoSectionTabs />

      <div class="grid min-w-0 gap-8 lg:grid-cols-[210px_minmax(0,1fr)] lg:gap-12">
        <aside class="hidden lg:block">
          <nav class="sticky top-6 border-l border-gray-200 pl-4 dark:border-dark-700" :aria-label="t('video.apiDocs.tocLabel')">
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
            <SectionHeading :title="t('video.apiDocs.quick.title')" :description="t('video.apiDocs.quick.description')" />
            <dl class="mt-5 grid gap-px overflow-hidden rounded-lg border border-gray-200 bg-gray-200 sm:grid-cols-3 dark:border-dark-700 dark:bg-dark-700">
              <div class="min-w-0 bg-white p-4 dark:bg-dark-900">
                <dt class="text-xs font-medium uppercase text-gray-400 dark:text-dark-400">{{ t('video.apiDocs.quick.baseUrl') }}</dt>
                <dd class="mt-2 break-all font-mono text-sm text-gray-800 dark:text-gray-100">{{ baseUrl }}</dd>
              </div>
              <div class="min-w-0 bg-white p-4 dark:bg-dark-900">
                <dt class="text-xs font-medium uppercase text-gray-400 dark:text-dark-400">{{ t('video.apiDocs.quick.auth') }}</dt>
                <dd class="mt-2 font-mono text-sm text-gray-800 dark:text-gray-100">Bearer API Key</dd>
              </div>
              <div class="min-w-0 bg-white p-4 dark:bg-dark-900">
                <dt class="text-xs font-medium uppercase text-gray-400 dark:text-dark-400">{{ t('video.apiDocs.quick.mode') }}</dt>
                <dd class="mt-2 font-mono text-sm text-gray-800 dark:text-gray-100">Prefer: respond-async</dd>
              </div>
            </dl>
            <div class="mt-5">
              <ApiCodeBlock :label="t('video.apiDocs.examples.textToVideo')" :code="textRequestExample" />
            </div>
            <div class="mt-4">
              <ApiCodeBlock :label="t('video.apiDocs.examples.acceptedResponse')" :code="acceptedResponseExample" />
            </div>
          </section>

          <section id="upload-image" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.upload.title')" :description="t('video.apiDocs.upload.description')" />
            <EndpointTitle method="POST" path="/v1/videos/uploads" class="mt-5" />
            <p class="mt-3 text-sm leading-6 text-gray-600 dark:text-dark-300">{{ t('video.apiDocs.upload.constraints') }}</p>
            <div class="mt-4">
              <ApiCodeBlock :label="t('video.apiDocs.examples.upload')" :code="uploadExample" />
            </div>
            <div class="mt-4">
              <ApiCodeBlock :label="t('video.apiDocs.examples.uploadResponse')" :code="uploadResponseExample" />
            </div>
          </section>

          <section id="create-job" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.create.title')" :description="t('video.apiDocs.create.description')" />
            <EndpointTitle method="POST" path="/v1/videos/generations" class="mt-5" />
            <div class="mt-5 overflow-x-auto rounded-lg border border-gray-200 dark:border-dark-700">
              <table class="min-w-[720px] w-full text-left text-sm">
                <thead class="bg-gray-50 text-xs uppercase text-gray-500 dark:bg-dark-800 dark:text-dark-300">
                  <tr>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.fields.field') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.fields.type') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.fields.required') }}</th>
                    <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.fields.description') }}</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100 bg-white dark:divide-dark-700 dark:bg-dark-900">
                  <tr v-for="field in requestFields" :key="field.name">
                    <td class="whitespace-nowrap px-4 py-3 font-mono text-xs text-gray-900 dark:text-white">{{ field.name }}</td>
                    <td class="whitespace-nowrap px-4 py-3 text-gray-500 dark:text-dark-300">{{ field.type }}</td>
                    <td class="whitespace-nowrap px-4 py-3 text-gray-500 dark:text-dark-300">{{ t(field.required ? 'video.apiDocs.fields.yes' : 'video.apiDocs.fields.no') }}</td>
                    <td class="px-4 py-3 leading-6 text-gray-600 dark:text-dark-300">{{ t(field.description) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div class="mt-6 grid min-w-0 gap-4 xl:grid-cols-2">
              <ApiCodeBlock :label="t('video.apiDocs.examples.singleImage')" :code="singleImageExample" />
              <ApiCodeBlock :label="t('video.apiDocs.examples.framePair')" :code="framePairExample" />
              <ApiCodeBlock :label="t('video.apiDocs.examples.multiImage')" :code="multiImageExample" />
            </div>
          </section>

          <section id="job-operations" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.jobs.title')" :description="t('video.apiDocs.jobs.description')" />
            <div class="mt-5 divide-y divide-gray-200 border-y border-gray-200 dark:divide-dark-700 dark:border-dark-700">
              <div v-for="endpoint in jobEndpoints" :key="`${endpoint.method}-${endpoint.path}`" class="grid gap-2 py-4 md:grid-cols-[90px_minmax(280px,1fr)_minmax(220px,0.8fr)] md:items-center md:gap-4">
                <span class="font-mono text-xs font-semibold" :class="endpoint.method === 'DELETE' ? 'text-red-600 dark:text-red-300' : 'text-emerald-700 dark:text-emerald-300'">{{ endpoint.method }}</span>
                <code class="min-w-0 break-all text-xs text-gray-900 dark:text-white">{{ endpoint.path }}</code>
                <p class="text-sm text-gray-500 dark:text-dark-300">{{ t(endpoint.description) }}</p>
              </div>
            </div>
            <div class="mt-5">
              <ApiCodeBlock :label="t('video.apiDocs.examples.poll')" :code="pollExample" />
            </div>
          </section>

          <section id="status-errors" class="scroll-mt-6 border-t border-gray-200 pt-10 dark:border-dark-700">
            <SectionHeading :title="t('video.apiDocs.status.title')" :description="t('video.apiDocs.status.description')" />
            <div class="mt-5 grid gap-8 xl:grid-cols-2">
              <div>
                <h3 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('video.apiDocs.status.lifecycle') }}</h3>
                <ol class="mt-4 border-l border-gray-200 pl-5 dark:border-dark-700">
                  <li v-for="status in statuses" :key="status.name" class="relative pb-5 last:pb-0">
                    <span class="absolute -left-[23px] top-1.5 h-2 w-2 rounded-full bg-primary-500 ring-4 ring-white dark:ring-dark-900"></span>
                    <code class="text-xs font-semibold text-gray-900 dark:text-white">{{ status.name }}</code>
                    <p class="mt-1 text-sm leading-6 text-gray-500 dark:text-dark-300">{{ t(status.description) }}</p>
                  </li>
                </ol>
              </div>
              <div>
                <h3 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('video.apiDocs.errors.title') }}</h3>
                <dl class="mt-4 divide-y divide-gray-200 border-y border-gray-200 dark:divide-dark-700 dark:border-dark-700">
                  <div v-for="error in errors" :key="error.code" class="grid grid-cols-[52px_minmax(0,1fr)] gap-3 py-3">
                    <dt class="font-mono text-xs font-semibold text-red-600 dark:text-red-300">{{ error.code }}</dt>
                    <dd class="text-sm leading-6 text-gray-500 dark:text-dark-300">{{ t(error.description) }}</dd>
                  </div>
                </dl>
              </div>
            </div>
            <div class="mt-8 flex items-start gap-3 border-l-2 border-primary-500 bg-primary-50 px-4 py-3 text-sm leading-6 text-primary-900 dark:bg-primary-900/20 dark:text-primary-100">
              <Icon name="shield" size="md" class="mt-0.5 flex-shrink-0" />
              <p>{{ t('video.apiDocs.privacy') }}</p>
            </div>
          </section>
        </main>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import ApiCodeBlock from '@/components/video/ApiCodeBlock.vue'
import EndpointTitle from '@/components/video/EndpointTitle.vue'
import SectionHeading from '@/components/video/SectionHeading.vue'
import VideoSectionTabs from '@/components/video/VideoSectionTabs.vue'
import { buildGatewayUrl } from '@/api/client'

const { t } = useI18n()

const baseUrl = (() => {
  try {
    const fallback = typeof window === 'undefined' ? 'https://your-sub2-domain.example' : window.location.origin
    return new URL(buildGatewayUrl('/v1/videos/generations'), fallback).origin
  } catch {
    return 'https://your-sub2-domain.example'
  }
})()

const navigation = [
  { href: '#quick-start', label: 'video.apiDocs.nav.quick' },
  { href: '#upload-image', label: 'video.apiDocs.nav.upload' },
  { href: '#create-job', label: 'video.apiDocs.nav.create' },
  { href: '#job-operations', label: 'video.apiDocs.nav.jobs' },
  { href: '#status-errors', label: 'video.apiDocs.nav.status' },
]

const requestFields = [
  { name: 'model', type: 'string', required: true, description: 'video.apiDocs.fields.model' },
  { name: 'prompt', type: 'string', required: true, description: 'video.apiDocs.fields.prompt' },
  { name: 'resolution', type: 'string', required: false, description: 'video.apiDocs.fields.resolution' },
  { name: 'duration', type: 'integer', required: false, description: 'video.apiDocs.fields.duration' },
  { name: 'aspect_ratio', type: 'string', required: false, description: 'video.apiDocs.fields.aspectRatio' },
  { name: 'audio', type: 'boolean', required: false, description: 'video.apiDocs.fields.audio' },
  { name: 'image_url', type: 'string', required: false, description: 'video.apiDocs.fields.imageUrl' },
  { name: 'start_frame_url', type: 'string', required: false, description: 'video.apiDocs.fields.startFrameUrl' },
  { name: 'end_frame_url', type: 'string', required: false, description: 'video.apiDocs.fields.endFrameUrl' },
  { name: 'guidances', type: 'object', required: false, description: 'video.apiDocs.fields.guidances' },
]

const jobEndpoints = [
  { method: 'GET', path: '/v1/videos/jobs?limit=50&status=running', description: 'video.apiDocs.jobs.list' },
  { method: 'GET', path: '/v1/videos/jobs/{job_id}', description: 'video.apiDocs.jobs.detail' },
  { method: 'DELETE', path: '/v1/videos/jobs/{job_id}', description: 'video.apiDocs.jobs.cancel' },
  { method: 'GET', path: '/v1/videos/jobs/{job_id}/content', description: 'video.apiDocs.jobs.content' },
]

const statuses = [
  { name: 'pending', description: 'video.apiDocs.status.pending' },
  { name: 'running', description: 'video.apiDocs.status.running' },
  { name: 'settling', description: 'video.apiDocs.status.settling' },
  { name: 'completed', description: 'video.apiDocs.status.completed' },
  { name: 'failed / canceled', description: 'video.apiDocs.status.terminal' },
]

const errors = [
  { code: '400', description: 'video.apiDocs.errors.badRequest' },
  { code: '401', description: 'video.apiDocs.errors.unauthorized' },
  { code: '402', description: 'video.apiDocs.errors.payment' },
  { code: '403', description: 'video.apiDocs.errors.forbidden' },
  { code: '404', description: 'video.apiDocs.errors.notFound' },
  { code: '409', description: 'video.apiDocs.errors.conflict' },
  { code: '422', description: 'video.apiDocs.errors.validation' },
  { code: '502', description: 'video.apiDocs.errors.upstream' },
]

const authHeaders = `-H "Authorization: Bearer $SUB2_API_KEY" \\\n  -H "Content-Type: application/json"`

const textRequestExample = `curl -X POST "${baseUrl}/v1/videos/generations" \\\n  ${authHeaders} \\\n  -H "Prefer: respond-async" \\\n  -d '{
    "model": "seedance-2.0",
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

const uploadExample = `curl -X POST "${baseUrl}/v1/videos/uploads" \\\n  -H "Authorization: Bearer $SUB2_API_KEY" \\\n  -F "image=@./reference.png"`

const uploadResponseExample = `{
  "upload_id": "local-input-token",
  "image_url": "http://127.0.0.1:8080/internal/video-inputs/local-input-token",
  "content_type": "image/png",
  "size": 428516
}`

const singleImageExample = `{
  "model": "seedance-2.0",
  "prompt": "Animate the camera through this scene",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "image_url": "https://media.example/start-frame.png"
}`

const framePairExample = `{
  "model": "seedance-2.0",
  "prompt": "Move smoothly from the opening composition to the closing composition",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "start_frame_url": "https://media.example/start-frame.png",
  "end_frame_url": "https://media.example/end-frame.png",
  "guidances": {
    "image_reference": [
      {
        "image": { "url": "https://media.example/character.png", "type": "UPLOADED" },
        "strength": "MID",
        "order": 0
      }
    ]
  }
}`

const multiImageExample = `{
  "model": "seedance-2.0",
  "prompt": "Keep the character and product consistent",
  "resolution": "720p",
  "duration": 8,
  "aspect_ratio": "16:9",
  "guidances": {
    "image_reference": [
      {
        "image": { "url": "https://media.example/character.png", "type": "UPLOADED" },
        "strength": "MID",
        "order": 0
      },
      {
        "image": { "url": "https://media.example/product.png", "type": "UPLOADED" },
        "strength": "MID",
        "order": 1
      }
    ]
  }
}`

const pollExample = `curl "${baseUrl}/v1/videos/jobs/vidjob_example" \\\n  -H "Authorization: Bearer $SUB2_API_KEY"

curl -o output.mp4 "${baseUrl}/v1/videos/jobs/vidjob_example/content" \\\n  -H "Authorization: Bearer $SUB2_API_KEY"`
</script>
