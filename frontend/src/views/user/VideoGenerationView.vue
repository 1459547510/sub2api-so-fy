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

      <VideoSectionTabs />

      <div class="grid min-w-0 gap-4 xl:grid-cols-[400px_minmax(0,1fr)]">
        <section class="card min-w-0 p-5" data-testid="video-settings">
          <div class="flex items-center justify-between gap-3 border-b border-gray-100 pb-4 dark:border-dark-700">
            <div>
              <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('video.settings') }}</h2>
              <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('video.settingsHint') }}</p>
            </div>
            <Icon name="sparkles" size="lg" class="text-primary-500" />
          </div>

          <form class="mt-5 space-y-4" @submit.prevent="submitJob">
            <div class="space-y-2">
              <span class="input-label">{{ t('video.apiKey') }}</span>
              <div class="grid grid-cols-2 gap-1 rounded-lg bg-gray-100 p-1 dark:bg-dark-800" role="tablist">
                <button type="button" class="rounded-md px-3 py-2 text-xs font-medium transition-colors" :class="apiKeyMode === 'saved' ? 'bg-white text-primary-700 shadow-sm dark:bg-dark-600 dark:text-primary-200' : 'text-gray-500 hover:text-gray-900 dark:text-dark-300 dark:hover:text-white'" :disabled="submitting || uploading" data-testid="key-mode-saved" @click="apiKeyMode = 'saved'">
                  {{ t('video.savedApiKey') }}
                </button>
                <button type="button" class="rounded-md px-3 py-2 text-xs font-medium transition-colors" :class="apiKeyMode === 'custom' ? 'bg-white text-primary-700 shadow-sm dark:bg-dark-600 dark:text-primary-200' : 'text-gray-500 hover:text-gray-900 dark:text-dark-300 dark:hover:text-white'" :disabled="submitting || uploading" data-testid="key-mode-custom" @click="apiKeyMode = 'custom'">
                  {{ t('video.customApiKey') }}
                </button>
              </div>
              <select v-if="apiKeyMode === 'saved'" v-model="selectedKeyId" class="input" :disabled="loadingKeys || !leoKeys.length || submitting || uploading" data-testid="video-api-key">
                <option v-for="key in leoKeys" :key="key.id" :value="key.id">{{ key.name || `API Key #${key.id}` }}</option>
              </select>
              <div v-else class="relative">
                <input v-model="customApiKey" :type="showCustomApiKey ? 'text' : 'password'" class="input pr-11" :placeholder="t('video.customApiKeyPlaceholder')" autocomplete="off" spellcheck="false" :disabled="submitting || uploading" data-testid="video-custom-api-key" />
                <button type="button" class="absolute inset-y-0 right-0 flex items-center pr-3.5 text-gray-400 transition-colors hover:text-gray-600 dark:hover:text-dark-300" :title="showCustomApiKey ? t('video.hideApiKey') : t('video.showApiKey')" :disabled="submitting || uploading" data-testid="toggle-custom-api-key" @click="showCustomApiKey = !showCustomApiKey">
                  <Icon :name="showCustomApiKey ? 'eyeOff' : 'eye'" size="md" />
                </button>
              </div>
            </div>

            <div v-if="apiKeyMode === 'saved' && !leoKeys.length && !loadingKeys" class="rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800 dark:border-amber-900/60 dark:bg-amber-900/20 dark:text-amber-200" data-testid="video-no-key">
              <div class="flex items-start gap-3">
                <Icon name="key" size="md" class="mt-0.5 flex-shrink-0" />
                <div>
                  <p class="font-medium">{{ t('video.noKey') }}</p>
                  <p class="mt-1 text-xs opacity-80">{{ t('video.noKeyHint') }}</p>
                </div>
              </div>
            </div>

            <label class="block">
              <span class="input-label">{{ t('video.prompt') }}</span>
              <textarea
                v-model="prompt"
                rows="5"
                class="input min-h-[128px] resize-y"
                :placeholder="t('video.promptPlaceholder')"
                :maxlength="currentModelCapability?.maxPromptLength"
                data-testid="video-prompt"
                @keydown.meta.enter.prevent="submitJob"
                @keydown.ctrl.enter.prevent="submitJob"
              ></textarea>
            </label>

            <div class="grid grid-cols-2 gap-3">
              <label class="block">
                <span class="input-label">{{ t('video.model') }}</span>
                <select v-model="model" class="input" :disabled="submitting || uploading" data-testid="video-model">
                  <option v-for="option in modelOptions" :key="option" :value="option">{{ option }}</option>
                </select>
              </label>
              <label class="block">
                <span class="input-label">{{ t('video.resolution') }}</span>
                <select v-model="resolution" class="input" :disabled="submitting || uploading" data-testid="video-resolution">
                  <option v-for="option in resolutionOptions" :key="option" :value="option">{{ option }}</option>
                </select>
              </label>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <label class="block">
                <span class="input-label">{{ t('video.aspectRatio') }}</span>
                <select v-model="aspectRatio" class="input" :disabled="submitting || uploading" data-testid="video-aspect-ratio">
                  <option v-for="option in aspectRatioOptions" :key="option" :value="option">{{ option }}</option>
                </select>
              </label>
              <label class="block">
                <span class="input-label">{{ t('video.duration') }}</span>
                <select v-model.number="duration" class="input" :disabled="submitting || uploading" data-testid="video-duration">
                  <option v-for="option in durationOptions" :key="option" :value="option">{{ option }}</option>
                </select>
              </label>
            </div>

            <label class="flex cursor-pointer items-center justify-between rounded-lg border border-gray-200 px-3 py-2.5 dark:border-dark-700">
              <span>
                <span class="block text-sm font-medium text-gray-800 dark:text-gray-100">{{ t('video.audio') }}</span>
                <span class="block text-xs text-gray-500 dark:text-dark-400">{{ t('video.audioHint') }}</span>
              </span>
              <input v-model="audio" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500" :disabled="submitting || uploading" data-testid="video-audio" />
            </label>

            <label v-if="currentModelCapability?.supportsPromptEnhance" class="block">
              <span class="input-label">{{ t('video.promptEnhance') }}</span>
              <select v-model="promptEnhance" class="input" :disabled="submitting || uploading" data-testid="video-prompt-enhance">
                <option value="">{{ t('video.promptEnhanceAuto') }}</option>
                <option value="AUTO">AUTO</option>
                <option value="ON">ON</option>
                <option value="OFF">OFF</option>
              </select>
            </label>

            <div>
              <span class="input-label">{{ t('video.imageInput', { count: maxImageReferences }) }}</span>
              <div class="grid grid-cols-3 gap-1 rounded-lg bg-gray-100 p-1 dark:bg-dark-800" role="tablist">
                <button v-for="option in imageModes" :key="option.value" type="button" class="rounded-md px-2 py-2 text-xs font-medium transition-colors" :class="imageMode === option.value ? 'bg-white text-primary-700 shadow-sm dark:bg-dark-600 dark:text-primary-200' : 'text-gray-500 hover:text-gray-900 dark:text-dark-300 dark:hover:text-white'" :data-testid="`mode-${option.value}`" @click="imageMode = option.value">
                  {{ t(option.label) }}
                </button>
              </div>
            </div>

            <div v-if="imageMode === 'local'" class="space-y-3">
              <label
                class="flex items-center justify-center gap-2 rounded-lg border border-dashed border-gray-300 px-3 py-4 text-sm text-gray-600 transition-colors dark:border-dark-600 dark:text-dark-300"
                :class="hasLocalFrames || maxImageReferences === 0 ? 'cursor-not-allowed opacity-50' : 'cursor-pointer hover:border-primary-400 hover:text-primary-600 dark:hover:border-primary-500 dark:hover:text-primary-300'"
                :title="hasLocalFrames ? t('video.imageInputConflict') : maxImageReferences === 0 ? t('video.modelGuidanceUnsupported') : undefined"
              >
                <Icon name="upload" size="sm" />
                <span>{{ localFiles.length ? t('video.addImage') : t('video.chooseImage') }}</span>
                <input type="file" accept="image/png,image/jpeg,image/webp" class="sr-only" :disabled="hasLocalFrames || maxImageReferences === 0" data-testid="video-image-file" @change="onFileChange" />
              </label>
              <div v-if="previewUrls.length" class="grid grid-cols-2 gap-2">
                <div v-for="(preview, index) in previewUrls" :key="preview" class="relative aspect-video overflow-hidden rounded-lg border border-gray-200 dark:border-dark-700">
                  <img :src="preview" alt="" class="h-full w-full object-cover" />
                  <button type="button" class="absolute right-1.5 top-1.5 inline-flex h-7 w-7 items-center justify-center rounded-md bg-black/60 text-white hover:bg-black/80" :title="t('video.removeImage')" @click="removeLocalImage(index)">
                    <Icon name="x" size="sm" />
                  </button>
                </div>
              </div>
              <div class="grid grid-cols-2 gap-3">
                <label
                  class="relative flex min-h-[92px] flex-col items-center justify-center gap-1 rounded-lg border border-dashed border-gray-300 px-2 py-3 text-center text-xs text-gray-600 transition-colors dark:border-dark-600 dark:text-dark-300"
                  :class="hasLocalReferences ? 'cursor-not-allowed opacity-50' : 'cursor-pointer hover:border-primary-400 hover:text-primary-600 dark:hover:border-primary-500 dark:hover:text-primary-300'"
                  :title="hasLocalReferences ? t('video.imageInputConflict') : undefined"
                >
                  <Icon name="upload" size="sm" class="relative z-10" />
                  <span class="relative z-10">{{ startFrameFile ? t('video.replaceStartFrame') : t('video.chooseStartFrame') }}</span>
                  <input type="file" accept="image/png,image/jpeg,image/webp" class="sr-only" :disabled="hasLocalReferences" data-testid="video-start-frame-file" @change="onStartFrameChange" />
                  <img v-if="startFramePreviewUrl" :src="startFramePreviewUrl" alt="" class="pointer-events-none absolute inset-1 h-[calc(100%-8px)] w-[calc(100%-8px)] rounded-md object-cover opacity-25" />
                </label>
                <label
                  class="relative flex min-h-[92px] flex-col items-center justify-center gap-1 rounded-lg border border-dashed border-gray-300 px-2 py-3 text-center text-xs text-gray-600 transition-colors dark:border-dark-600 dark:text-dark-300"
                  :class="hasLocalReferences ? 'cursor-not-allowed opacity-50' : 'cursor-pointer hover:border-primary-400 hover:text-primary-600 dark:hover:border-primary-500 dark:hover:text-primary-300'"
                  :title="hasLocalReferences ? t('video.imageInputConflict') : undefined"
                >
                  <Icon name="upload" size="sm" class="relative z-10" />
                  <span class="relative z-10">{{ endFrameFile ? t('video.replaceEndFrame') : t('video.chooseEndFrame') }}</span>
                  <input type="file" accept="image/png,image/jpeg,image/webp" class="sr-only" :disabled="hasLocalReferences || currentModelCapability?.maxEndFrames === 0" data-testid="video-end-frame-file" @change="onEndFrameChange" />
                  <img v-if="endFramePreviewUrl" :src="endFramePreviewUrl" alt="" class="pointer-events-none absolute inset-1 h-[calc(100%-8px)] w-[calc(100%-8px)] rounded-md object-cover opacity-25" />
                </label>
              </div>
              <div v-if="startFrameFile || endFrameFile" class="flex flex-wrap gap-2 text-xs text-gray-500 dark:text-dark-400">
                <button v-if="startFrameFile" type="button" class="inline-flex items-center gap-1 text-primary-600 hover:text-primary-700 dark:text-primary-300" data-testid="remove-start-frame" @click="removeStartFrame">
                  <Icon name="x" size="sm" /> {{ t('video.removeStartFrame') }}
                </button>
                <button v-if="endFrameFile" type="button" class="inline-flex items-center gap-1 text-primary-600 hover:text-primary-700 dark:text-primary-300" data-testid="remove-end-frame" @click="removeEndFrame">
                  <Icon name="x" size="sm" /> {{ t('video.removeEndFrame') }}
                </button>
              </div>
            </div>

            <label v-else-if="imageMode === 'url'" class="block" :title="hasRemoteFrames ? t('video.imageInputConflict') : undefined">
              <span class="sr-only">{{ t('video.imageUrl') }}</span>
              <textarea v-model="imageUrlText" rows="4" class="input resize-y disabled:cursor-not-allowed disabled:opacity-50" :placeholder="t('video.imageUrlPlaceholder')" :disabled="hasRemoteFrames || maxImageReferences === 0" data-testid="video-image-url"></textarea>
            </label>

            <div v-if="imageMode === 'url'" class="grid grid-cols-2 gap-3">
              <label class="block" :title="hasRemoteReferences ? t('video.imageInputConflict') : undefined">
                <span class="input-label">{{ t('video.startFrameUrl') }}</span>
                <input v-model="startFrameUrlText" type="url" class="input disabled:cursor-not-allowed disabled:opacity-50" :placeholder="t('video.frameUrlPlaceholder')" :disabled="hasRemoteReferences" data-testid="video-start-frame-url" />
              </label>
              <label class="block" :title="hasRemoteReferences ? t('video.imageInputConflict') : undefined">
                <span class="input-label">{{ t('video.endFrameUrl') }}</span>
                <input v-model="endFrameUrlText" type="url" class="input disabled:cursor-not-allowed disabled:opacity-50" :placeholder="t('video.frameUrlPlaceholder')" :disabled="hasRemoteReferences || currentModelCapability?.maxEndFrames === 0" data-testid="video-end-frame-url" />
              </label>
            </div>

            <p v-if="imageMode !== 'none'" class="text-xs text-gray-500 dark:text-dark-400">{{ t('video.imageInputExclusive') }}</p>

            <div class="space-y-3 border-t border-gray-100 pt-4 dark:border-dark-700">
              <div>
                <span class="input-label">{{ t('video.videoReference') }}</span>
                <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('video.videoReferenceHint') }}</p>
              </div>
              <label
                class="flex items-center justify-center gap-2 rounded-lg border border-dashed border-gray-300 px-3 py-4 text-sm text-gray-600 transition-colors dark:border-dark-600 dark:text-dark-300"
                :class="referenceVideoFiles.length >= (currentModelCapability?.maxVideoRefs || 0) ? 'cursor-not-allowed opacity-50' : 'cursor-pointer hover:border-primary-400 hover:text-primary-600 dark:hover:border-primary-500 dark:hover:text-primary-300'"
                :title="referenceVideoFiles.length >= (currentModelCapability?.maxVideoRefs || 0) ? t('video.modelGuidanceUnsupported') : undefined"
              >
                <Icon name="upload" size="sm" />
                <span>{{ t('video.chooseVideoReference') }}</span>
                <input type="file" accept=".mp4,.mov,video/mp4,video/quicktime" multiple class="sr-only" :disabled="referenceVideoFiles.length >= (currentModelCapability?.maxVideoRefs || 0) || submitting || uploading" data-testid="video-reference-file" @change="onReferenceVideoChange" />
              </label>
              <div v-if="referenceVideoPreviewUrls.length" class="grid gap-2 sm:grid-cols-3">
                <div v-for="(preview, index) in referenceVideoPreviewUrls" :key="preview" class="relative overflow-hidden rounded-lg border border-gray-200 bg-black dark:border-dark-700">
                  <video :src="preview" controls muted playsinline class="aspect-video w-full object-cover"></video>
                  <button type="button" class="absolute right-1.5 top-1.5 inline-flex h-7 w-7 items-center justify-center rounded-md bg-black/60 text-white hover:bg-black/80" :title="t('video.removeVideoReference')" @click="removeReferenceVideo(index)">
                    <Icon name="x" size="sm" />
                  </button>
                </div>
              </div>

              <div>
                <span class="input-label">{{ t('video.audioReference') }}</span>
                <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('video.audioReferenceHint') }}</p>
              </div>
              <label class="flex items-center justify-center gap-2 rounded-lg border border-dashed border-gray-300 px-3 py-4 text-sm text-gray-600 transition-colors dark:border-dark-600 dark:text-dark-300" :class="referenceAudioFile ? 'cursor-pointer hover:border-primary-400 hover:text-primary-600 dark:hover:border-primary-500 dark:hover:text-primary-300' : 'cursor-pointer hover:border-primary-400 hover:text-primary-600 dark:hover:border-primary-500 dark:hover:text-primary-300'">
                <Icon name="upload" size="sm" />
                <span>{{ referenceAudioFile ? t('video.replaceAudioReference') : t('video.chooseAudioReference') }}</span>
                <input type="file" accept=".mp3,.wav,audio/mpeg,audio/wav,audio/x-wav" class="sr-only" :disabled="currentModelCapability?.maxAudioRefs === 0 || submitting || uploading" data-testid="audio-reference-file" @change="onReferenceAudioChange" />
              </label>
              <div v-if="referenceAudioPreviewUrl" class="flex items-center gap-2 rounded-lg border border-gray-200 p-2 dark:border-dark-700">
                <audio :src="referenceAudioPreviewUrl" controls class="min-w-0 flex-1" data-testid="audio-reference-preview"></audio>
                <button type="button" class="inline-flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-md text-gray-500 hover:bg-gray-100 hover:text-red-600 dark:text-dark-300 dark:hover:bg-dark-800 dark:hover:text-red-300" :title="t('video.removeAudioReference')" @click="removeReferenceAudio">
                  <Icon name="x" size="sm" />
                </button>
              </div>
            </div>

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
              <button type="button" class="btn btn-secondary btn-sm" :disabled="loadingJobs || !effectiveApiKey" :title="t('common.refresh')" data-testid="refresh-video-jobs" @click="loadJobs">
                <Icon name="refresh" size="sm" :class="loadingJobs ? 'animate-spin' : ''" />
              </button>
            </div>
          </div>

          <div class="border-b border-gray-100 bg-gray-50/70 p-5 dark:border-dark-700 dark:bg-dark-900/40">
            <div v-if="selectedVideoUrl && selectedVideoKey === currentVideoKey" class="overflow-hidden rounded-lg bg-black">
              <video :src="selectedVideoUrl" controls playsinline class="aspect-video max-h-[440px] w-full" data-testid="video-preview" @error="onVideoError"></video>
            </div>
            <div v-else class="flex min-h-[230px] flex-col items-center justify-center rounded-lg border border-dashed border-gray-200 bg-white px-6 text-center dark:border-dark-700 dark:bg-dark-900/60">
              <Icon name="play" size="xl" class="mb-3 text-gray-300 dark:text-dark-600" />
              <p class="text-sm font-medium text-gray-700 dark:text-gray-200">{{ t('video.previewEmpty') }}</p>
              <p class="mt-1 max-w-sm text-xs text-gray-500 dark:text-dark-400">{{ t('video.previewEmptyHint') }}</p>
            </div>
            <div v-if="videoPreviewError" class="mt-2 flex items-center justify-between gap-3 text-xs text-red-600 dark:text-red-300" data-testid="video-preview-error">
              <p>{{ videoPreviewError }}</p>
              <button type="button" class="font-medium hover:text-red-700 dark:hover:text-red-200" data-testid="retry-video-preview" @click="retrySelectedVideo">{{ t('video.retryPreview') }}</button>
            </div>
            <div v-if="selectedJob" class="mt-3 flex flex-wrap items-center justify-between gap-2 text-xs text-gray-500 dark:text-dark-400">
              <span class="truncate">{{ selectedJob.prompt }}</span>
              <a v-if="selectedVideoUrl && selectedVideoKey === currentVideoKey" :href="selectedVideoUrl" :download="`${selectedJob.job_id}.mp4`" class="inline-flex items-center gap-1 font-medium text-primary-600 hover:text-primary-700 dark:text-primary-300" :title="t('video.download')" data-testid="download-video-output">
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
              <span v-else-if="job.status === 'completed'" class="inline-flex h-8 w-8 flex-shrink-0 items-center justify-center text-primary-600 dark:text-primary-300">
                <Icon name="play" size="sm" />
              </span>
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
import VideoSectionTabs from '@/components/video/VideoSectionTabs.vue'
import { keysAPI } from '@/api'
import type { ApiKey } from '@/types'
import { useAppStore } from '@/stores/app'
import {
  cancelVideoJob,
  createVideoJob,
  downloadVideoOutput,
  listVideoJobs,
  uploadVideoInput,
  type VideoGenerationRequest,
  type VideoJob,
} from '@/api/videoGeneration'

const { t } = useI18n()
const appStore = useAppStore()
const keys = ref<ApiKey[]>([])
const selectedKeyId = ref<number | null>(null)
const apiKeyMode = ref<'saved' | 'custom'>('saved')
const customApiKey = ref('')
const showCustomApiKey = ref(false)
const jobs = ref<VideoJob[]>([])
const selectedJobId = ref('')
const prompt = ref('')
const model = ref('seedance-2.0')
const resolution = ref<VideoResolution>('720p')
const duration = ref(8)
const aspectRatio = ref<VideoAspectRatio>('16:9')
const audio = ref(false)
const promptEnhance = ref<'AUTO' | 'ON' | 'OFF' | ''>('')
const imageMode = ref<'none' | 'local' | 'url'>('none')
const imageUrlText = ref('')
const localFiles = ref<File[]>([])
const previewUrls = ref<string[]>([])
const startFrameFile = ref<File | null>(null)
const endFrameFile = ref<File | null>(null)
const startFramePreviewUrl = ref('')
const endFramePreviewUrl = ref('')
const startFrameUrlText = ref('')
const endFrameUrlText = ref('')
const referenceVideoFiles = ref<File[]>([])
const referenceVideoPreviewUrls = ref<string[]>([])
const referenceAudioFile = ref<File | null>(null)
const referenceAudioPreviewUrl = ref('')
const loadingKeys = ref(false)
const loadingJobs = ref(false)
const submitting = ref(false)
const uploading = ref(false)
const selectedVideoUrl = ref('')
const selectedVideoKey = ref('')
const videoPreviewError = ref('')
let pollTimer: ReturnType<typeof setInterval> | null = null
let videoOutputRequest = 0
let jobsRequest = 0

type VideoResolution = '400p' | '480p' | '544p' | '720p' | '960p' | '1080p' | '1440p' | '2160p'
type VideoAspectRatio = '16:9' | '9:16' | '1:1' | '4:3' | '3:4' | '21:9' | '9:21'

interface VideoModelCapability {
  resolutions: readonly VideoResolution[]
  defaultResolution: VideoResolution
  durations: readonly number[]
  maxDurationByResolution?: Partial<Record<VideoResolution, number>>
  defaultDuration: number
  aspectsByResolution: Partial<Record<VideoResolution, readonly VideoAspectRatio[]>>
  defaultAspectRatio: VideoAspectRatio
  maxPromptLength: number
  maxStartFrames: number
  maxEndFrames: number
  maxImageRefs: number
  maxVideoRefs: number
  maxAudioRefs: number
  requiresStartFrame?: boolean
  supportsPromptEnhance?: boolean
  rejectsPromptEnhanceOnStartFrame?: boolean
}

const allDurationOptions = [4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
const shortVideoDurationOptions = [3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
const ltxProDurationOptions = [6, 8, 10]
const ltxFastDurationOptions = [6, 8, 10, 12, 14, 16, 18, 20]
const allAspectRatioOptions: readonly VideoAspectRatio[] = ['16:9', '9:16', '1:1', '4:3', '3:4', '21:9', '9:21']
const hdAspectRatioOptions: readonly VideoAspectRatio[] = ['16:9', '9:16', '1:1', '4:3', '3:4', '21:9']
const videoReferenceMaxBytes = 100 * 1024 * 1024
const audioReferenceMaxBytes = 15 * 1024 * 1024
const imageReferenceMaxBytes = 10 * 1024 * 1024
const videoModelCapabilities: Record<string, VideoModelCapability> = {
  'seedance-2.0': {
    resolutions: ['480p', '720p', '1080p'],
    defaultResolution: '720p',
    durations: allDurationOptions,
    maxDurationByResolution: { '1080p': 12 },
    defaultDuration: 8,
    aspectsByResolution: {
      '480p': allAspectRatioOptions,
      '720p': hdAspectRatioOptions,
      '1080p': allAspectRatioOptions,
    },
    defaultAspectRatio: '16:9',
    maxPromptLength: 5000,
    maxStartFrames: 1,
    maxEndFrames: 1,
    maxImageRefs: 4,
    maxVideoRefs: 3,
    maxAudioRefs: 1,
  },
  'seedance-2.0-fast': {
    resolutions: ['480p', '720p'],
    defaultResolution: '720p',
    durations: allDurationOptions,
    defaultDuration: 8,
    aspectsByResolution: {
      '480p': allAspectRatioOptions,
      '720p': hdAspectRatioOptions,
    },
    defaultAspectRatio: '16:9',
    maxPromptLength: 5000,
    maxStartFrames: 1,
    maxEndFrames: 1,
    maxImageRefs: 4,
    maxVideoRefs: 3,
    maxAudioRefs: 1,
  },
  'seedance-2.0-mini': {
    resolutions: ['480p', '720p'],
    defaultResolution: '720p',
    durations: allDurationOptions,
    defaultDuration: 8,
    aspectsByResolution: {
      '480p': ['16:9', '1:1', '9:16'],
      '720p': ['16:9', '1:1', '9:16'],
    },
    defaultAspectRatio: '16:9',
    maxPromptLength: 5000,
    maxStartFrames: 1,
    maxEndFrames: 1,
    maxImageRefs: 4,
    maxVideoRefs: 3,
    maxAudioRefs: 1,
  },
  'happy-horse-1.1': {
    resolutions: ['720p', '1080p'],
    defaultResolution: '1080p',
    durations: shortVideoDurationOptions,
    defaultDuration: 5,
    aspectsByResolution: {
      '720p': ['16:9', '4:3', '1:1', '3:4', '9:16'],
      '1080p': ['16:9', '4:3', '1:1', '3:4', '9:16'],
    },
    defaultAspectRatio: '16:9',
    maxPromptLength: 2500,
    maxStartFrames: 1,
    maxEndFrames: 0,
    maxImageRefs: 9,
    maxVideoRefs: 0,
    maxAudioRefs: 0,
    supportsPromptEnhance: true,
  },
  'grok-imagine-1.5': {
    resolutions: ['400p', '544p', '720p', '960p'],
    defaultResolution: '720p',
    durations: shortVideoDurationOptions,
    defaultDuration: 6,
    aspectsByResolution: {
      '400p': ['16:9', '9:16'],
      '544p': ['1:1'],
      '720p': ['16:9', '9:16'],
      '960p': ['1:1'],
    },
    defaultAspectRatio: '16:9',
    maxPromptLength: 5000,
    maxStartFrames: 1,
    maxEndFrames: 0,
    maxImageRefs: 0,
    maxVideoRefs: 0,
    maxAudioRefs: 0,
    requiresStartFrame: true,
  },
  'ltxv-2.3-pro': {
    resolutions: ['1080p', '1440p', '2160p'],
    defaultResolution: '1080p',
    durations: ltxProDurationOptions,
    defaultDuration: 6,
    aspectsByResolution: {
      '1080p': ['16:9'],
      '1440p': ['16:9'],
      '2160p': ['16:9'],
    },
    defaultAspectRatio: '16:9',
    maxPromptLength: 5000,
    maxStartFrames: 1,
    maxEndFrames: 1,
    maxImageRefs: 0,
    maxVideoRefs: 0,
    maxAudioRefs: 0,
    supportsPromptEnhance: true,
  },
  'ltxv-2.3-fast': {
    resolutions: ['1080p', '1440p', '2160p'],
    defaultResolution: '1080p',
    durations: ltxFastDurationOptions,
    defaultDuration: 6,
    aspectsByResolution: {
      '1080p': ['16:9'],
      '1440p': ['16:9'],
      '2160p': ['16:9'],
    },
    defaultAspectRatio: '16:9',
    maxPromptLength: 5000,
    maxStartFrames: 1,
    maxEndFrames: 1,
    maxImageRefs: 0,
    maxVideoRefs: 0,
    maxAudioRefs: 0,
    supportsPromptEnhance: true,
  },
}
const modelOptions = Object.keys(videoModelCapabilities)
const imageModes = [
  { value: 'none' as const, label: 'video.imageNone' },
  { value: 'local' as const, label: 'video.imageLocal' },
  { value: 'url' as const, label: 'video.imageRemote' },
]

const leoKeys = computed(() => keys.value.filter((key) => key.status === 'active' && key.group?.platform === 'leo' && key.group?.allow_image_generation === true))
const selectedKey = computed(() => leoKeys.value.find((key) => key.id === selectedKeyId.value) || null)
const effectiveApiKey = computed(() => apiKeyMode.value === 'custom' ? customApiKey.value.trim() : selectedKey.value?.key || '')
const currentModelCapability = computed(() => videoModelCapabilities[model.value])
const resolutionOptions = computed(() => currentModelCapability.value?.resolutions || [])
const durationOptions = computed(() => supportedDurations(currentModelCapability.value, resolution.value))
const aspectRatioOptions = computed(() => currentModelCapability.value?.aspectsByResolution[resolution.value] || [])
const maxImageReferences = computed(() => currentModelCapability.value?.maxImageRefs || 0)
const activeJobs = computed(() => jobs.value.filter((job) => ['pending', 'running', 'settling'].includes(job.status)))
const selectedJob = computed(() => jobs.value.find((job) => job.job_id === selectedJobId.value) || jobs.value[0] || null)
const remoteImageUrls = computed(() => imageUrlText.value.split(/[\r\n,]+/).map((value) => value.trim()).filter(Boolean))
const remoteStartFrameUrl = computed(() => startFrameUrlText.value.trim())
const remoteEndFrameUrl = computed(() => endFrameUrlText.value.trim())
const hasLocalReferences = computed(() => localFiles.value.length > 0)
const hasLocalFrames = computed(() => Boolean(startFrameFile.value || endFrameFile.value))
const hasRemoteReferences = computed(() => remoteImageUrls.value.length > 0)
const hasRemoteFrames = computed(() => Boolean(remoteStartFrameUrl.value || remoteEndFrameUrl.value))
const hasAudioVisualReference = computed(() => Boolean(
  referenceVideoFiles.value.length || hasLocalReferences.value || hasRemoteReferences.value,
))
const hasMixedImageInputs = computed(() => (
  imageMode.value === 'local'
    ? hasLocalReferences.value && hasLocalFrames.value
    : imageMode.value === 'url' && hasRemoteReferences.value && hasRemoteFrames.value
))
const hasUnsupportedModelGuidances = computed(() => {
  const capability = currentModelCapability.value
  if (!capability) return true
  const startCount = Number(Boolean(startFrameFile.value || remoteStartFrameUrl.value))
  const endCount = Number(Boolean(endFrameFile.value || remoteEndFrameUrl.value))
  const imageCount = imageMode.value === 'local' ? localFiles.value.length : imageMode.value === 'url' ? remoteImageUrls.value.length : 0
  return startCount > capability.maxStartFrames ||
    endCount > capability.maxEndFrames ||
    imageCount > capability.maxImageRefs ||
    referenceVideoFiles.value.length > capability.maxVideoRefs ||
    Number(Boolean(referenceAudioFile.value)) > capability.maxAudioRefs ||
    (Boolean(promptEnhance.value) && !capability.supportsPromptEnhance) ||
    (capability.rejectsPromptEnhanceOnStartFrame && promptEnhance.value === 'ON' && startCount > 0) ||
    Boolean(capability.requiresStartFrame && !startCount)
})
const currentVideoKey = computed(() => {
  const job = selectedJob.value
  return job && effectiveApiKey.value ? `${effectiveApiKey.value}\u0000${job.job_id}\u0000${job.status}` : ''
})
const hasValidModelParameters = computed(() => supportsModelParameters(model.value, resolution.value, duration.value, aspectRatio.value))
const canSubmit = computed(() => Boolean(
  effectiveApiKey.value && prompt.value.trim() && hasValidModelParameters.value &&
  !hasMixedImageInputs.value &&
  (!referenceAudioFile.value || hasAudioVisualReference.value) &&
  (imageMode.value !== 'url' || (
    (remoteImageUrls.value.length > 0 || remoteStartFrameUrl.value || remoteEndFrameUrl.value) &&
    remoteImageUrls.value.length <= maxImageReferences.value &&
    remoteImageUrls.value.every((value) => /^https?:\/\/\S+$/i.test(value)) &&
    (!remoteStartFrameUrl.value || /^https?:\/\/\S+$/i.test(remoteStartFrameUrl.value)) &&
    (!remoteEndFrameUrl.value || /^https?:\/\/\S+$/i.test(remoteEndFrameUrl.value))
  )) && !hasUnsupportedModelGuidances.value
))

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
  const apiKey = effectiveApiKey.value
  if (!apiKey) {
    jobs.value = []
    stopPolling()
    return
  }
  const request = ++jobsRequest
  loadingJobs.value = true
  try {
    const result = await listVideoJobs(apiKey, { limit: 50 })
    if (request !== jobsRequest || apiKey !== effectiveApiKey.value) return
    jobs.value = result.data || []
    if (!selectedJobId.value || !jobs.value.some((job) => job.job_id === selectedJobId.value)) selectedJobId.value = jobs.value[0]?.job_id || ''
    updatePolling()
  } catch (error) {
    if (request === jobsRequest) appStore.showError(errorMessage(error))
  } finally {
    if (request === jobsRequest) loadingJobs.value = false
  }
}

async function submitJob() {
  const apiKey = effectiveApiKey.value
  const requestParameters = {
    model: model.value,
    resolution: resolution.value,
    duration: duration.value,
    aspectRatio: aspectRatio.value,
    audio: audio.value,
  }
  if (!supportsModelParameters(requestParameters.model, requestParameters.resolution, requestParameters.duration, requestParameters.aspectRatio)) {
    appStore.showError(t('video.invalidModelParameters'))
    return
  }
  if (hasMixedImageInputs.value) {
    appStore.showError(t('video.imageInputConflict'))
    return
  }
  if (hasUnsupportedModelGuidances.value) {
    appStore.showError(t('video.modelGuidanceUnsupported'))
    return
  }
  if (referenceAudioFile.value && !hasAudioVisualReference.value) {
    appStore.showError(t('video.audioNeedsVisualReference'))
    return
  }
  if (!canSubmit.value || !apiKey || submitting.value || uploading.value) return
  submitting.value = true
  try {
    let selectedImageUrls: string[] = []
    let startFrameUrl = imageMode.value === 'url' ? remoteStartFrameUrl.value : ''
    let endFrameUrl = imageMode.value === 'url' ? remoteEndFrameUrl.value : ''
    let selectedVideoUrls: string[] = []
    let selectedAudioUrl = ''
    const referenceFiles = [...localFiles.value]
    const startFile = startFrameFile.value
    const endFile = endFrameFile.value
    if (imageMode.value === 'local' && (referenceFiles.length || startFile || endFile)) {
      uploading.value = true
      const uploadFiles = [
        ...referenceFiles,
        ...(startFile ? [startFile] : []),
        ...(endFile ? [endFile] : []),
      ]
      const uploaded = await Promise.all(uploadFiles.map((file) => uploadVideoInput(apiKey, file)))
      selectedImageUrls = uploaded.slice(0, referenceFiles.length).map(uploadedMediaUrl).filter(Boolean)
      const startUploadIndex = referenceFiles.length
      if (startFile) startFrameUrl = uploadedMediaUrl(uploaded[startUploadIndex])
      if (endFile) endFrameUrl = uploadedMediaUrl(uploaded[startUploadIndex + (startFile ? 1 : 0)])
      uploading.value = false
    } else if (imageMode.value === 'url') {
      selectedImageUrls = remoteImageUrls.value
    }
    if (referenceVideoFiles.value.length || referenceAudioFile.value) {
      uploading.value = true
      const [uploadedVideos, uploadedAudio] = await Promise.all([
        Promise.all(referenceVideoFiles.value.map((file) => uploadVideoInput(apiKey, file, 'video'))),
        referenceAudioFile.value ? uploadVideoInput(apiKey, referenceAudioFile.value, 'audio') : Promise.resolve(null),
      ])
      selectedVideoUrls = uploadedVideos.map(uploadedMediaUrl).filter(Boolean)
      selectedAudioUrl = uploadedMediaUrl(uploadedAudio || undefined)
      uploading.value = false
    }
    const payload: VideoGenerationRequest = {
      model: requestParameters.model,
      prompt: prompt.value.trim(),
      resolution: requestParameters.resolution,
      duration: requestParameters.duration,
      aspect_ratio: requestParameters.aspectRatio,
      audio: requestParameters.audio,
    }
    if (promptEnhance.value) payload.prompt_enhance = promptEnhance.value
    if (startFrameUrl) payload.start_frame_url = startFrameUrl
    if (endFrameUrl) payload.end_frame_url = endFrameUrl
    if (selectedImageUrls.length === 1 && !startFrameUrl && !endFrameUrl && !selectedAudioUrl) {
      payload.image_url = selectedImageUrls[0]
    } else if (selectedImageUrls.length) {
      payload.guidances = {
        image_reference: selectedImageUrls.map((url, order) => ({
          image: { url, type: 'UPLOADED' },
          strength: 'MID',
          order,
        })),
      }
    }
    if (selectedVideoUrls.length) {
      payload.guidances = {
        ...payload.guidances,
        video_reference_base: selectedVideoUrls.map((url) => ({ video: { url, type: 'UPLOADED' } })),
      }
    }
    if (selectedAudioUrl) {
      payload.guidances = {
        ...payload.guidances,
        audio_reference: [{ audio: { url: selectedAudioUrl, type: 'UPLOADED' } }],
      }
    }
    const accepted = await createVideoJob(apiKey, payload)
    const now = new Date().toISOString()
    const job: VideoJob = { ...accepted, model: accepted.model || requestParameters.model, prompt: accepted.prompt || payload.prompt, created_at: accepted.created_at || now, updated_at: accepted.updated_at || now }
    jobs.value = [job, ...jobs.value.filter((item) => item.job_id !== job.job_id)].slice(0, 50)
    selectedJobId.value = job.job_id
    updatePolling()
    appStore.showSuccess(t('video.submitSuccess'))
    if (imageMode.value === 'local') clearImageInputs()
    clearMediaInputs()
  } catch (error) {
    appStore.showError(errorMessage(error))
  } finally {
    uploading.value = false
    submitting.value = false
  }
}

async function cancelJob(job: VideoJob) {
  const apiKey = effectiveApiKey.value
  if (!apiKey || job.status !== 'pending') return
  try {
    const canceled = await cancelVideoJob(apiKey, job.job_id)
    if (apiKey !== effectiveApiKey.value) return
    const index = jobs.value.findIndex((item) => item.job_id === job.job_id)
    if (index >= 0) jobs.value[index] = { ...jobs.value[index], ...canceled }
    updatePolling()
  } catch (error) {
    appStore.showError(errorMessage(error))
  }
}

function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const files = Array.from(input.files || [])
  input.value = ''
  if (!files.length) return
  if (hasLocalFrames.value) {
    appStore.showError(t('video.imageInputConflict'))
    return
  }
  if (localFiles.value.length >= maxImageReferences.value) {
    appStore.showError(t('video.tooManyImages', { count: maxImageReferences.value }))
    return
  }
  const file = files[0]
  if (file.size > imageReferenceMaxBytes) {
    appStore.showError(t('video.imageTooLarge'))
    return
  }
  if (!['image/png', 'image/jpeg', 'image/webp'].includes(file.type)) {
    appStore.showError(t('video.invalidImage'))
    return
  }
  localFiles.value = [...localFiles.value, file]
  previewUrls.value = [...previewUrls.value, URL.createObjectURL(file)]
}

function onStartFrameChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  if (hasLocalReferences.value) {
    appStore.showError(t('video.imageInputConflict'))
    return
  }
  if (file.size > imageReferenceMaxBytes) {
    appStore.showError(t('video.imageTooLarge'))
    return
  }
  if (!['image/png', 'image/jpeg', 'image/webp'].includes(file.type)) {
    appStore.showError(t('video.invalidImage'))
    return
  }
  removeStartFrame()
  startFrameFile.value = file
  startFramePreviewUrl.value = URL.createObjectURL(file)
}

function onEndFrameChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  if (currentModelCapability.value?.maxEndFrames === 0) {
    appStore.showError(t('video.modelGuidanceUnsupported'))
    return
  }
  if (hasLocalReferences.value) {
    appStore.showError(t('video.imageInputConflict'))
    return
  }
  if (!['image/png', 'image/jpeg', 'image/webp'].includes(file.type)) {
    appStore.showError(t('video.invalidImage'))
    return
  }
  if (file.size > imageReferenceMaxBytes) {
    appStore.showError(t('video.imageTooLarge'))
    return
  }
  removeEndFrame()
  endFrameFile.value = file
  endFramePreviewUrl.value = URL.createObjectURL(file)
}

async function onReferenceVideoChange(event: Event) {
  const input = event.target as HTMLInputElement
  const files = Array.from(input.files || [])
  input.value = ''
  if (!files.length) return
  if (currentModelCapability.value?.maxVideoRefs === 0) {
    appStore.showError(t('video.modelGuidanceUnsupported'))
    return
  }
  const remaining = (currentModelCapability.value?.maxVideoRefs || 0) - referenceVideoFiles.value.length
  if (remaining <= 0) {
    appStore.showError(t('video.tooManyVideoReferences'))
    return
  }
  for (const file of files.slice(0, remaining)) {
    if (file.size > videoReferenceMaxBytes) {
      appStore.showError(t('video.videoTooLarge'))
      continue
    }
    if (!isSupportedVideoFile(file) || !await canReadVideoFile(file)) {
      appStore.showError(t('video.invalidVideo'))
      continue
    }
    referenceVideoFiles.value = [...referenceVideoFiles.value, file]
    referenceVideoPreviewUrls.value = [...referenceVideoPreviewUrls.value, URL.createObjectURL(file)]
  }
  if (files.length > remaining) appStore.showError(t('video.tooManyVideoReferences'))
}

async function onReferenceAudioChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  if (currentModelCapability.value?.maxAudioRefs === 0) {
    appStore.showError(t('video.modelGuidanceUnsupported'))
    return
  }
  if (!isSupportedAudioFile(file)) {
    appStore.showError(t('video.invalidAudio'))
    return
  }
  if (file.size > audioReferenceMaxBytes) {
    appStore.showError(t('video.audioTooLarge'))
    return
  }
  if (file.name.toLowerCase().endsWith('.wav') && !await isSupportedPCMWave(file)) {
    appStore.showError(t('video.invalidAudio'))
    return
  }
  const durationSeconds = await readAudioDuration(file)
  if (durationSeconds !== null && (durationSeconds < 2 || durationSeconds > 30)) {
    appStore.showError(t('video.audioDuration'))
    return
  }
  removeReferenceAudio()
  referenceAudioFile.value = file
  referenceAudioPreviewUrl.value = URL.createObjectURL(file)
  if (!hasAudioVisualReference.value) appStore.showError(t('video.audioNeedsVisualReference'))
}

function removeReferenceVideo(index: number) {
  const preview = referenceVideoPreviewUrls.value[index]
  if (preview) URL.revokeObjectURL(preview)
  referenceVideoPreviewUrls.value = referenceVideoPreviewUrls.value.filter((_, itemIndex) => itemIndex !== index)
  referenceVideoFiles.value = referenceVideoFiles.value.filter((_, itemIndex) => itemIndex !== index)
}

function removeReferenceAudio() {
  if (referenceAudioPreviewUrl.value) URL.revokeObjectURL(referenceAudioPreviewUrl.value)
  referenceAudioPreviewUrl.value = ''
  referenceAudioFile.value = null
}

function clearMediaInputs() {
  for (const preview of referenceVideoPreviewUrls.value) URL.revokeObjectURL(preview)
  referenceVideoPreviewUrls.value = []
  referenceVideoFiles.value = []
  removeReferenceAudio()
}

function isSupportedVideoFile(file: File) {
  const extension = file.name.toLowerCase().slice(file.name.lastIndexOf('.'))
  return ['.mp4', '.mov'].includes(extension)
}

function isSupportedAudioFile(file: File) {
  const extension = file.name.toLowerCase().slice(file.name.lastIndexOf('.'))
  return ['.mp3', '.wav'].includes(extension)
}

function canReadVideoFile(file: File): Promise<boolean> {
  if (typeof document === 'undefined') return Promise.resolve(true)
  return new Promise((resolve) => {
    const url = URL.createObjectURL(file)
    const video = document.createElement('video')
    let settled = false
    let timer: number | undefined
    const finish = (value: boolean) => {
      if (settled) return
      settled = true
      if (timer !== undefined) window.clearTimeout(timer)
      URL.revokeObjectURL(url)
      resolve(value)
    }
    video.preload = 'metadata'
    video.onloadedmetadata = () => finish(Number.isFinite(video.duration) || video.duration === Infinity)
    video.onerror = () => finish(false)
    video.src = url
    timer = window.setTimeout(() => finish(false), 3000)
    try {
      video.load()
    } catch {
      finish(false)
    }
  })
}

async function isSupportedPCMWave(file: File): Promise<boolean> {
  try {
    const buffer = await file.arrayBuffer()
    const view = new DataView(buffer)
    if (view.byteLength < 12 || readAscii(view, 0, 4) !== 'RIFF' || readAscii(view, 8, 4) !== 'WAVE') return false
    let offset = 12
    while (offset + 8 <= view.byteLength) {
      const chunkSize = view.getUint32(offset + 4, true)
      const chunkStart = offset + 8
      if (chunkStart + chunkSize > view.byteLength) return false
      if (readAscii(view, offset, 4) === 'fmt ' && chunkSize >= 16) {
        return view.getUint16(chunkStart, true) === 1 && [16, 24].includes(view.getUint16(chunkStart + 14, true))
      }
      offset = chunkStart + chunkSize + (chunkSize % 2)
    }
  } catch {
    return false
  }
  return false
}

function readAudioDuration(file: File): Promise<number | null> {
  if (typeof Audio === 'undefined') return Promise.resolve(null)
  return new Promise((resolve) => {
    const url = URL.createObjectURL(file)
    const audioElement = new Audio()
    let settled = false
    let timer: number | undefined
    const finish = (value: number | null) => {
      if (settled) return
      settled = true
      if (timer !== undefined) window.clearTimeout(timer)
      URL.revokeObjectURL(url)
      resolve(value)
    }
    audioElement.onloadedmetadata = () => finish(Number.isFinite(audioElement.duration) ? audioElement.duration : null)
    audioElement.onerror = () => finish(null)
    audioElement.src = url
    timer = window.setTimeout(() => finish(null), 3000)
  })
}

function readAscii(view: DataView, offset: number, length: number) {
  let value = ''
  for (let index = 0; index < length; index += 1) value += String.fromCharCode(view.getUint8(offset + index))
  return value
}

function uploadedMediaUrl(item?: { media_url?: string; image_url?: string; video_url?: string; audio_url?: string }) {
  return item?.media_url || item?.image_url || item?.video_url || item?.audio_url || ''
}

function removeLocalImage(index: number) {
  const preview = previewUrls.value[index]
  if (preview) URL.revokeObjectURL(preview)
  previewUrls.value = previewUrls.value.filter((_, itemIndex) => itemIndex !== index)
  localFiles.value = localFiles.value.filter((_, itemIndex) => itemIndex !== index)
}

function removeLocalImages() {
  for (const preview of previewUrls.value) URL.revokeObjectURL(preview)
  previewUrls.value = []
  localFiles.value = []
}

function removeStartFrame() {
  if (startFramePreviewUrl.value) URL.revokeObjectURL(startFramePreviewUrl.value)
  startFramePreviewUrl.value = ''
  startFrameFile.value = null
}

function removeEndFrame() {
  if (endFramePreviewUrl.value) URL.revokeObjectURL(endFramePreviewUrl.value)
  endFramePreviewUrl.value = ''
  endFrameFile.value = null
}

function clearImageInputs() {
  removeLocalImages()
  removeStartFrame()
  removeEndFrame()
}

function clearSelectedVideo() {
  if (selectedVideoUrl.value) URL.revokeObjectURL(selectedVideoUrl.value)
  selectedVideoUrl.value = ''
  selectedVideoKey.value = ''
  videoPreviewError.value = ''
}

async function loadSelectedVideo() {
  const job = selectedJob.value
  const apiKey = effectiveApiKey.value
  const key = job && apiKey ? `${apiKey}\u0000${job.job_id}\u0000${job.status}` : ''
  const request = ++videoOutputRequest
  if (!job || !apiKey || job.status !== 'completed') {
    clearSelectedVideo()
    return
  }
  if (selectedVideoKey.value === key && selectedVideoUrl.value) return
  videoPreviewError.value = ''
  try {
    const blob = await downloadVideoOutput(apiKey, job.job_id)
    if (request !== videoOutputRequest) return
    const nextUrl = URL.createObjectURL(blob)
    const previousUrl = selectedVideoUrl.value
    selectedVideoUrl.value = nextUrl
    selectedVideoKey.value = key
    if (previousUrl) URL.revokeObjectURL(previousUrl)
  } catch (error) {
    if (request === videoOutputRequest) {
      videoPreviewError.value = errorMessage(error)
      appStore.showError(videoPreviewError.value)
    }
  }
}

function onVideoError() {
  videoPreviewError.value = t('video.previewError')
}

function retrySelectedVideo() {
  const previousUrl = selectedVideoUrl.value
  selectedVideoUrl.value = ''
  selectedVideoKey.value = ''
  videoPreviewError.value = ''
  if (previousUrl) URL.revokeObjectURL(previousUrl)
  void loadSelectedVideo()
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

function resetKeyScopedState() {
  jobsRequest++
  videoOutputRequest++
  loadingJobs.value = false
  jobs.value = []
  selectedJobId.value = ''
  stopPolling()
  clearSelectedVideo()
}

function errorMessage(error: unknown) {
  return error instanceof Error && error.message ? error.message : t('common.error')
}

function supportedDurations(capability: VideoModelCapability | undefined, resolutionValue: VideoResolution) {
  if (!capability) return []
  const maxDuration = capability.maxDurationByResolution?.[resolutionValue]
  return maxDuration ? capability.durations.filter((value) => value <= maxDuration) : capability.durations
}

function defaultAspectRatioFor(capability: VideoModelCapability, resolutionValue: VideoResolution): VideoAspectRatio {
  const supported = capability.aspectsByResolution[resolutionValue] || []
  return supported.includes(capability.defaultAspectRatio) ? capability.defaultAspectRatio : (supported[0] || capability.defaultAspectRatio)
}

function supportsModelParameters(modelValue: string, resolutionValue: VideoResolution, durationValue: number, aspectRatioValue: VideoAspectRatio) {
  const capability = videoModelCapabilities[modelValue]
  return Boolean(
    capability &&
    capability.resolutions.includes(resolutionValue) &&
    supportedDurations(capability, resolutionValue).includes(durationValue) &&
    capability.aspectsByResolution[resolutionValue]?.includes(aspectRatioValue)
  )
}

watch(model, () => {
  const capability = currentModelCapability.value
  if (!capability) return
  resolution.value = capability.defaultResolution
  duration.value = capability.defaultDuration
  aspectRatio.value = defaultAspectRatioFor(capability, capability.defaultResolution)
  promptEnhance.value = ''
}, { flush: 'sync' })
watch(resolution, () => {
  const capability = currentModelCapability.value
  if (!capability) return
  const supportedAspects = capability.aspectsByResolution[resolution.value] || []
  if (!supportedDurations(capability, resolution.value).includes(duration.value)) duration.value = capability.defaultDuration
  if (!supportedAspects.includes(aspectRatio.value)) aspectRatio.value = defaultAspectRatioFor(capability, resolution.value)
}, { flush: 'sync' })

watch(selectedKeyId, () => {
  if (apiKeyMode.value !== 'saved') return
  resetKeyScopedState()
  if (effectiveApiKey.value) void loadJobs()
})
watch(apiKeyMode, () => {
  resetKeyScopedState()
  if (apiKeyMode.value === 'saved' && effectiveApiKey.value) void loadJobs()
})
watch(customApiKey, () => {
  if (apiKeyMode.value === 'custom') resetKeyScopedState()
})
watch(activeJobs, updatePolling)
watch([effectiveApiKey, () => selectedJob.value?.job_id, () => selectedJob.value?.status], () => void loadSelectedVideo())

onMounted(async () => {
  await loadKeys()
})

onBeforeUnmount(() => {
  videoOutputRequest++
  stopPolling()
  clearSelectedVideo()
  clearImageInputs()
  clearMediaInputs()
})
</script>
