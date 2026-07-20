<template>
  <figure class="min-w-0 overflow-hidden rounded-lg border border-gray-800 bg-gray-950 dark:border-dark-600">
    <figcaption class="flex h-10 items-center justify-between border-b border-gray-800 px-3 text-xs text-gray-400">
      <span class="truncate font-medium">{{ label }}</span>
      <button
        type="button"
        class="inline-flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-md text-gray-400 transition-colors hover:bg-gray-800 hover:text-white focus:outline-none focus-visible:ring-2 focus-visible:ring-primary-400"
        :title="copied ? t('common.copied') : t('common.copy')"
        :aria-label="copied ? t('common.copied') : t('common.copy')"
        @click="copyCode"
      >
        <Icon :name="copied ? 'check' : 'copy'" size="sm" />
      </button>
    </figcaption>
    <pre class="max-h-[440px] overflow-auto p-4 text-xs leading-6 text-gray-100"><code>{{ code }}</code></pre>
  </figure>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'

const props = defineProps<{
  label: string
  code: string
}>()

const { t } = useI18n()
const copied = ref(false)
let resetTimer: ReturnType<typeof setTimeout> | undefined

async function copyCode() {
  await navigator.clipboard.writeText(props.code)
  copied.value = true
  if (resetTimer) clearTimeout(resetTimer)
  resetTimer = setTimeout(() => { copied.value = false }, 1600)
}
</script>
