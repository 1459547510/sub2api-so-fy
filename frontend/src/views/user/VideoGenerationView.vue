<template>
  <AppLayout>
    <div class="mx-auto max-w-[1500px] space-y-5">
      <header class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-sm font-medium text-primary-600 dark:text-primary-400">{{ t('video.eyebrow') }}</p>
          <h1 class="mt-1 text-2xl font-semibold tracking-tight text-gray-900 dark:text-white">{{ t('video.title') }}</h1>
          <p class="mt-1 max-w-2xl text-sm text-gray-500 dark:text-dark-300">{{ t('video.description') }}</p>
        </div>
        <div class="flex items-center gap-2 text-xs text-gray-500 dark:text-dark-400">
          <Icon name="cloud" size="sm" />
          <span>{{ activeJobs.length ? t('video.polling') : t('video.queueReady') }}</span>
        </div>
      </header>

      <div class="grid min-w-0 gap-4 xl:grid-cols-[400px_minmax(0,1fr)]">
        <section class="card min-w-0 p-5" data-testid="video-settings">
          <div class="flex items-center justify-between gap-3 border-b border-gray-100 pb-4 dark:border-dark-700">
            <div>
              <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('video.settings') }}</h2>
              <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('video.settingsHint') }}</p>
            </div>
            <Icon name="sparkles" size="lg" class="text-primary-500" />
          </div>

          <div v-if="!leoKeys.length && !loadingKeys" class="mt-5 rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800 dark:border-amber-900/60 dark:bg-amber-900/20 dark:text-amber-200" data-testid="video-no-key">
            <div class="flex items-start gap-3">
              <Icon name="key" size="md" class="mt-0.5 flex-shrink-0" />
              <div>
                <p class="font-medium">{{ t('video.noKey') }}</p>
                <p class="mt-1 text-xs opacity-80">{{ t('video.noKeyHint') }}</p>
              </div>
            </div>
          </div>

          <form class="mt-5 space-y-4" @submit.prevent="submitJob">
            <label class="block">
              <span class="input-label">{{ t('video.apiKey') }}</span>
              <select v-model="selectedKeyId" class="input" :disabled="loadingKeys || !leoKeys.length" data-testid="video-api-key">
                <option v-for="key in leoKeys" :key="key.id" :value="key.id">{{ key.name || `API Key #${key.id}` }}</option>
              </select>
            </label>

            <label class="block">
              <span class="input-label">{{ t('video.prompt') }}</span>
              <textarea
                v-model="prompt"
                rows="5"
                class="input min-h-[128px] resize-y"
                :placeholder="t('video.promptPlaceholder')"
                data-testid="video-prompt"
                @keydown.meta.enter.prevent="submitJob"
                @keydown.ctrl.enter.prevent="submitJob"
              ></textarea>
            </label>

            <div class="grid grid-cols-2 gap-3">
              <label class="block">
                <span class="input-label">{{ t('video.model') }}</span>
                <select v-model="model" class="input" data-testid="video-model">
                  <option v-for="option in modelOptions" :key="option" :value="option">{{ option }}</option>
                </select>
              </label>
              <label class="block">
                <span class="input-label">{{ t('video.resolution') }}</span>
                <select v-model="resolution" class="input" data-testid="video-resolution">
                  <option value="480p">480p</option>
                  <option value="720p">720p</option>
                  <option value="1080p">1080p</option>
                </select>
              </label>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <label class="block">
                <span class="input-label">{{ t('video.aspectRatio') }}</span>
                <select v-model="aspectRatio" class="input" data-testid="video-aspect-ratio">
                  <option v-for="option in aspectRatioOptions" :key="option" :value="option">{{ option }}</option>
                </select>
              </label>
              <label class="block">
                <span class="input-label">{{ t('video.duration') }}</span>
                <input v-model.number="duration" type="number" min="4" max="15" step="1" class="input" data-testid="video-duration" />
              </label>
            </div>

            <label class="flex cursor-pointer items-center justify-between rounded-lg border border-gray-200 px-3 py-2.5 dark:border-dark-700">
              <span>
                <span class="block text-sm font-medium text-gray-800 dark:text-gray-100">{{ t('video.audio') }}</span>
                <span class="block text-xs text-gray-500 dark:text-dark-400">{{ t('video.audioHint') }}</span>
              </span>
              <input v-model="audio" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500" data-testid="video-audio" />
            </label>

            <div>
              <span class="input-label">{{ t('video.imageInput') }}</span>
              <div class="grid grid-cols-3 gap-1 rounded-lg bg-gray-100 p-1 dark:bg-dark-800" role="tablist">
                <button v-for="option in imageModes" :key="option.value" type="button" class="rounded-md px-2 py-2 text-xs font-medium transition-colors" :class="imageMode === option.value ? 'bg-white text-primary-700 shadow-sm dark:bg-dark-600 dark:text-primary-200' : 'text-gray-500 hover:text-gray-900 dark:text-dark-300 dark:hover:text-white'" :data-testid="`mode-${option.value}`" @click="imageMode = option.value">
                  {{ t(option.label) }}
                </button>
              </div>
            </div>

            <div v-if="imageMode === 'local'" class="space-y-3">
              <label class="flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-dashed border-gray-300 px-3 py-4 text-sm text-gray-600 transition-colors hover:border-primary-400 hover:text-primary-600 dark:border-dark-600 dark:text-dark-300 dark:hover:border-primary-500 dark:hover:text-primary-300">
                <Icon name="upload" size="sm" />
                <span>{{ localFile ? t('video.replaceImage') : t('video.chooseImage') }}</span>
                <input ref="fileInput" type="file" accept="image/png,image/jpeg,image/webp" class="sr-only" data-testid="video-image-file" @change="onFileChange" />
              </label>
              <div v-if="previewUrl" class="relative overflow-hidden rounded-lg border border-gray-200 dark:border-dark-700">
                <img :src="previewUrl" alt="" class="max-h-40 w-full object-cover" />
                <button type="button" class="absolute right-2 top-2 inline-flex h-8 w-8 items-center justify-center rounded-md bg-black/60 text-white hover:bg-black/80" :title="t('video.removeImage')" @click="removeLocalImage">
                  <Icon name="x" size="sm" />
                </button>
              </div>
            </div>

            <label v-else-if="imageMode === 'url'" class="block">
              <span class="sr-only">{{ t('video.imageUrl') }}</span>
              <input v-model="imageUrl" type="url" class="input" :placeholder="t('video.imageUrlPlaceholder')" data-testid="video-image-url" />
            </label>

            <button type="submit" class="btn btn-primary flex w-full items-center justify-center gap-2" :disabled="!canSubmit || submitting || uploading" data-testid="submit-video">
              <Icon :name="submitting || uploading ? 'refresh' : 'play'" size="sm" :class="submitting || uploading ? 'animate-spin' : ''" />
              {{ submitting ? t('video.submitting') : uploading ? t('video.uploading') : t('video.submit') }}
            </button>
          </form>
        </section>

        <section class="card min-w-0 overflow-hidden" data-testid="video-results">
          <div class="border-b border-gray-100 px-5 py-4 dark:border-dark-700">
            <div class="flex items-center justify-between gap-3">
              <div>
                <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('video.results') }}</h2>
                <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('video.resultsHint') }}</p>
              </div>
              <button type="button" class="btn btn-secondary btn-sm" :disabled="loadingJobs || !selectedKey" :title="t('common.refresh')" @click="loadJobs">
                <Icon name="refresh" size="sm" :class="loadingJobs ? 'animate-spin' : ''" />
              </button>
            </div>
          </div>

          <div class="border-b border-gray-100 bg-gray-50/70 p-5 dark:border-dark-700 dark:bg-dark-900/40">
            <div v-if="selectedVideoUrl" class="overflow-hidden rounded-lg bg-black">
              <video :src="selectedVideoUrl" controls playsinline class="aspect-video max-h-[440px] w-full" data-testid="video-preview"></video>
            </div>
            <div v-else class="flex min-h-[230px] flex-col items-center justify-center rounded-lg border border-dashed border-gray-200 bg-white px-6 text-center dark:border-dark-700 dark:bg-dark-900/60">
              <Icon name="play" size="xl" class="mb-3 text-gray-300 dark:text-dark-600" />
              <p class="text-sm font-medium text-gray-700 dark:text-gray-200">{{ t('video.previewEmpty') }}</p>
              <p class="mt-1 max-w-sm text-xs text-gray-500 dark:text-dark-400">{{ t('video.previewEmptyHint') }}</p>
            </div>
            <div v-if="selectedJob" class="mt-3 flex flex-wrap items-center justify-between gap-2 text-xs text-gray-500 dark:text-dark-400">
              <span class="truncate">{{ selectedJob.prompt }}</span>
              <a v-if="selectedVideoUrl" :href="selectedVideoUrl" download class="inline-flex items-center gap-1 font-medium text-primary-600 hover:text-primary-700 dark:text-primary-300" :title="t('video.download')">
                <Icon name="download" size="sm" />
                {{ t('video.download') }}
              </a>
            </div>
          </div>

          <div class="divide-y divide-gray-100 dark:divide-dark-700">
            <div v-if="!jobs.length" class="flex min-h-[180px] flex-col items-center justify-center px-6 text-center">
              <Icon name="inbox" size="xl" class="mb-3 text-gray-300 dark:text-dark-600" />
              <p class="text-sm font-medium text-gray-700 dark:text-gray-200">{{ t('video.noJobs') }}</p>
              <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('video.noJobsHint') }}</p>
            </div>
            <article v-for="job in jobs" :key="job.job_id" class="flex items-center gap-3 px-5 py-4 transition-colors hover:bg-gray-50 dark:hover:bg-dark-900/40" :class="selectedJob?.job_id === job.job_id ? 'bg-primary-50/60 dark:bg-primary-900/10' : ''" @click="selectedJobId = job.job_id">
              <div class="min-w-0 flex-1">
                <div class="flex min-w-0 items-center gap-2">
                  <span class="truncate text-sm font-medium text-gray-800 dark:text-gray-100">{{ job.prompt || job.job_id }}</span>
                  <span class="badge flex-shrink-0" :class="statusClass(job.status)">{{ statusLabel(job.status) }}</span>
                </div>
                <div class="mt-1 flex flex-wrap items-center gap-x-2 gap-y-1 text-xs text-gray-500 dark:text-dark-400">
                  <span class="font-mono">{{ job.job_id }}</span>
                  <span>{{ job.model || model }}</span>
                  <span>{{ formatDate(job.updated_at || job.created_at) }}</span>
                </div>
                <p v-if="job.status === 'failed' && job.error?.message" class="mt-1 truncate text-xs text-red-600 dark:text-red-300">{{ job.error.message }}</p>
              </div>
              <button v-if="job.status === 'pending'" type="button" class="inline-flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-md text-gray-500 hover:bg-gray-200 hover:text-red-600 dark:text-dark-300 dark:hover:bg-dark-700 dark:hover:text-red-300" :title="t('video.cancel')" :data-testid="`cancel-${job.job_id}`" @click.stop="cancelJob(job)">
                <Icon name="x" size="sm" />
              </button>
              <a v-else-if="videoUrl(job)" :href="videoUrl(job)" target="_blank" rel="noreferrer" class="inline-flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-md text-gray-500 hover:bg-gray-200 hover:text-primary-600 dark:text-dark-300 dark:hover:bg-dark-700 dark:hover:text-primary-300" :title="t('video.open')" @click.stop>
                <Icon name="externalLink" size="sm" />
              </a>
            </article>
          </div>
        </section>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { keysAPI } from '@/api'
import type { ApiKey } from '@/types'
import { useAppStore } from '@/stores/app'
import {
  cancelVideoJob,
  createVideoJob,
  listVideoJobs,
  uploadVideoInput,
  type VideoGenerationRequest,
  type VideoJob,
} from '@/api/videoGeneration'

const { t } = useI18n()
const appStore = useAppStore()
const keys = ref<ApiKey[]>([])
const selectedKeyId = ref<number | null>(null)
const jobs = ref<VideoJob[]>([])
const selectedJobId = ref('')
const prompt = ref('')
const model = ref('seedance-2.0')
const resolution = ref<'480p' | '720p' | '1080p'>('720p')
const duration = ref(8)
const aspectRatio = ref('16:9')
const audio = ref(false)
const imageMode = ref<'none' | 'local' | 'url'>('none')
const imageUrl = ref('')
const localFile = ref<File | null>(null)
const previewUrl = ref('')
const loadingKeys = ref(false)
const loadingJobs = ref(false)
const submitting = ref(false)
const uploading = ref(false)
let pollTimer: ReturnType<typeof setInterval> | null = null

const modelOptions = ['seedance-2.0', 'seedance-2.0-fast']
const aspectRatioOptions = ['16:9', '9:16', '1:1', '4:3', '3:4', '21:9', '9:21']
const imageModes = [
  { value: 'none' as const, label: 'video.imageNone' },
  { value: 'local' as const, label: 'video.imageLocal' },
  { value: 'url' as const, label: 'video.imageRemote' },
]

const leoKeys = computed(() => keys.value.filter((key) => key.status === 'active' && key.group?.platform === 'leo' && key.group?.allow_image_generation === true))
const selectedKey = computed(() => leoKeys.value.find((key) => key.id === selectedKeyId.value) || null)
const activeJobs = computed(() => jobs.value.filter((job) => ['pending', 'running', 'settling'].includes(job.status)))
const selectedJob = computed(() => jobs.value.find((job) => job.job_id === selectedJobId.value) || jobs.value[0] || null)
const selectedVideoUrl = computed(() => (selectedJob.value ? videoUrl(selectedJob.value) : ''))
const canSubmit = computed(() => Boolean(selectedKey.value && prompt.value.trim() && model.value && duration.value >= 4 && duration.value <= 15 && (imageMode.value !== 'url' || /^https?:\/\/\S+$/i.test(imageUrl.value.trim()))))

async function loadKeys() {
  loadingKeys.value = true
  try {
    const result = await keysAPI.list(1, 100, { status: 'active' })
    keys.value = result.items || []
    if (!selectedKey.value) selectedKeyId.value = leoKeys.value[0]?.id || null
  } catch (error) {
    appStore.showError(errorMessage(error))
  } finally {
    loadingKeys.value = false
  }
}

async function loadJobs() {
  if (!selectedKey.value) {
    jobs.value = []
    stopPolling()
    return
  }
  loadingJobs.value = true
  try {
    const result = await listVideoJobs(selectedKey.value.key, { limit: 50 })
    jobs.value = result.data || []
    if (!selectedJobId.value || !jobs.value.some((job) => job.job_id === selectedJobId.value)) selectedJobId.value = jobs.value[0]?.job_id || ''
    updatePolling()
  } catch (error) {
    appStore.showError(errorMessage(error))
  } finally {
    loadingJobs.value = false
  }
}

async function submitJob() {
  if (!canSubmit.value || !selectedKey.value || submitting.value || uploading.value) return
  submitting.value = true
  try {
    let selectedImageUrl = ''
    if (imageMode.value === 'local' && localFile.value) {
      uploading.value = true
      const uploaded = await uploadVideoInput(selectedKey.value.key, localFile.value)
      selectedImageUrl = uploaded.image_url
      uploading.value = false
    } else if (imageMode.value === 'url') {
      selectedImageUrl = imageUrl.value.trim()
    }
    const payload: VideoGenerationRequest = {
      model: model.value,
      prompt: prompt.value.trim(),
      resolution: resolution.value,
      duration: Math.round(duration.value),
      aspect_ratio: aspectRatio.value,
      audio: audio.value,
    }
    if (selectedImageUrl) payload.image_url = selectedImageUrl
    const accepted = await createVideoJob(selectedKey.value.key, payload)
    const now = new Date().toISOString()
    const job: VideoJob = { ...accepted, model: accepted.model || model.value, prompt: accepted.prompt || payload.prompt, created_at: accepted.created_at || now, updated_at: accepted.updated_at || now }
    jobs.value = [job, ...jobs.value.filter((item) => item.job_id !== job.job_id)].slice(0, 50)
    selectedJobId.value = job.job_id
    updatePolling()
    appStore.showSuccess(t('video.submitSuccess'))
    if (imageMode.value === 'local') removeLocalImage()
  } catch (error) {
    appStore.showError(errorMessage(error))
  } finally {
    uploading.value = false
    submitting.value = false
  }
}

async function cancelJob(job: VideoJob) {
  if (!selectedKey.value || job.status !== 'pending') return
  try {
    const canceled = await cancelVideoJob(selectedKey.value.key, job.job_id)
    const index = jobs.value.findIndex((item) => item.job_id === job.job_id)
    if (index >= 0) jobs.value[index] = { ...jobs.value[index], ...canceled }
    updatePolling()
  } catch (error) {
    appStore.showError(errorMessage(error))
  }
}

function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  if (!['image/png', 'image/jpeg', 'image/webp'].includes(file.type)) {
    appStore.showError(t('video.invalidImage'))
    input.value = ''
    return
  }
  removeLocalImage()
  localFile.value = file
  previewUrl.value = URL.createObjectURL(file)
}

function removeLocalImage() {
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
  previewUrl.value = ''
  localFile.value = null
}

function videoUrl(job: VideoJob): string {
  const result = job.result
  const first = result?.data?.[0]
  if (first) {
    for (const key of ['mp4_url', 'video_url', 'url']) {
      if (typeof first[key] === 'string' && first[key]) return first[key] as string
    }
  }
  if (result && typeof result.url === 'string') return result.url
  return ''
}

function statusLabel(status: string) {
  return t(`video.status.${status}`)
}

function statusClass(status: string) {
  return {
    'bg-blue-50 text-blue-700 dark:bg-blue-900/20 dark:text-blue-300': status === 'pending',
    'bg-amber-50 text-amber-700 dark:bg-amber-900/20 dark:text-amber-300': status === 'running' || status === 'settling',
    'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300': status === 'completed',
    'bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-300': status === 'failed',
    'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-dark-300': status === 'canceled',
  }
}

function formatDate(value?: string) {
  if (!value) return ''
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString()
}

function updatePolling() {
  if (activeJobs.value.length && !pollTimer) pollTimer = setInterval(() => void loadJobs(), 2000)
  if (!activeJobs.value.length) stopPolling()
}

function stopPolling() {
  if (pollTimer) clearInterval(pollTimer)
  pollTimer = null
}

function errorMessage(error: unknown) {
  return error instanceof Error && error.message ? error.message : t('common.error')
}

watch(selectedKeyId, () => void loadJobs())
watch(activeJobs, updatePolling)

onMounted(async () => {
  await loadKeys()
  await loadJobs()
})

onBeforeUnmount(() => {
  stopPolling()
  removeLocalImage()
})
</script>
