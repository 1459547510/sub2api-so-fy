<template>
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
          <ApiCodeBlock :label="t('video.apiDocs.v2.examples.imageReference')" :code="imageReferenceExample" />
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
          <table class="min-w-[1120px] w-full text-left text-sm">
            <thead class="bg-gray-50 text-xs uppercase text-gray-500 dark:bg-dark-800 dark:text-dark-300">
              <tr>
                <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.matrix.model') }}</th>
                <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.matrix.resolution') }}</th>
                <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.matrix.duration') }}</th>
                <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.matrix.aspectRatio') }}</th>
                <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.matrix.promptLimit') }}</th>
                <th class="px-4 py-3 font-medium">{{ t('video.apiDocs.matrix.references') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 bg-white dark:divide-dark-700 dark:bg-dark-900">
              <tr v-for="row in videoModelMatrixRows" :key="row.model">
                <td class="whitespace-nowrap px-4 py-3 font-mono text-xs font-semibold text-gray-900 dark:text-white">{{ row.model }}</td>
                <td class="px-4 py-3 leading-6 text-gray-600 dark:text-dark-300">{{ t(row.resolution) }}</td>
                <td class="px-4 py-3 leading-6 text-gray-600 dark:text-dark-300">{{ t(row.duration) }}</td>
                <td class="px-4 py-3 leading-6 text-gray-600 dark:text-dark-300">{{ t(row.aspectRatio) }}</td>
                <td class="whitespace-nowrap px-4 py-3 text-gray-600 dark:text-dark-300">{{ t(row.promptLimit) }}</td>
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
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import ApiCodeBlock from '@/components/video/ApiCodeBlock.vue'
import EndpointTitle from '@/components/video/EndpointTitle.vue'
import SectionHeading from '@/components/video/SectionHeading.vue'
import {
  buildV2VideoModelExamples,
  resolveVideoApiDocsBaseUrl,
  v2VideoModelMatrixRows,
  videoApiAuthHeaders,
} from '@/utils/videoApiDocs'

const { t } = useI18n()
const baseUrl = resolveVideoApiDocsBaseUrl()
const authHeaders = videoApiAuthHeaders()
const modelExamples = buildV2VideoModelExamples(baseUrl)
const videoModelMatrixRows = v2VideoModelMatrixRows

const navigation = [
  { href: '#quick-start', label: 'video.apiDocs.v2.nav.quick' },
  { href: '#compatibility', label: 'video.apiDocs.v2.nav.compatibility' },
  { href: '#model-discovery', label: 'video.apiDocs.v2.nav.discovery' },
  { href: '#image-api', label: 'video.apiDocs.v2.nav.images' },
  { href: '#upload-media', label: 'video.apiDocs.v2.nav.upload' },
  { href: '#create-job', label: 'video.apiDocs.v2.nav.video' },
  { href: '#model-matrix', label: 'video.apiDocs.v2.nav.matrix' },
  { href: '#model-examples', label: 'video.apiDocs.v2.nav.examples' },
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
  { name: 'image_urls', type: 'string[]', required: false, description: 'video.apiDocs.v2.imageFields.imageUrls' },
  { name: 'images[] / mask', type: 'object / file', required: false, description: 'video.apiDocs.v2.imageFields.edits' },
]

const videoFields = [
  { name: 'model', type: 'string', required: true, description: 'video.apiDocs.fields.model' },
  { name: 'prompt', type: 'string', required: true, description: 'video.apiDocs.fields.prompt' },
  { name: 'resolution', type: 'string', required: false, description: 'video.apiDocs.fields.resolution' },
  { name: 'duration', type: 'integer', required: false, description: 'video.apiDocs.fields.duration' },
  { name: 'aspect_ratio', type: 'string', required: false, description: 'video.apiDocs.fields.aspectRatio' },
  { name: 'audio', type: 'boolean', required: false, description: 'video.apiDocs.fields.audio' },
  { name: 'prompt_enhance', type: 'string', required: false, description: 'video.apiDocs.fields.promptEnhance' },
  { name: 'image_url', type: 'string', required: false, description: 'video.apiDocs.fields.imageUrl' },
  { name: 'start_frame_url', type: 'string', required: false, description: 'video.apiDocs.fields.startFrameUrl' },
  { name: 'end_frame_url', type: 'string', required: false, description: 'video.apiDocs.fields.endFrameUrl' },
  { name: 'guidances', type: 'object', required: false, description: 'video.apiDocs.fields.guidancesPublic' },
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

const textRequestExample = `curl -X POST "${baseUrl}/v1/videos/generations" \\
  ${authHeaders} \\
  -H "Prefer: respond-async" \\
  -d '{
    "model": "seedance-2.0",
    "prompt": "A slow aerial shot over a coastal city at sunrise",
    "resolution": "720p",
    "duration": 5,
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
    { "id": "seedance-2.0", "type": "model" },
    { "id": "gpt-image-2", "type": "model" }
  ]
}`

const imageGenerationExample = `curl -X POST "${baseUrl}/v1/images/generations" \\
  ${authHeaders} \\
  -d '{
    "model": "gpt-image-2",
    "prompt": "A clean product photo on a white background",
    "n": 1,
    "size": "1024x1024",
    "response_format": "b64_json"
  }'`

const imageReferenceExample = `curl -X POST "${baseUrl}/v1/images/generations" \\
  ${authHeaders} \\
  -d '{
    "model": "<image-model-id>",
    "prompt": "Keep the product and place it on a marble counter",
    "n": 1,
    "aspect_ratio": "1:1",
    "image_urls": ["https://media.example/reference.png"]
  }'`

const imageEditExample = `curl -X POST "${baseUrl}/v1/images/edits" \\
  ${authHeaders} \\
  -d '{
    "model": "gpt-image-2",
    "prompt": "Replace the background with a quiet studio",
    "images": [{ "image_url": "https://media.example/input.png" }],
    "mask": { "image_url": "https://media.example/mask.png" }
  }'`

const imageMultipartExample = `curl -X POST "${baseUrl}/v1/images/edits" \\
  -H "Authorization: Bearer $SUB2_API_KEY" \\
  -F "model=gpt-image-2" \\
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
  "model": "seedance-2.0",
  "prompt": "Animate the camera through this scene",
  "resolution": "720p",
  "duration": 5,
  "aspect_ratio": "16:9",
  "start_frame_url": "https://media.example/start-frame.png"
}`

const videoFramePairExample = `{
  "model": "seedance-2.0",
  "prompt": "Move smoothly from the opening composition to the closing composition",
  "resolution": "720p",
  "duration": 5,
  "aspect_ratio": "16:9",
  "start_frame_url": "https://media.example/start-frame.png",
  "end_frame_url": "https://media.example/end-frame.png"
}`

const videoReferencesExample = `{
  "model": "seedance-2.0",
  "prompt": "Keep the subject and product consistent",
  "resolution": "720p",
  "duration": 5,
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
  "model": "seedance-2.0",
  "prompt": "Match the rhythm and pacing of the reference audio",
  "resolution": "720p",
  "duration": 5,
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
</script>
